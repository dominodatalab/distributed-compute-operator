package mpi

import (
	"fmt"
	"path/filepath"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

var (
	workerCommand = []string{
		"/opt/domino/bin/mpi-worker-start.sh",
	}
	sidecarCommand = []string{
		"/opt/domino/bin/rsync-start.sh",
	}
	customizerCommand = []string{
		"tar", "-C", "/", "-xf", "/root/worker-utils.tgz",
	}
)

// Key of the shared Secret object that contains client-side SSH public key
const publicKeyField = "ssh-publickey"

func StatefulSet(initImage, syncImage string) core.OwnedComponent {
	return &statefulSetComponent{
		InitImage: initImage,
		SyncImage: syncImage,
	}
}

type statefulSetComponent struct {
	InitImage string
	SyncImage string
}

func (c statefulSetComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	workerImage, err := util.ParseImageDefinition(cr.Spec.Image)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot parse workerImage: %w", err)
	}

	err = assureSharedKey(ctx, cr)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("invalid shared key: %w", err)
	}

	worker := cr.Spec.Worker
	labels := meta.StandardLabelsWithComponent(cr, ComponentWorker, worker.Labels)
	serviceAccount := selectServiceAccount(cr)

	initVolumes, initMounts := initVolumes()
	secretVolumes, secretMounts := secretVolumes(cr)

	allVolumes := make([]corev1.Volume, 0)
	allVolumes = append(allVolumes, worker.Volumes...)
	allVolumes = append(allVolumes, secretVolumes...)
	allVolumes = append(allVolumes, initVolumes...)

	workerMounts := make([]corev1.VolumeMount, 0)
	workerMounts = append(workerMounts, worker.VolumeMounts...)
	workerMounts = append(workerMounts, secretMounts...)
	workerMounts = append(workerMounts, initMounts...)

	sidecarMounts := make([]corev1.VolumeMount, 0)
	sidecarMounts = append(sidecarMounts, worker.VolumeMounts...)
	sidecarMounts = append(sidecarMounts, secretMounts...) // TODO: to be removed

	initContainers := make([]corev1.Container, 0)
	initContainers = append(initContainers, worker.InitContainers...)
	initContainers = append(initContainers, createInitContainer(cr, c.InitImage, initMounts))

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workerStatefulSetName(cr),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    worker.Replicas,
			ServiceName: serviceName(cr, ComponentWorker),
			Selector: &metav1.LabelSelector{
				MatchLabels: meta.MatchLabelsWithComponent(cr, ComponentWorker),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: worker.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccount,
					NodeSelector:       worker.NodeSelector,
					Affinity:           worker.Affinity,
					Tolerations:        worker.Tolerations,
					InitContainers:     initContainers,
					ImagePullSecrets:   cr.Spec.ImagePullSecrets,
					SecurityContext:    cr.Spec.PodSecurityContext,
					Volumes:            allVolumes,
					Containers: []corev1.Container{
						createWorkerContainer(cr, workerImage, workerMounts),
						createSidecarContainer(cr, c.SyncImage, sidecarMounts),
					},
				},
			},
			VolumeClaimTemplates: persistentVolumeClaims(worker.VolumeClaimTemplates),
			PodManagementPolicy:  appsv1.ParallelPodManagement,
		},
	}

	err = actions.CreateOrUpdateOwnedResource(ctx, cr, sts)
	if err != nil {
		err = fmt.Errorf("cannot reconcile statefulset: %w", err)
	}

	return ctrl.Result{}, err
}

func (c statefulSetComponent) Finalize(ctx *core.Context) (ctrl.Result, bool, error) {
	cr := objToMPICluster(ctx.Object)

	pvcListOpts := []client.ListOption{
		client.InNamespace(cr.Namespace),
		client.MatchingLabels(meta.MatchLabels(cr)),
	}
	err := actions.DeleteStorage(ctx, pvcListOpts)
	if err != nil {
		return ctrl.Result{RequeueAfter: finalizerRetryPeriod}, false,
			fmt.Errorf("cannot delete storage: %w", err)
	}

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workerStatefulSetName(cr),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabelsWithComponent(cr, ComponentWorker, cr.Spec.Worker.Labels),
		},
	}
	err = actions.DeleteIfExists(ctx, sts)
	if err != nil {
		return ctrl.Result{RequeueAfter: finalizerRetryPeriod}, false,
			fmt.Errorf("cannot delete workers: %w", err)
	}

	return ctrl.Result{}, true, nil
}

func (c statefulSetComponent) Kind() client.Object {
	return &appsv1.StatefulSet{}
}

func initVolumes() ([]corev1.Volume, []corev1.VolumeMount) {
	const customUtilVolume = "custom-util-volume"
	volumes := []corev1.Volume{
		{
			Name: customUtilVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
	mounts := []corev1.VolumeMount{
		{
			Name:      customUtilVolume,
			MountPath: customUtilPath,
			SubPath:   filepath.Base(customUtilPath),
		},
	}
	return volumes, mounts
}

func secretVolumes(cr *dcv1alpha1.MPICluster) ([]corev1.Volume, []corev1.VolumeMount) {
	const authorizedKeysVolume = "authorized-keys-volume"
	authorizedKeysModeCopy := int32(authorizedKeysMode)
	authorizedKeysName := filepath.Base(authorizedKeysPath)
	volumes := []corev1.Volume{
		{
			Name: authorizedKeysVolume,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: sshSecretName(cr),
					Items: []corev1.KeyToPath{
						{
							Key:  publicKeyField,
							Path: authorizedKeysName,
						},
					},
					DefaultMode: &authorizedKeysModeCopy,
				},
			},
		},
	}
	mounts := []corev1.VolumeMount{
		{
			Name:      authorizedKeysVolume,
			MountPath: authorizedKeysPath,
			SubPath:   authorizedKeysName,
		},
	}
	return volumes, mounts
}

func persistentVolumeClaims(vcts []dcv1alpha1.PersistentVolumeClaimTemplate) (pvcs []corev1.PersistentVolumeClaim) {
	mode := corev1.PersistentVolumeFilesystem

	for _, tmpl := range vcts {
		spec := tmpl.Spec.DeepCopy()
		spec.VolumeMode = &mode

		pvcs = append(pvcs, corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: tmpl.Name,
			},
			Spec: tmpl.Spec,
		})
	}

	return
}

func assureSharedKey(ctx *core.Context, cr *dcv1alpha1.MPICluster) error {
	secretName := sshSecretName(cr)
	objKey := client.ObjectKey{
		Name:      secretName,
		Namespace: cr.Namespace,
	}
	var sec corev1.Secret
	err := ctx.Client.Get(ctx, objKey, &sec)
	if err != nil {
		return fmt.Errorf("shared secret '%s' not found", secretName)
	}
	_, hasKey := sec.Data[publicKeyField]
	if !hasKey {
		return fmt.Errorf("shared secret '%s' doesn't contain a key", secretName)
	}
	return nil
}

func workerEnvironmentExtras(cr *dcv1alpha1.MPICluster) []corev1.EnvVar {
	userID := int64(defaultUserID)
	if cr.Spec.Worker.UserID != nil {
		userID = *cr.Spec.Worker.UserID
	}
	userName := defaultUserName
	if cr.Spec.Worker.UserName != "" {
		userName = cr.Spec.Worker.UserName
	}
	groupID := int64(defaultGroupID)
	if cr.Spec.Worker.GroupID != nil {
		groupID = *cr.Spec.Worker.GroupID
	}
	groupName := defaultGroupName
	if cr.Spec.Worker.GroupName != "" {
		groupName = cr.Spec.Worker.GroupName
	}
	return []corev1.EnvVar{
		{
			Name:  "DOMINO_SSH_PORT",
			Value: strconv.FormatInt(sshdPort, 10),
		},
		{
			Name:  "DOMINO_UID",
			Value: strconv.FormatInt(userID, 10),
		},
		{
			Name:  "DOMINO_GID",
			Value: strconv.FormatInt(groupID, 10),
		},
		{
			Name:  "DOMINO_USER",
			Value: userName,
		},
		{
			Name:  "DOMINO_GROUP",
			Value: groupName,
		},
		{
			Name:  "DOMINO_KEYS_PATH",
			Value: authorizedKeysPath,
		},
	}
}

func createWorkerContainer(cr *dcv1alpha1.MPICluster, image string, mounts []corev1.VolumeMount) corev1.Container {
	environment := make([]corev1.EnvVar, 0)
	environment = append(environment, cr.Spec.EnvVars...)
	environment = append(environment, workerEnvironmentExtras(cr)...)

	probe := &corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(sshdPort),
			},
		},
	}

	return corev1.Container{
		Name:            ApplicationName,
		Command:         workerCommand,
		Image:           image,
		ImagePullPolicy: cr.Spec.Image.PullPolicy,
		Ports: []corev1.ContainerPort{
			{
				Name:          sshdPortName,
				ContainerPort: sshdPort,
			},
		},
		Env:            environment,
		VolumeMounts:   mounts,
		Resources:      cr.Spec.Worker.Resources,
		LivenessProbe:  probe,
		ReadinessProbe: probe,
	}
}

func createSidecarContainer(cr *dcv1alpha1.MPICluster, image string, mounts []corev1.VolumeMount) corev1.Container {
	if image == "" {
		image = defaultSyncImage
	}

	probe := &corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(rsyncPort),
			},
		},
	}

	user := int64(rsyncUserID)
	group := int64(rsyncGroupID)

	return corev1.Container{
		Name:            RsyncSidecarName,
		Command:         sidecarCommand,
		Image:           image,
		ImagePullPolicy: cr.Spec.Image.PullPolicy,
		Ports: []corev1.ContainerPort{
			{
				Name:          rsyncPortName,
				ContainerPort: rsyncPort,
			},
		},
		Env: []corev1.EnvVar{
			{
				Name:  "RSYNC_PORT",
				Value: strconv.FormatInt(rsyncPort, 10),
			},
		},
		VolumeMounts: mounts,
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:  &user,
			RunAsGroup: &group,
		},
		LivenessProbe:  probe,
		ReadinessProbe: probe,
	}
}

func createInitContainer(cr *dcv1alpha1.MPICluster, image string, mounts []corev1.VolumeMount) corev1.Container {
	if image == "" {
		image = defaultInitImage
	}

	return corev1.Container{
		Name:            ApplicationName + "-init",
		Command:         customizerCommand,
		Image:           image,
		ImagePullPolicy: cr.Spec.Image.PullPolicy,
		VolumeMounts:    mounts,
	}
}

package mpi

import (
	"fmt"
	"path/filepath"

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
	workerPorts = []corev1.ContainerPort{
		{
			Name:          sshdPortName,
			ContainerPort: sshdPort,
		},
	}
	workerProbe = &corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(int(sshdPort)),
			},
		},
	}
	workerCommand = []string{
		"/usr/sbin/sshd",
		"-De",
	}
)

func StatefulSet() core.OwnedComponent {
	return &statefulSetComponent{}
}

type statefulSetComponent struct{}

func (c statefulSetComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	image, err := util.ParseImageDefinition(cr.Spec.Image)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot parse image: %w", err)
	}

	err = assureSharedKey(ctx, cr)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("invalid shared key: %w", err)
	}

	worker := cr.Spec.Worker
	labels := meta.StandardLabelsWithComponent(cr, ComponentWorker, worker.Labels)
	serviceAccount := selectServiceAccount(cr)
	volumes, volumeMounts := buildWorkerVolumesAndMounts(cr)

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workerStatefulSetName(cr),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    worker.Replicas,
			ServiceName: serviceName(cr),
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
					InitContainers:     worker.InitContainers,
					ImagePullSecrets:   cr.Spec.ImagePullSecrets,
					SecurityContext:    cr.Spec.PodSecurityContext,
					Volumes:            volumes,
					Containers: []corev1.Container{
						{
							Name:            ApplicationName,
							Command:         workerCommand,
							Image:           image,
							ImagePullPolicy: cr.Spec.Image.PullPolicy,
							Ports:           workerPorts,
							Env:             cr.Spec.EnvVars,
							VolumeMounts:    volumeMounts,
							Resources:       worker.Resources,
							LivenessProbe:   workerProbe,
							ReadinessProbe:  workerProbe,
						},
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
		return ctrl.Result{}, false, fmt.Errorf("cannot delete storage: %w", err)
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
		return ctrl.Result{}, false, fmt.Errorf("cannot delete workers: %w", err)
	}

	return ctrl.Result{}, true, nil
}

func (c statefulSetComponent) Kind() client.Object {
	return &appsv1.StatefulSet{}
}

func buildWorkerVolumesAndMounts(cr *dcv1alpha1.MPICluster) ([]corev1.Volume, []corev1.VolumeMount) {
	volumes := []corev1.Volume{
		{
			Name: sshVolumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: sshSecretName(cr),
					Items: []corev1.KeyToPath{
						{
							Key:  sshAuthPublicKey,
							Path: sshAuthorizedKeysFilename,
						},
					},
				},
			},
		},
		{
			Name: configVolumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMapName(cr),
					},
					Items: []corev1.KeyToPath{
						{
							Key:  sshdConfigFilename,
							Path: sshdConfigFilename,
						},
					},
				},
			},
		},
	}
	volumes = append(volumes, cr.Spec.Worker.Volumes...)

	mounts := []corev1.VolumeMount{
		{
			Name:      sshVolumeName,
			MountPath: sshAuthorizedKeysPath,
			SubPath:   sshAuthorizedKeysFilename,
		},
		{
			Name:      configVolumeName,
			MountPath: filepath.Join(sshRootMountPath, sshdConfigFilename),
			SubPath:   sshdConfigFilename,
		},
	}
	mounts = append(mounts, cr.Spec.Worker.VolumeMounts...)

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
	_, hasKey := sec.Data[sshAuthPublicKey]
	if !hasKey {
		return fmt.Errorf("shared secret '%s' doesn't contain a key", secretName)
	}
	return nil
}

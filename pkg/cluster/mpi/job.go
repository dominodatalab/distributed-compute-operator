package mpi

import (
	"fmt"
	"path/filepath"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

const fileModeUserRO = 0400

var launcherEnv = []corev1.EnvVar{
	{
		Name:  "OMPI_MCA_orte_keep_fqdn_hostnames",
		Value: "true",
	},
	{
		Name:  "OMPI_MCA_orte_default_hostfile",
		Value: hostfilePath,
	},
}

func Job() core.OwnedComponent {
	return &jobComponent{}
}

type jobComponent struct{}

func (c jobComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPIJob(ctx.Object)

	ready, err := workersAvailable(ctx, cr)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !ready {
		ctx.Log.Info("Delaying job creation until all workers ready")
		return ctrl.Result{RequeueAfter: 500 * time.Millisecond}, nil
	}

	job, err := buildJob(cr)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = actions.CreateOrUpdateOwnedResource(ctx, cr, job)
	if err != nil {
		err = fmt.Errorf("cannot reconcile job: %w", err)
	}

	return ctrl.Result{}, err
}

func (c jobComponent) Kind() client.Object {
	return &batchv1.Job{}
}

func buildJob(cr *dcv1alpha1.MPIJob) (*batchv1.Job, error) {
	image, err := util.ParseImageDefinition(cr.Spec.Image)
	if err != nil {
		return nil, fmt.Errorf("cannot parse image: %w", err)
	}

	launcher := cr.Spec.Launcher
	env := buildEnv(cr)
	labels := meta.StandardLabelsWithComponent(cr, ComponentLauncher, launcher.Labels)
	serviceAccount := selectServiceAccount(cr)
	volumes, volumeMounts := buildLauncherVolumesAndMounts(cr)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName(cr),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: launcher.Annotations,
				},
				Spec: corev1.PodSpec{
					RestartPolicy:      corev1.RestartPolicyOnFailure,
					ServiceAccountName: serviceAccount,
					NodeSelector:       launcher.NodeSelector,
					Affinity:           launcher.Affinity,
					Tolerations:        launcher.Tolerations,
					InitContainers:     launcher.InitContainers,
					ImagePullSecrets:   cr.Spec.ImagePullSecrets,
					SecurityContext:    cr.Spec.PodSecurityContext,
					Volumes:            volumes,
					Containers: []corev1.Container{
						{
							Name:            ApplicationName,
							Command:         cr.Spec.Launcher.Command,
							Image:           image,
							ImagePullPolicy: cr.Spec.Image.PullPolicy,
							Env:             env,
							VolumeMounts:    volumeMounts,
							Resources:       launcher.Resources,
						},
					},
				},
			},
		},
	}, nil
}

func buildEnv(cr *dcv1alpha1.MPIJob) []corev1.EnvVar {
	return append(launcherEnv, cr.Spec.EnvVars...)
}

func buildLauncherVolumesAndMounts(cr *dcv1alpha1.MPIJob) ([]corev1.Volume, []corev1.VolumeMount) {
	volumes := []corev1.Volume{
		{
			Name: sshVolumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: pointer.Int32(fileModeUserRO),
					SecretName:  sshSecretName(cr),
					Items: []corev1.KeyToPath{
						{
							Key:  corev1.SSHAuthPrivateKey,
							Path: sshPrivateKeyFilename,
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
							Key:  hostfileFilename,
							Path: hostfileFilename,
						},
						{
							Key:  sshConfigFilename,
							Path: sshConfigFilename,
						},
					},
				},
			},
		},
	}
	volumes = append(volumes, cr.Spec.Launcher.Volumes...)

	mounts := []corev1.VolumeMount{
		{
			Name:      sshVolumeName,
			MountPath: sshPrivateKeyPath,
			SubPath:   sshPrivateKeyFilename,
		},
		{
			Name:      configVolumeName,
			MountPath: hostfilePath,
			SubPath:   hostfileFilename,
		},
		{
			Name:      configVolumeName,
			MountPath: filepath.Join(sshRootMountPath, sshConfigFilename),
			SubPath:   sshConfigFilename,
		},
	}
	mounts = append(mounts, cr.Spec.Launcher.VolumeMounts...)

	return volumes, mounts
}

func workersAvailable(ctx *core.Context, cr *dcv1alpha1.MPIJob) (bool, error) {
	var ep corev1.Endpoints
	objKey := client.ObjectKey{
		Name:      serviceName(cr),
		Namespace: cr.Namespace,
	}

	if err := ctx.Client.Get(ctx, objKey, &ep); err != nil && !apierrors.IsNotFound(err) {
		return false, fmt.Errorf("cannot fetch endpoints: %w", err)
	}

	actual := 0
	expected := int(*cr.Spec.Worker.Replicas)

	if len(ep.Subsets) == 1 {
		actual = len(ep.Subsets[0].Addresses)

		if actual == expected {
			ctx.Log.Info("All workers ready", "available", actual, "desired", expected)
			return true, nil
		}
	}

	ctx.Log.Info("Insufficient workers", "available", actual, "desired", expected)
	return false, nil
}

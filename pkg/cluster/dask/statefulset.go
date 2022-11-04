package dask

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

func StatefulSetScheduler() core.OwnedComponent {
	return components.StatefulSet(func(obj client.Object) components.StatefulSetDataSource {
		dc := daskCluster(obj)
		tc := &schedulerConfig{dc: dc}

		return &statefulSetDS{tc, dc, ComponentScheduler}
	})
}

func StatefulSetWorker() core.OwnedComponent {
	return components.StatefulSet(func(obj client.Object) components.StatefulSetDataSource {
		dc := daskCluster(obj)
		tc := &workerConfig{dc: dc}

		return &statefulSetDS{tc, dc, ComponentWorker}
	})
}

type statefulSetDS struct {
	tc   typeConfig
	dc   *dcv1alpha1.DaskCluster
	comp metadata.Component
}

func (s *statefulSetDS) StatefulSet() (*appsv1.StatefulSet, error) {
	imageDef, err := util.ParseImageDefinition(s.image())
	if err != nil {
		return nil, fmt.Errorf("cannot parse image: %w", err)
	}

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.name(),
			Namespace: s.namespace(),
			Labels:    s.labels(),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    s.replicas(),
			ServiceName: s.serviceName(),
			Selector: &metav1.LabelSelector{
				MatchLabels: s.matchLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      s.labels(),
					Annotations: s.podAnnotations(),
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: s.serviceAccountName(),
					NodeSelector:       s.nodeSelector(),
					Affinity:           s.affinity(),
					Tolerations:        s.tolerations(),
					InitContainers:     s.initContainers(),
					ImagePullSecrets:   s.imagePullSecrets(),
					SecurityContext:    s.podSecurityContext(),
					Volumes:            s.volumes(),
					Containers: []corev1.Container{
						{
							Name:            s.applicationName(),
							Args:            s.commandArgs(),
							Image:           imageDef,
							ImagePullPolicy: s.image().PullPolicy,
							Ports:           s.ports(),
							Env:             s.env(),
							VolumeMounts:    s.volumeMounts(),
							Resources:       s.resources(),
							LivenessProbe:   s.probe(),
							ReadinessProbe:  s.probe(),
							SecurityContext: s.securityContext(),
						},
					},
				},
			},
			VolumeClaimTemplates: s.pvcTemplates(),
			PodManagementPolicy:  appsv1.ParallelPodManagement,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
		},
	}

	return sts, nil
}

func (s *statefulSetDS) PVCListOpts() []client.ListOption {
	return []client.ListOption{
		client.InNamespace(s.dc.Namespace),
		client.MatchingLabels(meta.MatchLabels(s.dc)),
	}
}

func (s *statefulSetDS) applicationName() string {
	return ApplicationName
}

func (s *statefulSetDS) name() string {
	return meta.InstanceName(s.dc, s.comp)
}

func (s *statefulSetDS) namespace() string {
	return s.dc.Namespace
}

func (s *statefulSetDS) labels() map[string]string {
	return meta.StandardLabelsWithComponent(s.dc, s.comp, s.tc.podConfig().Labels)
}

func (s *statefulSetDS) securityContext() *corev1.SecurityContext {
	switch s.comp {
	case ComponentScheduler:
		return s.dc.Spec.Scheduler.SecurityContext
	case ComponentWorker:
		return s.dc.Spec.Worker.SecurityContext
	default:
		return nil
	}
}

func (s *statefulSetDS) matchLabels() map[string]string {
	return meta.MatchLabelsWithComponent(s.dc, s.comp)
}

func (s *statefulSetDS) serviceName() string {
	return meta.InstanceName(s.dc, s.comp)
}

func (s *statefulSetDS) serviceAccountName() string {
	if s.dc.Spec.ServiceAccount.Name != "" {
		return s.dc.Spec.ServiceAccount.Name
	}

	return meta.InstanceName(s.dc, metadata.ComponentNone)
}

func (s *statefulSetDS) image() *dcv1alpha1.OCIImageDefinition {
	return s.dc.Spec.Image
}

func (s *statefulSetDS) imagePullSecrets() []corev1.LocalObjectReference {
	return s.dc.Spec.ImagePullSecrets
}

func (s *statefulSetDS) podSecurityContext() *corev1.PodSecurityContext {
	return s.dc.Spec.PodSecurityContext
}

func (s *statefulSetDS) env() []corev1.EnvVar {
	envvars := s.dc.Spec.EnvVars
	envvars = append(envvars, s.tc.containerEnv()...)

	return envvars
}

func (s *statefulSetDS) replicas() *int32 {
	return s.tc.replicas()
}

func (s *statefulSetDS) commandArgs() []string {
	return s.tc.commandArgs()
}

func (s *statefulSetDS) ports() []corev1.ContainerPort {
	return s.tc.containerPorts()
}

func (s *statefulSetDS) podAnnotations() map[string]string {
	return s.tc.podConfig().Annotations
}

func (s *statefulSetDS) nodeSelector() map[string]string {
	return s.tc.podConfig().NodeSelector
}

func (s *statefulSetDS) affinity() *corev1.Affinity {
	return s.tc.podConfig().Affinity
}

func (s *statefulSetDS) tolerations() []corev1.Toleration {
	return s.tc.podConfig().Tolerations
}

func (s *statefulSetDS) initContainers() []corev1.Container {
	return s.tc.podConfig().InitContainers
}

func (s *statefulSetDS) volumes() []corev1.Volume {
	volumes := s.tc.podConfig().Volumes
	if s.dc.Spec.KerberosKeytab != nil {
		volumes = append(volumes, corev1.Volume{
			Name: "kerberos",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: meta.InstanceName(s.dc, metadata.ComponentNone),
					},
				},
			},
		})
	}

	return volumes
}

func (s *statefulSetDS) volumeMounts() []corev1.VolumeMount {
	mounts := s.tc.podConfig().VolumeMounts
	if s.dc.Spec.KerberosKeytab != nil {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "kerberos",
			ReadOnly:  true,
			MountPath: s.dc.Spec.KerberosKeytab.MountPath,
		})
	}

	return mounts
}

func (s *statefulSetDS) resources() corev1.ResourceRequirements {
	return s.tc.podConfig().Resources
}

func (s *statefulSetDS) probe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.FromString("dashboard"),
			},
		},
	}
}

func (s *statefulSetDS) pvcTemplates() (tmpls []corev1.PersistentVolumeClaim) {
	mode := corev1.PersistentVolumeFilesystem

	for _, vct := range s.tc.podConfig().VolumeClaimTemplates {
		spec := vct.Spec.DeepCopy()
		spec.VolumeMode = &mode

		tmpls = append(tmpls, corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: vct.Name,
			},
			Spec: vct.Spec,
		})
	}

	return
}

type typeConfig interface {
	podConfig() dcv1alpha1.WorkloadConfig
	replicas() *int32
	commandArgs() []string
	containerEnv() []corev1.EnvVar
	containerPorts() []corev1.ContainerPort
}

type schedulerConfig struct {
	dc *dcv1alpha1.DaskCluster
}

func (c *schedulerConfig) podConfig() dcv1alpha1.WorkloadConfig {
	return c.dc.Spec.Scheduler
}

func (c *schedulerConfig) replicas() *int32 {
	return pointer.Int32Ptr(1)
}

func (c *schedulerConfig) commandArgs() []string {
	return []string{
		"dask-scheduler",
		fmt.Sprintf("--port=%d", c.dc.Spec.SchedulerPort),
		fmt.Sprintf("--dashboard-address=:%d", c.dc.Spec.DashboardPort),
	}
}

func (c *schedulerConfig) containerEnv() []corev1.EnvVar {
	return nil
}

func (c *schedulerConfig) containerPorts() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "serve",
			ContainerPort: c.dc.Spec.SchedulerPort,
		},
		{
			Name:          "dashboard",
			ContainerPort: c.dc.Spec.DashboardPort,
		},
	}
}

type workerConfig struct {
	dc *dcv1alpha1.DaskCluster
}

func (c *workerConfig) podConfig() dcv1alpha1.WorkloadConfig {
	return c.dc.Spec.Worker.WorkloadConfig
}

func (c *workerConfig) replicas() *int32 {
	return c.dc.Spec.Worker.Replicas
}

func (c *workerConfig) commandArgs() []string {
	return []string{
		"dask-worker",
		"--name=$(MY_POD_NAME)",
		"--local-directory=/tmp",
		// NOTE: it looks like the dask worker can infer its threads/memory from resource.limits
		// "--nthreads=$(MY_CPU_LIMIT)",
		// "--memory=$(MY_MEM_LIMIT)",
		fmt.Sprintf("--worker-port=%d", c.dc.Spec.WorkerPort),
		fmt.Sprintf("--nanny-port=%d", c.dc.Spec.NannyPort),
		fmt.Sprintf("--dashboard-address=:%d", c.dc.Spec.DashboardPort),
		fmt.Sprintf("%s:%d", meta.InstanceName(c.dc, ComponentScheduler), c.dc.Spec.SchedulerPort),
	}
}

func (c *workerConfig) containerEnv() []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: "MY_POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
	}
}

func (c *workerConfig) containerPorts() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "worker",
			ContainerPort: c.dc.Spec.WorkerPort,
		},
		{
			Name:          "nanny",
			ContainerPort: c.dc.Spec.NannyPort,
		},
		{
			Name:          "dashboard",
			ContainerPort: c.dc.Spec.DashboardPort,
		},
	}
}

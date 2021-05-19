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
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

type statefulSetDS struct {
	tc   typeConfig
	dc   *dcv1alpha1.DaskCluster
	comp metadata.Component
}

func SchedulerStatefulSet(obj client.Object) components.StatefulSetDataSource {
	dc := obj.(*dcv1alpha1.DaskCluster)
	tc := &schedulerConfig{dc: dc}

	return &statefulSetDS{tc, dc, ComponentScheduler}
}

func WorkerStatefulSet(obj client.Object) components.StatefulSetDataSource {
	dc := obj.(*dcv1alpha1.DaskCluster)
	tc := &workerConfig{dc: dc}

	return &statefulSetDS{tc, dc, ComponentWorker}
}

func (s *statefulSetDS) GetStatefulSet() (*appsv1.StatefulSet, error) {
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
					SecurityContext:    s.securityContext(),
					Volumes:            s.volumes(),
					Containers: []corev1.Container{
						{
							Name:            s.applicationName(),
							Command:         s.command(),
							Args:            s.commandArgs(),
							Image:           imageDef,
							ImagePullPolicy: s.image().PullPolicy,
							Ports:           s.ports(),
							Env:             s.env(),
							VolumeMounts:    s.volumeMounts(),
							Resources:       s.resources(),
							LivenessProbe:   s.probe(),
							ReadinessProbe:  s.probe(),
						},
					},
				},
			},
			PodManagementPolicy: appsv1.ParallelPodManagement,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
		},
	}

	return sts, nil
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
	labels := meta.StandardLabelsWithComponent(s.dc, s.comp)
	return util.MergeStringMaps(s.tc.PodConfig().Labels, labels)
}

func (s *statefulSetDS) matchLabels() map[string]string {
	return meta.MatchLabelsWithComponent(s.dc, s.comp)
}

func (s *statefulSetDS) serviceName() string {
	return meta.InstanceName(s.dc, s.comp)
}

func (s *statefulSetDS) serviceAccountName() string {
	return meta.InstanceName(s.dc, metadata.ComponentNone)
}

func (s *statefulSetDS) image() *dcv1alpha1.OCIImageDefinition {
	return s.dc.Spec.Image
}

func (s *statefulSetDS) imagePullSecrets() []corev1.LocalObjectReference {
	return s.dc.Spec.ImagePullSecrets
}

func (s *statefulSetDS) securityContext() *corev1.PodSecurityContext {
	return s.dc.Spec.PodSecurityContext
}

func (s *statefulSetDS) env() []corev1.EnvVar {
	return s.dc.Spec.EnvVars
}

func (s *statefulSetDS) replicas() *int32 {
	return pointer.Int32Ptr(s.tc.replicas())
}

func (s *statefulSetDS) podAnnotations() map[string]string {
	return s.tc.PodConfig().Annotations
}

func (s *statefulSetDS) command() []string {
	return s.tc.command()
}

func (s *statefulSetDS) commandArgs() []string {
	return s.tc.commandArgs()
}

func (s *statefulSetDS) nodeSelector() map[string]string {
	return s.tc.PodConfig().NodeSelector
}

func (s *statefulSetDS) affinity() *corev1.Affinity {
	return s.tc.PodConfig().Affinity
}

func (s *statefulSetDS) tolerations() []corev1.Toleration {
	return s.tc.PodConfig().Tolerations
}

func (s *statefulSetDS) initContainers() []corev1.Container {
	return s.tc.PodConfig().InitContainers
}

func (s *statefulSetDS) ports() []corev1.ContainerPort {
	return s.tc.containerPorts()
}

func (s *statefulSetDS) volumes() []corev1.Volume {
	return s.tc.PodConfig().Volumes
}

func (s *statefulSetDS) volumeMounts() []corev1.VolumeMount {
	return s.tc.PodConfig().VolumeMounts
}

func (s *statefulSetDS) resources() corev1.ResourceRequirements {
	return s.tc.PodConfig().Resources
}

func (s *statefulSetDS) probe() *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.FromString("dashboard"),
			},
		},
	}
}

type typeConfig interface {
	replicas() int32
	command() []string
	commandArgs() []string
	containerPorts() []corev1.ContainerPort
	PodConfig() dcv1alpha1.WorkloadConfig
}

type schedulerConfig struct {
	dc *dcv1alpha1.DaskCluster
}

func (c *schedulerConfig) replicas() int32 {
	return 1
}

func (c *schedulerConfig) PodConfig() dcv1alpha1.WorkloadConfig {
	return c.dc.Spec.Scheduler
}

func (c *schedulerConfig) command() []string {
	return []string{"dask-scheduler"}
}

func (c *schedulerConfig) commandArgs() []string {
	return []string{
		fmt.Sprintf("--port=%d", c.dc.Spec.SchedulerPort),
		fmt.Sprintf("--dashboard-address=:%d", c.dc.Spec.DashboardPort),
	}
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

func (c *workerConfig) replicas() int32 {
	return c.dc.Spec.Worker.Replicas
}

func (c *workerConfig) PodConfig() dcv1alpha1.WorkloadConfig {
	return c.dc.Spec.Worker.WorkloadConfig
}

func (c *workerConfig) command() []string {
	return []string{"dask-worker"}
}

func (c *workerConfig) commandArgs() []string {
	return []string{
		fmt.Sprintf("--worker-port=%d", c.dc.Spec.WorkerPort),
		fmt.Sprintf("--nanny-port=%d", c.dc.Spec.NannyPort),
		fmt.Sprintf("--dashboard-address=:%d", c.dc.Spec.DashboardPort),
		fmt.Sprintf("example-dask-scheduler:%d", c.dc.Spec.SchedulerPort),
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

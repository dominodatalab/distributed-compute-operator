package ray

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

var (
	defaultCmd = []string{ApplicationName}
	defaultEnv = []corev1.EnvVar{
		{
			Name: "MY_POD_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
		{
			Name: "MY_CPU_REQUEST",
			ValueFrom: &corev1.EnvVarSource{
				ResourceFieldRef: &corev1.ResourceFieldSelector{
					Resource: "requests.cpu",
				},
			},
		},
	}
	defaultVolumes = []corev1.Volume{
		{
			Name: sharedMemoryVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: corev1.StorageMediumMemory,
				},
			},
		},
	}
	defaultVolumeMounts = []corev1.VolumeMount{
		{
			Name:      sharedMemoryVolumeName,
			MountPath: "/dev/shm",
		},
	}
)

const sharedMemoryVolumeName = "dshm"

func NewStatefulSet(rc *dcv1alpha1.RayCluster, comp Component) (*appsv1.StatefulSet, error) {
	p, err := newConfigProcessor(rc, comp)
	if err != nil {
		return nil, err
	}

	imageRef, err := util.ParseImageDefinition(rc.Spec.Image)
	if err != nil {
		return nil, err
	}

	serviceAccountName := InstanceObjectName(rc.Name, ComponentNone)
	if rc.Spec.ServiceAccountName != "" {
		serviceAccountName = rc.Spec.ServiceAccountName
	}

	replicas := p.replicas()
	nodeAttrs := p.nodeAttributes()
	args := p.processArgs()
	ports := p.processPorts()
	labels := p.processLabels()
	serviceName := p.processServiceName()
	envVars := append(defaultEnv, rc.Spec.EnvVars...)
	volumes := append(defaultVolumes, nodeAttrs.Volumes...)
	volumeMounts := append(defaultVolumeMounts, nodeAttrs.VolumeMounts...)
	pvcTemplates := processPVCTemplates(nodeAttrs.VolumeClaimTemplates)

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, comp),
			Namespace: rc.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: serviceName,
			Replicas:    pointer.Int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: SelectorLabelsWithComponent(rc, comp),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: nodeAttrs.Annotations,
				},

				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					NodeSelector:       nodeAttrs.NodeSelector,
					Affinity:           nodeAttrs.Affinity,
					Tolerations:        nodeAttrs.Tolerations,
					InitContainers:     nodeAttrs.InitContainers,
					ImagePullSecrets:   rc.Spec.ImagePullSecrets,
					SecurityContext:    rc.Spec.PodSecurityContext,
					Containers: []corev1.Container{
						{
							Name:            ApplicationName,
							Command:         defaultCmd,
							Args:            args,
							Image:           imageRef,
							ImagePullPolicy: rc.Spec.Image.PullPolicy,
							Ports:           ports,
							Env:             envVars,
							VolumeMounts:    volumeMounts,
							Resources:       nodeAttrs.Resources,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(int(rc.Spec.NodeManagerPort)),
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(int(rc.Spec.NodeManagerPort)),
									},
								},
							},
						},
					},
					Volumes: volumes,
				},
			},
			VolumeClaimTemplates: pvcTemplates,
			PodManagementPolicy:  appsv1.ParallelPodManagement,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
		},
	}

	return sts, nil
}

type configProcessor interface {
	replicas() int32
	nodeAttributes() *dcv1alpha1.RayClusterNode
	processArgs() []string
	processPorts() []corev1.ContainerPort
	processLabels() map[string]string
	processServiceName() string
}

func newConfigProcessor(rc *dcv1alpha1.RayCluster, comp Component) (configProcessor, error) {
	switch comp {
	case ComponentHead:
		return &headProcessor{rc: rc}, nil
	case ComponentWorker:
		return &workerProcessor{rc: rc}, nil
	default:
		return nil, fmt.Errorf("invalid ray component: %q", comp)
	}
}

type headProcessor struct {
	rc *dcv1alpha1.RayCluster
}

func (p *headProcessor) replicas() int32 {
	return 1
}

func (p *headProcessor) nodeAttributes() *dcv1alpha1.RayClusterNode {
	return &p.rc.Spec.Head.RayClusterNode
}

func (p *headProcessor) processArgs() []string {
	rc := p.rc
	headArgs := []string{
		"--head",
		fmt.Sprintf("--ray-client-server-port=%d", rc.Spec.ClientServerPort),
		fmt.Sprintf("--port=%d", rc.Spec.Port),
		fmt.Sprintf("--redis-shard-ports=%s", strings.Join(util.IntsToStrings(rc.Spec.RedisShardPorts), ",")),
		fmt.Sprintf("--gcs-server-port=%d", rc.Spec.GCSServerPort),
	}

	if util.BoolPtrIsTrue(rc.Spec.EnableDashboard) {
		dashArgs := []string{
			"--include-dashboard=true",
			"--dashboard-host=0.0.0.0",
			fmt.Sprintf("--dashboard-port=%d", rc.Spec.DashboardPort),
		}
		headArgs = append(headArgs, dashArgs...)
	}

	return append(processArgs(rc), headArgs...)
}

func (p *headProcessor) processPorts() []corev1.ContainerPort {
	rc := p.rc

	ports := []corev1.ContainerPort{
		{
			Name:          "client",
			ContainerPort: rc.Spec.ClientServerPort,
		},
		{
			Name:          "object-manager",
			ContainerPort: rc.Spec.ObjectManagerPort,
		},
		{
			Name:          "node-manager",
			ContainerPort: rc.Spec.NodeManagerPort,
		},
		{
			Name:          "gcs-server",
			ContainerPort: rc.Spec.GCSServerPort,
		},
		{
			Name:          "redis-primary",
			ContainerPort: rc.Spec.Port,
		},
	}
	for idx, port := range rc.Spec.RedisShardPorts {
		ports = append(ports, corev1.ContainerPort{
			Name:          fmt.Sprintf("redis-shard-%d", idx),
			ContainerPort: port,
		})
	}
	for idx, port := range rc.Spec.WorkerPorts {
		ports = append(ports, corev1.ContainerPort{
			Name:          fmt.Sprintf("worker-%d", idx),
			ContainerPort: port,
		})
	}

	if util.BoolPtrIsTrue(rc.Spec.EnableDashboard) {
		ports = append(ports, corev1.ContainerPort{
			Name:          "dashboard",
			ContainerPort: rc.Spec.DashboardPort,
		})
	}

	return ports
}

func (p *headProcessor) processLabels() map[string]string {
	return processLabels(p.rc, ComponentHead, p.rc.Spec.Head.RayClusterNode.Labels)
}

func (p *headProcessor) processServiceName() string {
	return HeadlessHeadServiceName(p.rc.Name)
}

type workerProcessor struct {
	rc *dcv1alpha1.RayCluster
}

func (p *workerProcessor) replicas() int32 {
	return *p.rc.Spec.Worker.Replicas
}

func (p *workerProcessor) nodeAttributes() *dcv1alpha1.RayClusterNode {
	return &p.rc.Spec.Worker.RayClusterNode
}

func (p *workerProcessor) processArgs() []string {
	rc := p.rc
	headNodeAddr := fmt.Sprintf("%s-0", InstanceObjectName(rc.Name, ComponentHead))

	return append(
		processArgs(rc),
		fmt.Sprintf("--address=%s.%s:%d", headNodeAddr, HeadlessHeadServiceName(rc.Name), rc.Spec.Port),
	)
}

func (p *workerProcessor) processPorts() []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          "object-manager",
			ContainerPort: p.rc.Spec.ObjectManagerPort,
		},
		{
			Name:          "node-manager",
			ContainerPort: p.rc.Spec.NodeManagerPort,
		},
	}
	for idx, port := range p.rc.Spec.WorkerPorts {
		ports = append(ports, corev1.ContainerPort{
			Name:          fmt.Sprintf("worker-%d", idx),
			ContainerPort: port,
		})
	}

	return ports
}

func (p *workerProcessor) processLabels() map[string]string {
	return processLabels(p.rc, ComponentWorker, p.rc.Spec.Worker.RayClusterNode.Labels)
}

func (p *workerProcessor) processServiceName() string {
	return HeadlessWorkerServiceName(p.rc.Name)
}

// common head/worker command arguments
func processArgs(rc *dcv1alpha1.RayCluster) []string {
	args := []string{
		"start",
		"--block",
		"--node-ip-address=$(MY_POD_IP)",
		"--num-cpus=$(MY_CPU_REQUEST)",
		fmt.Sprintf("--object-manager-port=%d", rc.Spec.ObjectManagerPort),
		fmt.Sprintf("--node-manager-port=%d", rc.Spec.NodeManagerPort),
	}

	if rc.Spec.WorkerPorts != nil {
		args = append(args, fmt.Sprintf("--worker-port-list=%s", strings.Join(util.IntsToStrings(rc.Spec.WorkerPorts), ",")))
	}

	if rc.Spec.ObjectStoreMemoryBytes != nil {
		args = append(args, fmt.Sprintf("--object-store-memory=%d", *rc.Spec.ObjectStoreMemoryBytes))
	}

	return args
}

// common head/worker labels
func processLabels(rc *dcv1alpha1.RayCluster, comp Component, extraLabels map[string]string) map[string]string {
	labels := MetadataLabelsWithComponent(rc, comp)
	if extraLabels != nil {
		labels = util.MergeStringMaps(extraLabels, labels)
	}

	return labels
}

func processPVCTemplates(vcts []dcv1alpha1.PersistentVolumeClaimTemplate) (pvcTmpls []corev1.PersistentVolumeClaim) {
	mode := corev1.PersistentVolumeFilesystem

	for _, vct := range vcts {
		spec := vct.Spec.DeepCopy()
		spec.VolumeMode = &mode

		pvcTmpls = append(pvcTmpls, corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: vct.Name,
			},
			Spec: vct.Spec,
		})
	}

	return
}

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
			Name:  "GLOG_logtostderr",
			Value: "1",
		},
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

// NewDeployment generates a Deployment configured to manage Ray cluster nodes.
// The configuration is based the provided spec and the desired Component workload.
func NewDeployment(rc *dcv1alpha1.RayCluster, comp Component) (*appsv1.Deployment, error) {
	// TODO: need to figure out how the update strategy for head and worker pods
	// 	- do we need to "recreate" the head to ensure only 1 instance is running at a time?
	// 	- how the workers behave when the head goes down and comes back up?

	var replicaCount int32
	switch comp {
	case ComponentHead:
		replicaCount = 1
	case ComponentWorker:
		replicaCount = rc.Spec.WorkerReplicaCount
	default:
		return nil, fmt.Errorf("invalid ray component: %q", comp)
	}

	imageRef, err := util.ParseImageDefinition(rc.Spec.Image)
	if err != nil {
		return nil, err
	}

	args := processArgs(rc, comp)
	ports := processPorts(rc, comp)
	labels := processLabels(rc, comp)
	volumes := append(defaultVolumes, rc.Spec.Volumes...)
	volumeMounts := append(defaultVolumeMounts, rc.Spec.VolumeMounts...)

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, comp),
			Namespace: rc.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(replicaCount),
			Selector: &metav1.LabelSelector{
				MatchLabels: SelectorLabelsWithComponent(rc, comp),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: rc.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: rc.Name,
					NodeSelector:       rc.Spec.NodeSelector,
					Affinity:           rc.Spec.Affinity,
					Tolerations:        rc.Spec.Tolerations,
					InitContainers:     rc.Spec.InitContainers,
					Containers: []corev1.Container{
						{
							Name:            ApplicationName,
							Command:         defaultCmd,
							Args:            args,
							Image:           imageRef,
							ImagePullPolicy: rc.Spec.Image.PullPolicy,
							Ports:           ports,
							Env:             defaultEnv,
							VolumeMounts:    volumeMounts,
							Resources:       rc.Spec.Resources,
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
		},
	}

	return deploy, nil
}

func processArgs(rc *dcv1alpha1.RayCluster, comp Component) []string {
	cmdArgs := []string{
		"start",
		"--block",
		"--node-ip-address=$(MY_POD_IP)",
		"--num-cpus=$(MY_CPU_REQUEST)",
		fmt.Sprintf("--object-manager-port=%d", rc.Spec.ObjectManagerPort),
		fmt.Sprintf("--node-manager-port=%d", rc.Spec.NodeManagerPort),
	}

	if rc.Spec.ObjectStoreMemoryBytes != 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--object-store-memory=%d", rc.Spec.ObjectStoreMemoryBytes))
	}

	switch comp {
	case ComponentHead:
		cmdArgs = addHeadCmdArgs(rc, cmdArgs)
	case ComponentWorker:
		cmdArgs = addWorkerCmdArgs(rc, cmdArgs)
	}

	return cmdArgs
}

func processPorts(rc *dcv1alpha1.RayCluster, comp Component) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          "object-manager",
			ContainerPort: rc.Spec.ObjectManagerPort,
		},
		{
			Name:          "node-manager",
			ContainerPort: rc.Spec.NodeManagerPort,
		},
	}

	if comp == ComponentHead {
		ports = addHeadContainerPorts(rc, ports)
	}

	return ports
}

func processLabels(rc *dcv1alpha1.RayCluster, comp Component) map[string]string {
	labels := MetadataLabelsWithComponent(rc, comp)
	if rc.Spec.Labels != nil {
		labels = util.MergeStringMaps(rc.Spec.Labels, labels)
	}

	return labels
}

func addHeadContainerPorts(rc *dcv1alpha1.RayCluster, ports []corev1.ContainerPort) []corev1.ContainerPort {
	redisPorts := []corev1.ContainerPort{
		{
			Name:          "redis-primary",
			ContainerPort: rc.Spec.Port,
		},
	}
	for idx, port := range rc.Spec.RedisShardPorts {
		redisPorts = append(redisPorts, corev1.ContainerPort{
			Name:          fmt.Sprintf("redis-shard-%d", idx),
			ContainerPort: port,
		})
	}

	return append(ports, redisPorts...)
}

func addHeadCmdArgs(rc *dcv1alpha1.RayCluster, args []string) []string {
	headArgs := []string{
		"--head",
		fmt.Sprintf("--ray-client-server-port=%d", rc.Spec.ClientServerPort),
		fmt.Sprintf("--port=%d", rc.Spec.Port),
		fmt.Sprintf("--redis-shard-ports=%s", strings.Join(util.IntsToStrings(rc.Spec.RedisShardPorts), ",")),
	}

	if rc.Spec.EnableDashboard {
		dashArgs := []string{
			"--include-dashboard=true",
			"--dashboard-host=0.0.0.0",
			fmt.Sprintf("--dashboard-port=%d", rc.Spec.DashboardPort),
		}
		headArgs = append(headArgs, dashArgs...)
	}

	return append(args, headArgs...)
}

func addWorkerCmdArgs(rc *dcv1alpha1.RayCluster, args []string) []string {
	return append(args, fmt.Sprintf("--address=%s:%d", HeadServiceName(rc.Name), rc.Spec.Port))
}

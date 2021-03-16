package spark

import (
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

var (
	//defaultCmd = []string{"sp"}
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

// NewStatefulSet generates a Deployment configured to manage Spark cluster nodes.
// The configuration is based the provided spec and the desired Component workload.
func NewStatefulSet(rc *dcv1alpha1.SparkCluster, comp Component) (*appsv1.StatefulSet, error) {
	var replicas int32
	var nodeAttrs dcv1alpha1.SparkClusterNode

	switch comp {
	case ComponentHead:
		replicas = 1
		nodeAttrs = rc.Spec.Head.SparkClusterNode
	case ComponentWorker:
		replicas = *rc.Spec.Worker.Replicas
		nodeAttrs = rc.Spec.Worker.SparkClusterNode
	default:
		return nil, fmt.Errorf("invalid spark component: %q", comp)
	}

	imageRef, err := util.ParseImageDefinition(rc.Spec.Image)
	if err != nil {
		return nil, err
	}

	ports := processPorts(rc)
	labels := processLabels(rc, comp, nodeAttrs.Labels)
	envVars := append(append(defaultEnv, componentEnvVars(rc, comp)...), rc.Spec.EnvVars...)
	volumes := append(defaultVolumes, nodeAttrs.Volumes...)
	volumeMounts := append(defaultVolumeMounts, nodeAttrs.VolumeMounts...)

	volumeClaimTemplates, err := processVolumeClaimTemplates(nodeAttrs.AdditionalStorage)
	if err != nil {
		return nil, err
	}
	serviceAccountName := InstanceObjectName(rc.Name, ComponentNone)
	if rc.Spec.ServiceAccountName != "" {
		serviceAccountName = rc.Spec.ServiceAccountName
	}

	annotations := make(map[string]string)
	if nodeAttrs.Annotations != nil {
		for k, v := range nodeAttrs.Annotations {
			annotations[k] = v
		}
	}
	//TODO: Discuss with @Sonny and @Po about proper place for this
	annotations["sidecar.istio.io/inject"] = "false"
	//TODO: Chart defaults a specific security context if enabled. Always setting for now
	context := rc.Spec.PodSecurityContext
	if context == nil {
		context = &corev1.PodSecurityContext{
			RunAsUser: pointer.Int64Ptr(1001),
			FSGroup:   pointer.Int64Ptr(1001),
		}
	}

	statefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, comp),
			Namespace: rc.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: InstanceObjectName(rc.Name, comp),
			Replicas:    pointer.Int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: SelectorLabelsWithComponent(rc, comp),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					NodeSelector:       nodeAttrs.NodeSelector,
					Affinity:           nodeAttrs.Affinity,
					Tolerations:        nodeAttrs.Tolerations,
					InitContainers:     nodeAttrs.InitContainers,
					ImagePullSecrets:   rc.Spec.ImagePullSecrets,
					SecurityContext:    context,
					Containers: []corev1.Container{
						{
							Name: string(ApplicationName + "-" + comp),
							//Command:         defaultCmd,
							Image:           imageRef,
							ImagePullPolicy: rc.Spec.Image.PullPolicy,
							Ports:           ports,
							Env:             envVars,
							VolumeMounts:    volumeMounts,
							Resources:       nodeAttrs.Resources,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(int(rc.Spec.DashboardPort)),
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(int(rc.Spec.DashboardPort)),
									},
								},
							},
						},
					},
					Volumes: volumes,
				},
			},
			VolumeClaimTemplates: volumeClaimTemplates,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
		},
	}

	return statefulSet, nil
}

func processVolumeClaimTemplates(storage []dcv1alpha1.SparkAdditionalStorage) ([]corev1.PersistentVolumeClaim, error) {
	pvcs := make([]corev1.PersistentVolumeClaim, len(storage))
	for i, as := range storage {
		quantity, err := resource.ParseQuantity(as.Size)
		if err != nil {
			return nil, err
		}
		pvcs[i] = corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: as.Name,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: as.AccessModes,
				Resources: corev1.ResourceRequirements{
					Requests: map[corev1.ResourceName]resource.Quantity{
						corev1.ResourceMemory: quantity,
					},
				},
				StorageClassName: &as.StorageClass,
			},
		}
	}
	return pvcs, nil
}

func componentEnvVars(rc *dcv1alpha1.SparkCluster, comp Component) []corev1.EnvVar {
	var envVar []corev1.EnvVar
	if comp == ComponentHead {
		envVar = []corev1.EnvVar{
			{
				Name:  "SPARK_MASTER_PORT",
				Value: strconv.Itoa(int(rc.Spec.ClusterPort)),
			},
			{
				Name:  "SPARK_MASTER_WEBUI_PORT",
				Value: strconv.Itoa(int(rc.Spec.DashboardPort)),
			},
			{
				Name:  "SPARK_MODE",
				Value: "master",
			},
		}
	} else if comp == ComponentWorker {
		envVar = []corev1.EnvVar{
			{
				Name:  "SPARK_MASTER_URL",
				Value: "spark://" + HeadServiceName(rc.Name) + ":" + strconv.Itoa(int(rc.Spec.ClusterPort)),
			},
			{
				Name: "SPARK_WORKER_WEBUI_PORT",
				//TODO talk to @Po about whether this makes sense. Spark defaults to 8081 as default for this
				Value: strconv.Itoa(int(rc.Spec.DashboardPort)),
			},
			{
				Name:  "SPARK_MODE",
				Value: "worker",
			},
		}
	}
	return envVar
}

func processPorts(rc *dcv1alpha1.SparkCluster) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: rc.Spec.DashboardPort,
		},
		{
			Name:          "cluster",
			ContainerPort: rc.Spec.ClusterPort,
		},
	}

	//if comp == ComponentHead {
	//	ports = addHeadContainerPorts(rc, ports)
	//}

	return ports
}

func processLabels(rc *dcv1alpha1.SparkCluster, comp Component, extraLabels map[string]string) map[string]string {
	labels := MetadataLabelsWithComponent(rc, comp)
	if extraLabels != nil {
		labels = util.MergeStringMaps(extraLabels, labels)
	}

	return labels
}

//
//func addHeadContainerPorts(rc *dcv1alpha1.SparkCluster, ports []corev1.ContainerPort) []corev1.ContainerPort {
//	redisPorts := []corev1.ContainerPort{
//		{
//			Name:          "redis-primary",
//			ContainerPort: rc.Spec.Port,
//		},
//	}
//
//	return append(ports, redisPorts...)
//}

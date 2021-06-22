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

// NewStatefulSet generates a Deployment configured to manage Spark cluster nodes.
// The configuration is based the provided spec and the desired Component workload.
func NewStatefulSet(sc *dcv1alpha1.SparkCluster, comp Component) (*appsv1.StatefulSet, error) {
	var replicas int32
	var nodeAttrs dcv1alpha1.SparkClusterNode
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	var ports []corev1.ContainerPort

	switch comp {
	case ComponentMaster:
		replicas = 1
		nodeAttrs = sc.Spec.Master.SparkClusterNode
		ports = []corev1.ContainerPort{
			{
				Name:          "http",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: sc.Spec.TCPMasterWebPort,
			},
			{
				Name:          "cluster",
				ContainerPort: sc.Spec.ClusterPort,
			},
		}
	case ComponentWorker:
		replicas = *sc.Spec.Worker.Replicas
		nodeAttrs = sc.Spec.Worker.SparkClusterNode
		ports = []corev1.ContainerPort{
			{
				Name:          "http",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: sc.Spec.TCPWorkerWebPort,
			},
		}
	default:
		return nil, fmt.Errorf("invalid spark component: %q", comp)
	}

	imageRef, err := util.ParseImageDefinition(sc.Spec.Image)
	if err != nil {
		return nil, err
	}

	labels := processLabels(sc, comp, nodeAttrs.Labels)
	envVars := append(componentEnvVars(sc, comp), sc.Spec.EnvVars...)
	volumes = nodeAttrs.Volumes
	volumeMounts = nodeAttrs.VolumeMounts

	if nodeAttrs.FrameworkConfig != nil {
		cmVolume := getConfigMapVolume("spark-config", FrameworkConfigMapName(sc.Name, ComponentNone))
		cmVolumeMount := getConfigMapVolumeMount("spark-config", nodeAttrs.FrameworkConfig.Path, string(comp))

		volumes = append(volumes, cmVolume)
		volumeMounts = append(volumeMounts, cmVolumeMount)
	}
	if nodeAttrs.KeyTabConfig != nil {
		cmVolume := getConfigMapVolume("keytab", KeyTabConfigMapName(sc.Name, ComponentNone))
		cmVolumeMount := getConfigMapVolumeMount("keytab", nodeAttrs.KeyTabConfig.Path, string(comp))

		volumes = append(volumes, cmVolume)
		volumeMounts = append(volumeMounts, cmVolumeMount)
	}
	volumeClaimTemplates, err := processVolumeClaimTemplates(nodeAttrs.AdditionalStorage)
	if err != nil {
		return nil, err
	}
	serviceAccountName := InstanceObjectName(sc.Name, ComponentNone)
	if sc.Spec.ServiceAccountName != "" {
		serviceAccountName = sc.Spec.ServiceAccountName
	}
	annotations := make(map[string]string)
	if nodeAttrs.Annotations != nil {
		for k, v := range nodeAttrs.Annotations {
			annotations[k] = v
		}
	}
	context := sc.Spec.PodSecurityContext //TODO: Chart defaults a specific security context if enabled. Always setting for now
	if context == nil {
		const DefaultUser = 1001
		const DefaultFSGroup = 1001
		context = &corev1.PodSecurityContext{
			RunAsUser: pointer.Int64Ptr(DefaultUser),
			FSGroup:   pointer.Int64Ptr(DefaultFSGroup),
		}
	}
	podSpec := getPodSpec(sc,
		comp,
		serviceAccountName,
		nodeAttrs,
		context,
		imageRef,
		ports,
		envVars,
		volumeMounts,
		volumes)
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(sc.Name, comp),
			Namespace: sc.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: InstanceObjectName(sc.Name, comp),
			Replicas:    pointer.Int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: SelectorLabelsWithComponent(sc, comp),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: podSpec,
			},
			VolumeClaimTemplates: volumeClaimTemplates,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			PodManagementPolicy: appsv1.ParallelPodManagement,
		},
	}
	return statefulSet, nil
}

func getConfigMapVolumeMount(name string, path string, subPath string) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      name,
		MountPath: path,
		SubPath:   subPath,
	}
}

func getConfigMapVolume(name string, cmName string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cmName,
				},
			},
		},
	}
}

func getPodSpec(sc *dcv1alpha1.SparkCluster,
	comp Component,
	serviceAccountName string,
	nodeAttrs dcv1alpha1.SparkClusterNode,
	context *corev1.PodSecurityContext,
	imageRef string,
	ports []corev1.ContainerPort,
	envVars []corev1.EnvVar,
	volumeMounts []corev1.VolumeMount,
	volumes []corev1.Volume) corev1.PodSpec {
	return corev1.PodSpec{
		ServiceAccountName: serviceAccountName,
		NodeSelector:       nodeAttrs.NodeSelector,
		Affinity:           nodeAttrs.Affinity,
		Tolerations:        nodeAttrs.Tolerations,
		InitContainers:     nodeAttrs.InitContainers,
		ImagePullSecrets:   sc.Spec.ImagePullSecrets,
		SecurityContext:    context,
		Containers: []corev1.Container{
			{
				Name:            InstanceObjectName(sc.Name, comp),
				Image:           imageRef,
				ImagePullPolicy: sc.Spec.Image.PullPolicy,
				Ports:           ports,
				Env:             envVars,
				VolumeMounts:    volumeMounts,
				Resources:       nodeAttrs.Resources,
				LivenessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/",
							Port: intstr.FromInt(int(sc.Spec.TCPWorkerWebPort)),
						},
					},
				},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/",
							Port: intstr.FromInt(int(sc.Spec.TCPWorkerWebPort)),
						},
					},
				},
			},
		},
		Volumes: volumes,
	}
}

func processVolumeClaimTemplates(storage []dcv1alpha1.SparkAdditionalStorage) ([]corev1.PersistentVolumeClaim, error) {
	pvcs := make([]corev1.PersistentVolumeClaim, len(storage))
	for i, as := range storage {
		quantity, err := resource.ParseQuantity(as.Size)
		if err != nil {
			return nil, err
		}
		fs := corev1.PersistentVolumeFilesystem
		pvcs[i] = corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: as.Name,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: as.AccessModes,
				Resources: corev1.ResourceRequirements{
					Requests: map[corev1.ResourceName]resource.Quantity{
						corev1.ResourceStorage: quantity,
					},
				},
				StorageClassName: &as.StorageClass,
				VolumeMode:       &fs,
			},
		}
	}
	return pvcs, nil
}

func componentEnvVars(sc *dcv1alpha1.SparkCluster, comp Component) []corev1.EnvVar {
	var envVar []corev1.EnvVar
	if comp == ComponentMaster {
		envVar = []corev1.EnvVar{
			{
				Name:  "SPARK_MASTER_PORT",
				Value: strconv.Itoa(int(sc.Spec.ClusterPort)),
			},
			{
				Name:  "SPARK_MASTER_WEBUI_PORT",
				Value: strconv.Itoa(int(sc.Spec.TCPMasterWebPort)),
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
				Value: "spark://" + HeadServiceName(sc.Name) + ":" + strconv.Itoa(int(sc.Spec.ClusterPort)),
			},
			{
				Name:  "SPARK_WORKER_WEBUI_PORT",
				Value: strconv.Itoa(int(sc.Spec.TCPWorkerWebPort)),
			},
			{
				Name:  "SPARK_WORKER_PORT",
				Value: strconv.Itoa(int(sc.Spec.ClusterPort)),
			},
			{
				Name:  "SPARK_MODE",
				Value: "worker",
			},
		}
	}
	return envVar
}

func processLabels(sc *dcv1alpha1.SparkCluster, comp Component, extraLabels map[string]string) map[string]string {
	labels := MetadataLabelsWithComponent(sc, comp)
	if extraLabels != nil {
		labels = util.MergeStringMaps(extraLabels, labels)
	}
	return labels
}

package spark

import (
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

const frameworkConfigMountPath = "/opt/bitnami/spark/conf/spark-defaults.conf"

// NewStatefulSet generates a Deployment configured to manage Spark cluster nodes.
// The configuration is based the provided spec and the desired Component workload.
func NewStatefulSet(sc *dcv1alpha1.SparkCluster, comp Component) (*appsv1.StatefulSet, error) {
	var replicas int32
	var nodeAttrs dcv1alpha1.SparkClusterNode
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	var ports []corev1.ContainerPort
	var securityContext *corev1.SecurityContext

	switch comp {
	case ComponentMaster:
		replicas = 1
		nodeAttrs = sc.Spec.Master
		ports = []corev1.ContainerPort{
			{
				Name:          "http",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: sc.Spec.MasterWebPort,
			},
			{
				Name:          "cluster",
				ContainerPort: sc.Spec.ClusterPort,
			},
		}
		securityContext = sc.Spec.Master.SecurityContext
	case ComponentWorker:
		replicas = *sc.Spec.Worker.Replicas
		nodeAttrs = sc.Spec.Worker.SparkClusterNode
		ports = []corev1.ContainerPort{
			{
				Name:          "http",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: sc.Spec.WorkerWebPort,
			},
		}
		securityContext = sc.Spec.Worker.SecurityContext
	default:
		return nil, fmt.Errorf("invalid spark component: %q", comp)
	}

	imageRef, err := util.ParseImageDefinition(sc.Spec.Image)
	if err != nil {
		return nil, err
	}

	labels := AddGlobalLabels(MetadataLabelsWithComponent(sc, comp), nodeAttrs.Labels)
	envVars := append(componentEnvVars(sc, comp), sc.Spec.EnvVars...)
	volumes = nodeAttrs.Volumes
	volumeMounts = nodeAttrs.VolumeMounts
	volumeClaimTemplates := processPVCTemplates(sc, nodeAttrs.VolumeClaimTemplates)

	if nodeAttrs.DefaultConfiguration != nil {
		cmVolume := getConfigMapVolume("spark-config", FrameworkConfigMapName(sc.Name, ComponentNone))
		cmVolumeMount := getConfigMapVolumeMount("spark-config", frameworkConfigMountPath, string(comp))

		volumes = append(volumes, cmVolume)
		volumeMounts = append(volumeMounts, cmVolumeMount)
	}
	if sc.Spec.KerberosKeytab != nil {
		cmVolume := getConfigMapVolume("keytab", KeyTabConfigMapName(sc.Name, ComponentNone))
		cmVolumeMount := getConfigMapVolumeMount("keytab", sc.Spec.KerberosKeytab.MountPath, string(comp))

		volumes = append(volumes, cmVolume)
		volumeMounts = append(volumeMounts, cmVolumeMount)
	}

	serviceAccountName := InstanceObjectName(sc.Name, ComponentNone)
	if sc.Spec.ServiceAccount.Name != "" {
		serviceAccountName = sc.Spec.ServiceAccount.Name
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
			RunAsUser: pointer.Int64(DefaultUser),
			FSGroup:   pointer.Int64(DefaultFSGroup),
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
		volumes,
		securityContext)

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(sc.Name, comp),
			Namespace: sc.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: InstanceObjectName(sc.Name, comp),
			Replicas:    pointer.Int32(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: SelectorLabelsWithComponent(sc, comp),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      util.MergeStringMaps(sc.Spec.EnvoyFilterLabels, labels),
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
	volumes []corev1.Volume,
	securityContext *corev1.SecurityContext) corev1.PodSpec {
	var port intstr.IntOrString

	switch comp {
	case ComponentMaster:
		port = intstr.FromInt(int(sc.Spec.MasterWebPort))
	case ComponentWorker:
		port = intstr.FromInt(int(sc.Spec.WorkerWebPort))
	case ComponentNone:
		port = intstr.FromInt(int(sc.Spec.MasterWebPort))
	}

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
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/",
							Port: port,
						},
					},
				},
				ReadinessProbe: &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/",
							Port: port,
						},
					},
				},
				SecurityContext: securityContext,
			},
		},
		Volumes: volumes,
	}
}

func processPVCTemplates(
	sc *dcv1alpha1.SparkCluster,
	vcts []dcv1alpha1.PersistentVolumeClaimTemplate,
) (pvcTmpls []corev1.PersistentVolumeClaim) {
	mode := corev1.PersistentVolumeFilesystem

	for _, vct := range vcts {
		spec := vct.Spec.DeepCopy()
		spec.VolumeMode = &mode

		pvcTmpls = append(pvcTmpls, corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:   vct.Name,
				Labels: sc.Spec.GlobalLabels,
			},
			Spec: vct.Spec,
		})
	}

	return
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
				Value: strconv.Itoa(int(sc.Spec.MasterWebPort)),
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
				Value: "spark://" + MasterServiceName(sc.Name) + ":" + strconv.Itoa(int(sc.Spec.ClusterPort)),
			},
			{
				Name:  "SPARK_WORKER_WEBUI_PORT",
				Value: strconv.Itoa(int(sc.Spec.WorkerWebPort)),
			},
			{
				Name:  "SPARK_WORKER_PORT",
				Value: strconv.Itoa(int(sc.Spec.ClusterPort)),
			},
			{
				Name:  "SPARK_MODE",
				Value: "worker",
			},
			{
				Name:  "SPARK_WORKER_MEMORY",
				Value: sc.Spec.WorkerMemoryLimit,
			},
			{
				Name:  "SPARK_WORKER_CORES",
				Value: sc.Spec.Worker.Resources.Requests.Cpu().String(),
			},
		}
	}
	return envVar
}

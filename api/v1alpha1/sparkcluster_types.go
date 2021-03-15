package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SparkClusterNode defines attributes common to all spark node types.
type SparkClusterNode struct {
	// Labels applied to spark pods in addition to stock labels.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations applied to spark pods.
	Annotations map[string]string `json:"annotations,omitempty"`

	// NodeSelector applied to spark pods.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Affinity applied to spark pods.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Tolerations applied to spark pods.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// InitContainers added to spark pods.
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// Volumes added to spark pods.
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// VolumeMounts added to spark containers.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// Resources are the requests and limits applied to spark containers.
	Resources corev1.ResourceRequirements `json:"resources"`

	AdditionalStorage []SparkAdditionalStorage `json:"additionalStorage,omitempty"`
}

type SparkAdditionalStorage struct {
	AccessModes  []corev1.PersistentVolumeAccessMode `json:"accessModes""`
	Size         string                              `json:"size"`
	StorageClass string                              `json:"size"`
	Name         string                              `json:"name"`
}

// SparkClusterHead defines head-specific pod settings.
type SparkClusterHead struct {
	SparkClusterNode `json:",inline"`
}

// SparkClusterWorker defines worker-specific pod settings.
type SparkClusterWorker struct {
	SparkClusterNode `json:",inline"`

	// Replicas configures the total number of workers in the cluster.
	// This field behaves differently when Autoscaling is enabled. If Autoscaling.MinReplicas is unspecified, then the
	// minimum number of replicas will be set to this value. Additionally, you can specify an "initial cluster size" by
	// setting this field to some value above the minimum number of replicas.
	Replicas *int32 `json:"replicas,omitempty"`
}

// SparkClusterSpec defines the desired state of a SparkCluster resource.
type SparkClusterSpec struct {
	// Image used to launch head and worker nodes.
	Image *OCIImageDefinition `json:"image,omitempty"`

	// ImagePullSecrets are references to secrets with credentials to private registries used to pull images.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Autoscaling parameters used to scale up/down spark worker nodes.
	Autoscaling *Autoscaling `json:"autoscaling,omitempty"`

	// Port is the port of the head spark process.
	Port int32 `json:"port,omitempty"`

	// ClientServerPort is the port number to which the spark client server will bind.
	ClientServerPort int32 `json:"clientServerPort,omitempty"`

	// HttpPort is the port on which spark pods expose http
	HttpPort int32 `json:"httpPort,omitempty"`

	// Cluster port is the port on which the spark protocol communicates
	ClusterPort int32 `json:"clusterPort,omitempty"`

	// ObjectStoreMemoryBytes is initial amount of memory with which to start the object store.
	ObjectStoreMemoryBytes *int64 `json:"objectStoreMemoryBytes,omitempty"`

	// DashboardPort is the port used by the dashboard server.
	DashboardPort int32 `json:"dashboardPort,omitempty"`

	// EnableDashboard starts the dashboard web UI.
	EnableDashboard *bool `json:"enableDashboard,omitempty"`

	// EnableNetworkPolicy will create a network policy that restricts ingress to the cluster.
	EnableNetworkPolicy *bool `json:"enableNetworkPolicy,omitempty"`

	// NetworkPolicyClientLabels will create a pod selector clause for each set of labels.
	// This is used to grant ingress access to one or more groups of external pods and is
	// only applicable when EnableNetworkPolicy is true.
	NetworkPolicyClientLabels []map[string]string `json:"networkPolicyClientLabels,omitempty"`

	// PodSecurityPolicy name can be provided to govern execution of the spark processes within pods.
	PodSecurityPolicy string `json:"podSecurityPolicy,omitempty"`

	// PodSecurityContext added to every spark pod.
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ServiceAccountName will disable the creation of a dedicated cluster service account.
	// The service account referenced by the provided name will be used instead.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// EnvVars added to every spark pod container.
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`

	// Head node configuration parameters.
	Head SparkClusterHead `json:"head,omitempty"`

	// Worker node configuration parameters.
	Worker SparkClusterWorker `json:"worker,omitempty"`
}

// SparkClusterStatus defines the observed state of a SparkCluster resource.
type SparkClusterStatus struct {
	// Nodes that comprise the cluster.
	Nodes []string `json:"nodes,omitempty"`

	// WorkerReplicas is the scale.status.replicas subresource field.
	WorkerReplicas int32 `json:"workerReplicas,omitempty"`
	// WorkerSelector is the scale.status.selector subresource field.
	WorkerSelector string `json:"workerSelector,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=spark
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.worker.replicas,statuspath=.status.workerReplicas,selectorpath=.status.workerSelector
//+kubebuilder:printcolumn:name="Worker Replicas",type=integer,JSONPath=".spec.worker.replicas"
//+kubebuilder:printcolumn:name="Image",type=string,JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

// SparkCluster is the Schema for the sparkclusters API.
type SparkCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SparkClusterSpec   `json:"spec,omitempty"`
	Status SparkClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SparkClusterList contains a list of SparkCluster resources.
type SparkClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SparkCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SparkCluster{}, &SparkClusterList{})
}

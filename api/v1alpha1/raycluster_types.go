package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RayClusterNode defines attributes common to all ray node types.
type RayClusterNode struct {
	// Labels applied to ray pods in addition to stock labels.
	//+kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations applied to ray pods.
	//+kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// NodeSelector applied to ray pods.
	//+kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Affinity applied to ray pods.
	//+kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Tolerations applied to ray pods.
	//+kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// InitContainers added to ray pods.
	//+kubebuilder:validation:Optional
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// Volumes added to ray pods.
	//+kubebuilder:validation:Optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// VolumeMounts added to ray containers.
	//+kubebuilder:validation:Optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// Resources are the requests and limits applied to ray containers.
	//+kubebuilder:validation:Optional
	Resources corev1.ResourceRequirements `json:"resources"`
}

// RayClusterHead defines head-specific pod settings.
type RayClusterHead struct {
	RayClusterNode `json:",inline"`
}

// RayClusterWorker defines worker-specific pod settings.
type RayClusterWorker struct {
	RayClusterNode `json:",inline"`

	// Replicas configures the total number of workers in the cluster.
	// This field behaves differently when Autoscaling is enabled. If Autoscaling.MinReplicas is unspecified, then the
	// minimum number of replicas will be set to this value. Additionally, you can specify an "initial cluster size" by
	// setting this field to some value above the minimum number of replicas.
	//+kubebuilder:default=1
	//+kubebuilder:validation:Optional
	Replicas int32 `json:"replicas"`
}

// RayClusterSpec defines the desired state of a RayCluster resource.
type RayClusterSpec struct {
	// Image used to launch head and worker nodes.
	//+kubebuilder:default={repository: "rayproject/ray", tag: "1.2.0-cpu"}
	//+kubebuilder:validation:Optional
	Image *OCIImageDefinition `json:"image,omitempty"`

	// ImagePullSecrets are references to secrets with credentials to private registries used to pull images.
	//+kubebuilder:validation:Optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Autoscaling parameters used to scale up/down ray worker nodes.
	//+kubebuilder:validation:Optional
	Autoscaling *Autoscaling `json:"autoscaling,omitempty"`

	// Port is the port of the head ray process.
	//+kubebuilder:default=6379
	//+kubebuilder:validation:Optional
	Port int32 `json:"port"`

	// RedisShardPorts is a list of ports for non-primary Redis shards.
	//+kubebuilder:default={6380,6381}
	//+kubebuilder:validation:Optional
	RedisShardPorts []int32 `json:"redisShardPorts,omitempty"`

	// ClientServerPort is the port number to which the ray client server will bind.
	//+kubebuilder:default=10001
	//+kubebuilder:validation:Optional
	ClientServerPort int32 `json:"clientServerPort"`

	// ObjectManagerPort is the raylet port for the object manager.
	//+kubebuilder:default=2384
	//+kubebuilder:validation:Optional
	ObjectManagerPort int32 `json:"objectManagerPort"`

	// NodeManagerPort is the raylet port for the node manager.
	//+kubebuilder:default=2385
	//+kubebuilder:validation:Optional
	NodeManagerPort int32 `json:"nodeManagerPort"`

	// ObjectStoreMemoryBytes is initial amount of memory with which to start the object store.
	//+kubebuilder:validation:Optional
	ObjectStoreMemoryBytes *int64 `json:"objectStoreMemoryBytes,omitempty"`

	// DashboardPort is the port used by the dashboard server.
	//+kubebuilder:default=8265
	//+kubebuilder:validation:Optional
	DashboardPort int32 `json:"dashboardPort"`

	// EnableDashboard starts the dashboard web UI.
	//+kubebuilder:default=true
	//+kubebuilder:validation:Optional
	EnableDashboard bool `json:"enableDashboard"`

	// EnableNetworkPolicy will create a network policy that restricts ingress to the cluster.
	//+kubebuilder:default=true
	//+kubebuilder:validation:Optional
	EnableNetworkPolicy bool `json:"enableNetworkPolicy"`

	// NetworkPolicyClientLabels will create a pod selector clause for each set of labels.
	// This is used to grant ingress access to one or more groups of external pods and is
	// only applicable when EnableNetworkPolicy is true.
	//+kubebuilder:default={{"ray-client": "true"}}
	//+kubebuilder:validation:Optional
	NetworkPolicyClientLabels []map[string]string `json:"networkPolicyClientLabels,omitempty"`

	// PodSecurityPolicy name can be provided to govern execution of the ray processes within pods.
	//+kubebuilder:validation:Optional
	PodSecurityPolicy string `json:"podSecurityPolicy"`

	// PodSecurityContext added to every ray pod.
	//+kubebuilder:validation:Optional
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ServiceAccountName will disable the creation of a dedicated cluster service account.
	// The service account referenced by the provided name will be used instead.
	//+kubebuilder:validation:Optional
	ServiceAccountName string `json:"serviceAccountName"`

	// EnvVars added to all every ray pod container.
	//+kubebuilder:validation:Optional
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`

	// Head node configuration parameters.
	//+kubebuilder:validation:Optional
	Head RayClusterHead `json:"head"`

	// Worker node configuration parameters.
	//+kubebuilder:default={replicas: 1}
	//+kubebuilder:validation:Optional
	Worker RayClusterWorker `json:"worker"`
}

// RayClusterStatus defines the observed state of a RayCluster resource.
type RayClusterStatus struct {
	// Nodes that comprise the cluster.
	Nodes []string `json:"nodes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=ray
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Worker Count",type=integer,JSONPath=".spec.worker.replicas"
//+kubebuilder:printcolumn:name="Image",type=string,JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

// RayCluster is the Schema for the rayclusters API.
type RayCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RayClusterSpec   `json:"spec,omitempty"`
	Status RayClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RayClusterList contains a list of RayCluster resources.
type RayClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RayCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RayCluster{}, &RayClusterList{})
}

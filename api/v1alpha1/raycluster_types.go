package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RayClusterNode defines attributes common to all ray node types.
type RayClusterNode struct {
	// Labels applied to ray pods in addition to stock labels.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations applied to ray pods.
	Annotations map[string]string `json:"annotations,omitempty"`

	// NodeSelector applied to ray pods.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Affinity applied to ray pods.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Tolerations applied to ray pods.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// InitContainers added to ray pods.
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// Volumes added to ray pods.
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// VolumeMounts added to ray containers.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// VolumeClaimTemplates is a list of claims that ray pods are allowed to
	// reference. You can enable dynamic provisioning of additional storage
	// on-demand by using a storage class provisioner.
	VolumeClaimTemplates []PersistentVolumeClaimTemplate `json:"volumeClaimTemplates,omitempty"`

	// Resources are the requests and limits applied to ray containers.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// RayClusterHead defines head-specific pod settings.
type RayClusterHead struct {
	RayClusterNode `json:",inline"`
}

// RayClusterWorker defines worker-specific pod settings.
type RayClusterWorker struct {
	RayClusterNode `json:",inline"`

	// Replicas configures the total number of workers in the cluster. This
	// field behaves differently when Autoscaling is enabled. If
	// Autoscaling.MinReplicas is unspecified, then the minimum number of
	// replicas will be set to this value. Additionally, you can specify an
	// "initial cluster size" by setting this field to some value above the
	// minimum number of replicas.
	Replicas *int32 `json:"replicas,omitempty"`
}

// RayClusterNetworkPolicy defines network policy configuration options.
type RayClusterNetworkPolicy struct {
	// Enabled controls the creation of network policies that limit and provide
	// ingress access to the cluster nodes.
	Enabled *bool `json:"enabled,omitempty"`

	// ClientServerLabels defines the pod selector clause that grant ingress
	// access to the head client server port.
	ClientServerLabels map[string]string `json:"clientServerLabels,omitempty"`

	// DashboardLabels defines the pod selector clause used to grant ingress
	// access to the head dashboard port.
	DashboardLabels map[string]string `json:"dashboardLabels,omitempty"`
}

// RayClusterSpec defines the desired state of a RayCluster resource.
type RayClusterSpec struct {
	// Image used to launch head and worker nodes.
	Image *OCIImageDefinition `json:"image,omitempty"`

	// ImagePullSecrets are references to secrets with credentials to private
	// registries used to pull images.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Autoscaling parameters used to scale up/down ray worker nodes.
	Autoscaling *Autoscaling `json:"autoscaling,omitempty"`

	// NetworkPolicy parameters that grant intra-cluster and external network
	// access to cluster nodes.
	NetworkPolicy RayClusterNetworkPolicy `json:"networkPolicy,omitempty"`

	// Port is the port of the head ray process.
	Port int32 `json:"port,omitempty"`

	// RedisShardPorts is a list of ports for non-primary Redis shards.
	RedisShardPorts []int32 `json:"redisShardPorts,omitempty"`

	// ClientServerPort is the port number to which the ray client server will
	// bind. This port is used by external clients to submit work.
	ClientServerPort int32 `json:"clientServerPort,omitempty"`

	// ObjectManagerPort is the raylet port for the object manager.
	ObjectManagerPort int32 `json:"objectManagerPort,omitempty"`

	// NodeManagerPort is the raylet port for the node manager.
	NodeManagerPort int32 `json:"nodeManagerPort,omitempty"`

	// ObjectStoreMemoryBytes is initial amount of memory with which to start
	// the object store.
	ObjectStoreMemoryBytes *int64 `json:"objectStoreMemoryBytes,omitempty"`

	// DashboardPort is the port used by the dashboard server.
	DashboardPort int32 `json:"dashboardPort,omitempty"`

	// EnableDashboard starts the dashboard web UI.
	EnableDashboard *bool `json:"enableDashboard,omitempty"`

	// PodSecurityPolicy name can be provided to govern execution of the ray
	// processes within pods.
	PodSecurityPolicy string `json:"podSecurityPolicy,omitempty"`

	// PodSecurityContext added to every ray pod.
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ServiceAccountName will disable the creation of a dedicated cluster
	// service account. The service account referenced by the provided name
	// will be used instead.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// EnvVars added to all every ray pod container.
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`

	// Head node configuration parameters.
	Head RayClusterHead `json:"head,omitempty"`

	// Worker node configuration parameters.
	Worker RayClusterWorker `json:"worker,omitempty"`
}

// RayClusterStatus defines the observed state of a RayCluster resource.
type RayClusterStatus struct {
	// Nodes that comprise the cluster.
	Nodes []string `json:"nodes,omitempty"`

	// WorkerReplicas is the scale.status.replicas subresource field.
	WorkerReplicas int32 `json:"workerReplicas,omitempty"`

	// WorkerSelector is the scale.status.selector subresource field.
	WorkerSelector string `json:"workerSelector,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=ray
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.worker.replicas,statuspath=.status.workerReplicas,selectorpath=.status.workerSelector
//+kubebuilder:printcolumn:name="Worker Replicas",type=integer,JSONPath=".spec.worker.replicas"
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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RayClusterSpec defines the desired state of a RayCluster resource.
type RayClusterSpec struct {
	// Image used to launch head and worker nodes.
	// +kubebuilder:default={repository: "ray-project/ray-ml", tag: "1.2.0-cpu"}
	Image OCIImageDefinition `json:"image"`

	// WorkerReplicaCount configures the total number of workers in the cluster.
	// +kubebuilder:default=2
	WorkerReplicaCount int `json:"workerReplicaCount"`

	// HeadPort is the port of the head ray process.
	// +kubebuilder:default=6379
	HeadPort int `json:"port"`

	// RedisShardPorts is a list of ports for non-primary Redis shards.
	// +kubebuilder:default={6380,6381}
	RedisShardPorts []int `json:"redisShardPorts"`

	// ObjectManagerPort is the raylet port for the object manager.
	// +kubebuilder:default=2384
	ObjectManagerPort int `json:"objectManagerPort"`

	// NodeManagerPort is the raylet port for the node manager.
	// +kubebuilder:default=2385
	NodeManagerPort int `json:"nodeManagerPort"`

	// ObjectStoreMemoryBytes is amount of memory with which to start the object store.
	ObjectStoreMemoryBytes int `json:"objectStoreMemoryBytes"`

	// EnableDashboard starts the dashboard web UI.
	// +kubebuilder:default=true
	EnableDashboard bool `json:"enableDashboard"`

	// DashboardPort is the port used by the dashboard server.
	// +kubebuilder:default=8265
	DashboardPort int `json:"dashboardPort"`

	// EnableNetworkPolicy will create a network policy that restricts ingress to the cluster.
	// +kubebuilder:default=true
	EnableNetworkPolicy bool `json:"enableNetworkPolicy"`

	// EnablePodSecurityPolicy will create a pod security policy that enforces secure execution of the ray process.
	// +kubebuilder:default=true
	EnablePodSecurityPolicy bool `json:"enablePodSecurityPolicy"`

	// Labels additionally applied to cluster nodes.
	Labels map[string]string `json:"labels"`

	// Annotations applied to cluster nodes.
	Annotations map[string]string `json:"annotations"`

	// Resources requested and limits applied to cluster nodes.
	// +kubebuilder:default={requests: {cpu: "100m", memory: "512Mi"}}
	Resources corev1.ResourceRequirements `json:"resources"`

	// Tolerations applied to cluster nodes.
	Tolerations []corev1.Toleration `json:"tolerations"`

	// InitContainers added to cluster node pods.
	InitContainers []corev1.Container `json:"initContainers"`
}

// RayClusterStatus defines the observed state of a RayCluster resource.
type RayClusterStatus struct {
	// Nodes that comprise the cluster.
	Nodes []string `json:"nodes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=ray
// +kubebuilder:subresource:status

// RayCluster is the Schema for the rayclusters API.
type RayCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RayClusterSpec   `json:"spec,omitempty"`
	Status RayClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RayClusterList contains a list of RayCluster resources.
type RayClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RayCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RayCluster{}, &RayClusterList{})
}

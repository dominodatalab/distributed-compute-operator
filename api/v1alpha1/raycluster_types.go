package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RayClusterSpec defines the desired state of a RayCluster resource.
type RayClusterSpec struct {
	// Image used to launch head and worker nodes.
	// +kubebuilder:default={repository: "rayproject/ray-ml", tag: "1.2.0-cpu"}
	// +kubebuilder:validation:Optional
	Image OCIImageDefinition `json:"image"`

	// WorkerReplicaCount configures the total number of workers in the cluster.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	WorkerReplicaCount int32 `json:"workerReplicaCount"`

	// HeadPort is the port of the head ray process.
	// +kubebuilder:default=6379
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65353
	HeadPort int32 `json:"port"`

	// RedisShardPorts is a list of ports for non-primary Redis shards.
	// +kubebuilder:default={6380,6381}
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	RedisShardPorts []int32 `json:"redisShardPorts"`

	// ObjectManagerPort is the raylet port for the object manager.
	// +kubebuilder:default=2384
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65353
	ObjectManagerPort int32 `json:"objectManagerPort"`

	// NodeManagerPort is the raylet port for the node manager.
	// +kubebuilder:default=2385
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65353
	NodeManagerPort int32 `json:"nodeManagerPort"`

	// ObjectStoreMemoryBytes is initial amount of memory with which to start the object store.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	ObjectStoreMemoryBytes int32 `json:"objectStoreMemoryBytes"`

	// DashboardPort is the port used by the dashboard server.
	// +kubebuilder:default=8265
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65353
	DashboardPort int32 `json:"dashboardPort"`

	// EnableDashboard starts the dashboard web UI.
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	EnableDashboard bool `json:"enableDashboard"`

	// EnableNetworkPolicy will create a network policy that restricts ingress to the cluster.
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	EnableNetworkPolicy bool `json:"enableNetworkPolicy"`

	// EnablePodSecurityPolicy will create a pod security policy that enforces secure execution of the ray process.
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	EnablePodSecurityPolicy bool `json:"enablePodSecurityPolicy"`

	// Labels applied to cluster resources in addition to stock labels.
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels"`

	// Annotations applied to cluster pods.
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations"`

	// NodeSelector applied to cluster pods.
	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector"`

	// Affinity applied to cluster pods.
	// +kubebuilder:validation:Optional
	Affinity *corev1.Affinity `json:"affinity"`

	// Resources requested and limits applied to all ray containers.
	// +kubebuilder:default={requests: {cpu: "100m", memory: "512Mi"}}
	// +kubebuilder:validation:Optional
	Resources corev1.ResourceRequirements `json:"resources"`

	// Tolerations applied to cluster pods.
	// +kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations"`

	// InitContainers added to cluster pods.
	// +kubebuilder:validation:Optional
	InitContainers []corev1.Container `json:"initContainers"`

	// Volumes added to cluster pods.
	// +kubebuilder:validation:Optional
	Volumes []corev1.Volume `json:"volumes"`

	// VolumeMounts added to all ray containers.
	// +kubebuilder:validation:Optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts"`
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

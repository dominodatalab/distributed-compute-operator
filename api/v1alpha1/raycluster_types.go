package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// RayClusterWorker defines worker-specific pod settings.
type RayClusterWorker struct {
	WorkloadConfig `json:",inline"`

	// Replicas configures the total number of workers in the cluster. This
	// field behaves differently when Autoscaling is enabled. If
	// Autoscaling.MinReplicas is unspecified, then the minimum number of
	// replicas will be set to this value. Additionally, you can specify an
	// "initial cluster size" by setting this field to some value above the
	// minimum number of replicas.
	Replicas *int32 `json:"replicas,omitempty"`
}

// RayClusterSpec defines the desired state of a RayCluster resource.
type RayClusterSpec struct {
	ScalableClusterConfig `json:",inline"`

	// Head node configuration parameters.
	Head WorkloadConfig `json:"head,omitempty"`
	// Worker node configuration parameters.
	Worker RayClusterWorker `json:"worker,omitempty"`

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
	// GCSServerPort is the port for the global control store.
	GCSServerPort int32 `json:"gcsServerPort,omitempty"`
	// WorkerPorts specifies the range of ports used by worker processes.
	WorkerPorts []int32 `json:"workerPorts,omitempty"`
	// ObjectStoreMemoryBytes is initial amount of memory with which to start
	// the object store.
	ObjectStoreMemoryBytes *int64 `json:"objectStoreMemoryBytes,omitempty"`
	// DashboardPort is the port used by the dashboard server.
	DashboardPort int32 `json:"dashboardPort,omitempty"`
	// EnableDashboard starts the dashboard web UI.
	EnableDashboard *bool `json:"enableDashboard,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=ray
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.worker.replicas,statuspath=.status.workerReplicas,selectorpath=.status.workerSelector
//+kubebuilder:printcolumn:name="Workers",type=integer,JSONPath=".spec.worker.replicas"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Image",type=string,JSONPath=".status.image"
//+kubebuilder:printcolumn:name="Network Policy",type=boolean,JSONPath=".spec.networkPolicy.enabled",priority=10
//+kubebuilder:printcolumn:name="Pods",type=string,JSONPath=".status.nodes",priority=10

// RayCluster is the Schema for the rayclusters API.
type RayCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RayClusterSpec      `json:"spec,omitempty"`
	Status ClusterStatusConfig `json:"status,omitempty"`
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

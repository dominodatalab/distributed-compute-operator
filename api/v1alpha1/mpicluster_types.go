package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MPIClusterWorker defines worker-specific workload settings.
type MPIClusterWorker struct {
	WorkloadConfig  `json:",inline"`
	Replicas        *int32 `json:"replicas,omitempty"`
	SharedSSHSecret string `json:"sharedSSHSecret"`
	UserName        string `json:"userName,omitempty"`
	UserID          *int64 `json:"userID,omitempty"`
	GroupName       string `json:"groupName,omitempty"`
	GroupID         *int64 `json:"groupID,omitempty"`
	HomeDir         string `json:"homeDir,omitempty"`
}

// MPIClusterSpec defines the desired state of MPICluster.
type MPIClusterSpec struct {
	ClusterConfig `json:",inline"`
	Worker        MPIClusterWorker `json:"worker,omitempty"`

	// WorkerPorts specifies the range of ports used by worker processes for communication.
	WorkerPorts []int32 `json:"workerPorts,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=mpi
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Workers",type=integer,JSONPath=".spec.worker.replicas"
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=".status.clusterStatus"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Image",type=string,JSONPath=".status.image",priority=10
//+kubebuilder:printcolumn:name="Bound PSP",type=string,JSONPath=".spec.podSecurityPolicy",priority=10
//+kubebuilder:printcolumn:name="Network Policy",type=boolean,JSONPath=".spec.networkPolicy.enabled",priority=10
//+kubebuilder:printcolumn:name="Pods",type=string,JSONPath=".status.nodes",priority=10

// MPICluster is the Schema for the MPI Clusters API.
type MPICluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MPIClusterSpec      `json:"spec,omitempty"`
	Status            ClusterStatusConfig `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MPIClusterList contains a list of MPICluster.
type MPIClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MPICluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MPICluster{}, &MPIClusterList{})
}

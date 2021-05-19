package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaskClusterWorker defines worker-specific workload settings.
type DaskClusterWorker struct {
	WorkloadConfig `json:",inline"`
	Replicas       int32 `json:"replicas,omitempty"`
}

// DaskClusterSpec defines the desired state of DaskCluster.
type DaskClusterSpec struct {
	ClusterConfig `json:",inline"`

	Scheduler WorkloadConfig    `json:"scheduler,omitempty"`
	Worker    DaskClusterWorker `json:"worker,omitempty"`

	SchedulerPort int32 `json:"schedulerPort,omitempty"`
	DashboardPort int32 `json:"dashboardPort,omitempty"`
	WorkerPort    int32 `json:"workerPort,omitempty"`
	NannyPort     int32 `json:"nannyPort,omitempty"`
}

// DaskClusterStatus defines the observed state of DaskCluster
type DaskClusterStatus struct {
	// TODO: add worker scaling fields (replicas/selector)
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=dask
//+kubebuilder:subresource:status

// DaskCluster is the Schema for the daskclusters API.
type DaskCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DaskClusterSpec   `json:"spec,omitempty"`
	Status DaskClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DaskClusterList contains a list of DaskCluster.
type DaskClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DaskCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DaskCluster{}, &DaskClusterList{})
}

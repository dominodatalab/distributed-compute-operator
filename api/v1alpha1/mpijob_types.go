package v1alpha1

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MPIJobLauncher defines launcher-specific workload settings.
type MPIJobLauncher struct {
	WorkloadConfig `json:",inline"`
	Command        []string `json:"command"`
}

// MPIJobWorker defines worker-specific workload settings.
type MPIJobWorker struct {
	WorkloadConfig `json:",inline"`
	Slots          *int32 `json:"slots,omitempty"`
	Replicas       *int32 `json:"replicas,omitempty"`
}

// MPIJobSpec defines the desired state of MPIJob.
type MPIJobSpec struct {
	ClusterConfig `json:",inline"`

	Launcher MPIJobLauncher `json:"launcher,omitempty"`
	Worker   MPIJobWorker   `json:"worker,omitempty"`
}

// MPIJobStatus defines the observed state of MPIJob.
type MPIJobStatus struct {
	LauncherStatus batchv1.JobConditionType `json:"launcherStatus"`
	StartTime      *metav1.Time             `json:"startTime,omitempty"`
	CompletionTime *metav1.Time             `json:"completionTime,omitempty"`

	// Image is the canonical reference url to the cluster container image.
	Image string `json:"image,omitempty"`
	// Nodes are pods that comprise the cluster.
	Nodes []string `json:"nodes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=mpi
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Workers",type=integer,JSONPath=".spec.worker.replicas"
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=".status.launcherStatus"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Image",type=string,JSONPath=".status.image",priority=10
//+kubebuilder:printcolumn:name="Bound PSP",type=string,JSONPath=".spec.podSecurityPolicy",priority=10
//+kubebuilder:printcolumn:name="Network Policy",type=boolean,JSONPath=".spec.networkPolicy.enabled",priority=10
//+kubebuilder:printcolumn:name="Pods",type=string,JSONPath=".status.nodes",priority=10

// MPIJob is the Schema for the mpijobs API.
type MPIJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MPIJobSpec   `json:"spec,omitempty"`
	Status MPIJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MPIJobList contains a list of MPIJob.
type MPIJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MPIJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MPIJob{}, &MPIJobList{})
}

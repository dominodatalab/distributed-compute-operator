package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MPIJobSpec defines the desired state of MPIJob
type MPIJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of MPIJob. Edit mpijob_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// MPIJobStatus defines the observed state of MPIJob
type MPIJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MPIJob is the Schema for the mpijobs API
type MPIJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MPIJobSpec   `json:"spec,omitempty"`
	Status MPIJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MPIJobList contains a list of MPIJob
type MPIJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MPIJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MPIJob{}, &MPIJobList{})
}

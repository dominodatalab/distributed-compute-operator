package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DaskClusterWorker struct {
	PodConfig `json:",inline"`

	Replicas int32 `json:"replicas,omitempty"`
}

// DaskClusterSpec defines the desired state of DaskCluster.
type DaskClusterSpec struct {
	Image              *OCIImageDefinition           `json:"image,omitempty"`
	ImagePullSecrets   []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	PodSecurityContext *corev1.PodSecurityContext    `json:"podSecurityContext,omitempty"`
	ServiceAccount     ServiceAccountConfig          `json:"serviceAccount,omitempty"`
	EnvVars            []corev1.EnvVar               `json:"envVars,omitempty"`

	SchedulerPort int32 `json:"schedulerPort,omitempty"` // 8786
	DashboardPort int32 `json:"dashboardPort,omitempty"` // 8787
	WorkerPort    int32 `json:"workerPort,omitempty"`    // 3000
	NannyPort     int32 `json:"nannyPort,omitempty"`     // 4000

	Scheduler PodConfig         `json:"scheduler,omitempty"`
	Worker    DaskClusterWorker `json:"worker,omitempty"`
}

// DaskClusterStatus defines the observed state of DaskCluster
type DaskClusterStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=dask
//+kubebuilder:subresource:status

// DaskCluster is the Schema for the daskclusters API
type DaskCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DaskClusterSpec   `json:"spec,omitempty"`
	Status DaskClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DaskClusterList contains a list of DaskCluster
type DaskClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DaskCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DaskCluster{}, &DaskClusterList{})
}

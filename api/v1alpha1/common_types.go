package v1alpha1

import corev1 "k8s.io/api/core/v1"

// Autoscaling configuration for scalable workloads.
type Autoscaling struct {
	// MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.
	// This value must be greater than zero and less than the MaxReplicas.
	//+kubebuilder:validation:Optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// MaxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// This value cannot be less than the replica count for your workload.
	MaxReplicas int32 `json:"maxReplicas"`

	// AverageUtilization is the target value of the average of the resource metric across all relevant pods.
	// This is represented as a percentage of the requested value of the resource for the pods.
	AverageUtilization int32 `json:"averageUtilization"`

	// ScaleDownStabilizationWindowSeconds is the number of seconds for which past recommendations should be considered
	// when scaling down. A shorter window will trigger scale down events quicker, but too short a window may cause
	// replica flapping when metrics used for scaling keep fluctuating.
	//+kubebuilder:validation:Optional
	ScaleDownStabilizationWindowSeconds *int32 `json:"scaleDownStabilizationWindowSeconds,omitempty"`
}

// OCIImageDefinition describes where and when to fetch a container image.
type OCIImageDefinition struct {
	// Registry where the container image is hosted.
	//+kubebuilder:validation:Optional
	Registry string `json:"registry"`

	// Repository where the container image is stored.
	Repository string `json:"repository"`

	// Tag points to a specific container image variant.
	Tag string `json:"tag"`

	// PullPolicy used to fetch container image.
	//+kubebuilder:validation:Optional
	PullPolicy corev1.PullPolicy `json:"pullPolicy"`
}

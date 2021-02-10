package v1alpha1

import corev1 "k8s.io/api/core/v1"

// OCIImageDefinition describes where and when to fetch a container image.
type OCIImageDefinition struct {
	// Registry where the container image is hosted.
	Registry string `json:"registry,omitempty"`

	// Repository where the container image is stored.
	Repository string `json:"repository"`

	// Tag points to a specific container image variant.
	Tag string `json:"tag"`

	// PullPolicy used to fetch container image.
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

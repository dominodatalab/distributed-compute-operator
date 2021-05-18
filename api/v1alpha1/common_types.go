package v1alpha1

import corev1 "k8s.io/api/core/v1"

// Autoscaling configuration for scalable workloads.
type Autoscaling struct {
	// MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.
	// This value must be greater than zero and less than the MaxReplicas.
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// MaxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// This value cannot be less than the replica count for your workload.
	MaxReplicas int32 `json:"maxReplicas"`

	// AverageCPUUtilization is the target value of the average of the resource metric across all relevant pods.
	// This is represented as a percentage of the requested value of the resource for the pods.
	AverageCPUUtilization *int32 `json:"averageCPUUtilization,omitempty"`

	// ScaleDownStabilizationWindowSeconds is the number of seconds for which past recommendations should be considered
	// when scaling down. A shorter window will trigger scale down events quicker, but too short a window may cause
	// replica flapping when metrics used for scaling keep fluctuating.
	ScaleDownStabilizationWindowSeconds *int32 `json:"scaleDownStabilizationWindowSeconds,omitempty"`
}

// IstioConfig defines operator configuration parameters.
type IstioConfig struct {
	// MutualTLSMode will be used to create a workload-specific peer
	// authentication policy that takes precedence over a global and/or
	// namespace-wide policy.
	MutualTLSMode string `json:"istioMutualTLSMode,omitempty"`
}

// OCIImageDefinition describes where and when to fetch a container image.
type OCIImageDefinition struct {
	// Registry where the container image is hosted.
	Registry string `json:"registry,omitempty"`

	// Repository where the container image is stored.
	Repository string `json:"repository,omitempty"`

	// Tag points to a specific container image variant.
	Tag string `json:"tag,omitempty"`

	// PullPolicy used to fetch container image.
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

// PersistentVolumeClaimTemplate describes a claim that pods are allowed to
// reference. These can either pre-exist or leverage storage classes to provide
// dynamic provisioning.
type PersistentVolumeClaimTemplate struct {
	// Name is the unique metadata ID of the volume claim.
	Name string `json:"name"`

	// Spec describes the storage attributes of the underlying claim.
	Spec corev1.PersistentVolumeClaimSpec `json:"spec"`
}

type PodConfig struct {
	Labels               map[string]string               `json:"labels,omitempty"`
	Annotations          map[string]string               `json:"annotations,omitempty"`
	NodeSelector         map[string]string               `json:"nodeSelector,omitempty"`
	Affinity             *corev1.Affinity                `json:"affinity,omitempty"`
	Tolerations          []corev1.Toleration             `json:"tolerations,omitempty"`
	InitContainers       []corev1.Container              `json:"initContainers,omitempty"`
	Volumes              []corev1.Volume                 `json:"volumes,omitempty"`
	VolumeMounts         []corev1.VolumeMount            `json:"volumeMounts,omitempty"`
	VolumeClaimTemplates []PersistentVolumeClaimTemplate `json:"volumeClaimTemplates,omitempty"`
	Resources            corev1.ResourceRequirements     `json:"resources,omitempty"`
}

type ServiceAccountConfig struct {
	Name                         string `json:"name,omitempty"`
	AutomountServiceAccountToken bool   `json:"automountServiceAccountToken,omitempty"`
}

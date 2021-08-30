package v1alpha1

import corev1 "k8s.io/api/core/v1"

// Autoscaling configuration for scalable workloads.
type Autoscaling struct {
	// MinReplicas is the lower limit for the number of replicas to which the
	// autoscaler can scale down. This value must be greater than zero and less
	// than the MaxReplicas.
	MinReplicas *int32 `json:"minReplicas,omitempty"`
	// MaxReplicas is the upper limit for the number of replicas to which the
	// autoscaler can scale up. This value cannot be less than the replica
	// count for your workload.
	MaxReplicas int32 `json:"maxReplicas"`
	// AverageCPUUtilization is the target value of the average of the resource
	// cpu metric across all relevant pods. This is represented as a percentage
	// of the requested value of the resource for the pods.
	AverageCPUUtilization *int32 `json:"averageCPUUtilization,omitempty"`
	// AverageMemoryUtilization is the target value of the average of the
	// resource memory metric across all relevant pods. This is represented as
	// a percentage of the requested value of the resource for the pods.
	AverageMemoryUtilization *int32 `json:"averageMemoryUtilization,omitempty"`
	// ScaleDownStabilizationWindowSeconds is the number of seconds for which
	// past recommendations should be considered when scaling down. A shorter
	// window will trigger scale down events quicker, but too short a window
	// may cause replica flapping when metrics used for scaling keep fluctuating.
	ScaleDownStabilizationWindowSeconds *int32 `json:"scaleDownStabilizationWindowSeconds,omitempty"`
}

// IstioConfig defines Istio configuration parameters.
type IstioConfig struct {
	// MutualTLSMode will be used to create a workload-specific peer
	// authentication policy that takes precedence over a global and/or
	// namespace-wide policy.
	MutualTLSMode string `json:"istioMutualTLSMode,omitempty"`
}

// NetworkPolicyConfig defines network policy configuration options.
type NetworkPolicyConfig struct {
	// Enabled controls the creation of network policies that limit and provide
	// ingress access to the cluster nodes.
	Enabled *bool `json:"enabled,omitempty"`
	// ClientLabels defines the pod selector clause that grants ingress access
	// to the cluster client port(s).
	ClientLabels map[string]string `json:"clientLabels,omitempty"`
	// DashboardLabels defines the pod selector clause that grants ingress
	// access to the cluster dashboard port.
	DashboardLabels map[string]string `json:"dashboardLabels,omitempty"`
}

// OCIImageDefinition describes where and how to fetch a container image.
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

// ServiceAccountConfig defines service account configuration parameters.
type ServiceAccountConfig struct {
	// Name of an existing service account used by cluster workloads. This
	// field will disable the creation of a dedicated cluster service account.
	Name string `json:"name,omitempty"`
	// AutomountServiceAccountToken into workload pods. This field is only used
	// when creating a dedicted cluster service account.
	AutomountServiceAccountToken bool `json:"automountServiceAccountToken,omitempty"`
}

// KerberosKeytabConfig defines kerberos key table configuration options.
type KerberosKeytabConfig struct {
	// Contents are the binary data stored in the keytab file.
	Contents []byte `json:"contents,omitempty"`
	// MountPath where the keytab file should be created inside a pod.
	MountPath string `json:"mountPath,omitempty"`
}

// ClusterConfig defines high-level cluster options.
type ClusterConfig struct {
	// IstioConfig overrides for a cluster.
	IstioConfig `json:",inline"`
	// GlobalLabels applied to all resources in addition to stock labels.
	GlobalLabels map[string]string `json:"globalLabels,omitempty"`
	// Image used to launch cluster nodes.
	Image *OCIImageDefinition `json:"image,omitempty"`
	// Autoscaling parameters used to scale up/down cluster nodes.
	Autoscaling *Autoscaling `json:"autoscaling,omitempty"`
	// NetworkPolicy parameters used to IP traffic flow.
	NetworkPolicy NetworkPolicyConfig `json:"networkPolicy,omitempty"`
	// ServiceAccount parameters used to override default behavior.
	ServiceAccount ServiceAccountConfig `json:"serviceAccount,omitempty"`
	// KerberosKeytab parameters used to add kerberos authentication.
	KerberosKeytab *KerberosKeytabConfig `json:"kerberosKeytab,omitempty"`
	// ImagePullSecrets are references to secrets with pull credentials to
	// private registries where cluster image are stored.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// PodSecurityContext added to every cluster pod.
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`
	// EnvVars added to all every cluster container.
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`
	// PodSecurityPolicy name can be provided to restrict and/or provide
	// execution permissions to processes running within cluster pods.
	PodSecurityPolicy string `json:"podSecurityPolicy,omitempty"`
}

// WorkloadConfig defines options common to all cluster nodes.
type WorkloadConfig struct {
	// Labels applied to cluster pods in addition to stock labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations applied to cluster pods.
	Annotations map[string]string `json:"annotations,omitempty"`
	// NodeSelector applied to cluster pods.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Resources are the requests and limits applied to cluster containers.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Affinity applied to cluster pods.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Tolerations applied to cluster pods.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// InitContainers added to cluster pods.
	InitContainers []corev1.Container `json:"initContainers,omitempty"`
	// Volumes added to cluster pods.
	Volumes []corev1.Volume `json:"volumes,omitempty"`
	// VolumeMounts added to cluster containers.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// VolumeClaimTemplates is a list of claims that cluster pods are allowed
	// to reference. You can enable dynamic provisioning of additional storage
	// on-demand by using a storage class provisioner.
	VolumeClaimTemplates []PersistentVolumeClaimTemplate `json:"volumeClaimTemplates,omitempty"`
}

// ClusterStatusConfig defines the observed state of a given cluster. The
// controllers will generate and populate these fields during reconciliation.
type ClusterStatusConfig struct {
	// Image is the canonical reference url to the cluster container image.
	Image string `json:"image,omitempty"`
	// Nodes are pods that comprise the cluster.
	Nodes []string `json:"nodes,omitempty"`
	// WorkerReplicas is the `scale.status.replicas` subresource field.
	WorkerReplicas int32 `json:"workerReplicas,omitempty"`
	// WorkerSelector is the `scale.status.selector` subresource field.
	WorkerSelector string `json:"workerSelector,omitempty"`
}

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SparkClusterNode defines attributes common to all spark node types.
type SparkClusterNode struct {
	// Labels applied to spark pods in addition to stock labels.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations applied to spark pods.
	Annotations map[string]string `json:"annotations,omitempty"`

	// NodeSelector applied to spark pods.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Affinity applied to spark pods.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Tolerations applied to spark pods.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// InitContainers added to spark pods.
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// Volumes added to spark pods.
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// VolumeMounts added to spark containers.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// Resources are the requests and limits applied to spark containers.
	Resources corev1.ResourceRequirements `json:"resources"`

	// VolumeClaimTemplates is a list of claims that Spark pods are allowed to
	// reference. You can enable dynamic provisioning of additional storage
	// on-demand by using a storage class provisioner.
	VolumeClaimTemplates []PersistentVolumeClaimTemplate `json:"volumeClaimTemplates,omitempty"`

	// FrameworkConfig is the extra framework-specific configuration for this cluster
	// For spark this means we'll generate a spark-defaults.conf configmap
	// and mount it at the requested location
	FrameworkConfig *FrameworkConfig `json:"frameworkConfig,omitempty"`
}

type FrameworkConfig struct {
	// Configs includes the configuration values to include in the configmap
	Configs map[string]string `json:"configs"`
}

// SparkClusterMaster defines master-specific pod settings.
type SparkClusterMaster struct {
	SparkClusterNode `json:",inline"`
}

// SparkClusterWorker defines worker-specific pod settings.
type SparkClusterWorker struct {
	SparkClusterNode `json:",inline"`

	// Replicas configures the total number of workers in the cluster.
	// This field behaves differently when Autoscaling is enabled. If Autoscaling.MinReplicas is unspecified, then the
	// minimum number of replicas will be set to this value. Additionally, you can specify an "initial cluster size" by
	// setting this field to some value above the minimum number of replicas.
	Replicas *int32 `json:"replicas,omitempty"`

	// WorkerMemoryLimit configures the SPARK_WORKER_MEMORY envVar
	WorkerMemoryLimit string `json:"workerMemoryLimit,omitempty"`
}

// SparkClusterSpec defines the desired state of a SparkCluster resource.
type SparkClusterSpec struct {
	// Image used to launch master and worker nodes.
	Image *OCIImageDefinition `json:"image,omitempty"`

	// ImagePullSecrets are references to secrets with credentials to private registries used to pull images.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	GlobalLabels map[string]string `json:"globalLabels,omitempty"`

	EnvoyFilterLabels map[string]string `json:"envoyFilterLabels,omitempty"`

	// ClusterPort is the port on which the spark protocol communicates
	ClusterPort int32 `json:"clusterPort,omitempty"`

	// These two are meant for istio compatibility on spark
	TCPMasterWebPort int32 `json:"tcpMasterWebPort,omitempty"`
	TCPWorkerWebPort int32 `json:"tcpWorkerWebPort,omitempty"`

	// IstioConfig parameters for Spark clusters.
	IstioConfig `json:",inline"`

	// Driver configures the SparkCluster to communicate with the Spark Driver
	Driver SparkClusterDriver `json:"sparkClusterDriver,omitempty"`

	// DashboardPort is the port used by the dashboard server.
	DashboardPort int32 `json:"dashboardPort,omitempty"`

	// EnableDashboard starts the dashboard web UI.
	EnableDashboard *bool `json:"enableDashboard,omitempty"`

	// NetworkPolicy will create a pod selector clause for each set of labels.
	// This is used to grant ingress access to one or more groups of external pods and is
	// only applicable when EnableNetworkPolicy is true.
	NetworkPolicy SparkClusterNetworkPolicy `json:"networkPolicy,omitempty"`

	// PodSecurityPolicy name can be provided to govern execution of the spark processes within pods.
	PodSecurityPolicy string `json:"podSecurityPolicy,omitempty"`

	// PodSecurityContext added to every spark pod.
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ServiceAccountName will disable the creation of a dedicated cluster service account.
	// The service account referenced by the provided name will be used instead.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// EnvVars added to every spark pod container.
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`

	// Master node configuration parameters.
	Master SparkClusterMaster `json:"master,omitempty"`

	// Worker node configuration parameters.
	Worker SparkClusterWorker `json:"worker,omitempty"`

	// Autoscaling parameters used to scale up/down spark worker nodes.
	Autoscaling *Autoscaling `json:"autoscaling,omitempty"`

	// KerberosKeytab configures the Kerberos Keytab for Spark
	KerberosKeytab *KerberosKeytabConfig `json:"kerberosKeytab,omitempty"`
}

// SparkClusterDriver defines the configuration for the driver service
type SparkClusterDriver struct {
	SparkClusterName           string `json:"sparkClusterName,omitempty"`
	ExecutionName              string `json:"executionName,omitempty"`
	DriverPortName             string `json:"driverPortName,omitempty"`
	DriverPort                 int32  `json:"driverPort,omitempty"`
	DriverBlockManagerPortName string `json:"driverBlockManagerPortName,omitempty"`
	DriverBlockManagerPort     int32  `json:"driverBlockManagerPort,omitempty"`
	DriverUIPortName           string `json:"driverUIPortName,omitempty"`
	DriverUIPort               int32  `json:"driverUIPort,omitempty"`
}

// SparkClusterNetworkPolicy defines network policy configuration options.
type SparkClusterNetworkPolicy struct {
	// Enabled controls the creation of network policies that limit and provide
	// ingress access to the cluster nodes.
	Enabled *bool `json:"enabled,omitempty"`

	// ClientServerLabels defines the pod selector clause that grant ingress
	// access to the master client server port.
	ClientServerLabels map[string]string `json:"clientServerLabels,omitempty"`

	// DashboardLabels defines the pod selector clause used to grant ingress
	// access to the master dashboard port.
	DashboardLabels map[string]string `json:"dashboardLabels,omitempty"`

	// ExternalPolicyEnabled controls creation of network policies which deal with two way traffic to pods that are
	// external to the cluster entirely. The spark driver, for example, is generally going to live outside the
	// cluster
	ExternalPolicyEnabled *bool `json:"externalPolicyEnabled,omitempty"`

	// ExternalPodLabels defines the pod selector clause for used to granted unfettered
	// access to cluster resources.
	ExternalPodLabels map[string]string `json:"clusterLabels,omitempty"`
}

// SparkClusterStatus defines the observed state of a SparkCluster resource.
type SparkClusterStatus struct {
	// Nodes that comprise the cluster.
	Nodes []string `json:"nodes,omitempty"`

	// WorkerReplicas is the scale.status.replicas subresource field.
	WorkerReplicas int32 `json:"workerReplicas,omitempty"`
	// WorkerSelector is the scale.status.selector subresource field.
	WorkerSelector string `json:"workerSelector,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=spark
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.worker.replicas,statuspath=.status.workerReplicas,selectorpath=.status.workerSelector
//+kubebuilder:printcolumn:name="Worker Replicas",type=integer,JSONPath=".spec.worker.replicas"
//+kubebuilder:printcolumn:name="Image",type=string,JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

// SparkCluster is the Schema for the sparkclusters API.
type SparkCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SparkClusterSpec   `json:"spec,omitempty"`
	Status SparkClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SparkClusterList contains a list of SparkCluster resources.
type SparkClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SparkCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SparkCluster{}, &SparkClusterList{})
}

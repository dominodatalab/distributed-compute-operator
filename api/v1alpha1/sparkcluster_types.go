package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SparkClusterNode defines attributes common to all spark node types.
type SparkClusterNode struct {
	WorkloadConfig `json:",inline"`

	// DefaultConfiguration can be used to tune the execution environment for
	// your Spark applications. The values provided will be used to construct
	// the spark-defaults.conf file.
	DefaultConfiguration map[string]string `json:"defaultConfiguration,omitempty"`
}

// SparkClusterWorker defines worker-specific pod settings.
type SparkClusterWorker struct {
	SparkClusterNode `json:",inline"`

	// Replicas configures the total number of workers in the cluster. This
	// field behaves differently when Autoscaling is enabled. If
	// Autoscaling.MinReplicas is unspecified, then the minimum number of
	// replicas will be set to this value. Additionally, you can specify an
	// "initial cluster size" by setting this field to some value above the
	// minimum number of replicas.
	Replicas *int32 `json:"replicas,omitempty"`

	// Obsolete value of WorkerMemoryLimit used in previous
	// versions; used here only to check compatibility of CRDs.
	ObsoleteWorkerMemoryLimit string `json:"workerMemoryLimit,omitempty"`
}

// SparkClusterDriver defines the configuration for the external driver.
type SparkClusterDriver struct {
	// Port used for communication by the driver.
	Port int32 `json:"port,omitempty"`
	// UIPort used by the driver.
	UIPort int32 `json:"uiPort,omitempty"`
	// BlockManagerPort used by the driver.
	BlockManagerPort int32 `json:"blockManagerPort,omitempty"`
	// Selector labels for driver pod(s).
	Selector map[string]string `json:"selector,omitempty"`
}

// SparkClusterSpec defines the desired state of a SparkCluster resource.
type SparkClusterSpec struct {
	ClusterConfig `json:",inline"`

	// Master node configuration parameters.
	Master SparkClusterNode `json:"master,omitempty"`
	// Worker node configuration parameters.
	Worker SparkClusterWorker `json:"worker,omitempty"`
	// Driver configures the SparkCluster to communicate with the Spark Driver.
	Driver SparkClusterDriver `json:"driver,omitempty"`

	// EnvoyFilterLabels are specific labels that must already exist on the
	// spark-driver so that users can set idle_timeout properly using the
	// EnvoyFilter resource.
	EnvoyFilterLabels map[string]string `json:"envoyFilterLabels,omitempty"`
	// WorkerMemoryLimit configures the SPARK_WORKER_MEMORY envVar.
	WorkerMemoryLimit string `json:"workerMemoryLimit,omitempty"`
	// ClusterPort is the port used for master/worker/driver communication.
	ClusterPort int32 `json:"clusterPort,omitempty"`
	// MasterWebPort is the port for the master web UI.
	MasterWebPort int32 `json:"masterWebPort,omitempty"`
	// WorkerWebPort is the port for the worker web UI.
	WorkerWebPort int32 `json:"workerWebPort,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=spark
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.worker.replicas,statuspath=.status.workerReplicas,selectorpath=.status.workerSelector
//+kubebuilder:printcolumn:name="Workers",type=integer,JSONPath=".spec.worker.replicas"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Image",type=string,JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Network Policy",type=boolean,JSONPath=".spec.networkPolicy.enabled",priority=10
//+kubebuilder:printcolumn:name="Pods",type=string,JSONPath=".status.nodes",priority=10

// SparkCluster is the Schema for the sparkclusters API.
type SparkCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SparkClusterSpec    `json:"spec,omitempty"`
	Status ClusterStatusConfig `json:"status,omitempty"`
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

// IsIncompatibleVersion checks if the provided instance of SparkCluster struct
// has a version compatible with the current one.
func (sc *SparkCluster) IsIncompatibleVersion() bool {
	// Unfortunately we can't rely on e.g. Spec.TypeMeta.APIVersion, because of
	// past breaking changes for which the meta information hasn't been updated.
	// We're checking on a field that has been mandatory in the old version,
	// but is currently removed.
	fmt.Printf(">>> >>> >>> %v\n", sc.Spec.Worker.ObsoleteWorkerMemoryLimit)
	return sc.Spec.Worker.ObsoleteWorkerMemoryLimit != ""
}

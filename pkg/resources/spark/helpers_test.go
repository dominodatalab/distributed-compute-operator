package spark

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// sparkClusterFixture should be used for all spark unit testing.
func sparkClusterFixture() *dcv1alpha1.SparkCluster {
	return &dcv1alpha1.SparkCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SparkCluster",
			APIVersion: "distributed-compute.dominodatalab.com/v1test1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id",
			Namespace: "fake-ns",
		},
		Spec: dcv1alpha1.SparkClusterSpec{
			Image: &dcv1alpha1.OCIImageDefinition{
				Registry:   "fake-reg",
				Repository: "fake-repo",
				Tag:        "fake-tag",
				PullPolicy: v1.PullIfNotPresent,
			},
			ClusterPort:          7077,
			TCPMasterWebPort:     8080,
			TCPWorkerWebPort:     8081,
			DashboardPort:        8080,
			DashboardServicePort: 80,
			Worker: dcv1alpha1.SparkClusterWorker{
				Replicas:          pointer.Int32Ptr(5),
				WorkerMemoryLimit: "4505m",
			},
		},
	}
}

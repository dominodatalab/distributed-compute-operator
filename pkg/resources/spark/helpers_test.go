package spark

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

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
			ScalableClusterConfig: dcv1alpha1.ScalableClusterConfig{
				ClusterConfig: dcv1alpha1.ClusterConfig{
					Image: &dcv1alpha1.OCIImageDefinition{
						Registry:   "fake-reg",
						Repository: "fake-repo",
						Tag:        "fake-tag",
						PullPolicy: corev1.PullIfNotPresent,
					},
				},
			},
			ClusterPort:       7077,
			MasterWebPort:     8080,
			WorkerWebPort:     8081,
			WorkerMemoryLimit: "4505m",
			Worker: dcv1alpha1.SparkClusterWorker{
				Replicas: ptr.To(int32(5)),
			},
		},
	}
}

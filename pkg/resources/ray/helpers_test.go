package ray

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// rayClusterFixture should be used for all ray unit testing.
func rayClusterFixture() *dcv1alpha1.RayCluster {
	return &dcv1alpha1.RayCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RayCluster",
			APIVersion: "distributed-compute.dominodatalab.com/v1test1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id",
			Namespace: "fake-ns",
		},
		Spec: dcv1alpha1.RayClusterSpec{
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
			Port: 6379,
			RedisShardPorts: []int32{
				6380,
				6381,
			},
			ClientServerPort:  10001,
			ObjectManagerPort: 2384,
			NodeManagerPort:   2385,
			GCSServerPort:     2386,
			WorkerPorts:       []int32{11000, 11001},
			DashboardPort:     8265,
			Worker: dcv1alpha1.RayClusterWorker{
				Replicas: pointer.Int32Ptr(5),
			},
		},
	}
}

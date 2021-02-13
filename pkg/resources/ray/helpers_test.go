package ray

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func rayClusterFixture() *dcv1alpha1.RayCluster {
	return &dcv1alpha1.RayCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id",
			Namespace: "fake-ns",
		},
		Spec: dcv1alpha1.RayClusterSpec{
			Image: &dcv1alpha1.OCIImageDefinition{
				Registry:   "fake-reg",
				Repository: "fake-repo",
				Tag:        "fake-tag",
				PullPolicy: v1.PullIfNotPresent,
			},
			WorkerReplicaCount: 5,
			Port:               6379,
			RedisShardPorts: []int32{
				6380,
				6381,
			},
			ClientServerPort:  10001,
			ObjectManagerPort: 2384,
			NodeManagerPort:   2385,
			DashboardPort:     8265,
		},
	}
}

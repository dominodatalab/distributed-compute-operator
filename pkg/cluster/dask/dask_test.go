package dask

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func testDaskCluster() *dcv1alpha1.DaskCluster {
	return &dcv1alpha1.DaskCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "ns",
		},
		Spec: dcv1alpha1.DaskClusterSpec{
			ClusterConfig: dcv1alpha1.ClusterConfig{
				Image: &dcv1alpha1.OCIImageDefinition{
					Registry:   "",
					Repository: "daskdev/dask",
					Tag:        "test-tag",
				},
				NetworkPolicy: dcv1alpha1.NetworkPolicyConfig{
					ClientLabels: map[string]string{
						"test-client": "true",
					},
					DashboardLabels: map[string]string{
						"test-ui-client": "true",
					},
				},
				PodSecurityPolicy: "privileged",
			},
			Scheduler:     dcv1alpha1.WorkloadConfig{},
			Worker:        dcv1alpha1.DaskClusterWorker{},
			SchedulerPort: 8786,
			DashboardPort: 8787,
			WorkerPort:    3000,
			NannyPort:     3001,
		},
	}
}

package dask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestHorizontalPodAutoscaler(t *testing.T) {
	dc := testDaskCluster()

	t.Run("basic", func(t *testing.T) {
		ds := horizontalPodAutoscalerDS{dc: dc}
		dc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{}

		actual := ds.HorizontalPodAutoscaler()
		expected := &autoscalingv2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-dask",
				Namespace: "ns",
				Labels: map[string]string{
					"app.kubernetes.io/instance":   "test",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
					"app.kubernetes.io/name":       "dask",
					"app.kubernetes.io/version":    "test-tag",
				},
			},
			Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
					Kind:       "DaskCluster",
					Name:       "test",
					APIVersion: "distributed-compute.dominodatalab.com/v1test1",
				},
				MinReplicas: nil,
				MaxReplicas: 0,
				Metrics:     nil,
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("min_replicas", func(t *testing.T) {
		expected := pointer.Int32Ptr(7)
		dc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			MinReplicas: expected,
		}

		assert.Equal(t, expected, dc.Spec.Autoscaling.MinReplicas)
	})

	t.Run("max_replicas", func(t *testing.T) {
		var expected int32 = 10
		dc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			MaxReplicas: expected,
		}

		assert.Equal(t, expected, dc.Spec.Autoscaling.MaxReplicas)
	})

	t.Run("avg_cpu_util", func(t *testing.T) {
		dc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			AverageCPUUtilization: pointer.Int32Ptr(75),
		}

		ds := horizontalPodAutoscalerDS{dc: dc}
		hpa := ds.HorizontalPodAutoscaler()

		expected := []autoscalingv2beta2.MetricSpec{
			{
				Type: "Resource",
				Resource: &autoscalingv2beta2.ResourceMetricSource{
					Name: "cpu",
					Target: autoscalingv2beta2.MetricTarget{
						Type:               "Utilization",
						AverageUtilization: pointer.Int32Ptr(75),
					},
				},
			},
		}
		assert.Equal(t, expected, hpa.Spec.Metrics)
	})

	t.Run("avg_memory_util", func(t *testing.T) {
		dc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			AverageMemoryUtilization: pointer.Int32Ptr(75),
		}

		ds := horizontalPodAutoscalerDS{dc: dc}
		hpa := ds.HorizontalPodAutoscaler()

		expected := []autoscalingv2beta2.MetricSpec{
			{
				Type: "Resource",
				Resource: &autoscalingv2beta2.ResourceMetricSource{
					Name: "memory",
					Target: autoscalingv2beta2.MetricTarget{
						Type:               "Utilization",
						AverageUtilization: pointer.Int32Ptr(75),
					},
				},
			},
		}
		assert.Equal(t, expected, hpa.Spec.Metrics)
	})

	t.Run("scale_down_behavior", func(t *testing.T) {
		dc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			ScaleDownStabilizationWindowSeconds: pointer.Int32Ptr(60),
		}

		ds := horizontalPodAutoscalerDS{dc: dc}
		hpa := ds.HorizontalPodAutoscaler()

		expected := &autoscalingv2beta2.HorizontalPodAutoscalerBehavior{
			ScaleDown: &autoscalingv2beta2.HPAScalingRules{
				StabilizationWindowSeconds: pointer.Int32Ptr(60),
			},
		}
		assert.Equal(t, expected, hpa.Spec.Behavior)
	})
}

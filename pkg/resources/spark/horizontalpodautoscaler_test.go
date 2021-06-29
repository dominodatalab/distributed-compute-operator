package spark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestNewHorizontalPodAutoscaler(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{}
		actual, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		expected := &autoscalingv2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-id",
				Namespace: "fake-ns",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "spark",
					"app.kubernetes.io/instance":   "test-id",
					"app.kubernetes.io/version":    "fake-tag",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
				},
			},
			Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
					Kind:       "SparkCluster",
					Name:       "test-id",
					APIVersion: "distributed-compute.dominodatalab.com/v1test1",
				},
				MinReplicas: nil,
				MaxReplicas: 0,
				Metrics: []autoscalingv2beta2.MetricSpec{
					{
						Type: "Resource",
						Resource: &autoscalingv2beta2.ResourceMetricSource{
							Name: "cpu",
							Target: autoscalingv2beta2.MetricTarget{
								Type:               "Utilization",
								AverageUtilization: nil,
							},
						},
					},
				},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("min_replicas", func(t *testing.T) {
		rc := sparkClusterFixture()

		expected := pointer.Int32Ptr(7)
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			MinReplicas: expected,
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		assert.Equal(t, expected, hpa.Spec.MinReplicas)
	})

	t.Run("max_replicas", func(t *testing.T) {
		rc := sparkClusterFixture()

		var expected int32 = 10
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			MaxReplicas: expected,
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		assert.Equal(t, expected, hpa.Spec.MaxReplicas)
	})

	t.Run("avg_cpu_util", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			AverageCPUUtilization: pointer.Int32Ptr(75),
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

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

	t.Run("scale_down_behavior", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			ScaleDownStabilizationWindowSeconds: pointer.Int32Ptr(60),
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		expected := &autoscalingv2beta2.HorizontalPodAutoscalerBehavior{
			ScaleDown: &autoscalingv2beta2.HPAScalingRules{
				StabilizationWindowSeconds: pointer.Int32Ptr(60),
			},
		}
		assert.Equal(t, expected, hpa.Spec.Behavior)
	})

	t.Run("error", func(t *testing.T) {
		rc := sparkClusterFixture()
		_, err := NewHorizontalPodAutoscaler(rc)

		assert.Error(t, err)
	})
}

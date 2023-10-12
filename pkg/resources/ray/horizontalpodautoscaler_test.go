package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestNewHorizontalPodAutoscaler(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{}
		actual, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		expected := &autoscalingv2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-id-ray",
				Namespace: "fake-ns",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "ray",
					"app.kubernetes.io/instance":   "test-id",
					"app.kubernetes.io/version":    "fake-tag",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
				},
			},
			Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
					Kind:       "RayCluster",
					Name:       "test-id",
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
		rc := rayClusterFixture()

		expected := ptr.To(int32(7))
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			MinReplicas: expected,
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		assert.Equal(t, expected, hpa.Spec.MinReplicas)
	})

	t.Run("max_replicas", func(t *testing.T) {
		rc := rayClusterFixture()

		var expected int32 = 10
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			MaxReplicas: expected,
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		assert.Equal(t, expected, hpa.Spec.MaxReplicas)
	})

	t.Run("avg_cpu_util", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			AverageCPUUtilization: ptr.To(int32(75)),
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		expected := []autoscalingv2.MetricSpec{
			{
				Type: "Resource",
				Resource: &autoscalingv2.ResourceMetricSource{
					Name: "cpu",
					Target: autoscalingv2.MetricTarget{
						Type:               "Utilization",
						AverageUtilization: ptr.To(int32(75)),
					},
				},
			},
		}
		assert.Equal(t, expected, hpa.Spec.Metrics)
	})

	t.Run("avg_memory_util", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			AverageMemoryUtilization: ptr.To(int32(75)),
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		expected := []autoscalingv2.MetricSpec{
			{
				Type: "Resource",
				Resource: &autoscalingv2.ResourceMetricSource{
					Name: "memory",
					Target: autoscalingv2.MetricTarget{
						Type:               "Utilization",
						AverageUtilization: ptr.To(int32(75)),
					},
				},
			},
		}
		assert.Equal(t, expected, hpa.Spec.Metrics)
	})

	t.Run("scale_down_behavior", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			ScaleDownStabilizationWindowSeconds: ptr.To(int32(60)),
		}

		hpa, err := NewHorizontalPodAutoscaler(rc)
		require.NoError(t, err)

		expected := &autoscalingv2.HorizontalPodAutoscalerBehavior{
			ScaleDown: &autoscalingv2.HPAScalingRules{
				StabilizationWindowSeconds: ptr.To(int32(60)),
			},
		}
		assert.Equal(t, expected, hpa.Spec.Behavior)
	})

	t.Run("error", func(t *testing.T) {
		rc := rayClusterFixture()
		_, err := NewHorizontalPodAutoscaler(rc)

		assert.Error(t, err)
	})
}

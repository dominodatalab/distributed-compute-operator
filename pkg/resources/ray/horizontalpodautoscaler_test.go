package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestNewHorizontalPodAutoscaler(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			MaxReplicas:        10,
			AverageUtilization: 50,
		}
		actual := NewHorizontalPodAutoscaler(rc)

		expected := &autoscalingv2beta2.HorizontalPodAutoscaler{
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
			Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
					Kind:       "RayCluster",
					Name:       "test-id",
					APIVersion: "distributed-compute.dominodatalab.com/v1test1",
				},
				MinReplicas: pointer.Int32Ptr(5),
				MaxReplicas: 10,
				Metrics: []autoscalingv2beta2.MetricSpec{
					{
						Type: "Resource",
						Resource: &autoscalingv2beta2.ResourceMetricSource{
							Name: "cpu",
							Target: autoscalingv2beta2.MetricTarget{
								Type:               "Utilization",
								AverageUtilization: pointer.Int32Ptr(50),
							},
						},
					},
				},
			},
		}
		assert.Equal(t, expected, actual)
	})

	t.Run("scale_down_behavior", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			ScaleDownStabilizationWindowSeconds: pointer.Int32Ptr(60),
		}
		hpa := NewHorizontalPodAutoscaler(rc)

		expected := &autoscalingv2beta2.HorizontalPodAutoscalerBehavior{
			ScaleDown: &autoscalingv2beta2.HPAScalingRules{
				StabilizationWindowSeconds: pointer.Int32Ptr(60),
			},
		}
		assert.Equal(t, expected, hpa.Spec.Behavior)
	})

	t.Run("min_replicas", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := pointer.Int32Ptr(7)
		rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
			MinReplicas: expected,
		}
		hpa := NewHorizontalPodAutoscaler(rc)

		assert.Equal(t, expected, hpa.Spec.MinReplicas)
	})
}

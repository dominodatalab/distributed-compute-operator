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
	rc := rayClusterFixture()
	rc.Spec.Autoscaling = &dcv1alpha1.Autoscaling{
		MaxReplicas:                         10,
		AverageUtilization:                  50,
		ScaleDownStabilizationWindowSeconds: pointer.Int32Ptr(60),
	}
	hpa := NewHorizontalPodAutoscaler(rc)

	expected := &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id",
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
				Kind:       "Deployment",
				Name:       "test-id-ray-worker",
				APIVersion: "apps/v1",
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
			Behavior: &autoscalingv2beta2.HorizontalPodAutoscalerBehavior{
				ScaleDown: &autoscalingv2beta2.HPAScalingRules{
					StabilizationWindowSeconds: pointer.Int32Ptr(60),
				},
			},
		},
	}
	assert.Equal(t, expected, hpa)
}

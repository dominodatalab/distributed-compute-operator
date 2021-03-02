package ray

import (
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

const defaultAverageUtilization int32 = 50

// NewHorizontalPodAutoscaler generates an HPA that targets a RayCluster resource.
//
// The metrics-server needs to be launched separately and the worker deployment
// requires cpu resource requests in order for this object to have any effect.
func NewHorizontalPodAutoscaler(rc *dcv1alpha1.RayCluster) *autoscalingv2beta2.HorizontalPodAutoscaler {
	var behavior *autoscalingv2beta2.HorizontalPodAutoscalerBehavior
	minReplicas := pointer.Int32Ptr(rc.Spec.Worker.Replicas)
	maxReplicas := *minReplicas
	avgUtilization := defaultAverageUtilization

	if autoscaling := rc.Spec.Autoscaling; autoscaling != nil {
		maxReplicas = autoscaling.MaxReplicas
		avgUtilization = autoscaling.AverageUtilization

		if autoscaling.MinReplicas != nil {
			minReplicas = rc.Spec.Autoscaling.MinReplicas
		}

		if autoscaling.ScaleDownStabilizationWindowSeconds != nil {
			behavior = &autoscalingv2beta2.HorizontalPodAutoscalerBehavior{
				ScaleDown: &autoscalingv2beta2.HPAScalingRules{
					StabilizationWindowSeconds: autoscaling.ScaleDownStabilizationWindowSeconds,
				},
			}
		}
	}

	return &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, ComponentNone),
			Namespace: rc.Namespace,
			Labels:    MetadataLabels(rc),
		},
		Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				APIVersion: rc.APIVersion,
				Kind:       rc.Kind,
				Name:       rc.Name,
			},
			MinReplicas: minReplicas,
			MaxReplicas: maxReplicas,
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Type: autoscalingv2beta2.ResourceMetricSourceType,
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: corev1.ResourceCPU,
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: pointer.Int32Ptr(avgUtilization),
						},
					},
				},
			},
			Behavior: behavior,
		},
	}
}

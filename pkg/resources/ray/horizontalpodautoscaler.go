package ray

import (
	"fmt"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// NewHorizontalPodAutoscaler generates an HPA that targets a RayCluster resource.
//
// The metrics-server needs to be launched separately and the worker stateful
// set requires cpu resource requests in order for this object to have any
// effect.
func NewHorizontalPodAutoscaler(rc *dcv1alpha1.RayCluster) (*autoscalingv2beta2.HorizontalPodAutoscaler, error) {
	autoscaling := rc.Spec.Autoscaling
	if autoscaling == nil {
		return nil, fmt.Errorf("cannot build HPA without autoscaling config")
	}

	var behavior *autoscalingv2beta2.HorizontalPodAutoscalerBehavior
	if autoscaling.ScaleDownStabilizationWindowSeconds != nil {
		behavior = &autoscalingv2beta2.HorizontalPodAutoscalerBehavior{
			ScaleDown: &autoscalingv2beta2.HPAScalingRules{
				StabilizationWindowSeconds: autoscaling.ScaleDownStabilizationWindowSeconds,
			},
		}
	}

	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: HorizontalPodAutoscalerObjectMeta(rc),
		Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				APIVersion: rc.APIVersion,
				Kind:       rc.Kind,
				Name:       rc.Name,
			},
			MinReplicas: autoscaling.MinReplicas,
			MaxReplicas: autoscaling.MaxReplicas,
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Type: autoscalingv2beta2.ResourceMetricSourceType,
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: corev1.ResourceCPU,
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: autoscaling.AverageCPUUtilization,
						},
					},
				},
			},
			Behavior: behavior,
		},
	}

	return hpa, nil
}

// HorizontalPodAutoscalerObjectMeta returns the ObjectMeta object used to identify new HPA objects.
func HorizontalPodAutoscalerObjectMeta(rc *dcv1alpha1.RayCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      InstanceObjectName(rc.Name, ComponentNone),
		Namespace: rc.Namespace,
		Labels:    MetadataLabels(rc),
	}
}

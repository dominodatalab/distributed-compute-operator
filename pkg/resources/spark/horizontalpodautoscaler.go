package spark

import (
	"fmt"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// NewHorizontalPodAutoscaler generates an HPA that targets a SparkCluster resource.
//
// The metrics-server needs to be launched separately and the worker deployment
// requires cpu resource requests in order for this object to have any effect.
func NewHorizontalPodAutoscaler(sc *dcv1alpha1.SparkCluster) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	autoscaling := sc.Spec.Autoscaling
	if autoscaling == nil {
		return nil, fmt.Errorf("cannot build HPA without autoscaling config")
	}

	var behavior *autoscalingv2.HorizontalPodAutoscalerBehavior
	if autoscaling.ScaleDownStabilizationWindowSeconds != nil {
		behavior = &autoscalingv2.HorizontalPodAutoscalerBehavior{
			ScaleDown: &autoscalingv2.HPAScalingRules{
				StabilizationWindowSeconds: autoscaling.ScaleDownStabilizationWindowSeconds,
			},
		}
	}

	var metrics []autoscalingv2.MetricSpec
	if autoscaling.AverageCPUUtilization != nil {
		metrics = append(metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceCPU,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: autoscaling.AverageCPUUtilization,
				},
			},
		})
	}
	if autoscaling.AverageMemoryUtilization != nil {
		metrics = append(metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceMemory,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: autoscaling.AverageMemoryUtilization,
				},
			},
		})
	}

	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: HorizontalPodAutoscalerObjectMeta(sc),
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: sc.APIVersion,
				Kind:       sc.Kind,
				Name:       sc.Name,
			},
			MinReplicas: autoscaling.MinReplicas,
			MaxReplicas: autoscaling.MaxReplicas,
			Metrics:     metrics,
			Behavior:    behavior,
		},
	}

	return hpa, nil
}

// HorizontalPodAutoscalerObjectMeta returns the ObjectMeta object used to identify new HPA objects.
func HorizontalPodAutoscalerObjectMeta(sc *dcv1alpha1.SparkCluster) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      InstanceObjectName(sc.Name, ComponentNone),
		Namespace: sc.Namespace,
		Labels:    AddGlobalLabels(MetadataLabels(sc), sc.Spec.GlobalLabels),
	}
}

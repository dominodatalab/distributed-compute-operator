package dask

import (
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func HorizontalPodAutoscaler() core.OwnedComponent {
	return components.HorizontalPodAutoscaler(func(obj client.Object) components.HorizontalPodAutoscalerDataSource {
		return &horizontalPodAutoscalerDS{dc: daskCluster(obj)}
	})
}

type horizontalPodAutoscalerDS struct {
	dc *dcv1alpha1.DaskCluster
}

func (s *horizontalPodAutoscalerDS) HorizontalPodAutoscaler() *autoscalingv2beta2.HorizontalPodAutoscaler {
	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(s.dc, metadata.ComponentNone),
			Namespace: s.dc.Namespace,
			Labels:    meta.StandardLabels(s.dc),
		},
	}

	as := s.dc.Spec.Autoscaling
	if as == nil {
		return hpa
	}

	var behavior *autoscalingv2beta2.HorizontalPodAutoscalerBehavior
	if as.ScaleDownStabilizationWindowSeconds != nil {
		behavior = &autoscalingv2beta2.HorizontalPodAutoscalerBehavior{
			ScaleDown: &autoscalingv2beta2.HPAScalingRules{
				StabilizationWindowSeconds: as.ScaleDownStabilizationWindowSeconds,
			},
		}
	}

	var metrics []autoscalingv2beta2.MetricSpec
	if as.AverageCPUUtilization != nil {
		metrics = append(metrics, autoscalingv2beta2.MetricSpec{
			Type: autoscalingv2beta2.ResourceMetricSourceType,
			Resource: &autoscalingv2beta2.ResourceMetricSource{
				Name: corev1.ResourceCPU,
				Target: autoscalingv2beta2.MetricTarget{
					Type:               autoscalingv2beta2.UtilizationMetricType,
					AverageUtilization: s.dc.Spec.Autoscaling.AverageCPUUtilization,
				},
			},
		})
	}
	if as.AverageMemoryUtilization != nil {
		metrics = append(metrics, autoscalingv2beta2.MetricSpec{
			Type: autoscalingv2beta2.ResourceMetricSourceType,
			Resource: &autoscalingv2beta2.ResourceMetricSource{
				Name: corev1.ResourceMemory,
				Target: autoscalingv2beta2.MetricTarget{
					Type:               autoscalingv2beta2.UtilizationMetricType,
					AverageUtilization: s.dc.Spec.Autoscaling.AverageMemoryUtilization,
				},
			},
		})
	}

	hpa.Spec = autoscalingv2beta2.HorizontalPodAutoscalerSpec{
		ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
			Kind:       s.dc.Kind,
			Name:       s.dc.Name,
			APIVersion: s.dc.APIVersion,
		},
		MinReplicas: s.dc.Spec.Autoscaling.MinReplicas,
		MaxReplicas: s.dc.Spec.Autoscaling.MaxReplicas,
		Metrics:     metrics,
		Behavior:    behavior,
	}

	return hpa
}

func (s *horizontalPodAutoscalerDS) Delete() bool {
	return s.dc.Spec.Autoscaling == nil
}

package ray

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

const (
	scaleTargetKind            = "Deployment"
	stabilizationWindowSeconds = 60
)

// NewHorizontalPodAutoscaler generates an HPA that targets the ray worker deployment.
// The Autoscaling config from the spec is used to set max replicas and the target average utilization.
//
// The metrics-server needs to be launched separately and the worker deployment
// requires cpu resource requests in order for this object to have any effect.
func NewHorizontalPodAutoscaler(rc *dcv1alpha1.RayCluster) *autoscalingv2beta2.HorizontalPodAutoscaler {
	return &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rc.Name,
			Namespace: rc.Namespace,
			Labels:    MetadataLabels(rc),
		},
		Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				APIVersion: appsv1.SchemeGroupVersion.String(),
				Kind:       scaleTargetKind,
				Name:       InstanceObjectName(rc.Name, ComponentWorker),
			},
			MinReplicas: pointer.Int32Ptr(rc.Spec.WorkerReplicas),
			MaxReplicas: rc.Spec.Autoscaling.MaxReplicas,
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Type: autoscalingv2beta2.ResourceMetricSourceType,
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: v1.ResourceCPU,
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: pointer.Int32Ptr(rc.Spec.Autoscaling.AverageUtilization),
						},
					},
				},
			},
			Behavior: &autoscalingv2beta2.HorizontalPodAutoscalerBehavior{
				ScaleDown: &autoscalingv2beta2.HPAScalingRules{
					StabilizationWindowSeconds: pointer.Int32Ptr(stabilizationWindowSeconds),
				},
			},
		},
	}
}

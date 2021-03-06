package ray

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

// NewHeadService creates a ClusterIP service that points to the head node.
// Dashboard port is exposed when enabled.
func NewHeadService(rc *dcv1alpha1.RayCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name: "client",
			Port: rc.Spec.ClientServerPort,
		},
		{
			Name: "redis-primary",
			Port: rc.Spec.Port,
		},
	}

	if util.BoolPtrIsTrue(rc.Spec.EnableDashboard) {
		ports = append(ports, corev1.ServicePort{
			Name: "dashboard",
			Port: rc.Spec.DashboardPort,
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HeadServiceName(rc.Name),
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentHead),
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: SelectorLabelsWithComponent(rc, ComponentHead),
		},
	}
}

package spark

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

// NewMasterService creates a ClusterIP service that points to the head node.
// Dashboard port is exposed when enabled.
func NewMasterService(sc *dcv1alpha1.SparkCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name: "cluster",
			Port: sc.Spec.ClusterPort,
			TargetPort: intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "cluster",
			},
		},
	}
	if util.BoolPtrIsTrue(sc.Spec.EnableDashboard) {
		ports = append(ports, corev1.ServicePort{
			Name:     "tcp", // named tcp to prevent istio from sniffing for Host
			Port:     sc.Spec.DashboardPort,
			Protocol: corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "http",
			},
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HeadServiceName(sc.Name),
			Namespace: sc.Namespace,
			Labels:    MetadataLabelsWithComponent(sc, ComponentMaster),
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Ports:    ports,
			Selector: SelectorLabelsWithComponent(sc, ComponentMaster),
		},
	}
}

// NewHeadlessService creates a headless service that points to worker nodes
func NewHeadlessService(sc *dcv1alpha1.SparkCluster) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HeadlessServiceName(sc.Name),
			Namespace: sc.Namespace,
			Labels:    MetadataLabelsWithComponent(sc, ComponentMaster),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector:  SelectorLabels(sc),
			Ports: []corev1.ServicePort{
				{
					Name: "cluster",
					Port: sc.Spec.ClusterPort,
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "cluster",
					},
				},
			},
		},
	}
}

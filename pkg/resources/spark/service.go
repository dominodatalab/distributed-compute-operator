package spark

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// NewHeadService creates a ClusterIP service that points to the head node.
// Dashboard port is exposed when enabled.
func NewHeadService(rc *dcv1alpha1.SparkCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name: "cluster",
			Port: rc.Spec.ClusterPort,
			TargetPort: intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "cluster",
			},
		},
		{
			Name:     "tcp", // named tcp to prevent istio from sniffing for Host
			Port:     rc.Spec.DashboardPort,
			Protocol: corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "http",
			},
		},
	}

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      HeadServiceName(rc.Name),
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentHead),
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeClusterIP,
			Ports: ports,
			Selector: SelectorLabelsWithComponent(rc, ComponentHead),
		},
	}
}

func NewHeadlessService(rc *dcv1alpha1.SparkCluster) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			//Name:      HeadlessServiceName(rc.Name),
			Name:      "example-spark-worker",
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentHead),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector:  SelectorLabels(rc),
		},
	}
}

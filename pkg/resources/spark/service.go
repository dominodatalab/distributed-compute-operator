package spark

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// NewMasterService creates a ClusterIP service that points to the head node.
// Dashboard port is exposed when enabled.
func NewMasterService(sc *dcv1alpha1.SparkCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name: "tcp-cluster",
			Port: sc.Spec.ClusterPort,
			TargetPort: intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "cluster",
			},
		},
		{
			Name:     "tcp",
			Port:     sc.Spec.MasterWebPort,
			Protocol: corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "http",
			},
		},
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MasterServiceName(sc.Name),
			Namespace: sc.Namespace,
			Labels:    AddGlobalLabels(MetadataLabelsWithComponent(sc, ComponentMaster), sc.Spec.GlobalLabels),
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
			Labels:    AddGlobalLabels(MetadataLabelsWithComponent(sc, ComponentWorker), sc.Spec.GlobalLabels),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector:  SelectorLabels(sc),
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp-cluster",
					Port:       sc.Spec.ClusterPort,
					TargetPort: intstr.FromString("cluster"),
				},
				{
					Name:       "tcp-master-webport",
					Port:       sc.Spec.MasterWebPort,
					TargetPort: intstr.FromString("http"),
					Protocol:   corev1.ProtocolTCP,
				}, {
					Name:       "tcp-worker-webport",
					Port:       sc.Spec.WorkerWebPort,
					TargetPort: intstr.FromString("http"),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-driver-block-manager",
					Port:     sc.Spec.Driver.BlockManagerPort,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

// NewSparkDriverService creates a ClusterIP service that exposes the driver UI port.
func NewSparkDriverService(sc *dcv1alpha1.SparkCluster) *corev1.Service {
	targetPort := intstr.FromInt(int(sc.Spec.Driver.UIPort))

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DriverServiceName(sc.Name),
			Namespace: sc.Namespace,
			Labels:    AddGlobalLabels(MetadataLabels(sc), sc.Spec.GlobalLabels),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Type:      "ClusterIP",
			Selector:  sc.Spec.Driver.Selector,
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp-ui",
					Port:       sc.Spec.Driver.UIPort,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: targetPort,
				},
				{
					Name:     "tcp-driver",
					Port:     sc.Spec.Driver.Port,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-block-manager",
					Port:     sc.Spec.Driver.BlockManagerPort,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

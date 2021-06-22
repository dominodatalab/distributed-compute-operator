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
			// deliberately named tcp to prevent istio from sniffing for Host
			Name:     "tcp",
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
			Labels:    MetadataLabelsWithComponent(sc, ComponentWorker),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector:  SelectorLabelsWithComponent(sc, ComponentWorker),
			Ports:     []corev1.ServicePort{},
			// {
			//	Name:       "cluster",
			//	Port:       sc.Spec.ClusterPort,
			//	TargetPort: intstr.FromString("cluster"),
			// },
			// {
			//	Name:       "tcp-master-webport",
			//	Port:       sc.Spec.TCPMasterWebPort,
			//	TargetPort: intstr.FromString("http"),
			//	Protocol:   corev1.ProtocolTCP,
			// }, {
			//	Name:       "tcp-worker-webport",
			//	Port:       sc.Spec.TCPWorkerWebPort,
			//	TargetPort: intstr.FromString("http"),
			//	Protocol:   corev1.ProtocolTCP,
			// },
			// },
		},
	}
}

// NewSparkDriverService creates a ClusterIP service that exposes the driver UI port.
func NewSparkDriverService(sc *dcv1alpha1.SparkCluster) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DriverServiceName(sc.Spec.Driver.SparkClusterName),
			Namespace: sc.Namespace,
			Labels:    MetadataLabelsWithComponent(sc, ComponentNone),
			// not sure if we need annotations on this service
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Type:      "ClusterIP",
			Selector:  map[string]string{"app.kubernetes.io/instance": sc.Spec.Driver.ExecutionName},
			Ports: []corev1.ServicePort{
				{
					Name:     sc.Spec.Driver.DriverUIPortName,
					Port:     sc.Spec.Driver.DriverUIPort,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     sc.Spec.Driver.DriverBlockManagerPortName,
					Port:     sc.Spec.Driver.DriverBlockManagerPort,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     sc.Spec.Driver.DriverPortName,
					Port:     sc.Spec.Driver.DriverPort,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

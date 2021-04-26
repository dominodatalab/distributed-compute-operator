package ray

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

// NewClientService creates a ClusterIP service that points to the head
// node that exposes the client server port, and dashboard port when enabled.
func NewClientService(rc *dcv1alpha1.RayCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name:       "tcp-client",
			Port:       rc.Spec.ClientServerPort,
			TargetPort: intstr.FromInt(int(rc.Spec.ClientServerPort)),
		},
	}

	if util.BoolPtrIsTrue(rc.Spec.EnableDashboard) {
		ports = append(ports, corev1.ServicePort{
			Name:       "http-dashboard",
			Port:       rc.Spec.DashboardPort,
			TargetPort: intstr.FromInt(int(rc.Spec.DashboardPort)),
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, "client"),
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentHead),
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: SelectorLabelsWithComponent(rc, ComponentHead),
		},
	}
}

// NewHeadlessHeadService creates a headless service that points to the head
// node and exposes cluster communication ports.
func NewHeadlessHeadService(rc *dcv1alpha1.RayCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name:       "tcp-gcs-server",
			Port:       rc.Spec.GCSServerPort,
			TargetPort: intstr.FromInt(int(rc.Spec.GCSServerPort)),
		},
		{
			Name:       "tcp-redis-primary",
			Port:       rc.Spec.Port,
			TargetPort: intstr.FromInt(int(rc.Spec.Port)),
		},
	}
	for idx, port := range rc.Spec.RedisShardPorts {
		ports = append(ports, corev1.ServicePort{
			Name:       fmt.Sprintf("tcp-redis-shard-%d", idx),
			Port:       port,
			TargetPort: intstr.FromInt(int(port)),
		})
	}
	ports = append(ports, workerPorts(rc)...)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HeadlessHeadServiceName(rc.Name),
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentHead),
		},
		Spec: corev1.ServiceSpec{
			Ports:     ports,
			Selector:  SelectorLabelsWithComponent(rc, ComponentHead),
			ClusterIP: corev1.ClusterIPNone,
		},
	}
}

// NewHeadlessWorkerService creates a headless service that points to the
// worker nodes and exposes cluster communication ports.
func NewHeadlessWorkerService(rc *dcv1alpha1.RayCluster) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HeadlessWorkerServiceName(rc.Name),
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentWorker),
		},
		Spec: corev1.ServiceSpec{
			Ports:     workerPorts(rc),
			Selector:  SelectorLabelsWithComponent(rc, ComponentWorker),
			ClusterIP: corev1.ClusterIPNone,
		},
	}
}

func workerPorts(rc *dcv1alpha1.RayCluster) []corev1.ServicePort {
	ports := []corev1.ServicePort{
		{
			Name:       "tcp-object-manager",
			Port:       rc.Spec.ObjectManagerPort,
			TargetPort: intstr.FromInt(int(rc.Spec.ObjectManagerPort)),
		},
		{
			Name:       "tcp-node-manager",
			Port:       rc.Spec.NodeManagerPort,
			TargetPort: intstr.FromInt(int(rc.Spec.NodeManagerPort)),
		},
	}
	for idx, port := range rc.Spec.WorkerPorts {
		ports = append(ports, corev1.ServicePort{
			Name:       fmt.Sprintf("tcp-worker-port-%d", idx),
			Port:       port,
			TargetPort: intstr.FromInt(int(port)),
		})
	}

	return ports
}

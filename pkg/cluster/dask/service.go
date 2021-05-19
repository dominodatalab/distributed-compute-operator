package dask

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
)

type serviceDS struct {
	dc   *dcv1alpha1.DaskCluster
	comp metadata.Component
}

func SchedulerService(obj client.Object) components.ServiceDataSource {
	return &serviceDS{
		dc:   obj.(*dcv1alpha1.DaskCluster),
		comp: ComponentScheduler,
	}
}

func WorkerService(obj client.Object) components.ServiceDataSource {
	return &serviceDS{
		dc:   obj.(*dcv1alpha1.DaskCluster),
		comp: ComponentWorker,
	}
}

func (s *serviceDS) GetService() *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(s.dc, s.comp),
			Namespace: s.dc.Namespace,
			Labels:    meta.StandardLabels(s.dc),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector:  meta.MatchLabelsWithComponent(s.dc, s.comp),
			Ports:     s.ports(),
		},
	}

	return svc
}

func (s *serviceDS) ports() []corev1.ServicePort {
	var ports []corev1.ServicePort

	if s.comp == ComponentScheduler {
		ports = []corev1.ServicePort{
			{
				Name:       "serve",
				Port:       s.dc.Spec.SchedulerPort,
				TargetPort: intstr.FromString("serve"),
			},
			{
				Name:       "dashboard",
				Port:       s.dc.Spec.DashboardPort,
				TargetPort: intstr.FromString("dashboard"),
			},
		}
	} else {
		ports = []corev1.ServicePort{
			{
				Name:       "worker",
				Port:       s.dc.Spec.WorkerPort,
				TargetPort: intstr.FromString("worker"),
			},
			{
				Name:       "nanny",
				Port:       s.dc.Spec.NannyPort,
				TargetPort: intstr.FromString("nanny"),
			},
			{
				Name:       "dashboard",
				Port:       s.dc.Spec.DashboardPort,
				TargetPort: intstr.FromString("dashboard"),
			},
		}
	}

	return ports
}

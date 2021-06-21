package dask

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func ServiceScheduler() core.OwnedComponent {
	return components.Service(func(obj client.Object) components.ServiceDataSource {
		return &serviceDS{dc: daskCluster(obj), comp: ComponentScheduler}
	})
}

func ServiceWorker() core.OwnedComponent {
	return components.Service(func(obj client.Object) components.ServiceDataSource {
		return &serviceDS{dc: daskCluster(obj), comp: ComponentWorker}
	})
}

type serviceDS struct {
	dc   *dcv1alpha1.DaskCluster
	comp metadata.Component
}

func (s *serviceDS) Service() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(s.dc, s.comp),
			Namespace: s.dc.Namespace,
			Labels:    meta.StandardLabelsWithComponent(s.dc, s.comp),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector:  meta.MatchLabelsWithComponent(s.dc, s.comp),
			Ports:     s.ports(),
		},
	}
}

func (s *serviceDS) ports() []corev1.ServicePort {
	if s.comp == ComponentScheduler {
		return []corev1.ServicePort{
			{
				Name:       "tcp-serve",
				Port:       s.dc.Spec.SchedulerPort,
				TargetPort: intstr.FromString("serve"),
			},
			{
				Name:       "tcp-dashboard",
				Port:       s.dc.Spec.DashboardPort,
				TargetPort: intstr.FromString("dashboard"),
			},
		}
	}

	return []corev1.ServicePort{
		{
			Name:       "tcp-worker",
			Port:       s.dc.Spec.WorkerPort,
			TargetPort: intstr.FromString("worker"),
		},
		{
			Name:       "tcp-nanny",
			Port:       s.dc.Spec.NannyPort,
			TargetPort: intstr.FromString("nanny"),
		},
		{
			Name:       "tcp-dashboard",
			Port:       s.dc.Spec.DashboardPort,
			TargetPort: intstr.FromString("dashboard"),
		},
	}
}

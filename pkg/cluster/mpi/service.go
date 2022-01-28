package mpi

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func ServiceWorker() core.OwnedComponent {
	return &serviceComponent{}
}

type serviceComponent struct{}

func (c serviceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(cr),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabelsWithComponent(cr, ComponentWorker, cr.Spec.Worker.Labels),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector:  meta.MatchLabelsWithComponent(cr, ComponentWorker),
			Ports: []corev1.ServicePort{
				{
					Name:       sshdPortName,
					Port:       sshdPort,
					TargetPort: intstr.FromString(sshdPortName),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       rsyncPortName,
					Port:       rsyncPort,
					TargetPort: intstr.FromString(rsyncPortName),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, cr, svc)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

func (c serviceComponent) Kind() client.Object {
	return &corev1.Service{}
}

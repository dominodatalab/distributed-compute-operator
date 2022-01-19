package mpi

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func ServiceWorker() core.OwnedComponent {
	return &serviceComponent{
		comp: ComponentWorker,
	}
}

func ServiceClient() core.OwnedComponent {
	return &serviceComponent{
		comp: ComponentClient,
	}
}

type serviceComponent struct {
	comp metadata.Component
}

func (c serviceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	ports := []corev1.ServicePort{}
	var selector map[string]string
	var clusterIP string
	var extraLabels map[string]string
	switch c.comp {
	case ComponentClient:
		selector = cr.Spec.NetworkPolicy.ClientLabels
		clusterIP = ""
		extraLabels = map[string]string{}
	case ComponentWorker:
		ports = append(ports, corev1.ServicePort{
			Name:       sshdPortName,
			Port:       sshdPort,
			TargetPort: intstr.FromString(sshdPortName),
			Protocol:   corev1.ProtocolTCP,
		})

		selector = meta.MatchLabelsWithComponent(cr, c.comp)
		clusterIP = corev1.ClusterIPNone
		extraLabels = cr.Spec.Worker.Labels
	case metadata.ComponentNone:
		err := errors.New("unknown component for NetworkPolicy")
		return ctrl.Result{}, err
	}

	ports = append(ports, mpiPorts(cr)...)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(cr, c.comp),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabelsWithComponent(cr, c.comp, extraLabels),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: clusterIP,
			Selector:  selector,
			Ports:     ports,
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

func mpiPorts(cr *v1alpha1.MPICluster) []corev1.ServicePort {
	ports := []corev1.ServicePort{}
	for idx, port := range cr.Spec.WorkerPorts {
		ports = append(ports, corev1.ServicePort{
			Name:       fmt.Sprintf("tcp-mpi-%d", idx),
			Port:       port,
			TargetPort: intstr.FromInt(int(port)),
		})
	}
	return ports
}

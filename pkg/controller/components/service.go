package components

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type ServiceDataSource interface {
	GetService() *corev1.Service
}

type serviceDataSourceFactory func(object client.Object) ServiceDataSource

type ServiceComponent struct {
	factory serviceDataSourceFactory
}

func Service(f serviceDataSourceFactory) *ServiceComponent {
	return &ServiceComponent{factory: f}
}

func (comp *ServiceComponent) Kind() client.Object {
	return &corev1.Service{}
}

func (comp *ServiceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := comp.factory(ctx.Object)

	svc := ds.GetService()
	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, svc)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}
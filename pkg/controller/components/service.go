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
	Service() *corev1.Service
}

type ServiceDataSourceFactory func(client.Object) ServiceDataSource

func Service(f ServiceDataSourceFactory) core.OwnedComponent {
	return &serviceComponent{factory: f}
}

type serviceComponent struct {
	factory ServiceDataSourceFactory
}

func (c *serviceComponent) Kind() client.Object {
	return &corev1.Service{}
}

func (c *serviceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)
	svc := ds.Service()

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, svc)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

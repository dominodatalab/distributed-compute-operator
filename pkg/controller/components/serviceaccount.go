//nolint:dupl
package components

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type ServiceAccountDataSource interface {
	ServiceAccount() *corev1.ServiceAccount
	Delete() bool
}

type ServiceAccountDataSourceFactory func(client.Object) ServiceAccountDataSource

func ServiceAccount(f ServiceAccountDataSourceFactory) core.OwnedComponent {
	return &serviceAccountComponent{factory: f}
}

type serviceAccountComponent struct {
	factory ServiceAccountDataSourceFactory
}

func (c *serviceAccountComponent) Kind() client.Object {
	return &corev1.ServiceAccount{}
}

func (c *serviceAccountComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)
	sa := ds.ServiceAccount()

	if ds.Delete() {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, sa)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, sa)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service account: %w", err)
	}

	return ctrl.Result{}, err
}

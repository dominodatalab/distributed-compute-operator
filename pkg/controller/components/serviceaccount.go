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
	GetServiceAccount() *corev1.ServiceAccount
	IsDelete() bool
}

type serviceAccountDataSourceFactory func(client.Object) ServiceAccountDataSource

type serviceAccountComponent struct {
	factory serviceAccountDataSourceFactory
}

func ServiceAccount(f serviceAccountDataSourceFactory) *serviceAccountComponent {
	return &serviceAccountComponent{factory: f}
}

func (comp *serviceAccountComponent) Kind() client.Object {
	return &corev1.ServiceAccount{}
}

func (comp *serviceAccountComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := comp.factory(ctx.Object)
	sa := ds.GetServiceAccount()

	if ds.IsDelete() {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, sa)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, sa)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service account: %w", err)
	}

	return ctrl.Result{}, err
}

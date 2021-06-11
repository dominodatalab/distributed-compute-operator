//nolint:dupl
package components

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type NetworkPolicyDataSource interface {
	NetworkPolicy() *networkingv1.NetworkPolicy
	Delete() bool
}

type NetworkPolicyDataSourceFactory func(client.Object) NetworkPolicyDataSource

func NetworkPolicy(f NetworkPolicyDataSourceFactory) core.OwnedComponent {
	return &networkPolicyComponent{factory: f}
}

type networkPolicyComponent struct {
	factory NetworkPolicyDataSourceFactory
}

func (c *networkPolicyComponent) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

func (c *networkPolicyComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)
	netpol := ds.NetworkPolicy()

	if ds.Delete() {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, netpol)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, netpol)
	if err != nil {
		err = fmt.Errorf("cannot reconcile network policy: %w", err)
	}

	return ctrl.Result{}, err
}

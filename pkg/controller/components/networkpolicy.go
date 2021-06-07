package components

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type DeletePredicateFn func(client.Object) bool
type NetworkPolicyFactory func(client.Object) *networkingv1.NetworkPolicy

type networkPolicy struct {
	factory NetworkPolicyFactory
	delete  DeletePredicateFn
}

func NetworkPolicy(f NetworkPolicyFactory, fn DeletePredicateFn) core.Component {
	return &networkPolicy{
		factory: f,
		delete:  fn,
	}
}

func (c *networkPolicy) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

func (c *networkPolicy) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	netpol := c.factory(ctx.Object)

	if c.delete(ctx.Object) {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, netpol)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, netpol)
	if err != nil {
		err = fmt.Errorf("cannot reconcile network policy: %w", err)
	}

	return ctrl.Result{}, err
}

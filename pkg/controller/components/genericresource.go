package components

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type resourceFactory func(ctx *core.Context) (client.Object, bool, error)

// NOTE: this is just a POC to illustrate we can leverage a single component
type GenericResourceComponent struct {
	kind    client.Object
	factory resourceFactory
}

func Resource(k client.Object, f resourceFactory) *GenericResourceComponent {
	// NOTE: we require an additional param for ownership registration. the factory method is designed to ONLY return an
	// 	object during runtime execution.
	return &GenericResourceComponent{
		kind:    k,
		factory: f,
	}
}

func (comp *GenericResourceComponent) Kind() client.Object {
	return comp.kind
}

func (comp *GenericResourceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj, destroy, err := comp.factory(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}
	if destroy {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, obj)
	}
	if obj == nil {
		return ctrl.Result{}, nil
	}

	err = actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, obj)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service account: %w", err)
	}

	return ctrl.Result{}, err
}

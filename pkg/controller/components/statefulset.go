package components

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type StatefulSetDataSource interface {
	GetStatefulSet() (*appsv1.StatefulSet, error)
}

type statefulSetDataSourceFactory func(object client.Object) StatefulSetDataSource

type statefulSetComponent struct {
	factory statefulSetDataSourceFactory
}

func StatefulSet(f statefulSetDataSourceFactory) *statefulSetComponent {
	return &statefulSetComponent{factory: f}
}

func (comp *statefulSetComponent) Kind() client.Object {
	return &appsv1.StatefulSet{}
}

func (comp *statefulSetComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := comp.factory(ctx.Object)

	sts, err := ds.GetStatefulSet()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build statefulset: %w", err)
	}

	err = actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, sts)
	if err != nil {
		err = fmt.Errorf("cannot reconcile stateful set: %w", err)
	}

	return ctrl.Result{}, err
}

func (comp *statefulSetComponent) Finalize(ctx *core.Context) (ctrl.Result, bool, error) {
	return ctrl.Result{}, true, nil
}

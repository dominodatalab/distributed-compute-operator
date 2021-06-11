package components

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type StatefulSetDataSource interface {
	StatefulSet() (*appsv1.StatefulSet, error)
	PVCListOpts() []client.ListOption
}

type StatefulSetDataSourceFactory func(client.Object) StatefulSetDataSource

func StatefulSet(f StatefulSetDataSourceFactory) core.OwnedComponent {
	return &statefulSetComponent{factory: f}
}

type statefulSetComponent struct {
	factory StatefulSetDataSourceFactory
}

func (c *statefulSetComponent) Kind() client.Object {
	return &appsv1.StatefulSet{}
}

func (c *statefulSetComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)

	sts, err := ds.StatefulSet()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build statefulset: %w", err)
	}

	err = actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, sts)
	if err != nil {
		err = fmt.Errorf("cannot reconcile stateful set: %w", err)
	}

	return ctrl.Result{}, err
}

func (c *statefulSetComponent) Finalize(ctx *core.Context) (ctrl.Result, bool, error) {
	ds := c.factory(ctx.Object)

	pvcList := &corev1.PersistentVolumeClaimList{}
	listOpts := ds.PVCListOpts()

	ctx.Log.Info("Querying for persistent volume claims", "opts", listOpts)
	if err := ctx.Client.List(ctx, pvcList, listOpts...); err != nil {
		ctx.Log.Error(err, "cannot list persistent volume claims")
		return ctrl.Result{}, false, err
	}

	for idx := range pvcList.Items {
		pvc := &pvcList.Items[idx]
		key := client.ObjectKeyFromObject(pvc)

		ctx.Log.Info("Deleting persistent volume claim", "claim", key)
		if err := ctx.Client.Delete(ctx, pvc); err != nil {
			ctx.Log.Error(err, "cannot delete persistent volume claim", "claim", key)
			return ctrl.Result{}, false, err
		}
	}

	return ctrl.Result{}, true, nil
}

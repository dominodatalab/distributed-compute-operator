package components

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type ConfigMapDataSource interface {
	ConfigMap() *corev1.ConfigMap
	Delete() bool
}

type ConfigMapDataSourceFactory func(object client.Object) ConfigMapDataSource

func ConfigMap(f ConfigMapDataSourceFactory) core.OwnedComponent {
	return &configMapComponent{factory: f}
}

type configMapComponent struct {
	factory ConfigMapDataSourceFactory
}

func (c *configMapComponent) Kind() client.Object {
	return &corev1.ConfigMap{}
}

func (c *configMapComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)
	cm := ds.ConfigMap()

	if ds.Delete() {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, cm)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, cm)
	if err != nil {
		err = fmt.Errorf("cannot reconcile config map: %w", err)
	}

	return ctrl.Result{}, err
}

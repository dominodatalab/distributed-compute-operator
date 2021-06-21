//nolint:dupl
package components

import (
	"fmt"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type HorizontalPodAutoscalerDataSource interface {
	HorizontalPodAutoscaler() *autoscalingv2beta2.HorizontalPodAutoscaler
	Delete() bool
}

type HorizontalPodAutoscalerDataSourceFactory func(client.Object) HorizontalPodAutoscalerDataSource

func HorizontalPodAutoscaler(f HorizontalPodAutoscalerDataSourceFactory) core.OwnedComponent {
	return &horizontalPodAutoscaler{factory: f}
}

type horizontalPodAutoscaler struct {
	factory HorizontalPodAutoscalerDataSourceFactory
}

func (c *horizontalPodAutoscaler) Kind() client.Object {
	return &autoscalingv2beta2.HorizontalPodAutoscaler{}
}

func (c *horizontalPodAutoscaler) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)
	hpa := ds.HorizontalPodAutoscaler()

	if ds.Delete() {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, hpa)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, hpa)
	if err != nil {
		err = fmt.Errorf("cannot reconcile horizontal pod autoscaler: %w", err)
	}

	return ctrl.Result{}, err
}

package components

import (
	"fmt"
	"reflect"
	"sort"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type ClusterStatusUpdateDataSource interface {
	ListOpts() []client.ListOption
	StatefulSet() *appsv1.StatefulSet
	ClusterStatusConfig() *dcv1alpha1.ClusterStatusConfig
}

type ClusterStatusUpdateDataSourceFactory func(client.Object) ClusterStatusUpdateDataSource

func ClusterStatusUpdate(f ClusterStatusUpdateDataSourceFactory) core.Component {
	return &clusterStatusUpdateComponent{factory: f}
}

type clusterStatusUpdateComponent struct {
	factory func(client.Object) ClusterStatusUpdateDataSource
}

func (c *clusterStatusUpdateComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	var modified bool

	ds := c.factory(ctx.Object)
	csc := ds.ClusterStatusConfig()

	// modify node list field
	podList := &corev1.PodList{}
	listOpts := ds.ListOpts()
	if err := ctx.Client.List(ctx, podList, listOpts...); err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot list cluster pods: %w", err)
	}

	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	sort.Strings(podNames)

	if !reflect.DeepEqual(podNames, csc.Nodes) {
		csc.Nodes = podNames
		modified = true
	}

	// modify scale subresource fields
	sts := ds.StatefulSet()
	err := ctx.Client.Get(ctx, client.ObjectKeyFromObject(sts), sts)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	selector, err := metav1.LabelSelectorAsSelector(sts.Spec.Selector)
	if err != nil {
		return ctrl.Result{}, err
	}

	if csc.WorkerSelector != selector.String() {
		csc.WorkerSelector = selector.String()
		modified = true
	}
	if csc.WorkerReplicas != *sts.Spec.Replicas { // NOTE: panic: runtime error: invalid memory address or nil pointer dereference
		csc.WorkerReplicas = *sts.Spec.Replicas
		modified = true
	}

	if modified {
		err = ctx.Client.Status().Update(ctx, ctx.Object)
	}

	return ctrl.Result{}, err
}

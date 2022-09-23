package components

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	v1 "k8s.io/api/batch/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

type ClusterStatusUpdateDataSource interface {
	ListOpts() []client.ListOption
	StatefulSet() *appsv1.StatefulSet
	ClusterStatusConfig() *dcv1alpha1.ClusterStatusConfig
	Image() *dcv1alpha1.OCIImageDefinition
}

const finalizerRetryPeriod = 1 * time.Second

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
	var masterPod corev1.Pod
	masterPodFound := false
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
		if !masterPodFound && strings.Contains(pod.Name, "-scheduler") {
			masterPodFound = true
			masterPod = pod
		}
	}
	sort.Strings(podNames)

	if !reflect.DeepEqual(podNames, csc.Nodes) {
		csc.Nodes = podNames
		modified = true
	}

	// modify scale subresource fields
	sts := ds.StatefulSet()
	if err := ctx.Client.Get(ctx, client.ObjectKeyFromObject(sts), sts); client.IgnoreNotFound(err) != nil {
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

	// store canonical image reference
	image, err := util.ParseImageDefinition(ds.Image())
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot build cluster image: %w", err)
	}
	if csc.Image != image {
		csc.Image = image
		modified = true
	}

	var status v1.JobConditionType
	if masterPodFound {
		if dcv1alpha1.IsPodReady(masterPod) {
			status = dcv1alpha1.RunningStatus
		} else {
			status = dcv1alpha1.StartingStatus
		}
	} else {
		status = dcv1alpha1.PendingStatus
	}
	if csc.ClusterStatus != status && ctx.Object.GetDeletionTimestamp() == nil {
		modified = true
		csc.ClusterStatus = status
		if status == dcv1alpha1.RunningStatus {
			tt := metav1.Now()
			csc.StartTime = &tt
		} else {
			csc.StartTime = nil
		}
	}

	// only update when fields have changed
	if modified {
		err = ctx.Client.Status().Update(ctx, ctx.Object)
	}

	return ctrl.Result{}, err
}

func (c clusterStatusUpdateComponent) Finalize(ctx *core.Context) (ctrl.Result, bool, error) {
	ds := c.factory(ctx.Object)
	csc := ds.ClusterStatusConfig()

	if csc.ClusterStatus != dcv1alpha1.StoppingStatus {
		csc.ClusterStatus = dcv1alpha1.StoppingStatus
		csc.StartTime = nil
		err := ctx.Client.Status().Update(ctx, ctx.Object)
		if err != nil {
			return ctrl.Result{RequeueAfter: finalizerRetryPeriod}, false,
				fmt.Errorf("cannot update cluster status: %w", err)
		}
	}

	return ctrl.Result{}, true, nil
}

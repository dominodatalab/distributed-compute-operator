package mpi

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	v1 "k8s.io/api/batch/v1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

const (
	PendingStatus  v1.JobConditionType = "Pending"
	StartingStatus v1.JobConditionType = "Starting"
	RunningStatus  v1.JobConditionType = "Running"
	StoppingStatus v1.JobConditionType = "Stopping"
)

func StatusUpdate() core.Component {
	return &statusUpdateComponent{}
}

type statusUpdateComponent struct{}

func (c statusUpdateComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	var modified bool

	image, err := util.ParseImageDefinition(cr.Spec.Image)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot build cluster image: %w", err)
	}
	if cr.Status.Image != image {
		cr.Status.Image = image
		modified = true
	}

	pods, err := getPods(ctx, cr)
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("cannot list cluster pods: %w", err)
	}
	var podNames []string
	var runningPodCnt = 0
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
		if pod.Status.Phase == corev1.PodRunning {
			runningPodCnt++
		}
	}
	sort.Strings(podNames)
	if !reflect.DeepEqual(podNames, cr.Status.Nodes) {
		cr.Status.Nodes = podNames
		modified = true
	}

	expectedPodCnt := int(*cr.Spec.Worker.Replicas)

	var status v1.JobConditionType
	if runningPodCnt >= expectedPodCnt {
		status = RunningStatus
	} else if runningPodCnt == 0 {
		status = PendingStatus
	} else {
		status = StartingStatus
	}

	if cr.Status.ClusterStatus != status {
		modified = true
		cr.Status.ClusterStatus = status
		if status == RunningStatus {
			tt := metav1.Now()
			cr.Status.StartTime = &tt
		} else {
			cr.Status.StartTime = nil
		}
	}

	if modified {
		err = ctx.Client.Status().Update(ctx, cr)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("cannot update cluster status: %w", err)
		}
	}

	return ctrl.Result{}, nil
}

func (c statusUpdateComponent) Finalize(ctx *core.Context) (ctrl.Result, bool, error) {
	cr := objToMPICluster(ctx.Object)

	if cr.Status.ClusterStatus != StoppingStatus {
		cr.Status.ClusterStatus = StoppingStatus
		err := ctx.Client.Status().Update(ctx, cr)
		if err != nil {
			return ctrl.Result{}, false, fmt.Errorf("cannot update cluster status: %w", err)
		}
	}

	pods, err := getPods(ctx, cr)
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{}, false, fmt.Errorf("cannot list cluster pods: %w", err)
	}
	podCnt := len(pods)
	if podCnt == 0 {
		return ctrl.Result{}, true, nil
	} else {
		return ctrl.Result{RequeueAfter: time.Second}, false, nil
	}
}

func getPods(ctx *core.Context, cr *dcv1alpha1.MPICluster) ([]corev1.Pod, error) {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(cr.Namespace),
		client.MatchingLabels(meta.StandardLabels(cr)),
	}
	err := ctx.Client.List(ctx, podList, listOpts...)
	return podList.Items, err // first item is an empty array when err == nil
}

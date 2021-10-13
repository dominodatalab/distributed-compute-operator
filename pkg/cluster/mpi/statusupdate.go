package mpi

import (
	"fmt"
	"reflect"
	"sort"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

func StatusUpdate() core.Component {
	return &statusUpdateComponent{}
}

type statusUpdateComponent struct{}

func (c statusUpdateComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPIJob(ctx.Object)

	var modified bool

	// update image reference
	image, err := util.ParseImageDefinition(cr.Spec.Image)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot build cluster image: %w", err)
	}
	if cr.Status.Image != image {
		cr.Status.Image = image
		modified = true
	}

	// update node list
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(cr.Namespace),
		client.MatchingLabels(meta.StandardLabels(cr)),
	}
	if lErr := ctx.Client.List(ctx, podList, listOpts...); lErr != nil {
		return ctrl.Result{}, fmt.Errorf("cannot list cluster pods: %w", lErr)
	}

	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	sort.Strings(podNames)

	if !reflect.DeepEqual(podNames, cr.Status.Nodes) {
		cr.Status.Nodes = podNames
		modified = true
	}

	// update job details
	var job batchv1.Job
	objKey := client.ObjectKey{
		Name:      jobName(cr),
		Namespace: cr.Namespace,
	}

	if gErr := ctx.Client.Get(ctx, objKey, &job); gErr != nil {
		if apierrors.IsNotFound(gErr) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, gErr
	}

	if cr.Status.StartTime != job.Status.StartTime {
		cr.Status.StartTime = job.Status.StartTime
		modified = true
	}
	if cr.Status.CompletionTime != job.Status.CompletionTime {
		cr.Status.CompletionTime = job.Status.CompletionTime
		modified = true
	}

	oldStatus := cr.Status.LauncherStatus
	switch {
	case job.Status.Conditions != nil:
		cr.Status.LauncherStatus = job.Status.Conditions[0].Type
	case job.Status.Active == 1:
		cr.Status.LauncherStatus = "Active"
	default:
		cr.Status.LauncherStatus = "Pending"
	}

	if cr.Status.LauncherStatus != oldStatus {
		modified = true
	}

	if modified {
		err = ctx.Client.Status().Update(ctx, cr)
	}

	// TODO: scale down workers after job is "Complete" or "Failed"
	// 	this component focuses on updating the "status" so the downscaling
	// 	logic probably belongs elsewhere

	return ctrl.Result{}, err
}

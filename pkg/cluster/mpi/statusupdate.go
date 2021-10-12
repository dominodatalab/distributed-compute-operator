package mpi

import (
	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func StatusUpdate() core.Component {
	return &statusUpdateComponent{}
}

type statusUpdateComponent struct{}

func (c statusUpdateComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPIJob(ctx.Object)

	var job batchv1.Job
	objKey := client.ObjectKey{
		Name:      jobName(cr),
		Namespace: cr.Namespace,
	}

	if err := ctx.Client.Get(ctx, objKey, &job); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	cr.Status.StartTime = job.Status.StartTime
	cr.Status.CompletionTime = job.Status.CompletionTime

	switch {
	case job.Status.Conditions != nil:
		cr.Status.LauncherStatus = job.Status.Conditions[0].Type
	case job.Status.Active == 1:
		cr.Status.LauncherStatus = "Active"
	default:
		cr.Status.LauncherStatus = "Pending"
	}

	// add image, scale down sts

	return ctrl.Result{}, ctx.Client.Status().Update(ctx, cr)
}

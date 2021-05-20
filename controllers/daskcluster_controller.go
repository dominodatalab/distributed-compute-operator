package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/dask"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=daskclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=daskclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=daskclusters/finalizers,verbs=update

func DaskCluster(mgr ctrl.Manager) error {
	return core.NewReconciler(mgr).
		For(&dcv1alpha1.DaskCluster{}).
		Component("serviceaccount", components.ServiceAccount(dask.ServiceAccount)).
		Component("svc-scheduler", components.Service(dask.SchedulerService)).
		Component("sts-scheduler", components.StatefulSet(dask.SchedulerStatefulSet)).
		Component("svc-workers", components.Service(dask.WorkerService)).
		Component("sts-workers", components.StatefulSet(dask.WorkerStatefulSet)).
		WithWebhooks().
		Complete()
}

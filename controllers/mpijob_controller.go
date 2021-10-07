package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/mpi"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpijobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpijobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpijobs/finalizers,verbs=update

// MPIJob builds a controller that reconciles MPIJob objects and registers it with the manager.
func MPIJob(mgr ctrl.Manager, webhooksEnabled, istioEnabled bool) error {
	reconciler := core.NewReconciler(mgr).
		For(&dcv1alpha1.MPIJob{}).
		Component("serviceaccount", mpi.ServiceAccount()).
		Component("secret", mpi.Secret()).
		Component("configmap", mpi.ConfigMap()).
		Component("service", mpi.ServiceWorker()).
		Component("workers", mpi.StatefulSet()).
		Component("launcher", mpi.Job())

	if webhooksEnabled {
		reconciler.WithWebhooks()
	}
	if istioEnabled {
		panic("implement istio support")
	}

	return reconciler.Complete()
}

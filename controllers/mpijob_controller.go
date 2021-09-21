package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpijobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpijobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpijobs/finalizers,verbs=update

// MPIJob builds a controller that reconciles MPIJob objects and registers it with the manager.
func MPIJob(mgr ctrl.Manager, webhooksEnabled, istioEnabled bool) error {
	return core.NewReconciler(mgr).
		For(&dcv1alpha1.MPIJob{}).
		Complete()
}

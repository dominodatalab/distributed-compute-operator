package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/mpi"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpiclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpiclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=mpiclusters/finalizers,verbs=update

// MPICluster builds a controller that reconciles MPICluster objects and registers it with the manager.
func MPICluster(mgr ctrl.Manager, webhooksEnabled, istioEnabled bool) error {
	reconciler := core.NewReconciler(mgr).
		For(&dcv1alpha1.MPICluster{}).
		Component("istio-peerauthentication", mpi.IstioPeerAuthentication(istioEnabled)).
		Component("serviceaccount", mpi.ServiceAccount()).
		Component("role", mpi.RolePodSecurityPolicy()).
		Component("rolebinding", mpi.RoleBindingPodSecurityPolicy()).
		Component("configmap", mpi.ConfigMap()).
		Component("service-worker", mpi.ServiceWorker()).
		Component("service-client", mpi.ServiceClient()).
		Component("networkpolicy-worker", mpi.NetworkPolicyWorker()).
		Component("networkpolicy-client", mpi.NetworkPolicyClient()).
		Component("workers", mpi.StatefulSet()).
		Component("statusupdate", mpi.StatusUpdate())

	if webhooksEnabled {
		reconciler.WithWebhooks()
	}
	return reconciler.Complete()
}

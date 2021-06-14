package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/dask"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=daskclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=daskclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=daskclusters/finalizers,verbs=update

func DaskCluster(mgr ctrl.Manager, istioEnabled bool) error {
	return core.NewReconciler(mgr).
		For(&dcv1alpha1.DaskCluster{}).
		Component("istio-peerauthentication", dask.IstioPeerAuthentication(istioEnabled)).
		Component("serviceaccount", dask.ServiceAccount()).
		Component("role-podsecuritypolicy", dask.RolePodSecurityPolicy()).
		Component("rolebinding-podsecuritypolicy", dask.RoleBindingPodSecurityPolicy()).
		Component("service-scheduler", dask.ServiceScheduler()).
		Component("service-worker", dask.ServiceWorker()).
		Component("networkpolicy-scheduler", dask.NetworkPolicyScheduler()).
		Component("networkpolicy-worker", dask.NetworkPolicyWorker()).
		Component("statefulset-scheduler", dask.StatefulSetScheduler()).
		Component("statefulset-worker", dask.StatefulSetWorker()).
		Component("horizontalpodautoscaler", dask.HorizontalPodAutoscaler()).
		Component("statusupdate", dask.ClusterStatusUpdate()).
		WithWebhooks().
		Complete()
}

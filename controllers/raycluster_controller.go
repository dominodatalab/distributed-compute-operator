package controllers

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/ray"
)

const LastAppliedConfig = "distributed-compute-operator.dominodatalab.com/last-applied"

var (
	PatchAnnotator = patch.NewAnnotator(LastAppliedConfig)
	PatchMaker     = patch.NewPatchMaker(PatchAnnotator)
)

// RayClusterReconciler reconciles RayCluster objects.
type RayClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager creates and registers this controller with the manager.
func (r *RayClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dcv1alpha1.RayCluster{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.Deployment{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&autoscalingv2beta2.HorizontalPodAutoscaler{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=rayclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=rayclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=list;watch
//+kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=create;update;list;watch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=create;update;list;watch
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=create;update;delete;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=create;update;delete;list;watch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=create;update;delete;list;watch

// Reconcile implements state reconciliation logic for RayCluster objects.
func (r *RayClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("raycluster", req.NamespacedName)

	rc := &dcv1alpha1.RayCluster{}
	if err := r.Get(ctx, req.NamespacedName, rc); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("RayCluster resource not found, assuming object was deleted")
			return ctrl.Result{}, nil
		}

		log.Error(err, "Failed to get RayCluster")
		return ctrl.Result{}, err
	}

	if err := r.reconcileResources(ctx, rc); err != nil {
		log.Error(err, "Failed to reconcile cluster resources")
		return ctrl.Result{}, err
	}

	// NOTE: this func will error out during certain update events because of
	// 	generation version conflicts. the desired state is eventually achieved
	// 	but this produces a large amount of noise in the logs. we should figure
	//  out how to remove these errors.
	// if err := r.updateStatus(ctx, rc); err != nil {
	// 	 log.Error(err, "Failed to update cluster status")
	// 	 return ctrl.Result{}, err
	// }

	return ctrl.Result{}, nil
}

// reconcileResources manages the creation and updates of resources that
// collectively comprise a Ray cluster. Each resource is controlled by a parent
// RayCluster object so that full cleanup occurs during a delete operation.
func (r *RayClusterReconciler) reconcileResources(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	if err := r.reconcileServiceAccount(ctx, rc); err != nil {
		return err
	}
	if err := r.reconcileHeadService(ctx, rc); err != nil {
		return err
	}
	if err := r.reconcileNetworkPolicies(ctx, rc); err != nil {
		return err
	}
	if err := r.reconcilePodSecurityPolicyRBAC(ctx, rc); err != nil {
		return err
	}
	if err := r.reconcileAutoscaler(ctx, rc); err != nil {
		return err
	}

	return r.reconcileDeployments(ctx, rc)
}

// reconcileServiceAccount creates a new dedicated service account for a Ray
// cluster unless a different service account name is provided in the spec.
func (r *RayClusterReconciler) reconcileServiceAccount(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	if rc.Spec.ServiceAccountName != "" {
		return nil
	}

	sa := ray.NewServiceAccount(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, sa); err != nil {
		return fmt.Errorf("failed to reconcile service account: %w", err)
	}

	return nil
}

// reconcileHeadService creates a service that points to the head Ray pod and
// applies updates when the parent CR changes.
func (r *RayClusterReconciler) reconcileHeadService(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	svc := ray.NewHeadService(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, svc); err != nil {
		return fmt.Errorf("failed to reconcile head service: %w", err)
	}

	return nil
}

// reconcileNetworkPolicies optionally creates network policies that control
// traffic flow between cluster nodes and external clients.
func (r RayClusterReconciler) reconcileNetworkPolicies(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	headNetpol := ray.NewHeadNetworkPolicy(rc)
	clusterNetpol := ray.NewClusterNetworkPolicy(rc)

	if !rc.Spec.EnableNetworkPolicy {
		return r.deleteIfExists(ctx, headNetpol, clusterNetpol)
	}

	if err := r.createOrUpdateOwnedResource(ctx, rc, clusterNetpol); err != nil {
		return fmt.Errorf("failed to reconcile cluster network policy: %w", err)
	}

	if err := r.createOrUpdateOwnedResource(ctx, rc, headNetpol); err != nil {
		return fmt.Errorf("failed to reconcile head network policy: %w", err)
	}

	return nil
}

// reconcilePodSecurityPolicyRBAC optionally creates a role and role binding
// that allows the Ray pods to "use" the specified pod security policy.
func (r *RayClusterReconciler) reconcilePodSecurityPolicyRBAC(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	role, binding := ray.NewPodSecurityPolicyRBAC(rc)

	if rc.Spec.PodSecurityPolicy == "" {
		return r.deleteIfExists(ctx, role, binding)
	}

	err := r.Get(ctx, types.NamespacedName{Name: rc.Spec.PodSecurityPolicy}, &policyv1beta1.PodSecurityPolicy{})
	if err != nil {
		return fmt.Errorf("cannot verify pod security policy: %w", err)
	}

	if err := r.createOrUpdateOwnedResource(ctx, rc, role); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	if err := r.createOrUpdateOwnedResource(ctx, rc, binding); err != nil {
		return fmt.Errorf("failed to create role binding: %w", err)
	}

	return nil
}

// reconcileAutoscaler optionally creates a horizontal pod autoscaler that
// targets Ray worker pods.
func (r *RayClusterReconciler) reconcileAutoscaler(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	hpa := ray.NewHorizontalPodAutoscaler(rc)

	if rc.Spec.Autoscaling == nil {
		return r.deleteIfExists(ctx, hpa)
	}

	if err := r.createOrUpdateOwnedResource(ctx, rc, hpa); err != nil {
		return fmt.Errorf("failed to reconcile horizontal pod autoscaler: %w", err)
	}

	return nil
}

// reconcileDeployments creates separate Ray head and worker deployments that
// will collectively comprise the execution agents of the cluster.
func (r *RayClusterReconciler) reconcileDeployments(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	head, err := ray.NewDeployment(rc, ray.ComponentHead)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, rc, head); err != nil {
		return fmt.Errorf("failed to create head deployment: %w", err)
	}

	worker, err := ray.NewDeployment(rc, ray.ComponentWorker)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, rc, worker); err != nil {
		return fmt.Errorf("failed to create worker deployment: %w", err)
	}

	selector, err := metav1.LabelSelectorAsSelector(worker.Spec.Selector)
	if err != nil {
		return err
	}

	// update autoscaling fields
	var updateStatus bool
	if rc.Status.WorkerReplicas != *worker.Spec.Replicas {
		rc.Status.WorkerReplicas = *worker.Spec.Replicas
		updateStatus = true
	}
	if rc.Status.WorkerSelector != selector.String() {
		rc.Status.WorkerSelector = selector.String()
		updateStatus = true
	}

	if updateStatus {
		err = r.Status().Update(ctx, rc)
	}

	return err
}

// createOrUpdateOwnedResource should be used to manage the lifecycle of namespace-scoped objects.
//
// The CR will become the "owner" of the "controlled" object and cleanup will
// occur automatically when the CR is deleted.
//
// The controller resource will be created if it's missing.
// The controller resource will be updated if any changes are applicable.
// Any unexpected api errors will be reported.
func (r *RayClusterReconciler) createOrUpdateOwnedResource(ctx context.Context, owner metav1.Object, controlled client.Object) error {
	if err := ctrl.SetControllerReference(owner, controlled, r.Scheme); err != nil {
		return err
	}

	found := controlled.DeepCopyObject().(client.Object)
	err := r.Get(ctx, client.ObjectKeyFromObject(controlled), found)

	if apierrors.IsNotFound(err) {
		if err = PatchAnnotator.SetLastAppliedAnnotation(controlled); err != nil {
			return err
		}

		return r.Create(ctx, controlled)
	}
	if err != nil {
		return err
	}

	patchResult, err := PatchMaker.Calculate(found, controlled, patch.IgnoreStatusFields())
	if err != nil {
		return err
	}
	if patchResult.IsEmpty() {
		return nil
	}

	if err = PatchAnnotator.SetLastAppliedAnnotation(controlled); err != nil {
		return err
	}
	controlled.SetResourceVersion(found.GetResourceVersion())

	if modified, ok := controlled.(*corev1.Service); ok {
		current := found.(*corev1.Service)
		modified.Spec.ClusterIP = current.Spec.ClusterIP
	}

	return r.Update(ctx, controlled)
}

// deleteIfExists will delete one or more Kubernetes objects if they exist.
func (r *RayClusterReconciler) deleteIfExists(ctx context.Context, objs ...client.Object) error {
	for _, obj := range objs {
		err := r.Get(ctx, client.ObjectKeyFromObject(obj), obj)
		if apierrors.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if err = r.Delete(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

//nolint
// updateStatus with a list of pods from both the head and worker deployments.
func (r *RayClusterReconciler) updateStatus(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(rc.Namespace),
		client.MatchingLabels(ray.MetadataLabels(rc)),
	}
	if err := r.List(ctx, podList, listOpts...); err != nil {
		return fmt.Errorf("cannot list ray pods: %w", err)
	}

	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	sort.Strings(podNames)

	if reflect.DeepEqual(podNames, rc.Status.Nodes) {
		return nil
	}

	rc.Status.Nodes = podNames
	if err := r.Status().Update(ctx, rc); err != nil {
		return fmt.Errorf("cannot update ray status nodes: %w", err)
	}

	return nil
}

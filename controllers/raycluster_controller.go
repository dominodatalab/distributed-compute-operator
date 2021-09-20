package controllers

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/logging"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/istio"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/ray"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

// RayClusterReconciler reconciles RayCluster objects.
type RayClusterReconciler struct {
	client.Client
	Log          logging.ContextLogger
	Scheme       *runtime.Scheme
	IstioEnabled bool
}

// SetupWithManager creates and registers this controller with the manager.
func (r *RayClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dcv1alpha1.RayCluster{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.StatefulSet{}).
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
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=create;update;list;watch
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=create;update;delete;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=create;update;delete;list;watch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=create;update;delete;list;watch

// Reconcile implements state reconciliation logic for RayCluster objects.
func (r *RayClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx, log := r.Log.NewContext(ctx, "raycluster", req.NamespacedName)

	log.V(2).Info("reconciliation loop triggered")

	rc := &dcv1alpha1.RayCluster{}
	if err := r.Get(ctx, req.NamespacedName, rc); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("resource not found, assuming object was deleted")
			return ctrl.Result{}, nil
		}

		log.Error(err, "failed to retrieve resource")
		return ctrl.Result{}, err
	}

	if updated, err := r.manageFinalization(ctx, rc); err != nil {
		return ctrl.Result{}, err
	} else if updated {
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.reconcileResources(ctx, rc); err != nil {
		log.Error(err, "failed to reconcile cluster resources")
		return ctrl.Result{}, err
	}

	if err := r.updateStatus(ctx, rc); err != nil {
		if strings.Contains(err.Error(), genericregistry.OptimisticLockErrorMsg) {
			log.V(1).Info("cannot update status on modified object, requeuing key for reprocessing")
			return ctrl.Result{RequeueAfter: 500 * time.Millisecond}, nil
		}

		log.Error(err, "failed to update cluster status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// manageFinalization will add a finalizer to new ray cluster resources if it's
// absent and remove it during a delete request after performing the required
// finalization steps.
func (r *RayClusterReconciler) manageFinalization(ctx context.Context, rc *dcv1alpha1.RayCluster) (bool, error) {
	log := r.Log.FromContext(ctx)
	registered := controllerutil.ContainsFinalizer(rc, DistributedComputeFinalizer)

	if rc.GetDeletionTimestamp().IsZero() && !registered {
		log.V(1).Info("registering finalizer", "name", DistributedComputeFinalizer)
		controllerutil.AddFinalizer(rc, DistributedComputeFinalizer)

		if err := r.Update(ctx, rc); err != nil {
			log.Error(err, "failed to register finalizer")
			return false, err
		}

		return true, nil
	}

	if rc.GetDeletionTimestamp() != nil && registered {
		log.Info("executing finalization steps")
		if err := r.deleteExternalStorage(ctx, rc); err != nil {
			log.Error(err, "failed to clean up storage")
			return false, err
		}

		log.V(1).Info("removing finalizer", "name", DistributedComputeFinalizer)
		controllerutil.RemoveFinalizer(rc, DistributedComputeFinalizer)

		if err := r.Update(ctx, rc); err != nil {
			log.Error(err, "failed to remove finalizer")
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// reconcileResources manages the creation and updates of resources that
// collectively comprise a Ray cluster. Each resource is controlled by a parent
// RayCluster object so that full cleanup occurs during a delete operation.
func (r *RayClusterReconciler) reconcileResources(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	if err := r.reconcileIstio(ctx, rc); err != nil {
		return err
	}
	if err := r.reconcileServiceAccount(ctx, rc); err != nil {
		return err
	}
	if err := r.reconcileServices(ctx, rc); err != nil {
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

	return r.reconcileStatefulSets(ctx, rc)
}

func (r *RayClusterReconciler) reconcileIstio(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	if !r.IstioEnabled {
		return nil
	}

	peerAuth := istio.NewPeerAuthentication(&istio.PeerAuthInfo{
		Name:      ray.InstanceObjectName(rc.Name, ray.ComponentNone),
		Namespace: rc.Namespace,
		Labels:    ray.MetadataLabels(rc),
		Selector:  ray.SelectorLabels(rc),
		Mode:      rc.Spec.IstioConfig.MutualTLSMode,
	})

	if rc.Spec.IstioConfig.MutualTLSMode == "" {
		return r.deleteIfExists(ctx, peerAuth)
	}
	if err := r.createOrUpdateOwnedResource(ctx, rc, peerAuth); err != nil {
		return fmt.Errorf("failed to reconcile peer authentication: %w", err)
	}

	return nil
}

// reconcileServiceAccount creates a new dedicated service account for a Ray
// cluster unless a different service account name is provided in the spec.
func (r *RayClusterReconciler) reconcileServiceAccount(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	if rc.Spec.ServiceAccount.Name != "" {
		return nil
	}

	sa := ray.NewServiceAccount(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, sa); err != nil {
		return fmt.Errorf("failed to reconcile service account: %w", err)
	}

	return nil
}

// reconcileServices creates services that point to head and worker pods and
// applies updates when the parent CR changes.
func (r *RayClusterReconciler) reconcileServices(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	svc := ray.NewClientService(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, svc); err != nil {
		return fmt.Errorf("failed to reconcile client service: %w", err)
	}

	svc = ray.NewHeadlessHeadService(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, svc); err != nil {
		return fmt.Errorf("failed to reconcile headless head service: %w", err)
	}

	svc = ray.NewHeadlessWorkerService(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, svc); err != nil {
		return fmt.Errorf("failed to reconcile headless worker service: %w", err)
	}

	return nil
}

// reconcileNetworkPolicies optionally creates network policies that control
// traffic flow between cluster nodes and external clients. Existing network
// policies will be deleted if enabled is set to false.
func (r RayClusterReconciler) reconcileNetworkPolicies(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	clusterNetpol := ray.NewClusterNetworkPolicy(rc)
	clientNetpol := ray.NewHeadClientNetworkPolicy(rc)
	dashboardNetpol := ray.NewHeadDashboardNetworkPolicy(rc)

	if util.BoolPtrIsNilOrFalse(rc.Spec.NetworkPolicy.Enabled) {
		return r.deleteIfExists(ctx, dashboardNetpol, clientNetpol, clusterNetpol)
	}

	if err := r.createOrUpdateOwnedResource(ctx, rc, clusterNetpol); err != nil {
		return fmt.Errorf("failed to reconcile cluster network policy: %w", err)
	}
	if err := r.createOrUpdateOwnedResource(ctx, rc, clientNetpol); err != nil {
		return fmt.Errorf("failed to reconcile head client network policy: %w", err)
	}
	if err := r.createOrUpdateOwnedResource(ctx, rc, dashboardNetpol); err != nil {
		return fmt.Errorf("failed to reconcile head dashboard network policy: %w", err)
	}

	return nil
}

// nolint:dupl
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
	if rc.Spec.Autoscaling == nil {
		// calling ray.NewHorizontalPodAutoscaler when autoscaling is nil will
		// result in error. so we leverage a shallow reference here instead.
		hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{
			ObjectMeta: ray.HorizontalPodAutoscalerObjectMeta(rc),
		}
		return r.deleteIfExists(ctx, hpa)
	}

	hpa, err := ray.NewHorizontalPodAutoscaler(rc)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, rc, hpa); err != nil {
		return fmt.Errorf("failed to reconcile horizontal pod autoscaler: %w", err)
	}

	return nil
}

// reconcileStatefulSets creates separate Ray head and worker stateful sets
// that will collectively comprise the execution agents of the cluster.
func (r *RayClusterReconciler) reconcileStatefulSets(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	head, err := ray.NewStatefulSet(rc, ray.ComponentHead, r.IstioEnabled)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, rc, head); err != nil {
		return fmt.Errorf("failed to create head stateful set: %w", err)
	}

	worker, err := ray.NewStatefulSet(rc, ray.ComponentWorker, r.IstioEnabled)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, rc, worker); err != nil {
		return fmt.Errorf("failed to create worker stateful set: %w", err)
	}

	return nil
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

	var gvks []schema.GroupVersionKind
	gvks, _, err := r.Scheme.ObjectKinds(controlled)
	if err != nil {
		return err
	}
	gvk := gvks[0]

	log := r.Log.FromContext(ctx)

	found := controlled.DeepCopyObject().(client.Object)
	if err = r.Get(ctx, client.ObjectKeyFromObject(controlled), found); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}

		if err = PatchAnnotator.SetLastAppliedAnnotation(controlled); err != nil {
			return err
		}

		log.Info("creating controlled object", "gvk", gvk, "object", controlled)
		return r.Create(ctx, controlled)
	}

	patchResult, err := PatchMaker.Calculate(found, controlled, PatchCalculateOpts...)
	if err != nil {
		return err
	}
	if patchResult.IsEmpty() {
		return nil
	}

	log.V(1).Info("applying patch to object", "gvk", gvk, "object", controlled, "patch", string(patchResult.Patch))
	if err = PatchAnnotator.SetLastAppliedAnnotation(controlled); err != nil {
		return err
	}

	controlled.SetResourceVersion(found.GetResourceVersion())
	if modified, ok := controlled.(*corev1.Service); ok {
		current := found.(*corev1.Service)
		modified.Spec.ClusterIP = current.Spec.ClusterIP
	}

	log.Info("updating controlled object", "gvk", gvk, "object", controlled)
	return r.Update(ctx, controlled)
}

// deleteIfExists will delete one or more Kubernetes objects if they exist.
func (r *RayClusterReconciler) deleteIfExists(ctx context.Context, objs ...client.Object) error {
	log := r.Log.FromContext(ctx)

	for _, obj := range objs {
		if err := r.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}

			return err
		}

		log.Info("deleting controlled object", "object", obj)
		if err := r.Delete(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

// updateStatus updates the RayCluster status subresource when changes occur.
func (r *RayClusterReconciler) updateStatus(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	mNodes, err := r.modifyStatusNodes(ctx, rc)
	if err != nil {
		return fmt.Errorf("cannot modify cluster status nodes: %w", err)
	}

	mWorkedFields, err := r.modifyStatusWorkerFields(ctx, rc)
	if err != nil {
		return fmt.Errorf("cannot modify cluster status worker fields: %w", err)
	}

	if mNodes || mWorkedFields {
		if err = r.Status().Update(ctx, rc); err != nil {
			return err
		}
	}

	return nil
}

// modifyStatusNodes will ensure the status contains an accurate list of all the pods in the cluster.
func (r *RayClusterReconciler) modifyStatusNodes(ctx context.Context, rc *dcv1alpha1.RayCluster) (bool, error) {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(rc.Namespace),
		client.MatchingLabels(ray.MetadataLabels(rc)),
	}
	if err := r.List(ctx, podList, listOpts...); err != nil {
		return false, fmt.Errorf("cannot list ray pods: %w", err)
	}

	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	sort.Strings(podNames)

	if reflect.DeepEqual(podNames, rc.Status.Nodes) {
		return false, nil
	}

	log := r.Log.FromContext(ctx)

	log.V(1).Info("modifying status", "path", ".status.nodes", "value", podNames)
	rc.Status.Nodes = podNames

	return true, nil
}

// modifyStatusWorkerFields syncs certain worker stateful set fields into the status.
func (r *RayClusterReconciler) modifyStatusWorkerFields(ctx context.Context, rc *dcv1alpha1.RayCluster) (bool, error) {
	worker, err := ray.NewStatefulSet(rc, ray.ComponentWorker, r.IstioEnabled)
	if err != nil {
		return false, err
	}

	err = r.Get(ctx, client.ObjectKeyFromObject(worker), worker)
	if client.IgnoreNotFound(err) != nil {
		return false, err
	}

	selector, err := metav1.LabelSelectorAsSelector(worker.Spec.Selector)
	if err != nil {
		return false, err
	}

	log := r.Log.FromContext(ctx)

	var modified bool
	if rc.Status.WorkerSelector != selector.String() {
		rc.Status.WorkerSelector = selector.String()
		modified = true

		log.V(1).Info("modifying status", "path", ".status.workerSelector", "value", rc.Status.WorkerSelector)
	}
	if rc.Status.WorkerReplicas != *worker.Spec.Replicas {
		rc.Status.WorkerReplicas = *worker.Spec.Replicas
		modified = true

		log.V(1).Info("modifying status", "path", ".status.workerReplicas", "value", rc.Status.WorkerReplicas)
	}

	return modified, nil
}

// deleteExternalStorage queries for all persistent volume claims belonging to
// a cluster instance using selector labels. this should find all the claims
// created by both the head and worker stateful sets.
func (r *RayClusterReconciler) deleteExternalStorage(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	log := r.Log.FromContext(ctx)

	ns := rc.Namespace
	labels := ray.SelectorLabels(rc)
	pvcList := &corev1.PersistentVolumeClaimList{}
	listOpts := []client.ListOption{
		client.InNamespace(ns),
		client.MatchingLabels(labels),
	}

	log.Info("querying for persistent volume claims", "namespace", ns, "labels", labels)
	if err := r.List(ctx, pvcList, listOpts...); err != nil {
		log.Error(err, "cannot list persistent volume claims")
		return err
	}

	for idx := range pvcList.Items {
		pvc := &pvcList.Items[idx]
		key := client.ObjectKeyFromObject(pvc)

		log.Info("deleting persistent volume claim", "claim", key)
		if err := r.Delete(ctx, pvc); err != nil {
			log.Error(err, "cannot delete persistent volume claim", "claim", key)
			return err
		}
	}

	return nil
}

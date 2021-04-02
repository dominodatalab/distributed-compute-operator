package controllers

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/logging"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/spark"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

// LastAppliedConfig is the annotation key used to store object state on owned components.
const SparkLastAppliedConfig = "distributed-compute-operator.dominodatalab.com/last-applied"

var (
	// SparkPatchAnnotator applies state annotations to owned components.
	SparkPatchAnnotator = patch.NewAnnotator(SparkLastAppliedConfig)
	// SparkPatchMaker calculates changes to state annotations on owned components.
	SparkPatchMaker = patch.NewPatchMaker(SparkPatchAnnotator)
)

// SparkClusterReconciler reconciles SparkCluster objects.
type SparkClusterReconciler struct {
	client.Client
	Log    logging.ContextLogger
	Scheme *runtime.Scheme
}

// nolint:dupl
// SetupWithManager creates and registers this controller with the manager.
func (r *SparkClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dcv1alpha1.SparkCluster{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&autoscalingv2beta2.HorizontalPodAutoscaler{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

const SparkFinalizerName = "distributed-compute.dominodatalab.com/dco-finalizer"

//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=sparkclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=sparkclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=sparkclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=list;watch
//+kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=create;update;list;watch
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=create;update;list;watch
//+kubebuilder:rbac:groups=apps,resources=configmaps,verbs=create;update;list;watch
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=create;update;delete;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=create;update;delete;list;watch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=create;update;delete;list;watch

// Reconcile implements state reconciliation logic for SparkCluster objects.
func (r *SparkClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx, log := r.setLogger(ctx, r.Log.WithValues("sparkcluster", req.NamespacedName))

	log.V(2).Info("reconciliation loop trigged")

	rc := &dcv1alpha1.SparkCluster{}
	if err := r.Get(ctx, req.NamespacedName, rc); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("resource not found, assuming object was deleted")
			return ctrl.Result{}, nil
		}

		log.Error(err, "failed to retrieve resource")
		return ctrl.Result{}, err
	}

	updated, err := r.processFinalizers(ctx, rc, log)
	if err != nil {
		log.Error(err, "failed to process finalizers")
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

func (r *SparkClusterReconciler) processFinalizers(ctx context.Context, sc *dcv1alpha1.SparkCluster, log logr.Logger) (bool, error) {
	// no finalizer and no deletion timestamp means this is a new object so we're going to set a finalizer
	containsFinalizer := controllerutil.ContainsFinalizer(sc, SparkFinalizerName)
	hasDeletionTimestamp := hasDeletionTimestamp(sc)
	if !containsFinalizer && !hasDeletionTimestamp {
		controllerutil.AddFinalizer(sc, SparkFinalizerName)
		err := r.Update(ctx, sc)
		if err != nil {
			log.Error(err, "failed to set finalizer")
			return false, err
		}
		return true, nil
		// if it has finalizer and has a deletion timestamp then we want to delete some stuff
	}
	if containsFinalizer && hasDeletionTimestamp {
		log.Info(fmt.Sprintf("%s has finalizer and deletion timestamp. looking for pvcs to delete", sc.Name))
		pvcs, err := r.getPvcsForDeletion(ctx, sc)
		if err != nil {
			return false, err
		}

		err = r.deletePvcs(ctx, sc, log, pvcs)
		if err != nil {
			return false, err
		}

		controllerutil.RemoveFinalizer(sc, SparkFinalizerName)
		err = r.Update(ctx, sc)
		if err != nil {
			log.Error(err, "unable to remove finalizer from resource")
			return false, err
		}
	}
	return false, nil
}

func (r *SparkClusterReconciler) deletePvcs(
	ctx context.Context,
	sc *dcv1alpha1.SparkCluster,
	log logr.Logger,
	pvcsToDelete []client.Object) error {
	if len(pvcsToDelete) > 0 {
		log.Info(fmt.Sprintf("deleting %d pvcs associated with %s", len(pvcsToDelete), sc.Name))
		err := r.deleteIfExists(ctx, pvcsToDelete...)
		if err != nil {
			log.Error(err, fmt.Sprintf("unable to delete %s pvcs", sc.Name))
			return err
		}
	}
	return nil
}

func (r *SparkClusterReconciler) getPvcsForDeletion(
	ctx context.Context,
	sc *dcv1alpha1.SparkCluster) ([]client.Object, error) {
	var pvcsToDelete []client.Object
	claims := &corev1.PersistentVolumeClaimList{}
	var selectors client.MatchingLabels = spark.SelectorLabels(sc)
	err := r.List(ctx, claims, client.InNamespace(sc.Namespace), selectors)
	if err != nil {
		return pvcsToDelete, err
	}
	for _, claim := range claims.Items {
		// dont bother deleting if its already been deleted
		if claim.DeletionTimestamp == nil {
			pvc := claim
			pvcsToDelete = append(pvcsToDelete, &pvc)
		}
	}
	return pvcsToDelete, err
}

func hasDeletionTimestamp(sc *dcv1alpha1.SparkCluster) bool {
	return sc.ObjectMeta.DeletionTimestamp != nil
}

// reconcileResources manages the creation and updates of resources that
// collectively comprise a Spark cluster. Each resource is controlled by a parent
// SparkCluster object so that full cleanup occurs during a delete operation.
func (r *SparkClusterReconciler) reconcileResources(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	if err := r.reconcileServiceAccount(ctx, sc); err != nil {
		return err
	}
	if err := r.reconcileHeadService(ctx, sc); err != nil {
		return err
	}
	if err := r.reconcileHeadlessService(ctx, sc); err != nil {
		return err
	}
	if err := r.reconcileNetworkPolicies(ctx, sc); err != nil {
		return err
	}
	if err := r.reconcilePodSecurityPolicyRBAC(ctx, sc); err != nil {
		return err
	}
	if err := r.reconcileAutoscaler(ctx, sc); err != nil {
		return err
	}
	if err := r.reconcileConfigMap(ctx, sc); err != nil {
		return err
	}

	return r.reconcileStatefulSets(ctx, sc)
}

func (r *SparkClusterReconciler) reconcileConfigMap(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	if sc.Spec.Master.FrameworkConfig == nil && sc.Spec.Worker.FrameworkConfig == nil {
		return nil
	}

	cm := spark.NewConfigMap(sc)
	if err := r.createOrUpdateOwnedResource(ctx, sc, cm); err != nil {
		return fmt.Errorf("failed to reconcile configmap: %w", err)
	}

	return nil
}

// reconcileServiceAccount creates a new dedicated service account for a Spark
// cluster unless a different service account name is provided in the spec.
func (r *SparkClusterReconciler) reconcileServiceAccount(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	if sc.Spec.ServiceAccountName != "" {
		return nil
	}

	sa := spark.NewServiceAccount(sc)
	if err := r.createOrUpdateOwnedResource(ctx, sc, sa); err != nil {
		return fmt.Errorf("failed to reconcile service account: %w", err)
	}

	return nil
}

// reconcileHeadService creates a service that points to the head Spark pod and
// applies updates when the parent CR changes.
func (r *SparkClusterReconciler) reconcileHeadService(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	svc := spark.NewMasterService(sc)
	if err := r.createOrUpdateOwnedResource(ctx, sc, svc); err != nil {
		return fmt.Errorf("failed to reconcile head service: %w", err)
	}

	return nil
}

// reconcileHeadService creates a service that points to the head Spark pod and
// applies updates when the parent CR changes.
func (r *SparkClusterReconciler) reconcileHeadlessService(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	svc := spark.NewHeadlessService(sc)
	if err := r.createOrUpdateOwnedResource(ctx, sc, svc); err != nil {
		return fmt.Errorf("failed to reconcile headless service: %w", err)
	}

	return nil
}

// reconcileNetworkPolicies optionally creates network policies that control
// traffic flow between cluster nodes and external clients.
func (r SparkClusterReconciler) reconcileNetworkPolicies(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	headNetpol := spark.NewHeadClientNetworkPolicy(sc)
	clusterNetpol := spark.NewClusterNetworkPolicy(sc)
	dashboardNetpol := spark.NewHeadDashboardNetworkPolicy(sc)

	if !util.BoolPtrIsTrue(sc.Spec.NetworkPolicy.Enabled) {
		return r.deleteIfExists(ctx, headNetpol, clusterNetpol)
	}

	if err := r.createOrUpdateOwnedResource(ctx, sc, clusterNetpol); err != nil {
		return fmt.Errorf("failed to reconcile cluster network policy: %w", err)
	}

	if err := r.createOrUpdateOwnedResource(ctx, sc, headNetpol); err != nil {
		return fmt.Errorf("failed to reconcile head network policy: %w", err)
	}

	if err := r.createOrUpdateOwnedResource(ctx, sc, dashboardNetpol); err != nil {
		return fmt.Errorf("failed to reconcile dashboard network policy: %w", err)
	}

	return nil
}

// nolint:dupl
// reconcilePodSecurityPolicyRBAC optionally creates a role and role binding
// that allows the Spark pods to "use" the specified pod security policy.
func (r *SparkClusterReconciler) reconcilePodSecurityPolicyRBAC(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	role, binding := spark.NewPodSecurityPolicyRBAC(sc)

	if sc.Spec.PodSecurityPolicy == "" {
		return r.deleteIfExists(ctx, role, binding)
	}

	err := r.Get(ctx, types.NamespacedName{Name: sc.Spec.PodSecurityPolicy}, &policyv1beta1.PodSecurityPolicy{})
	if err != nil {
		return fmt.Errorf("cannot verify pod security policy: %w", err)
	}

	if err := r.createOrUpdateOwnedResource(ctx, sc, role); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	if err := r.createOrUpdateOwnedResource(ctx, sc, binding); err != nil {
		return fmt.Errorf("failed to create role binding: %w", err)
	}

	return nil
}

// reconcileAutoscaler optionally creates a horizontal pod autoscaler that
// targets Spark worker pods.
func (r *SparkClusterReconciler) reconcileAutoscaler(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	if sc.Spec.Autoscaling == nil {
		hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{
			ObjectMeta: spark.HorizontalPodAutoscalerObjectMeta(sc),
		}
		return r.deleteIfExists(ctx, hpa)
	}

	hpa, err := spark.NewHorizontalPodAutoscaler(sc)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, sc, hpa); err != nil {
		return fmt.Errorf("failed to reconcile horizontal pod autoscaler: %w", err)
	}

	return nil
}

// reconcileStatefulSets creates separate Spark head and worker statefulsets that
// will collectively comprise the execution agents of the cluster.
func (r *SparkClusterReconciler) reconcileStatefulSets(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	head, err := spark.NewStatefulSet(sc, spark.ComponentMaster)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, sc, head); err != nil {
		return fmt.Errorf("failed to create head deployment: %w", err)
	}

	worker, err := spark.NewStatefulSet(sc, spark.ComponentWorker)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, sc, worker); err != nil {
		return fmt.Errorf("failed to create worker deployment: %w", err)
	}

	selector, err := metav1.LabelSelectorAsSelector(worker.Spec.Selector)
	if err != nil {
		return err
	}

	log := r.getLogger(ctx)

	// update autoscaling fields
	var updateStatus bool
	if sc.Status.WorkerReplicas != *worker.Spec.Replicas {
		sc.Status.WorkerReplicas = *worker.Spec.Replicas
		updateStatus = true

		log.V(1).Info("updating status", "path", ".status.workerReplicas", "value", sc.Status.WorkerReplicas)
	}
	if sc.Status.WorkerSelector != selector.String() {
		sc.Status.WorkerSelector = selector.String()
		updateStatus = true

		log.V(1).Info("updating status", "path", ".status.workerSelector", "value", sc.Status.WorkerSelector)
	}

	if updateStatus {
		err = r.Status().Update(ctx, sc)
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
func (r *SparkClusterReconciler) createOrUpdateOwnedResource(ctx context.Context, owner metav1.Object, controlled client.Object) error {
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

	patchResult, err := PatchMaker.Calculate(found, controlled, patch.IgnoreStatusFields())
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
func (r *SparkClusterReconciler) deleteIfExists(ctx context.Context, objs ...client.Object) error {
	log := r.getLogger(ctx)

	for _, obj := range objs {
		err := r.Get(ctx, client.ObjectKeyFromObject(obj), obj)

		if apierrors.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}

		log.Info("deleting controlled object", "object", obj)
		if err = r.Delete(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

// updateStatus with a list of pods from both the head and worker deployments.
func (r *SparkClusterReconciler) updateStatus(ctx context.Context, sc *dcv1alpha1.SparkCluster) error {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(sc.Namespace),
		client.MatchingLabels(spark.MetadataLabels(sc)),
	}
	if err := r.List(ctx, podList, listOpts...); err != nil {
		return fmt.Errorf("cannot list spark pods: %w", err)
	}

	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	sort.Strings(podNames)

	if reflect.DeepEqual(podNames, sc.Status.Nodes) {
		return nil
	}
	sc.Status.Nodes = podNames

	if err := r.Status().Update(ctx, sc); err != nil {
		return fmt.Errorf("cannot update spark status nodes: %w", err)
	}

	return nil
}

type loggerKeyType int

const loggerKey loggerKeyType = iota

func (r *SparkClusterReconciler) setLogger(ctx context.Context, logger logr.Logger) (context.Context, logr.Logger) {
	return context.WithValue(ctx, loggerKey, logger), logger
}

func (r *SparkClusterReconciler) getLogger(ctx context.Context) logr.Logger {
	if ctx == nil {
		return r.Log
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(logr.Logger); ok {
		return ctxLogger
	}

	return r.Log
}

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
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/spark"
)

//
//// LastAppliedConfig is the annotation key used to store object state on owned components.
//const LastAppliedConfig = "distributed-compute-operator.dominodatalab.com/last-applied"
//
//var (
//	// PatchAnnotator applies state annotations to owned components.
//	PatchAnnotator = patch.NewAnnotator(LastAppliedConfig)
//	// PatchMaker calculates changes to state annotations on owned components.
//	PatchMaker = patch.NewPatchMaker(PatchAnnotator)
//)

// SparkClusterReconciler reconciles SparkCluster objects.
type SparkClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager creates and registers this controller with the manager.
func (r *SparkClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dcv1alpha1.SparkCluster{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.Deployment{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&autoscalingv2beta2.HorizontalPodAutoscaler{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=sparkclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=sparkclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=sparkclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=list;watch
//+kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=create;update;list;watch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=create;update;list;watch
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

	if err := r.reconcileResources(ctx, rc); err != nil {
		log.Error(err, "failed to reconcile cluster resources")
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
// collectively comprise a Spark cluster. Each resource is controlled by a parent
// SparkCluster object so that full cleanup occurs during a delete operation.
func (r *SparkClusterReconciler) reconcileResources(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	if err := r.reconcileServiceAccount(ctx, rc); err != nil {
		return err
	}
	if err := r.reconcileHeadService(ctx, rc); err != nil {
		return err
	}
	if err := r.reconcileHeadLessService(ctx, rc); err != nil {
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

// reconcileServiceAccount creates a new dedicated service account for a Spark
// cluster unless a different service account name is provided in the spec.
func (r *SparkClusterReconciler) reconcileServiceAccount(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	if rc.Spec.ServiceAccountName != "" {
		return nil
	}

	sa := spark.NewServiceAccount(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, sa); err != nil {
		return fmt.Errorf("failed to reconcile service account: %w", err)
	}

	return nil
}

// reconcileHeadService creates a service that points to the head Spark pod and
// applies updates when the parent CR changes.
func (r *SparkClusterReconciler) reconcileHeadService(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	svc := spark.NewHeadService(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, svc); err != nil {
		return fmt.Errorf("failed to reconcile head service: %w", err)
	}

	return nil
}

// reconcileHeadService creates a service that points to the head Spark pod and
// applies updates when the parent CR changes.
func (r *SparkClusterReconciler) reconcileHeadLessService(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	svc := spark.NewHeadlessService(rc)
	if err := r.createOrUpdateOwnedResource(ctx, rc, svc); err != nil {
		return fmt.Errorf("failed to reconcile headless service: %w", err)
	}

	return nil
}

// reconcileNetworkPolicies optionally creates network policies that control
// traffic flow between cluster nodes and external clients.
func (r SparkClusterReconciler) reconcileNetworkPolicies(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	headNetpol := spark.NewHeadNetworkPolicy(rc)
	clusterNetpol := spark.NewClusterNetworkPolicy(rc)

	if rc.Spec.EnableNetworkPolicy == nil || !*rc.Spec.EnableNetworkPolicy {
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
// that allows the Spark pods to "use" the specified pod security policy.
func (r *SparkClusterReconciler) reconcilePodSecurityPolicyRBAC(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	role, binding := spark.NewPodSecurityPolicyRBAC(rc)

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
// targets Spark worker pods.
func (r *SparkClusterReconciler) reconcileAutoscaler(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	if rc.Spec.Autoscaling == nil {
		hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{
			ObjectMeta: spark.HorizontalPodAutoscalerObjectMeta(rc),
		}
		return r.deleteIfExists(ctx, hpa)
	}

	hpa, err := spark.NewHorizontalPodAutoscaler(rc)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, rc, hpa); err != nil {
		return fmt.Errorf("failed to reconcile horizontal pod autoscaler: %w", err)
	}

	return nil
}

// reconcileDeployments creates separate Spark head and worker deployments that
// will collectively comprise the execution agents of the cluster.
func (r *SparkClusterReconciler) reconcileDeployments(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	head, err := spark.NewStatefulSet(rc, spark.ComponentHead)
	if err != nil {
		return err
	}
	if err = r.createOrUpdateOwnedResource(ctx, rc, head); err != nil {
		return fmt.Errorf("failed to create head deployment: %w", err)
	}

	worker, err := spark.NewStatefulSet(rc, spark.ComponentWorker)
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

	log := r.getLogger(ctx)

	// update autoscaling fields
	var updateStatus bool
	if rc.Status.WorkerReplicas != *worker.Spec.Replicas {
		rc.Status.WorkerReplicas = *worker.Spec.Replicas
		updateStatus = true

		log.V(1).Info("updating status", "path", ".status.workerReplicas", "value", rc.Status.WorkerReplicas)
	}
	if rc.Status.WorkerSelector != selector.String() {
		rc.Status.WorkerSelector = selector.String()
		updateStatus = true

		log.V(1).Info("updating status", "path", ".status.workerSelector", "value", rc.Status.WorkerSelector)
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
func (r *SparkClusterReconciler) createOrUpdateOwnedResource(ctx context.Context, owner metav1.Object, controlled client.Object) error {
	if err := ctrl.SetControllerReference(owner, controlled, r.Scheme); err != nil {
		return err
	}

	log := r.getLogger(ctx)
	found := controlled.DeepCopyObject().(client.Object)
	err := r.Get(ctx, client.ObjectKeyFromObject(controlled), found)

	if apierrors.IsNotFound(err) {
		if err = PatchAnnotator.SetLastAppliedAnnotation(controlled); err != nil {
			return err
		}

		log.Info("creating controlled object", "object", controlled)
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

	log.V(1).Info("applying patch to object", "object", controlled, "patch", string(patchResult.Patch))
	if err = PatchAnnotator.SetLastAppliedAnnotation(controlled); err != nil {
		return err
	}

	controlled.SetResourceVersion(found.GetResourceVersion())
	if modified, ok := controlled.(*corev1.Service); ok {
		current := found.(*corev1.Service)
		modified.Spec.ClusterIP = current.Spec.ClusterIP
	}

	log.Info("updating controlled object", "object", controlled)
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

//nolint
// updateStatus with a list of pods from both the head and worker deployments.
func (r *SparkClusterReconciler) updateStatus(ctx context.Context, rc *dcv1alpha1.SparkCluster) error {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(rc.Namespace),
		client.MatchingLabels(spark.MetadataLabels(rc)),
	}
	if err := r.List(ctx, podList, listOpts...); err != nil {
		return fmt.Errorf("cannot list spark pods: %w", err)
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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/ray"
)

// RayClusterReconciler reconciles RayCluster objects.
type RayClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=rayclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=distributed-compute.dominodatalab.com,resources=rayclusters/finalizers,verbs=update

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

	if err := r.ProcessResources(ctx, rc); err != nil {
		log.Error(err, "Failed to reconcile cluster resources")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// ProcessResources manages the creation and updates of resources that collectively comprise a Ray cluster.
// Each resource is controlled by a parent RayCluster object so that full cleanup occurs during a delete operation.
func (r *RayClusterReconciler) ProcessResources(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	// manage supporting resources
	sa := ray.NewServiceAccount(rc)
	if err := r.createOwnedResource(ctx, rc, sa); err != nil {
		return fmt.Errorf("failed to create service account: %w", err)
	}

	// TODO: service, network policy, pod security policy

	// manage deployments

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RayClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dcv1alpha1.RayCluster{}).
		Complete(r)
}

func (r *RayClusterReconciler) createOwnedResource(ctx context.Context, owner v1.Object, controlled client.Object) error {
	objKey := types.NamespacedName{Name: controlled.GetName(), Namespace: controlled.GetNamespace()}
	err := r.Get(ctx, objKey, controlled)

	if err == nil {
		return nil
	}
	if !apierrors.IsNotFound(err) {
		return err
	}
	if err := ctrl.SetControllerReference(owner, controlled, r.Scheme); err != nil {
		return err
	}

	return r.Create(ctx, controlled)
}

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	if err := r.processResources(ctx, rc); err != nil {
		log.Error(err, "Failed to reconcile cluster resources")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RayClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dcv1alpha1.RayCluster{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.Deployment{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&autoscalingv2beta2.HorizontalPodAutoscaler{}).
		Complete(r)
}

// processResources manages the creation and updates of resources that
// collectively comprise a Ray cluster. Each resource is controlled by a parent
// RayCluster object so that full cleanup occurs during a delete operation.
func (r *RayClusterReconciler) processResources(ctx context.Context, rc *dcv1alpha1.RayCluster) error {
	sa := ray.NewServiceAccount(rc)
	if err := r.createOwnedResource(ctx, rc, sa); err != nil {
		return fmt.Errorf("failed to create service account: %w", err)
	}

	svc := ray.NewHeadService(rc)
	if err := r.createOwnedResource(ctx, rc, svc); err != nil {
		return fmt.Errorf("failed to create head service: %w", err)
	}

	if rc.Spec.EnableNetworkPolicy {
		netpol := ray.NewClusterNetworkPolicy(rc)
		if err := r.createOwnedResource(ctx, rc, netpol); err != nil {
			return fmt.Errorf("failed to create cluster network policy: %w", err)
		}

		netpol = ray.NewHeadNetworkPolicy(rc)
		if err := r.createOwnedResource(ctx, rc, netpol); err != nil {
			return fmt.Errorf("failed to create head network policy: %w", err)
		}
	}

	if rc.Spec.PodSecurityPolicy != "" {
		err := r.Get(ctx, types.NamespacedName{Name: rc.Spec.PodSecurityPolicy}, &policyv1beta1.PodSecurityPolicy{})
		if err != nil {
			return fmt.Errorf("cannot verify pod security policy: %w", err)
		}

		role, binding := ray.NewPodSecurityPolicyRBAC(rc)

		if err := r.createOwnedResource(ctx, rc, role); err != nil {
			return fmt.Errorf("failed to create role: %w", err)
		}

		if err := r.createOwnedResource(ctx, rc, binding); err != nil {
			return fmt.Errorf("failed to create role binding: %w", err)
		}
	}

	if rc.Spec.Autoscaling != nil {
		hpa := ray.NewHorizontalPodAutoscaler(rc)
		if err := r.createOwnedResource(ctx, rc, hpa); err != nil {
			return fmt.Errorf("failed to create horizontal pod autoscaler: %w", err)
		}
	}

	// manage deployments

	head, err := ray.NewDeployment(rc, ray.ComponentHead)
	if err != nil {
		return err
	}
	if err = r.createOwnedResource(ctx, rc, head); err != nil {
		return fmt.Errorf("failed to create head deployment: %w", err)
	}

	worker, err := ray.NewDeployment(rc, ray.ComponentWorker)
	if err != nil {
		return err
	}
	if err = r.createOwnedResource(ctx, rc, worker); err != nil {
		return fmt.Errorf("failed to create worker deployment: %w", err)
	}

	return nil
}

// createOwnedResource should be used to create namespace-scoped object.
// The CR will become the "owner" of the "controlled" object and cleanup will
// occur automatically when the CR is deleted.
func (r *RayClusterReconciler) createOwnedResource(ctx context.Context, owner metav1.Object, controlled client.Object) error {
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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var rayclusterlog = logf.Log.WithName("raycluster-resource")

func (r *RayCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-raycluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=create;update,versions=v1alpha1,name=mraycluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &RayCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *RayCluster) Default() {
	rayclusterlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// +kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-raycluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=create;update,versions=v1alpha1,name=vraycluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &RayCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *RayCluster) ValidateCreate() error {
	rayclusterlog.Info("validate create", "name", r.Name)

	return r.validateRayCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *RayCluster) ValidateUpdate(old runtime.Object) error {
	rayclusterlog.Info("validate update", "name", r.Name)

	return r.validateRayCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *RayCluster) ValidateDelete() error {
	return nil
}

func (r *RayCluster) validateRayCluster() error {
	return nil
}

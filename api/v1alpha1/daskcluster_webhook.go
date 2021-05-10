package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var daskclusterlog = logf.Log.WithName("daskcluster-resource")

func (r *DaskCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-daskcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=daskclusters,verbs=create;update,versions=v1alpha1,name=mdaskcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &DaskCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *DaskCluster) Default() {
	daskclusterlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-daskcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=daskclusters,verbs=create;update,versions=v1alpha1,name=vdaskcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &DaskCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *DaskCluster) ValidateCreate() error {
	daskclusterlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *DaskCluster) ValidateUpdate(old runtime.Object) error {
	daskclusterlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DaskCluster) ValidateDelete() error {
	daskclusterlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

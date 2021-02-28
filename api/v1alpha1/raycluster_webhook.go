package v1alpha1

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var rayclusterlog = logf.Log.WithName("raycluster-resource")

// SetupWebhookWithManager creates and registers this webhook with the manager.
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
}

// +kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-raycluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=create;update,versions=v1alpha1,name=vraycluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &RayCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *RayCluster) ValidateCreate() error {
	rayclusterlog.Info("validate create", "name", r.Name)

	return r.validateRayCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *RayCluster) ValidateUpdate(old runtime.Object) error {
	rayclusterlog.Info("validate update", "name", r.Name)

	return r.validateRayCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
// Not used, just here for interface compliance.
func (r *RayCluster) ValidateDelete() error {
	return nil
}

func (r *RayCluster) validateRayCluster() error {
	var allErrs field.ErrorList

	if err := r.validateWorkerReplicas(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateObjectStoreMemoryBytes(); err != nil {
		allErrs = append(allErrs, err)
	}
	if errs := r.validatePorts(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateAutoscaler(); errs != nil {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "distributed-compute.dominodatalab.com", Kind: "RayCluster"},
		r.Name,
		allErrs,
	)
}

func (r *RayCluster) validateWorkerReplicas() *field.Error {
	replicas := r.Spec.Worker.Replicas
	if replicas >= 1 {
		return nil
	}

	return field.Invalid(
		field.NewPath("spec").Child("worker").Child("replicas"),
		replicas,
		"should be greater than or equal to 1",
	)
}

func (r *RayCluster) validateObjectStoreMemoryBytes() *field.Error {
	memBytes := r.Spec.ObjectStoreMemoryBytes

	if memBytes == nil || *memBytes >= 78643200 {
		return nil
	}

	return field.Invalid(
		field.NewPath("spec").Child("objectStoreMemoryBytes"),
		memBytes,
		"should be greater than or equal to 78643200",
	)
}

func (r *RayCluster) validatePorts() field.ErrorList {
	var errs field.ErrorList

	if err := r.validatePort(r.Spec.Port, field.NewPath("spec").Child("port")); err != nil {
		errs = append(errs, err)
	}

	for idx, port := range r.Spec.RedisShardPorts {
		name := fmt.Sprintf("redisShardPorts[%d]", idx)
		if err := r.validatePort(port, field.NewPath("spec").Child(name)); err != nil {
			errs = append(errs, err)
		}
	}

	if err := r.validatePort(r.Spec.ClientServerPort, field.NewPath("spec").Child("clientServerPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.ObjectManagerPort, field.NewPath("spec").Child("objectManagerPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.NodeManagerPort, field.NewPath("spec").Child("nodeManagerPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.DashboardPort, field.NewPath("spec").Child("dashboardPort")); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (r *RayCluster) validatePort(port int32, fldPath *field.Path) *field.Error {
	if port >= 0 && port <= 65353 {
		return nil
	}

	return field.Invalid(fldPath, port, "should be greater than or equal to 0 and less than or equal to 65353")
}

func (r *RayCluster) validateAutoscaler() field.ErrorList {
	var errs field.ErrorList

	as := r.Spec.Autoscaling
	if as == nil {
		return nil
	}

	fldPath := field.NewPath("spec").Child("autoscaling")

	if as.MinReplicas != nil {
		if *as.MinReplicas <= 0 {
			errs = append(errs, field.Invalid(
				fldPath.Child("minReplicas"),
				as.MinReplicas,
				"should be greater than 0",
			))
		}

		if *as.MinReplicas > as.MaxReplicas {
			errs = append(errs, field.Invalid(
				fldPath.Child("maxReplicas"),
				as.MaxReplicas,
				"should be greater than spec.autoscaling.minReplicas",
			))
		}
	}

	if as.MaxReplicas <= 0 {
		errs = append(errs, field.Invalid(
			fldPath.Child("maxReplicas"),
			as.MaxReplicas,
			"should be greater than minReplicas",
		))
	}

	if as.AverageUtilization <= 0 || as.AverageUtilization > 100 {
		errs = append(errs, field.Invalid(
			fldPath.Child("averageUtilization"),
			as.AverageUtilization,
			"should be greater than 0 and less than or equal to 100",
		))
	}

	if as.ScaleDownStabilizationWindowSeconds != nil && *as.ScaleDownStabilizationWindowSeconds < 0 {
		errs = append(errs, field.Invalid(
			fldPath.Child("scaleDownStabilizationWindowSeconds"),
			as.ScaleDownStabilizationWindowSeconds,
			"should be greater than or equal to 0",
		))
	}

	return errs
}

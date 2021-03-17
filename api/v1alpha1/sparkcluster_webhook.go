package v1alpha1

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	minSparkValidPort int32 = 1024
	maxSparkValidPort int32 = 65535
)

var (
	defaultSparkDashboardPort             int32 = 8265
	defaultSparkEnableNetworkPolicy             = pointer.BoolPtr(true)
	defaultSparkWorkerReplicas                  = pointer.Int32Ptr(1)
	defaultSparkHttpPort                  int32 = 80
	defaultSparkClusterPort               int32 = 7077
	defaultSparkNetworkPolicyClientLabels       = map[string]string{
		"spark-client": "true",
	}
	defaultSparkImage = &OCIImageDefinition{
		Repository: "bitnami/spark",
		Tag:        "3.0.2-debian-10-r0",
	}
)

// logger is for webhook logging.
var sparkLogger = logf.Log.WithName("webhooks").WithName("SparkCluster")

// SetupWebhookWithManager creates and registers this webhook with the manager.
func (r *SparkCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-sparkcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=sparkclusters,verbs=create;update,versions=v1alpha1,name=msparkcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &SparkCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *SparkCluster) Default() {
	log := sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(r))
	log.Info("applying defaults")

	if r.Spec.ClusterPort == 0 {
		log.Info("setting default cluster port", "value", defaultSparkClusterPort)
		r.Spec.ClusterPort = defaultSparkClusterPort
	}
	if r.Spec.DashboardPort == 0 {
		log.Info("setting default dashboard port", "value", defaultSparkDashboardPort)
		r.Spec.DashboardPort = defaultSparkDashboardPort
	}
	if r.Spec.NetworkPolicy.Enabled == nil {
		log.Info("setting enable network policy flag", "value", *defaultSparkEnableNetworkPolicy)
		r.Spec.EnableNetworkPolicy = defaultSparkEnableNetworkPolicy
	}
	if r.Spec.NetworkPolicy.ClientServerLabels == nil {
		log.Info("setting default network policy client labels", "value", defaultSparkNetworkPolicyClientLabels)
		r.Spec.NetworkPolicy.ClientServerLabels = defaultSparkNetworkPolicyClientLabels
	}
	if r.Spec.NetworkPolicy.DashboardLabels == nil {
		log.Info("setting default network policy dashboard labels", "value", defaultSparkNetworkPolicyClientLabels)
		r.Spec.NetworkPolicy.DashboardLabels = defaultSparkNetworkPolicyClientLabels
	}
	if r.Spec.Worker.Replicas == nil {
		log.Info("setting default worker replicas", "value", *defaultSparkWorkerReplicas)
		r.Spec.Worker.Replicas = defaultSparkWorkerReplicas
	}

	if r.Spec.Image == nil {
		log.Info("setting default image", "value", *defaultSparkImage)
		r.Spec.Image = defaultSparkImage
	} else {
		if r.Spec.Image.Repository == "" {
			log.Info("setting default image repository", "value", defaultSparkImage.Repository)
			r.Spec.Image.Repository = defaultSparkImage.Repository
		}
		if r.Spec.Image.Tag == "" {
			log.Info("setting default image tag", "value", defaultSparkImage.Tag)
			r.Spec.Image.Tag = defaultSparkImage.Tag
		}
	}
}

//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-sparkcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=sparkclusters,verbs=create;update,versions=v1alpha1,name=vsparkcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &SparkCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *SparkCluster) ValidateCreate() error {
	sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(r)).Info("validating create")

	return r.validateSparkCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *SparkCluster) ValidateUpdate(old runtime.Object) error {
	sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(r)).Info("validating update")

	return r.validateSparkCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
// Not used, just here for interface compliance.
func (r *SparkCluster) ValidateDelete() error {
	return nil
}

func (r *SparkCluster) validateSparkCluster() error {
	var allErrs field.ErrorList

	if err := r.validateWorkerReplicas(); err != nil {
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
		schema.GroupKind{Group: "distributed-compute.dominodatalab.com", Kind: "SparkCluster"},
		r.Name,
		allErrs,
	)
}

func (r *SparkCluster) validateWorkerReplicas() *field.Error {
	replicas := r.Spec.Worker.Replicas
	if replicas == nil || *replicas >= 0 {
		return nil
	}

	return field.Invalid(
		field.NewPath("spec").Child("worker").Child("replicas"),
		replicas,
		"should be greater than or equal to 0",
	)
}

func (r *SparkCluster) validatePorts() field.ErrorList {
	var errs field.ErrorList

	if err := r.validatePort(r.Spec.ClusterPort, field.NewPath("spec").Child("clusterPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.DashboardPort, field.NewPath("spec").Child("dashboardPort")); err != nil {
		errs = append(errs, err)
	}

	// TODO: add validation to prevent port values overlap

	return errs
}

func (r *SparkCluster) validatePort(port int32, fldPath *field.Path) *field.Error {
	if port < minSparkValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be greater than or equal to %d", minSparkValidPort))
	}
	if port > maxSparkValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be less than or equal to %d", maxSparkValidPort))
	}

	return nil
}

func (r *SparkCluster) validateAutoscaler() field.ErrorList {
	var errs field.ErrorList

	as := r.Spec.Autoscaling
	if as == nil {
		return nil
	}

	fldPath := field.NewPath("spec").Child("autoscaling")

	if as.MinReplicas != nil {
		if *as.MinReplicas < 1 {
			errs = append(errs, field.Invalid(
				fldPath.Child("minReplicas"),
				as.MinReplicas,
				"must be greater than or equal to 1",
			))
		}

		if *as.MinReplicas > as.MaxReplicas {
			errs = append(errs, field.Invalid(
				fldPath.Child("maxReplicas"),
				as.MaxReplicas,
				"cannot be less than spec.autoscaling.minReplicas",
			))
		}
	}

	if as.MaxReplicas < 1 {
		errs = append(errs, field.Invalid(
			fldPath.Child("maxReplicas"),
			as.MaxReplicas,
			"must be greater than or equal to 1",
		))
	}

	if as.AverageCPUUtilization != nil && *as.AverageCPUUtilization <= 0 {
		errs = append(errs, field.Invalid(
			fldPath.Child("averageUtilization"),
			as.AverageCPUUtilization,
			"must be greater than 0",
		))
	}

	if as.ScaleDownStabilizationWindowSeconds != nil && *as.ScaleDownStabilizationWindowSeconds < 0 {
		errs = append(errs, field.Invalid(
			fldPath.Child("scaleDownStabilizationWindowSeconds"),
			as.ScaleDownStabilizationWindowSeconds,
			"must be greater than or equal to 0",
		))
	}

	return errs
}

package v1alpha1

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
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
	minValidPort int32 = 1024
	maxValidPort int32 = 65535
)

var (
	defaultPort                int32 = 6379
	defaultRedisShardPorts           = []int32{6380, 6381}
	defaultClientServerPort    int32 = 10001
	defaultObjectManagerPort   int32 = 2384
	defaultNodeManagerPort     int32 = 2385
	defaultDashboardPort       int32 = 8265
	defaultEnableDashboard           = pointer.BoolPtr(true)
	defaultEnableNetworkPolicy       = pointer.BoolPtr(true)
	defaultWorkerReplicas            = pointer.Int32Ptr(1)

	defaultNetworkPolicyClientLabels = []map[string]string{
		{"ray-client": "true"},
	}

	defaultImage = &OCIImageDefinition{
		Repository: "rayproject/ray",
		Tag:        "1.2.0-cpu",
	}
)

// logger is for webhook logging.
var logger = logf.Log.WithName("webhooks").WithName("RayCluster")

// SetupWebhookWithManager creates and registers this webhook with the manager.
func (r *RayCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-raycluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=create;update,versions=v1alpha1,name=mraycluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &RayCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *RayCluster) Default() {
	log := logger.WithValues("raycluster", client.ObjectKeyFromObject(r))
	log.Info("applying defaults")

	if r.Spec.Port == 0 {
		log.Info("setting default port", "value", defaultPort)
		r.Spec.Port = defaultPort
	}
	if r.Spec.RedisShardPorts == nil {
		log.Info("setting default redis shard ports", "value", defaultRedisShardPorts)
		r.Spec.RedisShardPorts = defaultRedisShardPorts
	}
	if r.Spec.ClientServerPort == 0 {
		log.Info("setting default client server port", "value", defaultClientServerPort)
		r.Spec.ClientServerPort = defaultClientServerPort
	}
	if r.Spec.ObjectManagerPort == 0 {
		log.Info("setting default object manager port", "value", defaultObjectManagerPort)
		r.Spec.ObjectManagerPort = defaultObjectManagerPort
	}
	if r.Spec.NodeManagerPort == 0 {
		log.Info("setting default node manager port", "value", defaultNodeManagerPort)
		r.Spec.NodeManagerPort = defaultNodeManagerPort
	}
	if r.Spec.DashboardPort == 0 {
		log.Info("setting default dashboard port", "value", defaultDashboardPort)
		r.Spec.DashboardPort = defaultDashboardPort
	}
	if r.Spec.EnableDashboard == nil {
		log.Info("setting enable dashboard flag", "value", *defaultEnableDashboard)
		r.Spec.EnableDashboard = defaultEnableDashboard
	}
	if r.Spec.EnableNetworkPolicy == nil {
		log.Info("setting enable network policy flag", "value", *defaultEnableNetworkPolicy)
		r.Spec.EnableNetworkPolicy = defaultEnableNetworkPolicy
	}
	if r.Spec.NetworkPolicyClientLabels == nil {
		log.Info("setting default network policy client labels", "value", defaultNetworkPolicyClientLabels)
		r.Spec.NetworkPolicyClientLabels = defaultNetworkPolicyClientLabels
	}
	if r.Spec.Worker.Replicas == nil {
		log.Info("setting default worker replicas", "value", *defaultWorkerReplicas)
		r.Spec.Worker.Replicas = defaultWorkerReplicas
	}

	if r.Spec.Image == nil {
		log.Info("setting default image", "value", *defaultImage)
		r.Spec.Image = defaultImage
	} else {
		if r.Spec.Image.Repository == "" {
			log.Info("setting default image repository", "value", defaultImage.Repository)
			r.Spec.Image.Repository = defaultImage.Repository
		}
		if r.Spec.Image.Tag == "" {
			log.Info("setting default image tag", "value", defaultImage.Tag)
			r.Spec.Image.Tag = defaultImage.Tag
		}
	}
}

//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-raycluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=create;update,versions=v1alpha1,name=vraycluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &RayCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *RayCluster) ValidateCreate() error {
	logger.WithValues("raycluster", client.ObjectKeyFromObject(r)).Info("validating create")

	return r.validateRayCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *RayCluster) ValidateUpdate(old runtime.Object) error {
	logger.WithValues("raycluster", client.ObjectKeyFromObject(r)).Info("validating update")

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
	if err := r.validateWorkerResourceRequestsCPU(); err != nil {
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
	if replicas == nil || *replicas >= 0 {
		return nil
	}

	return field.Invalid(
		field.NewPath("spec").Child("worker").Child("replicas"),
		replicas,
		"should be greater than or equal to 0",
	)
}

func (r *RayCluster) validateWorkerResourceRequestsCPU() *field.Error {
	if r.Spec.Autoscaling == nil {
		return nil
	}
	if _, ok := r.Spec.Worker.Resources.Requests[v1.ResourceCPU]; ok {
		return nil
	}

	return field.Required(
		field.NewPath("spec").Child("worker").Child("resources").Child("requests").Child("cpu"),
		"is mandatory when autoscaling is enabled",
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

	// TODO: add validation to prevent port values overlap

	return errs
}

func (r *RayCluster) validatePort(port int32, fldPath *field.Path) *field.Error {
	if port < minValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be greater than or equal to %d", minValidPort))
	}
	if port > maxValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be less than or equal to %d", maxValidPort))
	}

	return nil
}

func (r *RayCluster) validateAutoscaler() field.ErrorList {
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

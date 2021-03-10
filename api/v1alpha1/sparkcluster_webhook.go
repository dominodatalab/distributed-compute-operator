package v1alpha1

import (
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

//const (
//	minValidPort int32 = 1024
//	maxValidPort int32 = 65535
//)

var (
	//defaultPort                int32 = 6379
	//defaultRedisShardPorts           = []int32{6380, 6381}
	//defaultClientServerPort    int32 = 10001
	//defaultHttpPort   		   int32 = 80
	//defaultClusterPort         int32 = 7077
	//defaultDashboardPort       int32 = 8265
	//defaultEnableDashboard           = pointer.BoolPtr(true)
	//defaultEnableNetworkPolicy       = pointer.BoolPtr(true)
	//defaultWorkerReplicas            = pointer.Int32Ptr(1)
	//
	//defaultNetworkPolicyClientLabels = []map[string]string{
	//	{"spark-client": "true"},
	//}
	//
	//defaultImage = &OCIImageDefinition{
	//	Repository: "sparkproject/spark",
	//	Tag:        "1.2.0-cpu",
	//}
	defaultHttpPort            int32 = 80
	defaultClusterPort         int32 = 7077
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

	if r.Spec.Port == 0 {
		log.Info("setting default port", "value", defaultPort)
		r.Spec.Port = defaultPort
	}
	//if r.Spec.RedisShardPorts == nil {
	//	log.Info("setting default redis shard ports", "value", defaultRedisShardPorts)
	//	r.Spec.RedisShardPorts = defaultRedisShardPorts
	//}
	if r.Spec.ClientServerPort == 0 {
		log.Info("setting default client server port", "value", defaultClientServerPort)
		r.Spec.ClientServerPort = defaultClientServerPort
	}
	if r.Spec.HttpPort == 0 {
		log.Info("setting default http port", "value", defaultHttpPort)
		r.Spec.HttpPort = defaultHttpPort
	}
	if r.Spec.ClusterPort == 0 {
		log.Info("setting default cluster port", "value", defaultClusterPort)
		r.Spec.ClusterPort = defaultClusterPort
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

func (r *SparkCluster) validateObjectStoreMemoryBytes() *field.Error {
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

func (r *SparkCluster) validatePorts() field.ErrorList {
	var errs field.ErrorList

	if err := r.validatePort(r.Spec.Port, field.NewPath("spec").Child("port")); err != nil {
		errs = append(errs, err)
	}

	if err := r.validatePort(r.Spec.ClientServerPort, field.NewPath("spec").Child("clientServerPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.HttpPort, field.NewPath("spec").Child("httpPort")); err != nil {
		errs = append(errs, err)
	}
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
	if port < minValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be greater than or equal to %d", minValidPort))
	}
	if port > maxValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be less than or equal to %d", maxValidPort))
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

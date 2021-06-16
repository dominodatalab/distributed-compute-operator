package v1alpha1

import (
	"fmt"

	securityv1beta1 "istio.io/api/security/v1beta1"

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
	sparkMinValidPort int32 = 80
	sparkMaxValidPort int32 = 65535
)

var (
	sparkDefaultDashboardPort               int32 = 8265
	sparkDefaultClusterPort                 int32 = 7077
	sparkDefaultMasterWebPort               int32 = 80
	sparkDefaultWorkerWebPort               int32 = 8081
	sparkDefaultDriverUIPort                int32 = 4040
	sparkDefaultDriverPort                  int32 = 4041
	sparkDefaultDriverBlockManagerPort      int32 = 4042
	sparkDefaultDriverBlockManagerPortName        = "spark-block-manager-port"
	sparkDefaultDriverPortName                    = "spark-driver-port"
	sparkDefaultDriverUIPortName                  = "spark-ui-port"
	sparkDefaultEnableNetworkPolicy               = pointer.BoolPtr(true)
	sparkDefaultEnableExternalNetworkPolicy       = pointer.BoolPtr(false)
	sparkDefaultWorkerReplicas                    = pointer.Int32Ptr(1)
	sparkDefaultEnableDashboard                   = pointer.BoolPtr(true)
	sparkDefaultNetworkPolicyClientLabels         = map[string]string{
		"spark-client": "true",
	}
	sparkDefaultImage = &OCIImageDefinition{
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
// nolint:funlen,gocyclo
func (r *SparkCluster) Default() {
	log := sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(r))
	log.Info("applying defaults")

	if r.Spec.ClusterPort == 0 {
		log.Info("setting default cluster port", "value", sparkDefaultClusterPort)
		r.Spec.ClusterPort = sparkDefaultClusterPort
	}
	if r.Spec.TCPWorkerWebPort == 0 {
		log.Info("setting default worker web port", "value", sparkDefaultWorkerWebPort)
		r.Spec.TCPWorkerWebPort = sparkDefaultWorkerWebPort
	}
	if r.Spec.TCPMasterWebPort == 0 {
		log.Info("setting default master web port", "value", sparkDefaultMasterWebPort)
		r.Spec.TCPMasterWebPort = sparkDefaultMasterWebPort
	}
	if r.Spec.DashboardPort == 0 {
		log.Info("setting default dashboard port", "value", sparkDefaultDashboardPort)
		r.Spec.DashboardPort = sparkDefaultDashboardPort
	}
	if r.Spec.EnableDashboard == nil {
		log.Info("setting enable dashboard flag", "value", *sparkDefaultEnableDashboard)
		r.Spec.EnableDashboard = sparkDefaultEnableDashboard
	}
	if r.Spec.NetworkPolicy.Enabled == nil {
		log.Info("setting enable network policy flag", "value", *sparkDefaultEnableNetworkPolicy)
		r.Spec.NetworkPolicy.Enabled = sparkDefaultEnableNetworkPolicy
	}
	if r.Spec.NetworkPolicy.ExternalPolicyEnabled == nil {
		log.Info("setting enable external network policy flag", "value", *sparkDefaultEnableExternalNetworkPolicy)
		r.Spec.NetworkPolicy.ExternalPolicyEnabled = sparkDefaultEnableExternalNetworkPolicy
	}
	if r.Spec.NetworkPolicy.ClientServerLabels == nil {
		log.Info("setting default network policy client labels", "value", sparkDefaultNetworkPolicyClientLabels)
		r.Spec.NetworkPolicy.ClientServerLabels = sparkDefaultNetworkPolicyClientLabels
	}
	if r.Spec.NetworkPolicy.DashboardLabels == nil {
		log.Info("setting default network policy dashboard labels", "value", sparkDefaultNetworkPolicyClientLabels)
		r.Spec.NetworkPolicy.DashboardLabels = sparkDefaultNetworkPolicyClientLabels
	}
	if r.Spec.Worker.Replicas == nil {
		log.Info("setting default worker replicas", "value", *sparkDefaultWorkerReplicas)
		r.Spec.Worker.Replicas = sparkDefaultWorkerReplicas
	}
	if r.Spec.Driver.DriverPort == 0 {
		log.Info("setting default driver port", "value", sparkDefaultDriverPort)
		r.Spec.Driver.DriverPort = sparkDefaultDriverPort
	}
	if r.Spec.Driver.DriverPortName == "" {
		log.Info("setting default driver port name", "value", sparkDefaultDriverPortName)
		r.Spec.Driver.DriverPortName = sparkDefaultDriverPortName
	}
	if r.Spec.Driver.DriverBlockManagerPortName == "" {
		log.Info("setting default driver block manager port name", "value", sparkDefaultDriverBlockManagerPortName)
		r.Spec.Driver.DriverBlockManagerPortName = sparkDefaultDriverBlockManagerPortName
	}
	if r.Spec.Driver.DriverBlockManagerPort == 0 {
		log.Info("setting default driver block manager port", "value", sparkDefaultDriverBlockManagerPort)
		r.Spec.Driver.DriverBlockManagerPort = sparkDefaultDriverBlockManagerPort
	}
	if r.Spec.Driver.DriverUIPortName == "" {
		log.Info("setting default driver ui port name", "value", sparkDefaultDriverUIPortName)
		r.Spec.Driver.DriverUIPortName = sparkDefaultDriverUIPortName
	}
	if r.Spec.Driver.DriverUIPort == 0 {
		log.Info("setting default driver ui port", "value", sparkDefaultDriverUIPort)
		r.Spec.Driver.DriverUIPort = sparkDefaultDriverUIPort
	}
	if r.Spec.Image == nil {
		log.Info("setting default image", "value", *sparkDefaultImage)
		r.Spec.Image = sparkDefaultImage
	}

	nodes := []*SparkClusterNode{&r.Spec.Master.SparkClusterNode, &r.Spec.Worker.SparkClusterNode}
	for i := range nodes {
		node := nodes[i]
		if node.Annotations == nil {
			node.Annotations = make(map[string]string)
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

	if err := r.validateMutualTLSMode(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateWorkerReplicas(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateWorkerResourceRequestsCPU(); err != nil {
		allErrs = append(allErrs, err)
	}
	if errs := r.validatePorts(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateDriverConfigs(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateAutoscaler(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateImage(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateExtraConfigs(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateKeyTabConfigs(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if err := r.validateNetworkPolicies(); err != nil {
		allErrs = append(allErrs, err)
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

func (r *SparkCluster) validateNetworkPolicies() *field.Error {
	if r.Spec.NetworkPolicy.ExternalPolicyEnabled != nil &&
		*r.Spec.NetworkPolicy.ExternalPolicyEnabled &&
		len(r.Spec.NetworkPolicy.ExternalPodLabels) == 0 {
		return field.Invalid(
			field.NewPath("spec").Child("NetworkPolicy").Child("ExternalPodLabels"),
			r.Spec.NetworkPolicy,
			"should have at least one item if the policy is enabled",
		)
	}

	return nil
}

func (r *SparkCluster) validateExtraConfigs() field.ErrorList {
	var errs field.ErrorList

	if err := r.validateExtraConfig(r.Spec.Master.FrameworkConfig, "master"); err != nil {
		errs = append(errs, err...)
	}

	if err := r.validateExtraConfig(r.Spec.Worker.FrameworkConfig, "worker"); err != nil {
		errs = append(errs, err...)
	}

	return errs
}

func (r *SparkCluster) validateKeyTabConfigs() field.ErrorList {
	var errs field.ErrorList

	if err := r.validateKeyTabConfig(r.Spec.Master.KeyTabConfig, "master"); err != nil {
		errs = append(errs, err...)
	}

	if err := r.validateKeyTabConfig(r.Spec.Worker.KeyTabConfig, "worker"); err != nil {
		errs = append(errs, err...)
	}

	return errs
}

func (r *SparkCluster) validateExtraConfig(config *FrameworkConfig, comp string) field.ErrorList {
	var errs field.ErrorList
	if config == nil {
		return nil
	}
	if len(config.Configs) == 0 {
		errs = append(errs, field.Invalid(
			field.NewPath("spec").Child(comp).Child("frameworkConfig").Child("configs"),
			config.Configs,
			"should have at least one item",
		))
	}

	if config.Path == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec").Child(comp).Child("frameworkConfig").Child("path"),
			config.Path,
			"should be non-empty",
		))
	}
	return errs
}

func (r *SparkCluster) validateKeyTabConfig(config *KeyTabConfig, comp string) field.ErrorList {
	var errs field.ErrorList
	if config == nil {
		return nil
	}

	if config.Path == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec").Child(comp).Child("keyTabConfig").Child("path"),
			config.Path,
			"should be non-empty",
		))
	}

	if len(config.KeyTab) == 0 {
		errs = append(errs, field.Invalid(
			field.NewPath("spec").Child(comp).Child("keyTabConfig").Child("keytab"),
			config.KeyTab,
			"should have at least one item",
		))
	}

	return errs
}

func (r *SparkCluster) validateMutualTLSMode() *field.Error {
	if r.Spec.MutualTLSMode == "" {
		return nil
	}
	if _, ok := securityv1beta1.PeerAuthentication_MutualTLS_Mode_value[r.Spec.MutualTLSMode]; ok {
		return nil
	}

	var validModes []string
	for s := range securityv1beta1.PeerAuthentication_MutualTLS_Mode_value {
		validModes = append(validModes, s)
	}

	return field.Invalid(
		field.NewPath("spec").Child("istioMutualTLSMode"),
		r.Spec.MutualTLSMode,
		fmt.Sprintf("mode must be one of the following: %v", validModes),
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

func (r *SparkCluster) validateDriverConfigs() field.ErrorList {
	var errs field.ErrorList

	// validate driver ports
	if err := r.validatePort(r.Spec.Driver.DriverPort, field.NewPath("spec").Child("sparkClusterDriver").Child("driverPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.Driver.DriverBlockManagerPort,
		field.NewPath("spec").Child("sparkClusterDriver").Child("driverBlockManagerPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.Driver.DriverUIPort, field.NewPath("spec").Child("sparkClusterDriver").Child("driverUIPort")); err != nil {
		errs = append(errs, err)
	}

	// validate driver name configurations
	// if r.Spec.Driver.ExecutionName == "" {
	//	errs = append(errs, field.Invalid(
	//		field.NewPath("spec").Child("sparkClusterDriver").Child("executionName"),
	//		r.Spec.Driver.ExecutionName,
	//		"should be non-empty",
	//	))
	// }

	// if r.Spec.Driver.SparkClusterName == "" {
	//	errs = append(errs, field.Invalid(
	//		field.NewPath("spec").Child("sparkClusterDriver").Child("sparkClusterName"),
	//		r.Spec.Driver.SparkClusterName,
	//		"should be non-empty",
	//	))
	// }

	if r.Spec.Driver.DriverUIPortName == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec").Child("sparkClusterDriver").Child("driverUIPortName"),
			r.Spec.Driver.DriverUIPortName,
			"should be non-empty",
		))
	}

	if r.Spec.Driver.DriverPortName == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec").Child("sparkClusterDriver").Child("driverPortName"),
			r.Spec.Driver.DriverPortName,
			"should be non-empty",
		))
	}

	if r.Spec.Driver.DriverBlockManagerPortName == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec").Child("sparkClusterDriver").Child("driverBlockManagerPortName"),
			r.Spec.Driver.DriverBlockManagerPort,
			"should be non-empty",
		))
	}

	return errs
}

func (r *SparkCluster) validatePorts() field.ErrorList {
	var errs field.ErrorList

	if err := r.validatePort(r.Spec.ClusterPort, field.NewPath("spec").Child("clusterPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.TCPMasterWebPort, field.NewPath("spec").Child("tcpMasterWebPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.TCPWorkerWebPort, field.NewPath("spec").Child("tcpWorkerWebPort")); err != nil {
		errs = append(errs, err)
	}
	if err := r.validatePort(r.Spec.DashboardPort, field.NewPath("spec").Child("dashboardPort")); err != nil {
		errs = append(errs, err)
	}

	// TODO: add validation to prevent port values overlap

	return errs
}

func (r *SparkCluster) validatePort(port int32, fldPath *field.Path) *field.Error {
	if port < sparkMinValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be greater than or equal to %d", sparkMinValidPort))
	}
	if port > sparkMaxValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be less than or equal to %d", sparkMaxValidPort))
	}

	return nil
}

// nolint:dupl
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

func (r *SparkCluster) validateWorkerResourceRequestsCPU() *field.Error {
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

func (r *SparkCluster) validateImage() field.ErrorList {
	var errs field.ErrorList
	fldPath := field.NewPath("spec").Child("image")

	if r.Spec.Image.Repository == "" {
		errs = append(errs, field.Required(fldPath.Child("repository"), "cannot be blank"))
	}
	if r.Spec.Image.Tag == "" {
		errs = append(errs, field.Required(fldPath.Child("tag"), "cannot be blank"))
	}

	return errs
}

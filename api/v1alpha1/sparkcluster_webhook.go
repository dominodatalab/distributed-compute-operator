package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"

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
	sparkDefaultDashboardServicePort        int32 = 80
	sparkDefaultClusterPort                 int32 = 7077
	sparkDefaultDashboardPort               int32 = 8080
	sparkDefaultMasterWebPort               int32 = 8080
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
func (sc *SparkCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(sc).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-sparkcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=sparkclusters,verbs=create;update,versions=v1alpha1,name=msparkcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &SparkCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (sc *SparkCluster) Default() {
	log := sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(sc))
	log.Info("applying defaults")

	if sc.Spec.ClusterPort == 0 {
		log.Info("setting default cluster port", "value", sparkDefaultClusterPort)
		sc.Spec.ClusterPort = sparkDefaultClusterPort
	}
	if sc.Spec.TCPWorkerWebPort == 0 {
		log.Info("setting default worker web port", "value", sparkDefaultWorkerWebPort)
		sc.Spec.TCPWorkerWebPort = sparkDefaultWorkerWebPort
	}
	if sc.Spec.TCPMasterWebPort == 0 {
		log.Info("setting default master web port", "value", sparkDefaultMasterWebPort)
		sc.Spec.TCPMasterWebPort = sparkDefaultMasterWebPort
	}
	if sc.Spec.DashboardPort == 0 {
		log.Info("setting default dashboard port", "value", sparkDefaultDashboardPort)
		sc.Spec.DashboardPort = sparkDefaultDashboardPort
	}
	if sc.Spec.DashboardServicePort == 0 {
		log.Info("setting default dashboard service port", "value", sparkDefaultDashboardServicePort)
		sc.Spec.DashboardServicePort = sparkDefaultDashboardServicePort
	}
	if sc.Spec.EnableDashboard == nil {
		log.Info("setting enable dashboard flag", "value", *sparkDefaultEnableDashboard)
		sc.Spec.EnableDashboard = sparkDefaultEnableDashboard
	}
	if sc.Spec.NetworkPolicy.Enabled == nil {
		log.Info("setting enable network policy flag", "value", *sparkDefaultEnableNetworkPolicy)
		sc.Spec.NetworkPolicy.Enabled = sparkDefaultEnableNetworkPolicy
	}
	if sc.Spec.NetworkPolicy.ExternalPolicyEnabled == nil {
		log.Info("setting enable external network policy flag", "value", *sparkDefaultEnableExternalNetworkPolicy)
		sc.Spec.NetworkPolicy.ExternalPolicyEnabled = sparkDefaultEnableExternalNetworkPolicy
	}
	if sc.Spec.NetworkPolicy.ClientServerLabels == nil {
		log.Info("setting default network policy client labels", "value", sparkDefaultNetworkPolicyClientLabels)
		sc.Spec.NetworkPolicy.ClientServerLabels = sparkDefaultNetworkPolicyClientLabels
	}
	if sc.Spec.NetworkPolicy.DashboardLabels == nil {
		log.Info("setting default network policy dashboard labels", "value", sparkDefaultNetworkPolicyClientLabels)
		sc.Spec.NetworkPolicy.DashboardLabels = sparkDefaultNetworkPolicyClientLabels
	}
	if sc.Spec.Worker.Replicas == nil {
		log.Info("setting default worker replicas", "value", *sparkDefaultWorkerReplicas)
		sc.Spec.Worker.Replicas = sparkDefaultWorkerReplicas
	}
	if sc.Spec.Driver.DriverPort == 0 {
		log.Info("setting default driver port", "value", sparkDefaultDriverPort)
		sc.Spec.Driver.DriverPort = sparkDefaultDriverPort
	}
	if sc.Spec.Driver.DriverPortName == "" {
		log.Info("setting default driver port name", "value", sparkDefaultDriverPortName)
		sc.Spec.Driver.DriverPortName = sparkDefaultDriverPortName
	}
	if sc.Spec.Driver.DriverBlockManagerPortName == "" {
		log.Info("setting default driver block manager port name", "value", sparkDefaultDriverBlockManagerPortName)
		sc.Spec.Driver.DriverBlockManagerPortName = sparkDefaultDriverBlockManagerPortName
	}
	if sc.Spec.Driver.DriverBlockManagerPort == 0 {
		log.Info("setting default driver block manager port", "value", sparkDefaultDriverBlockManagerPort)
		sc.Spec.Driver.DriverBlockManagerPort = sparkDefaultDriverBlockManagerPort
	}
	if sc.Spec.Driver.DriverUIPortName == "" {
		log.Info("setting default driver ui port name", "value", sparkDefaultDriverUIPortName)
		sc.Spec.Driver.DriverUIPortName = sparkDefaultDriverUIPortName
	}
	if sc.Spec.Driver.DriverUIPort == 0 {
		log.Info("setting default driver ui port", "value", sparkDefaultDriverUIPort)
		sc.Spec.Driver.DriverUIPort = sparkDefaultDriverUIPort
	}
	if sc.Spec.Image == nil {
		log.Info("setting default image", "value", *sparkDefaultImage)
		sc.Spec.Image = sparkDefaultImage
	}

	nodes := []*SparkClusterNode{&sc.Spec.Master.SparkClusterNode, &sc.Spec.Worker.SparkClusterNode}
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
func (sc *SparkCluster) ValidateCreate() error {
	sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(sc)).Info("validating create")

	return sc.validateSparkCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (sc *SparkCluster) ValidateUpdate(old runtime.Object) error {
	sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(sc)).Info("validating update")

	return sc.validateSparkCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
// Not used, just here for interface compliance.
func (sc *SparkCluster) ValidateDelete() error {
	return nil
}

func (sc *SparkCluster) validateSparkCluster() error {
	var allErrs field.ErrorList

	if err := sc.validateMutualTLSMode(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := sc.validateWorkerReplicas(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := sc.validateWorkerMemoryLimit(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := sc.validateWorkerResourceRequestsCPU(); err != nil {
		allErrs = append(allErrs, err)
	}
	if errs := sc.validatePorts(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := sc.validateDriverConfigs(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := sc.validateAutoscaler(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := sc.validateImage(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := sc.validateFrameworkConfigs(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if errs := sc.validateKeyTabConfigs(); errs != nil {
		allErrs = append(allErrs, errs...)
	}
	if err := sc.validateNetworkPolicies(); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "distributed-compute.dominodatalab.com", Kind: "SparkCluster"},
		sc.Name,
		allErrs,
	)
}

func (sc *SparkCluster) validateNetworkPolicies() *field.Error {
	if sc.Spec.NetworkPolicy.ExternalPolicyEnabled != nil &&
		*sc.Spec.NetworkPolicy.ExternalPolicyEnabled &&
		len(sc.Spec.NetworkPolicy.ExternalPodLabels) == 0 {
		return field.Invalid(
			field.NewPath("spec", "networkPolicy", "externalPodLabels"),
			sc.Spec.NetworkPolicy,
			"should have at least one item if the policy is enabled",
		)
	}

	return nil
}

func (sc *SparkCluster) validateFrameworkConfigs() field.ErrorList {
	var errs field.ErrorList

	if err := sc.validateFrameworkConfig(sc.Spec.Master.FrameworkConfig, "master"); err != nil {
		errs = append(errs, err...)
	}

	if err := sc.validateFrameworkConfig(sc.Spec.Worker.FrameworkConfig, "worker"); err != nil {
		errs = append(errs, err...)
	}

	return errs
}

func (sc *SparkCluster) validateKeyTabConfigs() field.ErrorList {
	var errs field.ErrorList

	if err := sc.validateKeyTabConfig(sc.Spec.Master.KeyTabConfig, "master"); err != nil {
		errs = append(errs, err...)
	}

	if err := sc.validateKeyTabConfig(sc.Spec.Worker.KeyTabConfig, "worker"); err != nil {
		errs = append(errs, err...)
	}

	return errs
}

func (sc *SparkCluster) validateFrameworkConfig(config *FrameworkConfig, comp string) field.ErrorList {
	var errs field.ErrorList

	if config == nil {
		return nil
	}

	if len(config.Configs) == 0 {
		errs = append(errs, field.Invalid(
			field.NewPath("spec", comp, "frameworkConfig", "configs"),
			config.Configs,
			"should have at least one item",
		))
	}

	return errs
}

func (sc *SparkCluster) validateKeyTabConfig(config *KeyTabConfig, comp string) field.ErrorList {
	var errs field.ErrorList
	if config == nil {
		return nil
	}

	if config.Path == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec", comp, "keyTabConfig", "path"),
			config.Path,
			"should be non-empty",
		))
	}

	if len(config.KeyTab) == 0 {
		errs = append(errs, field.Invalid(
			field.NewPath("spec", comp, "keyTabConfig", "keytab"),
			config.KeyTab,
			"should have at least one item",
		))
	}

	return errs
}

func (sc *SparkCluster) validateMutualTLSMode() *field.Error {
	if sc.Spec.MutualTLSMode == "" {
		return nil
	}
	if _, ok := securityv1beta1.PeerAuthentication_MutualTLS_Mode_value[sc.Spec.MutualTLSMode]; ok {
		return nil
	}

	var validModes []string
	for s := range securityv1beta1.PeerAuthentication_MutualTLS_Mode_value {
		validModes = append(validModes, s)
	}

	return field.Invalid(
		field.NewPath("spec", "istioMutualTLSMode"),
		sc.Spec.MutualTLSMode,
		fmt.Sprintf("mode must be one of the following: %v", validModes),
	)
}

func (sc *SparkCluster) validateWorkerReplicas() *field.Error {
	replicas := sc.Spec.Worker.Replicas
	if replicas == nil || *replicas >= 0 {
		return nil
	}

	return field.Invalid(
		field.NewPath("spec", "worker", "replicas"),
		replicas,
		"should be greater than or equal to 0",
	)
}

func (sc *SparkCluster) validateWorkerMemoryLimit() *field.Error {
	request := sc.Spec.Worker.WorkerMemoryLimit

	if request == "" {
		return field.Invalid(
			field.NewPath("spec", "worker", "workerMemoryLimit"),
			request,
			"should be non-empty",
		)
	}

	if _, err := resource.ParseQuantity(request); err != nil {
		return field.Invalid(
			field.NewPath("spec", "worker", "workerMemoryLimit"),
			request,
			"should be a valid parsable quantity",
		)
	}

	return nil
}

func (sc *SparkCluster) validateDriverConfigs() field.ErrorList {
	var errs field.ErrorList

	// validate driver ports
	if err := sc.validatePort(sc.Spec.Driver.DriverPort, field.NewPath("spec", "sparkClusterDriver", "driverPort")); err != nil {
		errs = append(errs, err)
	}
	if err := sc.validatePort(sc.Spec.Driver.DriverBlockManagerPort,
		field.NewPath("spec", "sparkClusterDriver", "driverBlockManagerPort")); err != nil {
		errs = append(errs, err)
	}
	if err := sc.validatePort(sc.Spec.Driver.DriverUIPort,
		field.NewPath("spec", "sparkClusterDriver", "driverUIPort")); err != nil {
		errs = append(errs, err)
	}

	if sc.Spec.Driver.DriverUIPortName == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec", "sparkClusterDriver", "driverUIPortName"),
			sc.Spec.Driver.DriverUIPortName,
			"should be non-empty",
		))
	}

	if sc.Spec.Driver.DriverPortName == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec", "sparkClusterDriver", "driverPortName"),
			sc.Spec.Driver.DriverPortName,
			"should be non-empty",
		))
	}

	if sc.Spec.Driver.DriverBlockManagerPortName == "" {
		errs = append(errs, field.Invalid(
			field.NewPath("spec", "sparkClusterDriver", "driverBlockManagerPortName"),
			sc.Spec.Driver.DriverBlockManagerPort,
			"should be non-empty",
		))
	}

	return errs
}

func (sc *SparkCluster) validatePorts() field.ErrorList {
	var errs field.ErrorList

	if err := sc.validatePort(sc.Spec.ClusterPort, field.NewPath("spec", "clusterPort")); err != nil {
		errs = append(errs, err)
	}
	if err := sc.validatePort(sc.Spec.TCPMasterWebPort, field.NewPath("spec", "tcpMasterWebPort")); err != nil {
		errs = append(errs, err)
	}
	if err := sc.validatePort(sc.Spec.TCPWorkerWebPort, field.NewPath("spec", "tcpWorkerWebPort")); err != nil {
		errs = append(errs, err)
	}
	if err := sc.validatePort(sc.Spec.DashboardPort, field.NewPath("spec", "dashboardPort")); err != nil {
		errs = append(errs, err)
	}
	if err := sc.validatePort(sc.Spec.DashboardServicePort, field.NewPath("spec", "dashboardServicePort")); err != nil {
		errs = append(errs, err)
	}

	// TODO: add validation to prevent port values overlap

	return errs
}

func (sc *SparkCluster) validatePort(port int32, fldPath *field.Path) *field.Error {
	if port < sparkMinValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be greater than or equal to %d", sparkMinValidPort))
	}
	if port > sparkMaxValidPort {
		return field.Invalid(fldPath, port, fmt.Sprintf("must be less than or equal to %d", sparkMaxValidPort))
	}

	return nil
}

func (sc *SparkCluster) validateAutoscaler() field.ErrorList {
	var errs field.ErrorList

	as := sc.Spec.Autoscaling
	if as == nil {
		return nil
	}

	fldPath := field.NewPath("spec", "autoscaling")

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

func (sc *SparkCluster) validateWorkerResourceRequestsCPU() *field.Error {
	if sc.Spec.Autoscaling == nil {
		return nil
	}

	if _, ok := sc.Spec.Worker.Resources.Requests[v1.ResourceCPU]; ok {
		return nil
	}

	return field.Required(
		field.NewPath("spec", "worker", "resources", "requests", "cpu"),
		"is mandatory when autoscaling is enabled",
	)
}

func (sc *SparkCluster) validateImage() field.ErrorList {
	var errs field.ErrorList
	fldPath := field.NewPath("spec", "image")

	if sc.Spec.Image.Repository == "" {
		errs = append(errs, field.Required(fldPath.Child("repository"), "cannot be blank"))
	}
	if sc.Spec.Image.Tag == "" {
		errs = append(errs, field.Required(fldPath.Child("tag"), "cannot be blank"))
	}

	return errs
}

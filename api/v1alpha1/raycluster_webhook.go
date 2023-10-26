package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	rayDefaultRedisPort              int32 = 6379
	rayDefaultClientServerPort       int32 = 10001
	rayDefaultObjectManagerPort      int32 = 2384
	rayDefaultNodeManagerPort        int32 = 2385
	rayDefaultGCSServerPort          int32 = 2386
	rayDefaultDashboardPort          int32 = 8265
	rayDefaultRedisShardPorts              = []int32{6380, 6381}
	rayDefaultEnableDashboard              = ptr.To(true)
	rayDefaultEnableNetworkPolicy          = ptr.To(true)
	rayDefaultWorkerReplicas               = ptr.To(int32(1))
	rayDefaultNetworkPolicyPodLabels       = map[string]string{
		"ray-client": "true",
	}

	rayDefaultImage = &OCIImageDefinition{
		Repository: "rayproject/ray",
		Tag:        "1.6.0-cpu",
	}

	rayLogger = logf.Log.WithName("webhooks").WithName("RayCluster")
)

// SetupWebhookWithManager creates and registers this webhook with the manager.
// TODO: delete
func (rc *RayCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(rc).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-raycluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=create;update,versions=v1alpha1,name=mraycluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &RayCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (rc *RayCluster) Default() {
	log := rayLogger.WithValues("raycluster", client.ObjectKeyFromObject(rc))
	log.Info("Applying defaults")

	spec := &rc.Spec
	if spec.Port == 0 {
		log.Info("Setting default port", "value", rayDefaultRedisPort)
		rc.Spec.Port = rayDefaultRedisPort
	}
	if spec.RedisShardPorts == nil {
		log.Info("Setting default redis shard ports", "value", rayDefaultRedisShardPorts)
		rc.Spec.RedisShardPorts = rayDefaultRedisShardPorts
	}
	if spec.ClientServerPort == 0 {
		log.Info("Setting default client server port", "value", rayDefaultClientServerPort)
		rc.Spec.ClientServerPort = rayDefaultClientServerPort
	}
	if spec.ObjectManagerPort == 0 {
		log.Info("Setting default object manager port", "value", rayDefaultObjectManagerPort)
		rc.Spec.ObjectManagerPort = rayDefaultObjectManagerPort
	}
	if spec.GCSServerPort == 0 {
		log.Info("Setting default gcs server port", "value", rayDefaultGCSServerPort)
		rc.Spec.GCSServerPort = rayDefaultGCSServerPort
	}
	if spec.NodeManagerPort == 0 {
		log.Info("Setting default node manager port", "value", rayDefaultNodeManagerPort)
		rc.Spec.NodeManagerPort = rayDefaultNodeManagerPort
	}
	if spec.DashboardPort == 0 {
		log.Info("Setting default dashboard port", "value", rayDefaultDashboardPort)
		rc.Spec.DashboardPort = rayDefaultDashboardPort
	}
	if spec.EnableDashboard == nil {
		log.Info("Setting enable dashboard flag", "value", *rayDefaultEnableDashboard)
		rc.Spec.EnableDashboard = rayDefaultEnableDashboard
	}
	if spec.NetworkPolicy.Enabled == nil {
		log.Info("Setting enable network policy flag", "value", *rayDefaultEnableNetworkPolicy)
		rc.Spec.NetworkPolicy.Enabled = rayDefaultEnableNetworkPolicy
	}
	if spec.NetworkPolicy.ClientLabels == nil {
		log.Info("Setting default network policy client server labels", "value", rayDefaultNetworkPolicyPodLabels)
		rc.Spec.NetworkPolicy.ClientLabels = rayDefaultNetworkPolicyPodLabels
	}
	if spec.NetworkPolicy.DashboardLabels == nil {
		log.Info("Setting default network policy dashboard pod labels", "value", rayDefaultNetworkPolicyPodLabels)
		rc.Spec.NetworkPolicy.DashboardLabels = rayDefaultNetworkPolicyPodLabels
	}
	if spec.Worker.Replicas == nil {
		log.Info("Setting default worker replicas", "value", *rayDefaultWorkerReplicas)
		rc.Spec.Worker.Replicas = rayDefaultWorkerReplicas
	}
	if spec.Image == nil {
		log.Info("Setting default image", "value", *rayDefaultImage)
		rc.Spec.Image = rayDefaultImage
	}
}

//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-raycluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=rayclusters,verbs=create;update,versions=v1alpha1,name=vraycluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &RayCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (rc *RayCluster) ValidateCreate() (admission.Warnings, error) {
	rayLogger.WithValues("raycluster", client.ObjectKeyFromObject(rc)).Info("Validating create")
	return admission.Warnings{}, rc.validateRayCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (rc *RayCluster) ValidateUpdate(runtime.Object) (admission.Warnings, error) {
	rayLogger.WithValues("raycluster", client.ObjectKeyFromObject(rc)).Info("Validating update")
	return admission.Warnings{}, rc.validateRayCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (rc *RayCluster) ValidateDelete() (admission.Warnings, error) {
	// NOTE: not used, just here for interface compliance.
	return admission.Warnings{}, nil
}

func (rc *RayCluster) validateRayCluster() error {
	var errList field.ErrorList

	if err := validateIstioMutualTLSMode(rc.Spec.MutualTLSMode); err != nil {
		errList = append(errList, err)
	}
	if errs := validateImage(rc.Spec.Image); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateAutoscaler(rc.Spec.Autoscaling); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateKerberosKeytab(rc.Spec.KerberosKeytab); errs != nil {
		errList = append(errList, errs...)
	}
	if err := validateWorkerReplicas(rc.Spec.Worker.Replicas); err != nil {
		errList = append(errList, err)
	}
	if errs := validateWorkerResourceRequests(rc.Spec.Worker.Resources); errs != nil {
		errList = append(errList, errs...)
	}
	if err := validateObjectStoreMemoryBytes(rc.Spec.ObjectStoreMemoryBytes); err != nil {
		errList = append(errList, err)
	}

	ports := map[string]int32{
		"port":              rc.Spec.Port,
		"clientServerPort":  rc.Spec.ClientServerPort,
		"objectManagerPort": rc.Spec.ObjectManagerPort,
		"nodeManagerPort":   rc.Spec.NodeManagerPort,
		"gcsServerPort":     rc.Spec.GCSServerPort,
		"dashboardPort":     rc.Spec.DashboardPort,
	}
	for idx, port := range rc.Spec.RedisShardPorts {
		ports[fmt.Sprintf("redisShardPorts[%d]", idx)] = port
	}
	for idx, port := range rc.Spec.WorkerPorts {
		ports[fmt.Sprintf("workerPorts[%d]", idx)] = port
	}
	if errs := validatePorts(ports); errs != nil {
		errList = append(errList, errs...)
	}

	return invalidIfNotEmpty("RayCluster", rc.Name, errList)
}

func validateObjectStoreMemoryBytes(memBytes *int64) *field.Error {
	if memBytes == nil || *memBytes >= 78643200 {
		return nil
	}

	return field.Invalid(
		field.NewPath("spec", "objectStoreMemoryBytes"),
		memBytes,
		"should be greater than or equal to 78643200",
	)
}

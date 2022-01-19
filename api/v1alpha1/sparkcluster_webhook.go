package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	sparkDefaultClusterPort               int32 = 7077
	sparkDefaultMasterWebPort             int32 = 8080
	sparkDefaultWorkerWebPort             int32 = 8081
	sparkDefaultDriverUIPort              int32 = 4040
	sparkDefaultDriverPort                int32 = 4041
	sparkDefaultDriverBlockManagerPort    int32 = 4042
	sparkDefaultEnableNetworkPolicy             = pointer.BoolPtr(true)
	sparkDefaultWorkerReplicas                  = pointer.Int32Ptr(1)
	sparkDefaultNetworkPolicyClientLabels       = map[string]string{
		"spark-client": "true",
	}
	sparkDefaultImage = &OCIImageDefinition{
		Repository: "bitnami/spark",
		Tag:        "3.0.2-debian-10-r0",
	}

	sparkLogger = logf.Log.WithName("webhooks").WithName("SparkCluster")
)

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
	log.Info("Applying defaults")

	spec := &sc.Spec
	if spec.ClusterPort == 0 {
		log.Info("Setting default cluster port", "value", sparkDefaultClusterPort)
		spec.ClusterPort = sparkDefaultClusterPort
	}
	if spec.WorkerWebPort == 0 {
		log.Info("Setting default worker web port", "value", sparkDefaultWorkerWebPort)
		spec.WorkerWebPort = sparkDefaultWorkerWebPort
	}
	if spec.MasterWebPort == 0 {
		log.Info("Setting default master web port", "value", sparkDefaultMasterWebPort)
		spec.MasterWebPort = sparkDefaultMasterWebPort
	}
	if spec.NetworkPolicy.Enabled == nil {
		log.Info("Setting enable network policy flag", "value", *sparkDefaultEnableNetworkPolicy)
		spec.NetworkPolicy.Enabled = sparkDefaultEnableNetworkPolicy
	}
	if spec.NetworkPolicy.ClientLabels == nil {
		log.Info("Setting default network policy client labels", "value", sparkDefaultNetworkPolicyClientLabels)
		spec.NetworkPolicy.ClientLabels = sparkDefaultNetworkPolicyClientLabels
	}
	if spec.NetworkPolicy.DashboardLabels == nil {
		log.Info("Setting default network policy dashboard labels", "value", sparkDefaultNetworkPolicyClientLabels)
		spec.NetworkPolicy.DashboardLabels = sparkDefaultNetworkPolicyClientLabels
	}
	if spec.Worker.Replicas == nil {
		log.Info("Setting default worker replicas", "value", *sparkDefaultWorkerReplicas)
		spec.Worker.Replicas = sparkDefaultWorkerReplicas
	}
	if spec.Driver.Port == 0 {
		log.Info("Setting default driver port", "value", sparkDefaultDriverPort)
		spec.Driver.Port = sparkDefaultDriverPort
	}
	if spec.Driver.BlockManagerPort == 0 {
		log.Info("Setting default driver block manager port", "value", sparkDefaultDriverBlockManagerPort)
		spec.Driver.BlockManagerPort = sparkDefaultDriverBlockManagerPort
	}
	if spec.Driver.UIPort == 0 {
		log.Info("Setting default driver ui port", "value", sparkDefaultDriverUIPort)
		spec.Driver.UIPort = sparkDefaultDriverUIPort
	}
	if spec.Image == nil {
		log.Info("Setting default image", "value", *sparkDefaultImage)
		spec.Image = sparkDefaultImage
	}

	nodes := []*SparkClusterNode{&sc.Spec.Master, &sc.Spec.Worker.SparkClusterNode}
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
	sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(sc)).Info("Validating create")
	return sc.validateSparkCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (sc *SparkCluster) ValidateUpdate(runtime.Object) error {
	sparkLogger.WithValues("sparkcluster", client.ObjectKeyFromObject(sc)).Info("Validating update")
	return sc.validateSparkCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (sc *SparkCluster) ValidateDelete() error {
	// NOTE: not used, just here for interface compliance.
	return nil
}

func (sc *SparkCluster) validateSparkCluster() error {
	var errList field.ErrorList

	if err := validateIstioMutualTLSMode(sc.Spec.MutualTLSMode); err != nil {
		errList = append(errList, err)
	}
	if err := validateWorkerReplicas(sc.Spec.Worker.Replicas); err != nil {
		errList = append(errList, err)
	}
	if errs := validateImage(sc.Spec.Image); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateAutoscaler(sc.Spec.Autoscaling); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateWorkerResourceRequests(sc.Spec.Worker.Resources); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateKerberosKeytab(sc.Spec.KerberosKeytab); errs != nil {
		errList = append(errList, errs...)
	}

	if err := sc.validateWorkerMemoryLimit(); err != nil {
		errList = append(errList, err)
	}
	if errs := sc.validateDriverConfigs(); errs != nil {
		errList = append(errList, errs...)
	}

	ports := map[string]int32{
		"clusterPort":   sc.Spec.ClusterPort,
		"masterWebPort": sc.Spec.MasterWebPort,
		"workerWebPort": sc.Spec.WorkerWebPort,
	}
	if errs := validatePorts(ports); errs != nil {
		errList = append(errList, errs...)
	}

	return invalidIfNotEmpty("SparkCluster", sc.Name, errList)
}

func (sc *SparkCluster) validateWorkerMemoryLimit() *field.Error {
	request := sc.Spec.WorkerMemoryLimit
	fp := field.NewPath("spec", "workerMemoryLimit")

	if request == "" {
		return field.Required(fp, "cannot be blank")
	}
	if _, err := resource.ParseQuantity(request); err != nil {
		return field.Invalid(fp, request, "must be a parsable quantity")
	}

	return nil
}

func (sc *SparkCluster) validateDriverConfigs() field.ErrorList {
	var errs field.ErrorList

	driver := sc.Spec.Driver
	ports := map[string]int32{
		"driver.port":             driver.Port,
		"driver.uiPort":           driver.UIPort,
		"driver.blockManagerPort": driver.BlockManagerPort,
	}
	if pErrs := validatePorts(ports); pErrs != nil {
		errs = append(errs, pErrs...)
	}

	if driver.Selector == nil || len(driver.Selector) == 0 {
		errs = append(errs, field.Required(
			field.NewPath("spec", "driver", "selector"),
			"must provide a label selector for driver pods",
		))
	}

	return errs
}

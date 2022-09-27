package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const nobodyUID int64 = 65534

var (
	daskDefaultSchedulerPort int32 = 8786
	daskDefaultDashboardPort int32 = 8787
	daskDefaultWorkerPort    int32 = 3000
	daskDefaultNannyPort     int32 = 3001

	daskDefaultWorkerReplicas      = pointer.Int32(1)
	daskDefaultEnableNetworkPolicy = pointer.Bool(true)
	daskDefaultNetworkPolicyLabels = map[string]string{
		"dask-client": "true",
	}

	daskDefaultPodSecurityContext = &corev1.PodSecurityContext{
		RunAsUser: pointer.Int64Ptr(nobodyUID),
	}

	daskDefaultImage = &OCIImageDefinition{
		Repository: "daskdev/dask",
		Tag:        "2021.7.2",
		PullPolicy: corev1.PullIfNotPresent,
	}

	daskLogger = logf.Log.WithName("webhooks").WithName("DaskCluster")
)

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-daskcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=daskclusters,verbs=create;update,versions=v1alpha1,name=mdaskcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &DaskCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (dc *DaskCluster) Default() {
	log := daskLogger.WithValues("daskcluster", client.ObjectKeyFromObject(dc))
	log.Info("Applying defaults")

	spec := &dc.Spec
	if spec.SchedulerPort == 0 {
		log.Info("Setting default scheduler port", "value", daskDefaultSchedulerPort)
		spec.SchedulerPort = daskDefaultSchedulerPort
	}
	if spec.DashboardPort == 0 {
		log.Info("Setting default dashboard port", "value", daskDefaultDashboardPort)
		spec.DashboardPort = daskDefaultDashboardPort
	}
	if spec.WorkerPort == 0 {
		log.Info("Setting default worker port", "value", daskDefaultWorkerPort)
		spec.WorkerPort = daskDefaultWorkerPort
	}
	if spec.NannyPort == 0 {
		log.Info("Setting default nanny port", "value", daskDefaultNannyPort)
		spec.NannyPort = daskDefaultNannyPort
	}
	if spec.Image == nil {
		log.Info("Setting default image", "value", daskDefaultImage)
		spec.Image = daskDefaultImage
	}
	if spec.Worker.Replicas == nil {
		log.Info("Setting default worker replicas", "value", *daskDefaultWorkerReplicas)
		spec.Worker.Replicas = daskDefaultWorkerReplicas
	}
	if spec.PodSecurityContext == nil {
		log.Info("Setting default pod security context", "value", daskDefaultPodSecurityContext)
		spec.PodSecurityContext = daskDefaultPodSecurityContext
	}
	if spec.NetworkPolicy.Enabled == nil {
		log.Info("Setting enable network policy flag", "value", pointer.Bool(true))
		spec.NetworkPolicy.Enabled = daskDefaultEnableNetworkPolicy
	}
	if spec.NetworkPolicy.ClientLabels == nil {
		log.Info("Setting default network policy client labels", "values", daskDefaultNetworkPolicyLabels)
		spec.NetworkPolicy.ClientLabels = daskDefaultNetworkPolicyLabels
	}
	if spec.NetworkPolicy.DashboardPodLabels == nil {
		log.Info("Setting default network policy dashboard labels", "values", daskDefaultNetworkPolicyLabels)
		spec.NetworkPolicy.DashboardPodLabels = daskDefaultNetworkPolicyLabels
	}
}

//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-daskcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=daskclusters,verbs=create;update,versions=v1alpha1,name=vdaskcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &DaskCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (dc *DaskCluster) ValidateCreate() error {
	daskLogger.WithValues("daskcluster", client.ObjectKeyFromObject(dc)).Info("Validating create")
	return dc.validateDaskCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (dc *DaskCluster) ValidateUpdate(old runtime.Object) error {
	daskLogger.WithValues("daskcluster", client.ObjectKeyFromObject(dc)).Info("Validating update")
	return dc.validateDaskCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (dc *DaskCluster) ValidateDelete() error {
	// NOTE: not used, just here for interface compliance.
	return nil
}

func (dc *DaskCluster) validateDaskCluster() error {
	var errList field.ErrorList

	if err := validateIstioMutualTLSMode(dc.Spec.MutualTLSMode); err != nil {
		errList = append(errList, err)
	}
	if err := validateWorkerReplicas(dc.Spec.Worker.Replicas); err != nil {
		errList = append(errList, err)
	}
	if errs := validateImage(dc.Spec.Image); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateAutoscaler(dc.Spec.Autoscaling); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateWorkerResourceRequests(dc.Spec.Worker.Resources); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateKerberosKeytab(dc.Spec.KerberosKeytab); errs != nil {
		errList = append(errList, errs...)
	}

	ports := map[string]int32{
		"schedulerPort": dc.Spec.SchedulerPort,
		"workerPort":    dc.Spec.WorkerPort,
		"dashboardPort": dc.Spec.DashboardPort,
		"nannyPort":     dc.Spec.NannyPort,
	}
	if errs := validatePorts(ports); errs != nil {
		errList = append(errList, errs...)
	}

	return invalidIfNotEmpty("DaskCluster", dc.Name, errList)
}

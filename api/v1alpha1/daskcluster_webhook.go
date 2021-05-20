package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	daskDefaultImage               = &OCIImageDefinition{
		Repository: "daskdev/dask",
		Tag:        "2021.5.0",
		PullPolicy: corev1.PullIfNotPresent,
	}
	daskDefaultPodSecurityContext = &corev1.PodSecurityContext{
		RunAsUser: pointer.Int64Ptr(nobodyUID),
	}
)

// log is for logging in this package.
var daskclusterlog = logf.Log.WithName("daskcluster-resource")

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-daskcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=daskclusters,verbs=create;update,versions=v1alpha1,name=mdaskcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &DaskCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (dc *DaskCluster) Default() {
	// daskclusterlog.Info("default", "name", dc.Name)
	log := daskclusterlog.WithValues("daskcluster", client.ObjectKeyFromObject(dc))

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
	if spec.PodSecurityContext == nil {
		log.Info("Setting default pod security context", "value", daskDefaultPodSecurityContext)
		spec.PodSecurityContext = daskDefaultPodSecurityContext
	}
}

//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-daskcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=daskclusters,verbs=create;update,versions=v1alpha1,name=vdaskcluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &DaskCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (dc *DaskCluster) ValidateCreate() error {
	daskclusterlog.Info("validate create", "name", dc.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (dc *DaskCluster) ValidateUpdate(old runtime.Object) error {
	daskclusterlog.Info("validate update", "name", dc.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (dc *DaskCluster) ValidateDelete() error {
	daskclusterlog.Info("validate delete", "name", dc.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

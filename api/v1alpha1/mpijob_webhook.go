package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var mpijobLogger = logf.Log.WithName("webhooks").WithName("MPIJob")

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-mpijob,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=mpijobs,verbs=create;update,versions=v1alpha1,name=mmpijob.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &MPIJob{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (j *MPIJob) Default() {
	log := mpijobLogger.WithValues("mpijob", client.ObjectKeyFromObject(j))
	log.Info("Applying defaults")

	// TODO: add defaulting logic here.
}

//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-mpijob,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=mpijobs,verbs=create;update,versions=v1alpha1,name=vmpijob.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MPIJob{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (j *MPIJob) ValidateCreate() error {
	mpijobLogger.WithValues("mpijob", client.ObjectKeyFromObject(j)).Info("Validating create")
	return j.validateMPIJob()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (j *MPIJob) ValidateUpdate(old runtime.Object) error {
	mpijobLogger.WithValues("mpijob", client.ObjectKeyFromObject(j)).Info("Validating update")
	return j.validateMPIJob()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (j *MPIJob) ValidateDelete() error {
	// NOTE: not used, just here for interface compliance.
	return nil
}

func (j *MPIJob) validateMPIJob() error {
	return nil
}

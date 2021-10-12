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

var (
	mpiDefaultSlotsPerWorker = pointer.Int32(1)
	mpiDefaultWorkerReplicas = pointer.Int32(1)

	mpiDefaultImage = &OCIImageDefinition{
		Repository: "horovod/horovod",
		Tag:        "0.22.1",
		PullPolicy: corev1.PullIfNotPresent,
	}

	mpijobLogger = logf.Log.WithName("webhooks").WithName("MPIJob")
)

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-mpijob,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=mpijobs,verbs=create;update,versions=v1alpha1,name=mmpijob.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &MPIJob{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (j *MPIJob) Default() {
	log := mpijobLogger.WithValues("mpijob", client.ObjectKeyFromObject(j))
	log.Info("Applying defaults")

	spec := &j.Spec
	if spec.SlotsPerWorker == nil {
		log.Info("Setting default slots per worker", "value", mpiDefaultSlotsPerWorker)
		spec.SlotsPerWorker = mpiDefaultSlotsPerWorker
	}
	if spec.Worker.Replicas == nil {
		log.Info("Setting default worker replicas", "value", *mpiDefaultWorkerReplicas)
		spec.Worker.Replicas = mpiDefaultWorkerReplicas
	}
	if spec.Image == nil {
		log.Info("Setting default image", "value", mpiDefaultImage)
		spec.Image = mpiDefaultImage
	}
}

//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-mpijob,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=mpijobs,verbs=create;update,versions=v1alpha1,name=vmpijob.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MPIJob{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (j *MPIJob) ValidateCreate() error {
	mpijobLogger.WithValues("mpijob", client.ObjectKeyFromObject(j)).Info("Validating create")

	var errList field.ErrorList

	if err := validateIstioMutualTLSMode(j.Spec.MutualTLSMode); err != nil {
		errList = append(errList, err)
	}
	if err := validateWorkerReplicas(j.Spec.Worker.Replicas); err != nil {
		errList = append(errList, err)
	}
	if errs := validateImage(j.Spec.Image); errs != nil {
		errList = append(errList, errs...)
	}
	if errs := validateKerberosKeytab(j.Spec.KerberosKeytab); errs != nil {
		errList = append(errList, errs...)
	}

	slots := j.Spec.SlotsPerWorker
	if slots == nil || *slots < 1 {
		errList = append(errList, field.Invalid(
			field.NewPath("spec", "slotsPerWorker"),
			slots,
			"should be greate than or equal to 1",
		))
	}

	if j.Spec.Launcher.Command == nil {
		errList = append(errList, field.Required(
			field.NewPath("spec", "launcher", "command"),
			"must be provided",
		))
	}

	return invalidIfNotEmpty("MPIJob", j.Name, errList)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (j *MPIJob) ValidateUpdate(old runtime.Object) error {
	mpijobLogger.WithValues("mpijob", client.ObjectKeyFromObject(j)).Info("Validating update")

	// TODO: reject all updates to spec, or certain fields?
	// if equality.Semantic.DeepDerivative(j.Spec, old.(*MPIJob).Spec) {
	// 	return nil
	// }
	//
	// return apierrors.NewForbidden(schema.GroupResource{}, j.Name, errors.New(""))

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (j *MPIJob) ValidateDelete() error {
	// NOTE: not used, just here for interface compliance.
	return nil
}

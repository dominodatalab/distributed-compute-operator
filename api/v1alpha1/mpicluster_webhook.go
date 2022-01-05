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
	mpiDefaultSlotsPerWorker      = pointer.Int32(1)
	mpiDefaultWorkerReplicas      = pointer.Int32(1)
	mpiDefaultEnableNetworkPolicy = pointer.Bool(true)

	mpiDefaultImage = &OCIImageDefinition{
		Repository: "horovod/horovod",
		Tag:        "0.22.1",
		PullPolicy: corev1.PullIfNotPresent,
	}

	mpiClusterLogger = logf.Log.WithName("webhooks").WithName("MPICluster")
)

//+kubebuilder:webhook:path=/mutate-distributed-compute-dominodatalab-com-v1alpha1-mpicluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=mpiclusters,verbs=create;update,versions=v1alpha1,name=mmpicluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &MPICluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (j *MPICluster) Default() {
	log := mpiClusterLogger.WithValues("mpicluster", client.ObjectKeyFromObject(j))
	log.Info("Applying defaults")

	spec := &j.Spec
	if spec.Worker.Slots == nil {
		log.Info("Setting default worker slots", "value", mpiDefaultSlotsPerWorker)
		spec.Worker.Slots = mpiDefaultSlotsPerWorker
	}
	if spec.Worker.Replicas == nil {
		log.Info("Setting default worker replicas", "value", *mpiDefaultWorkerReplicas)
		spec.Worker.Replicas = mpiDefaultWorkerReplicas
	}
	if spec.Image == nil {
		log.Info("Setting default image", "value", mpiDefaultImage)
		spec.Image = mpiDefaultImage
	}
	if spec.NetworkPolicy.Enabled == nil {
		log.Info("Setting enable network policy flag", "value", *mpiDefaultEnableNetworkPolicy)
		spec.NetworkPolicy.Enabled = mpiDefaultEnableNetworkPolicy
	}
}

//+kubebuilder:webhook:path=/validate-distributed-compute-dominodatalab-com-v1alpha1-mpicluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=distributed-compute.dominodatalab.com,resources=mpiclusters,verbs=create;update,versions=v1alpha1,name=vmpicluster.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MPICluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (j *MPICluster) ValidateCreate() error {
	mpiClusterLogger.WithValues("mpicluster", client.ObjectKeyFromObject(j)).Info("Validating create")

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

	slots := j.Spec.Worker.Slots
	if slots == nil || *slots < 1 {
		errList = append(errList, field.Invalid(
			field.NewPath("spec", "worker", "slots"),
			slots,
			"should be greater than or equal to 1",
		))
	}

	return invalidIfNotEmpty("MPICluster", j.Name, errList)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (j *MPICluster) ValidateUpdate(old runtime.Object) error {
	mpiClusterLogger.WithValues("mpicluster", client.ObjectKeyFromObject(j)).Info("Validating update")

	// TODO: reject all updates to spec, or certain fields?
	// if equality.Semantic.DeepDerivative(j.Spec, old.(*MPICluster).Spec) {
	// 	return nil
	// }
	//
	// return apierrors.NewForbidden(schema.GroupResource{}, j.Name, errors.New(""))

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (j *MPICluster) ValidateDelete() error {
	// NOTE: not used, just here for interface compliance.
	return nil
}

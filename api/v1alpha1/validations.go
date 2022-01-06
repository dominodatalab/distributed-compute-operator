package v1alpha1

import (
	"fmt"

	securityv1beta1 "istio.io/api/security/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	minValidPort int32 = 1024
	maxValidPort int32 = 65535
)

func validateIstioMutualTLSMode(mode string) *field.Error {
	if mode == "" {
		return nil
	}
	if _, ok := securityv1beta1.PeerAuthentication_MutualTLS_Mode_value[mode]; ok {
		return nil
	}

	var validModes []string
	for v := range securityv1beta1.PeerAuthentication_MutualTLS_Mode_value {
		validModes = append(validModes, v)
	}

	return field.Invalid(
		field.NewPath("spec", "istioMutualTLSMode"),
		mode,
		fmt.Sprintf("mode must be one of the following: %v", validModes),
	)
}

func validateWorkerReplicas(replicas *int32) *field.Error {
	if replicas == nil || *replicas >= 0 {
		return nil
	}

	return field.Invalid(
		field.NewPath("spec", "worker", "replicas"),
		replicas,
		"should be greater than or equal to 0",
	)
}

func validateWorkerResourceRequests(resources corev1.ResourceRequirements) field.ErrorList {
	var errs field.ErrorList
	fp := field.NewPath("spec", "worker", "resources", "requests")

	if _, ok := resources.Requests[corev1.ResourceCPU]; !ok {
		errs = append(errs, field.Required(fp.Child("cpu"), "is mandatory"))
	}
	if _, ok := resources.Requests[corev1.ResourceMemory]; !ok {
		errs = append(errs, field.Required(fp.Child("memory"), "is mandatory"))
	}

	return errs
}

func validateImage(image *OCIImageDefinition) field.ErrorList {
	var errs field.ErrorList
	fp := field.NewPath("spec", "image")

	if image.Repository == "" {
		errs = append(errs, field.Required(fp.Child("repository"), "cannot be blank"))
	}
	if image.Tag == "" {
		errs = append(errs, field.Required(fp.Child("tag"), "cannot be blank"))
	}

	return errs
}

func validateAutoscaler(as *Autoscaling) field.ErrorList {
	if as == nil {
		return nil
	}

	var errs field.ErrorList
	fp := field.NewPath("spec", "autoscaling")

	if as.MinReplicas != nil {
		if *as.MinReplicas < 1 {
			errs = append(errs, field.Invalid(
				fp.Child("minReplicas"),
				as.MinReplicas,
				"must be greater than or equal to 1",
			))
		}

		if *as.MinReplicas > as.MaxReplicas {
			errs = append(errs, field.Invalid(
				fp.Child("maxReplicas"),
				as.MaxReplicas,
				"cannot be less than spec.autoscaling.minReplicas",
			))
		}
	}

	if as.MaxReplicas < 1 {
		errs = append(errs, field.Invalid(
			fp.Child("maxReplicas"),
			as.MaxReplicas,
			"must be greater than or equal to 1",
		))
	}

	if as.AverageCPUUtilization != nil && *as.AverageCPUUtilization <= 0 {
		errs = append(errs, field.Invalid(
			fp.Child("averageCPUUtilization"),
			as.AverageCPUUtilization,
			"must be greater than 0",
		))
	}

	if as.AverageMemoryUtilization != nil && *as.AverageMemoryUtilization <= 0 {
		errs = append(errs, field.Invalid(
			fp.Child("averageMemoryUtilization"),
			as.AverageMemoryUtilization,
			"must be greater than 0",
		))
	}

	if as.ScaleDownStabilizationWindowSeconds != nil && *as.ScaleDownStabilizationWindowSeconds < 0 {
		errs = append(errs, field.Invalid(
			fp.Child("scaleDownStabilizationWindowSeconds"),
			as.ScaleDownStabilizationWindowSeconds,
			"must be greater than or equal to 0",
		))
	}

	return errs
}

func validatePorts(portMap map[string]int32) field.ErrorList {
	var errs field.ErrorList
	fp := field.NewPath("spec")

	ports := sets.Int32{}
	for name, port := range portMap {
		if ports.Has(port) {
			errs = append(errs, field.Duplicate(fp.Child(name), port))
		} else {
			ports.Insert(port)
		}

		if port < minValidPort {
			errs = append(errs, field.Invalid(fp.Child(name), port, fmt.Sprintf("must be greater than or equal to %d", minValidPort)))
		} else if port > maxValidPort {
			errs = append(errs, field.Invalid(fp.Child(name), port, fmt.Sprintf("must be less than or equal to %d", maxValidPort)))
		}
	}

	return errs
}

func validateKerberosKeytab(keytab *KerberosKeytabConfig) field.ErrorList {
	if keytab == nil {
		return nil
	}

	var errs field.ErrorList
	fp := field.NewPath("spec", "kerberosKeytab")

	if len(keytab.Contents) == 0 {
		errs = append(errs, field.Required(fp.Child("contents"), "must contain file contents"))
	}
	if keytab.MountPath == "" {
		errs = append(errs, field.Required(fp.Child("mountPath"), "must be a valid file path"))
	}

	return errs
}

func validateSharedSSHSecret(secret string) field.ErrorList {
	var errs field.ErrorList
	fp := field.NewPath("spec", "workers", "sharedSSHSecret")
	if secret == "" {
		errs = append(errs, field.Required(fp, "must contain a shared secret name"))
	}

	return errs
}

func invalidIfNotEmpty(kind, name string, errList field.ErrorList) error {
	if len(errList) == 0 {
		return nil
	}

	return apierrors.NewInvalid(GroupVersion.WithKind(kind).GroupKind(), name, errList)
}

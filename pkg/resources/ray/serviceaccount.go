package ray

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// NewServiceAccount generates a service account resource without API access.
func NewServiceAccount(rc *dcv1alpha1.RayCluster) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, ComponentNone),
			Namespace: rc.Namespace,
			Labels:    AddGlobalLabels(MetadataLabels(rc), rc.Spec.GlobalLabels),
		},
		AutomountServiceAccountToken: pointer.Bool(false),
	}
}

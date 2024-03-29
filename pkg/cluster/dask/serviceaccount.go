package dask

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func ServiceAccount() core.OwnedComponent {
	factory := func(obj client.Object) components.ServiceAccountDataSource {
		return &serviceAccountDS{dc: daskCluster(obj)}
	}

	return components.ServiceAccount(factory)
}

type serviceAccountDS struct {
	dc *dcv1alpha1.DaskCluster
}

func (s *serviceAccountDS) ServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(s.dc, metadata.ComponentNone),
			Namespace: s.dc.Namespace,
			Labels:    meta.StandardLabels(s.dc),
		},
		AutomountServiceAccountToken: pointer.Bool(s.dc.Spec.ServiceAccount.AutomountServiceAccountToken),
	}
}

func (s *serviceAccountDS) Delete() bool {
	return s.dc.Spec.ServiceAccount.Name != ""
}

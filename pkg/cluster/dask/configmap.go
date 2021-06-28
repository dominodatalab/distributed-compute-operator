package dask

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func ConfigMapKeyTab() core.OwnedComponent {
	return components.ConfigMap(func(obj client.Object) components.ConfigMapDataSource {
		return &configMapDS{dc: daskCluster(obj)}
	})
}

type configMapDS struct {
	dc *dcv1alpha1.DaskCluster
}

func (s *configMapDS) ConfigMap() *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(s.dc, metadata.ComponentNone),
			Namespace: s.dc.Namespace,
			Labels:    meta.StandardLabels(s.dc),
		},
	}

	if s.dc.Spec.KerberosKeytab == nil {
		return cm
	}
	cm.BinaryData = map[string][]byte{"keytab": s.dc.Spec.KerberosKeytab.Contents}

	return cm
}

func (s *configMapDS) Delete() bool {
	return s.dc.Spec.KerberosKeytab == nil
}

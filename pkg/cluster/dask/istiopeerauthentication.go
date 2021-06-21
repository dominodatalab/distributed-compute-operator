package dask

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/istio"
)

func IstioPeerAuthentication(enabled bool) core.Component {
	return components.IstioPeerAuthentication(func(obj client.Object) components.IstioPeerAuthenticationDataSource {
		return &istioPeerAuthenticationDS{dc: daskCluster(obj), enabled: enabled}
	})
}

type istioPeerAuthenticationDS struct {
	dc      *dcv1alpha1.DaskCluster
	enabled bool
}

func (s *istioPeerAuthenticationDS) PeerAuthInfo() *istio.PeerAuthInfo {
	return &istio.PeerAuthInfo{
		Name:      meta.InstanceName(s.dc, metadata.ComponentNone),
		Namespace: s.dc.Namespace,
		Labels:    meta.StandardLabels(s.dc),
		Selector:  meta.MatchLabels(s.dc),
		Mode:      s.dc.Spec.MutualTLSMode,
	}
}

func (s *istioPeerAuthenticationDS) Enabled() bool {
	return s.enabled
}

func (s *istioPeerAuthenticationDS) Delete() bool {
	return s.dc.Spec.MutualTLSMode == ""
}

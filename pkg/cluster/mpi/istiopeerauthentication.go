package mpi

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
		return &istioPeerAuthenticationDS{mpi: objToMPICluster(obj), enabled: enabled}
	})
}

type istioPeerAuthenticationDS struct {
	mpi     *dcv1alpha1.MPICluster
	enabled bool
}

func (s *istioPeerAuthenticationDS) PeerAuthInfo() *istio.PeerAuthInfo {
	return &istio.PeerAuthInfo{
		Name:      meta.InstanceName(s.mpi, metadata.ComponentNone),
		Namespace: s.mpi.Namespace,
		Labels:    meta.StandardLabels(s.mpi),
		Selector:  meta.MatchLabels(s.mpi),
		Mode:      s.mpi.Spec.MutualTLSMode,
	}
}

func (s *istioPeerAuthenticationDS) Enabled() bool {
	return s.enabled
}

func (s *istioPeerAuthenticationDS) Delete() bool {
	return s.mpi.Spec.MutualTLSMode == ""
}

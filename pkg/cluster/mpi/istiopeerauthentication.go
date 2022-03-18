package mpi

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/istio"
	authenticationv1alpha1 "istio.io/api/authentication/v1alpha1"
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

func IstioClientPeerAuthentication(enabled bool) core.Component {
	return components.IstioPeerAuthentication(func(obj client.Object) components.IstioPeerAuthenticationDataSource {
		return &istioClientPeerAuthenticationDS{mpi: objToMPICluster(obj), enabled: enabled}
	})
}

type istioClientPeerAuthenticationDS struct {
	mpi     *dcv1alpha1.MPICluster
	enabled bool
}

func (s *istioClientPeerAuthenticationDS) PeerAuthInfo() *istio.PeerAuthInfo {
	return &istio.PeerAuthInfo{
		Name:      meta.InstanceName(s.mpi, ComponentClient),
		Namespace: s.mpi.Namespace,
		Labels:    meta.StandardLabels(s.mpi),
		Selector:  s.mpi.Spec.NetworkPolicy.ClientLabels,
		Mode:      authenticationv1alpha1.MutualTls_PERMISSIVE.String(),
	}
}

func (s *istioClientPeerAuthenticationDS) Enabled() bool {
	return s.enabled && s.mpi.Spec.Worker.Annotations["sidecar.istio.io/inject"] == "false"
}

func (s *istioClientPeerAuthenticationDS) Delete() bool {
	return false
}

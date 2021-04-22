package istio

import (
	securityv1beta1 "istio.io/api/security/v1beta1"
	"istio.io/api/type/v1beta1"
	istio "istio.io/client-go/pkg/apis/security/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PeerAuthInfo defines fields used to generate Istio PeerAuthentication objects.
type PeerAuthInfo struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Selector  map[string]string
	Mode      string
}

// NewPeerAuthentication uses PeerAuthInfo to generate and return a new PeerAuthentication object.
func NewPeerAuthentication(info *PeerAuthInfo) *istio.PeerAuthentication {
	modeVal := securityv1beta1.PeerAuthentication_MutualTLS_Mode_value[info.Mode]

	return &istio.PeerAuthentication{
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.Name,
			Namespace: info.Namespace,
			Labels:    info.Labels,
		},
		Spec: securityv1beta1.PeerAuthentication{
			Selector: &v1beta1.WorkloadSelector{
				MatchLabels: info.Selector,
			},
			Mtls: &securityv1beta1.PeerAuthentication_MutualTLS{
				Mode: securityv1beta1.PeerAuthentication_MutualTLS_Mode(modeVal),
			},
		},
	}
}

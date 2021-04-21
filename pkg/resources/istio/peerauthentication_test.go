package istio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	securityv1beta1 "istio.io/api/security/v1beta1"
	"istio.io/api/type/v1beta1"
	istio "istio.io/client-go/pkg/apis/security/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewPeerAuthentication(t *testing.T) {
	testcases := []struct {
		smode string
		mode  securityv1beta1.PeerAuthentication_MutualTLS_Mode
	}{
		{"UNSET", securityv1beta1.PeerAuthentication_MutualTLS_UNSET},
		{"DISABLE", securityv1beta1.PeerAuthentication_MutualTLS_DISABLE},
		{"PERMISSIVE", securityv1beta1.PeerAuthentication_MutualTLS_PERMISSIVE},
		{"STRICT", securityv1beta1.PeerAuthentication_MutualTLS_STRICT},
		{"GARBAGE", securityv1beta1.PeerAuthentication_MutualTLS_UNSET},
	}
	for _, tc := range testcases {
		info := &PeerAuthInfo{
			Name:      "cluster",
			Namespace: "ns",
			Labels: map[string]string{
				"awesome": "true",
			},
			Selector: map[string]string{
				"app.kubernetes.io/name": "compute-r",
			},
			Mode: tc.smode,
		}
		actual := NewPeerAuthentication(info)

		expected := &istio.PeerAuthentication{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cluster",
				Namespace: "ns",
				Labels: map[string]string{
					"awesome": "true",
				},
			},
			Spec: securityv1beta1.PeerAuthentication{
				Selector: &v1beta1.WorkloadSelector{
					MatchLabels: map[string]string{
						"app.kubernetes.io/name": "compute-r",
					},
				},
				Mtls: &securityv1beta1.PeerAuthentication_MutualTLS{
					Mode: tc.mode,
				},
			},
		}

		assert.Equal(t, expected, actual)
	}
}

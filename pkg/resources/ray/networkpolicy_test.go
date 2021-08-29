package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestNewClusterNetworkPolicy(t *testing.T) {
	rc := rayClusterFixture()
	netpol := NewClusterNetworkPolicy(rc)

	expected := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-ray-cluster",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "ray",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
			Annotations: map[string]string{
				"distributed-compute.dominodatalab.com/description": "Allows all ingress traffic between cluster nodes",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":     "ray",
					"app.kubernetes.io/instance": "test-id",
				},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/name":     "ray",
									"app.kubernetes.io/instance": "test-id",
								},
							},
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				"Ingress",
			},
		},
	}
	assert.Equal(t, expected, netpol)
}

func TestNewHeadClientNetworkPolicy(t *testing.T) {
	rc := rayClusterFixture()
	rc.Spec.NetworkPolicy = dcv1alpha1.NetworkPolicyConfig{
		ClientLabels: map[string]string{
			"server-client": "true",
		},
	}
	netpol := NewHeadClientNetworkPolicy(rc)

	tcpProto := v1.ProtocolTCP
	clientPort := intstr.FromInt(10001)
	expected := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-ray-client",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "ray",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/component":  "head",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
			Annotations: map[string]string{
				"distributed-compute.dominodatalab.com/description": "Allows client ingress traffic to head client server port",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":      "ray",
					"app.kubernetes.io/instance":  "test-id",
					"app.kubernetes.io/component": "head",
				},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &tcpProto,
							Port:     &clientPort,
						},
					},
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"server-client": "true",
								},
							},
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				"Ingress",
			},
		},
	}
	assert.Equal(t, expected, netpol)
}

func TestNewHeadDashboardNetworkPolicy(t *testing.T) {
	rc := rayClusterFixture()
	rc.Spec.NetworkPolicy = dcv1alpha1.NetworkPolicyConfig{
		DashboardLabels: map[string]string{
			"dashboard-client": "true",
		},
	}
	netpol := NewHeadDashboardNetworkPolicy(rc)

	tcpProto := v1.ProtocolTCP
	dashboardPort := intstr.FromInt(8265)
	expected := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-ray-dashboard",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "ray",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/component":  "head",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
			Annotations: map[string]string{
				"distributed-compute.dominodatalab.com/description": "Allows client ingress traffic to head dashboard port",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":      "ray",
					"app.kubernetes.io/instance":  "test-id",
					"app.kubernetes.io/component": "head",
				},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &tcpProto,
							Port:     &dashboardPort,
						},
					},
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"dashboard-client": "true",
								},
							},
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				"Ingress",
			},
		},
	}
	assert.Equal(t, expected, netpol)
}

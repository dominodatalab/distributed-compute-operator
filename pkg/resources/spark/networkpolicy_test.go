package spark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestNewClusterNetworkPolicy(t *testing.T) {
	rc := sparkClusterFixture()
	netpol := NewClusterNetworkPolicy(rc)

	expected := &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkPolicy",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-spark-cluster",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
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
					"app.kubernetes.io/name":     "spark",
					"app.kubernetes.io/instance": "test-id",
				},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/name":     "spark",
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

func TestNewHeadNetworkPolicy(t *testing.T) {
	rc := sparkClusterFixture()
	rc.Spec.NetworkPolicyClientLabels = []map[string]string{
		{
			"spark-client": "true",
		},
	}
	netpol := NewHeadNetworkPolicy(rc)

	tcpProto := v1.ProtocolTCP
	clientPort := intstr.FromInt(10001)
	dashboardPort := intstr.FromInt(8265)
	expected := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-spark-client",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/component":  "head",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
			Annotations: map[string]string{
				"distributed-compute.dominodatalab.com/description": "Allows client ingress traffic to head node",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":      "spark",
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
						{
							Protocol: &tcpProto,
							Port:     &dashboardPort,
						},
					},
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"spark-client": "true",
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

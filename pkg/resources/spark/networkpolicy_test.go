package spark

import (
	"testing"

	"github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"

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

func TestNewHeadClientNetworkPolicy(t *testing.T) {
	rc := sparkClusterFixture()
	labels := map[string]string{
		"spark-client": "true",
	}
	rc.Spec.NetworkPolicy = v1alpha1.SparkClusterNetworkPolicy{
		ClientServerLabels: labels,
		DashboardLabels:    labels,
	}
	netpol := NewHeadClientNetworkPolicy(rc)

	tcpProto := v1.ProtocolTCP
	clusterPort := intstr.FromInt(7077)
	expected := getNetworkPolicy(tcpProto, clusterPort, "test-id-spark-client", "Allows client ingress traffic to head client server port")
	assert.Equal(t, expected, netpol)
}

func TestNewHeadDashboardNetworkPolicy(t *testing.T) {
	rc := sparkClusterFixture()
	labels := map[string]string{
		"spark-client": "true",
	}
	rc.Spec.NetworkPolicy = v1alpha1.SparkClusterNetworkPolicy{
		ClientServerLabels: labels,
		DashboardLabels:    labels,
	}
	netpol := NewHeadDashboardNetworkPolicy(rc)

	tcpProto := v1.ProtocolTCP
	clusterPort := intstr.FromInt(8265)
	expected := getNetworkPolicy(tcpProto, clusterPort, "test-id-spark-dashboard", "Allows client ingress traffic to head dashboard port")
	assert.Equal(t, expected, netpol)
}

func getNetworkPolicy(tcpProto v1.Protocol, clusterPort intstr.IntOrString, name string, description string) *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkPolicy",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/component":  "master",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
			Annotations: map[string]string{
				"distributed-compute.dominodatalab.com/description": description,
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":      "spark",
					"app.kubernetes.io/instance":  "test-id",
					"app.kubernetes.io/component": "master",
				},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &tcpProto,
							Port:     &clusterPort,
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
}

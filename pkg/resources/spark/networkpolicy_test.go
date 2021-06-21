package spark

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestNewClusterDriverNetworkPolicy(t *testing.T) {

}

func TestNewClusterMasterNetworkPolicy(t *testing.T) {

}

func TestNewClusterWorkerNetworkPolicy(t *testing.T) {

}

// func TestNewClusterNetworkPolicy(t *testing.T) {
//	rc := sparkClusterFixture()
//	netpol := NewClusterNetworkPolicy(rc)
//
//	expected := &networkingv1.NetworkPolicy{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "test-id-spark-cluster",
//			Namespace: "fake-ns",
//			Labels: map[string]string{
//				"app.kubernetes.io/name":       "spark",
//				"app.kubernetes.io/instance":   "test-id",
//				"app.kubernetes.io/version":    "fake-tag",
//				"app.kubernetes.io/managed-by": "distributed-compute-operator",
//			},
//			Annotations: map[string]string{
//				"distributed-compute.dominodatalab.com/description": "Allows all ingress traffic between cluster nodes",
//			},
//		},
//		Spec: networkingv1.NetworkPolicySpec{
//			PodSelector: metav1.LabelSelector{
//				MatchLabels: map[string]string{
//					"app.kubernetes.io/name":     "spark",
//					"app.kubernetes.io/instance": "test-id",
//				},
//			},
//			Ingress: []networkingv1.NetworkPolicyIngressRule{
//				{
//					From: []networkingv1.NetworkPolicyPeer{
//						{
//							PodSelector: &metav1.LabelSelector{
//								MatchLabels: map[string]string{
//									"app.kubernetes.io/name":     "spark",
//									"app.kubernetes.io/instance": "test-id",
//								},
//							},
//						},
//					},
//				},
//			},
//			PolicyTypes: []networkingv1.PolicyType{
//				"Ingress",
//			},
//		},
//	}
//	assert.Equal(t, expected, netpol)
// }
//
// func TestNewClusterNetworkPolicyWithDriver(t *testing.T) {
//	rc := sparkClusterFixture()
//	rc.Spec.NetworkPolicy.ExternalPolicyEnabled = pointer.BoolPtr(true)
//	rc.Spec.NetworkPolicy.ExternalPodLabels = map[string]string{"app.kubernetes.io/instance": "spark-driver"}
//
//	netpol := NewClusterNetworkPolicy(rc)
//
//	expected := &networkingv1.NetworkPolicy{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "test-id-spark-cluster",
//			Namespace: "fake-ns",
//			Labels: map[string]string{
//				"app.kubernetes.io/name":       "spark",
//				"app.kubernetes.io/instance":   "test-id",
//				"app.kubernetes.io/version":    "fake-tag",
//				"app.kubernetes.io/managed-by": "distributed-compute-operator",
//			},
//			Annotations: map[string]string{
//				"distributed-compute.dominodatalab.com/description": "Allows all ingress traffic between cluster nodes",
//			},
//		},
//		Spec: networkingv1.NetworkPolicySpec{
//			PodSelector: metav1.LabelSelector{
//				MatchLabels: map[string]string{
//					"app.kubernetes.io/name":     "spark",
//					"app.kubernetes.io/instance": "test-id",
//				},
//			},
//			Ingress: []networkingv1.NetworkPolicyIngressRule{
//				{
//					From: []networkingv1.NetworkPolicyPeer{
//						{
//							PodSelector: &metav1.LabelSelector{
//								MatchLabels: map[string]string{
//									"app.kubernetes.io/name":     "spark",
//									"app.kubernetes.io/instance": "test-id",
//								},
//							},
//						},
//						{
//							PodSelector: &metav1.LabelSelector{
//								MatchLabels: map[string]string{
//									"app.kubernetes.io/instance": "spark-driver",
//								},
//							},
//						},
//					},
//				},
//			},
//			PolicyTypes: []networkingv1.PolicyType{
//				"Ingress",
//			},
//		},
//	}
//	assert.Equal(t, expected, netpol)
// }
//
// func TestNewClusterExternalNetworkPolicy(t *testing.T) {
//	rc := sparkClusterFixture()
//	labels := map[string]string{
//		"instance": "spark-external",
//	}
//	rc.Spec.NetworkPolicy = v1alpha1.SparkClusterNetworkPolicy{
//		ExternalPodLabels:     labels,
//		ExternalPolicyEnabled: pointer.BoolPtr(true),
//	}
//
//	netpol := NewClusterExternalNetworkPolicy(rc)
//
//	expected := &networkingv1.NetworkPolicy{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "test-id-spark-external",
//			Namespace: "fake-ns",
//			Labels: map[string]string{
//				"app.kubernetes.io/name":       "spark",
//				"app.kubernetes.io/instance":   "test-id",
//				"app.kubernetes.io/version":    "fake-tag",
//				"app.kubernetes.io/managed-by": "distributed-compute-operator",
//			},
//			Annotations: map[string]string{
//				"distributed-compute.dominodatalab.com/description": "Allows all ingress traffic between cluster and external nodes",
//			},
//		},
//		Spec: networkingv1.NetworkPolicySpec{
//			PodSelector: metav1.LabelSelector{
//				MatchLabels: map[string]string{
//					"instance": "spark-external",
//				},
//			},
//			Ingress: []networkingv1.NetworkPolicyIngressRule{
//				{
//					From: []networkingv1.NetworkPolicyPeer{
//						{
//							PodSelector: &metav1.LabelSelector{
//								MatchLabels: map[string]string{
//									"app.kubernetes.io/name":     "spark",
//									"app.kubernetes.io/instance": "test-id",
//								},
//							},
//						},
//					},
//				},
//			},
//			PolicyTypes: []networkingv1.PolicyType{
//				"Ingress",
//			},
//		},
//	}
//
//	assert.Equal(t, expected, netpol)
// }
//
// func TestNewHeadClientNetworkPolicy(t *testing.T) {
//	rc := sparkClusterFixture()
//	labels := map[string]string{
//		"spark-client": "true",
//	}
//	rc.Spec.NetworkPolicy = v1alpha1.SparkClusterNetworkPolicy{
//		ClientServerLabels: labels,
//		DashboardLabels:    labels,
//	}
//	netpol := NewHeadClientNetworkPolicy(rc)
//
//	tcpProto := v1.ProtocolTCP
//	clusterPort := intstr.FromInt(7077)
//	expected := getNetworkPolicy(tcpProto, clusterPort, "test-id-spark-client", "Allows client ingress traffic to head client server port")
//	assert.Equal(t, expected, netpol)
// }
//
// func TestNewHeadDashboardNetworkPolicy(t *testing.T) {
//	rc := sparkClusterFixture()
//	labels := map[string]string{
//		"spark-client": "true",
//	}
//	rc.Spec.NetworkPolicy = v1alpha1.SparkClusterNetworkPolicy{
//		ClientServerLabels: labels,
//		DashboardLabels:    labels,
//	}
//	netpol := NewHeadDashboardNetworkPolicy(rc)
//
//	tcpProto := v1.ProtocolTCP
//	clusterPort := intstr.FromInt(8265)
//	expected := getNetworkPolicy(tcpProto, clusterPort, "test-id-spark-dashboard", "Allows client ingress traffic to head dashboard port")
//	assert.Equal(t, expected, netpol)
// }

func getNetworkPolicy(tcpProto v1.Protocol, clusterPort intstr.IntOrString, name string, description string) *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
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

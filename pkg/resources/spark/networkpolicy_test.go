package spark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestNewClusterDriverNetworkPolicy(t *testing.T) {
	rc := sparkClusterFixture()
	rc.Spec.Driver.UIPort = 4040
	rc.Spec.NetworkPolicy.ClientLabels = map[string]string{"app.kubernetes.io/instance": "spark-driver"}

	netpol := NewClusterDriverNetworkPolicy(rc)

	protocol := corev1.ProtocolTCP
	driverUIPort := intstr.FromInt(4040)

	expected := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-spark-driver",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
			Annotations: map[string]string{
				"distributed-compute.dominodatalab.com/description": "driver network policy",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/instance": "spark-driver",
				},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/component": "master",
									"app.kubernetes.io/name":      "spark",
									"app.kubernetes.io/instance":  "test-id",
								},
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &protocol,
							Port:     &driverUIPort,
						},
					},
				},
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/component": "worker",
									"app.kubernetes.io/name":      "spark",
									"app.kubernetes.io/instance":  "test-id",
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

func TestNewClusterWorkerNetworkPolicy(t *testing.T) {
	rc := sparkClusterFixture()
	rc.Spec.Driver.UIPort = 4040
	rc.Spec.NetworkPolicy.ClientLabels = map[string]string{"app.kubernetes.io/instance": "spark-driver"}

	netpol := NewClusterWorkerNetworkPolicy(rc)

	expected := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-spark-worker",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
			Annotations: map[string]string{
				"distributed-compute.dominodatalab.com/description": "worker network policy",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/component":  "worker",
					"app.kubernetes.io/name":       "spark",
					"app.kubernetes.io/instance":   "test-id",
					"app.kubernetes.io/version":    "fake-tag",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
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
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/instance": "spark-driver",
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

func TestNewClusterMasterNetworkPolicy(t *testing.T) {
	rc := sparkClusterFixture()
	rc.Spec.NetworkPolicy.ClientLabels = map[string]string{"app.kubernetes.io/instance": "spark-driver"}
	rc.Spec.NetworkPolicy.DashboardPodLabels = map[string]string{"spark-client": "true"}
	rc.Spec.NetworkPolicy.DashboardNamespaceLabels = map[string]string{"domino-platform": "true"}

	netpol := NewClusterMasterNetworkPolicy(rc)

	protocol := corev1.ProtocolTCP
	masterDashboardPort := intstr.FromInt(8080)
	clusterPort := intstr.FromInt(7077)

	expected := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-spark-master",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
			Annotations: map[string]string{
				"distributed-compute.dominodatalab.com/description": "master network policy",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/component":  "master",
					"app.kubernetes.io/name":       "spark",
					"app.kubernetes.io/instance":   "test-id",
					"app.kubernetes.io/version":    "fake-tag",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
				},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/component":  "worker",
									"app.kubernetes.io/name":       "spark",
									"app.kubernetes.io/instance":   "test-id",
									"app.kubernetes.io/version":    "fake-tag",
									"app.kubernetes.io/managed-by": "distributed-compute-operator",
								},
							},
						},
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/instance": "spark-driver",
								},
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &protocol,
							Port:     &clusterPort,
						},
					},
				},
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"spark-client": "true",
								},
							},
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"domino-platform": "true",
								},
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &protocol,
							Port:     &masterDashboardPort,
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

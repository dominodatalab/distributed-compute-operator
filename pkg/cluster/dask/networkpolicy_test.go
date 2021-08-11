package dask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestNetworkPolicyDataSource(t *testing.T) {
	dc := testDaskCluster()
	tcpProto := corev1.ProtocolTCP
	dashboardPort := intstr.FromInt(8787)

	t.Run("scheduler", func(t *testing.T) {
		ds := networkPolicyDS{dc: dc, comp: ComponentScheduler}
		schedulerPort := intstr.FromInt(8786)

		actual := ds.NetworkPolicy()
		expected := &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-dask-scheduler",
				Namespace: "ns",
				Labels: map[string]string{
					"app.kubernetes.io/component":  "scheduler",
					"app.kubernetes.io/instance":   "test",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
					"app.kubernetes.io/name":       "dask",
					"app.kubernetes.io/version":    "test-tag",
				},
			},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app.kubernetes.io/component": "scheduler",
						"app.kubernetes.io/instance":  "test",
						"app.kubernetes.io/name":      "dask",
					},
				},
				Ingress: []networkingv1.NetworkPolicyIngressRule{
					{
						Ports: []networkingv1.NetworkPolicyPort{
							{
								Protocol: &tcpProto,
								Port:     &schedulerPort,
							},
						},
						From: []networkingv1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"test-client": "true",
									},
								},
							},
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/component": "worker",
										"app.kubernetes.io/instance":  "test",
										"app.kubernetes.io/name":      "dask",
									},
								},
							},
						},
					},
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
										"test-ui-client": "true",
									},
								},
							},
						},
					},
				},
				PolicyTypes: []networkingv1.PolicyType{
					networkingv1.PolicyTypeIngress,
				},
			},
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("worker", func(t *testing.T) {
		ds := networkPolicyDS{dc: dc, comp: ComponentWorker}
		workerPort := intstr.FromInt(3000)
		nannyPort := intstr.FromInt(3001)

		actual := ds.NetworkPolicy()
		expected := &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-dask-worker",
				Namespace: "ns",
				Labels: map[string]string{
					"app.kubernetes.io/component":  "worker",
					"app.kubernetes.io/instance":   "test",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
					"app.kubernetes.io/name":       "dask",
					"app.kubernetes.io/version":    "test-tag",
				},
			},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app.kubernetes.io/component": "worker",
						"app.kubernetes.io/instance":  "test",
						"app.kubernetes.io/name":      "dask",
					},
				},
				Ingress: []networkingv1.NetworkPolicyIngressRule{
					{
						Ports: []networkingv1.NetworkPolicyPort{
							{
								Protocol: &tcpProto,
								Port:     &workerPort,
							},
						},
						From: []networkingv1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/component": "scheduler",
										"app.kubernetes.io/instance":  "test",
										"app.kubernetes.io/name":      "dask",
									},
								},
							},
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/component": "worker",
										"app.kubernetes.io/instance":  "test",
										"app.kubernetes.io/name":      "dask",
									},
								},
							},
						},
					},
					{
						Ports: []networkingv1.NetworkPolicyPort{
							{
								Protocol: &tcpProto,
								Port:     &nannyPort,
							},
						},
						From: []networkingv1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"app.kubernetes.io/component": "scheduler",
										"app.kubernetes.io/instance":  "test",
										"app.kubernetes.io/name":      "dask",
									},
								},
							},
						},
					},
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
										"app.kubernetes.io/component": "scheduler",
										"app.kubernetes.io/instance":  "test",
										"app.kubernetes.io/name":      "dask",
									},
								},
							},
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"test-ui-client": "true",
									},
								},
							},
						},
					},
				},
				PolicyTypes: []networkingv1.PolicyType{
					networkingv1.PolicyTypeIngress,
				},
			},
		}

		assert.Equal(t, expected, actual)
	})
}

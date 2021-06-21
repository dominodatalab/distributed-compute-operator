package spark

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources"
)

// const (
//	descriptionCluster   = "Allows all ingress traffic between cluster nodes"
//	descriptionExternal  = "Allows all ingress traffic between cluster and external nodes"
//	descriptionClient    = "Allows client ingress traffic to head client server port"
//	descriptionDashboard = "Allows client ingress traffic to head dashboard port"
// )

func NewClusterWorkerNetworkPolicy(sc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
	workerSelector := metav1.LabelSelector{
		MatchLabels: MetadataLabelsWithComponent(sc, ComponentWorker),
	}

	driverSelector := metav1.LabelSelector{
		MatchLabels: sc.Spec.NetworkPolicy.ExternalPodLabels,
	}

	clusterSelector := metav1.LabelSelector{
		MatchLabels: SelectorLabels(sc),
	}

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(sc.Name, ComponentWorker),
			Namespace: sc.Namespace,
			Labels:    MetadataLabels(sc),
			Annotations: map[string]string{
				resources.DescriptionAnnotationKey: "worker network policy",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: workerSelector,
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &clusterSelector,
						},
						{
							PodSelector: &driverSelector,
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}
}

func NewClusterDriverNetworkPolicy(sc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
	driverSelector := metav1.LabelSelector{
		MatchLabels: sc.Spec.NetworkPolicy.ExternalPodLabels,
	}

	clusterSelector := metav1.LabelSelector{
		MatchLabels: SelectorLabels(sc),
	}

	protocol := corev1.ProtocolTCP
	driverPort := intstr.FromInt(int(sc.Spec.Driver.DriverPort))

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(sc.Name, "driver"),
			Namespace: sc.Namespace,
			Labels:    MetadataLabels(sc),
			Annotations: map[string]string{
				resources.DescriptionAnnotationKey: "driver network policy",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: driverSelector,
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &clusterSelector,
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &protocol,
							Port:     &driverPort,
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}
}

func NewClusterMasterNetworkPolicy(sc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
	driverSelector := metav1.LabelSelector{
		MatchLabels: sc.Spec.NetworkPolicy.ExternalPodLabels,
	}

	masterSelector := metav1.LabelSelector{
		MatchLabels: MetadataLabelsWithComponent(sc, ComponentMaster),
	}

	workerSelector := metav1.LabelSelector{
		MatchLabels: MetadataLabelsWithComponent(sc, ComponentWorker),
	}

	dashboardSelector := metav1.LabelSelector{
		MatchLabels: sc.Spec.NetworkPolicy.DashboardLabels,
	}

	namespaceSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"domino-platform": "true",
		},
	}

	protocol := corev1.ProtocolTCP
	dashboardPort := intstr.FromInt(int(sc.Spec.DashboardPort))
	clusterPort := intstr.FromInt(int(sc.Spec.ClusterPort))

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(sc.Name, "master"),
			Namespace: sc.Namespace,
			Labels:    MetadataLabels(sc),
			Annotations: map[string]string{
				resources.DescriptionAnnotationKey: "master network policy",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: masterSelector,
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &workerSelector,
						},
						{
							PodSelector: &driverSelector,
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
							NamespaceSelector: &namespaceSelector,
							PodSelector:       &dashboardSelector,
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &protocol,
							Port:     &dashboardPort,
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}
}

// func NewClusterExternalNetworkPolicy(sc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
//	driverSelector := metav1.LabelSelector{
//		MatchLabels: sc.Spec.NetworkPolicy.ExternalPodLabels,
//	}
//
//	clusterSelector := metav1.LabelSelector{
//		MatchLabels: SelectorLabels(sc),
//	}
//
//	return &networkingv1.NetworkPolicy{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      InstanceObjectName(sc.Name, Component("external")),
//			Namespace: sc.Namespace,
//			Labels:    MetadataLabels(sc),
//			Annotations: map[string]string{
//				resources.DescriptionAnnotationKey: descriptionExternal,
//			},
//		},
//		Spec: networkingv1.NetworkPolicySpec{
//			PodSelector: driverSelector,
//			Ingress: []networkingv1.NetworkPolicyIngressRule{
//				{
//					From: []networkingv1.NetworkPolicyPeer{
//						{
//							PodSelector: &clusterSelector,
//						},
//					},
//				},
//			},
//			PolicyTypes: []networkingv1.PolicyType{
//				networkingv1.PolicyTypeIngress,
//			},
//		},
//	}
// }

// NewClusterNetworkPolicy generates a network policy that allows all nodes
// within a single cluster to communicate on all ports. Optionally, it will also
// allow pods external to the cluster itself to communicate in and out of cluster
// func NewClusterNetworkPolicy(sc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
//	labelSelector := metav1.LabelSelector{
//		MatchLabels: SelectorLabels(sc),
//	}
//
//	peers := []networkingv1.NetworkPolicyPeer{
//		{
//			PodSelector: &labelSelector,
//		},
//	}
//
//	if sc.Spec.NetworkPolicy.ExternalPolicyEnabled != nil {
//		peers = append(peers, networkingv1.NetworkPolicyPeer{
//			PodSelector: &metav1.LabelSelector{
//				MatchLabels: sc.Spec.NetworkPolicy.ExternalPodLabels,
//			},
//		})
//	}
//
//	return &networkingv1.NetworkPolicy{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      InstanceObjectName(sc.Name, Component("cluster")),
//			Namespace: sc.Namespace,
//			Labels:    MetadataLabels(sc),
//			Annotations: map[string]string{
//				resources.DescriptionAnnotationKey: descriptionCluster,
//			},
//		},
//		Spec: networkingv1.NetworkPolicySpec{
//			PodSelector: labelSelector,
//			Ingress: []networkingv1.NetworkPolicyIngressRule{
//				{
//					From: peers,
//				},
//			},
//			PolicyTypes: []networkingv1.PolicyType{
//				networkingv1.PolicyTypeIngress,
//			},
//		},
//	}
// }

// NewHeadClientNetworkPolicy generates a network policy that allows client
// access to any pods that have been appointed with the configured client
// server labels.
// func NewHeadClientNetworkPolicy(sc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
//	return headNetworkPolicy(
//		sc,
//		sc.Spec.ClusterPort,
//		sc.Spec.NetworkPolicy.ClientServerLabels,
//		"client",
//		descriptionClient,
//	)
// }

// NewHeadDashboardNetworkPolicy generates a network policy that allows
// dashboard access to any pods that have been appointed with configured
// dashboard labels.
// func NewHeadDashboardNetworkPolicy(sc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
//	return headNetworkPolicy(
//		sc,
//		sc.Spec.DashboardPort,
//		sc.Spec.NetworkPolicy.DashboardLabels,
//		"dashboard",
//		descriptionDashboard,
//	)
// }

// func headNetworkPolicy(sc *dcv1alpha1.SparkCluster, p int32, l map[string]string, c Component, desc string) *networkingv1.NetworkPolicy {
//	proto := corev1.ProtocolTCP
//	targetPort := intstr.FromInt(int(p))
//
//	return &networkingv1.NetworkPolicy{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      InstanceObjectName(sc.Name, c),
//			Namespace: sc.Namespace,
//			Labels:    MetadataLabelsWithComponent(sc, ComponentMaster),
//			Annotations: map[string]string{
//				resources.DescriptionAnnotationKey: desc,
//			},
//		},
//		Spec: networkingv1.NetworkPolicySpec{
//			PodSelector: metav1.LabelSelector{
//				MatchLabels: SelectorLabelsWithComponent(sc, ComponentMaster),
//			},
//			Ingress: []networkingv1.NetworkPolicyIngressRule{
//				{
//					Ports: []networkingv1.NetworkPolicyPort{
//						{
//							Protocol: &proto,
//							Port:     &targetPort,
//						},
//					},
//					From: []networkingv1.NetworkPolicyPeer{
//						{
//							PodSelector: &metav1.LabelSelector{
//								MatchLabels: l,
//							},
//						},
//					},
//				},
//			},
//			PolicyTypes: []networkingv1.PolicyType{
//				networkingv1.PolicyTypeIngress,
//			},
//		},
//	}
// }

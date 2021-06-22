package spark

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources"
)

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
	driverUIPort := intstr.FromInt(int(sc.Spec.Driver.DriverUIPort))

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
							Port:     &driverUIPort,
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
	masterDashboardPort := intstr.FromInt(int(sc.Spec.TCPMasterWebPort))
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
							Port:     &masterDashboardPort,
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

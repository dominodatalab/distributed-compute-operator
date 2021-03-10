package spark

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources"
)

// NewClusterNetworkPolicy generates a network policy that allows all nodes
// within a single cluster to communicate on all ports.
func NewClusterNetworkPolicy(rc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
	labelSelector := metav1.LabelSelector{
		MatchLabels: SelectorLabels(rc),
	}

	return &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkPolicy",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, Component("cluster")),
			Namespace: rc.Namespace,
			Labels:    MetadataLabels(rc),
			Annotations: map[string]string{
				resources.DescriptionAnnotationKey: "Allows all ingress traffic between cluster nodes",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: labelSelector,
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &labelSelector,
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

// NewHeadNetworkPolicy generates a network policy that allows client/dashboard
// port access to any pods that have been appointed with the ClientAccessLabels.
func NewHeadNetworkPolicy(rc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
	proto := corev1.ProtocolTCP
	clientPort := intstr.FromInt(int(rc.Spec.ClientServerPort))
	dashboardPort := intstr.FromInt(int(rc.Spec.DashboardPort))

	var policyPeers []networkingv1.NetworkPolicyPeer
	for _, labels := range rc.Spec.NetworkPolicyClientLabels {
		policyPeers = append(policyPeers, networkingv1.NetworkPolicyPeer{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		})
	}

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, Component("client")),
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentHead),
			Annotations: map[string]string{
				resources.DescriptionAnnotationKey: "Allows client ingress traffic to head node",
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: SelectorLabelsWithComponent(rc, ComponentHead),
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &proto,
							Port:     &clientPort,
						},
						{
							Protocol: &proto,
							Port:     &dashboardPort,
						},
					},
					From: policyPeers,
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}
}

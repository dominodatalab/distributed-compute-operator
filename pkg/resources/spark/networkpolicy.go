package spark

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources"
)

const (
	descriptionCluster   = "Allows all ingress traffic between cluster nodes"
	descriptionClient    = "Allows client ingress traffic to head client server port"
	descriptionDashboard = "Allows client ingress traffic to head dashboard port"
)

// NewClusterNetworkPolicy generates a network policy that allows all nodes
// within a single cluster to communicate on all ports.
func NewClusterNetworkPolicy(rc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
	labelSelector := metav1.LabelSelector{
		MatchLabels: SelectorLabels(rc),
	}

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, Component("cluster")),
			Namespace: rc.Namespace,
			Labels:    MetadataLabels(rc),
			Annotations: map[string]string{
				resources.DescriptionAnnotationKey: descriptionCluster,
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

// NewHeadClientNetworkPolicy generates a network policy that allows client
// access to any pods that have been appointed with the configured client
// server labels.
func NewHeadClientNetworkPolicy(rc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
	return headNetworkPolicy(
		rc,
		rc.Spec.ClusterPort,
		rc.Spec.NetworkPolicy.ClientServerLabels,
		Component("client"),
		descriptionClient,
	)
}

// NewHeadDashboardNetworkPolicy generates a network policy that allows
// dashboard access to any pods that have been appointed with configured
// dashboard labels.
func NewHeadDashboardNetworkPolicy(rc *dcv1alpha1.SparkCluster) *networkingv1.NetworkPolicy {
	return headNetworkPolicy(
		rc,
		rc.Spec.DashboardPort,
		rc.Spec.NetworkPolicy.DashboardLabels,
		Component("dashboard"),
		descriptionDashboard,
	)
}

func headNetworkPolicy(rc *dcv1alpha1.SparkCluster, p int32, l map[string]string, c Component, desc string) *networkingv1.NetworkPolicy {
	proto := corev1.ProtocolTCP
	targetPort := intstr.FromInt(int(p))

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(rc.Name, c),
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentMaster),
			Annotations: map[string]string{
				resources.DescriptionAnnotationKey: desc,
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: SelectorLabelsWithComponent(rc, ComponentMaster),
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &proto,
							Port:     &targetPort,
						},
					},
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: l,
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
}

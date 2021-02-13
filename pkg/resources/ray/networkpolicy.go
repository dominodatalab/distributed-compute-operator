package ray

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources"
)

func NewNetworkPolicy(rc *dcv1alpha1.RayCluster) *networkingv1.NetworkPolicy {
	labelSelector := metav1.LabelSelector{
		MatchLabels: SelectorLabels(rc),
	}

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-cluster", rc.Name),
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

package dask

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

func NetworkPolicyScheduler() core.Component {
	return components.NetworkPolicy(
		schedulerNetworkPolicy,
		shouldDeleteNetworkPolicy,
	)
}

func NetworkPolicyWorker() core.Component {
	return components.NetworkPolicy(
		workerNetworkPolicy,
		shouldDeleteNetworkPolicy,
	)
}

func schedulerNetworkPolicy(obj client.Object) *networkingv1.NetworkPolicy {
	dc := obj.(*dcv1alpha1.DaskCluster)
	sPort := intstr.FromInt(int(dc.Spec.SchedulerPort))
	dPort := intstr.FromInt(int(dc.Spec.DashboardPort))
	proto := corev1.ProtocolTCP

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(dc, ComponentScheduler),
			Namespace: dc.Namespace,
			Labels:    meta.StandardLabelsWithComponent(dc, ComponentScheduler),
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: meta.MatchLabelsWithComponent(dc, ComponentScheduler),
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: dc.Spec.NetworkPolicy.ClientLabels,
							},
						},
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: meta.MatchLabelsWithComponent(dc, ComponentWorker),
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Port:     &sPort,
							Protocol: &proto,
						},
					},
				},
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: dc.Spec.NetworkPolicy.DashboardLabels,
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Port:     &dPort,
							Protocol: &proto,
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

func workerNetworkPolicy(obj client.Object) *networkingv1.NetworkPolicy {
	dc := obj.(*dcv1alpha1.DaskCluster)
	wPort := intstr.FromInt(int(dc.Spec.WorkerPort))
	nPort := intstr.FromInt(int(dc.Spec.NannyPort))
	dPort := intstr.FromInt(int(dc.Spec.DashboardPort))
	proto := corev1.ProtocolTCP

	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(dc, ComponentWorker),
			Namespace: dc.Namespace,
			Labels:    meta.StandardLabelsWithComponent(dc, ComponentWorker),
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: meta.MatchLabelsWithComponent(dc, ComponentWorker),
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: meta.MatchLabelsWithComponent(dc, ComponentScheduler),
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Port:     &wPort,
							Protocol: &proto,
						},
						{
							Port:     &nPort,
							Protocol: &proto,
						},
						{
							Port:     &dPort,
							Protocol: &proto,
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

func shouldDeleteNetworkPolicy(obj client.Object) bool {
	dc := obj.(*dcv1alpha1.DaskCluster)
	return util.BoolPtrIsNilOrFalse(dc.Spec.NetworkPolicy.Enabled)
}

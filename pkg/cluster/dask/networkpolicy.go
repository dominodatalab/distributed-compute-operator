package dask

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util"
)

func NetworkPolicyScheduler() core.OwnedComponent {
	return components.NetworkPolicy(func(obj client.Object) components.NetworkPolicyDataSource {
		return &networkPolicyDS{dc: daskCluster(obj), comp: ComponentScheduler}
	})
}

func NetworkPolicyWorker() core.OwnedComponent {
	return components.NetworkPolicy(func(obj client.Object) components.NetworkPolicyDataSource {
		return &networkPolicyDS{dc: daskCluster(obj), comp: ComponentWorker}
	})
}

type networkPolicyDS struct {
	dc   *dcv1alpha1.DaskCluster
	comp metadata.Component
}

func (s *networkPolicyDS) NetworkPolicy() *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(s.dc, s.comp),
			Namespace: s.dc.Namespace,
			Labels:    meta.StandardLabelsWithComponent(s.dc, s.comp, nil),
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: meta.MatchLabelsWithComponent(s.dc, s.comp),
			},
			Ingress: s.ingressRules(),
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}
}

func (s *networkPolicyDS) Delete() bool {
	return util.BoolPtrIsNilOrFalse(s.dc.Spec.NetworkPolicy.Enabled)
}

func (s *networkPolicyDS) ingressRules() []networkingv1.NetworkPolicyIngressRule {
	tcpProto := corev1.ProtocolTCP
	dashboardPort := intstr.FromInt(int(s.dc.Spec.DashboardPort))

	if s.comp == ComponentScheduler {
		sPort := intstr.FromInt(int(s.dc.Spec.SchedulerPort))

		return []networkingv1.NetworkPolicyIngressRule{
			{
				From: []networkingv1.NetworkPolicyPeer{
					{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: s.dc.Spec.NetworkPolicy.ClientLabels,
						},
					},
					{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: meta.MatchLabelsWithComponent(s.dc, ComponentWorker),
						},
					},
				},
				Ports: []networkingv1.NetworkPolicyPort{
					{
						Port:     &sPort,
						Protocol: &tcpProto,
					},
				},
			},
			{
				From: []networkingv1.NetworkPolicyPeer{
					{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: s.dc.Spec.NetworkPolicy.DashboardLabels,
						},
						NamespaceSelector: &metav1.LabelSelector{
							MatchLabels: s.dc.Spec.NetworkPolicy.DashboardNamespaceLabels,
						},
					},
				},
				Ports: []networkingv1.NetworkPolicyPort{
					{
						Port:     &dashboardPort,
						Protocol: &tcpProto,
					},
				},
			},
		}
	}

	workerPort := intstr.FromInt(int(s.dc.Spec.WorkerPort))
	nannyPort := intstr.FromInt(int(s.dc.Spec.NannyPort))

	return []networkingv1.NetworkPolicyIngressRule{
		{
			From: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: meta.MatchLabelsWithComponent(s.dc, ComponentScheduler),
					},
				},
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: meta.MatchLabelsWithComponent(s.dc, ComponentWorker),
					},
				},
			},
			Ports: []networkingv1.NetworkPolicyPort{
				{
					Port:     &workerPort,
					Protocol: &tcpProto,
				},
			},
		},
		{
			From: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: meta.MatchLabelsWithComponent(s.dc, ComponentScheduler),
					},
				},
			},
			Ports: []networkingv1.NetworkPolicyPort{
				{
					Port:     &nannyPort,
					Protocol: &tcpProto,
				},
			},
		},
		{
			From: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: meta.MatchLabelsWithComponent(s.dc, ComponentScheduler),
					},
				},
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: s.dc.Spec.NetworkPolicy.DashboardLabels,
					},
					NamespaceSelector: &metav1.LabelSelector{
						MatchLabels: s.dc.Spec.NetworkPolicy.DashboardNamespaceLabels,
					},
				},
			},
			Ports: []networkingv1.NetworkPolicyPort{
				{
					Port:     &dashboardPort,
					Protocol: &tcpProto,
				},
			},
		},
	}
}

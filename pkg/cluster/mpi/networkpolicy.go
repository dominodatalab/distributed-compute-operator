package mpi

import (
	"errors"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func WorkerNetworkPolicy() core.OwnedComponent {
	return &networkPolicyComponent{
		comp: ComponentWorker,
	}
}

func DriverNetworkPolicy() core.OwnedComponent {
	return &networkPolicyComponent{
		comp: ComponentDriver,
	}
}

type networkPolicyComponent struct {
	comp metadata.Component
}

func (c networkPolicyComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	matchLabels := meta.MatchLabels(cr)

	var podSelectorMatchLabels map[string]string
	var ingressRules []networkingv1.NetworkPolicyPeer

	switch c.comp {
	case ComponentWorker:
		podSelectorMatchLabels = matchLabels
		ingressRules = []networkingv1.NetworkPolicyPeer{
			{
				PodSelector: &metav1.LabelSelector{
					MatchLabels: matchLabels,
				},
			},
			{
				PodSelector: &metav1.LabelSelector{
					MatchLabels: cr.Spec.NetworkPolicy.ClientLabels,
				},
			},
		}
	case ComponentDriver:
		podSelectorMatchLabels = cr.Spec.NetworkPolicy.ClientLabels
		ingressRules = []networkingv1.NetworkPolicyPeer{
			{
				PodSelector: &metav1.LabelSelector{
					MatchLabels: matchLabels,
				},
			},
		}
	case metadata.ComponentNone:
		err := errors.New("unknown component for NetworkPolicy")
		return ctrl.Result{}, err
	}

	netpol := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(cr, c.comp),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabels(cr),
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: podSelectorMatchLabels,
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: ingressRules,
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, cr, netpol)
	if err != nil {
		err = fmt.Errorf("cannot reconcile networkpolicy: %w", err)
	}

	return ctrl.Result{}, err
}

func (c networkPolicyComponent) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

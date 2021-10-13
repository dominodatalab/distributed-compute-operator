package mpi

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func NetworkPolicy() core.OwnedComponent {
	return &networkPolicyComponent{}
}

type networkPolicyComponent struct{}

func (c networkPolicyComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPIJob(ctx.Object)

	matchLabels := meta.MatchLabels(cr)
	netpol := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(cr, metadata.ComponentNone),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabels(cr),
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: matchLabels,
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: matchLabels,
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

	err := actions.CreateOrUpdateOwnedResource(ctx, cr, netpol)
	if err != nil {
		err = fmt.Errorf("cannot reconcile networkpolicy: %w", err)
	}

	return ctrl.Result{}, err
}

func (c networkPolicyComponent) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

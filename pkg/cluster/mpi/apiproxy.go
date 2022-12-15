package mpi

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func APIProxyService() core.OwnedComponent {
	return components.APIProxyServiceComponent{
		ClientPorts: func(obj *client.Object) []corev1.ServicePort {
			return objToMPICluster(*obj).Spec.AdditionalClientPorts
		},
		ClientLabels: func(obj *client.Object) map[string]string {
			return objToMPICluster(*obj).Spec.NetworkPolicy.ClientLabels
		},
		Meta: meta,
	}
}

func APIProxyNetworkPolicy() core.OwnedComponent {
	return components.APIProxyNetworkPolicyComponent{
		ClientPorts: func(obj *client.Object) []corev1.ServicePort {
			return objToMPICluster(*obj).Spec.AdditionalClientPorts
		},
		ClientLabels: func(obj *client.Object) map[string]string {
			return objToMPICluster(*obj).Spec.NetworkPolicy.ClientLabels
		},
		Meta: meta,
	}
}

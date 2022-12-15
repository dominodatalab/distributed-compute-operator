package dask

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func ClientPortsService() core.OwnedComponent {
	return components.ClientPortsServiceComponent{
		ClientPorts: func(obj *client.Object) []corev1.ServicePort {
			return daskCluster(*obj).Spec.AdditionalClientPorts
		},
		ClientLabels: func(obj *client.Object) map[string]string {
			return daskCluster(*obj).Spec.NetworkPolicy.ClientLabels
		},
		Meta: meta,
	}
}

func ClientPortsNetworkPolicy() core.OwnedComponent {
	return components.ClientPortsNetworkPolicyComponent{
		ClientPorts: func(obj *client.Object) []corev1.ServicePort {
			return daskCluster(*obj).Spec.AdditionalClientPorts
		},
		ClientLabels: func(obj *client.Object) map[string]string {
			return daskCluster(*obj).Spec.NetworkPolicy.ClientLabels
		},
		Meta: meta,
	}
}

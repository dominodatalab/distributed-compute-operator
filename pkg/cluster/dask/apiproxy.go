package dask

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func APIProxyService() core.OwnedComponent {
	return components.APIProxyServiceComponent{
		APIProxyPort: func(obj *client.Object) int32 {
			return daskCluster(*obj).Spec.APIProxyPort
		},
		ClientLabels: func(obj *client.Object) map[string]string {
			return daskCluster(*obj).Spec.NetworkPolicy.ClientLabels
		},
		Meta: meta,
	}
}

func APIProxyNetworkPolicy() core.OwnedComponent {
	return components.APIProxyNetworkPolicyComponent{
		APIProxyPort: func(obj *client.Object) int32 {
			return daskCluster(*obj).Spec.APIProxyPort
		},
		ClientLabels: func(obj *client.Object) map[string]string {
			return daskCluster(*obj).Spec.NetworkPolicy.ClientLabels
		},
		Meta: meta,
	}
}

package spark

import (
	"fmt"

	protobuftypes "github.com/gogo/protobuf/types"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	apinetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

const filterName = "envoy.filters.network.tcp_proxy"

// NewEnvoyFilter creates a new EnvoyFilter resource to set idle_timeout for Istio-enabled deployments
func NewEnvoyFilter(sc *dcv1alpha1.SparkCluster) *apinetworkingv1alpha3.EnvoyFilter {
	match := networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
		Context: networkingv1alpha3.EnvoyFilter_ANY,
		ObjectTypes: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
			Listener: &networkingv1alpha3.EnvoyFilter_ListenerMatch{
				FilterChain: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
					Filter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
						Name: filterName,
					},
				},
			},
		},
	}

	patch := networkingv1alpha3.EnvoyFilter_Patch{
		Operation: networkingv1alpha3.EnvoyFilter_Patch_MERGE,
		Value: &protobuftypes.Struct{
			Fields: map[string]*protobuftypes.Value{
				"name": {
					Kind: &protobuftypes.Value_StringValue{
						StringValue: "envoy.filters.network.tcp_proxy",
					},
				},
				"typed_config": {
					Kind: &protobuftypes.Value_StructValue{
						StructValue: &protobuftypes.Struct{
							Fields: map[string]*protobuftypes.Value{
								"@type": {
									Kind: &protobuftypes.Value_StringValue{
										StringValue: "type.googleapis.com/envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy",
									},
								},
								"idle_timeout": {
									Kind: &protobuftypes.Value_StringValue{
										StringValue: "0s",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	configPatches := []*networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
		{
			ApplyTo: networkingv1alpha3.EnvoyFilter_NETWORK_FILTER,
			Match:   &match,
			Patch:   &patch,
		},
	}

	workloadSelector := networkingv1alpha3.WorkloadSelector{
		Labels: sc.Spec.EnvoyFilterLabels,
	}

	envoyFilter := &apinetworkingv1alpha3.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", InstanceObjectName(sc.Name, ComponentNone), "envoyfilter"),
			Namespace: sc.Namespace,
			Labels:    AddGlobalLabels(MetadataLabels(sc), sc.Labels),
		},
		Spec: networkingv1alpha3.EnvoyFilter{
			WorkloadSelector: &workloadSelector,
			ConfigPatches:    configPatches,
		},
	}

	return envoyFilter
}

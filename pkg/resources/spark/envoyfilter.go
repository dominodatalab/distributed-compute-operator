package spark

import (
	"fmt"

	protobuftypes "github.com/gogo/protobuf/types"
	"istio.io/api/networking/v1alpha3"
	v1alpha32 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

const (
	filterName = "envoy.filters.network.tcp_proxy"
)

// NewEnvoyFilter creates a new EnvoyFilter resource to set idle_timeout for Istio-enabled deployments
func NewEnvoyFilter(sc *dcv1alpha1.SparkCluster) (v1alpha32.EnvoyFilter, error) {
	match := v1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
		Context: v1alpha3.EnvoyFilter_ANY,
		ObjectTypes: &v1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
			Listener: &v1alpha3.EnvoyFilter_ListenerMatch{
				FilterChain: &v1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
					Filter: &v1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
						Name: filterName,
					},
				},
			},
		},
	}

	// matchOutbound := v1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
	//	Context: v1alpha3.EnvoyFilter_SIDECAR_OUTBOUND,
	//	ObjectTypes: &v1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
	//		Listener: &v1alpha3.EnvoyFilter_ListenerMatch{
	//			FilterChain: &v1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
	//				Filter: &v1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
	//					Name: filterName,
	//				},
	//			},
	//		},
	//	},
	// }

	patch := v1alpha3.EnvoyFilter_Patch{
		Operation: v1alpha3.EnvoyFilter_Patch_MERGE,
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

	// create private function to take in match, one for in/outbound
	configPatches := []*v1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
		{
			ApplyTo: v1alpha3.EnvoyFilter_NETWORK_FILTER,
			Match:   &match,
			Patch:   &patch,
		},
	}

	workloadSelector := v1alpha3.WorkloadSelector{
		Labels: sc.Spec.GlobalLabels,
	}

	envoyFilter := v1alpha32.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", InstanceObjectName(sc.Name, ComponentNone), "envoyfilter"),
			Namespace: sc.Namespace,
			Labels:    AddGlobalLabels(MetadataLabels(sc), sc.Labels),
		},
		Spec: v1alpha3.EnvoyFilter{
			WorkloadSelector: &workloadSelector,
			ConfigPatches:    configPatches,
		},
	}

	return envoyFilter, nil
}

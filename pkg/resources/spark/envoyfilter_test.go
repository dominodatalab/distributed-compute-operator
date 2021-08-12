package spark

import (
	"testing"

	protobuftypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	apinetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewEnvoyFilter(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		sc := sparkClusterFixture()
		actual := NewEnvoyFilter(sc)

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
				Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
					Context: networkingv1alpha3.EnvoyFilter_ANY,
					ObjectTypes: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
						Listener: &networkingv1alpha3.EnvoyFilter_ListenerMatch{
							FilterChain: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
								Filter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
									Name: "envoy.filters.network.tcp_proxy",
								},
							},
						},
					},
				},
				Patch: &patch,
			},
		}

		workloadSelector := networkingv1alpha3.WorkloadSelector{
			Labels: sc.Spec.EnvoyFilterLabels,
		}

		expected := &apinetworkingv1alpha3.EnvoyFilter{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-id-spark-envoyfilter",
				Namespace: sc.Namespace,
				Labels:    MetadataLabels(sc),
			},
			Spec: networkingv1alpha3.EnvoyFilter{
				WorkloadSelector: &workloadSelector,
				ConfigPatches:    configPatches,
			},
		}

		assert.Equal(t, expected, actual, "Istio EnvoyFilter not correctly generated")
	})
}

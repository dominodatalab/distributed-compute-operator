package spark

import (
	"testing"

	protobuftypes "github.com/gogo/protobuf/types"
	v1alpha32 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"istio.io/api/networking/v1alpha3"
)

func TestNewEnvoyFilter(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		sc := sparkClusterFixture()
		actual, err := NewEnvoyFilter(sc)
		require.NoError(t, err)

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
											StringValue: "30s",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		configPatches := []*v1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
			{
				ApplyTo: v1alpha3.EnvoyFilter_NETWORK_FILTER,
				Match: &v1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
					Context: v1alpha3.EnvoyFilter_SIDECAR_INBOUND,
					ObjectTypes: &v1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
						Listener: &v1alpha3.EnvoyFilter_ListenerMatch{
							FilterChain: &v1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
								Filter: &v1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
									Name: "envoy.filters.network.tcp_proxy",
								},
							},
						},
					},
				},
				Patch: &patch,
			},
			{
				ApplyTo: v1alpha3.EnvoyFilter_NETWORK_FILTER,
				Match: &v1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
					Context: v1alpha3.EnvoyFilter_SIDECAR_OUTBOUND,
					ObjectTypes: &v1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
						Listener: &v1alpha3.EnvoyFilter_ListenerMatch{
							FilterChain: &v1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
								Filter: &v1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
									Name: "envoy.filters.network.tcp_proxy",
								},
							},
						},
					},
				},
				Patch: &patch,
			},
		}

		workloadSelector := v1alpha3.WorkloadSelector{
			Labels: sc.Labels,
		}

		expected := v1alpha32.EnvoyFilter{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-id-spark-envoyfilter",
				Namespace: sc.Namespace,
				Labels:    MetadataLabels(sc),
			},
			Spec: v1alpha3.EnvoyFilter{
				WorkloadSelector: &workloadSelector,
				ConfigPatches:    configPatches,
			},
		}

		assert.Equal(t, expected, actual, "Istio EnvoyFilter not correctly generated")
	})
}

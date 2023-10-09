package spark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	spb "google.golang.org/protobuf/types/known/structpb"
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
			Value: &spb.Struct{
				Fields: map[string]*spb.Value{
					"name": {
						Kind: &spb.Value_StringValue{
							StringValue: "envoy.filters.network.tcp_proxy",
						},
					},
					"typed_config": {
						Kind: &spb.Value_StructValue{
							StructValue: &spb.Struct{
								Fields: map[string]*spb.Value{
									"@type": {
										Kind: &spb.Value_StringValue{
											StringValue: "type.googleapis.com/envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy",
										},
									},
									"idle_timeout": {
										Kind: &spb.Value_StringValue{
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
			{
				ApplyTo: networkingv1alpha3.EnvoyFilter_NETWORK_FILTER,
				Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
					Context: networkingv1alpha3.EnvoyFilter_ANY,
					ObjectTypes: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
						Listener: &networkingv1alpha3.EnvoyFilter_ListenerMatch{
							FilterChain: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
								Filter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
									Name: "envoy.filters.network.http_connection_manager",
								},
							},
						},
					},
				},
				Patch: &networkingv1alpha3.EnvoyFilter_Patch{
					Operation: networkingv1alpha3.EnvoyFilter_Patch_MERGE,
					Value: &spb.Struct{
						Fields: map[string]*spb.Value{
							"name": {
								Kind: &spb.Value_StringValue{
									StringValue: "envoy.filters.network.http_connection_manager",
								},
							},
							"typed_config": {
								Kind: &spb.Value_StructValue{
									StructValue: &spb.Struct{
										Fields: map[string]*spb.Value{
											"@type": {
												Kind: &spb.Value_StringValue{
													StringValue: "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
												},
											},
											"idle_timeout": {
												Kind: &spb.Value_StringValue{
													StringValue: "0s",
												},
											},
										},
									},
								},
							},
						},
					},
				},
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

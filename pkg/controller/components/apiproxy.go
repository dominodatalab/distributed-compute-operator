package components

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

const (
	defaultAPIProxyPort = 8899
	targetServicePort   = 8899
	apiProxyPortName    = "http-api-proxy"
	component           = "api-proxy"
)

func newResourceMeta(obj *client.Object, componentMeta *metadata.Provider) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      componentMeta.InstanceName(*obj, component),
		Namespace: (*obj).GetNamespace(),
		Labels:    componentMeta.StandardLabelsWithComponent(*obj, component, nil),
	}
}

func NewAPIProxyServiceComponent(
	obj *client.Object,
	port int32,
	clientLabels map[string]string,
	meta *metadata.Provider) *corev1.Service {
	var apiProxyPort int32
	if port != 0 {
		apiProxyPort = port
	} else {
		apiProxyPort = defaultAPIProxyPort
	}

	ports := []corev1.ServicePort{{
		Name:       apiProxyPortName,
		Port:       apiProxyPort,
		TargetPort: intstr.FromInt(targetServicePort),
		Protocol:   corev1.ProtocolTCP,
	}}

	return &corev1.Service{
		ObjectMeta: newResourceMeta(obj, meta),
		Spec: corev1.ServiceSpec{
			Selector: clientLabels,
			Ports:    ports,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}

type APIProxyServiceComponent struct {
	APIProxyPort func(obj *client.Object) int32
	ClientLabels func(obj *client.Object) map[string]string
	Meta         *metadata.Provider
}

func (c APIProxyServiceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object
	svc := NewAPIProxyServiceComponent(&obj, c.APIProxyPort(&obj), c.ClientLabels(&obj), c.Meta)

	err := actions.CreateOrUpdateOwnedResource(ctx, obj, svc)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

func (c APIProxyServiceComponent) Kind() client.Object {
	return &corev1.Service{}
}

type APIProxyNetworkPolicyComponent struct {
	APIProxyPort func(obj *client.Object) int32
	ClientLabels func(obj *client.Object) map[string]string
	Meta         *metadata.Provider
}

func NewAPIProxyNetworkPolicyComponent(
	obj *client.Object, port int32,
	clientLabels map[string]string,
	meta *metadata.Provider) *networkingv1.NetworkPolicy {
	var apiProxyPort intstr.IntOrString
	if port != 0 {
		apiProxyPort = intstr.FromInt(int(port))
	} else {
		apiProxyPort = intstr.FromInt(defaultAPIProxyPort)
	}

	tcpProto := corev1.ProtocolTCP

	ingressRules := []networkingv1.NetworkPolicyIngressRule{{
		From: []networkingv1.NetworkPolicyPeer{{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: meta.MatchLabels(*obj),
			},
		}},
		Ports: []networkingv1.NetworkPolicyPort{{
			Port:     &apiProxyPort,
			Protocol: &tcpProto,
		}},
	}}

	return &networkingv1.NetworkPolicy{
		ObjectMeta: newResourceMeta(obj, meta),
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: clientLabels,
			},
			Ingress: ingressRules,
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}
}

func (c APIProxyNetworkPolicyComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object
	netPol := NewAPIProxyNetworkPolicyComponent(&obj, c.APIProxyPort(&obj), c.ClientLabels(&obj), c.Meta)

	err := actions.CreateOrUpdateOwnedResource(ctx, obj, netPol)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

func (c APIProxyNetworkPolicyComponent) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

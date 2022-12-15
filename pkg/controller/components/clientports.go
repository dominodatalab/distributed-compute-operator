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
	component             = "proxy"
	defaultClientPort     = 8899
	defaultClientPortName = "http-proxy"
)

func newResourceMeta(obj *client.Object, componentMeta *metadata.Provider) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      componentMeta.InstanceName(*obj, component),
		Namespace: (*obj).GetNamespace(),
		Labels:    componentMeta.StandardLabelsWithComponent(*obj, component, nil),
	}
}

func NewClientPortsServiceComponent(
	obj *client.Object,
	ports []corev1.ServicePort,
	clientLabels map[string]string,
	meta *metadata.Provider) *corev1.Service {
	var clientPorts []corev1.ServicePort
	if len(ports) > 0 {
		clientPorts = ports
	} else {
		// TODO: Perhaps remove after Nucleus is set up to supply the API proxy port
		clientPorts = []corev1.ServicePort{{
			Name:       defaultClientPortName,
			Port:       defaultClientPort,
			TargetPort: intstr.FromInt(defaultClientPort),
			Protocol:   corev1.ProtocolTCP,
		}}
	}

	return &corev1.Service{
		ObjectMeta: newResourceMeta(obj, meta),
		Spec: corev1.ServiceSpec{
			Selector: clientLabels,
			Ports:    clientPorts,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}

type ClientPortsServiceComponent struct {
	ClientPorts  func(obj *client.Object) []corev1.ServicePort
	ClientLabels func(obj *client.Object) map[string]string
	Meta         *metadata.Provider
}

func (c ClientPortsServiceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object
	svc := NewClientPortsServiceComponent(&obj, c.ClientPorts(&obj), c.ClientLabels(&obj), c.Meta)

	err := actions.CreateOrUpdateOwnedResource(ctx, obj, svc)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

func (c ClientPortsServiceComponent) Kind() client.Object {
	return &corev1.Service{}
}

type ClientPortsNetworkPolicyComponent struct {
	ClientPorts  func(obj *client.Object) []corev1.ServicePort
	ClientLabels func(obj *client.Object) map[string]string
	Meta         *metadata.Provider
}

func NewClientPortsNetworkPolicyComponent(
	obj *client.Object,
	ports []corev1.ServicePort,
	clientLabels map[string]string,
	meta *metadata.Provider) *networkingv1.NetworkPolicy {
	var clientPorts []corev1.ServicePort
	if len(ports) > 0 {
		clientPorts = ports
	} else {
		// TODO: Perhaps remove after Nucleus is set up to supply the API proxy port
		clientPorts = []corev1.ServicePort{{
			Name:       defaultClientPortName,
			Port:       defaultClientPort,
			TargetPort: intstr.FromInt(defaultClientPort),
			Protocol:   corev1.ProtocolTCP,
		}}
	}

	portCount := len(clientPorts)
	networkPorts := make([]networkingv1.NetworkPolicyPort, portCount)
	for i := 0; i < portCount; i++ {
		networkPorts[i] = networkingv1.NetworkPolicyPort{
			Port:     &(clientPorts[i].TargetPort),
			Protocol: &(clientPorts[i].Protocol),
		}
	}

	ingressRules := []networkingv1.NetworkPolicyIngressRule{{
		From: []networkingv1.NetworkPolicyPeer{{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: meta.MatchLabels(*obj),
			},
		}},
		Ports: networkPorts,
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

func (c ClientPortsNetworkPolicyComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object
	netPol := NewClientPortsNetworkPolicyComponent(&obj, c.ClientPorts(&obj), c.ClientLabels(&obj), c.Meta)

	err := actions.CreateOrUpdateOwnedResource(ctx, obj, netPol)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

func (c ClientPortsNetworkPolicyComponent) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

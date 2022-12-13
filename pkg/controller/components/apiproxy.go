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
	defaultAPIProxyPort        = 8899
	targetServicePort          = 8899
	apiProxyPortName           = "api-proxy"
	component                  = "api-proxy"
	executionIDLabel           = "dominodatalab.com/execution-id"
	projectIDLabel             = "dominodatalab.com/project-id"
	datasourceProxyClientLabel = "datasource-proxy-client"
)

func executionID(obj *client.Object) string {
	return (*obj).GetLabels()[executionIDLabel]
}

func runPodSelector(obj *client.Object) map[string]string {
	runPodName := fmt.Sprintf("run-%s", executionID(obj))
	return map[string]string{
		metadata.ApplicationInstanceLabelKey: runPodName,
	}
}

func newResourceMeta(obj *client.Object, componentMeta *metadata.Provider) metav1.ObjectMeta {
	instanceName := fmt.Sprintf("%s-%s", component, executionID(obj))
	return metav1.ObjectMeta{
		Name:      instanceName,
		Namespace: (*obj).GetNamespace(),
		Labels:    componentMeta.StandardLabelsWithComponent(*obj, component, nil),
	}
}

func NewAPIProxyServiceComponent(obj *client.Object, port int32, meta *metadata.Provider) *corev1.Service {
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

	internalTrafficPolicy := corev1.ServiceInternalTrafficPolicyCluster
	ipFamilyPolicy := corev1.IPFamilyPolicySingleStack

	return &corev1.Service{
		ObjectMeta: newResourceMeta(obj, meta),
		Spec: corev1.ServiceSpec{
			Selector:              runPodSelector(obj),
			Ports:                 ports,
			InternalTrafficPolicy: &internalTrafficPolicy,
			Type:                  corev1.ServiceTypeClusterIP,
			IPFamilies:            []corev1.IPFamily{corev1.IPv4Protocol},
			IPFamilyPolicy:        &ipFamilyPolicy,
			SessionAffinity:       corev1.ServiceAffinityNone,
		},
	}
}

type APIProxyServiceComponent struct {
	APIProxyPort func(obj *client.Object) int32
	Meta         *metadata.Provider
}

func (c APIProxyServiceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object
	svc := NewAPIProxyServiceComponent(&obj, c.APIProxyPort(&obj), c.Meta)

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
	Meta         *metadata.Provider
}

func NewAPIProxyNetworkPolicyComponent(obj *client.Object, port int32, meta *metadata.Provider) *networkingv1.NetworkPolicy {
	var apiProxyPort intstr.IntOrString
	if port != 0 {
		apiProxyPort = intstr.FromInt(int(port))
	} else {
		apiProxyPort = intstr.FromInt(defaultAPIProxyPort)
	}

	targetSelector := map[string]string{
		executionIDLabel:           (*obj).GetLabels()[executionIDLabel],
		projectIDLabel:             (*obj).GetLabels()[projectIDLabel],
		datasourceProxyClientLabel: "true",
	}

	tcpProto := corev1.ProtocolTCP

	ingressRules := []networkingv1.NetworkPolicyIngressRule{{
		From: []networkingv1.NetworkPolicyPeer{{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: targetSelector,
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
				MatchLabels: runPodSelector(obj),
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
	netPol := NewAPIProxyNetworkPolicyComponent(&obj, c.APIProxyPort(&obj), c.Meta)

	err := actions.CreateOrUpdateOwnedResource(ctx, obj, netPol)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

func (c APIProxyNetworkPolicyComponent) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

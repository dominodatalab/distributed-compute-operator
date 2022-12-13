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

type APIProxyServiceComponent struct {
	APIProxyPort func(obj *client.Object) int32
	Meta         *metadata.Provider
}

func executionID(obj client.Object) string {
	return obj.GetLabels()[executionIDLabel]
}

func projectID(obj client.Object) string {
	return obj.GetLabels()[projectIDLabel]
}

func instanceName(obj client.Object) string {
	return fmt.Sprintf("%s-%s", component, executionID(obj))
}

func runPodName(obj client.Object) string {
	return fmt.Sprintf("run-%s", executionID(obj))
}

func runPodSelector(obj client.Object) map[string]string {
	return map[string]string{
		metadata.ApplicationInstanceLabelKey: runPodName(obj),
	}
}

func (c APIProxyServiceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object

	internalTrafficPolicy := corev1.ServiceInternalTrafficPolicyCluster
	ipFamilyPolicy := corev1.IPFamilyPolicySingleStack

	apiProxyPort := c.APIProxyPort(&obj)
	if apiProxyPort == 0 {
		apiProxyPort = defaultAPIProxyPort
	}

	ports := []corev1.ServicePort{{
		Name:       apiProxyPortName,
		Port:       apiProxyPort,
		TargetPort: intstr.FromInt(targetServicePort),
		Protocol:   corev1.ProtocolTCP,
	}}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instanceName(obj),
			Namespace: obj.GetNamespace(),
			Labels:    c.Meta.StandardLabelsWithComponent(obj, component, nil),
		},
		Spec: corev1.ServiceSpec{
			//			ClusterIP:             corev1.ClusterIPNone,
			Selector:              runPodSelector(obj),
			Ports:                 ports,
			InternalTrafficPolicy: &internalTrafficPolicy,
			Type:                  corev1.ServiceTypeClusterIP,
			IPFamilies:            []corev1.IPFamily{corev1.IPv4Protocol},
			IPFamilyPolicy:        &ipFamilyPolicy,
			SessionAffinity:       corev1.ServiceAffinityNone,
		},
	}

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

func (c APIProxyNetworkPolicyComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object
	tcpProto := corev1.ProtocolTCP

	p := c.APIProxyPort(&obj)
	if p == 0 {
		p = defaultAPIProxyPort
	}
	apiProxyPort := intstr.FromInt(int(p))

	targetSelector := map[string]string{
		executionIDLabel:           executionID(obj),
		projectIDLabel:             projectID(obj),
		datasourceProxyClientLabel: "true",
	}

	ingressRules := []networkingv1.NetworkPolicyIngressRule{
		{
			From: []networkingv1.NetworkPolicyPeer{
				{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: targetSelector,
					},
				},
			},
			Ports: []networkingv1.NetworkPolicyPort{
				{
					Port:     &apiProxyPort,
					Protocol: &tcpProto,
				},
			},
		},
	}

	svc := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instanceName(obj),
			Namespace: obj.GetNamespace(),
			Labels:    c.Meta.StandardLabelsWithComponent(obj, component, nil),
		},
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

	err := actions.CreateOrUpdateOwnedResource(ctx, obj, svc)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

func (c APIProxyNetworkPolicyComponent) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

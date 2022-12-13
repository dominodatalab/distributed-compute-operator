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
	defaultApiProxyPort        = 8899
	targetServicePort          = 8899
	apiProxyPortName           = "api-proxy"
	component                  = "api-proxy"
	executionIdLabel           = "dominodatalab.com/execution-id"
	projectIdLabel             = "dominodatalab.com/project-id"
	datasourceProxyClientLabel = "datasource-proxy-client"
)

func APIProxyService(port int32, metaProvider *metadata.Provider) core.OwnedComponent {
	var apiProxyPort int32
	if port != 0 {
		apiProxyPort = port
	} else {
		apiProxyPort = defaultApiProxyPort
	}
	return &apiProxyServiceComponent{
		apiProxyPort,
		metaProvider,
	}
}

type apiProxyServiceComponent struct {
	apiProxyPort int32
	metaProvider *metadata.Provider
}

func executionId(obj client.Object) string {
	return obj.GetLabels()[executionIdLabel]
}

func projectId(obj client.Object) string {
	return obj.GetLabels()[projectIdLabel]
}

func instanceName(obj client.Object) string {
	return fmt.Sprintf("%s-%s", component, executionId(obj))
}

func runPodName(obj client.Object) string {
	return fmt.Sprintf("run-%s", executionId(obj))
}

func runPodSelector(obj client.Object) map[string]string {
	return map[string]string{
		metadata.ApplicationInstanceLabelKey: runPodName(obj),
	}
}

func (c apiProxyServiceComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object

	ports := []corev1.ServicePort{{
		Name:       apiProxyPortName,
		Port:       c.apiProxyPort,
		TargetPort: intstr.FromInt(targetServicePort),
		Protocol:   corev1.ProtocolTCP,
	}}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instanceName(obj),
			Namespace: obj.GetNamespace(),
			Labels:    c.metaProvider.StandardLabelsWithComponent(obj, component, nil),
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector:  runPodSelector(obj),
			Ports:     ports,
			//			InternalTrafficPolicy: &corev1.ServiceInternalTrafficPolicyCluster,
			Type:       corev1.ServiceTypeClusterIP,
			IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol},
			//IPFamilyPolicy: corev1.IPFamilyPolicySingleStack,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, obj, svc)
	if err != nil {
		err = fmt.Errorf("cannot reconcile service: %w", err)
	}

	return ctrl.Result{}, err
}

func (c apiProxyServiceComponent) Kind() client.Object {
	return &corev1.Service{}
}

type apiProxyNetworkPolicyComponent struct {
	apiProxyPort int32
	metaProvider *metadata.Provider
}

func (c *apiProxyNetworkPolicyComponent) ingressRules(ctx *core.Context) []networkingv1.NetworkPolicyIngressRule {
	tcpProto := corev1.ProtocolTCP
	apiProxyPort := intstr.FromInt(int(c.apiProxyPort))
	obj := ctx.Object
	targetSelector := map[string]string{
		executionIdLabel:           executionId(obj),
		projectIdLabel:             projectId(obj),
		datasourceProxyClientLabel: "true",
	}

	return []networkingv1.NetworkPolicyIngressRule{
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
}

func (c apiProxyNetworkPolicyComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	obj := ctx.Object
	tcpProto := corev1.ProtocolTCP
	apiProxyPort := intstr.FromInt(int(c.apiProxyPort))
	targetSelector := map[string]string{
		executionIdLabel:           executionId(obj),
		projectIdLabel:             projectId(obj),
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
			Labels:    c.metaProvider.StandardLabelsWithComponent(obj, component, nil),
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

func (c apiProxyNetworkPolicyComponent) Kind() client.Object {
	return &networkingv1.NetworkPolicy{}
}

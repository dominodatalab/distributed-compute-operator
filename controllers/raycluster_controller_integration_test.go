package controllers

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

var _ = Describe("RayCluster Controller", func() {

	Describe("Processing a new RayCluster resource", func() {

		It("should create a functional cluster", func() {
			ctx := context.Background()
			timeout := time.Second * 10

			psp := &policyv1beta1.PodSecurityPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "it",
				},
				Spec: policyv1beta1.PodSecurityPolicySpec{
					SELinux: policyv1beta1.SELinuxStrategyOptions{
						Rule: policyv1beta1.SELinuxStrategyRunAsAny,
					},
					RunAsUser: policyv1beta1.RunAsUserStrategyOptions{
						Rule: policyv1beta1.RunAsUserStrategyMustRunAsNonRoot,
					},
					SupplementalGroups: policyv1beta1.SupplementalGroupsStrategyOptions{
						Rule: policyv1beta1.SupplementalGroupsStrategyRunAsAny,
					},
					FSGroup: policyv1beta1.FSGroupStrategyOptions{
						Rule: policyv1beta1.FSGroupStrategyRunAsAny,
					},
				},
			}
			rayCluster := &dcv1alpha1.RayCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "it",
					Namespace: "default",
				},
				Spec: dcv1alpha1.RayClusterSpec{
					Image: &dcv1alpha1.OCIImageDefinition{
						Repository: "foo",
						Tag:        "bar",
					},
					Autoscaling: &dcv1alpha1.Autoscaling{
						MinReplicas:           pointer.Int32Ptr(1),
						MaxReplicas:           1,
						AverageCPUUtilization: pointer.Int32Ptr(50),
					},
					NetworkPolicy: dcv1alpha1.RayClusterNetworkPolicy{
						Enabled: pointer.BoolPtr(true),
					},
					Worker: dcv1alpha1.RayClusterWorker{
						Replicas: pointer.Int32Ptr(1),
					},
					Port:              6379,
					ClientServerPort:  10001,
					ObjectManagerPort: 2384,
					NodeManagerPort:   2385,
					DashboardPort:     8265,
					PodSecurityPolicy: psp.Name,
				},
			}

			Expect(k8sClient.Create(ctx, psp)).To(Succeed())
			Expect(k8sClient.Create(ctx, rayCluster)).To(Succeed())

			testcases := []struct {
				desc string
				name string
				obj  client.Object
			}{
				{"service account", "it-ray", &corev1.ServiceAccount{}},
				{"head service", "it-ray-head", &corev1.Service{}},
				{"cluster network policy", "it-ray-cluster", &networkingv1.NetworkPolicy{}},
				{"client network policy", "it-ray-client", &networkingv1.NetworkPolicy{}},
				{"dashboard network policy", "it-ray-dashboard", &networkingv1.NetworkPolicy{}},
				{"pod security policy role", "it-ray", &rbacv1.Role{}},
				{"pod security policy role binding", "it-ray", &rbacv1.RoleBinding{}},
				{"horizontal pod autoscaler", "it-ray", &autoscalingv2beta2.HorizontalPodAutoscaler{}},
				{"head stateful set", "it-ray-head", &appsv1.StatefulSet{}},
				{"worker stateful set", "it-ray-worker", &appsv1.StatefulSet{}},
			}
			for _, tc := range testcases {
				By(fmt.Sprintf("Creating a new %s", tc.desc))

				key := types.NamespacedName{
					Name:      tc.name,
					Namespace: "default",
				}
				obj := tc.obj

				Eventually(func() error {
					return k8sClient.Get(ctx, key, obj)
				}, timeout).Should(Succeed())
			}
		})
	})
})

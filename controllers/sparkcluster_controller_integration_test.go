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

var _ = Describe("SparkCluster Controller", func() {

	Describe("Processing a new SparkCluster resource", func() {

		It("should create a functional cluster", func() {
			ctx := context.Background()
			timeout := time.Second * 10

			psp := &policyv1beta1.PodSecurityPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: "it-spark",
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
			sparkCluster := &dcv1alpha1.SparkCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "it",
					Namespace: "default",
				},
				Spec: dcv1alpha1.SparkClusterSpec{
					Image: &dcv1alpha1.OCIImageDefinition{
						Repository: "foo",
						Tag:        "bar",
					},
					Autoscaling: &dcv1alpha1.Autoscaling{
						MinReplicas:           pointer.Int32Ptr(1),
						MaxReplicas:           1,
						AverageCPUUtilization: pointer.Int32Ptr(50),
					},
					NetworkPolicy: dcv1alpha1.SparkClusterNetworkPolicy{
						Enabled: pointer.BoolPtr(true),
					},
					Worker: dcv1alpha1.SparkClusterWorker{
						Replicas: pointer.Int32Ptr(1),
					},
					ClusterPort:       7077,
					DashboardPort:     8265,
					PodSecurityPolicy: psp.Name,
				},
			}

			Expect(k8sClient.Create(ctx, psp)).To(Succeed())
			Expect(k8sClient.Create(ctx, sparkCluster)).To(Succeed())

			testcases := []struct {
				desc string
				name string
				obj  client.Object
			}{
				{"service account", "it-spark", &corev1.ServiceAccount{}},
				{"head service", "it-spark-master", &corev1.Service{}},
				{"cluster network policy", "it-spark-cluster", &networkingv1.NetworkPolicy{}},
				{"client network policy", "it-spark-client", &networkingv1.NetworkPolicy{}},
				{"dashboard network policy", "it-spark-dashboard", &networkingv1.NetworkPolicy{}},
				{"pod security policy role", "it-spark", &rbacv1.Role{}},
				{"pod security policy role binding", "it-spark", &rbacv1.RoleBinding{}},
				{"horizontal pod autoscaler", "it-spark", &autoscalingv2beta2.HorizontalPodAutoscaler{}},
				{"head statefulset", "it-spark-master", &appsv1.StatefulSet{}},
				{"worker statefulset", "it-spark-worker", &appsv1.StatefulSet{}},
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
			cluster := dcv1alpha1.SparkCluster{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "default",
				Name:      "it",
			}, &cluster)).Should(Succeed())
			Expect(len(cluster.Finalizers)).Should(Equal(1))
			Expect(cluster.Finalizers[0]).Should(Equal(SparkFinalizerName))
		})
	})
})

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
			name := "functional"
			createAndBasicTest(ctx, name)
			timeout := time.Second * 10
			testcases := []struct {
				desc string
				name string
				obj  client.Object
			}{
				{"service account", name + "-spark", &corev1.ServiceAccount{}},
				{"head service", name + "-spark-master", &corev1.Service{}},
				{"worker headless service", name + "-spark-worker", &corev1.Service{}},
				{"cluster network policy", name + "-spark-cluster", &networkingv1.NetworkPolicy{}},
				{"client network policy", name + "-spark-client", &networkingv1.NetworkPolicy{}},
				{"dashboard network policy", name + "-spark-dashboard", &networkingv1.NetworkPolicy{}},
				{"driver network policy", name + "-spark-external", &networkingv1.NetworkPolicy{}},
				{"pod security policy role", name + "-spark", &rbacv1.Role{}},
				{"pod security policy role binding", name + "-spark", &rbacv1.RoleBinding{}},
				{"horizontal pod autoscaler", name + "-spark", &autoscalingv2beta2.HorizontalPodAutoscaler{}},
				{"head statefulset", name + "-spark-master", &appsv1.StatefulSet{}},
				{"worker statefulset", name + "-spark-worker", &appsv1.StatefulSet{}},
				{"framework configmap", name + "-framework-spark", &corev1.ConfigMap{}},
				{"keytab configmap", name + "-keytab-spark", &corev1.ConfigMap{}},
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

		It("should tear down gracefully", func() {
			ctx := context.Background()
			createAndBasicTest(ctx, "teardown")
			timeout := time.Second * 10
			cluster := dcv1alpha1.SparkCluster{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      "teardown"},
					&cluster)
			}, timeout).Should(Succeed())

			Expect(k8sClient.Delete(ctx, &cluster)).To(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      "teardown"},
					&cluster)
			}, timeout).ShouldNot(Succeed())

		})

		It("should add a finalizer", func() {
			ctx := context.Background()
			name := "finalizer"
			timeout := time.Second * 10
			createAndBasicTest(ctx, name)
			cluster := dcv1alpha1.SparkCluster{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      name,
				}, &cluster)
			}, timeout).Should(Succeed())
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      name,
				}, &cluster)).To(Succeed())
				return len(cluster.Finalizers) == 1 && cluster.Finalizers[0] == SparkFinalizerName
			}, timeout).Should(BeTrue())
		})

		It("should delete the finalizer", func() {
			ctx := context.Background()
			name := "finalizer-delete"
			timeout := time.Second * 10
			createAndBasicTest(ctx, name)
			cluster := dcv1alpha1.SparkCluster{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      name,
				}, &cluster)
			}, timeout).Should(Succeed())

			cluster.Finalizers = append(cluster.Finalizers, "test-finalizer")
			Expect(k8sClient.Update(ctx, &cluster)).To(Succeed())
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      name,
				}, &cluster)).To(Succeed())
				return len(cluster.Finalizers) == 2
			}, timeout).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, &cluster)).To(Succeed())

			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "default",
					Name:      name,
				}, &cluster)).To(Succeed())
				return len(cluster.Finalizers) == 1 && cluster.Finalizers[0] == "test-finalizer"
			}, timeout).Should(BeTrue())
		})
	})
})

func createAndBasicTest(ctx context.Context, name string) {
	psp := &policyv1beta1.PodSecurityPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
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
			Name:      name,
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
				Enabled:               pointer.BoolPtr(true),
				ExternalPodLabels:     map[string]string{"app.kubernetes.io/instance": "spark-driver"},
				ExternalPolicyEnabled: pointer.BoolPtr(true),
			},
			Master: dcv1alpha1.SparkClusterHead{
				SparkClusterNode: dcv1alpha1.SparkClusterNode{
					FrameworkConfig: &dcv1alpha1.FrameworkConfig{
						Path: "/opt/bitnami/spark/conf/spark-defaults.conf",
						Configs: map[string]string{
							"m1": "v1",
						},
					},
					KeyTabConfig: &dcv1alpha1.KeyTabConfig{
						Path:   "/etc/security/keytabs/kerberos.conf",
						KeyTab: []byte{'m', 'a', 's', 't', 'e', 'r'},
					},
				},
			},
			Worker: dcv1alpha1.SparkClusterWorker{
				SparkClusterNode: dcv1alpha1.SparkClusterNode{
					FrameworkConfig: &dcv1alpha1.FrameworkConfig{
						Path: "/opt/bitnami/spark/conf/spark-defaults.conf",
						Configs: map[string]string{
							"w1": "v1",
						},
					},
					KeyTabConfig: &dcv1alpha1.KeyTabConfig{
						Path:   "/etc/security/keytabs/kerberos.conf",
						KeyTab: []byte{'w', 'o', 'r', 'k', 'e', 'r'},
					},
				},
				Replicas: pointer.Int32Ptr(1),
			},
			ClusterPort:       7077,
			TCPMasterWebPort:  80,
			TCPWorkerWebPort:  8081,
			DashboardPort:     8265,
			PodSecurityPolicy: psp.Name,
		},
	}

	Expect(k8sClient.Create(ctx, psp)).To(Succeed())
	Expect(k8sClient.Create(ctx, sparkCluster)).To(Succeed())
}

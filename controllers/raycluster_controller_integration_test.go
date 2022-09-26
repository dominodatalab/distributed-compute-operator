package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

var _ = Describe("RayCluster Controller", func() {
	ctx := context.Background()
	timeout := time.Second * 10

	Describe("Processing a new RayCluster resource", func() {
		It("should create a functional cluster", func() {
			clusterKey, cluster := createCluster(ctx, "it")

			testcases := []struct {
				desc string
				name string
				obj  client.Object
			}{
				{"service account", "it-ray", &corev1.ServiceAccount{}},
				{"client service", "it-ray-client", &corev1.Service{}},
				{"headless head service", "it-ray-head", &corev1.Service{}},
				{"headless worker service", "it-ray-worker", &corev1.Service{}},
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
					Namespace: cluster.Namespace,
				}
				obj := tc.obj

				Eventually(func() error {
					return k8sClient.Get(ctx, key, obj)
				}, timeout).Should(Succeed())

				Expect(obj.GetOwnerReferences()).To(ConsistOf(metav1.OwnerReference{
					APIVersion:         dcv1alpha1.GroupVersion.Identifier(),
					Kind:               reflect.TypeOf(*cluster).Name(),
					Name:               cluster.Name,
					UID:                cluster.UID,
					Controller:         pointer.BoolPtr(true),
					BlockOwnerDeletion: pointer.BoolPtr(true),
				}))
			}

			By("Adding a finalizer")
			Eventually(func() []string {
				cluster := &dcv1alpha1.RayCluster{}
				if err := k8sClient.Get(ctx, clusterKey, cluster); err != nil {
					return nil
				}
				return cluster.Finalizers
			}, timeout).Should(ContainElement(DistributedComputeFinalizer))

			By("Updating the status with worker metadata")
			Eventually(func() dcv1alpha1.ClusterStatusConfig {
				cluster := &dcv1alpha1.RayCluster{}
				if err := k8sClient.Get(ctx, clusterKey, cluster); err != nil {
					return dcv1alpha1.ClusterStatusConfig{}
				}
				return cluster.Status
			}, timeout).Should(Equal(dcv1alpha1.ClusterStatusConfig{
				ClusterStatus:  dcv1alpha1.PendingStatus,
				Nodes:          nil,
				WorkerReplicas: 1,
				WorkerSelector: "app.kubernetes.io/component=worker,app.kubernetes.io/instance=it,app.kubernetes.io/name=ray",
			}))
		})
	})

	Describe("Updating an existing RayCluster resource", func() {
		It("should reconcile state changes", func() {
			clusterKey, cluster := createCluster(ctx, "update")

			By("promoting metadata onto all owned resources")
			Eventually(func() error {
				rc := &dcv1alpha1.RayCluster{}
				if err := k8sClient.Get(ctx, clusterKey, rc); err != nil {
					return err
				}

				rc.Spec.Image.Tag = "baz"
				return k8sClient.Update(ctx, rc)
			}, timeout).Should(Succeed())

			By("deleting the horizontal pod autoscaler when disabled")
			key := types.NamespacedName{Name: "update-ray", Namespace: cluster.Namespace}
			autoscaler := &autoscalingv2beta2.HorizontalPodAutoscaler{}

			Eventually(func() error {
				return k8sClient.Get(ctx, key, autoscaler)
			}, timeout).Should(Succeed())

			Eventually(func() error {
				rc := &dcv1alpha1.RayCluster{}
				if err := k8sClient.Get(ctx, clusterKey, rc); err != nil {
					return err
				}

				rc.Spec.Autoscaling = nil
				return k8sClient.Update(ctx, rc)
			}, timeout).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, key, autoscaler)
			}, timeout).ShouldNot(Succeed())

			By("deleting network policies when disabled")

			By("deleting pod security policy rbac resources when disabled")
		})
	})

	Describe("Deleting a RayCluster resource", func() {
		It("should delete external persistent volume claims created by stateful sets", func() {
			clusterKey, cluster := createCluster(ctx, "delete")

			pvc := &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "whatever",
					Namespace: cluster.Namespace,
					Labels: map[string]string{
						"apps.kubernetes.io/name":     "ray",
						"apps.kubernetes.io/instance": cluster.Name,
					},
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse("1Gi"),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, pvc)).To(Succeed())
			Expect(k8sClient.Delete(ctx, cluster)).To(Succeed())

			Eventually(func() error {
				return k8sClient.Get(context.Background(), clusterKey, pvc)
			}, timeout).ShouldNot(Succeed())
		})
	})
})

func createCluster(ctx context.Context, name string) (client.ObjectKey, *dcv1alpha1.RayCluster) {
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
	cluster := &dcv1alpha1.RayCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: dcv1alpha1.RayClusterSpec{
			ScalableClusterConfig: dcv1alpha1.ScalableClusterConfig{
				ClusterConfig: dcv1alpha1.ClusterConfig{
					Image: &dcv1alpha1.OCIImageDefinition{
						Repository: "foo",
						Tag:        "bar",
					},
					NetworkPolicy: dcv1alpha1.NetworkPolicyConfig{
						Enabled: pointer.BoolPtr(true),
					},
					PodSecurityPolicy: psp.Name,
				},
				Autoscaling: &dcv1alpha1.Autoscaling{
					MinReplicas:              pointer.Int32Ptr(1),
					MaxReplicas:              1,
					AverageCPUUtilization:    pointer.Int32Ptr(50),
					AverageMemoryUtilization: pointer.Int32Ptr(50),
				},
			},
			Worker: dcv1alpha1.RayClusterWorker{
				Replicas: pointer.Int32Ptr(1),
			},
			Port:              6379,
			ClientServerPort:  10001,
			ObjectManagerPort: 2384,
			NodeManagerPort:   2385,
			GCSServerPort:     2386,
			WorkerPorts:       []int32{11000, 11001},
			DashboardPort:     8265,
		},
	}
	clusterKey := client.ObjectKeyFromObject(cluster)

	Expect(k8sClient.Create(ctx, psp)).To(Succeed())
	Expect(k8sClient.Create(ctx, cluster)).To(Succeed())

	return clusterKey, cluster
}

package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

var _ = Describe("RayCluster controller", func() {
	timeout := time.Second * 10
	interval := time.Millisecond * 250

	Context("New RayCluster resource", func() {
		It("should create a functional cluster", func() {
			ctx := context.Background()
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
					Port:              6379,
					ClientServerPort:  10001,
					ObjectManagerPort: 2384,
					NodeManagerPort:   2385,
					DashboardPort:     8265,
					Worker: dcv1alpha1.RayClusterWorker{
						Replicas: pointer.Int32Ptr(1),
					},
				},
			}

			Expect(k8sClient.Create(ctx, rayCluster)).To(Succeed())

			By("Creating a service account")
			key := types.NamespacedName{
				Name:      "it-ray",
				Namespace: "default",
			}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, key, &corev1.ServiceAccount{})
				return err == nil
			}, timeout, interval).Should(BeTrue())

			By("Creating a head service")
			key = types.NamespacedName{
				Name:      "it-ray-head",
				Namespace: "default",
			}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, key, &corev1.Service{})
				return err == nil
			}, timeout, interval).Should(BeTrue())

			By("Creating network policies")

			By("Creating pod security policies")

			By("Creating a horizontal pod autoscaler")

			By("Creating a head deployment")

			By("Creating a worker deployment")
		})
	})
})

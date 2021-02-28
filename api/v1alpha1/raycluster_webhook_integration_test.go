package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func fixture(nsName string) *RayCluster {
	return &RayCluster{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-",
			Namespace:    nsName,
		},
		Spec: RayClusterSpec{
			Worker: RayClusterWorker{
				Replicas: 1,
			},
		},
	}
}

var _ = Describe("RayCluster", func() {
	Describe("Validation", func() {
		var testNS *v1.Namespace

		BeforeEach(func() {
			testNS = &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
				},
			}
			Expect(k8sClient.Create(ctx, testNS)).To(Succeed())
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, testNS)).To(Succeed())
		})

		It("passes when object is valid", func() {
			rc := fixture(testNS.Name)
			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
		})

		It("requires a positive worker replica count", func() {
			rc := fixture(testNS.Name)
			rc.Spec.Worker.Replicas = 0

			Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
		})

		It("requires a minimum of 75MB for object store memory", func() {
			rc := fixture(testNS.Name)
			rc.Spec.ObjectStoreMemoryBytes = pointer.Int64Ptr(74 * 1024 * 1024)

			Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
		})

		table.DescribeTable("(networking ports)",
			func(portSetter func(*RayCluster, int32)) {
				rc := fixture(testNS.Name)

				portSetter(rc, -1)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())

				portSetter(rc, 65354)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			},
			table.Entry("rejects an invalid port",
				func(rc *RayCluster, val int32) { rc.Spec.Port = val },
			),
			table.Entry("rejects invalid redis shard ports",
				func(rc *RayCluster, val int32) { rc.Spec.RedisShardPorts = append(rc.Spec.RedisShardPorts, val) },
			),
			table.Entry("rejects an invalid client server port",
				func(rc *RayCluster, val int32) { rc.Spec.ClientServerPort = val },
			),
			table.Entry("rejects an invalid object manager port",
				func(rc *RayCluster, val int32) { rc.Spec.ObjectManagerPort = val },
			),
			table.Entry("rejects an invalid node manager port",
				func(rc *RayCluster, val int32) { rc.Spec.NodeManagerPort = val },
			),
			table.Entry("rejects an invalid dashboard port",
				func(rc *RayCluster, val int32) { rc.Spec.DashboardPort = val },
			),
		)

		Context("With autoscaling enabled", func() {
			clusterWithAS := func() *RayCluster {
				rc := fixture(testNS.Name)
				rc.Spec.Autoscaling = &Autoscaling{
					MinReplicas:                         pointer.Int32Ptr(1),
					MaxReplicas:                         1,
					AverageUtilization:                  50,
					ScaleDownStabilizationWindowSeconds: pointer.Int32Ptr(50),
				}

				return rc
			}

			It("passes when valid", func() {
				rc := clusterWithAS()
				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			})

			It("does not require min replicas", func() {
				rc := clusterWithAS()
				rc.Spec.Autoscaling.MinReplicas = nil

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			})

			It("requires min replicas to be > 0 when provided", func() {
				rc := clusterWithAS()
				rc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(0)

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("requires max replicas to be > 0", func() {
				rc := clusterWithAS()
				rc.Spec.Autoscaling.MinReplicas = nil
				rc.Spec.Autoscaling.MaxReplicas = 0

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("requires max replicas to be > min replicas", func() {
				rc := clusterWithAS()
				rc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(2)
				rc.Spec.Autoscaling.MaxReplicas = 1

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("requires average utilization to be > 0 and <= 100", func() {
				rc := clusterWithAS()

				rc.Spec.Autoscaling.AverageUtilization = 0
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())

				rc.Spec.Autoscaling.AverageUtilization = 101
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("does not require scale down stabilization", func() {
				rc := clusterWithAS()
				rc.Spec.Autoscaling.ScaleDownStabilizationWindowSeconds = nil

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			})

			It("requires scale down stabilization to be >= 0 when provided", func() {
				rc := clusterWithAS()
				rc.Spec.Autoscaling.ScaleDownStabilizationWindowSeconds = pointer.Int32Ptr(-1)

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})
		})
	})
})

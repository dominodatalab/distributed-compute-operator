package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func rayFixture(nsName string) *RayCluster {
	return &RayCluster{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-",
			Namespace:    nsName,
		},
	}
}

var _ = Describe("RayCluster", func() {
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

	Describe("Defaulting", func() {
		It("sets expected values on an empty object", func() {
			rc := rayFixture(testNS.Name)
			Expect(k8sClient.Create(ctx, rc)).To(Succeed())

			Expect(rc.Spec.Port).To(
				BeNumerically("==", 6379),
				"port should equal 6379",
			)
			Expect(rc.Spec.RedisShardPorts).To(
				Equal([]int32{6380, 6381}),
				"redis shard ports should equal [6380, 6381]",
			)
			Expect(rc.Spec.ClientServerPort).To(
				BeNumerically("==", 10001),
				"client server port should equal 10001",
			)
			Expect(rc.Spec.ObjectManagerPort).To(
				BeNumerically("==", 2384),
				"object manager port should equal 2384",
			)
			Expect(rc.Spec.NodeManagerPort).To(
				BeNumerically("==", 2385),
				"node manager port should equal 2385",
			)
			Expect(rc.Spec.GCSServerPort).To(
				BeNumerically("==", 2386),
				"gcs server port should equal 2386",
			)
			Expect(rc.Spec.DashboardPort).To(
				BeNumerically("==", 8265),
				"dashboard port should equal 8265",
			)
			Expect(rc.Spec.EnableDashboard).To(
				PointTo(Equal(true)),
				"enable dashboard should point to true",
			)
			Expect(rc.Spec.NetworkPolicy.Enabled).To(
				PointTo(Equal(true)),
				"enable network policy should point to true",
			)
			Expect(rc.Spec.NetworkPolicy.ClientServerLabels).To(
				Equal(map[string]string{"ray-client": "true"}),
				`network policy client labels should equal [{"ray-client": "true"}]`,
			)
			Expect(rc.Spec.NetworkPolicy.DashboardLabels).To(
				Equal(map[string]string{"ray-client": "true"}),
				`network policy dashboard labels should equal [{"ray-client": "true"}]`,
			)
			Expect(rc.Spec.Worker.Replicas).To(
				PointTo(BeNumerically("==", 1)),
				"worker replicas should point to 1",
			)
			Expect(rc.Spec.Image).To(
				Equal(&OCIImageDefinition{Repository: "rayproject/ray", Tag: "1.3.0-cpu"}),
				`image reference should equal "rayproject/ray:1.3.0-cpu"`,
			)
		})

		It("does not set the port when present", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.Port = 3000

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.Port).To(BeNumerically("==", 3000))
		})

		It("does not set redis shard ports when present", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.RedisShardPorts = []int32{5000}

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.RedisShardPorts).To(Equal([]int32{5000}))
		})

		It("does not set the client server port when present", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.ClientServerPort = 12000

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.ClientServerPort).To(BeNumerically("==", 12000))
		})

		It("does not set the object manager port when present", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.ObjectManagerPort = 4832

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.ObjectManagerPort).To(BeNumerically("==", 4832))
		})

		It("does not set the node manager port when present", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.NodeManagerPort = 5832

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.NodeManagerPort).To(BeNumerically("==", 5832))
		})

		It("does not set the dashboard port when present", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.DashboardPort = 5555

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.DashboardPort).To(BeNumerically("==", 5555))
		})

		It("does not enable the dashboard when false", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.EnableDashboard = pointer.BoolPtr(false)

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.EnableDashboard).To(PointTo(Equal(false)))
		})

		Context("Network policies", func() {
			It("are not enabled when false", func() {
				rc := rayFixture(testNS.Name)
				rc.Spec.NetworkPolicy.Enabled = pointer.BoolPtr(false)

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
				Expect(rc.Spec.NetworkPolicy.Enabled).To(PointTo(Equal(false)))
			})

			It("use provided client server labels", func() {
				rc := rayFixture(testNS.Name)

				expected := map[string]string{"server-client": "true"}
				rc.Spec.NetworkPolicy.ClientServerLabels = expected

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
				Expect(rc.Spec.NetworkPolicy.ClientServerLabels).To(Equal(expected))
			})

			It("use provided dashboard labels", func() {
				rc := rayFixture(testNS.Name)

				expected := map[string]string{"dashboard-client": "true"}
				rc.Spec.NetworkPolicy.DashboardLabels = expected

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
				Expect(rc.Spec.NetworkPolicy.DashboardLabels).To(Equal(expected))
			})
		})
	})

	Describe("Validation", func() {
		It("passes when object is valid", func() {
			rc := rayFixture(testNS.Name)
			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
		})

		It("requires a positive worker replica count", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.Worker.Replicas = pointer.Int32Ptr(-1)

			Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
		})

		It("requires a minimum of 75MB for object store memory", func() {
			rc := rayFixture(testNS.Name)
			rc.Spec.ObjectStoreMemoryBytes = pointer.Int64Ptr(74 * 1024 * 1024)

			Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
		})

		DescribeTable("(networking ports)",
			func(portSetter func(*RayCluster, int32)) {
				rc := rayFixture(testNS.Name)

				portSetter(rc, 1023)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())

				portSetter(rc, 65536)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			},
			Entry("rejects an invalid port",
				func(rc *RayCluster, val int32) { rc.Spec.Port = val },
			),
			Entry("rejects invalid redis shard ports",
				func(rc *RayCluster, val int32) { rc.Spec.RedisShardPorts = append(rc.Spec.RedisShardPorts, val) },
			),
			Entry("rejects an invalid client server port",
				func(rc *RayCluster, val int32) { rc.Spec.ClientServerPort = val },
			),
			Entry("rejects an invalid object manager port",
				func(rc *RayCluster, val int32) { rc.Spec.ObjectManagerPort = val },
			),
			Entry("rejects an invalid node manager port",
				func(rc *RayCluster, val int32) { rc.Spec.NodeManagerPort = val },
			),
			Entry("rejects an invalid gcs server port",
				func(rc *RayCluster, val int32) { rc.Spec.GCSServerPort = val },
			),
			Entry("rejects invalid worker ports",
				func(rc *RayCluster, val int32) { rc.Spec.WorkerPorts = append(rc.Spec.WorkerPorts, val) },
			),
			Entry("rejects an invalid dashboard port",
				func(rc *RayCluster, val int32) { rc.Spec.DashboardPort = val },
			),
		)

		Context("With a provided image", func() {
			It("requires a non-blank image registry", func() {
				rc := rayFixture(testNS.Name)
				rc.Spec.Image = &OCIImageDefinition{Tag: "test-tag"}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("requires a non-blank image tag", func() {
				rc := rayFixture(testNS.Name)
				rc.Spec.Image = &OCIImageDefinition{Repository: "test-repo"}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})
		})

		Context("With autoscaling enabled", func() {
			clusterWithAutoscaling := func() *RayCluster {
				rc := rayFixture(testNS.Name)
				rc.Spec.Autoscaling = &Autoscaling{
					MaxReplicas: 1,
				}
				rc.Spec.Worker.Resources.Requests = v1.ResourceList{
					v1.ResourceCPU: resource.MustParse("100m"),
				}

				return rc
			}

			It("passes when valid", func() {
				rc := clusterWithAutoscaling()
				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			})

			It("requires min replicas to be > 0 when provided", func() {
				rc := clusterWithAutoscaling()

				rc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(0)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())

				rc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(1)
				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			})

			It("requires max replicas to be > 0", func() {
				rc := clusterWithAutoscaling()
				rc.Spec.Autoscaling.MaxReplicas = 0

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("requires max replicas to be > min replicas", func() {
				rc := clusterWithAutoscaling()

				rc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(2)
				rc.Spec.Autoscaling.MaxReplicas = 1
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())

				rc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(1)
				rc.Spec.Autoscaling.MaxReplicas = 2
				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			})

			It("requires average utilization to be > 0", func() {
				rc := clusterWithAutoscaling()

				rc.Spec.Autoscaling.AverageCPUUtilization = pointer.Int32Ptr(0)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())

				rc.Spec.Autoscaling.AverageCPUUtilization = pointer.Int32Ptr(75)
				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			})

			It("requires scale down stabilization to be >= 0 when provided", func() {
				rc := clusterWithAutoscaling()

				rc.Spec.Autoscaling.ScaleDownStabilizationWindowSeconds = pointer.Int32Ptr(-1)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())

				rc.Spec.Autoscaling.ScaleDownStabilizationWindowSeconds = pointer.Int32Ptr(0)
				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			})

			It("requires cpu resource requests for worker", func() {
				rc := clusterWithAutoscaling()
				rc.Spec.Worker.Resources.Requests = nil

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})
		})

		DescribeTable("With mutal tls mode set",
			func(smode string, expectErr bool) {
				rc := rayFixture(testNS.Name)
				rc.Spec.MutualTLSMode = smode

				if expectErr {
					Expect(k8sClient.Create(ctx, rc)).To(HaveOccurred())
				} else {
					Expect(k8sClient.Create(ctx, rc)).NotTo(HaveOccurred())
				}
			},
			Entry("empty string is valid", "", false),
			Entry("UNSET is valid", "UNSET", false),
			Entry("DISABLE is valid", "DISABLE", false),
			Entry("PERMISSIVE is valid", "PERMISSIVE", false),
			Entry("STRICT is valid", "STRICT", false),
			Entry("GARBAGE is not valid", "GARBAGE", true),
		)
	})
})

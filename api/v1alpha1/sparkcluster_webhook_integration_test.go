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

func sparkFixture(nsName string) *SparkCluster {
	return &SparkCluster{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-",
			Namespace:    nsName,
		},
		// need this value to be preset to pass webhook tests
		Spec: SparkClusterSpec{
			Worker: SparkClusterWorker{WorkerMemoryLimit: "4505m"},
		},
	}
}

var _ = Describe("SparkCluster", func() {
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
			sc := sparkFixture(testNS.Name)
			Expect(k8sClient.Create(ctx, sc)).To(Succeed())

			Expect(sc.Spec.ClusterPort).To(
				BeNumerically("==", 7077),
				"port should equal 7077",
			)
			Expect(sc.Spec.TCPWorkerWebPort).To(
				BeNumerically("==", 8081),
				"worker web port should equal 8081",
			)
			Expect(sc.Spec.TCPMasterWebPort).To(
				BeNumerically("==", 8080),
				"master web port should equal 8080",
			)
			Expect(sc.Spec.DashboardPort).To(
				BeNumerically("==", 8080),
				"dashboard port should equal 8080",
			)
			Expect(sc.Spec.DashboardServicePort).To(
				BeNumerically("==", 80),
				"port should equal 80",
			)
			Expect(sc.Spec.EnableDashboard).To(
				PointTo(Equal(true)),
				"enable dashboard should point to true",
			)
			Expect(sc.Spec.Driver.DriverUIPort).To(
				BeNumerically("==", 4040),
				"driver ui port should equal 4040",
			)
			Expect(sc.Spec.Driver.DriverPort).To(
				BeNumerically("==", 4041),
				"driver port should equal 4041",
			)
			Expect(sc.Spec.Driver.DriverBlockManagerPort).To(
				BeNumerically("==", 4042),
				"driver block manager port should equal 4042",
			)
			Expect(sc.Spec.Driver.DriverPortName).To(
				Equal("spark-driver-port"),
				"driver port name should equal spark-driver-port",
			)
			Expect(sc.Spec.Driver.DriverUIPortName).To(
				Equal("spark-ui-port"),
				"driver port name should equal spark-ui-port",
			)
			Expect(sc.Spec.Driver.DriverBlockManagerPortName).To(
				Equal("spark-block-manager-port"),
				"driver block manager port name should equal spark-block-manager-port",
			)
			Expect(sc.Spec.NetworkPolicy.Enabled).To(
				PointTo(Equal(true)),
				"enable network policy should point to true",
			)
			Expect(sc.Spec.NetworkPolicy.ClientServerLabels).To(
				Equal(map[string]string{"spark-client": "true"}),
				`network policy client labels should equal [{"spark-client": "true"}]`,
			)
			Expect(sc.Spec.NetworkPolicy.DashboardLabels).To(
				Equal(map[string]string{"spark-client": "true"}),
				`network policy dashboard labels should equal [{"spark-client": "true"}]`,
			)
			Expect(sc.Spec.Worker.Replicas).To(
				PointTo(BeNumerically("==", 1)),
				"worker replicas should point to 1",
			)
			Expect(sc.Spec.Image).To(
				Equal(&OCIImageDefinition{Repository: "bitnami/spark", Tag: "3.0.2-debian-10-r0"}),
				`image reference should equal "bitnami/spark:3.0.2-debian-10-r0"`,
			)
		})

		It("does not set the cluster port when present", func() {
			sc := sparkFixture(testNS.Name)
			sc.Spec.ClusterPort = 3000

			Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			Expect(sc.Spec.ClusterPort).To(BeNumerically("==", 3000))
		})

		It("does not set the worker web port when present", func() {
			sc := sparkFixture(testNS.Name)
			sc.Spec.TCPWorkerWebPort = 3000

			Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			Expect(sc.Spec.TCPWorkerWebPort).To(BeNumerically("==", 3000))
		})

		It("does not set the master web port when present", func() {
			sc := sparkFixture(testNS.Name)
			sc.Spec.TCPMasterWebPort = 3000

			Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			Expect(sc.Spec.TCPMasterWebPort).To(BeNumerically("==", 3000))
		})

		It("does not set the dashboard port when present", func() {
			sc := sparkFixture(testNS.Name)
			sc.Spec.DashboardPort = 5555

			Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			Expect(sc.Spec.DashboardPort).To(BeNumerically("==", 5555))
		})

		It("does not enable the dashboard when false", func() {
			sc := sparkFixture(testNS.Name)
			sc.Spec.EnableDashboard = pointer.BoolPtr(false)

			Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			Expect(sc.Spec.EnableDashboard).To(PointTo(Equal(false)))
		})

		Context("Network policies", func() {
			It("are not enabled when false", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.NetworkPolicy.Enabled = pointer.BoolPtr(false)

				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
				Expect(sc.Spec.NetworkPolicy.Enabled).To(PointTo(Equal(false)))
			})

			It("use provided client server labels", func() {
				sc := sparkFixture(testNS.Name)

				expected := map[string]string{"server-client": "true"}
				sc.Spec.NetworkPolicy.ClientServerLabels = expected

				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
				Expect(sc.Spec.NetworkPolicy.ClientServerLabels).To(Equal(expected))
			})

			It("use provided dashboard labels", func() {
				sc := sparkFixture(testNS.Name)

				expected := map[string]string{"dashboard-client": "true"}
				sc.Spec.NetworkPolicy.DashboardLabels = expected

				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
				Expect(sc.Spec.NetworkPolicy.DashboardLabels).To(Equal(expected))
			})

			It("use provided cluster labels", func() {
				sc := sparkFixture(testNS.Name)

				expected := map[string]string{"instance": "spark-driver"}
				sc.Spec.NetworkPolicy.ExternalPodLabels = expected

				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
				Expect(sc.Spec.NetworkPolicy.ExternalPodLabels).To(Equal(expected))
			})
		})
	})

	Describe("Validation", func() {
		It("passes when object is valid", func() {
			sc := sparkFixture(testNS.Name)
			Expect(k8sClient.Create(ctx, sc)).To(Succeed())
		})

		It("requires a positive worker replica count", func() {
			sc := sparkFixture(testNS.Name)
			sc.Spec.Worker.Replicas = pointer.Int32Ptr(-10)

			Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
		})

		DescribeTable("(networking ports)",
			func(portSetter func(*SparkCluster, int32)) {
				sc := sparkFixture(testNS.Name)

				portSetter(sc, 79)
				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())

				portSetter(sc, 65536)
				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			},
			Entry("rejects an invalid cluster port",
				func(sc *SparkCluster, val int32) { sc.Spec.ClusterPort = val },
			),
			Entry("rejects an invalid worker web port",
				func(sc *SparkCluster, val int32) { sc.Spec.TCPWorkerWebPort = val },
			),
			Entry("rejects an invalid master web port",
				func(sc *SparkCluster, val int32) { sc.Spec.TCPMasterWebPort = val },
			),
			Entry("rejects an invalid dashboard port",
				func(sc *SparkCluster, val int32) { sc.Spec.DashboardPort = val },
			),
			Entry("rejects an invalid dashboard service port",
				func(sc *SparkCluster, val int32) { sc.Spec.DashboardServicePort = val },
			),
		)

		DescribeTable("With mutual tls mode set",
			func(smode string, expectErr bool) {
				sc := sparkFixture(testNS.Name)
				sc.Spec.MutualTLSMode = smode

				if expectErr {
					Expect(k8sClient.Create(ctx, sc)).To(HaveOccurred())
				} else {
					Expect(k8sClient.Create(ctx, sc)).NotTo(HaveOccurred())
				}
			},
			Entry("empty string is valid", "", false),
			Entry("UNSET is valid", "UNSET", false),
			Entry("DISABLE is valid", "DISABLE", false),
			Entry("PERMISSIVE is valid", "PERMISSIVE", false),
			Entry("STRICT is valid", "STRICT", false),
			Entry("GARBAGE is not valid", "GARBAGE", true),
		)

		Context("With a provided image", func() {
			It("requires a non-blank image registry", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.Image = &OCIImageDefinition{Tag: "test-tag"}

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})

			It("requires a non-blank image tag", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.Image = &OCIImageDefinition{Repository: "test-repo"}

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})
		})

		Context("With autoscaling enabled", func() {
			clusterWithAutoscaling := func() *SparkCluster {
				sc := sparkFixture(testNS.Name)
				sc.Spec.Autoscaling = &Autoscaling{
					MaxReplicas: 1,
				}
				sc.Spec.Worker.Resources.Requests = v1.ResourceList{
					v1.ResourceCPU: resource.MustParse("100m"),
				}

				return sc
			}

			It("passes when valid", func() {
				sc := clusterWithAutoscaling()
				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			})

			It("requires min replicas to be > 0 when provided", func() {
				sc := clusterWithAutoscaling()

				sc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(0)
				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())

				sc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(1)
				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			})

			It("requires max replicas to be > 0", func() {
				sc := clusterWithAutoscaling()
				sc.Spec.Autoscaling.MaxReplicas = 0

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})

			It("requires max replicas to be > min replicas", func() {
				sc := clusterWithAutoscaling()

				sc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(2)
				sc.Spec.Autoscaling.MaxReplicas = 1
				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())

				sc.Spec.Autoscaling.MinReplicas = pointer.Int32Ptr(1)
				sc.Spec.Autoscaling.MaxReplicas = 2
				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			})

			It("requires average utilization to be > 0", func() {
				sc := clusterWithAutoscaling()

				sc.Spec.Autoscaling.AverageCPUUtilization = pointer.Int32Ptr(0)
				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())

				sc.Spec.Autoscaling.AverageCPUUtilization = pointer.Int32Ptr(75)
				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			})

			It("requires scale down stabilization to be >= 0 when provided", func() {
				sc := clusterWithAutoscaling()

				sc.Spec.Autoscaling.ScaleDownStabilizationWindowSeconds = pointer.Int32Ptr(-1)
				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())

				sc.Spec.Autoscaling.ScaleDownStabilizationWindowSeconds = pointer.Int32Ptr(0)
				Expect(k8sClient.Create(ctx, sc)).To(Succeed())
			})

			It("requires cpu resource requests for worker", func() {
				sc := clusterWithAutoscaling()
				sc.Spec.Worker.Resources.Requests = nil

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})
		})

		Context("worker memory limit", func() {
			It("rejects an empty limit", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.Worker.WorkerMemoryLimit = ""

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})
			It("rejects an invalid limit", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.Worker.WorkerMemoryLimit = "blah"

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})
		})

		Context("framework configs", func() {
			It("rejects a config with no data set", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.Worker.FrameworkConfig = &FrameworkConfig{
					Configs: nil,
				}

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})
		})

		Context("keytab configs", func() {
			It("rejects a config with only path set", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.KerberosKeytab = &KerberosKeytabConfig{
					MountPath: "test/path/",
					Contents:  nil,
				}

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})

			It("rejects a config with only data set", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.KerberosKeytab = &KerberosKeytabConfig{
					MountPath: "",
					Contents:  []byte{'c', 'o', 'n', 'f', 'i', 'g'},
				}

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})
		})

		Context("external network policies", func() {
			It("rejects when policy is enabled but no values are set", func() {
				sc := sparkFixture(testNS.Name)
				sc.Spec.NetworkPolicy.ExternalPolicyEnabled = pointer.BoolPtr(true)
				sc.Spec.NetworkPolicy.ExternalPodLabels = map[string]string{}

				Expect(k8sClient.Create(ctx, sc)).ToNot(Succeed())
			})
		})
	})
})

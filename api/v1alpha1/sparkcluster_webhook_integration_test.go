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
			rc := sparkFixture(testNS.Name)
			Expect(k8sClient.Create(ctx, rc)).To(Succeed())

			Expect(rc.Spec.ClusterPort).To(
				BeNumerically("==", 7077),
				"port should equal 7077",
			)
			Expect(rc.Spec.TCPWorkerWebPort).To(
				BeNumerically("==", 8081),
				"worker web port should equal 8081",
			)
			Expect(rc.Spec.TCPMasterWebPort).To(
				BeNumerically("==", 80),
				"master web port should equal 80",
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
				Equal(map[string]string{"spark-client": "true"}),
				`network policy client labels should equal [{"spark-client": "true"}]`,
			)
			Expect(rc.Spec.NetworkPolicy.DashboardLabels).To(
				Equal(map[string]string{"spark-client": "true"}),
				`network policy dashboard labels should equal [{"spark-client": "true"}]`,
			)
			Expect(rc.Spec.Worker.Replicas).To(
				PointTo(BeNumerically("==", 1)),
				"worker replicas should point to 1",
			)
			Expect(rc.Spec.Image).To(
				Equal(&OCIImageDefinition{Repository: "bitnami/spark", Tag: "3.0.2-debian-10-r0"}),
				`image reference should equal "bitnami/spark:3.0.2-debian-10-r0"`,
			)
		})

		It("does not set the cluster port when present", func() {
			rc := sparkFixture(testNS.Name)
			rc.Spec.ClusterPort = 3000

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.ClusterPort).To(BeNumerically("==", 3000))
		})

		It("does not set the worker web port when present", func() {
			rc := sparkFixture(testNS.Name)
			rc.Spec.TCPWorkerWebPort = 3000

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.TCPWorkerWebPort).To(BeNumerically("==", 3000))
		})

		It("does not set the master web port when present", func() {
			rc := sparkFixture(testNS.Name)
			rc.Spec.TCPMasterWebPort = 3000

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.TCPMasterWebPort).To(BeNumerically("==", 3000))
		})

		It("does not set the dashboard port when present", func() {
			rc := sparkFixture(testNS.Name)
			rc.Spec.DashboardPort = 5555

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.DashboardPort).To(BeNumerically("==", 5555))
		})

		It("does not enable the dashboard when false", func() {
			rc := sparkFixture(testNS.Name)
			rc.Spec.EnableDashboard = pointer.BoolPtr(false)

			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			Expect(rc.Spec.EnableDashboard).To(PointTo(Equal(false)))
		})

		Context("Network policies", func() {
			It("are not enabled when false", func() {
				rc := sparkFixture(testNS.Name)
				rc.Spec.NetworkPolicy.Enabled = pointer.BoolPtr(false)

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
				Expect(rc.Spec.NetworkPolicy.Enabled).To(PointTo(Equal(false)))
			})

			It("use provided client server labels", func() {
				rc := sparkFixture(testNS.Name)

				expected := map[string]string{"server-client": "true"}
				rc.Spec.NetworkPolicy.ClientServerLabels = expected

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
				Expect(rc.Spec.NetworkPolicy.ClientServerLabels).To(Equal(expected))
			})

			It("use provided dashboard labels", func() {
				rc := sparkFixture(testNS.Name)

				expected := map[string]string{"dashboard-client": "true"}
				rc.Spec.NetworkPolicy.DashboardLabels = expected

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
				Expect(rc.Spec.NetworkPolicy.DashboardLabels).To(Equal(expected))
			})

			It("use provided cluster labels", func() {
				rc := sparkFixture(testNS.Name)

				expected := map[string]string{"instance": "spark-driver"}
				rc.Spec.NetworkPolicy.ExternalPodLabels = expected

				Expect(k8sClient.Create(ctx, rc)).To(Succeed())
				Expect(rc.Spec.NetworkPolicy.ExternalPodLabels).To(Equal(expected))
			})
		})

		Context("Annotations", func() {
			//It("add istio annotation to provided annotations", func() {
			//	rc := sparkFixture(testNS.Name)
			//	provided := map[string]string{"annotation": "test"}
			//	rc.Spec.Master.Annotations = provided
			//	rc.Spec.Worker.Annotations = provided
			//
			//	expected := map[string]string{
			//		"annotation":              "test",
			//		"sidecar.istio.io/inject": "false",
			//	}
			//	Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			//	Expect(rc.Spec.Master.Annotations).To(Equal(expected))
			//	Expect(rc.Spec.Worker.Annotations).To(Equal(expected))
			//})

			//It("override provided istio annotation", func() {
			//	rc := sparkFixture(testNS.Name)
			//	provided := map[string]string{"sidecar.istio.io/inject": "true"}
			//	rc.Spec.Master.Annotations = provided
			//	rc.Spec.Worker.Annotations = provided
			//
			//	expected := map[string]string{"sidecar.istio.io/inject": "false"}
			//	Expect(k8sClient.Create(ctx, rc)).To(Succeed())
			//	Expect(rc.Spec.Master.Annotations).To(Equal(expected))
			//	Expect(rc.Spec.Worker.Annotations).To(Equal(expected))
			//})
		})
	})

	Describe("Validation", func() {
		It("passes when object is valid", func() {
			rc := sparkFixture(testNS.Name)
			Expect(k8sClient.Create(ctx, rc)).To(Succeed())
		})

		It("requires a positive worker replica count", func() {
			rc := sparkFixture(testNS.Name)
			rc.Spec.Worker.Replicas = pointer.Int32Ptr(-10)

			Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
		})

		DescribeTable("(networking ports)",
			func(portSetter func(*SparkCluster, int32)) {
				rc := sparkFixture(testNS.Name)

				portSetter(rc, 79)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())

				portSetter(rc, 65536)
				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			},
			Entry("rejects an invalid cluster port",
				func(rc *SparkCluster, val int32) { rc.Spec.ClusterPort = val },
			),
			Entry("rejects an invalid worker web port",
				func(rc *SparkCluster, val int32) { rc.Spec.TCPWorkerWebPort = val },
			),
			Entry("rejects an invalid master web port",
				func(rc *SparkCluster, val int32) { rc.Spec.TCPMasterWebPort = val },
			),
			Entry("rejects an invalid dashboard port",
				func(rc *SparkCluster, val int32) { rc.Spec.DashboardPort = val },
			),
		)

		DescribeTable("With mutual tls mode set",
			func(smode string, expectErr bool) {
				rc := sparkFixture(testNS.Name)
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

		Context("With a provided image", func() {
			It("requires a non-blank image registry", func() {
				rc := sparkFixture(testNS.Name)
				rc.Spec.Image = &OCIImageDefinition{Tag: "test-tag"}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("requires a non-blank image tag", func() {
				rc := sparkFixture(testNS.Name)
				rc.Spec.Image = &OCIImageDefinition{Repository: "test-repo"}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})
		})

		Context("With autoscaling enabled", func() {
			clusterWithAutoscaling := func() *SparkCluster {
				rc := sparkFixture(testNS.Name)
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

		Context("framework configs", func() {
			It("rejects a config with only path set", func() {
				rc := sparkFixture(testNS.Name)
				rc.Spec.Worker.FrameworkConfig = &FrameworkConfig{
					Path:    "test/path/",
					Configs: nil,
				}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("rejects a config with only data set", func() {
				rc := sparkFixture(testNS.Name)
				rc.Spec.Worker.FrameworkConfig = &FrameworkConfig{
					Path:    "",
					Configs: map[string]string{"test": "config"},
				}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})
		})

		Context("keytab configs", func() {
			It("rejects a config with only path set", func() {
				rc := sparkFixture(testNS.Name)
				rc.Spec.Worker.KeyTabConfig = &KeyTabConfig{
					Path:   "test/path/",
					KeyTab: nil,
				}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})

			It("rejects a config with only data set", func() {
				rc := sparkFixture(testNS.Name)
				rc.Spec.Worker.KeyTabConfig = &KeyTabConfig{
					Path:   "",
					KeyTab: []byte{'c', 'o', 'n', 'f', 'i', 'g'},
				}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})
		})

		Context("external network policies", func() {
			It("rejects when policy is enabled but no values are set", func() {
				rc := sparkFixture(testNS.Name)
				rc.Spec.NetworkPolicy.ExternalPolicyEnabled = pointer.BoolPtr(true)
				rc.Spec.NetworkPolicy.ExternalPodLabels = map[string]string{}

				Expect(k8sClient.Create(ctx, rc)).ToNot(Succeed())
			})
		})
	})
})

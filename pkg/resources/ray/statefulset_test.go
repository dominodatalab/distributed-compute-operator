package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestNewStatefulSet(t *testing.T) {
	t.Run("invalid_component", func(t *testing.T) {
		rc := rayClusterFixture()
		_, err := NewStatefulSet(rc, Component("garbage"), false)
		assert.Error(t, err)
	})

	t.Run("head", func(t *testing.T) {
		testCommonFeatures(t, ComponentHead)

		t.Run("default_values", func(t *testing.T) {
			rc := rayClusterFixture()
			actual, err := NewStatefulSet(rc, ComponentHead, false)
			require.NoError(t, err)

			expected := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-id-ray-head",
					Namespace: "fake-ns",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "ray",
						"app.kubernetes.io/instance":   "test-id",
						"app.kubernetes.io/component":  "head",
						"app.kubernetes.io/version":    "fake-tag",
						"app.kubernetes.io/managed-by": "distributed-compute-operator",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "test-id-ray-head",
					Replicas:    pointer.Int32Ptr(1),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":      "ray",
							"app.kubernetes.io/instance":  "test-id",
							"app.kubernetes.io/component": "head",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":       "ray",
								"app.kubernetes.io/instance":   "test-id",
								"app.kubernetes.io/component":  "head",
								"app.kubernetes.io/version":    "fake-tag",
								"app.kubernetes.io/managed-by": "distributed-compute-operator",
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "test-id-ray",
							Containers: []corev1.Container{
								{
									Name:            "ray",
									Image:           "docker.io/fake-reg/fake-repo:fake-tag",
									ImagePullPolicy: corev1.PullIfNotPresent,
									Command:         []string{"ray"},
									Args: []string{
										"start",
										"--block",
										"--node-ip-address=$(MY_POD_IP)",
										"--num-cpus=$(MY_CPU_REQUEST)",
										"--object-manager-port=2384",
										"--node-manager-port=2385",
										"--worker-port-list=11000,11001",
										"--head",
										"--ray-client-server-port=10001",
										"--port=6379",
										"--redis-shard-ports=6380,6381",
										"--gcs-server-port=2386",
									},
									Env: []corev1.EnvVar{
										{
											Name: "MY_POD_IP",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													FieldPath: "status.podIP",
												},
											},
										},
										{
											Name: "MY_CPU_REQUEST",
											ValueFrom: &corev1.EnvVarSource{
												ResourceFieldRef: &corev1.ResourceFieldSelector{
													Resource: "requests.cpu",
												},
											},
										},
									},
									Ports: []corev1.ContainerPort{
										{
											Name:          "client",
											ContainerPort: 10001,
										},
										{
											Name:          "object-manager",
											ContainerPort: 2384,
										},
										{
											Name:          "node-manager",
											ContainerPort: 2385,
										},
										{
											Name:          "gcs-server",
											ContainerPort: 2386,
										},
										{
											Name:          "redis-primary",
											ContainerPort: 6379,
										},
										{
											Name:          "redis-shard-0",
											ContainerPort: 6380,
										},
										{
											Name:          "redis-shard-1",
											ContainerPort: 6381,
										},
										{
											Name:          "worker-0",
											ContainerPort: 11000,
										},
										{
											Name:          "worker-1",
											ContainerPort: 11001,
										},
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      sharedMemoryVolumeName,
											MountPath: "/dev/shm",
										},
									},
									LivenessProbe: &corev1.Probe{
										Handler: corev1.Handler{
											TCPSocket: &corev1.TCPSocketAction{
												Port: intstr.FromInt(2385),
											},
										},
									},
									ReadinessProbe: &corev1.Probe{
										Handler: corev1.Handler{
											TCPSocket: &corev1.TCPSocketAction{
												Port: intstr.FromInt(2385),
											},
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "dshm",
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{
											Medium: corev1.StorageMediumMemory,
										},
									},
								},
							},
						},
					},
					PodManagementPolicy: appsv1.ParallelPodManagement,
					UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
						Type: appsv1.RollingUpdateStatefulSetStrategyType,
					},
				},
			}
			assert.Equal(t, expected, actual, "head stateful set not correctly generated")
		})

		t.Run("enable_dashboard", func(t *testing.T) {
			rc := rayClusterFixture()
			rc.Spec.EnableDashboard = pointer.BoolPtr(true)
			rc.Spec.DashboardPort = 8265

			actual, err := NewStatefulSet(rc, ComponentHead, false)
			require.NoError(t, err)

			expected := []string{
				"--include-dashboard=true",
				"--dashboard-host=0.0.0.0",
				"--dashboard-port=8265",
			}
			assert.Subset(t, actual.Spec.Template.Spec.Containers[0].Args, expected)
		})
	})

	t.Run("worker", func(t *testing.T) {
		testCommonFeatures(t, ComponentWorker)

		t.Run("default_values", func(t *testing.T) {
			rc := rayClusterFixture()
			actual, err := NewStatefulSet(rc, ComponentWorker, false)
			require.NoError(t, err)

			expected := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-id-ray-worker",
					Namespace: "fake-ns",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "ray",
						"app.kubernetes.io/instance":   "test-id",
						"app.kubernetes.io/component":  "worker",
						"app.kubernetes.io/version":    "fake-tag",
						"app.kubernetes.io/managed-by": "distributed-compute-operator",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "test-id-ray-worker",
					Replicas:    pointer.Int32Ptr(5),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":      "ray",
							"app.kubernetes.io/instance":  "test-id",
							"app.kubernetes.io/component": "worker",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":       "ray",
								"app.kubernetes.io/instance":   "test-id",
								"app.kubernetes.io/component":  "worker",
								"app.kubernetes.io/version":    "fake-tag",
								"app.kubernetes.io/managed-by": "distributed-compute-operator",
							},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "test-id-ray",
							Containers: []corev1.Container{
								{
									Name:            "ray",
									Image:           "docker.io/fake-reg/fake-repo:fake-tag",
									ImagePullPolicy: corev1.PullIfNotPresent,
									Command:         []string{"ray"},
									Args: []string{
										"start",
										"--block",
										"--node-ip-address=$(MY_POD_IP)",
										"--num-cpus=$(MY_CPU_REQUEST)",
										"--object-manager-port=2384",
										"--node-manager-port=2385",
										"--worker-port-list=11000,11001",
										"--address=test-id-ray-head-0.test-id-ray-head:6379",
									},
									Env: []corev1.EnvVar{
										{
											Name: "MY_POD_IP",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													FieldPath: "status.podIP",
												},
											},
										},
										{
											Name: "MY_CPU_REQUEST",
											ValueFrom: &corev1.EnvVarSource{
												ResourceFieldRef: &corev1.ResourceFieldSelector{
													Resource: "requests.cpu",
												},
											},
										},
									},
									Ports: []corev1.ContainerPort{
										{
											Name:          "object-manager",
											ContainerPort: 2384,
										},
										{
											Name:          "node-manager",
											ContainerPort: 2385,
										},
										{
											Name:          "worker-0",
											ContainerPort: 11000,
										},
										{
											Name:          "worker-1",
											ContainerPort: 11001,
										},
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      sharedMemoryVolumeName,
											MountPath: "/dev/shm",
										},
									},
									LivenessProbe: &corev1.Probe{
										Handler: corev1.Handler{
											TCPSocket: &corev1.TCPSocketAction{
												Port: intstr.FromInt(2385),
											},
										},
									},
									ReadinessProbe: &corev1.Probe{
										Handler: corev1.Handler{
											TCPSocket: &corev1.TCPSocketAction{
												Port: intstr.FromInt(2385),
											},
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "dshm",
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{
											Medium: corev1.StorageMediumMemory,
										},
									},
								},
							},
						},
					},
					PodManagementPolicy: appsv1.ParallelPodManagement,
					UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
						Type: appsv1.RollingUpdateStatefulSetStrategyType,
					},
				},
			}
			assert.Equal(t, expected, actual, "worker stateful set not correctly generated")
		})
	})
}

func testCommonFeatures(t *testing.T, comp Component) {
	t.Helper()

	t.Run("invalid_image", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.Image = &dcv1alpha1.OCIImageDefinition{}

		_, err := NewStatefulSet(rc, comp, false)
		assert.Error(t, err)
	})

	t.Run("object_store_memory", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.ObjectStoreMemoryBytes = pointer.Int64Ptr(100 * 1 << 20)

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Contains(t, actual.Spec.Template.Spec.Containers[0].Args, "--object-store-memory=104857600")
	})

	t.Run("extra_labels", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := map[string]string{
			"thou": "shalt write tests",
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.Labels = expected
		case ComponentWorker:
			rc.Spec.Worker.Labels = expected
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		for _, labels := range []map[string]string{actual.Labels, actual.Spec.Template.Labels} {
			v, ok := labels["thou"]

			assert.True(t, ok)
			assert.Equal(t, "shalt write tests", v)
		}
	})

	t.Run("annotations", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := map[string]string{
			"dominodatalab.com/inject-tooling": "true",
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.Annotations = expected
		case ComponentWorker:
			rc.Spec.Worker.Annotations = expected
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Annotations)
	})

	t.Run("annotations_with_istio", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := map[string]string{
			"dominodatalab.com/inject-tooling": "true",
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.Annotations = expected
			expected["traffic.sidecar.istio.io/includeInboundPorts"] = "10001,8265,2384,2385,2386,6379,6380,6381,11000,11001"
		case ComponentWorker:
			rc.Spec.Worker.Annotations = expected
			expected["traffic.sidecar.istio.io/includeInboundPorts"] = "2384,2385,11000,11001"
		}

		actual, err := NewStatefulSet(rc, comp, true)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Annotations)
	})

	t.Run("volumes_and_mounts", func(t *testing.T) {
		rc := rayClusterFixture()

		expectedVols := []corev1.Volume{
			{
				Name: "extra-vol",
			},
		}
		expectedVolMounts := []corev1.VolumeMount{
			{
				Name: "extra-vol-mount",
			},
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.Volumes = expectedVols
			rc.Spec.Head.VolumeMounts = expectedVolMounts
		case ComponentWorker:
			rc.Spec.Worker.Volumes = expectedVols
			rc.Spec.Worker.VolumeMounts = expectedVolMounts
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Subset(t, actual.Spec.Template.Spec.Volumes, expectedVols)
		assert.Subset(t, actual.Spec.Template.Spec.Containers[0].VolumeMounts, expectedVolMounts)
	})

	t.Run("volume_claim_templates", func(t *testing.T) {
		rc := rayClusterFixture()

		elem := dcv1alpha1.PersistentVolumeClaimTemplate{
			Name: "stuffz",
			Spec: corev1.PersistentVolumeClaimSpec{
				StorageClassName: pointer.StringPtr("test-gpu"),
			},
		}
		input := []dcv1alpha1.PersistentVolumeClaimTemplate{elem}

		switch comp {
		case ComponentHead:
			rc.Spec.Head.VolumeClaimTemplates = input
		case ComponentWorker:
			rc.Spec.Worker.VolumeClaimTemplates = input
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		expected := []corev1.PersistentVolumeClaim{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: elem.Name,
				},
				Spec: elem.Spec,
			},
		}
		assert.Equal(t, expected, actual.Spec.VolumeClaimTemplates)
	})

	t.Run("resource_requirements", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("512Mi"),
			},
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.Resources = expected
		case ComponentWorker:
			rc.Spec.Worker.Resources = expected
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.Containers[0].Resources)
	})

	t.Run("node_selector", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := map[string]string{
			"nodeType": "gpu",
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.NodeSelector = expected
		case ComponentWorker:
			rc.Spec.Worker.NodeSelector = expected
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.NodeSelector)
	})

	t.Run("affinity", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					{
						Weight: 1,
						PodAffinityTerm: corev1.PodAffinityTerm{
							LabelSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      "instance",
										Operator: metav1.LabelSelectorOpIn,
										Values: []string{
											"test-ray",
										},
									},
								},
							},
							TopologyKey: "kubernetes.io/hostname",
						},
					},
				},
			},
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.Affinity = expected
		case ComponentWorker:
			rc.Spec.Worker.Affinity = expected
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.Affinity)
	})

	t.Run("tolerations", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := []corev1.Toleration{
			{
				Key:      "test-key",
				Value:    "test-value",
				Effect:   corev1.TaintEffectNoSchedule,
				Operator: corev1.TolerationOpEqual,
			},
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.Tolerations = expected
		case ComponentWorker:
			rc.Spec.Worker.Tolerations = expected
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.Tolerations)
	})

	t.Run("init_containers", func(t *testing.T) {
		rc := rayClusterFixture()

		expected := []corev1.Container{
			{
				Name: "ray-init",
			},
		}
		switch comp {
		case ComponentHead:
			rc.Spec.Head.InitContainers = expected
		case ComponentWorker:
			rc.Spec.Worker.InitContainers = expected
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.InitContainers)
	})

	t.Run("extra_env_vars", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.EnvVars = []corev1.EnvVar{
			{
				Name:  "foo",
				Value: "bar",
			},
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Subset(t, actual.Spec.Template.Spec.Containers[0].Env, rc.Spec.EnvVars)
	})

	t.Run("security_context", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.PodSecurityContext = &corev1.PodSecurityContext{
			RunAsUser: pointer.Int64Ptr(0),
		}

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.PodSecurityContext, actual.Spec.Template.Spec.SecurityContext)
	})

	t.Run("service_account_override", func(t *testing.T) {
		rc := rayClusterFixture()
		rc.Spec.ServiceAccount.Name = "user-managed-sa"

		actual, err := NewStatefulSet(rc, comp, false)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.ServiceAccount.Name, actual.Spec.Template.Spec.ServiceAccountName)
	})
}

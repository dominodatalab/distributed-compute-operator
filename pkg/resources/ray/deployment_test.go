package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestNewDeployment(t *testing.T) {
	t.Run("invalid_component", func(t *testing.T) {
		rc := testFixture()
		_, err := NewDeployment(rc, Component("garbage"))
		assert.Error(t, err)
	})

	t.Run("head", func(t *testing.T) {
		testCommonFeatures(t, ComponentHead)

		t.Run("default_values", func(t *testing.T) {
			rc := testFixture()
			actual, err := NewDeployment(rc, ComponentHead)
			require.NoError(t, err)

			expected := &appsv1.Deployment{
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
				Spec: appsv1.DeploymentSpec{
					Replicas: pointer.Int32Ptr(1),
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
							ServiceAccountName: "test-id",
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
										"--head",
										"--port=6379",
										"--redis-shard-ports=6380,6381",
									},
									Env: []corev1.EnvVar{
										{
											Name:  "GLOG_logtostderr",
											Value: "1",
										},
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
				},
			}
			assert.Equal(t, expected, actual, "head deployment not correctly generated")
		})

		t.Run("enable_dashboard", func(t *testing.T) {
			rc := testFixture()
			rc.Spec.EnableDashboard = true
			rc.Spec.DashboardPort = 8265

			actual, err := NewDeployment(rc, ComponentHead)
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
			rc := testFixture()
			actual, err := NewDeployment(rc, ComponentWorker)
			require.NoError(t, err)

			expected := &appsv1.Deployment{
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
				Spec: appsv1.DeploymentSpec{
					Replicas: pointer.Int32Ptr(5),
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
							ServiceAccountName: "test-id",
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
										"--address=test-id-ray-head:6379",
									},
									Env: []corev1.EnvVar{
										{
											Name:  "GLOG_logtostderr",
											Value: "1",
										},
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
				},
			}
			assert.Equal(t, expected, actual, "worker deployment not correctly generated")
		})
	})
}

func testFixture() *dcv1alpha1.RayCluster {
	return &dcv1alpha1.RayCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id",
			Namespace: "fake-ns",
		},
		Spec: dcv1alpha1.RayClusterSpec{
			Image: &dcv1alpha1.OCIImageDefinition{
				Registry:   "fake-reg",
				Repository: "fake-repo",
				Tag:        "fake-tag",
				PullPolicy: v1.PullIfNotPresent,
			},
			WorkerReplicaCount: 5,
			HeadPort:           6379,
			RedisShardPorts: []int32{
				6380,
				6381,
			},
			ObjectManagerPort: 2384,
			NodeManagerPort:   2385,
		},
	}
}

func testCommonFeatures(t *testing.T, comp Component) {
	t.Helper()

	t.Run("invalid_image", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.Image = &dcv1alpha1.OCIImageDefinition{}

		_, err := NewDeployment(rc, comp)
		assert.Error(t, err)
	})

	t.Run("object_store_memory", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.ObjectStoreMemoryBytes = 100 * 1 << 20

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		assert.Contains(t, actual.Spec.Template.Spec.Containers[0].Args, "--object-store-memory=104857600")
	})

	t.Run("extra_labels", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.Labels = map[string]string{
			"thou": "shalt write tests",
		}

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		for _, labels := range []map[string]string{actual.Labels, actual.Spec.Template.Labels} {
			v, ok := labels["thou"]

			assert.True(t, ok)
			assert.Equal(t, "shalt write tests", v)
		}
	})

	t.Run("annotations", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.Annotations = map[string]string{
			"dominodatalab.com/inject-tooling": "true",
		}

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.Annotations, actual.Spec.Template.Annotations)
	})

	t.Run("volumes_and_mounts", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.Volumes = []corev1.Volume{
			{
				Name: "extra-vol",
			},
		}
		rc.Spec.VolumeMounts = []corev1.VolumeMount{
			{
				Name: "extra-vol-mount",
			},
		}

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		assert.Subset(t, actual.Spec.Template.Spec.Volumes, rc.Spec.Volumes)
		assert.Subset(t, actual.Spec.Template.Spec.Containers[0].VolumeMounts, rc.Spec.VolumeMounts)
	})

	t.Run("resource_requirements", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("1G"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("512Mi"),
			},
		}

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.Resources, actual.Spec.Template.Spec.Containers[0].Resources)
	})

	t.Run("node_selector", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.NodeSelector = map[string]string{
			"nodeType": "gpu",
		}

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.NodeSelector, actual.Spec.Template.Spec.NodeSelector)
	})

	t.Run("affinity", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.Affinity = &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
					{
						Weight: 1,
						PodAffinityTerm: v1.PodAffinityTerm{
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

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.Affinity, actual.Spec.Template.Spec.Affinity)
	})

	t.Run("tolerations", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.Tolerations = []corev1.Toleration{
			{
				Key:      "test-key",
				Value:    "test-value",
				Effect:   v1.TaintEffectNoSchedule,
				Operator: v1.TolerationOpEqual,
			},
		}

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.Tolerations, actual.Spec.Template.Spec.Tolerations)
	})

	t.Run("init_containers", func(t *testing.T) {
		rc := testFixture()
		rc.Spec.InitContainers = []corev1.Container{
			{
				Name: "ray-init",
			},
		}

		actual, err := NewDeployment(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.InitContainers, actual.Spec.Template.Spec.InitContainers)
	})
}

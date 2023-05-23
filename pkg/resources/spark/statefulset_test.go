package spark

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
		rc := sparkClusterFixture()
		_, err := NewStatefulSet(rc, "garbage")
		assert.Error(t, err)
	})

	t.Run("master", func(t *testing.T) {
		testCommonFeatures(t, ComponentMaster)

		t.Run("default_values", func(t *testing.T) {
			rc := sparkClusterFixture()
			actual, err := NewStatefulSet(rc, ComponentMaster)
			require.NoError(t, err)

			expected := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-id-spark-master",
					Namespace: "fake-ns",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "spark",
						"app.kubernetes.io/instance":   "test-id",
						"app.kubernetes.io/component":  "master",
						"app.kubernetes.io/version":    "fake-tag",
						"app.kubernetes.io/managed-by": "distributed-compute-operator",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "test-id-spark-master",
					Replicas:    pointer.Int32(1),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":      "spark",
							"app.kubernetes.io/instance":  "test-id",
							"app.kubernetes.io/component": "master",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":       "spark",
								"app.kubernetes.io/instance":   "test-id",
								"app.kubernetes.io/component":  "master",
								"app.kubernetes.io/version":    "fake-tag",
								"app.kubernetes.io/managed-by": "distributed-compute-operator",
							},
							Annotations: map[string]string{},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "test-id-spark",
							Containers: []corev1.Container{
								{
									Name:            "test-id-spark-master",
									Image:           "docker.io/fake-reg/fake-repo:fake-tag",
									ImagePullPolicy: corev1.PullIfNotPresent,
									Env: []corev1.EnvVar{
										{
											Name:  "SPARK_MASTER_PORT",
											Value: "7077",
										},
										{
											Name:  "SPARK_MASTER_WEBUI_PORT",
											Value: "8080",
										},
										{
											Name:  "SPARK_MODE",
											Value: "master",
										},
									},
									Ports: []corev1.ContainerPort{
										{
											Name:          "http",
											ContainerPort: 8080,
											Protocol:      "TCP",
										},
										{
											Name:          "cluster",
											ContainerPort: 7077,
										},
									},
									LivenessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Port: intstr.FromInt(8080),
												Path: "/",
											},
										},
									},
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Port: intstr.FromInt(8080),
												Path: "/",
											},
										},
									},
								},
							},
							SecurityContext: &corev1.PodSecurityContext{
								RunAsUser: pointer.Int64(1001),
								FSGroup:   pointer.Int64(1001),
							},
						},
					},
					UpdateStrategy:      appsv1.StatefulSetUpdateStrategy{Type: appsv1.RollingUpdateStatefulSetStrategyType},
					PodManagementPolicy: appsv1.ParallelPodManagement,
				},
			}
			assert.Equal(t, expected, actual, "master statefulset not correctly generated")
		})
	})

	t.Run("worker", func(t *testing.T) {
		testCommonFeatures(t, ComponentWorker)

		t.Run("default_values", func(t *testing.T) {
			rc := sparkClusterFixture()
			actual, err := NewStatefulSet(rc, ComponentWorker)
			require.NoError(t, err)

			expected := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-id-spark-worker",
					Namespace: "fake-ns",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "spark",
						"app.kubernetes.io/instance":   "test-id",
						"app.kubernetes.io/component":  "worker",
						"app.kubernetes.io/version":    "fake-tag",
						"app.kubernetes.io/managed-by": "distributed-compute-operator",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "test-id-spark-worker",
					Replicas:    pointer.Int32(5),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":      "spark",
							"app.kubernetes.io/instance":  "test-id",
							"app.kubernetes.io/component": "worker",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app.kubernetes.io/name":       "spark",
								"app.kubernetes.io/instance":   "test-id",
								"app.kubernetes.io/component":  "worker",
								"app.kubernetes.io/version":    "fake-tag",
								"app.kubernetes.io/managed-by": "distributed-compute-operator",
							},
							Annotations: map[string]string{},
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "test-id-spark",
							Containers: []corev1.Container{
								{
									Name:            "test-id-spark-worker",
									Image:           "docker.io/fake-reg/fake-repo:fake-tag",
									ImagePullPolicy: corev1.PullIfNotPresent,
									Env: []corev1.EnvVar{
										{
											Name:  "SPARK_MASTER_URL",
											Value: "spark://test-id-spark-master:7077",
										},
										{
											Name:  "SPARK_WORKER_WEBUI_PORT",
											Value: "8081",
										},
										{
											Name:  "SPARK_WORKER_PORT",
											Value: "7077",
										},
										{
											Name:  "SPARK_MODE",
											Value: "worker",
										},
										{
											Name:  "SPARK_WORKER_MEMORY",
											Value: "4505m",
										},
										{
											Name:  "SPARK_WORKER_CORES",
											Value: "0",
										},
									},
									Ports: []corev1.ContainerPort{
										{
											Name:          "http",
											ContainerPort: 8081,
											Protocol:      "TCP",
										},
									},
									LivenessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Port: intstr.FromInt(8081),
												Path: "/",
											},
										},
									},
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Port: intstr.FromInt(8081),
												Path: "/",
											},
										},
									},
								},
							},
							SecurityContext: &corev1.PodSecurityContext{
								RunAsUser: pointer.Int64(1001),
								FSGroup:   pointer.Int64(1001),
							},
						},
					},
					UpdateStrategy:      appsv1.StatefulSetUpdateStrategy{Type: appsv1.RollingUpdateStatefulSetStrategyType},
					PodManagementPolicy: appsv1.ParallelPodManagement,
				},
			}
			assert.Equal(t, expected, actual, "worker statefulset not correctly generated")
		})
	})
}

func testCommonFeatures(t *testing.T, comp Component) {
	t.Helper()

	t.Run("invalid_image", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.Image = &dcv1alpha1.OCIImageDefinition{}

		_, err := NewStatefulSet(rc, comp)
		assert.Error(t, err)
	})

	t.Run("extra_labels", func(t *testing.T) {
		rc := sparkClusterFixture()

		expected := map[string]string{
			"thou": "shalt write tests",
		}
		switch comp {
		case ComponentMaster:
			rc.Spec.Master.Labels = expected
		case ComponentWorker:
			rc.Spec.Worker.Labels = expected
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		for _, labels := range []map[string]string{actual.Labels, actual.Spec.Template.Labels} {
			v, ok := labels["thou"]

			assert.True(t, ok)
			assert.Equal(t, "shalt write tests", v)
		}
	})

	t.Run("annotations", func(t *testing.T) {
		rc := sparkClusterFixture()

		expected := map[string]string{
			"dominodatalab.com/inject-tooling": "true",
		}
		switch comp {
		case ComponentMaster:
			rc.Spec.Master.Annotations = expected
		case ComponentWorker:
			rc.Spec.Worker.Annotations = expected
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Annotations)
	})

	t.Run("volumes_and_mounts", func(t *testing.T) {
		rc := sparkClusterFixture()

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
		case ComponentMaster:
			rc.Spec.Master.Volumes = expectedVols
			rc.Spec.Master.VolumeMounts = expectedVolMounts
		case ComponentWorker:
			rc.Spec.Worker.Volumes = expectedVols
			rc.Spec.Worker.VolumeMounts = expectedVolMounts
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Subset(t, actual.Spec.Template.Spec.Volumes, expectedVols)
		assert.Subset(t, actual.Spec.Template.Spec.Containers[0].VolumeMounts, expectedVolMounts)
	})

	t.Run("resource_requirements", func(t *testing.T) {
		rc := sparkClusterFixture()

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
		case ComponentMaster:
			rc.Spec.Master.Resources = expected
		case ComponentWorker:
			rc.Spec.Worker.Resources = expected
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.Containers[0].Resources)
	})

	t.Run("node_selector", func(t *testing.T) {
		rc := sparkClusterFixture()

		expected := map[string]string{
			"nodeType": "gpu",
		}
		switch comp {
		case ComponentMaster:
			rc.Spec.Master.NodeSelector = expected
		case ComponentWorker:
			rc.Spec.Worker.NodeSelector = expected
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.NodeSelector)
	})

	t.Run("affinity", func(t *testing.T) {
		rc := sparkClusterFixture()

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
											"test",
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
		case ComponentMaster:
			rc.Spec.Master.Affinity = expected
		case ComponentWorker:
			rc.Spec.Worker.Affinity = expected
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.Affinity)
	})

	t.Run("tolerations", func(t *testing.T) {
		rc := sparkClusterFixture()

		expected := []corev1.Toleration{
			{
				Key:      "test-key",
				Value:    "test-value",
				Effect:   corev1.TaintEffectNoSchedule,
				Operator: corev1.TolerationOpEqual,
			},
		}
		switch comp {
		case ComponentMaster:
			rc.Spec.Master.Tolerations = expected
		case ComponentWorker:
			rc.Spec.Worker.Tolerations = expected
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.Tolerations)
	})

	t.Run("init_containers", func(t *testing.T) {
		rc := sparkClusterFixture()

		expected := []corev1.Container{
			{
				Name: "spark-init",
			},
		}
		switch comp {
		case ComponentMaster:
			rc.Spec.Master.InitContainers = expected
		case ComponentWorker:
			rc.Spec.Worker.InitContainers = expected
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, expected, actual.Spec.Template.Spec.InitContainers)
	})

	t.Run("extra_env_vars", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.EnvVars = []corev1.EnvVar{
			{
				Name:  "foo",
				Value: "bar",
			},
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Subset(t, actual.Spec.Template.Spec.Containers[0].Env, rc.Spec.EnvVars)
	})

	t.Run("security_context", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.PodSecurityContext = &corev1.PodSecurityContext{
			RunAsUser: pointer.Int64(0),
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.PodSecurityContext, actual.Spec.Template.Spec.SecurityContext)
	})

	t.Run("service_account_override", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.ServiceAccount.Name = "user-managed-sa"

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, rc.Spec.ServiceAccount.Name, actual.Spec.Template.Spec.ServiceAccountName)
	})

	t.Run("volume_claim_template", func(t *testing.T) {
		sc := sparkClusterFixture()

		elem := dcv1alpha1.PersistentVolumeClaimTemplate{
			Name: "stuffz",
			Spec: corev1.PersistentVolumeClaimSpec{
				StorageClassName: pointer.String("test-gpu"),
			},
		}
		input := []dcv1alpha1.PersistentVolumeClaimTemplate{elem}

		switch comp {
		case ComponentMaster:
			sc.Spec.Master.VolumeClaimTemplates = input
		case ComponentWorker:
			sc.Spec.Worker.VolumeClaimTemplates = input
		}

		actual, err := NewStatefulSet(sc, comp)
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

	t.Run("default_configuration", func(t *testing.T) {
		rc := sparkClusterFixture()
		fcMaster := map[string]string{
			"m1": "v1",
		}

		fcWorker := map[string]string{
			"w1": "v1",
		}

		rc.Spec.Master = dcv1alpha1.SparkClusterNode{
			DefaultConfiguration: fcMaster,
		}
		rc.Spec.Worker = dcv1alpha1.SparkClusterWorker{
			SparkClusterNode: dcv1alpha1.SparkClusterNode{
				DefaultConfiguration: fcWorker,
			},
			Replicas: pointer.Int32(2),
		}

		expectedVolumes := []corev1.Volume{
			{
				Name: "spark-config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "test-id-framework-spark",
						},
					},
				},
			},
		}

		expectedVolumeMounts := []corev1.VolumeMount{
			{
				Name:      "spark-config",
				ReadOnly:  false,
				MountPath: "/opt/bitnami/spark/conf/spark-defaults.conf",
				SubPath:   string(comp),
			},
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, expectedVolumes, actual.Spec.Template.Spec.Volumes)
		assert.Equal(t, expectedVolumeMounts, actual.Spec.Template.Spec.Containers[0].VolumeMounts)
	})

	t.Run("keytab_config", func(t *testing.T) {
		rc := sparkClusterFixture()

		keyTab := dcv1alpha1.KerberosKeytabConfig{
			MountPath: "/test/path",
			Contents:  []byte{'t', 'e', 's', 't', 'e', 'r'},
		}

		rc.Spec.KerberosKeytab = &keyTab

		expectedVolumes := []corev1.Volume{
			{
				Name: "keytab",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "test-id-keytab-spark",
						},
					},
				},
			},
		}

		expectedVolumeMounts := []corev1.VolumeMount{
			{
				Name:      "keytab",
				ReadOnly:  false,
				MountPath: "/test/path",
				SubPath:   string(comp),
			},
		}

		actual, err := NewStatefulSet(rc, comp)
		require.NoError(t, err)

		assert.Equal(t, expectedVolumes, actual.Spec.Template.Spec.Volumes)
		assert.Equal(t, expectedVolumeMounts, actual.Spec.Template.Spec.Containers[0].VolumeMounts)
	})
}

package spark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestNewFrameworkConfigMap(t *testing.T) {
	t.Run("fully loaded", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.Master.FrameworkConfig = &v1alpha1.FrameworkConfig{
			Configs: map[string]string{
				"m1": "v1",
				"m2": "v2",
			},
		}
		rc.Spec.Worker.FrameworkConfig = &v1alpha1.FrameworkConfig{
			Configs: map[string]string{
				"w1": "v1",
				"w2": "v2",
			},
		}
		cm := NewFrameworkConfigMap(rc)

		expected := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-id-framework",
				Namespace: "fake-ns",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "spark",
					"app.kubernetes.io/instance":   "test-id",
					"app.kubernetes.io/version":    "fake-tag",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
				},
			},
			Data: map[string]string{
				"master": `m1 v1
m2 v2
`,
				"worker": `w1 v1
w2 v2
`,
			},
		}
		assert.Equal(t, expected, cm)
	})

	t.Run("only one node type", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.Master.FrameworkConfig = &v1alpha1.FrameworkConfig{
			Configs: map[string]string{
				"m1": "v1",
				"m2": "v2",
			},
		}
		cm := NewFrameworkConfigMap(rc)

		expected := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-id-framework",
				Namespace: "fake-ns",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "spark",
					"app.kubernetes.io/instance":   "test-id",
					"app.kubernetes.io/version":    "fake-tag",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
				},
			},
			Data: map[string]string{
				"master": `m1 v1
m2 v2
`,
			},
		}
		assert.Equal(t, expected, cm)
	})
	t.Run("no nodes", func(t *testing.T) {
		rc := sparkClusterFixture()
		cm := NewFrameworkConfigMap(rc)
		assert.Nil(t, cm)
	})
}

func TestGenerateSparkDefaults(t *testing.T) {
	actual := generateSparkDefaults(map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
		"e": "5",
	})

	expected := `a 1
b 2
c 3
d 4
e 5
`

	assert.Equal(t, expected, actual)
}

func TestNewKeyTabConfigMap(t *testing.T) {
	t.Run("fully loaded", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.Master.KeyTabConfig = &v1alpha1.KeyTabConfig{
			Path:   "ignore-me",
			KeyTab: []byte{'m', 'a', 's', 't', 'e', 'r'},
		}
		rc.Spec.Worker.KeyTabConfig = &v1alpha1.KeyTabConfig{
			Path:   "ignore-me",
			KeyTab: []byte{'w', 'o', 'r', 'k', 'e', 'r'},
		}
		cm := NewKeyTabConfigMap(rc)

		expected := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-id-keytab",
				Namespace: "fake-ns",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "spark",
					"app.kubernetes.io/instance":   "test-id",
					"app.kubernetes.io/version":    "fake-tag",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
				},
			},
			BinaryData: map[string][]byte{
				"master": {'m', 'a', 's', 't', 'e', 'r'},
				"worker": {'w', 'o', 'r', 'k', 'e', 'r'},
			},
		}
		assert.Equal(t, expected, cm)
	})

	t.Run("only one node type", func(t *testing.T) {
		rc := sparkClusterFixture()
		rc.Spec.Master.KeyTabConfig = &v1alpha1.KeyTabConfig{
			Path:   "ignore-me",
			KeyTab: []byte{'m', 'a', 's', 't', 'e', 'r'},
		}
		cm := NewKeyTabConfigMap(rc)

		expected := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-id-keytab",
				Namespace: "fake-ns",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "spark",
					"app.kubernetes.io/instance":   "test-id",
					"app.kubernetes.io/version":    "fake-tag",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
				},
			},
			BinaryData: map[string][]byte{
				"master": {'m', 'a', 's', 't', 'e', 'r'},
			},
		}
		assert.Equal(t, expected, cm)
	})
	t.Run("no nodes", func(t *testing.T) {
		rc := sparkClusterFixture()
		cm := NewKeyTabConfigMap(rc)
		assert.Nil(t, cm)
	})
}

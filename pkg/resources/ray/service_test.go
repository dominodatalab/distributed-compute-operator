package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewHeadService(t *testing.T) {
	rc := rayClusterFixture()
	svc := NewHeadService(rc)

	expected := &corev1.Service{
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
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "client",
					Port: 10001,
				},
				{
					Name: "redis-primary",
					Port: 6379,
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/name":      "ray",
				"app.kubernetes.io/instance":  "test-id",
				"app.kubernetes.io/component": "head",
			},
		},
	}
	assert.Equal(t, expected, svc)

	t.Run("with_dashboard", func(t *testing.T) {
		rc.Spec.EnableDashboard = true
		svc := NewHeadService(rc)

		expected.Spec.Ports = append(expected.Spec.Ports, corev1.ServicePort{
			Name: "dashboard",
			Port: 8265,
		})

		assert.Equal(t, expected, svc)
	})
}

package spark

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestNewHeadService(t *testing.T) {
	rc := sparkClusterFixture()
	svc := NewMasterService(rc)

	expected := &corev1.Service{
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
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "cluster",
					Port:       7077,
					TargetPort: intstr.FromString("cluster"),
				},
			},
			Type: "ClusterIP",
			Selector: map[string]string{
				"app.kubernetes.io/name":      "spark",
				"app.kubernetes.io/instance":  "test-id",
				"app.kubernetes.io/component": "master",
			},
		},
	}
	assert.Equal(t, expected, svc)

	t.Run("with_dashboard_enabled", func(t *testing.T) {
		rc.Spec.EnableDashboard = pointer.BoolPtr(true)
		svc := NewMasterService(rc)

		expected.Spec.Ports = append(expected.Spec.Ports, corev1.ServicePort{
			Name:       "tcp",
			Protocol:   corev1.ProtocolTCP,
			Port:       8265,
			TargetPort: intstr.FromString("http"),
		})

		assert.Equal(t, expected, svc)
	})
}

func TestNewHeadlessService(t *testing.T) {
	rc := sparkClusterFixture()
	svc := NewHeadlessService(rc)

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-spark-worker",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/component":  "master",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector: map[string]string{
				"app.kubernetes.io/name":      "spark",
				"app.kubernetes.io/instance":  "test-id",
			},
		},
	}
	assert.Equal(t, expected, svc)
}

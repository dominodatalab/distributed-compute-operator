package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

func TestNewClientService(t *testing.T) {
	rc := rayClusterFixture()
	svc := NewClientService(rc)

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-ray-client",
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
					Name:       "tcp-client",
					Port:       10001,
					TargetPort: intstr.FromInt(10001),
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

	t.Run("with_dashboard_enabled", func(t *testing.T) {
		rc.Spec.EnableDashboard = ptr.To(true)
		svc := NewClientService(rc)

		expected.Spec.Ports = append(expected.Spec.Ports, corev1.ServicePort{
			Name:       "tcp-dashboard",
			Port:       8265,
			TargetPort: intstr.FromInt(8265),
		})

		assert.Equal(t, expected, svc)
	})
}

func TestNewHeadlessHeadService(t *testing.T) {
	rc := rayClusterFixture()
	svc := NewHeadlessHeadService(rc)

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
					Name:       "tcp-gcs-server",
					Port:       2386,
					TargetPort: intstr.FromInt(2386),
				},
				{
					Name:       "tcp-redis-primary",
					Port:       6379,
					TargetPort: intstr.FromInt(6379),
				},
				{
					Name:       "tcp-redis-shard-0",
					Port:       6380,
					TargetPort: intstr.FromInt(6380),
				},
				{
					Name:       "tcp-redis-shard-1",
					Port:       6381,
					TargetPort: intstr.FromInt(6381),
				},
				{
					Name:       "tcp-object-manager",
					Port:       2384,
					TargetPort: intstr.FromInt(2384),
				},
				{
					Name:       "tcp-node-manager",
					Port:       2385,
					TargetPort: intstr.FromInt(2385),
				},
				{
					Name:       "tcp-worker-port-0",
					Port:       11000,
					TargetPort: intstr.FromInt(11000),
				},
				{
					Name:       "tcp-worker-port-1",
					Port:       11001,
					TargetPort: intstr.FromInt(11001),
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/name":      "ray",
				"app.kubernetes.io/instance":  "test-id",
				"app.kubernetes.io/component": "head",
			},
			ClusterIP: corev1.ClusterIPNone,
		},
	}
	assert.Equal(t, expected, svc)
}

func TestNewHeadlessWorkerService(t *testing.T) {
	rc := rayClusterFixture()
	svc := NewHeadlessWorkerService(rc)

	expected := &corev1.Service{
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
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp-object-manager",
					Port:       2384,
					TargetPort: intstr.FromInt(2384),
				},
				{
					Name:       "tcp-node-manager",
					Port:       2385,
					TargetPort: intstr.FromInt(2385),
				},
				{
					Name:       "tcp-worker-port-0",
					Port:       11000,
					TargetPort: intstr.FromInt(11000),
				},
				{
					Name:       "tcp-worker-port-1",
					Port:       11001,
					TargetPort: intstr.FromInt(11001),
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/name":      "ray",
				"app.kubernetes.io/instance":  "test-id",
				"app.kubernetes.io/component": "worker",
			},
			ClusterIP: corev1.ClusterIPNone,
		},
	}
	assert.Equal(t, expected, svc)
}

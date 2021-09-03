package spark

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const SparkBlockManagerPortName = "spark-block-manager-port"

func TestNewMasterService(t *testing.T) {
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
					Name:       "tcp-cluster",
					Port:       7077,
					TargetPort: intstr.FromString("cluster"),
				},
				{
					Name:       "tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.FromString("http"),
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
}

func TestNewHeadlessService(t *testing.T) {
	rc := sparkClusterFixture()
	rc.Spec.MasterWebPort = 8080
	rc.Spec.Driver.BlockManagerPort = 4042

	svc := NewHeadlessService(rc)

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-spark-worker",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/component":  "worker",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector: map[string]string{
				"app.kubernetes.io/name":     "spark",
				"app.kubernetes.io/instance": "test-id",
			},
			Ports: []corev1.ServicePort{
				// these ports are exposed for Istio support
				{
					Name:       "tcp-cluster",
					Port:       7077,
					TargetPort: intstr.FromString("cluster"),
				},
				{
					Name:       "tcp-master-webport",
					Port:       8080,
					TargetPort: intstr.FromString("http"),
					Protocol:   corev1.ProtocolTCP,
				}, {
					Name:       "tcp-worker-webport",
					Port:       8081,
					TargetPort: intstr.FromString("http"),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-driver-block-manager",
					Port:     4042,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
	assert.Equal(t, expected, svc)
}

func TestNewSparkDriverService(t *testing.T) {
	rc := sparkClusterFixture()
	rc.Spec.Driver.Selector = map[string]string{
		"app.kubernetes.io/instance": "test-id",
	}
	rc.Spec.Driver.UIPort = 4040
	rc.Spec.Driver.Port = 4041
	rc.Spec.Driver.BlockManagerPort = 4042

	svc := NewSparkDriverService(rc)

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-spark-driver",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "spark",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
		},
		Spec: corev1.ServiceSpec{
			Type:      "ClusterIP",
			ClusterIP: "None",
			Selector: map[string]string{
				"app.kubernetes.io/instance": "test-id",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp-ui",
					Port:       4040,
					TargetPort: intstr.FromInt(4040),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-driver",
					Port:     4041,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-block-manager",
					Port:     4042,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
	assert.Equal(t, expected, svc)
}

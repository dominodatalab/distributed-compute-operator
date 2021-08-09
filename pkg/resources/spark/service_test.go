package spark

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

const (
	DriverUIPort              = 4040
	DriverUIPortName          = "spark-ui-port"
	DriverPort                = 4041
	DriverPortName            = "spark-driver-port"
	BlockManagerPort          = 4042
	SparkBlockManagerPortName = "spark-block-manager-port"
)

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
			Port:       8080,
			TargetPort: intstr.FromString("http"),
		})

		assert.Equal(t, expected, svc)
	})
}

func TestNewHeadlessService(t *testing.T) {
	rc := sparkClusterFixture()
	rc.Spec.TCPMasterWebPort = 8080
	rc.Spec.Driver.DriverUIPort = DriverUIPort
	rc.Spec.Driver.DriverUIPortName = DriverUIPortName
	rc.Spec.Driver.DriverPort = DriverPort
	rc.Spec.Driver.DriverPortName = DriverPortName
	rc.Spec.Driver.DriverBlockManagerPort = BlockManagerPort
	rc.Spec.Driver.DriverBlockManagerPortName = SparkBlockManagerPortName
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
					Name:       "tcp-spark-ui-port",
					Port:       4040,
					TargetPort: intstr.FromInt(4040),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-spark-driver-port",
					Port:     4041,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-spark-block-manager-port",
					Port:     4042,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
	assert.Equal(t, expected, svc)
}

func TestNewSparkDriverService(t *testing.T) {
	const clusterName = "test-id"

	rc := sparkClusterFixture()
	rc.Spec.Driver.SparkClusterName = clusterName
	rc.Spec.Driver.ExecutionName = clusterName
	rc.Spec.Driver.DriverUIPort = DriverUIPort
	rc.Spec.Driver.DriverUIPortName = DriverUIPortName
	rc.Spec.Driver.DriverPort = DriverPort
	rc.Spec.Driver.DriverPortName = DriverPortName
	rc.Spec.Driver.DriverBlockManagerPort = BlockManagerPort
	rc.Spec.Driver.DriverBlockManagerPortName = SparkBlockManagerPortName

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
					Name:       "tcp-spark-ui-port",
					Port:       4040,
					TargetPort: intstr.FromInt(4040),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-spark-driver-port",
					Port:     4041,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "tcp-spark-block-manager-port",
					Port:     4042,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
	assert.Equal(t, expected, svc)
}

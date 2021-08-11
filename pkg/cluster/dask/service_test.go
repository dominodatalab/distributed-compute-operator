package dask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestServiceDataSource_Service(t *testing.T) {
	dc := testDaskCluster()

	t.Run("scheduler", func(t *testing.T) {
		ds := serviceDS{dc: dc, comp: ComponentScheduler}

		actual := ds.Service()
		expected := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-dask-scheduler",
				Namespace: "ns",
				Labels: map[string]string{
					"app.kubernetes.io/component":  "scheduler",
					"app.kubernetes.io/instance":   "test",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
					"app.kubernetes.io/name":       "dask",
					"app.kubernetes.io/version":    "test-tag",
				},
			},
			Spec: corev1.ServiceSpec{
				ClusterIP: corev1.ClusterIPNone,
				Selector: map[string]string{
					"app.kubernetes.io/component": "scheduler",
					"app.kubernetes.io/instance":  "test",
					"app.kubernetes.io/name":      "dask",
				},
				Ports: []corev1.ServicePort{
					{
						Name:       "tcp-serve",
						Port:       8786,
						TargetPort: intstr.FromString("serve"),
					},
					{
						Name:       "tcp-dashboard",
						Port:       8787,
						TargetPort: intstr.FromString("dashboard"),
					},
				},
			},
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("worker", func(t *testing.T) {
		ds := serviceDS{dc: dc, comp: ComponentWorker}

		actual := ds.Service()
		expected := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-dask-worker",
				Namespace: "ns",
				Labels: map[string]string{
					"app.kubernetes.io/component":  "worker",
					"app.kubernetes.io/instance":   "test",
					"app.kubernetes.io/managed-by": "distributed-compute-operator",
					"app.kubernetes.io/name":       "dask",
					"app.kubernetes.io/version":    "test-tag",
				},
			},
			Spec: corev1.ServiceSpec{
				ClusterIP: corev1.ClusterIPNone,
				Selector: map[string]string{
					"app.kubernetes.io/component": "worker",
					"app.kubernetes.io/instance":  "test",
					"app.kubernetes.io/name":      "dask",
				},
				Ports: []corev1.ServicePort{
					{
						Name:       "tcp-worker",
						Port:       3000,
						TargetPort: intstr.FromString("worker"),
					},
					{
						Name:       "tcp-nanny",
						Port:       3001,
						TargetPort: intstr.FromString("nanny"),
					},
					{
						Name:       "tcp-dashboard",
						Port:       8787,
						TargetPort: intstr.FromString("dashboard"),
					},
				},
			},
		}

		assert.Equal(t, expected, actual)
	})
}

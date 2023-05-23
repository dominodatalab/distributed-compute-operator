package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestNewServiceAccount(t *testing.T) {
	rc := rayClusterFixture()
	sa := NewServiceAccount(rc)

	expected := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-id-ray",
			Namespace: "fake-ns",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "ray",
				"app.kubernetes.io/instance":   "test-id",
				"app.kubernetes.io/version":    "fake-tag",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
			},
		},
		AutomountServiceAccountToken: pointer.Bool(false),
	}
	assert.Equal(t, expected, sa)
}

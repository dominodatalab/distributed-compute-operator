package dask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestServiceAccountDS_ServiceAccount(t *testing.T) {
	dc := testDaskCluster()
	ds := serviceAccountDS{dc: dc}

	actual := ds.ServiceAccount()
	expected := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-dask",
			Namespace: "ns",
			Labels: map[string]string{
				"app.kubernetes.io/instance":   "test",
				"app.kubernetes.io/managed-by": "distributed-compute-operator",
				"app.kubernetes.io/name":       "dask",
				"app.kubernetes.io/version":    "test-tag",
			},
		},
		AutomountServiceAccountToken: ptr.To(false),
	}

	require.Equal(t, expected, actual)

	dc.Spec.ServiceAccount.AutomountServiceAccountToken = true
	actual = ds.ServiceAccount()

	assert.Equal(t, actual.AutomountServiceAccountToken, ptr.To(true))
}

func TestServiceAccountDS_Delete(t *testing.T) {
	dc := testDaskCluster()
	ds := serviceAccountDS{dc: dc}

	t.Run("empty_name", func(t *testing.T) {
		dc.Spec.ServiceAccount.Name = ""
		assert.False(t, ds.Delete())
	})

	t.Run("provided_name", func(t *testing.T) {
		dc.Spec.ServiceAccount.Name = "other"
		assert.True(t, ds.Delete())
	})
}

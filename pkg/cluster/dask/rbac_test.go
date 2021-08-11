package dask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPspDS_Role(t *testing.T) {
	dc := testDaskCluster()
	ds := pspDS{dc: dc}

	actual := ds.Role()
	expected := &rbacv1.Role{
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
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"policy"},
				Resources:     []string{"podsecuritypolicies"},
				Verbs:         []string{"use"},
				ResourceNames: []string{"privileged"},
			},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestPspDS_RoleBinding(t *testing.T) {
	dc := testDaskCluster()
	ds := pspDS{dc: dc}

	actual := ds.RoleBinding()
	expected := &rbacv1.RoleBinding{
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
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     "test-dask",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      "test-dask",
				Namespace: "ns",
			},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestPspDS_Delete(t *testing.T) {
	dc := testDaskCluster()
	ds := pspDS{dc: dc}

	t.Run("provided_name", func(t *testing.T) {
		dc.Spec.PodSecurityPolicy = "restricted"
		assert.False(t, ds.Delete())
	})

	t.Run("empty_name", func(t *testing.T) {
		dc.Spec.PodSecurityPolicy = ""
		assert.True(t, ds.Delete())
	})
}

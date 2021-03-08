package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewPodSecurityPolicyRBAC(t *testing.T) {
	rc := rayClusterFixture()
	rc.Spec.PodSecurityPolicy = "test-psp"
	role, roleBinding := NewPodSecurityPolicyRBAC(rc)

	t.Run("role", func(t *testing.T) {
		expected := &rbacv1.Role{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Role",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
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
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups:     []string{"policy"},
					Resources:     []string{"podsecuritypolicies"},
					Verbs:         []string{"use"},
					ResourceNames: []string{"test-psp"},
				},
			},
		}
		assert.Equal(t, expected, role)
	})

	t.Run("role_binding", func(t *testing.T) {
		expected := &rbacv1.RoleBinding{
			TypeMeta: metav1.TypeMeta{
				Kind:       "RoleBinding",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
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
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     "test-id-ray",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "test-id-ray",
					Namespace: "fake-ns",
				},
			},
		}
		assert.Equal(t, expected, roleBinding)
	})
}

package ray

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

var (
	policyAPIGroups            = []string{"policy"}
	podSecurityPolicyResources = []string{"podsecuritypolicies"}
	useVerbs                   = []string{"use"}
)

// NewPodSecurityPolicyRBAC generates the role and role binding required to use a pod security policy.
// The role is bound to the service account used by the ray cluster pods.
func NewPodSecurityPolicyRBAC(rc *dcv1alpha1.RayCluster) (*rbacv1.Role, *rbacv1.RoleBinding) {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rc.Name,
			Namespace: rc.Namespace,
			Labels:    MetadataLabels(rc),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     policyAPIGroups,
				Resources:     podSecurityPolicyResources,
				Verbs:         useVerbs,
				ResourceNames: []string{rc.Spec.PodSecurityPolicy},
			},
		},
	}

	binding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rc.Name,
			Namespace: rc.Namespace,
			Labels:    MetadataLabels(rc),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     role.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      rc.Name,
				Namespace: rc.Namespace,
			},
		},
	}

	return role, binding
}

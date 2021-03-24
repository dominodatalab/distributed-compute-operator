package spark

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
// The role is bound to the service account used by the spark cluster pods.
func NewPodSecurityPolicyRBAC(sc *dcv1alpha1.SparkCluster) (*rbacv1.Role, *rbacv1.RoleBinding) {
	name := InstanceObjectName(sc.Name, ComponentNone)

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: sc.Namespace,
			Labels:    MetadataLabels(sc),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     policyAPIGroups,
				Resources:     podSecurityPolicyResources,
				Verbs:         useVerbs,
				ResourceNames: []string{sc.Spec.PodSecurityPolicy},
			},
		},
	}

	binding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: sc.Namespace,
			Labels:    MetadataLabels(sc),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     role.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      InstanceObjectName(sc.Name, ComponentNone),
				Namespace: sc.Namespace,
			},
		},
	}

	return role, binding
}

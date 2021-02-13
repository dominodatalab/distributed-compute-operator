package ray

import (
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

const (
	capabilityAll = "ALL"
	roleKind      = "Role"
)

var (
	policyAPIGroups            = []string{"policy"}
	podSecurityPolicyResources = []string{"podsecuritypolicies"}
	useVerbs                   = []string{"use"}
)

func NewPodSecurityPolicy(rc *dcv1alpha1.RayCluster) (*policyv1beta1.PodSecurityPolicy, *rbacv1.Role, *rbacv1.RoleBinding) {
	psp := &policyv1beta1.PodSecurityPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:   rc.Name,
			Labels: MetadataLabels(rc),
		},
		Spec: policyv1beta1.PodSecurityPolicySpec{
			Privileged:               false,
			AllowPrivilegeEscalation: pointer.BoolPtr(false),
			RequiredDropCapabilities: []corev1.Capability{
				capabilityAll,
			},
			Volumes: []policyv1beta1.FSType{
				policyv1beta1.ConfigMap,
				policyv1beta1.EmptyDir,
				policyv1beta1.Projected,
				policyv1beta1.Secret,
				policyv1beta1.DownwardAPI,
				policyv1beta1.PersistentVolumeClaim,
			},
			HostNetwork: false,
			HostIPC:     false,
			HostPID:     false,
			RunAsUser: policyv1beta1.RunAsUserStrategyOptions{
				Rule: policyv1beta1.RunAsUserStrategyMustRunAsNonRoot,
			},
			SELinux: policyv1beta1.SELinuxStrategyOptions{
				Rule: policyv1beta1.SELinuxStrategyRunAsAny,
			},
			SupplementalGroups: policyv1beta1.SupplementalGroupsStrategyOptions{
				Rule: policyv1beta1.SupplementalGroupsStrategyMustRunAs,
				Ranges: []policyv1beta1.IDRange{
					{
						Min: 1,
						Max: 65535,
					},
				},
			},
			FSGroup: policyv1beta1.FSGroupStrategyOptions{
				Rule: policyv1beta1.FSGroupStrategyMustRunAs,
				Ranges: []policyv1beta1.IDRange{
					{
						Min: 1,
						Max: 65535,
					},
				},
			},
		},
	}

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
				ResourceNames: []string{psp.Name},
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
			Kind:     roleKind,
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

	return psp, role, binding
}

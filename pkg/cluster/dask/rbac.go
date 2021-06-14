package dask

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

var (
	policyAPIGroups            = []string{"policy"}
	podSecurityPolicyResources = []string{"podsecuritypolicies"}
	useVerbs                   = []string{"use"}
)

func RolePodSecurityPolicy() core.OwnedComponent {
	return components.Role(func(obj client.Object) components.RoleDataSource {
		return &pspDS{dc: daskCluster(obj)}
	})
}

func RoleBindingPodSecurityPolicy() core.OwnedComponent {
	return components.RoleBinding(func(obj client.Object) components.RoleBindingDataSource {
		return &pspDS{dc: daskCluster(obj)}
	})
}

type pspDS struct {
	dc *dcv1alpha1.DaskCluster
}

func (s *pspDS) Role() *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: s.objectMeta(),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     policyAPIGroups,
				Resources:     podSecurityPolicyResources,
				Verbs:         useVerbs,
				ResourceNames: []string{s.dc.Spec.PodSecurityPolicy},
			},
		},
	}
}

func (s *pspDS) RoleBinding() *rbacv1.RoleBinding {
	om := s.objectMeta()

	return &rbacv1.RoleBinding{
		ObjectMeta: om,
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     om.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      om.Name,
				Namespace: s.dc.Namespace,
			},
		},
	}
}

func (s *pspDS) Delete() bool {
	return s.dc.Spec.PodSecurityPolicy == ""
}

func (s *pspDS) objectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      meta.InstanceName(s.dc, metadata.ComponentNone),
		Namespace: s.dc.Namespace,
		Labels:    meta.StandardLabels(s.dc),
	}
}

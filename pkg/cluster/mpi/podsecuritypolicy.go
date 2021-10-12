package mpi

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

var (
	policyAPIGroups            = []string{"policy"}
	podSecurityPolicyResources = []string{"podsecuritypolicies"}
	useVerbs                   = []string{"use"}
)

func RolePodSecurityPolicy() core.OwnedComponent {
	return &podSecurityPolicyComponent{kind: &rbacv1.Role{}}
}

func RoleBindingPodSecurityPolicy() core.OwnedComponent {
	return &podSecurityPolicyComponent{kind: &rbacv1.RoleBinding{}}
}

type podSecurityPolicyComponent struct {
	kind client.Object
}

func (c podSecurityPolicyComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPIJob(ctx.Object)

	desc, resource := c.buildResource(cr)

	if cr.Spec.PodSecurityPolicy == "" {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, resource)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, resource)
	if err != nil {
		err = fmt.Errorf("cannot reconcile %s: %w", desc, err)
	}

	return ctrl.Result{}, err
}

func (c podSecurityPolicyComponent) Kind() client.Object {
	return c.kind
}

func (c podSecurityPolicyComponent) buildResource(cr *dcv1alpha1.MPIJob) (string, client.Object) {
	om := metav1.ObjectMeta{
		Name:      meta.InstanceName(cr, metadata.ComponentNone),
		Namespace: cr.Namespace,
		Labels:    meta.StandardLabels(cr),
	}

	switch c.Kind().(type) {
	case *rbacv1.Role:
		return "role", &rbacv1.Role{
			ObjectMeta: om,
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups:     policyAPIGroups,
					Resources:     podSecurityPolicyResources,
					Verbs:         useVerbs,
					ResourceNames: []string{cr.Spec.PodSecurityPolicy},
				},
			},
		}
	case *rbacv1.RoleBinding:
		return "role binding", &rbacv1.RoleBinding{
			ObjectMeta: om,
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "Role",
				Name:     om.Name,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      serviceAccountName(cr),
					Namespace: cr.Namespace,
				},
			},
		}
	}

	panic("you did something dumb")
}

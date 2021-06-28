package components

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

type RoleDataSource interface {
	Role() *rbacv1.Role
	Delete() bool
}

type RoleDataSourceFactory func(client.Object) RoleDataSource

func Role(f RoleDataSourceFactory) core.OwnedComponent {
	return &roleComponent{factory: f}
}

type roleComponent struct {
	factory RoleDataSourceFactory
}

func (c *roleComponent) Kind() client.Object {
	return &rbacv1.Role{}
}

func (c *roleComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)
	role := ds.Role()

	if ds.Delete() {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, role)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, role)
	if err != nil {
		err = fmt.Errorf("cannot reconcile role: %w", err)
	}

	return ctrl.Result{}, err
}

type RoleBindingDataSource interface {
	RoleBinding() *rbacv1.RoleBinding
	Delete() bool
}

type RoleBindingDataSourceFactory func(client.Object) RoleBindingDataSource

func RoleBinding(f RoleBindingDataSourceFactory) core.OwnedComponent {
	return &roleBindingComponent{factory: f}
}

type roleBindingComponent struct {
	factory RoleBindingDataSourceFactory
}

func (c *roleBindingComponent) Kind() client.Object {
	return &rbacv1.RoleBinding{}
}

func (c *roleBindingComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)
	rb := ds.RoleBinding()

	if ds.Delete() {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, rb)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, rb)
	if err != nil {
		err = fmt.Errorf("cannot reconcile role binding: %w", err)
	}

	return ctrl.Result{}, err
}

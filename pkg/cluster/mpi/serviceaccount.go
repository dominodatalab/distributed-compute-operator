package mpi

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func ServiceAccount() core.OwnedComponent {
	return &serviceAccountComponent{}
}

type serviceAccountComponent struct{}

func (c serviceAccountComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)
	conf := cr.Spec.ServiceAccount

	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName(cr),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabels(cr),
		},
		AutomountServiceAccountToken: ptr.To(conf.AutomountServiceAccountToken),
	}

	if conf.Name != "" {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, sa)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, cr, sa)
	if err != nil {
		err = fmt.Errorf("cannot reconcile serviceaccount: %w", err)
	}

	return ctrl.Result{}, err
}

func (c serviceAccountComponent) Kind() client.Object {
	return &corev1.ServiceAccount{}
}

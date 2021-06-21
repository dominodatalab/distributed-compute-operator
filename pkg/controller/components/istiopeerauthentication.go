package components

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources/istio"
)

type IstioPeerAuthenticationDataSource interface {
	PeerAuthInfo() *istio.PeerAuthInfo
	Enabled() bool
	Delete() bool
}

type IstioPeerAuthenticationDataSourceFactory func(client.Object) IstioPeerAuthenticationDataSource

func IstioPeerAuthentication(f IstioPeerAuthenticationDataSourceFactory) core.Component {
	return &istioPeerAuthenticationComponent{factory: f}
}

type istioPeerAuthenticationComponent struct {
	factory IstioPeerAuthenticationDataSourceFactory
}

func (c *istioPeerAuthenticationComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	ds := c.factory(ctx.Object)

	if !ds.Enabled() {
		return ctrl.Result{}, nil
	}

	peerAuth := istio.NewPeerAuthentication(ds.PeerAuthInfo())
	if ds.Delete() {
		return ctrl.Result{}, actions.DeleteIfExists(ctx, peerAuth)
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, ctx.Object, peerAuth)
	if err != nil {
		err = fmt.Errorf("cannot reconcile istio peer authentication: %w", err)
	}

	return ctrl.Result{}, err
}

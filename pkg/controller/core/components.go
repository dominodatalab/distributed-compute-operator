package core

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Component interface {
	Reconcile(*Context) (ctrl.Result, error)
}

type OwnedComponent interface {
	Component
	Kind() client.Object
}

type FinalizerComponent interface {
	Finalize(*Context) (ctrl.Result, bool, error)
}

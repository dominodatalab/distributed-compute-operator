package core

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Component interface {
	Kind() client.Object
	Reconcile(*Context) (ctrl.Result, error)
}

type FinalizerComponent interface {
	Finalize(*Context) (ctrl.Result, bool, error)
}

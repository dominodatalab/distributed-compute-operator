package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

type Builder func(manager ctrl.Manager, webhooksEnabled bool, cfg *Config) error

var BuilderFuncs = []Builder{
	DaskCluster,
	MPICluster,
}

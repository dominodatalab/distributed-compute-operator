package dask

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
)

const (
	ApplicationName                       = "dask"
	ComponentScheduler metadata.Component = "scheduler"
	ComponentWorker    metadata.Component = "worker"
)

var meta = metadata.NewProvider(
	ApplicationName,
	func(obj client.Object) string { return obj.(*dcv1alpha1.DaskCluster).Spec.Image.Tag },
)

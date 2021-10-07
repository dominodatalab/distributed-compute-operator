package mpi

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
)

const (
	ApplicationName                      = "mpi"
	ComponentLauncher metadata.Component = "launcher"
	ComponentWorker   metadata.Component = "worker"
)

var meta = metadata.NewProvider(
	ApplicationName,
	func(obj client.Object) string { return objToMPIJob(obj).Spec.Image.Tag },
	func(obj client.Object) map[string]string { return objToMPIJob(obj).Spec.GlobalLabels },
)

func objToMPIJob(obj client.Object) *dcv1alpha1.MPIJob {
	return obj.(*dcv1alpha1.MPIJob)
}

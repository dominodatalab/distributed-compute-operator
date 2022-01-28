package mpi

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
)

const (
	ApplicationName                     = "mpi"
	RsyncSidecarName                    = "rsync"
	ComponentWorker  metadata.Component = "worker"
	ComponentDriver  metadata.Component = "driver"
)

var meta = metadata.NewProvider(
	ApplicationName,
	func(obj client.Object) string { return objToMPICluster(obj).Spec.Image.Tag },
	func(obj client.Object) map[string]string { return objToMPICluster(obj).Spec.GlobalLabels },
)

func objToMPICluster(obj client.Object) *dcv1alpha1.MPICluster {
	return obj.(*dcv1alpha1.MPICluster)
}

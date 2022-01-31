package mpi

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
)

const (
	// SSH port used by MPI and a name of this port within the service
	sshdPort     = 2222
	sshdPortName = "tcp-ssh"

	// Locations of the mounted files and their modes
	authorizedKeysPath = "/etc/mpi/authorized_keys"
	authorizedKeysMode = 0444 // octal!
	launchScriptPath   = "/opt/domino/bin/mpi-worker-start.sh"
	launchScriptMode   = 0544 // octal!

	// Default parameters of a user account for executing MPI workload.
	defaultUserID    = 12574
	defaultUserName  = "domino"
	defaultGroupID   = 12574
	defaultGroupName = "domino"

	// Configmap key containing the host file
	hostFileName = "hostfile"

	// Period of rerunning resource finalizers
	finalizerRetryPeriod = 1 * time.Second
)

func configMapName(cr client.Object) string {
	return meta.InstanceName(cr, "config")
}

func selectServiceAccount(cr *dcv1alpha1.MPICluster) string {
	if cr.Spec.ServiceAccount.Name != "" {
		return cr.Spec.ServiceAccount.Name
	}

	return serviceAccountName(cr)
}

func serviceAccountName(cr client.Object) string {
	return meta.InstanceName(cr, metadata.ComponentNone)
}

func serviceName(cr client.Object, comp metadata.Component) string {
	return meta.InstanceName(cr, comp)
}

func sshSecretName(cr *dcv1alpha1.MPICluster) string {
	worker := cr.Spec.Worker
	return worker.SharedSSHSecret
}

func workerStatefulSetName(cr client.Object) string {
	return meta.InstanceName(cr, ComponentWorker)
}

package mpi

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
)

const (
	// SSH port used by MPI worker
	sshdPort     = 2222
	sshdPortName = "tcp-ssh"

	// Locations of the mounted files and their modes
	authorizedKeysPath = "/etc/mpi/authorized_keys"
	authorizedKeysMode = 0444 // octal!

	// Location of common Domino utilities
	customUtilPath = "/opt/domino"

	// Default parameters of a user account for executing MPI workload.
	defaultUserID    = 12574
	defaultUserName  = "domino"
	defaultGroupID   = 12574
	defaultGroupName = "domino"

	// SSH ports used by rsync sidecar
	rsyncPort     = 2223
	rsyncPortName = "tcp-rsync"

	// User and group for running the sidecar container;
	// they should match a user provisioned in the sidecar image.
	rsyncUserID  = 12574
	rsyncGroupID = 12574

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

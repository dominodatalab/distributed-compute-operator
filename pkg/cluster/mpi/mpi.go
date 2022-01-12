package mpi

import (
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/cluster/metadata"
)

const (
	sshdPort                  int32 = 2222
	sshdPortName                    = "sshd"
	sshVolumeName                   = "ssh-auth"
	sshRootMountPath                = "/etc/ssh"
	sshAuthorizedKeysFilename       = "authorized_keys"
	sshAuthPublicKey                = "ssh-publickey"

	configVolumeName    = "config"
	configRootMountPath = "/etc/mpi"

	hostFileFilename   = "hostfile"
	sshdConfigFilename = "sshd_config"

	finalizerRetryPeriod = 1 * time.Second
)

var (
	sshAuthorizedKeysPath = filepath.Join(configRootMountPath, sshAuthorizedKeysFilename)
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

func serviceName(cr client.Object) string {
	return meta.InstanceName(cr, ComponentWorker)
}

func sshSecretName(cr *dcv1alpha1.MPICluster) string {
	worker := cr.Spec.Worker
	return worker.SharedSSHSecret
}

func workerStatefulSetName(cr client.Object) string {
	return meta.InstanceName(cr, ComponentWorker)
}

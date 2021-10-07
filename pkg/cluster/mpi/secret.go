package mpi

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
	"github.com/dominodatalab/distributed-compute-operator/pkg/util/ssh"
)

const SSHAuthPublicKey = "ssh-publickey"

func Secret() core.OwnedComponent {
	return &secretComponent{}
}

type secretComponent struct{}

func (c secretComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPIJob(ctx.Object)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sshSecretName(cr),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabels(cr),
		},
		Immutable: pointer.Bool(true),
		Type:      corev1.SecretTypeSSHAuth,
	}

	var existing corev1.Secret
	switch err := ctx.Client.Get(ctx, client.ObjectKeyFromObject(secret), &existing); {
	case err == nil:
		secret.Data = existing.Data
	case apierrors.IsNotFound(err):
		privateKey, publicKey, kerr := ssh.GenerateECCKeyPair()
		if kerr != nil {
			return ctrl.Result{}, fmt.Errorf("cannot generate ECC keypair: %w", kerr)
		}

		secret.Data = map[string][]byte{
			corev1.SSHAuthPrivateKey: privateKey,
			SSHAuthPublicKey:         publicKey,
		}
	default:
		return ctrl.Result{}, err
	}

	err := actions.CreateOwnedResource(ctx, cr, secret)
	if err != nil {
		err = fmt.Errorf("cannot reconcile secret: %w", err)
	}

	return ctrl.Result{}, err
}

func (c secretComponent) Kind() client.Object {
	return &corev1.Secret{}
}

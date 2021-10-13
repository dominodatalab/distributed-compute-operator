package mpi

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

var (
	sshConfigTmpl = `
Host *
    Port %d
    IdentityFile %s
    UserKnownHostsFile /dev/null
    ConnectionAttempts 10
    StrictHostKeyChecking no
`
	sshdConfigTmpl = `
Port %d
AuthorizedKeysFile %s
`
)

func ConfigMap() core.OwnedComponent {
	return &configMapComponent{}
}

type configMapComponent struct{}

func (c configMapComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPIJob(ctx.Object)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName(cr),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabels(cr),
		},
		Data: map[string]string{
			hostfileFilename:   buildHostfile(cr),
			sshConfigFilename:  buildSSHConfig(),
			sshdConfigFilename: buildSSHDConfig(),
		},
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, cr, cm)
	if err != nil {
		err = fmt.Errorf("cannot reconcile configmap: %w", err)
	}

	return ctrl.Result{}, err
}

func (c configMapComponent) Kind() client.Object {
	return &corev1.ConfigMap{}
}

func buildHostfile(cr *dcv1alpha1.MPIJob) string {
	slots := *cr.Spec.Worker.Slots
	svcName := serviceName(cr)
	workerName := workerStatefulSetName(cr)
	workerReplicas := *cr.Spec.Worker.Replicas

	var builder strings.Builder
	for idx := 0; idx < int(workerReplicas); idx++ {
		entry := fmt.Sprintf("%s-%d.%s slots=%d\n", workerName, idx, svcName, slots)
		builder.WriteString(entry)
	}

	return builder.String()
}

func buildSSHConfig() string {
	return fmt.Sprintf(sshConfigTmpl, sshdPort, sshPrivateKeyPath)
}

func buildSSHDConfig() string {
	return fmt.Sprintf(sshdConfigTmpl, sshdPort, sshAuthorizedKeysPath)
}

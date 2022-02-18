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

func ConfigMap() core.OwnedComponent {
	return &configMapComponent{}
}

type configMapComponent struct{}

func (c configMapComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	hostFileConfig := createHostFileConfig(cr)
	err := actions.CreateOrUpdateOwnedResource(ctx, cr, hostFileConfig)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot reconcile hostfile configmap: %w", err)
	}

	keytabConfig := createKeytabConfig(cr)
	if keytabConfig != nil {
		err := actions.CreateOrUpdateOwnedResource(ctx, cr, keytabConfig)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("cannot reconcile keytab configmap: %w", err)
		}
	}

	return ctrl.Result{}, nil
}

func (c configMapComponent) Kind() client.Object {
	return &corev1.ConfigMap{}
}

func createHostFileConfig(cr *dcv1alpha1.MPICluster) *corev1.ConfigMap {
	svcName := serviceName(cr, ComponentWorker)
	workerName := workerStatefulSetName(cr)
	workerReplicas := *cr.Spec.Worker.Replicas

	var hostFileBuilder strings.Builder
	for idx := 0; idx < int(workerReplicas); idx++ {
		entry := fmt.Sprintf("%s-%d.%s\n", workerName, idx, svcName)
		hostFileBuilder.WriteString(entry)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName(cr) + "-" + hostFileName,
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabels(cr),
		},
		Data: map[string]string{
			hostFileName: hostFileBuilder.String(),
		},
	}
}

func createKeytabConfig(cr *dcv1alpha1.MPICluster) *corev1.ConfigMap {
	if cr.Spec.KerberosKeytab == nil {
		return nil
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName(cr) + "-" + keytabName,
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabels(cr),
		},
		BinaryData: map[string][]byte{
			keytabName: cr.Spec.KerberosKeytab.Contents,
		},
	}
}

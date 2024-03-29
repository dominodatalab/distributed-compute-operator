package spark

import (
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// NewFrameworkConfigMap generates a configmap which represents a spark-defaults.conf file out of provided config
func NewFrameworkConfigMap(sc *dcv1alpha1.SparkCluster) *corev1.ConfigMap {
	data := map[string]string{}
	if sc.Spec.Master.DefaultConfiguration != nil {
		data[string(ComponentMaster)] = generateSparkDefaults(sc.Spec.Master.DefaultConfiguration)
	}
	if sc.Spec.Worker.DefaultConfiguration != nil {
		data[string(ComponentWorker)] = generateSparkDefaults(sc.Spec.Worker.DefaultConfiguration)
	}
	if len(data) == 0 {
		return nil
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      FrameworkConfigMapName(sc.Name, ComponentNone),
			Namespace: sc.Namespace,
			Labels:    AddGlobalLabels(MetadataLabels(sc), sc.Spec.GlobalLabels),
		},
		Data: data,
	}
}

// NewKeyTabConfigMap generates a configmap which represents the Kerberos KeyTab configuration out of provided config
func NewKeyTabConfigMap(sc *dcv1alpha1.SparkCluster) *corev1.ConfigMap {
	binaryData := map[string][]byte{}

	if sc.Spec.KerberosKeytab != nil {
		binaryData["keytab"] = sc.Spec.KerberosKeytab.Contents
	}

	if len(binaryData) == 0 {
		return nil
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KeyTabConfigMapName(sc.Name, ComponentNone),
			Namespace: sc.Namespace,
			Labels:    AddGlobalLabels(MetadataLabels(sc), sc.Spec.GlobalLabels),
		},
		BinaryData: binaryData,
	}
}

// looks a little weird because map iteration isn't stable in go, but we want to provide a stable interface
// so we sort the keys and emit a config in sorted order
func generateSparkDefaults(defaults map[string]string) string {
	var keys []string
	for k := range defaults {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	b := strings.Builder{}
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("%s %s\n", k, defaults[k]))
	}
	return b.String()
}

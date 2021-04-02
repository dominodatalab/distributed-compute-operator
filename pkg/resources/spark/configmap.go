package spark

import (
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// NewConfigMap generates a configmap which represents a spark-defaults.conf file out of provided config
func NewConfigMap(sc *dcv1alpha1.SparkCluster) *corev1.ConfigMap {
	data := map[string]string{}
	if sc.Spec.Master.FrameworkConfig != nil {
		data[string(ComponentMaster)] = generateSparkDefaults(sc.Spec.Master.FrameworkConfig.Configs)
	}
	if (sc.Spec.Worker.FrameworkConfig) != nil {
		data[string(ComponentWorker)] = generateSparkDefaults(sc.Spec.Worker.FrameworkConfig.Configs)
	}
	if len(data) == 0 {
		return nil
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      InstanceObjectName(sc.Name, ComponentNone),
			Namespace: sc.Namespace,
			Labels:    MetadataLabels(sc),
		},
		Data: data,
	}
}

// looks a little weird because map iteration isn't stable in go, but we want to provide a stable interface
// so we sort the keys and emit a config in sorted order
func generateSparkDefaults(vals map[string]string) string {
	var keys []string
	for k := range vals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	b := strings.Builder{}
	for _, k := range keys {
		b.WriteString(k + " " + vals[k] + "\n")
	}
	return b.String()
}

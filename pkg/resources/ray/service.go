package ray

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func NewHeadService(rc *dcv1alpha1.RayCluster) *corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name: "redis-primary",
			Port: rc.Spec.HeadPort,
		},
		{
			Name: "object-manager",
			Port: rc.Spec.ObjectManagerPort,
		},
		{
			Name: "node-manager",
			Port: rc.Spec.NodeManagerPort,
		},
	}

	for idx, port := range rc.Spec.RedisShardPorts {
		ports = append(ports, corev1.ServicePort{
			Name: fmt.Sprintf("redis-shard-%d", idx),
			Port: port,
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HeadServiceName(rc.Name),
			Namespace: rc.Namespace,
			Labels:    MetadataLabelsWithComponent(rc, ComponentHead),
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: SelectorLabelsWithComponent(rc, ComponentHead),
		},
	}
}

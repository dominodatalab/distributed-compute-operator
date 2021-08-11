package dask

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/components"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

func ClusterStatusUpdate() core.Component {
	return components.ClusterStatusUpdate(func(obj client.Object) components.ClusterStatusUpdateDataSource {
		return &clusterStatusUpdateDS{dc: daskCluster(obj)}
	})
}

type clusterStatusUpdateDS struct {
	dc *dcv1alpha1.DaskCluster
}

func (c *clusterStatusUpdateDS) ListOpts() []client.ListOption {
	return []client.ListOption{
		client.InNamespace(c.dc.Namespace),
		client.MatchingLabels(meta.StandardLabels(c.dc)),
	}
}

func (c *clusterStatusUpdateDS) StatefulSet() *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.InstanceName(c.dc, ComponentWorker),
			Namespace: c.dc.Namespace,
		},
	}
}

func (c *clusterStatusUpdateDS) ClusterStatusConfig() *dcv1alpha1.ClusterStatusConfig {
	return &c.dc.Status.ClusterStatusConfig
}

func (c *clusterStatusUpdateDS) Image() *dcv1alpha1.OCIImageDefinition {
	return c.dc.Spec.Image
}

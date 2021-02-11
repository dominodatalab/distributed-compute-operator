package ray

import (
	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources"
)

const ApplicationName = "ray"

type Component string

const (
	ComponentHead   Component = "head"
	ComponentWorker Component = "worker"
)

func MetadataLabels(rc *dcv1alpha1.RayCluster) map[string]string {
	return resources.MetadataLabels(ApplicationName, rc.Name, rc.Spec.Image.Tag)
}

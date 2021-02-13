package ray

import (
	"fmt"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources"
)

type Component string

const (
	ComponentHead   Component = "head"
	ComponentWorker Component = "worker"
)

const ApplicationName = "ray"

func MetadataLabels(rc *dcv1alpha1.RayCluster) map[string]string {
	return resources.MetadataLabels(ApplicationName, rc.Name, rc.Spec.Image.Tag)
}

func MetadataLabelsWithComponent(rc *dcv1alpha1.RayCluster, comp Component) map[string]string {
	return resources.MetadataLabelsWithComponent(ApplicationName, rc.Name, rc.Spec.Image.Tag, string(comp))
}

func SelectorLabels(rc *dcv1alpha1.RayCluster) map[string]string {
	return resources.SelectorLabels(ApplicationName, rc.Name)
}

func SelectorLabelsWithComponent(rc *dcv1alpha1.RayCluster, comp Component) map[string]string {
	return resources.SelectorLabelsWithComponent(ApplicationName, rc.Name, string(comp))
}

func HeadServiceName(name string) string {
	return fmt.Sprintf("%s-%s-head", name, ApplicationName)
}

func InstanceObjectName(instance string, comp Component) string {
	return fmt.Sprintf("%s-%s-%s", instance, ApplicationName, comp)
}

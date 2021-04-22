package ray

import (
	"fmt"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/resources"
)

// Component is used to drive Kubernetes object generation for different ray types.
type Component string

const (
	// ComponentNone indicates a generic ray resource.
	ComponentNone Component = "none"
	// ComponentHead indicates a ray head resource.
	ComponentHead Component = "head"
	// ComponentWorker indicates a ray worker resource.
	ComponentWorker Component = "worker"
	// ApplicationName defines the static name used to generate ray object metadata.
	ApplicationName = "ray"
)

// InstanceObjectName returns the name that will be used to create most owned cluster resources.
func InstanceObjectName(instance string, comp Component) string {
	if comp == ComponentNone {
		return fmt.Sprintf("%s-%s", instance, ApplicationName)
	}

	return fmt.Sprintf("%s-%s-%s", instance, ApplicationName, comp)
}

// HeadlessHeadServiceName returns the name of the headless service used to
// register the head ray pod.
func HeadlessHeadServiceName(name string) string {
	return InstanceObjectName(name, ComponentHead)
}

// HeadlessWorkerServiceName returns the name of the headless service used to
// register ray worker pods.
func HeadlessWorkerServiceName(name string) string {
	return InstanceObjectName(name, ComponentWorker)
}

// MetadataLabels returns standard metadata for ray resources.
func MetadataLabels(rc *dcv1alpha1.RayCluster) map[string]string {
	return resources.MetadataLabels(ApplicationName, rc.Name, rc.Spec.Image.Tag)
}

// MetadataLabelsWithComponent returns standard component metadata for ray resources.
func MetadataLabelsWithComponent(rc *dcv1alpha1.RayCluster, comp Component) map[string]string {
	return resources.MetadataLabelsWithComponent(ApplicationName, rc.Name, rc.Spec.Image.Tag, string(comp))
}

// SelectorLabels returns a resource selector clause for ray resources.
func SelectorLabels(rc *dcv1alpha1.RayCluster) map[string]string {
	return resources.SelectorLabels(ApplicationName, rc.Name)
}

// SelectorLabelsWithComponent returns a resource component selector clause for ray resources.
func SelectorLabelsWithComponent(rc *dcv1alpha1.RayCluster, comp Component) map[string]string {
	return resources.SelectorLabelsWithComponent(ApplicationName, rc.Name, string(comp))
}

package metadata

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ApplicationNameLabelKey indicates the name of the application.
	ApplicationNameLabelKey = "app.kubernetes.io/name"
	// ApplicationInstanceLabelKey indicates a unique name identifying the instance of an application.
	ApplicationInstanceLabelKey = "app.kubernetes.io/instance"
	// ApplicationVersionLabelKey indicates the current version of the application.
	ApplicationVersionLabelKey = "app.kubernetes.io/version"
	// ApplicationComponentLabelKey indicates the component within the architecture of an application.
	ApplicationComponentLabelKey = "app.kubernetes.io/component"
	// ApplicationManagedByLabelKey indicates the tool being used to manage the operation of an application.
	ApplicationManagedByLabelKey = "app.kubernetes.io/managed-by"
	// ApplicationManagedByLabelValue is the specific tool being used to manage applications created by this project.
	ApplicationManagedByLabelValue = "distributed-compute-operator"
	// DescriptionAnnotationKey can be used to add extra information to a Kubernetes object via its annotations.
	DescriptionAnnotationKey = "distributed-compute.dominodatalab.com/description"
)

type Component string

const ComponentNone Component = "none"

type versionExtractor func(client.Object) string

type provider struct {
	application string
	version     versionExtractor
}

func NewProvider(name string, fn versionExtractor) *provider {
	return &provider{
		application: name,
		version:     fn,
	}
}

func (p *provider) InstanceName(obj client.Object, comp Component) string {
	if comp == ComponentNone {
		return fmt.Sprintf("%s-%s", obj.GetName(), p.application)
	}

	return fmt.Sprintf("%s-%s-%s", obj.GetName(), p.application, comp)
}

func (p *provider) StandardLabels(obj client.Object) map[string]string {
	return map[string]string{
		ApplicationNameLabelKey:      p.application,
		ApplicationInstanceLabelKey:  obj.GetName(),
		ApplicationVersionLabelKey:   p.version(obj),
		ApplicationManagedByLabelKey: ApplicationManagedByLabelValue,
	}
}

func (p *provider) StandardLabelsWithComponent(obj client.Object, comp Component) map[string]string {
	labels := p.StandardLabels(obj)
	labels[ApplicationComponentLabelKey] = string(comp)

	return labels
}

func (p *provider) MatchLabels(obj client.Object) map[string]string {
	return map[string]string{
		ApplicationNameLabelKey:     p.application,
		ApplicationInstanceLabelKey: obj.GetName(),
	}
}

func (p *provider) MatchLabelsWithComponent(obj client.Object, comp Component) map[string]string {
	labels := p.MatchLabels(obj)
	labels[ApplicationComponentLabelKey] = string(comp)

	return labels
}

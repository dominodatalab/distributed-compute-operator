package resources

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

// MetadataLabels returns a map used to label Kubernetes resources.
func MetadataLabels(name, instance, version string) map[string]string {
	return map[string]string{
		ApplicationNameLabelKey:      name,
		ApplicationInstanceLabelKey:  instance,
		ApplicationVersionLabelKey:   version,
		ApplicationManagedByLabelKey: ApplicationManagedByLabelValue,
	}
}

// MetadataLabelsWithComponent returns a map used to label Kubernetes resources that act as unique components.
func MetadataLabelsWithComponent(name, instance, version, component string) map[string]string {
	labels := MetadataLabels(name, instance, version)
	labels[ApplicationComponentLabelKey] = component

	return labels
}

// SelectorLabels returns a map used to select Kubernetes objects that have
// been labeled with output from MetadataLabels.
func SelectorLabels(name, instance string) map[string]string {
	return map[string]string{
		ApplicationNameLabelKey:     name,
		ApplicationInstanceLabelKey: instance,
	}
}

// SelectorLabelsWithComponent returns a map used to select Kubernetes objects
// that have been labeled with output from MetadataLabelsWithComponent.
func SelectorLabelsWithComponent(name, instance, component string) map[string]string {
	labels := SelectorLabels(name, instance)
	labels[ApplicationComponentLabelKey] = component

	return labels
}

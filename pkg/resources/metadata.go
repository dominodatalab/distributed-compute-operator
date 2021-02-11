package resources

const (
	ApplicationNameLabelKey      = "app.kubernetes.io/name"
	ApplicationInstanceLabelKey  = "app.kubernetes.io/instance"
	ApplicationVersionLabelKey   = "app.kubernetes.io/version"
	ApplicationComponentLabelKey = "app.kubernetes.io/component"
	ApplicationManagedByLabelKey = "app.kubernetes.io/managed-by"

	ManagerLabelValue = "distributed-compute-operator"
)

func MetadataLabels(name, instance, version string) map[string]string {
	return map[string]string{
		ApplicationNameLabelKey:      name,
		ApplicationInstanceLabelKey:  instance,
		ApplicationVersionLabelKey:   version,
		ApplicationManagedByLabelKey: ManagerLabelValue,
	}
}

func MetadataLabelsWithComponent(name, instance, version, component string) map[string]string {
	labels := MetadataLabels(name, instance, version)
	labels[ApplicationComponentLabelKey] = component

	return labels
}

func SelectorLabels(name, instance string) map[string]string {
	return map[string]string{
		ApplicationNameLabelKey:     name,
		ApplicationInstanceLabelKey: instance,
	}
}

func SelectorLabelsWithComponent(name, instance, component string) map[string]string {
	labels := SelectorLabels(name, instance)
	labels[ApplicationComponentLabelKey] = component

	return labels
}

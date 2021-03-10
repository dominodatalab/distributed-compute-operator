package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataLabels(t *testing.T) {
	actual := MetadataLabels("my-app", "inst", "v1.0.0")
	expected := map[string]string{
		"app.kubernetes.io/name":       "my-app",
		"app.kubernetes.io/instance":   "inst",
		"app.kubernetes.io/version":    "v1.0.0",
		"app.kubernetes.io/managed-by": "distributed-compute-operator",
	}

	assert.Equal(t, expected, actual)
}

func TestMetadataLabelsWithComponent(t *testing.T) {
	actual := MetadataLabelsWithComponent("my-app", "inst", "v1.0.0", "comp")
	expected := map[string]string{
		"app.kubernetes.io/name":       "my-app",
		"app.kubernetes.io/instance":   "inst",
		"app.kubernetes.io/version":    "v1.0.0",
		"app.kubernetes.io/managed-by": "distributed-compute-operator",
		"app.kubernetes.io/component":  "comp",
	}

	assert.Equal(t, expected, actual)
}

func TestSelectorLabels(t *testing.T) {
	actual := SelectorLabels("my-app", "inst")
	expected := map[string]string{
		"app.kubernetes.io/name":     "my-app",
		"app.kubernetes.io/instance": "inst",
	}

	assert.Equal(t, expected, actual)
}

func TestSelectorLabelsWithComponent(t *testing.T) {
	actual := SelectorLabelsWithComponent("my-app", "inst", "comp")
	expected := map[string]string{
		"app.kubernetes.io/name":      "my-app",
		"app.kubernetes.io/instance":  "inst",
		"app.kubernetes.io/component": "comp",
	}

	assert.Equal(t, expected, actual)
}

package ray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceObjectName(t *testing.T) {
	t.Run("with_component", func(t *testing.T) {
		comp := Component("test")
		actual := InstanceObjectName("steve-o", comp)
		assert.Equal(t, "steve-o-ray-test", actual)
	})

	t.Run("component_none", func(t *testing.T) {
		actual := InstanceObjectName("steve-o", ComponentNone)
		assert.Equal(t, "steve-o-ray", actual)
	})
}

func TestHeadlessHeadServiceName(t *testing.T) {
	actual := HeadlessHeadServiceName("steve-o")
	assert.Equal(t, "steve-o-ray-head", actual)
}

func TestHeadlessWorkerServiceName(t *testing.T) {
	actual := HeadlessWorkerServiceName("steve-o")
	assert.Equal(t, "steve-o-ray-worker", actual)
}

func TestMetadataLabels(t *testing.T) {
	rc := rayClusterFixture()
	actual := MetadataLabels(rc)

	expected := map[string]string{
		"app.kubernetes.io/name":       "ray",
		"app.kubernetes.io/instance":   "test-id",
		"app.kubernetes.io/version":    "fake-tag",
		"app.kubernetes.io/managed-by": "distributed-compute-operator",
	}
	assert.Equal(t, expected, actual)
}

func TestMetadataLabelsWithComponent(t *testing.T) {
	rc := rayClusterFixture()
	actual := MetadataLabelsWithComponent(rc, "something")

	expected := map[string]string{
		"app.kubernetes.io/name":       "ray",
		"app.kubernetes.io/instance":   "test-id",
		"app.kubernetes.io/version":    "fake-tag",
		"app.kubernetes.io/managed-by": "distributed-compute-operator",
		"app.kubernetes.io/component":  "something",
	}
	assert.Equal(t, expected, actual)
}

func TestSelectorLabels(t *testing.T) {
	rc := rayClusterFixture()
	actual := SelectorLabels(rc)

	expected := map[string]string{
		"app.kubernetes.io/name":     "ray",
		"app.kubernetes.io/instance": "test-id",
	}
	assert.Equal(t, expected, actual)
}

func TestSelectorLabelsWithComponent(t *testing.T) {
	rc := rayClusterFixture()
	actual := SelectorLabelsWithComponent(rc, "something")

	expected := map[string]string{
		"app.kubernetes.io/name":      "ray",
		"app.kubernetes.io/instance":  "test-id",
		"app.kubernetes.io/component": "something",
	}
	assert.Equal(t, expected, actual)
}

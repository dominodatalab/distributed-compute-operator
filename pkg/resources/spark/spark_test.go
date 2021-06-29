package spark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadServiceName(t *testing.T) {
	actual := MasterServiceName("steve-o")
	assert.Equal(t, "steve-o-master", actual)
}

func TestInstanceObjectName(t *testing.T) {
	t.Run("with_component", func(t *testing.T) {
		comp := Component("test")
		actual := InstanceObjectName("steve-o", comp)
		assert.Equal(t, "steve-o-test", actual)
	})

	t.Run("component_none", func(t *testing.T) {
		actual := InstanceObjectName("steve-o", ComponentNone)
		assert.Equal(t, "steve-o", actual)
	})
}

func TestMetadataLabels(t *testing.T) {
	rc := sparkClusterFixture()
	actual := MetadataLabels(rc)

	expected := map[string]string{
		"app.kubernetes.io/name":       "spark",
		"app.kubernetes.io/instance":   "test-id",
		"app.kubernetes.io/version":    "fake-tag",
		"app.kubernetes.io/managed-by": "distributed-compute-operator",
	}
	assert.Equal(t, expected, actual)
}

func TestMetadataLabelsWithComponent(t *testing.T) {
	rc := sparkClusterFixture()
	actual := MetadataLabelsWithComponent(rc, Component("something"))

	expected := map[string]string{
		"app.kubernetes.io/name":       "spark",
		"app.kubernetes.io/instance":   "test-id",
		"app.kubernetes.io/version":    "fake-tag",
		"app.kubernetes.io/managed-by": "distributed-compute-operator",
		"app.kubernetes.io/component":  "something",
	}
	assert.Equal(t, expected, actual)
}

func TestSelectorLabels(t *testing.T) {
	rc := sparkClusterFixture()
	actual := SelectorLabels(rc)

	expected := map[string]string{
		"app.kubernetes.io/name":     "spark",
		"app.kubernetes.io/instance": "test-id",
	}
	assert.Equal(t, expected, actual)
}

func TestSelectorLabelsWithComponent(t *testing.T) {
	rc := sparkClusterFixture()
	actual := SelectorLabelsWithComponent(rc, Component("something"))

	expected := map[string]string{
		"app.kubernetes.io/name":      "spark",
		"app.kubernetes.io/instance":  "test-id",
		"app.kubernetes.io/component": "something",
	}
	assert.Equal(t, expected, actual)
}

func TestFrameworkConfigMapName(t *testing.T) {
	rc := sparkClusterFixture()
	actual := FrameworkConfigMapName(rc.Name, Component("something"))

	expected := "test-id-framework-something"

	assert.Equal(t, expected, actual)
}

func TestKeyTabConfigMapName(t *testing.T) {
	rc := sparkClusterFixture()
	actual := KeyTabConfigMapName(rc.Name, Component("something"))

	expected := "test-id-keytab-something"

	assert.Equal(t, expected, actual)
}

package resources

import (
	"fmt"
	"strconv"

	"github.com/docker/distribution/reference"

	"github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

// IntsToStrings converts an integer slice into a string slice.
func IntsToStrings(is []int32) (ss []string) {
	for _, i := range is {
		ss = append(ss, strconv.Itoa(int(i)))
	}
	return
}

// ParseImageDefinition generates a fully-qualified image reference to an OCI image.
// An error will be returned when the image definition is invalid.
func ParseImageDefinition(def *v1alpha1.OCIImageDefinition) (string, error) {
	ref := def.Repository

	if def.Registry != "" {
		ref = fmt.Sprintf("%s/%s", def.Registry, ref)
	}
	if def.Tag != "" {
		ref = fmt.Sprintf("%s:%s", ref, def.Tag)
	}

	named, err := reference.ParseNormalizedNamed(ref)
	if err != nil {
		return "", fmt.Errorf("invalid OCIImageDefinition: %w", err)
	}
	named = reference.TagNameOnly(named)

	return named.String(), nil
}

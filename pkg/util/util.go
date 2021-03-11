package util

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

// MergeStringMaps merges the src map into the dst.
func MergeStringMaps(src, dst map[string]string) map[string]string {
	for k, v := range src {
		dst[k] = v
	}
	return dst
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

// BoolPtrIsTrue returns true if bool pointer is true. This returns false if
// pointer is false or nil.
func BoolPtrIsTrue(ptr *bool) bool {
	return ptr != nil && *ptr
}

// BoolPtrIsNilOrFalse returns true if bool pointer is nil or false, otherwise
// this returns false.
func BoolPtrIsNilOrFalse(ptr *bool) bool {
	return ptr == nil || !*ptr
}

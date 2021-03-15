package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
)

func TestIntsToStrings(t *testing.T) {
	testcases := []struct {
		in  []int32
		out []string
	}{
		{
			in:  []int32{},
			out: nil,
		},
		{
			in:  []int32{1, 2, 3},
			out: []string{"1", "2", "3"},
		},
		{
			in:  []int32{-5, 0, 1e3},
			out: []string{"-5", "0", "1000"},
		},
	}

	for _, tc := range testcases {
		actual := IntsToStrings(tc.in)
		assert.Equal(t, tc.out, actual)
	}
}

func TestMergeStringMaps(t *testing.T) {
	src := map[string]string{
		"one":    "two",
		"buckle": "my shoe",
	}
	dst := map[string]string{
		"three": "four",
		"knock": "at the door",
	}

	expected := map[string]string{
		"one":    "two",
		"buckle": "my shoe",
		"three":  "four",
		"knock":  "at the door",
	}
	actual := MergeStringMaps(src, dst)

	assert.Equal(t, expected, actual)
}

func TestParseImageDefinition(t *testing.T) {
	testcases := []struct {
		input    *dcv1alpha1.OCIImageDefinition
		expected string
		invalid  bool
	}{
		{
			input: &dcv1alpha1.OCIImageDefinition{
				Registry:   "test-reg:5000",
				Repository: "test-repo",
				Tag:        "test-tag",
			},
			expected: "test-reg:5000/test-repo:test-tag",
		},
		{
			input: &dcv1alpha1.OCIImageDefinition{
				Repository: "test-repo",
				Tag:        "test-tag",
			},
			expected: "docker.io/library/test-repo:test-tag",
		},
		{
			input: &dcv1alpha1.OCIImageDefinition{
				Repository: "test-repo",
			},
			expected: "docker.io/library/test-repo:latest",
		},
		{
			input: &dcv1alpha1.OCIImageDefinition{
				Registry: "test-reg:5000",
				Tag:      "test-tag",
			},
			invalid: true,
		},
		{
			input: &dcv1alpha1.OCIImageDefinition{
				Repository: "!*@~",
			},
			invalid: true,
		},
		{
			input:   &dcv1alpha1.OCIImageDefinition{},
			invalid: true,
		},
	}

	for _, tc := range testcases {
		actual, err := ParseImageDefinition(tc.input)

		if tc.invalid {
			assert.Error(t, err)
			return
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestBoolPtrIsTrue(t *testing.T) {
	tb, fb := true, false

	assert.True(t, BoolPtrIsTrue(&tb))
	assert.False(t, BoolPtrIsTrue(&fb))
	assert.False(t, BoolPtrIsTrue(nil))
}

func TestBoolPtrIsNilOrFalse(t *testing.T) {
	tb, fb := true, false

	assert.True(t, BoolPtrIsNilOrFalse(nil))
	assert.True(t, BoolPtrIsNilOrFalse(&fb))
	assert.False(t, BoolPtrIsNilOrFalse(&tb))
}

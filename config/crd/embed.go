package crd

import (
	"embed"
	"path/filepath"
	"strings"
)

// NOTE: If we start using conversion webhooks in the future and need to
//  "patch" our CRD bases with `kustomize', we can (1) pre-process CRDs during
//  build time and store them in "config/crd/processed", (2) git ignore that
//  directory, (3) and embed that directory instead of "bases".

//go:embed bases/*.yaml
var bases embed.FS

const (
	contentDir        = "bases"
	v1beta1FileSuffix = ".v1beta1.yaml"
)

// Definition represents the metadata and contents of a single custom resource definition.
type Definition struct {
	Filename    string
	Contents    []byte
	BetaVersion bool
}

// ReadAll returns a slice of Definition objects for all of the embedded custom resource definitions.
func ReadAll() (definitions []Definition, err error) {
	files, err := bases.ReadDir(contentDir)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		var contents []byte
		contents, err = bases.ReadFile(filepath.Join(contentDir, f.Name()))
		if err != nil {
			return
		}
		betaVersion := strings.HasSuffix(f.Name(), v1beta1FileSuffix)

		definitions = append(definitions, Definition{
			Filename:    f.Name(),
			Contents:    contents,
			BetaVersion: betaVersion,
		})
	}

	return definitions, nil
}

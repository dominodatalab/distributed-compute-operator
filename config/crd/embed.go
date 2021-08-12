package crd

import (
	"embed"
	"path/filepath"
)

// NOTE: If we start using conversion webhooks in the future and need to
//  "patch" our CRD bases with `kustomize', we can (1) pre-process CRDs during
//  build time and store them in "config/crd/processed", (2) git ignore that
//  directory, (3) and embed that directory instead of "bases".

//go:embed bases/*.yaml
var bases embed.FS

const contentDir = "bases"

// Definition represents the metadata and contents of a single custom resource definition.
type Definition struct {
	Filename string
	Contents []byte
}

// ReadAll returns a slice of custom resource Definition objects.
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

		definitions = append(definitions, Definition{
			Filename: f.Name(),
			Contents: contents,
		})
	}

	return definitions, nil
}

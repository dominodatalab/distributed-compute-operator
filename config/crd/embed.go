package crd

import (
	"embed"
	"path/filepath"
)

// NOTE: ON USING WEBHOOKS
//
// If we start using webhooks in the future and need to "patch" our CRD bases
// with `kustomize', we can (1) pre-process CRDs during build time and store
// them in "config/crd/processed", (2) git ignore that directory, (3) and embed
// that directory instead of "bases".

//go:embed bases/*.yaml
var bases embed.FS

const contentDir = "bases"

// ReadAll returns a slice containing the contents of all base custom resource definitions.
func ReadAll() (definitions [][]byte, err error) {
	files, err := bases.ReadDir(contentDir)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		var content []byte
		content, err = bases.ReadFile(filepath.Join(contentDir, f.Name()))
		if err != nil {
			return
		}

		definitions = append(definitions, content)
	}

	return definitions, nil
}

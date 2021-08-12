package test

import (
	"path/filepath"
	"runtime"
)

// MissingAssetsWarning is a hint as to why an envtest environment will not start.
const MissingAssetsWarning = "Ensure required testing binaries are present by running `make test-assets`"

// KubebuilderBinaryAssetsDir returns a path where control plane binaries required by envtest should be installed.
// TODO: figure out whether to remove this or update it; it no longer works.
func KubebuilderBinaryAssetsDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "..", "testbin", "bin")
}

package test

import (
	"os"
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

// HackBetaCRDs removes v1beta1 CRDs from the provided directory and returns a reset func.
func HackBetaCRDs(path string) (func() error, error) {
	betaFiles, err := filepath.Glob(filepath.Join(path, "*.v1beta1.yaml"))
	if err != nil {
		return nil, err
	}

	dir, err := os.MkdirTemp("", "beta-crds")
	if err != nil {
		return nil, err
	}

	for _, file := range betaFiles {
		if err = os.Rename(file, filepath.Join(dir, filepath.Base(file))); err != nil {
			return nil, err
		}
	}

	fn := func() error {
		for _, file := range betaFiles {
			if err := os.Rename(filepath.Join(dir, filepath.Base(file)), file); err != nil {
				return err
			}
		}

		return os.RemoveAll(dir)
	}

	return fn, nil
}

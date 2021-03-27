#!/usr/bin/env bash
#
# Performs actions required prior to a release.

set -eu

# ensure dependencies are in-sync prior to builds
go mod tidy

# ensure crds are up-to-date
make manifests

# copy crds into a known location for goreleaser
dir=custom-resource-definitions
mkdir -p $dir
cp config/crd/bases/*.yaml $dir

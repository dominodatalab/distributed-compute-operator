#!/usr/bin/env bash

set -euo pipefail

GENERATOR_IMAGE_NAME=ghcr.io/kubernetes-client/java/crd-model-gen
GENERATOR_IMAGE_TAG=v1.0.4

LOCAL_GENERATE_DIR="$(pwd)/client-java"
#LOCAL_MANIFEST_FILE="$LOCAL_GENERATE_DIR/distributed-compute-operator.crds.yaml"

mkdir -p "$LOCAL_GENERATE_DIR"

# generate unified crd manifest
#make -s render > "$LOCAL_MANIFEST_FILE"

# generate classes in separate packages

docker run \
  --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$(pwd)":"$(pwd)" \
  --network host \
  ${GENERATOR_IMAGE_NAME}:${GENERATOR_IMAGE_TAG} \
  /generate.sh \
  -u /Users/sonny/go/src/github.com/dominodatalab/distributed-compute-operator/config/crd/bases/distributed-compute.dominodatalab.com_daskclusters.yaml \
  -n com.dominodatalab.distributed-compute \
  -p com.dominodatalab.distributedcomputeoperator.dask \
  -o "$LOCAL_GENERATE_DIR"

docker run \
  --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$(pwd)":"$(pwd)" \
  --network host \
  ${GENERATOR_IMAGE_NAME}:${GENERATOR_IMAGE_TAG} \
  /generate.sh \
  -u /Users/sonny/go/src/github.com/dominodatalab/distributed-compute-operator/config/crd/bases/distributed-compute.dominodatalab.com_rayclusters.yaml \
  -n com.dominodatalab.distributed-compute \
  -p com.dominodatalab.distributedcomputeoperator.ray \
  -o "$LOCAL_GENERATE_DIR"

docker run \
  --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$(pwd)":"$(pwd)" \
  --network host \
  ${GENERATOR_IMAGE_NAME}:${GENERATOR_IMAGE_TAG} \
  /generate.sh \
  -u /Users/sonny/go/src/github.com/dominodatalab/distributed-compute-operator/config/crd/bases/distributed-compute.dominodatalab.com_sparkclusters.yaml \
  -n com.dominodatalab.distributed-compute \
  -p com.dominodatalab.distributedcomputeoperator.spark \
  -o "$LOCAL_GENERATE_DIR"

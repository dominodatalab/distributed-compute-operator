#!/bin/bash
IMAGE_NAME=${IMAGE_NAME:-"quay.io/domino/distributed-compute-operator"}
IMAGE_TAG_PREFIX=${IMAGE_TAG_PREFIX:-"dev-"}
latest_tag="$IMAGE_TAG_PREFIX$(date +%s)"
image="$IMAGE_NAME:$latest_tag"
make manifests generate docker-build IMG="$image"

declare -r COMPUTE_NAMESPACE=$(kubectl get namespaces -ojson | jq -rc '.items[] | select(.metadata.name | endswith("-compute")) | .metadata.name')

docker push $image

helm upgrade \
  distributed-compute-operator \
  deploy/helm/distributed-compute-operator \
  --install \
  -n $COMPUTE_NAMESPACE \
  --set image.registry="quay.io" \
  --set image.repository="domino/distributed-compute-operator" \
  --set image.tag="$latest_tag" \
  --set config.logDevelopmentMode=true \
  --set istio.enabled=true \
  --set istio.cniPluginInstalled=true \
  --set networkPolicy.enabled=true \

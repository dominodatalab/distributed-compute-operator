#!/usr/bin/env bash
#
# A collection of functions to aid with development efforts.

set -eo pipefail

MINIKUBE_PROFILE=${MINIKUBE_PROFILE:-"distributed-compute-operator"}
MINIKUBE_CPUS=${MINIKUBE_CPUS:-"6"}
MINIKUBE_MEMORY=${MINIKUBE_MEMORY:-"16384"}
MINIKUBE_DISK_SIZE=${MINIKUBE_DISK_SIZE:-"50000mb"}
IMAGE_NAME=${IMAGE_NAME:-"distributed-compute-operator"}
IMAGE_TAG_PREFIX=${IMAGE_TAG_PREFIX:-"dev-"}
ISTIOCTL_VERSION=${ISTIOCTL_VERSION:-1.10.2}

function dco::_info {
  echo -e "\033[0;32m[development-workflow]\033[0m INFO: $*"
}

function dco::_error {
  echo -e "\033[0;31m[development-workflow]\033[0m ERROR: $*"
}

function dco::minikube_setup() {
  if ! minikube profile list 2> /dev/null | grep -q "$MINIKUBE_PROFILE"; then
    dco::_info "Creating minikube cluster"
    minikube start \
      --profile="$MINIKUBE_PROFILE" \
      --cpus="$MINIKUBE_CPUS" \
      --memory="$MINIKUBE_MEMORY" \
      --disk-size="$MINIKUBE_DISK_SIZE" \
      --driver=hyperkit \
      --addons=pod-security-policy \
      --extra-config=apiserver.enable-admission-plugins=PodSecurityPolicy \
      --cni=calico \
      --network-plugin=cni
  elif minikube status --profile="$MINIKUBE_PROFILE" | grep -q 'host: Stopped'; then
    dco::_info "Restarting minikube cluster"
    minikube start --profile="$MINIKUBE_PROFILE"
  else
    dco::_info "Minikube cluster is running"
  fi

  if ! helm repo list | grep -q bitnami; then
    dco::_info "Adding bitnami helm repo"
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo update
  else
    dco::_info "Found bitnami helm repo"
  fi

  if ! helm repo list | grep -q jetstack; then
    dco::_info "Adding jetstack helm repo"
    helm repo add jetstack https://charts.jetstack.io
    helm repo update
  else
    dco::_info "Found jetstack helm repo"
  fi

  if ! helm list --namespace=kube-system | grep -q metrics-server; then
    dco::_info "Creating metrics-server helm release"
    helm install metrics-server bitnami/metrics-server \
    --namespace=kube-system \
    --version=v5.6.0 \
    --set=apiService.create=true \
    --set=extraArgs.kubelet-preferred-address-types=InternalIP,extraArgs.kubelet-insecure-tls=true,extraArgs.metric-resolution=5s \
    --wait
  else
    dco::_info "Found metrics-server helm release"
  fi

  if ! helm list --namespace=cert-manager | grep -q cert-manager; then
    dco::_info "Creating cert-manager helm release"
    helm install cert-manager jetstack/cert-manager \
    --namespace=cert-manager \
    --version=v1.2.0 \
    --set=installCRDs=true \
    --create-namespace \
    --wait
  else
    dco::_info "Found cert-manager helm release"
  fi

  dco::_info "Your development environment is ready to use."
}

function dco::minikube_teardown() {
  dco::_info "Tearing down development environment"
  minikube delete --profile="$MINIKUBE_PROFILE"
  dco::_info "Teardown complete"
}

function dco::docker_build() {
    image="$IMAGE_NAME:$IMAGE_TAG_PREFIX$(date +%s)"

    dco::_info "Building local development image '$image'"
    make manifests generate docker-build IMG="$image"

    dco::_info "Loading image '$image' into Minikube"
    minikube image load "$image" --profile="$MINIKUBE_PROFILE"

    dco::_info "Loading complete"
}

function dco::helm_install() {
  local chart_name ssh_key ip_addr latest_tag

  chart_name=distributed-compute-operator
  ssh_key=$(minikube ssh-key --profile="$MINIKUBE_PROFILE")
  ip_addr=$(minikube ip --profile="$MINIKUBE_PROFILE")
  latest_tag=$(
    ssh  -o LogLevel=ERROR -o StrictHostKeyChecking=no -o GlobalKnownHostsFile=/dev/null -o UserKnownHostsFile=/dev/null -i "$ssh_key" docker@"$ip_addr" \
      "docker image list $IMAGE_NAME:$IMAGE_TAG_PREFIX* --format '{{ .Tag }}|{{ .CreatedAt }}'" | \
        sort -r -t '|' -k 2 | head -n 1 | cut -d '|' -f 1
  )

  dco::_info "Deploying operator image '$IMAGE_NAME:$latest_tag'"

  helm upgrade \
    $chart_name \
    deploy/helm/$chart_name \
    --install \
    --set image.registry="" \
    --set image.repository="$IMAGE_NAME" \
    --set image.tag="$latest_tag" \
    --set config.logDevelopmentMode=true
}

dco::install_istio() {
  local bin=bin/istioctl

  if [[ -f $bin ]] && [[ $($bin version) =~ client\ version:\ $ISTIOCTL_VERSION ]]; then
    dco::_info "Istioctl is present"
  else
    dco::_info "Downloading istioctl"

    local osarch="osx"
    if [[ $(uname -s) == "Linux" ]]; then
      osarch="linux-amd64"
    fi

    curl -sLS "https://github.com/istio/istio/releases/download/$ISTIOCTL_VERSION/istioctl-$ISTIOCTL_VERSION-$osarch.tar.gz" \
      | tar -xz -C ./bin
  fi

  local operator_manifest=istio/operator-minimal.yaml
  if $bin verify-install -f $operator_manifest 1> /dev/null && [[ $($bin version) =~ control\ plane\ version:\ $ISTIOCTL_VERSION ]]; then
    dco::_info "Istio is configured"
  else
    dco::_info "Installing istio version $ISTIOCTL_VERSION"
    $bin install -y -f $operator_manifest
  fi

  dco::_info "Applying a global STRICT mTLS policy"
  kubectl apply -f istio/global-strict-mtls.yaml 1> /dev/null

  dco::_info "Adding default namespace to service mesh"
  kubectl label namespaces default istio-injection=enabled --overwrite 1> /dev/null
}

function dco::display_usage() {
  echo
  echo "Helper script that automates parts of the DCO developer workflow."
  echo
  echo "Usage: $(basename "$0") COMMAND"
  echo
  echo "Commands:"
  echo "  create   Creates Minikube instance configured for DCO development"
  echo "  istio    Deploy Istio service mesh into Minikube"
  echo "  build    Build image locally and load it into Minikube"
  echo "  deploy   Deploy Helm chart into Minikube using latest "
  echo "  teardown Destroy Minikube instance"
  echo "  help     Display usage"

  exit 1
}

function dco::main() {
  local command=$1

  case $command in
    create)
      dco::minikube_setup
      ;;
    istio)
      dco::install_istio
      ;;
    build)
      dco::docker_build
      ;;
    deploy)
      dco::helm_install
      ;;
    teardown)
      dco::minikube_teardown
      ;;
    ""|help)
      dco::display_usage
      ;;
    *)
      _error "Unknown command: $command"
      dco::display_usage
      ;;
  esac
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  if [[ $# -gt 1 ]]; then
    dco::_error "Only 1 command is allowed"
    dco::display_usage
  fi

  dco::main "${1:-""}"
fi

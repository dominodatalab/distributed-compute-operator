#!/usr/bin/env bash
#
# A collection of functions to aid with development efforts.

set -euo pipefail

NAME=distributed-compute-operator

function dco::_info {
  echo -e "\033[0;32m[development-workflow]\033[0m INFO: $*"
}

function dco::_error {
  echo -e "\033[0;31m[development-workflow]\033[0m ERROR: $*"
}

function dco::minikube_setup() {
  if ! minikube profile list 2> /dev/null | grep -q $NAME; then
    dco::_info "Creating minikube cluster"
    minikube start \
      --profile=$NAME \
      --driver=hyperkit \
      --cpus=6 \
      --memory=16384 \
      --disk-size=50000mb \
      --addons=pod-security-policy \
      --extra-config=apiserver.enable-admission-plugins=PodSecurityPolicy \
      --network-plugin=cni \
      --cni=calico
  elif minikube status --profile=$NAME | grep -q 'host: Stopped'; then
    dco::_info "Restarting minikube cluster"
    minikube start --profile=$NAME
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
  minikube delete --profile=$NAME
  dco::_info "Teardown complete"
}

function dco::docker_build() {
    image="distributed-compute-operator:dev-$(date +%s)"

    dco::_info "Building local development image '$image'"
    make docker-build IMG="$image"

    dco::_info "Loading image '$image' into Minikube"
    minikube image load "$image" --profile $NAME

    dco::_info "Loading complete"
}

function dco::helm_install() {
  local ssh_key ip_addr latest_tag

  ssh_key=$(minikube ssh-key --profile $NAME)
  ip_addr=$(minikube ip --profile $NAME)
  latest_tag=$(
    ssh -o StrictHostKeyChecking=no -i "$ssh_key" docker@"$ip_addr" \
      "docker image list distributed-compute-operator:dev-* --format '{{.Tag}}:{{.CreatedAt}}'" \
      | sort -k 2 | tail -n 1 | cut -d ':' -f 1
  )

  dco::_info "Deploying operator image 'distributed-compute-operator:$latest_tag'"

  helm upgrade \
    distributed-compute-operator \
    deploy/helm/distributed-compute-operator \
    --install \
    --set image.repository=distributed-compute-operator \
    --set image.tag="$latest_tag" \
    --set config.logDevelopmentMode=true
}

function dco::display_usage() {
  echo
  echo "Helper script that automates parts of the DCO developer workflow."
  echo
  echo "Usage: $(basename "$0") COMMAND"
  echo
  echo "Commands:"
  echo -e "  create    Creates Minikube instance configured for DCO development"
  echo -e "  build     Build image locally and load it into Minikube"
  echo -e "  deploy    Deploy Helm chart into Minikube using latest "
  echo -e "  teardown  Destroy Minikube instance"
  echo -e "  help      Display usage"

  exit 1
}

function dco::main() {
  local command=$1

  case $command in
    create)
      dco::minikube_setup
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

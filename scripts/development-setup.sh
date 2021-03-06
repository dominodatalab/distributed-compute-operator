#!/usr/bin/env bash

MINIKUBE_PROFILE_NAME=distributed-compute-operator

function _info {
  echo -e "\033[0;32m[development-setup]\033[0m INFO: $*"
}

if ! minikube profile list 2> /dev/null | grep -q $MINIKUBE_PROFILE_NAME; then
  _info "Creating minikube cluster"
  minikube start \
    --profile=$MINIKUBE_PROFILE_NAME \
    --driver=hyperkit \
    --cpus=6 \
    --memory=16384 \
    --disk-size=50000mb \
    --addons=pod-security-policy \
    --extra-config=apiserver.enable-admission-plugins=PodSecurityPolicy \
    --network-plugin=cni \
    --cni=calico
elif minikube status --profile=$MINIKUBE_PROFILE_NAME | grep -q 'host: Stopped'; then
  _info "Restarting minikube cluster"
  minikube start --profile=$MINIKUBE_PROFILE_NAME
else
  _info "Minikube cluster is running"
fi

if ! helm repo list | grep -q bitnami; then
  _info "Adding bitnami helm repo"
  helm repo add bitnami https://charts.bitnami.com/bitnami
  helm repo update
else
  _info "Found bitnami helm repo"
fi

if ! helm repo list | grep -q jetstack; then
  _info "Adding jetstack helm repo"
  helm repo add jetstack https://charts.jetstack.io
  helm repo update
else
  _info "Found jetstack helm repo"
fi

if ! helm list --namespace=kube-system | grep -q metrics-server; then
  _info "Creating metrics-server helm release"
  helm install metrics-server bitnami/metrics-server \
  --namespace=kube-system \
  --version=v5.6.0 \
  --set=apiService.create=true \
  --set=extraArgs.kubelet-preferred-address-types=InternalIP,extraArgs.kubelet-insecure-tls=true,extraArgs.metric-resolution=5s \
  --wait
else
  _info "Found metrics-server helm release"
fi

if ! helm list --namespace=cert-manager | grep -q cert-manager; then
  _info "Creating cert-manager helm release"
  helm install cert-manager jetstack/cert-manager \
  --namespace=cert-manager \
  --version=v1.2.0 \
  --set=installCRDs=true \
  --create-namespace \
  --wait
else
  _info "Found cert-manager helm release"
fi

echo -e "\nYour development environment is ready to use." \
        "\nRun \`minikube delete -p $MINIKUBE_PROFILE_NAME\` to teardown."

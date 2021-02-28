<p align="center">
  <img src="docs/img/logo.png" alt="Logo" />
</p>
<p align="center">
  <a href="https://github.com/dominodatalab/distributed-compute-operator/releases">
    <img src="https://img.shields.io/github/v/release/dominodatalab/distributed-compute-operator?include_prereleases&sort=semver" alt="GitHub release" />
  </a>
  <a href="https://github.com/dominodatalab/distributed-compute-operator/actions?query=workflow%3AGo">
    <img src="https://github.com/dominodatalab/distributed-compute-operator/workflows/Go/badge.svg" alt="Go workflow" />
  </a>
  <a href="https://goreportcard.com/report/github.com/dominodatalab/distributed-compute-operator">
    <img src="https://goreportcard.com/badge/github.com/dominodatalab/distributed-compute-operator" alt="Go report card" />
  </a>
  <a href="https://codecov.io/gh/dominodatalab/distributed-compute-operator">
    <img src="https://codecov.io/gh/dominodatalab/distributed-compute-operator/branch/main/graph/badge.svg?token=RY8FO9ITU6" alt="Codecov" />
  </a>
  <a href="https://pkg.go.dev/mod/github.com/dominodatalab/distributed-compute-operator">
    <img src="https://pkg.go.dev/badge/mod/github.com/dominodatalab/distributed-compute-operator" alt="PkgGoDev" />
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/dominodatalab/distributed-compute-operator?color=informational" alt="License" />
  </a>
</p>

# Distributed Compute Operator

Kubernetes operator providing Ray|Spark|Dask clusters on-demand via [Custom Resource Definitions][custom resources].

## Overview

TODO

## Deployment

TODO

## Development

The following instructions will help you create a local Kubernetes environment
that can be used to test every feature supported by this operator.

1. Install [minikube] and create a new cluster.

    ```shell
    # tested using minikube v1.17.1 and k8s v1.20.2
    $ minikube start \
        --cpus=6 --memory=16384 --driver=hyperkit \
        --extra-config=apiserver.enable-admission-plugins=PodSecurityPolicy \
        --addons=pod-security-policy
    ```

1. Install cert-manager
1. Install metrics-server
1. Launch operator

[custom resources]: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/
[minikube]: https://minikube.sigs.k8s.io/docs/

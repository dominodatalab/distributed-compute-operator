#!/usr/bin/env bash
#
# Packages helm chart and pushes it to a remote registry.

export HELM_EXPERIMENTAL_OCI=1

set -euo pipefail

HELM_BIN=${HELM_BIN:-helm}

function dco::helm::login() {
  local registry="$1"
  local username="$2"
  local password="$3"
  local namespace="$4"

  echo "$password" | $HELM_BIN registry login "$registry" \
    --namespace "$namespace" \
    --username "$username" \
    --password-stdin
}

function dco::helm::push() {
  local version=$1
  local ref=$2

  $HELM_BIN package deploy/helm/distributed-compute-operator \
    --destination chart-archives \
    --app-version "$version" \
    --version "$version"
  $HELM_BIN chart save "chart-archives/distributed-compute-operator-$version.tgz" "$ref"
  $HELM_BIN chart push "$ref"

  rm -rf chart-archives/
}

function dco::helm::main() {
  local command=$1
  shift

  case $command in
    login)
      local host=""
      local username=""
      local password=""
      local namespace=""
      local usage

      usage="usage: $(basename "$0") login -h HOST -u USERNAME -p PASSWORD [-n NAMESPACE]"
      while getopts h:u:p:n: opt; do
        case $opt in
          h)
            host=$OPTARG
            ;;
          u)
            username=$OPTARG
            ;;
          p)
            password=$OPTARG
            ;;
          n)
            namespace=$OPTARG
            ;;
          *)
            echo "$usage"
            exit 1
        esac
      done
      shift $((OPTIND -1))

      if [[ -z $host ]] || [[ -z $username ]] || [[ -z $password ]]; then
        echo "$usage"
        exit 1
      fi

      dco::helm::login "$host" "$username" "$password" "$namespace"
      ;;
    push)
      local version=""
      local ref=""
      local usage

      usage="usage: $(basename "$0") push -v VERSION -r REF"
      while getopts v:r: opt; do
        case $opt in
          v)
            version=$OPTARG
            ;;
          r)
            ref=$OPTARG
            ;;
          *)
            echo "$usage"
            exit 1
        esac
      done
      shift $((OPTIND -1))

      if [[ -z $version ]] || [[ -z $ref ]]; then
        echo "$usage"
        exit 1
      fi

      dco::helm::push "$version" "$ref"
      ;;
    ""|help)
      echo
      echo "Usage: $(basename "$0") COMMAND ARGS"
      echo
      echo "Commands:"
      echo "  login Authenticate with remote registry"
      echo "  push  Build and upload chart to a remote registry"
      echo "  help  Display usage"
      exit 1
  esac
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  dco::helm::main "${@:-""}"
fi

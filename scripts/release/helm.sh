#!/usr/bin/env bash
#
# Functions used to log into helm registries, and package/push project chart.

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
  local registry=$1
  local version=$2
  local semantic_version
  local chart_path

  if [[ $version =~ ^(pr-[[:digit:]]+|main)$ ]]; then
    semantic_version="0.0.0-$version"
  else
    semantic_version=$version
  fi

  $HELM_BIN package deploy/helm/distributed-compute-operator \
    --destination chart-archives \
    --app-version "$version" \
    --version "$semantic_version"

  chart_path="chart-archives/distributed-compute-operator-$semantic_version.tgz"

  $HELM_BIN push "$chart_path" oci://"$registry"

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
      local registry=""
      local version=""
      local usage

      usage="usage: $(basename "$0") push -r REGISTRY -v VERSION"
      while getopts r:v: opt; do
        case $opt in
          r)
            registry=$OPTARG
            ;;
          v)
            version=$OPTARG
            ;;
          *)
            echo "$usage"
            exit 1
        esac
      done
      shift $((OPTIND -1))

      if [[ -z $registry ]] || [[ -z $version ]]; then
        echo "$usage"
        exit 1
      fi

      dco::helm::push "$registry" "$version"
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

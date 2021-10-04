# HACK: Fixes "test" target when on a machine where "/bin/sh" is not bash.
# We can remove this once the controller-runtime setup-envtest.sh script is
# updated so that it doesn't require sourcing.
SHELL := /bin/bash
.SHELLFLAGS := -e -o pipefail -c

# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/dominodatalab/distributed-compute-operator:latest
# Produce CRDs that work with Kubernetes 1.16+ and supports defaulting, api
# version conversion, and field pruning.
CRD_OPTIONS ?= "crd:crdVersions=v1,maxDescLen=0"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

fmt: goimports ## Run formatter against code.
	@$(GOIMPORTS) -d -w -local github.com/dominodatalab/distributed-compute-operator \
		$(shell find . -type f -name '*.go' -not -iname "zz_generated.*" -not -path "./vendor/*")

lint: golangci-lint ## Run linters against code.
	$(GOLANGCI_LINT) run

ENVTEST_VERSION = 1.20.x!
ENVTEST_ASSETS_DIR = $(shell pwd)/testbin
test: setup-envtest manifests generate fmt ## Run full test suite.
	$(shell eval "$(SETUP_ENVTEST) --bin-dir $(ENVTEST_ASSETS_DIR) use --print env $(ENVTEST_VERSION)"); go test ./... -race -covermode atomic -coverprofile cover.out

##@ Build

build: generate fmt ## Build manager binary.
	go build -o bin/manager main.go

run: manifests generate fmt ## Run a controller from your host.
	go run ./main.go start

docker-build: ## Build docker image with the manager.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Deployment

render: manifests kustomize ## Maybe keep this; rendering is preferable but we need to figure out how to separate CRDs
	$(KUSTOMIZE) build config/crd

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -

##@ Tools

SETUP_ENVTEST = $(shell pwd)/bin/setup-envtest
setup-envtest: ## Download setup-envtest locally if necessary.
	$(call go-get-tool,$(SETUP_ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

GOIMPORTS = $(shell pwd)/bin/goimports
goimports: ## Download goimports locally if necessary.
	$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	@[ -f $(GOLANGCI_LINT) ] || { \
		set -e ;\
		echo "Installing golangci-lint" ;\
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_DIR)/bin v1.38.0 ;\
	}

HELM = $(shell pwd)/bin/helm
helm: ## Download helm locally if necessary.
	@[ -f $(HELM) ] || { \
		set -e ;\
		echo "Installing helm" ;\
		mkdir -p $(PROJECT_DIR)/bin ;\
		export HELM_INSTALL_DIR=$(PROJECT_DIR)/bin ;\
		curl -sSfL https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | $(SHELL) -s -- --no-sudo --version v3.6.3 ;\
	}

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

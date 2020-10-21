SHELL := bash

# Directory, where all required tools are located (absolute path required)
TOOLS_DIR ?= $(shell cd tools && pwd)

VERSION   ?= 0.0.1-local
KUBE_NAMESPACE ?= platform-monoskope-monoskope

export 

clean: go-clean kind-clean helm-clean tools-clean

# go

go-%:
	@$(MAKE) -f go.mk $*

# helm

HELM_PATH_MONOSKOPE 		?= build/package/helm/monoskope
HELM_VALUES_FILE_MONOSKOPE  ?= examples/00-monoskope-dev-values.yaml

helm-%:
	@$(MAKE) -f helm.mk $*

# kind

kind-%:
	@$(MAKE) -f kind.mk $*

# docs

diagrams:
	$(SHELL) ./build/ci/gen_charts.sh

# Phony target to install all required tools into ${TOOLS_DIR}
tools: kind-get go-ginkgo-get go-golangci-lint-get

tools-clean: kind-clean go-ginkgo-clean go-golangci-lint-clean
SHELL := bash

# Directory, where all required tools are located (absolute path required)
TOOLS_DIR ?= $(shell cd tools && pwd)

VERSION   ?= 0.0.1-local
KUBE_NAMESPACE ?= platform-monoskope-monoskope

export 

# go

go-%:
	@$(MAKE) -f go.mk $*

# helm

helm-%:
	@$(MAKE) -f helm.mk $*

# kind

kind-%:
	@$(MAKE) -f kind.mk $*

# Phony target to install all required tools into ${TOOLS_DIR}
tools: kind-get go-ginkgo-get go-golangci-lint-get

tools-clean: kind-clean go-ginkgo-clean go-golangci-lint-clean
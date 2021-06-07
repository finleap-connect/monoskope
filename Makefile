SHELL := bash

# Directory, where all required tools are located (absolute path required)
TOOLS_DIR ?= $(shell cd tools && pwd)

VERSION   ?= 0.0.1-local
KUBE_NAMESPACE ?= platform-monoskope-monoskope

LATEST_REV=$(shell git rev-list --tags --max-count=1)
LATEST_TAG=$(shell git describe --tags $(LATEST_REV))

export 

clean: go-clean helm-clean tools-clean

# go
include go.mk

# helm

HELM_PATH 		            ?= build/package/helm
HELM_VALUES_FILE            ?= examples/00-monoskope-dev-values.yaml

helm-%:
	@$(MAKE) -f helm.mk $*

# docs

diagrams:
	$(SHELL) ./build/ci/gen_charts.sh

# Phony target to install all required tools into ${TOOLS_DIR}
tools: go-ginkgo-get go-golangci-lint-get

tools-clean: go-ginkgo-clean go-golangci-lint-clean

get-latest:
	@echo LATEST_REV: $(LATEST_REV)
	@echo LATEST_TAG: $(LATEST_TAG)

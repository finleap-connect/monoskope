SHELL := bash

# Directory, where all required tools are located (absolute path required)
TOOLS_DIR ?= $(shell cd tools && pwd)
HACK_DIR ?= $(shell cd hack && pwd)

VERSION   ?= 0.0.1-local
COMMIT     	   := $(shell git rev-parse --short HEAD)
LATEST_TAG=$(shell git describe --tags $(LATEST_REV))

export

clean: go-clean helm-clean tools-clean ## clean up everything

# go
include go.mk

# helm
HELM_PATH 		            ?= build/package/helm
HELM_VALUES_FILE            ?= examples/00-monoskope-dev-values.yaml
include helm.mk

# tools
tools: go-tools ## Phony target to install all required tools into ${TOOLS_DIR}
tools-clean: go-tools-clean ## Phony target to clean all required tools

commit-hash: ## Echos the current commit hash
	@echo $(COMMIT)

latest-tag: ## Echos the latest tag
	@echo $(LATEST_TAG)

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

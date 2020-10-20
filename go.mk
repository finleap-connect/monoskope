GO             ?= go

GINKGO         ?= $(TOOLS_DIR)/ginkgo
GINKO_VERSION  ?= v1.12.0

LINTER 	   	   ?= $(TOOLS_DIR)/golangci-lint
LINTER_VERSION ?= v1.25.0

COMMIT     	   := $(shell git rev-parse --short HEAD)
LDFLAGS    	   += -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)"
BUILDFLAGS 	   += -installsuffix cgo --tags release
PROTOC     	   ?= protoc

VERSION    	   ?= 0.0.1-dev

BUILD_PATH ?= $(shell pwd)

define go-run
	$(GO) run $(LDFLAGS) cmd/$(1)/*.go $(ARGS)
endef

.PHONY: lint prepare fmt vet

prepare:
	$(GO) mod download

lint: golangci-lint-get
	$(GO) mod verify
	$(LINTER) run -v --no-config --deadline=5m

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

run-%:
	$(call go-run,$*)

ginkgo-get:
	$(shell $(TOOLS_DIR)/goget-wrapper github.com/onsi/ginkgo/ginkgo@$(GINKO_VERSION))

golangci-lint-get:
	$(shell curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR) $(LINTER_VERSION))

ginkgo-clean:
	rm -Rf $(TOOLS_DIR)/ginkgo

golangci-lint-clean:
	rm -Rf $(TOOLS_DIR)/golangci-lint

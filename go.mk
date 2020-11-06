BUILD_PATH ?= $(shell pwd)
GO_MODULE ?= gitlab.figo.systems/platform/monoskope/monoskope

GO             ?= go

GINKGO         ?= $(TOOLS_DIR)/ginkgo
GINKO_VERSION  ?= v1.12.0

LINTER 	   	   ?= $(TOOLS_DIR)/golangci-lint
LINTER_VERSION ?= v1.25.0

COMMIT     	   := $(shell git rev-parse --short HEAD)
LDFLAGS    	   += -X=$(GO_MODULE)/internal/metadata.Version=$(VERSION) -X=$(GO_MODULE)/internal/metadata.Commit=$(COMMIT)
BUILDFLAGS 	   += -installsuffix cgo --tags release
PROTOC     	   ?= protoc

CMD_MONOCTL = $(BUILD_PATH)/monoctl
CMD_MONOCTL_SRC = cmd/monoctl/*.go

CMD_GATEWAY = $(BUILD_PATH)/gateway
CMD_GATEWAY_SRC = cmd/gateway/*.go

define go-run
	$(GO) run -ldflags "$(LDFLAGS)" cmd/$(1)/*.go $(ARGS)
endef

.PHONY: lint mod fmt vet test clean

mod:
	$(GO) mod download
	$(GO) mod verify

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint:
	$(LINTER) run -v --no-config --deadline=5m

run-%:
	$(call go-run,$*)

test-kind: 
	$(GINKGO) -r -v -cover internal -- --with-kind --helm-chart-path "$(BUILD_PATH)/$(HELM_PATH_MONOSKOPE)" --helm-chart-values "$(BUILD_PATH)/$(HELM_VALUES_FILE_MONOSKOPE)"

test:
	$(GINKGO) -r -v -cover pkg/gateway -- --dex-conf-path "$(BUILD_PATH)/config/dex"
	$(GINKGO) -r -v -cover pkg/monoctl
	$(GINKGO) -r -v -cover pkg/util

ginkgo-get:
	$(shell $(TOOLS_DIR)/goget-wrapper github.com/onsi/ginkgo/ginkgo@$(GINKO_VERSION))

golangci-lint-get:
	$(shell curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR) $(LINTER_VERSION))

ginkgo-clean:
	rm -Rf $(TOOLS_DIR)/ginkgo

golangci-lint-clean:
	rm -Rf $(TOOLS_DIR)/golangci-lint

clean: ginkgo-clean golangci-lint-clean

protobuf:
	find ./api -name '*.go' -exec rm {} \;
	find ./api -name '*.proto' -exec $(PROTOC) --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. {} \;

$(CMD_MONOCTL):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_MONOCTL) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/pkg/logger.logMode=noop" $(CMD_MONOCTL_SRC)

$(CMD_GATEWAY):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_GATEWAY) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_GATEWAY_SRC)

build-clean: 
	rm -Rf $(CMD_MONOCTL)
	rm -Rf $(CMD_GATEWAY)

build-monoctl: $(CMD_MONOCTL)

build-gateway: $(CMD_GATEWAY)

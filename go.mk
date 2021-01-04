BUILD_PATH ?= $(shell pwd)
GO_MODULE ?= gitlab.figo.systems/platform/monoskope/monoskope

GO             ?= go

GINKGO         ?= $(TOOLS_DIR)/ginkgo
GINKO_VERSION  ?= v1.14.2

LINTER 	   	   ?= $(TOOLS_DIR)/golangci-lint
LINTER_VERSION ?= v1.33.0

COMMIT     	   := $(shell git rev-parse --short HEAD)
LDFLAGS    	   += -X=$(GO_MODULE)/internal/version.Version=$(VERSION) -X=$(GO_MODULE)/internal/version.Commit=$(COMMIT)
BUILDFLAGS 	   += -installsuffix cgo --tags release
PROTOC     	   ?= protoc

CMD_MONOCTL_LINUX = $(BUILD_PATH)/monoctl-linux-amd64
CMD_MONOCTL_OSX = $(BUILD_PATH)/monoctl-osx-amd64
CMD_MONOCTL_WIN = $(BUILD_PATH)/monoctl-win-amd64
CMD_MONOCTL_SRC = cmd/monoctl/*.go

CMD_GATEWAY = $(BUILD_PATH)/gateway
CMD_GATEWAY_SRC = cmd/gateway/*.go

CMD_EVENTSTORE = $(BUILD_PATH)/eventstore
CMD_EVENTSTORE_SRC = cmd/eventstore/*.go

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

test:
	@find . -name '*.coverprofile' -exec rm {} \;
	$(GINKGO) -r -v -cover pkg/*
	$(GINKGO) -r -v -cover internal/gateway -- --dex-conf-path "$(BUILD_PATH)/config/dex"
	$(GINKGO) -r -v -cover internal/monoctl
	$(GINKGO) -r -v -cover internal/eventstore
	@echo "mode: set" > ./monoskope.coverprofile
	@find ./pkg -name "*.coverprofile" -exec cat {} \; | grep -v mode: | sort -r >> ./monoskope.coverprofile   
	@find ./pkg -name '*.coverprofile' -exec rm {} \;
	@find ./internal -name "*.coverprofile" -exec cat {} \; | grep -v mode: | sort -r >> ./monoskope.coverprofile   
	@find ./internal -name '*.coverprofile' -exec rm {} \;

coverage:
	@find . -name '*.coverprofile' -exec go tool cover -func {} \;

loc:
	@gocloc .

ginkgo-get:
	$(shell $(TOOLS_DIR)/goget-wrapper github.com/onsi/ginkgo/ginkgo@$(GINKO_VERSION))

golangci-lint-get:
	$(shell curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR) $(LINTER_VERSION))

ginkgo-clean:
	rm -Rf $(TOOLS_DIR)/ginkgo

golangci-lint-clean:
	rm -Rf $(TOOLS_DIR)/golangci-lint

tools: golangci-lint-get ginkgo-get

clean: ginkgo-clean golangci-lint-clean build-clean
	rm -Rf reports/
	find . -name '*.coverprofile' -exec rm {} \;

protobuf:
	rm -rf $(BUILD_PATH)/pkg/api
	mkdir -p $(BUILD_PATH)/pkg/api
	find ./api -name '*.proto' -exec $(PROTOC) --go_opt=module=gitlab.figo.systems/platform/monoskope/monoskope --go-grpc_opt=module=gitlab.figo.systems/platform/monoskope/monoskope --go_out=. --go-grpc_out=. {} \;

$(CMD_GATEWAY):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_GATEWAY) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_GATEWAY_SRC)

$(CMD_EVENTSTORE):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_EVENTSTORE) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_EVENTSTORE_SRC)

$(CMD_MONOCTL_LINUX):
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_LINUX) $(CMD_MONOCTL_SRC)

$(CMD_MONOCTL_OSX):
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_OSX) $(CMD_MONOCTL_SRC)

$(CMD_MONOCTL_WIN):
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_WIN) $(CMD_MONOCTL_SRC)

build-clean: 
	rm -Rf $(CMD_GATEWAY)
	rm -Rf $(CMD_MONOCTL_LINUX)
	rm -Rf $(CMD_MONOCTL_OSX)
	rm -Rf $(CMD_MONOCTL_WIN)

build-monoctl: $(CMD_MONOCTL_LINUX) $(CMD_MONOCTL_OSX) $(CMD_MONOCTL_WIN)

build-gateway: $(CMD_GATEWAY)

build-eventstore: $(CMD_EVENTSTORE)

push-monoctl:
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_LINUX) "https://artifactory.figo.systems/artifactory/binaries/linux/monoctl-$(VERSION)"
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_OSX) "https://artifactory.figo.systems/artifactory/binaries/osx/monoctl-$(VERSION)"
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_WIN) "https://artifactory.figo.systems/artifactory/binaries/win/monoctl-$(VERSION)"

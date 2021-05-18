BUILD_PATH ?= $(shell pwd)
GO_MODULE ?= gitlab.figo.systems/platform/monoskope/monoskope

GO             ?= go

GINKGO         ?= $(TOOLS_DIR)/ginkgo
GINKO_VERSION  ?= v1.15.2

LINTER 	   	   ?= $(TOOLS_DIR)/golangci-lint
LINTER_VERSION ?= v1.39.0

PROTOC 	   	           ?= $(TOOLS_DIR)/protoc
PROTOC_IMPORTS_DIR          ?= $(BUILD_PATH)/include
PROTOC_VERSION             ?= 3.17.0
PROTOC_GEN_GO_VERSION      ?= v1.26
PROTOC_GEN_GO_GRPC_VERSION ?= v1.1

CURL          ?= curl

COMMIT     	   := $(shell git rev-parse --short HEAD)
LDFLAGS    	   += -X=$(GO_MODULE)/internal/version.Version=$(VERSION) -X=$(GO_MODULE)/internal/version.Commit=$(COMMIT)
BUILDFLAGS 	   += -installsuffix cgo --tags release

uname_S := $(shell sh -c 'uname -s 2>/dev/null || echo not')

ifeq ($(uname_S),Linux)
ARCH = linux-x86_64
endif
ifeq ($(uname_S),Darwin)
ARCH = osx-x86_64
endif

CMD_GATEWAY = $(BUILD_PATH)/gateway
CMD_GATEWAY_SRC = cmd/gateway/*.go

CMD_EVENTSTORE = $(BUILD_PATH)/eventstore
CMD_EVENTSTORE_SRC = cmd/eventstore/*.go

CMD_COMMANDHANDLER = $(BUILD_PATH)/commandhandler
CMD_COMMANDHANDLER_SRC = cmd/commandhandler/*.go

CMD_QUERYHANDLER = $(BUILD_PATH)/queryhandler
CMD_QUERYHANDLER_SRC = cmd/queryhandler/*.go

export DEX_CONFIG = $(BUILD_PATH)/config/dex
export M8_OPERATION_MODE = development

define go-run
	$(GO) run -ldflags "$(LDFLAGS)" cmd/$(1)/*.go $(ARGS)
endef

.PHONY: lint mod fmt vet test clean report

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

report:
	@echo
	@M8_OPERATION_MODE=cmdline $(GO) run -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/pkg/logger.logMode=noop" cmd/commandhandler/*.go report commands $(ARGS)
	@echo
	@M8_OPERATION_MODE=cmdline $(GO) run -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/pkg/logger.logMode=noop" cmd/commandhandler/*.go report permissions $(ARGS)
	@echo

test: 
	@find . -name '*.coverprofile' -exec rm {} \;
	$(GINKGO) -r -v -cover *
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

protoc-get:
	$(CURL) -LO "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(ARCH).zip"
	unzip protoc-$(PROTOC_VERSION)-$(ARCH).zip -d $(TOOLS_DIR)/.protoc-unpack
	mv $(TOOLS_DIR)/.protoc-unpack/bin/protoc $(TOOLS_DIR)/protoc
	mkdir -p $(PROTOC_IMPORTS_DIR)/
	cp -a $(TOOLS_DIR)/.protoc-unpack/include/* $(PROTOC_IMPORTS_DIR)/
	rm -rf $(TOOLS_DIR)/.protoc-unpack/
	$(shell $(TOOLS_DIR)/goget-wrapper google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION))
	$(shell $(TOOLS_DIR)/goget-wrapper google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION))

protobuf:
	rm -rf $(BUILD_PATH)/pkg/api
	mkdir -p $(BUILD_PATH)/pkg/api
	# generates server part
	export PATH="$$PATH:$(TOOLS_DIR)" ; find ./api -name '*.proto' -exec $(PROTOC) -I. -I$(PROTOC_IMPORTS_DIR) --go-grpc_opt=module=gitlab.figo.systems/platform/monoskope/monoskope --go-grpc_out=. {} \;
	# generates client part
	export PATH="$$PATH:$(TOOLS_DIR)" ; find ./api -name '*.proto' -exec $(PROTOC) -I. -I$(PROTOC_IMPORTS_DIR) --go_opt=module=gitlab.figo.systems/platform/monoskope/monoskope --go_out=. {} \;
	cp -a gitlab.figo.systems/platform/monoskope/monoskope/pkg/api pkg/
	rm -rf gitlab.figo.systems/platform/monoskope/monoskope/pkg/api

$(CMD_GATEWAY):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_GATEWAY) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_GATEWAY_SRC)

$(CMD_EVENTSTORE):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_EVENTSTORE) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_EVENTSTORE_SRC)

$(CMD_COMMANDHANDLER):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_COMMANDHANDLER) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_COMMANDHANDLER_SRC)

$(CMD_QUERYHANDLER):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_QUERYHANDLER) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_QUERYHANDLER_SRC)

build-clean: 
	rm -Rf $(CMD_GATEWAY)
	rm -Rf $(CMD_EVENTSTORE)
	rm -Rf $(CMD_COMMANDHANDLER)

build-gateway: $(CMD_GATEWAY)

build-eventstore: $(CMD_EVENTSTORE)

build-commandhandler: $(CMD_COMMANDHANDLER)

build-queryhandler: $(CMD_QUERYHANDLER)

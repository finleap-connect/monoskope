BUILD_PATH ?= $(shell pwd)
GO_MODULE ?= gitlab.figo.systems/platform/monoskope/monoskope

GO             ?= go

GOGET          ?= $(HACK_DIR)/goget-wrapper

GINKGO         ?= $(TOOLS_DIR)/ginkgo
GINKO_VERSION  ?= v1.15.2

LINTER 	   	   ?= $(TOOLS_DIR)/golangci-lint
LINTER_VERSION ?= v1.39.0

MOCKGEN         ?= $(TOOLS_DIR)/mockgen
GOMOCK_VERSION  ?= v1.5.0

PROTOC 	   	           	   ?= $(TOOLS_DIR)/protoc
PROTOC_IMPORTS_DIR         ?= $(BUILD_PATH)/include
PROTOC_VERSION             ?= 3.17.0
PROTOC_GEN_GO_VERSION      ?= v1.26.0
PROTOC_GEN_GO_GRPC_VERSION ?= v1.1.0
PROTO_FILES                != find api -name "*.proto"

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

CMD_CLBOREACTOR = $(BUILD_PATH)/clboreactor
CMD_CLBOREACTOR_SRC = cmd/clusterbootstrapreactor/*.go

export DEX_CONFIG = $(BUILD_PATH)/config/dex
export M8_OPERATION_MODE = development

define go-run
	$(GO) run -ldflags "$(LDFLAGS)" cmd/$(1)/*.go $(ARGS)
endef

.PHONY: go-lint go-mod go-fmt go-vet go-test go-clean go-report

go-mod:
	$(GO) mod download
	$(GO) mod verify

go-fmt:
	$(GO) fmt ./...

go-vet:
	$(GO) vet ./...

go-lint:
	$(LINTER) run -v --no-config --deadline=5m

go-run-%:
	$(call go-run,$*)

go-report:
	@echo
	@M8_OPERATION_MODE=cmdline $(GO) run -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/pkg/logger.logMode=noop" cmd/commandhandler/*.go report commands $(ARGS)
	@echo
	@M8_OPERATION_MODE=cmdline $(GO) run -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/pkg/logger.logMode=noop" cmd/commandhandler/*.go report permissions $(ARGS)
	@echo

.protobuf-deps: $(PROTO_FILES)
	for file in $$(find pkg/api/ -name "*.pb.go") ; do source=$$(awk '/^\/\/ source:/ { print $$3 }' $$file) ; echo "$$file: $$source"; done >.protobuf-deps
	echo -n "GENERATED_GO_FILES := " >>.protobuf-deps
	for file in $$(find pkg/api/ -name "*.pb.go") ; do echo -n " $$file"; done >>.protobuf-deps
	echo >>.protobuf-deps

include .protobuf-deps

go-test: $(GENERATED_GO_FILES) go-rebuild-mocks
	@find . -name '*.coverprofile' -exec rm {} \;
	$(GINKGO) -r -v -cover *
	@echo "mode: set" > ./monoskope.coverprofile
	@find ./pkg -name "*.coverprofile" -exec cat {} \; | grep -v mode: | sort -r >> ./monoskope.coverprofile   
	@find ./pkg -name '*.coverprofile' -exec rm {} \;
	@find ./internal -name "*.coverprofile" -exec cat {} \; | grep -v mode: | sort -r >> ./monoskope.coverprofile   
	@find ./internal -name '*.coverprofile' -exec rm {} \;

go-coverage:
	@find . -name '*.coverprofile' -exec go tool cover -func {} \;

go-loc:
	@gocloc .

ginkgo-get:
	$(shell $(GOGET) github.com/onsi/ginkgo/ginkgo@$(GINKO_VERSION))

golangci-lint-get:
	$(shell $(HACK_DIR)/golangci-lint.sh -b $(TOOLS_DIR) $(LINTER_VERSION))

gomock-get:
	$(shell $(GOGET) github.com/golang/mock/mockgen@$(GOMOCK_VERSION))

protoc-get:
	$(CURL) -LO "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(ARCH).zip"
	unzip protoc-$(PROTOC_VERSION)-$(ARCH).zip -d $(TOOLS_DIR)/.protoc-unpack
	mv $(TOOLS_DIR)/.protoc-unpack/bin/protoc $(TOOLS_DIR)/protoc
	mkdir -p $(PROTOC_IMPORTS_DIR)/
	cp -a $(TOOLS_DIR)/.protoc-unpack/include/* $(PROTOC_IMPORTS_DIR)/
	rm -rf $(TOOLS_DIR)/.protoc-unpack/ protoc-$(PROTOC_VERSION)-$(ARCH).zip
	$(shell $(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION))
	$(shell $(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION))

go-tools: golangci-lint-get ginkgo-get protoc-get gomock-get

go-clean: go-build-clean
	rm  .protobuf-deps
	rm -Rf reports/
	rm -Rf $(TOOLS_DIR)/
	mkdir $(TOOLS_DIR)
	find . -name '*.coverprofile' -exec rm {} \;

%.pb.go: .protobuf-deps
	rm -rf $(BUILD_PATH)/pkg/api
	mkdir -p $(BUILD_PATH)/pkg/api
	# generates server part
	export PATH="$(TOOLS_DIR):$$PATH:" ; find ./api -name '*.proto' -exec $(PROTOC) -I. -I$(PROTOC_IMPORTS_DIR) --go-grpc_opt=module=gitlab.figo.systems/platform/monoskope/monoskope --go-grpc_out=. {} \;
	# generates client part
	export PATH="$(TOOLS_DIR):$$PATH" ; find ./api -name '*.proto' -exec $(PROTOC) -I. -I$(PROTOC_IMPORTS_DIR) --go_opt=module=gitlab.figo.systems/platform/monoskope/monoskope --go_out=. {} \;

$(CMD_GATEWAY):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_GATEWAY) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_GATEWAY_SRC)

$(CMD_EVENTSTORE):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_EVENTSTORE) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_EVENTSTORE_SRC)

$(CMD_COMMANDHANDLER):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_COMMANDHANDLER) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_COMMANDHANDLER_SRC)

$(CMD_QUERYHANDLER):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_QUERYHANDLER) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_QUERYHANDLER_SRC)

$(CMD_CLBOREACTOR):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_CLBOREACTOR) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS)" $(CMD_CLBOREACTOR_SRC)

go-build-clean: 
	rm -Rf $(CMD_GATEWAY)
	rm -Rf $(CMD_EVENTSTORE)
	rm -Rf $(CMD_COMMANDHANDLER)
	rm -Rf $(CMD_CLBOREACTOR)

go-build-gateway: $(CMD_GATEWAY)

go-build-eventstore: $(CMD_EVENTSTORE)

go-build-commandhandler: $(CMD_COMMANDHANDLER)

go-build-queryhandler: $(CMD_QUERYHANDLER)

go-build-clboreactor: $(CMD_CLBOREACTOR)

go-rebuild-mocks: .protobuf-deps
	$(MOCKGEN) -package k8s -destination pkg/k8s/mock_client.go sigs.k8s.io/controller-runtime/pkg/client Client
	$(MOCKGEN) -package eventsourcing -destination pkg/api/eventsourcing/eventstore_client_mock.go gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing EventStoreClient,EventStore_StoreClient
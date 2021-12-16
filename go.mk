BUILD_PATH ?= $(shell pwd)
GO_MODULE ?= github.com/finleap-connect/monoskope

GO             ?= go

GOGET          ?= $(HACK_DIR)/goget-wrapper

GINKGO         ?= $(TOOLS_DIR)/ginkgo
GINKO_VERSION  ?= v1.16.4

LINTER 	   	   ?= $(TOOLS_DIR)/golangci-lint
LINTER_VERSION ?= v1.39.0

MOCKGEN         ?= $(TOOLS_DIR)/mockgen
GOMOCK_VERSION  ?= v1.5.0

PROTOC 	   	           	    ?= $(TOOLS_DIR)/protoc
PROTOC_IMPORTS_DIR          ?= $(BUILD_PATH)/include
PROTOC_VERSION              ?= 3.17.0
PROTOC_GEN_GO_VERSION       ?= v1.26.0
PROTOC_GEN_GO_GRPC_VERSION  ?= v1.1.0
PROTOC_GEN_VALIDATE_VERSION ?= 0.6.2
PROTO_FILES                 != find api -name "*.proto"

CURL          ?= curl

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

.PHONY: go-lint go-mod go-fmt go-vet go-test go-clean go-report go-protobuf ginkgo-get

##@ Go

go-mod: ## go mod download and verify
	$(GO) mod tidy
	$(GO) mod download
	$(GO) mod verify

go-fmt:  ## go fmt
	$(GO) fmt ./...

go-vet: ## go vet
	$(GO) vet ./...

go-lint: $(LINTER) ## go lint
	$(LINTER) run -v --no-config --deadline=5m

go-report: ## create report of commands and permission
	@echo
	@M8_OPERATION_MODE=cmdline $(GO) run -ldflags "$(LDFLAGS) cmd/commandhandler/*.go report commands $(ARGS)
	@echo
	@M8_OPERATION_MODE=cmdline $(GO) run -ldflags "$(LDFLAGS) cmd/commandhandler/*.go report permissions $(ARGS)
	@echo

.protobuf-deps: $(PROTO_FILES)
	for file in $$(find pkg/api/ -name "*.pb.go") ; do source=$$(awk '/^\/\/ source:/ { print $$3 }' $$file) ; echo "$$file: $$source"; done >.protobuf-deps
	echo -n "GENERATED_GO_FILES := " >>.protobuf-deps
	for file in $$(find pkg/api/ -name "*.pb.go") ; do echo -n " $$file"; done >>.protobuf-deps
	echo >>.protobuf-deps

go-test: ## run all tests
# https://onsi.github.io/ginkgo/#running-tests
	@$(GINKGO) -r -v -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoskope.coverprofile 

go-test-ci: ## run all tests in CICD
# https://onsi.github.io/ginkgo/#running-tests
	@$(GINKGO) -r -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoskope.coverprofile 

go-coverage: ## print coverage from coverprofiles
	@go tool cover -func monoskope.coverprofile

ginkgo-get $(GINKGO):
	$(shell $(GOGET) github.com/onsi/ginkgo/ginkgo@$(GINKO_VERSION))

golangci-lint-get $(LINTER):
	$(shell $(HACK_DIR)/golangci-lint.sh -b $(TOOLS_DIR) $(LINTER_VERSION))

gomock-get $(MOCKGEN):
	$(shell $(GOGET) github.com/golang/mock/mockgen@$(GOMOCK_VERSION))

protoc-get $(PROTOC):
	mkdir -p $(PROTOC_IMPORTS_DIR)/
	$(CURL) -LO "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(ARCH).zip"
	unzip protoc-$(PROTOC_VERSION)-$(ARCH).zip -d $(TOOLS_DIR)/.protoc-unpack
	mv $(TOOLS_DIR)/.protoc-unpack/bin/protoc $(TOOLS_DIR)/protoc
	cp -a $(TOOLS_DIR)/.protoc-unpack/include/* $(PROTOC_IMPORTS_DIR)/
	rm -rf $(TOOLS_DIR)/.protoc-unpack/ protoc-$(PROTOC_VERSION)-$(ARCH).zip
	$(CURL) -LO "https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v$(PROTOC_GEN_VALIDATE_VERSION).zip"
	unzip v$(PROTOC_GEN_VALIDATE_VERSION).zip -d $(TOOLS_DIR)/
	cp -a $(TOOLS_DIR)/protoc-gen-validate-$(PROTOC_GEN_VALIDATE_VERSION)/validate $(PROTOC_IMPORTS_DIR)/
	rm -rf $(TOOLS_DIR)/protoc-gen-validate-$(PROTOC_GEN_VALIDATE_VERSION) v$(PROTOC_GEN_VALIDATE_VERSION).zip
	$(shell $(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION))
	$(shell $(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION))
	$(shell $(GOGET) github.com/envoyproxy/protoc-gen-validate@v$(PROTOC_GEN_VALIDATE_VERSION))

go-tools: golangci-lint-get ginkgo-get protoc-get gomock-get ## download needed go tools

go-tools-clean:
	rm -Rf $(TOOLS_DIR)/
	mkdir $(TOOLS_DIR)

go-clean: go-build-clean ## clean up all go parts
	rm  .protobuf-deps
	rm -Rf reports/
	rm -Rf $(TOOLS_DIR)/
	mkdir $(TOOLS_DIR)
	find . -name '*.coverprofile' -exec rm {} \;

go-protobuf: .protobuf-deps
	rm -rf $(BUILD_PATH)/pkg/api
	mkdir -p $(BUILD_PATH)/pkg/api
	# generates server part
	export PATH="$(TOOLS_DIR):$$PATH:" ; find ./api -name '*.proto' -exec $(PROTOC) -I. -I$(PROTOC_IMPORTS_DIR) --go-grpc_opt=module=github.com/finleap-connect/monoskope --go-grpc_out=. --validate_out="lang=go,module=github.com/finleap-connect/monoskope:." {} \;
	# generates client part
	export PATH="$(TOOLS_DIR):$$PATH" ; find ./api -name '*.proto' -exec $(PROTOC) -I. -I$(PROTOC_IMPORTS_DIR) --go_opt=module=github.com/finleap-connect/monoskope --go_out=. --validate_out="lang=go,module=github.com/finleap-connect/monoskope:." {} \;

$(CMD_GATEWAY):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_GATEWAY) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/internal/version.Name=$(CMD_GATEWAY)" $(CMD_GATEWAY_SRC)

$(CMD_EVENTSTORE):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_EVENTSTORE) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/internal/version.Name=$(CMD_EVENTSTORE)" $(CMD_EVENTSTORE_SRC)

$(CMD_COMMANDHANDLER):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_COMMANDHANDLER) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/internal/version.Name=$(CMD_COMMANDHANDLER)" $(CMD_COMMANDHANDLER_SRC)

$(CMD_QUERYHANDLER):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_QUERYHANDLER) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/internal/version.Name=$(CMD_QUERYHANDLER)" $(CMD_QUERYHANDLER_SRC)

$(CMD_CLBOREACTOR):
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(CMD_CLBOREACTOR) -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE)/internal/version.Name=$(CMD_CLBOREACTOR)" $(CMD_CLBOREACTOR_SRC)

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

go-rebuild-mocks: .protobuf-deps $(MOCKGEN)
	$(MOCKGEN) -package k8s -destination test/k8s/mock_client.go sigs.k8s.io/controller-runtime/pkg/client Client
	$(MOCKGEN) -package eventsourcing -destination test/api/eventsourcing/eventstore_client_mock.go github.com/finleap-connect/monoskope/pkg/api/eventsourcing EventStoreClient,EventStore_StoreClient,EventStore_RetrieveClient
	$(MOCKGEN) -package eventsourcing -destination test/eventsourcing/mock_handler.go github.com/finleap-connect/monoskope/pkg/eventsourcing EventHandler
	$(MOCKGEN) -package domain -destination test/domain/repositories/repositories.go github.com/finleap-connect/monoskope/pkg/domain/repositories UserRepository,ClusterRepository
	$(MOCKGEN) -package eventsourcing -destination test/eventsourcing/aggregate_store.go github.com/finleap-connect/monoskope/pkg/eventsourcing AggregateStore

BUILD_PATH ?= $(shell pwd)

GO ?= go
CURL ?= curl

uname_S := $(shell sh -c 'uname -s 2>/dev/null || echo not')

ifeq ($(uname_S),Linux)
ARCH = linux-x86_64
endif
ifeq ($(uname_S),Darwin)
ARCH = osx-x86_64
endif

export DEX_CONFIG = $(BUILD_PATH)/config/dex
export M8_OPERATION_MODE = development

##@ Go

.PHONY: go-mod
go-mod: ## Run go mod tidy, download and verify
	$(GO) mod tidy
	$(GO) mod download
	$(GO) mod verify

.PHONY: go-fmt
go-fmt:  ## Run go fmt against code.
	$(GO) fmt ./...

.PHONY: go-vet
go-vet: ## Run go vet against code.
	$(GO) vet ./...

.PHONY: go-lint
go-lint: go-fmt go-vet golangcilint ## Run linter against code.
	$(GOLANGCILINT) run -v -E goconst -E misspell -E gofmt

.PHONY: go-test
go-test: ginkgo ## Run tests.
# https://onsi.github.io/ginkgo/#running-tests
	@find . -name '*.coverprofile' -exec rm {} \;
	@$(GINKGO) -r -v -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoskope.coverprofile 

.PHONY: go-test-ci
go-test-ci: ## Run tests in CI/CD.
# https://onsi.github.io/ginkgo/#running-tests
	@find . -name '*.coverprofile' -exec rm {} \;
	@$(GINKGO) --repeat=3 -r -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoskope.coverprofile 

.PHONY: go-coverage
go-coverage: ## Print coverage from coverprofiles.
	@go tool cover -func monoskope.coverprofile

.PHONY: go-protobuf
go-protobuf: .protobuf-deps
	rm -rf $(BUILD_PATH)/pkg/api
	mkdir -p $(BUILD_PATH)/pkg/api
	# generates server part
	export PATH="$(LOCALBIN):$$PATH:" ; find ./api -name '*.proto' -exec $(PROTOC) -I. -I$(PROTOC_IMPORTS_DIR) --go-grpc_opt=module=github.com/finleap-connect/monoskope --go-grpc_out=. --validate_out="lang=go,module=github.com/finleap-connect/monoskope:." {} \;
	# generates client part
	export PATH="$(LOCALBIN):$$PATH" ; find ./api -name '*.proto' -exec $(PROTOC) -I. -I$(PROTOC_IMPORTS_DIR) --go_opt=module=github.com/finleap-connect/monoskope --go_out=. --validate_out="lang=go,module=github.com/finleap-connect/monoskope:." {} \;

.PHONY: go-rebuild-mocks
go-rebuild-mocks: .protobuf-deps gomock
	$(MOCKGEN) -copyright_file hack/copyright.lic -destination internal/test/sigs.k8s.io/controller-runtime/pkg/client.go sigs.k8s.io/controller-runtime/pkg/client Client
	$(MOCKGEN) -copyright_file hack/copyright.lic -destination internal/test/api/eventsourcing/eventstore_client.go github.com/finleap-connect/monoskope/pkg/api/eventsourcing EventStoreClient,EventStore_StoreClient,EventStore_RetrieveClient
	$(MOCKGEN) -copyright_file hack/copyright.lic -destination internal/test/api/eventsourcing/commandhandler_client.go github.com/finleap-connect/monoskope/pkg/api/eventsourcing CommandHandlerClient
	$(MOCKGEN) -copyright_file hack/copyright.lic -destination internal/test/api/domain/user_client.go github.com/finleap-connect/monoskope/pkg/api/domain UserClient,User_GetAllClient
	$(MOCKGEN) -copyright_file hack/copyright.lic -destination internal/test/api/gateway/gateway_auth_client.go github.com/finleap-connect/monoskope/pkg/api/gateway GatewayAuthClient
	$(MOCKGEN) -copyright_file hack/copyright.lic -destination internal/test/eventsourcing/mock_handler.go github.com/finleap-connect/monoskope/pkg/eventsourcing EventHandler
	$(MOCKGEN) -copyright_file hack/copyright.lic -destination internal/test/eventsourcing/aggregate_store.go github.com/finleap-connect/monoskope/pkg/eventsourcing AggregateStore
	$(MOCKGEN) -copyright_file hack/copyright.lic -destination internal/test/domain/repositories/repositories.go github.com/finleap-connect/monoskope/pkg/domain/repositories UserRepository,ClusterRepository,ClusterAccessRepository

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: clean
clean: ## Clean up build dependencies
	rm -R $(LOCALBIN)

## Tool Binaries
MOCKGEN ?= $(LOCALBIN)/mockgen
GINKGO ?= $(LOCALBIN)/ginkgo
GOLANGCILINT ?= $(LOCALBIN)/golangci-lint
PROTOC ?= $(LOCALBIN)/protoc

## Tool Versions
GOMOCK_VERSION  ?= v1.5.0
GINKGO_VERSION ?= v1.16.5
GOLANGCILINT_VERSION ?= v1.48.0
PROTOC_VERSION ?= 3.17.0
PROTOC_GEN_GO_VERSION ?= v1.26.0
PROTOC_GEN_GO_GRPC_VERSION ?= v1.1.0
PROTOC_GEN_VALIDATE_VERSION ?= 0.6.2

## Tool Config
PROTOC_IMPORTS_DIR          ?= $(BUILD_PATH)/include
PROTO_FILES                 != find api -name "*.proto"

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo locally if necessary.
$(GINKGO): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/ginkgo@$(GINKGO_VERSION)

.PHONY: golangcilint
golangcilint: $(GOLANGCILINT) ## Download golangci-lint locally if necessary.
$(GOLANGCILINT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

.PHONY: gomock
gomock: $(MOCKGEN) ## Download mockgen locally if necessary.
$(MOCKGEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golang/mock/mockgen@$(GOMOCK_VERSION)

.protobuf-deps: protoc-get $(PROTO_FILES)
	for file in $$(find pkg/api/ -name "*.pb.go") ; do source=$$(awk '/^\/\/ source:/ { print $$3 }' $$file) ; echo "$$file: $$source"; done >.protobuf-deps
	echo -n "GENERATED_GO_FILES := " >>.protobuf-deps
	for file in $$(find pkg/api/ -name "*.pb.go") ; do echo -n " $$file"; done >>.protobuf-deps
	echo >>.protobuf-deps
	
protoc-get $(PROTOC):
	mkdir -p $(PROTOC_IMPORTS_DIR)/
	$(CURL) -LO "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(ARCH).zip"
	unzip protoc-$(PROTOC_VERSION)-$(ARCH).zip -d $(LOCALBIN)/.protoc-unpack
	mv $(LOCALBIN)/.protoc-unpack/bin/protoc $(LOCALBIN)/protoc
	cp -a $(LOCALBIN)/.protoc-unpack/include/* $(PROTOC_IMPORTS_DIR)/
	rm -rf $(LOCALBIN)/.protoc-unpack/ protoc-$(PROTOC_VERSION)-$(ARCH).zip
	$(CURL) -LO "https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v$(PROTOC_GEN_VALIDATE_VERSION).zip"
	unzip v$(PROTOC_GEN_VALIDATE_VERSION).zip -d $(LOCALBIN)/
	cp -a $(LOCALBIN)/protoc-gen-validate-$(PROTOC_GEN_VALIDATE_VERSION)/validate $(PROTOC_IMPORTS_DIR)/
	rm -rf $(LOCALBIN)/protoc-gen-validate-$(PROTOC_GEN_VALIDATE_VERSION) v$(PROTOC_GEN_VALIDATE_VERSION).zip
	GOBIN=$(LOCALBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	GOBIN=$(LOCALBIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	GOBIN=$(LOCALBIN) go install github.com/envoyproxy/protoc-gen-validate@v$(PROTOC_GEN_VALIDATE_VERSION)
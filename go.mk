BUILD_PATH ?= $(shell pwd)

COPYRIGHT_FILE ?= hack/copyright.lic

GO ?= go
CURL ?= curl

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
	@$(GINKGO) -r -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoskope.coverprofile 

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
	$(MOCKGEN) -copyright_file $(COPYRIGHT_FILE) -destination internal/test/sigs.k8s.io/controller-runtime/pkg/client.go sigs.k8s.io/controller-runtime/pkg/client Client
	$(MOCKGEN) -copyright_file $(COPYRIGHT_FILE) -destination internal/test/api/eventsourcing/eventstore_client.go github.com/finleap-connect/monoskope/pkg/api/eventsourcing EventStoreClient,EventStore_StoreClient,EventStore_RetrieveClient
	$(MOCKGEN) -copyright_file $(COPYRIGHT_FILE) -destination internal/test/api/eventsourcing/commandhandler_client.go github.com/finleap-connect/monoskope/pkg/api/eventsourcing CommandHandlerClient
	$(MOCKGEN) -copyright_file $(COPYRIGHT_FILE) -destination internal/test/api/domain/user_client.go github.com/finleap-connect/monoskope/pkg/api/domain UserClient,User_GetAllClient
	$(MOCKGEN) -copyright_file $(COPYRIGHT_FILE) -destination internal/test/api/gateway/gateway_auth_client.go github.com/finleap-connect/monoskope/pkg/api/gateway GatewayAuthClient
	$(MOCKGEN) -copyright_file $(COPYRIGHT_FILE) -destination internal/test/eventsourcing/mock_handler.go github.com/finleap-connect/monoskope/pkg/eventsourcing EventHandler
	$(MOCKGEN) -copyright_file $(COPYRIGHT_FILE) -destination internal/test/eventsourcing/aggregate_store.go github.com/finleap-connect/monoskope/pkg/eventsourcing AggregateStore
	$(MOCKGEN) -copyright_file $(COPYRIGHT_FILE) -destination internal/test/domain/repositories/repositories.go github.com/finleap-connect/monoskope/pkg/domain/repositories UserRepository,ClusterRepository,ClusterAccessRepository

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
GOLANGCILINT_VERSION ?= v1.46.1
PROTOC_VERSION ?= 21.2
PROTOC_GEN_GO_VERSION ?= v1.28
PROTOC_GEN_GO_GRPC_VERSION ?= v1.2
PROTOC_GEN_VALIDATE_VERSION ?= 0.6.7

GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
PROTO_ARCH := $(GOARCH)
PROTO_OS := $(GOOS)

ifeq ($(PROTO_OS),darwin)
PROTO_OS = osx
endif
ifeq ($(PROTO_ARCH),arm64)
PROTO_ARCH = aarch_64
endif
PROTO_ARCH_OS := $(PROTO_OS)-$(PROTO_ARCH)

## Tool Config
PROTOC_IMPORTS_DIR          ?= $(BUILD_PATH)/api_includes
PROTO_FILES                 != find api -name "*.proto"

ginkgo: $(GINKGO) ## Download ginkgo locally if necessary.
$(GINKGO): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/ginkgo@$(GINKGO_VERSION)

golangcilint: $(GOLANGCILINT) ## Download golangci-lint locally if necessary.
$(GOLANGCILINT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

gomock: $(MOCKGEN) ## Download mockgen locally if necessary.
$(MOCKGEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golang/mock/mockgen@$(GOMOCK_VERSION)

.protobuf-deps: protoc-get $(PROTO_FILES)
	@for file in $$(find pkg/api/ -name "*.pb.go") ; do source=$$(awk '/^\/\/ source:/ { print $$3 }' $$file) ; echo "$$file: $$source"; done >.protobuf-deps
	@echo -n "GENERATED_GO_FILES := " >>.protobuf-deps
	@for file in $$(find pkg/api/ -name "*.pb.go") ; do echo -n " $$file"; done >>.protobuf-deps
	@echo >>.protobuf-deps
	
protoc-get $(PROTOC):
	mkdir -p $(PROTOC_IMPORTS_DIR)/
	$(CURL) -fLO "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTO_ARCH_OS).zip"
	unzip protoc-$(PROTOC_VERSION)-$(PROTO_ARCH_OS).zip -d $(LOCALBIN)/.protoc-unpack
	mv -f $(LOCALBIN)/.protoc-unpack/bin/protoc $(LOCALBIN)/protoc
	cp -a $(LOCALBIN)/.protoc-unpack/include/* $(PROTOC_IMPORTS_DIR)/
	rm -rf $(LOCALBIN)/.protoc-unpack/ protoc-$(PROTOC_VERSION)-$(PROTO_ARCH_OS).zip
	$(CURL) -fLO "https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v$(PROTOC_GEN_VALIDATE_VERSION).zip"
	unzip v$(PROTOC_GEN_VALIDATE_VERSION).zip -d $(LOCALBIN)/
	mkdir -p $(PROTOC_IMPORTS_DIR)/validate/
	cp -a $(LOCALBIN)/protoc-gen-validate-$(PROTOC_GEN_VALIDATE_VERSION)/validate/validate.proto $(PROTOC_IMPORTS_DIR)/validate/
	rm -rf $(LOCALBIN)/protoc-gen-validate-$(PROTOC_GEN_VALIDATE_VERSION) v$(PROTOC_GEN_VALIDATE_VERSION).zip
	GOBIN=$(LOCALBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	GOBIN=$(LOCALBIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	GOBIN=$(LOCALBIN) go install github.com/envoyproxy/protoc-gen-validate@v$(PROTOC_GEN_VALIDATE_VERSION)
	
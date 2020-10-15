GO         ?= go
GINKGO     ?= ginkgo
MOCKGEN    ?= mockgen
LINTER     ?= golangci-lint
VERSION    ?= 0.0.1-dev
DOCKER     ?= docker
COMMIT     := $(shell git rev-parse --short HEAD)
LDFLAGS    += -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)"
BUILDFLAGS += -installsuffix cgo --tags release
PROTOC     ?= protoc

BUILD_PATH ?= $(shell pwd)

define go-run
	$(GO) run $(LDFLAGS) cmd/$(1)/*.go $(ARGS)
endef

.PHONY: lint clean prepare

clean:
	rm -f $(CMD_AUTHORITY)
	rm -f $(CMD_GATEWAY)
	rm -f $(CMD_SECRETMANAGER)
	rm -f $(CMD_HOOKMUX)

prepare:
	$(GO) mod download

lint:
	$(GO) mod verify
	$(LINTER) run -v --no-config --deadline=5m

run-%:
	$(call go-run,$*)


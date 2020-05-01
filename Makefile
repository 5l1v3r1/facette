# -*- Makefile -*-

VERSION := 1.0.0-dev
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
REVISION := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date +"%F %T")

GO ?= go

GO_BUILD_LDFLAGS := \
	-X 'facette.io/facette/pkg/internal/version.Version=$(VERSION)' \
	-X 'facette.io/facette/pkg/internal/version.Branch=$(BRANCH)' \
	-X 'facette.io/facette/pkg/internal/version.Revision=$(REVISION)' \
	-X 'facette.io/facette/pkg/internal/version.BuildDate=$(BUILD_DATE)'

ifeq ($(filter dev,$(TAGS)),)
GO_BUILD_LDFLAGS += -s -w
endif

GO_BUILD_ARGS := -ldflags "$(GO_BUILD_LDFLAGS)" -tags "$(TAGS)"
ifeq ($(filter dev,$(TAGS)),)
GO_BUILD_ARGS += -trimpath
endif

GO_TEST_ARGS := -cover

ADDLICENSE_ARGS := -c "The Facette Authors" -f docs/LICENSE.template

all: build

clean: clean-bin

clean-bin:
	@echo "==> Cleaning source tree"
	@rm -rf bin/

build: build-bin

build-bin: generate
	@echo "==> Building binaries"
	@$(GO) build $(GO_BUILD_ARGS) -o bin/ -v ./cmd/...

test: test-bin

test-bin:
	@echo "==> Testing binary sources"
	@$(GO) test $(GO_TEST_ARGS) -v ./...

lint: lint-bin

lint-bin:
	@echo "==> Linting binary sources"
	@golangci-lint run ./...

generate:
	@echo "==> Running code generation"
	@$(GO) generate ./...

license:
	@echo "==> Adding license headers"
	@addlicense $(ADDLICENSE_ARGS) cmd/ pkg/ ui/src ui/types

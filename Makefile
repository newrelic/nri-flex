PROJECT_NAME := $(shell basename $(shell pwd))
NATIVEOS     := $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH   := $(shell go version | awk -F '[ /]' '{print $$5}')
GO_PKGS      := $(shell go list ./... | grep -v -e "/vendor/" -e "/example")
SRCDIR       ?= .
BUILD_DIR    := ./bin/
COVERAGE_DIR := ./coverage/
COVERMODE     = atomic

GO_CMD        = go
GODOC         = godocdown
GOLINTER      = golangci-lint

GORELEASER_VERSION := v0.126.0
GORELEASER_SHA256 := 6c0145df61140ec1bffe4048b9ef3e105e18a89734816e7a64f342d3f9267691
GORELEASER_BIN ?= bin/goreleaser

# Determine packages by looking into pkg/*
ifneq ("$(wildcard ${SRCDIR}/pkg/*)","")
	PACKAGES  = $(wildcard ${SRCDIR}/pkg/*)
endif
ifneq ("$(wildcard ${SRCDIR}/internal/*)","")
	PACKAGES += $(wildcard ${SRCDIR}/internal/*)
endif

# Determine commands by looking into cmd/*
COMMANDS = $(wildcard ${SRCDIR}/cmd/*)

# Determine binary names by stripping out the dir names
BINS=$(foreach cmd,${COMMANDS},$(notdir ${cmd}))

all: build

# Humans running make:
build: check-version clean lint test-unit coverage compile document

# Build command for CI tooling
build-ci: check-version clean lint test-integration

clean:
	@echo "=== $(PROJECT_NAME) === [ clean            ]: removing binaries and coverage file..."
	@rm -rfv $(BUILD_DIR)/* $(COVERAGE_DIR)/*

bin:
	@mkdir -p bin

$(GORELEASER_BIN): bin
	@echo "=== $(PROJECT) === [ release/deps ]: Installing goreleaser"
	@(wget -qO /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/$(GORELEASER_VERSION)/goreleaser_$(GOOS)_x86_64.tar.gz)
	@(tar -xf  /tmp/goreleaser.tar.gz -C bin/)
	@(rm -f /tmp/goreleaser.tar.gz)

release/deps: $(GORELEASER_BIN)

release: release/deps
	@echo "=== $(PROJECT) === [ release ]: Releasing new version..."
	@$(GORELEASER_BIN) release

# Import fragments
include build/deps.mk
include build/compile.mk
include build/testing.mk
include build/util.mk
include build/document.mk
include build/docker.mk

.PHONY: all build build-ci

INTEGRATION      := flex
PROJECT_NAME     = nri-$(INTEGRATION)
NATIVEOS         := $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH       := $(shell go version | awk -F '[ /]' '{print $$5}')
SRCDIR           ?= .
BUILD_DIR        ?= $(CURDIR)/bin
COVERAGE_FILE    ?= coverage.out

GO_VERSION       ?= 1.25
GO_CMD           ?= go
GODOC            ?= godocdown

GOLINTER         = golangci-lint
GOLINTER_VERSION = v1.24.0

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
build: check-version clean deps lint test-unit compile document

# lint (golangci-lint) disabled as it does not support Go 1.18 yet https://github.com/golangci/golangci-lint/issues/2649
# Build command for CI tooling
build-ci: check-version clean deps test-coverage

clean:
	@echo "=== $(PROJECT_NAME) === [ clean ]: removing binaries and coverage file..."
	@rm -rfv $(BUILD_DIR)/* $(COVERAGE_FILE)

bin:
	@mkdir -p $(BUILD_DIR)

# Import fragments
include build/deps.mk
include build/compile.mk
include build/setup.mk
include build/testing.mk
include build/util.mk
include build/document.mk
include build/docker.mk
include build/ci.mk
include build/release.mk

.PHONY: all build build-ci

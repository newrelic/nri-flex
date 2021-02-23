#
# Makefile fragment for installing deps
#   Auto-detects between dep and govendor
#

# Go file to track tool deps with go modules
TOOL_DIR     ?= tools
TOOL_CONFIG  ?= $(TOOL_DIR)/tools.go

GOTOOLS		  ?=
GOTOOLS       += github.com/axw/gocov/gocov
GOTOOLS       += github.com/AlekSi/gocov-xml
GOTOOLS       += github.com/robertkrimen/godocdown/godocdown
GOTOOLS       += github.com/jandelgado/gcov2lcov


VENDOR_CMD	= $(GO_CMD) mod vendor

.PHONY: tools
tools: check-version
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@cd $(TOOL_DIR)
	@$(GO_CMD) get $(GOTOOLS)
	@$(GO_CMD) mod tidy

.PHONY: tools-update
tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update     ]: Updating tools required by the project..."
	@cd $(TOOL_DIR)
	@$(GO) get -u $(GOTOOLS)
	@$(VENDOR_CMD)

.PHONY: deps
deps: tools deps-only

.PHONY: deps-only
deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Using '$(VENDOR_CMD)'"
	@$(VENDOR_CMD)

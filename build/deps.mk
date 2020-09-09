#
# Makefile fragment for installing deps
#   Auto-detects between dep and govendor
#

GOTOOLS       = github.com/axw/gocov/gocov \
                github.com/AlekSi/gocov-xml \
                github.com/stretchr/testify/assert \
                github.com/robertkrimen/godocdown/godocdown

VENDOR_CMD	= $(GO_CMD) mod vendor

.PHONY: tools
tools: check-version
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@$(GO_CMD) get $(GOTOOLS)

.PHONY: tools-update
tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update     ]: Updating tools required by the project..."
	@$(GO_CMD) get -u $(GOTOOLS)

.PHONY: deps
deps: tools deps-only

.PHONY: deps-only
deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Using '$(VENDOR_CMD)'"
	@$(VENDOR_CMD)

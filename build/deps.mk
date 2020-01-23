#
# Makefile fragment for installing deps
#   Auto-detects between dep and govendor
#

GOTOOLS       = github.com/axw/gocov/gocov \
                github.com/AlekSi/gocov-xml \
                github.com/stretchr/testify/assert \
                github.com/robertkrimen/godocdown/godocdown \
                github.com/golangci/golangci-lint/cmd/golangci-lint

# Determine package dep manager
ifneq (,"$(wildcard go.mod)")
	VENDOR     = ${GO_CMD}
	VENDOR_CMD = ${VENDOR} mod download
else ifneq (,"$(wildcard Gopkg.toml)")
	VENDOR     = dep
	VENDOR_CMD = ${VENDOR} ensure
	GOTOOLS    += github.com/golang/dep
else
	VENDOR     = govendor
	VENDOR_CMD = ${VENDOR} sync
	GOTOOLS    += github.com/kardianos/govendor
endif

tools: check-version
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@$(GO_CMD) get $(GOTOOLS)

tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update     ]: Updating tools required by the project..."
	@$(GO_CMD) get -u $(GOTOOLS)

deps: tools deps-only

deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Using '$(VENDOR_CMD)'"
	@$(VENDOR_CMD)


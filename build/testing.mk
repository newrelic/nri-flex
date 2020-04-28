#
# Makefile fragment for Testing
#

GOLINTER_BIN = bin/$(GOLINTER)

$(GOLINTER_BIN): bin
	@echo "=== $(PROJECT_NAME) === [ lint ]: Installing $(GOLINTER)..."
	@(wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLINTER_VERSION))

lint: bin $(GOLINTER_BIN)
	@echo "=== $(PROJECT_NAME) === [ lint ]: Validating source code running $(GOLINTER)..."
	@$(GOLINTER_BIN) run ./...

test: test-only
test-only: test-unit

test-unit:
	@echo "=== $(PROJECT_NAME) === [ unit-test ]: running unit tests..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GO_CMD) test -tags unit -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/unit.tmp $(GO_PKGS)

test-integration: setup test
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integration tests..."
	@mkdir -p $(COVERAGE_DIR)
	@sh ./integration-test/ci-test.sh

cover-report:
	@echo "=== $(PROJECT_NAME) === [ cover-report ]: generating coverage results..."
	@mkdir -p $(COVERAGE_DIR)
	@echo 'mode: $(COVERMODE)' > $(COVERAGE_DIR)/coverage.out
	@cat $(COVERAGE_DIR)/*.tmp | grep -v 'mode: $(COVERMODE)' >> $(COVERAGE_DIR)/coverage.out || true
	@$(GO_CMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]:     $(COVERAGE_DIR)coverage.html"

cover-view: cover-report
	@$(GO_CMD) tool cover -html=$(COVERAGE_DIR)/coverage.out

.PHONY : test-linux
test-linux:
	@(echo "=== $(PROJECT_NAME) === [ unit-test-linux ]: running unit tests for Linux...")
	@(docker build -t nri-flex-test -f ./integration-test/Dockerfile .)
	@(docker run --rm -it nri-flex-test go test ./...)

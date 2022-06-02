#
# Makefile fragment for Testing
#

GOLINTER_BIN = bin/$(GOLINTER)
TEST_PATTERN ?=.
TEST_OPTIONS ?=

TEST_FLAGS += -failfast
TEST_FLAGS += -race

GO_TEST ?= test $(TEST_OPTIONS) $(TEST_FLAGS) ./... -run $(TEST_PATTERN) -timeout=10m

$(GOLINTER_BIN): bin
	@echo "=== $(PROJECT_NAME) === [ lint ]: Installing $(GOLINTER)..."
	@(wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLINTER_VERSION))

.PHONY: lint
lint: bin $(GOLINTER_BIN)
	@echo "=== $(PROJECT_NAME) === [ lint ]: Validating source code running $(GOLINTER)..."
	@$(GOLINTER_BIN) run ./...

.PHONY: test
test: test-only

.PHONY: test-only
test-only: test-unit

.PHONY: test-unit
test-unit:
	@echo "=== $(PROJECT_NAME) === [ unit-test ]: running unit tests..."
	@$(GO_CMD) $(GO_TEST)

.PHONY : test-coverage
test-coverage: TEST_FLAGS += -covermode=atomic -coverprofile=$(COVERAGE_FILE)
test-coverage:
	@echo "=== $(PROJECT_NAME) === [ test-coverage ]: running unit tests with coverage..."
	@$(GO_CMD) $(GO_TEST)

.PHONY : convert-coverage
convert-coverage:
	@(gcov2lcov -infile=$(COVERAGE_FILE) -outfile=lcov.info)

.PHONY: test-integration
test-integration: setup ci/deps
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integration tests..."
	# TODO restore
	#@sh ./integration-test/ci-test.sh
	@echo "SKIPPED as it remains loopin with: waiting for database deployment to be available"

.PHONY : test-linux
test-linux:
	@(echo "=== $(PROJECT_NAME) === [ unit-test-linux ]: running unit tests for Linux...")
	@(docker build -t nri-flex-test -f ./integration-test/Dockerfile .)
	@(docker run --rm -it nri-flex-test make test-unit)


.PHONY: test-e2e
test-e2e: bin/nri-flex
	@echo "=== $(PROJECT_NAME) === [ e2e-test ]: running e2e tests..."
	@cd ./test/testbed && sh ./launch_e2e.sh


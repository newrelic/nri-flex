#
# Makefile fragment for Testing
#



lint: deps
	@echo "=== $(PROJECT_NAME) === [ lint             ]: Validating source code running $(GOLINTER)..."
	@$(GOLINTER) run ./...

test: test-only
test-only: test-unit test-integration

test-unit:
	@echo "=== $(PROJECT_NAME) === [ unit-test        ]: running unit tests..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GO_CMD) test -tags unit -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/unit.tmp $(GO_PKGS)

test-integration: setup
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integration tests..."
	@mkdir -p $(COVERAGE_DIR)
	@sh ./integration-test/ci-test.sh

cover-report:
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]: generating coverage results..."
	@mkdir -p $(COVERAGE_DIR)
	@echo 'mode: $(COVERMODE)' > $(COVERAGE_DIR)/coverage.out
	@cat $(COVERAGE_DIR)/*.tmp | grep -v 'mode: $(COVERMODE)' >> $(COVERAGE_DIR)/coverage.out || true
	@$(GO_CMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]:     $(COVERAGE_DIR)coverage.html"

cover-view: cover-report
	@$(GO_CMD) tool cover -html=$(COVERAGE_DIR)/coverage.out

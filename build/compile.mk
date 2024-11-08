#
# Makefile Fragment for Compiling
#

LDFLAGS ?= -s -w

.PHONY: compile
compile: deps compile-only

.PHONY: compile-all
compile-all: compile-linux compile-darwin compile-windows

.PHONY: build-all
build-all: compile-linux compile-darwin compile-windows

.PHONY: compile-only
compile-only: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ compile          ]:     $(BUILD_DIR)$(GOOS)/$$b"; \
		BUILD_FILES="$(SRCDIR)/cmd/..." ; \
		$(GO_CMD) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(GOOS)/$$b $$BUILD_FILES ; \
	done

.PHONY: fmt
fmt:
	@($(GO_CMD) fmt ./...)

bin/nri-flex:
	@(CGO_ENABLED=0 $(GO_CMD) build -ldflags="$(LDFLAGS)" -trimpath -o ./bin/nri-flex ./cmd/nri-flex/)

.PHONY: build-linux
build-linux: compile-linux

.PHONY: compile-linux
compile-linux: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-linux    ]: building commands:"
	@mkdir -p $(BUILD_DIR)/linux
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)/linux/$$b" ; \
		echo "=== $(PROJECT_NAME) === [ compile-linux    ]:     $$OUTPUT_FILE"; \
		BUILD_FILES="$(SRCDIR)/cmd/..." ; \
		GOOS=linux $(GO_CMD) build -ldflags="$(LDFLAGS)" -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

.PHONY: compile-linux-fips
compile-linux-fips: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-linux-fips    ]: building commands:"
	@mkdir -p $(BUILD_DIR)/linux-fips
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)/linux-fips/$$b" ; \
		echo "=== $(PROJECT_NAME) === [ compile-linux-fips    ]:     $$OUTPUT_FILE"; \
		BUILD_FILES="$(SRCDIR)/cmd/..." ; \
		GOOS=linux GOFIPS=1 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOEXPERIMENT=boringcrypto $(GO_CMD) build -tags fips -ldflags="$(LDFLAGS)" -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

.PHONY: build-darwin
build-darwin: compile-darwin

.PHONY: compile-darwin
compile-darwin: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-darwin   ]: building commands:"
	@mkdir -p $(BUILD_DIR)/darwin
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)/darwin/$$b" ; \
		echo "=== $(PROJECT_NAME) === [ compile-darwin   ]:     $$OUTPUT_FILE"; \
		BUILD_FILES="$(SRCDIR)/cmd/..." ; \
		GOOS=darwin $(GO_CMD) build -ldflags="$(LDFLAGS)" -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

.PHONY: build-windows
build-windows: compile-windows

.PHONY: compile-windows
compile-windows: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-windows  ]: building commands:"
	@mkdir -p $(BUILD_DIR)/windows
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)/windows/$$b.exe" ; \
		echo "=== $(PROJECT_NAME) === [ compile-windows  ]:     $$OUTPUT_FILE"; \
		BUILD_FILES="$(SRCDIR)/cmd/..." ; \
		GOOS=windows $(GO_CMD) build -ldflags="$(LDFLAGS)" -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

.PHONY: build-windows32
build-windows32: compile-windows32

.PHONY: compile-windows32
compile-windows32: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-windows  ]: building commands:"
	@mkdir -p $(BUILD_DIR)/windows
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)/windows/$$b.exe" ; \
		echo "=== $(PROJECT_NAME) === [ compile-windows  ]:     $$OUTPUT_FILE"; \
		BUILD_FILES="$(SRCDIR)/cmd/..." ; \
		GOARCH=386 CGO_ENABLED=1 GOOS=windows $(GO_CMD) build -ldflags="$(LDFLAGS)" -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

.PHONY: compile-for-debug-linux
compile-for-debug-linux: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile-for-debug-linux    ]: building commands:"
	@mkdir -p $(BUILD_DIR)/linux
	@for b in $(BINS); do \
		OUTPUT_FILE="$(BUILD_DIR)/linux/$$b" ; \
		echo "=== $(PROJECT_NAME) === [ compile-for-debug-linux    ]:     $$OUTPUT_FILE"; \
		BUILD_FILES="$(SRCDIR)/cmd/..." ; \
		GOOS=linux $(GO_CMD) build -gcflags 'all=-N -l' -o $$OUTPUT_FILE $$BUILD_FILES ; \
	done

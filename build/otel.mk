#
# Makefile fragment for creating otel collector
#

# OpenTelemetry Flex Receiver targets
.PHONY: install-otel-builder
install-otel-builder:
	@echo "=== $(PROJECT_NAME) === [ install-otel-builder ]: installing OpenTelemetry Collector Builder..."
	@go install go.opentelemetry.io/collector/cmd/builder@latest

.PHONY: build-otel-flex
build-otel-flex: install-otel-builder
	@echo "=== $(PROJECT_NAME) === [ build-otel-flex ]: building custom OpenTelemetry Collector with Flex receiver..."
	@$(shell go env GOPATH)/bin/builder --config=otelcol-builder.yaml

.PHONY: run-otel-flex
run-otel-flex: build-otel-flex
	@echo "=== $(PROJECT_NAME) === [ run-otel-flex ]: running custom OpenTelemetry Collector with Flex receiver..."
	@./bin/otelcol-flex --config=configs/otel-collector-config.yaml

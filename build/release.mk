GORELEASER_VERSION  ?= v0.157.0
GORELEASER_BIN      ?= $(CURDIR)/bin/goreleaser
GORELEASER_CONFIG   ?= --config $(CURDIR)/build/goreleaser.yml
PKG_FLAGS			?= --rm-dist
IS_RELEASE			?= false # Default to safe mode which is pre-release

ifneq ($(IS_RELEASE), "false")
	PKG_FLAGS += --snapshot
endif

$(GORELEASER_BIN): bin
	@echo "=== $(PROJECT_NAME) === [ release/deps ]: Installing goreleaser"
	@(wget -qO /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/$(GORELEASER_VERSION)/goreleaser_$(GOOS)_x86_64.tar.gz)
	@(tar -xf  /tmp/goreleaser.tar.gz -C $(CURDIR)/bin/)
	@(rm -f /tmp/goreleaser.tar.gz)
	@echo "=== [$(GORELEASER_BIN)] goreleaser downloaded"

.PHONY : release/deps
release/deps: $(GORELEASER_BIN)

.PHONY : release
release: clean release/deps compile-only
	@echo "=== $(PROJECT_NAME) === [ release ]: Releasing new version..."
	$(GORELEASER_BIN) release $(GORELEASER_CONFIG) $(PKG_FLAGS)

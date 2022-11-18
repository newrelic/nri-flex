CI_BUILDER_TAG ?= nri-$(INTEGRATION)-builder

.PHONY : ci/deps
ci/deps:
	@docker build \
		-t $(CI_BUILDER_TAG) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		-f $(CURDIR)/build/Dockerfile $(CURDIR)

.PHONY : ci/snyk-test
ci/snyk-test:
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-snyk-test" \
		-v $(CURDIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
			-e SNYK_TOKEN \
			snyk/snyk:golang snyk test --severity-threshold=high

.PHONY : ci/pre-release
ci/pre-release: ci/deps
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-release" \
		-v $(CURDIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-e IS_RELEASE \
		-e GITHUB_TOKEN \
		-e TAG \
		$(CI_BUILDER_TAG) \
		make release

.PHONY : ci/test
ci/test: ci/deps
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-release" \
		-v $(CURDIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-e IS_RELEASE \
		-e GITHUB_TOKEN \
		-e TAG \
		$(CI_BUILDER_TAG) \
		make build-ci

.PHONY : ci/convert-coverage
ci/convert-coverage: ci/deps
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-release" \
		-v $(CURDIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-e IS_RELEASE \
		-e GITHUB_TOKEN \
		-e TAG \
		$(CI_BUILDER_TAG) \
		make convert-coverage

.PHONY : ci/snyk-test
ci/snyk-test:
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-snyk-test" \
		-v $(CURDIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
			-e SNYK_TOKEN \
			snyk/snyk:golang snyk test --severity-threshold=hig

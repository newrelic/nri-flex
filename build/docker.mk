#
# Makefile fragment for Docker actions
#

IMAGE_NAME     := newrelic-es/$(PROJECT_NAME)
DOCKER         := docker
DOCKER_COMPOSE := docker-compose -f scripts/docker-compose-build.yml
TEST_CMD := go test -v -coverprofile=coverage.txt -covermode=atomic ./...
LINE_BREAK := "--------------------------------------------------------------------"

docker-image: compile-linux
	$(DOCKER) build -t $(IMAGE_NAME):latest .

docker-clean:
	@echo "### Removing $(PROJECT_NAME) containers..."
	@echo ${LINE_BREAK}
	@$(DOCKER) rm -f $$($(DOCKER) ps -a | grep $(PROJECT_NAME) | cut -d' ' -f 1)
	@echo ${LINE_BREAK}

docker-run:
	@echo "### Running via docker-compose..."
	@read -p "Enter any CLI arguments, else hit enter to skip:" args; \
	$(DOCKER_COMPOSE) run golang sh -c "$(TEST_CMD) && $(GO_CMD) run cmd/nri-flex/nri-flex.go $$args"


#
# Testing within Docker
#
docker-test-setup:
	@echo "### If first run, this may take some time..."
	$(DOCKER_COMPOSE) -f scripts/docker-compose-build.yml build

docker-test: docker-test-setup
	@echo "### Testing via docker-compose (linux)"
	$(DOCKER_COMPOSE) run golang $(TEST_CMD)

docker-test-infra: docker-setup compile-linux docker-clean
	@echo "### Testing within NR Infra Container"
	@read -p "Enter Infrastructure License Key:" infrakey; \
	$(DOCKER) run -d --name $(BINARY_NAME) --network=host --cap-add=SYS_PTRACE \
	-v "/:/host:ro" -v "/var/run/docker.sock:/var/run/docker.sock" \
	-e NRIA_LICENSE_KEY=$$infrakey $(IMAGE_NAME):$(PROJECT_VER)
	$(DOCKER) ps -a | grep $(IMAGE_NAME)


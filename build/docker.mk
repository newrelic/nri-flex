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

docker-image-infra: compile-linux
	$(DOCKER) build -f Dockerfile.infra -t $(IMAGE_NAME)-infra:latest .

docker-image-run:
	@echo "### Running via docker run"
	@$(DOCKER) run --rm $(IMAGE_NAME):latest

docker-image-infra-run:
	@echo "### Running Flex with NR Infra agent via docker run"
	@$(DOCKER) run --rm $(IMAGE_NAME)-infra:latest

#
# Testing within Docker
#
docker-test-setup:
	@echo "### If first run, this may take some time..."
	$(DOCKER_COMPOSE) -f scripts/docker-compose-build.yml build

docker-test: docker-test-setup
	@echo "### Testing via docker-compose (linux)"
	$(DOCKER_COMPOSE) run golang $(TEST_CMD)

docker-test-infra: compile-linux docker-image-infra
	@echo "### Testing within NR Infra Container"
	@read -p "Enter Infrastructure License Key:" infrakey; \
	$(DOCKER) run --rm --name $(PROJECT_NAME) --network=host --cap-add=SYS_PTRACE \
	-v "/:/host:ro" -v "/var/run/docker.sock:/var/run/docker.sock" \
	-e NRIA_LICENSE_KEY=$$infrakey $(IMAGE_NAME)-infra:latest 


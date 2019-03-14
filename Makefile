#   
#	Description: Generic Makefile to easily build and package
#
#	Requirements: Golang, docker, docker-compose
#	
#

INTEGRATION  := $(shell basename $(shell pwd))
BINARY_NAME := $(INTEGRATION)
VERSION := $$(cat internal/load/load.go | grep IntegrationVersion | head -1 | cut -d'"' -f 2)
PKG_NAME_NIX := $(BINARY_NAME)_linux-v$(VERSION)
PKG_NAME_WIN := $(BINARY_NAME)_win-v$(VERSION)
PKG_NAME_DARWIN := $(BINARY_NAME)_darwin-v$(VERSION)
IMAGE_NAME := newrelic-es/$(BINARY_NAME)
ARCH ?= $$(uname -s | tr A-Z a-z)
GOARCH := amd64
DOCKER_COMPOSE := docker-compose -f scripts/docker-compose-build.yml
TEST_CMD := go test -v -coverprofile=coverage.txt -covermode=atomic ./...
LINE_BREAK := "--------------------------------------------------------------------"

setup: 
	@echo "### If first run, this may take some time..."
	@echo ${LINE_BREAK}
	$(DOCKER_COMPOSE) build
	dep ensure
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@echo ${LINE_BREAK}

run:
	@echo "### Running locally..."
	@echo ${LINE_BREAK}
	@read -p "Enter any CLI arguments, else hit enter to skip:" args; \
	go run cmd/flex/$(BINARY_NAME).go $$args
	# go run $$(ls -1 *.go | grep -v _test.go) $$args
	@echo ${LINE_BREAK}

version:
	@echo $(VERSION)

#default build for current os
build: setup
	@echo "### Building for Current OS"
	@echo ${LINE_BREAK}
	# docker-compose run golang sh -c "go build -o bin/$(ARCH)/$(BINARY_NAME) cmd/flex/$(BINARY_NAME).go"
	$(DOCKER_COMPOSE) run golang sh -c "GOOS=$(ARCH) go build -o bin/$(ARCH)/$(BINARY_NAME) cmd/flex/$(BINARY_NAME).go"
	@echo ${LINE_BREAK}

build-linux: setup
	@echo "### Building Linux binary"
	@echo ${LINE_BREAK}
	$(DOCKER_COMPOSE) run golang sh -c "GOOS=linux go build -o bin/linux/$(BINARY_NAME) cmd/flex/$(BINARY_NAME).go"
	@echo ${LINE_BREAK}
	# docker-compose run golang GOOS=linux GOARCH=$(GOARCH) go build -o bin/linux/$(BINARY_NAME)

build-darwin: setup
	@echo "### Building for Mac binary"
	@echo ${LINE_BREAK}
	$(DOCKER_COMPOSE) run golang sh -c "GOOS=darwin go build -o bin/darwin/$(BINARY_NAME) cmd/flex/$(BINARY_NAME).go"
	@echo ${LINE_BREAK}

build-windows: setup
	@echo "### Building for Windows binary"
	@echo ${LINE_BREAK}
	$(DOCKER_COMPOSE) run golang sh -c "GOOS=windows go build -o bin/windows/$(BINARY_NAME).exe cmd/flex/$(BINARY_NAME).go"
	@echo ${LINE_BREAK}

build-all: build-linux build-windows build-darwin

build-docker: build-linux
	docker build -t $(IMAGE_NAME):$(VERSION) .
	@echo ${LINE_BREAK}

clean-docker:
	@echo "### Removing $(BINARY_NAME) containers..."
	@echo ${LINE_BREAK}
	@docker rm -f $$(docker ps -a | grep $(BINARY_NAME) | cut -d' ' -f 1)
	@echo ${LINE_BREAK}

package-linux: build-linux
	@echo "### Packaging into $(PKG_NAME_NIX).tar"
	@echo ${LINE_BREAK}
	@rm -rf $(BINARY_NAME)_linux-*
	@mkdir $(PKG_NAME_NIX)
	@cp ./bin/linux/$(BINARY_NAME) ./$(PKG_NAME_NIX)/
	@cp ./README.md ./$(PKG_NAME_NIX)/
	@cp ./Dockerfile ./$(PKG_NAME_NIX)/
	@cp ./scripts/install_linux.sh ./$(PKG_NAME_NIX)/
	@cp ./examples/$(BINARY_NAME)-config.yml ./$(PKG_NAME_NIX)/
	@cp ./examples/$(BINARY_NAME)-def-nix.yml ./$(PKG_NAME_NIX)/
	@cp -av ./examples ./$(PKG_NAME_NIX)/
	@cp -av ./nrjmx ./$(PKG_NAME_NIX)/
	@tar -cvf $(PKG_NAME_NIX).tar $(PKG_NAME_NIX)/
	@rm -rf $(PKG_NAME_NIX)
	@echo "Completed packaging: $(PKG_NAME_NIX).tar"
	@echo ${LINE_BREAK}

package-windows: build-windows
	@echo "### Packaging into $(PKG_NAME_WIN).tar"
	@echo ${LINE_BREAK}
	@rm -rf $(BINARY_NAME)_win-*
	@mkdir $(PKG_NAME_WIN)
	@cp ./bin/windows/$(BINARY_NAME).exe ./$(PKG_NAME_WIN)/
	@cp ./README.md ./$(PKG_NAME_WIN)/
	@cp ./Dockerfile ./$(PKG_NAME_WIN)/
	@cp ./scripts/install_win.bat ./$(PKG_NAME_WIN)/
	@cp ./examples/$(BINARY_NAME)-config.yml ./$(PKG_NAME_WIN)/
	@cp ./examples/$(BINARY_NAME)-def-nix.yml ./$(PKG_NAME_WIN)/
	@cp -av ./examples ./$(PKG_NAME_WIN)/
	@cp -av ./nrjmx ./$(PKG_NAME_WIN)/
	@tar -cvf $(PKG_NAME_WIN).tar $(PKG_NAME_WIN)/
	@rm -rf $(PKG_NAME_WIN)
	@echo "Completed packaging: $(PKG_NAME_WIN).tar"
	@echo ${LINE_BREAK}

package-darwin: build-darwin
	@echo "### Packaging into $(PKG_NAME_DARWIN).tar"
	@echo ${LINE_BREAK}
	@rm -rf $(BINARY_NAME)_darwin-*
	@mkdir $(PKG_NAME_DARWIN)
	@cp ./bin/darwin/$(BINARY_NAME) ./$(PKG_NAME_DARWIN)/
	@cp ./README.md ./$(PKG_NAME_DARWIN)/
	@cp ./Dockerfile ./$(PKG_NAME_DARWIN)/
	@cp ./scripts/install_linux.sh ./$(PKG_NAME_DARWIN)/
	@cp ./examples/$(BINARY_NAME)-config.yml ./$(PKG_NAME_DARWIN)/
	@cp ./examples/$(BINARY_NAME)-def-nix.yml ./$(PKG_NAME_DARWIN)/
	@cp -av ./examples ./$(PKG_NAME_DARWIN)/
	@cp -av ./nrjmx ./$(PKG_NAME_DARWIN)/
	@tar -cvf $(PKG_NAME_DARWIN).tar $(PKG_NAME_DARWIN)/
	@rm -rf $(PKG_NAME_DARWIN)
	@echo "Completed packaging: $(PKG_NAME_DARWIN).tar"
	@echo ${LINE_BREAK}

package-mac: package-darwin

package-all: package-linux package-windows package-darwin
	@echo "Completed packaging for Linux, Windows & Mac"
	@echo ${LINE_BREAK}

clean: clean-docker
	@echo "### Removing folders: vendor, bin"
	@echo ${LINE_BREAK}
	rm -rf vendor bin
	@echo ${LINE_BREAK}

test: setup do-test lint clean-docker

lint: 
	@echo "### Running Linter"
	@echo ${LINE_BREAK}
	golangci-lint run -v
	@echo ${LINE_BREAK}

do-test:
	@echo "### Testing via docker-compose (linux)"
	@echo ${LINE_BREAK}
	-docker-compose run golang $(TEST_CMD)
	@echo ${LINE_BREAK}

view: 
	go tool cover -html=coverage.txt

run-docker: 
	@echo "### Running via docker-compose..."
	@echo ${LINE_BREAK}
	@read -p "Enter any CLI arguments, else hit enter to skip:" args; \
	docker-compose run golang sh -c "$(TEST_CMD) && go run cmd/flex/$(BINARY_NAME).go $$args"
	@echo ${LINE_BREAK}

run-docker-test: build-linux clean-docker
	@echo "### Testing within NR Infra Container"
	@echo ${LINE_BREAK}
	@read -p "Enter Infrastructure License Key:" infrakey; \
	docker run -d --name $(BINARY_NAME) --network=host --cap-add=SYS_PTRACE \
	-v "/:/host:ro" -v "/var/run/docker.sock:/var/run/docker.sock" \
	-e NRIA_LICENSE_KEY=$$infrakey $(IMAGE_NAME):$(VERSION)
	docker ps -a | grep $(IMAGE_NAME)
	@echo ${LINE_BREAK}

test-docker: run-docker-test

.PHONY : setup build-linux build-windows build-darwin
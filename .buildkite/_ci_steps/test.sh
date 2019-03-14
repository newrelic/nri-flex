#!/bin/sh
export CGO_ENABLED=0

echo -e "--- \033[33m Grabbing Dependencies \033[0m :golang::flying_saucer:"
dep ensure -v

echo -e "--- \033[33m Running Tests \033[0m :golang::hammer_and_wrench:"
go test -v -coverprofile=coverage.txt -covermode=atomic ./...

echo -e "--- \033[33m Running Linter \033[0m :golang::lint-remover:"
golangci-lint run -v
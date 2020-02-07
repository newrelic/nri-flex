#!/usr/bin/env sh

cd "$(dirname "$0")"

docker-compose build

docker-compose run --rm nri-flex-src \
    go test --tags=integration -covermode=atomic -coverprofile ./coverage/integration.tmp ./...

result=$?

docker-compose down

exit $result

#!/usr/bin/env sh

set -e

KIND=bin/kind
KCTL=bin/kubectl

cleanup() {
  echo '=== cleaning up'
  $KIND delete cluster
}

trap "cleanup" ERR

# create out KiD cluster
$KIND create cluster

echo '=== building https-server container image'
# build our http/s server
docker build . -f integration-test/Dockerfile_https -t newrelic/https-server:integration-test
# load the image into K8s cluster (we need it later when we apply the manifest)
$KIND load docker-image newrelic/https-server:integration-test
echo '=== building nri-flex container image'
# build the Flex image
docker build . -f integration-test/Dockerfile -t newrelic/nri-flex:integration-test
# load it into the K8s cluster (we need it to run the integration tests)
$KIND load docker-image newrelic/nri-flex:integration-test
echo '=== deploying to K8s'
# deploy our required services in K8s (right now just the http/s server we use for the https tests)
$KCTL apply -f integration-test/k8s.yaml
# make sure services are running
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'
until $KCTL -n default get pods -lapp=database-server -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 5;echo "waiting for database deployment to be available"; done
until $KCTL -n default get pods -lapp=https-server -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1;echo "waiting for https-server deployment to be available"; done
echo '=== running tests'
# run the integration tests and save the result. this works because we run with attach=true
$KCTL run nri-flex --rm --restart=Never --attach=true --image newrelic/nri-flex:integration-test -- go test --tags=integration ./...
result=$?

cleanup()

exit $result

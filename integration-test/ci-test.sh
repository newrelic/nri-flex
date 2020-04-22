#!/usr/bin/env sh

# create out KiD cluster
kind create cluster

# build our http/s server
docker build . -f integration-test/Dockerfile_https -t newrelic/https-server:integration-test
# load the image into K8s cluster (we need it later when we apply the manifest)
kind load docker-image newrelic/https-server:integration-test

# build the Flex image
docker build . -f integration-test/Dockerfile -t newrelic/nri-flex:integration-test
# load it into the K8s cluster (we need it to run the integration tests)
kind load docker-image newrelic/nri-flex:integration-test

# deploy our required services in K8s (right now just the http/s server we use for the https tests)
# TODO find a better way to guarantee a service is running
kubectl apply -f integration-test/k8s.yaml
# make sure services are running
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'
until kubectl -n default get pods -lapp=database-server -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 5;echo "waiting for database deployment to be available";  done
until kubectl -n default get pods -lapp=https-server -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1;echo "waiting for https-server deployment to be available"; done

# run the integration tests and save the result. this works because we run with attach=true
kubectl run nri-flex --rm --restart=Never --attach=true --image newrelic/nri-flex:integration-test -- go test --tags=integration ./...
result=$?
#cleanup
kind delete cluster

exit $result

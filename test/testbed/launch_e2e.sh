#!/bin/bash

go test -c scenarios/e2e_test.go
cp ../../bin/nri-flex .

docker build -t flex_e2e_tests --build-arg flex_bin=nri-flex --build-arg flex_tests_bin=scenarios.test -f Containerfile .
docker run --rm flex_e2e_tests

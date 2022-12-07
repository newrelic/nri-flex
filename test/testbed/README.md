# E2E tests

**IMPORTANT:** End-to-end are not meant to replace standard Go unit and integration tests, instead, they are an additional layer of tests to assure
we don't break any of the current supported features.

This directory contains the needed files to generate and run end-to-end tests for `nri-flex`. As the integration can
run any arbitrary command, the integration output can be very different depending on the underlying system/environment.
To overcome this issue, the tests on this directory **must** run in the container image specified in the [Containerfile](./Containerfile).
The image uses an `Ubuntu:22.04` as base image, thus assuring the same environment for the test execution.

## New test case

To add a new E2E test, we would need a nri-flex configuration and its exact output after running it in an `ubuntu:22.04` container. The configuration should make use of the new functionalities we want to test.

First, we would need to indentify which nri-flex API is using our configuration and its corresponding Go test file:

- `command`: [scenarios/fixtures/command_api.go](./scenarios/fixtures/command_api.go)
- `file`: [scenarios/fixtures/file_api.go](./scenarios/fixtures/file_api.go)
- `url`: [scenarios/fixtures/url_api.go](./scenarios/fixtures/url_api.go)

Every API test file contains an array structure to easily add new test cases, for example, the `command_api.go` has the following initialized variable:

```Go
var CommandTests = []struct {
	Name           string
	Config         string
	ExpectedStdout string
}{...}
```

The `Name` determines the test name, the `Config` the raw nri-flex configuration we want to add and the `ExpectedStdout` parameter the output of running that configuration in the `ubuntu:22.04` container. As the `CommandTests` is an array, adding a new test would only require appending a new object into it.

Note that each API test file has a different configuration depending on their needs. For example, the URL API queries the provided endpoint, as **no public endpoints are allowed**, the configuration allows to specify a payload to be returned for queries in the local endpoint.
If the test case we want to add requires additional changes in the execution environment not covered by the current APIs configuration, we can extend its configuration parameters and its logic in the [scenarios/e2e_test.go](./scenarios/e2e_test.go).


## Usage

From the root level of the repository:

```
make test-e2e
```

In the background, the make target builds and executes a script used to build and run the tests, from the current
directory:

```
bash launch_e2e.sh
```

Mainly, the script generates the [e2e_test.go](./scenarios/e2e_test.go) as binary files, builds the container image with
those binaries and executes the image entrypoint.

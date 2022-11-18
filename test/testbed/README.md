# E2E tests

This directory contains the needed files to generate and run end-to-end tests for `nri-flex`. As the integration can
run any arbitrary command, the integration output can be very different depending on the underlying system/environment.
To overcome this issue, the tests on this directory **must** run in the container image specified in the [Containerfile](./Containerfile).
The image uses an `Ubuntu:22.04` as base image, thus assuring the same environment for the test execution.


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
those binaries and nri-flex and executes the image.
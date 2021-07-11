# Integration tests

This folder contains integrations tests especially prepared to run on CI builds.

Some configuration files in the `configs` folder are very similar to those found in the `examples/flexConfigs` folders, 
but adapted to be used in integration tests.

The integration tests are grouped by specific use cases (like `run command in linux`).

It is recommended to run these from the root of the project by running the Make target:

```shell
$ make test-integration
```

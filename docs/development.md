# Development

Here you can learn how to build Flex from source, run it as a Docker image, or run it without the New Relic agent (development mode). Development mode is useful when developing and testing your config files.

* [Development mode](#developmentmode)
	* [Before you start](#Beforeyoustart)
	* [Installation](#Installation)
	* [Standard configuration](#Standardconfiguration)
	* [Testing](#Testing)
* [Build from source](#Compilefromsource)
	* [Requirements](#Requirements)
	* [Setup](#Setup)
	* [Build](#Build)
	* [Cross-compiling](#Cross-compiling)
	* [Packaging](#Packaging)
	* [Docker Related](#DockerRelated)
	* [Other Utility Commands](#OtherUtilityCommands)
* [Run as a Docker image](#Docker)

## Development mode

###  <a name='Beforeyoustart'></a>Before you start

- Flex outputs to the terminal/console, so you don't need to send the data to New Relic to see the results of running Flex against you config file.
- Flex runs everything by default in the `flexConfigs/` folder, next to the binary file.
- Browse `examples/flexConfigs` for Flex configurations that you can reuse.

### <a name='Installation'></a>Installation

1. Download the latest [release](https://github.com/newrelic/nri-flex/releases) for your development platform.
2. Unpack the file.
3. Run `./nri-flex -help` to see all available flags.

### <a name='Standardconfiguration'></a>Standard configuration

Flex looks for configuration files in a folder named `flexConfigs` by default.

You can use the following flags to instruct Flex to read configuration files from somewhere else than the default folder:

* `config_dir` `string` Specifies a directory of configurations files
* `config_file` (or `-config_path`) `string` Specifies a single config file

### <a name='Testing'></a>Testing your configuration

Running without any flags defaults to running all configs within `./flexConfigs`:

```bash
./nri-flex
```
To test a single Flex config file use `-config_file`:

```bash
./nri-flex -config_file "examples/flexConfigs/redis-cmd-raw-example.yml"
```

For additional logging, use `-verbose`:

```bash
./nri-flex -verbose
```

Once you've tested your configuration and you're ready to use in production, you can:

- Add your configuration to the integrations config file in `integrations.d`.
	```yaml
	integrations:                                    # OHI configuration starts here  
      - name: nri-flex                               # OHI to be executed by the Agent
        config:                                      # OHI configuration to be parsed by Flex
        # Actual Flex configuration starts here
	```
	or
- Use `config_template_path` to reference your Flex configuration file from the integrations config file:
	```yaml
	integrations:
       - name: nri-flex
          interval: 60s
          timeout: 5s
          config_template_path: /path/to/flex/integration.yml
	```

##  <a name='Compilefromsource'></a>Build from source

### <a name='Requirements'></a>Requirements

- Make
- Go 1.13 or higher
- [dep](https://github.com/golang/dep) - Dependency management tool (if not using `go mod`, which we advise you to use)
- [golangci-lint v1.22.2](https://github.com/golangci/golangci-lint)
- Docker Compose (for integration tests)

### <a name='Setup'></a>Setup

This assumes that you have a functional Go environment:

```bash
go get github.com/newrelic/nri-flex

cd ${GOPATH}/src/github.com/newrelic/nri-flex

# Ensure a clean start
make clean

# Download all required libraries
make dep
```

### <a name='Build'></a>Build

```bash
# Default command runs clean, linter, unit test, and compiles for the local OS
make

# run all tests + linter
make test

# run integration tests (requires docker-compose)
make test-integration

# run unit tests
make test-unit

# run only linter
make lint

# Create a coverage report
make cover

# Launch the coverage report into a web browser
make cover-view
```

### <a name='Cross-compiling'></a>Cross-compiling

```bash
# Build binary for current OS
make build

# Build binaries for all supported OSes
make build-all

# Build binaries for a specific OS
make build-darwin
make build-linux
make build-windows
```

### <a name='Packaging'></a>Packaging

To build tar.gz files for distribution:

```bash
# Create a package for the current OS
make package

# Create packages for all supported OSes
make package-all

# Create packages for a specific OS
make package-darwin
make package-linux
make package-windows
```

###  <a name='DockerRelated'></a>Docker related

```bash
# clean/remove any docker containers that have been created
make docker-clean

# Build a new docker image
make docker-image

# Run via docker-compose
make docker-run

# Testing within docker
make docker-test

# Testing with the Infrastructure Agent within Docker
make docker-test-infra
```

### <a name='OtherUtilityCommands'></a>Other utility Commands

```bash
# Use godocdown to create Markdown documentation for all commands and packages
# this is run by default.
make document
```

## <a name='Docker'></a>Run as a Docker image

- Set your configs, modify Dockerfile if need be.
- Build and run the image.

```bash
# BUILD
docker build -t nri-flex .

# RUN - standard
docker run -d --name nri-flex --network=host --cap-add=SYS_PTRACE -v "/:/host:ro" -v "/var/run/docker.sock:/var/run/docker.sock" -e NRIA_LICENSE_KEY="yourInfraLicenseKey" nri-flex:latest

# RUN - with container discovery reverse lookup (ensure -container_discovery is set to true nri-flex-config.yml)
docker run -d --name nri-flex --network=host --cap-add=SYS_PTRACE -l flexDiscoveryRedis="t=redis,c=redis,tt=img,tm=contains,r=true"  -v "/:/host:ro" -v "/var/run/docker.sock:/var/run/docker.sock" -e NRIA_LICENSE_KEY="yourInfraLicenseKey" nri-flex:latest

# Example: Run Redis with a flex discovery label
docker run -it -p 9696:6379 --label flexDiscoveryRedis="t=redis,c=redis,tt=img,tm=contains" --name redis-svr -d redis
```

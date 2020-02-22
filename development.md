# Development mode

Flex can run in isolation while you're developing and testing your configuration files.
The instructions below apply only when running in this mode.
For instructions in running Flex in within the New Relic Infrastructure agent see the [documentation](./docs/wiki/README.md)

## Before you start

- See inside examples/flexConfigs for configuration examples that you can reuse
- Flex will run everything by default in the default flexConfigs/ folder, next to the binary file.
- Flex provides two options for ingesting your events, via the New Relic Infrastructure agent & the New Relic Insights Event API.
  If you want so the results of running Flex against you configuration file, Flex will output the results to the terminal/console so you don't need to send the data during this period to New Relic.
- Flex has an Event Limiter built in - the event_limit argument is available and there to ensure you don't spam heaps of events unknowingly, the default is 500 per execution/run, which can be dialled up if required.

## Installation

- Go to [releases](https://github.com/newrelic/nri-flex/releases) and download the latest release for your development platform
- Unpack the archive
- Run ./nri-flex -help for all available flags.

## Standard Configuration

- Flex by default looks for configuration files in a folder named */flexConfigs*.
- If you have your configuration files somewhere else, you can use the following flags to instruct Flex to read configuration files from somwhere else than the defautl folder.

    -config_dir `string` Specifies a directory of configurations files

    -config_file (or -config_path) `string` Specifies a single config file

## New Relic Insights Event API

If you do want to see the data in New Relic while developping you can use the New Relie Insights Event APi in Flex. Note that it is not officially supported and we may remove this at any point.

To use the New Relic Insight Event API from within Flex, you will need to:

- create an Insert API Key from here: <https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/manage/api_keys>
- Use the below flags to configure
  
    -insights_api_key `string` : set the Insights API key you created before

    -insights_url `string` : set Insights API endpoint
    ex: <https://insights-collector.newrelic.com/v1/accounts/YOUR_ACCOUNT_ID/events>

    -insights_output `bool` : Output the events generated to standard out true/false

## Serverless

[Serverless README](examples/lambdaExample/README.md)

## Kubernetes

- Build your Docker Image, and deploy as a daemonset, view [examples/nri-flex-k8s.yml](examples/nri-flex-k8s.yml)
- View wiki for information on how to use [service discovery](https://github.com/newrelic/nri-flex/wiki/Service-Discovery)

## Testing

- For additional logging use the `-verbose` flag

- Testing a single config

```bash
./nri-flex -config_file "examples/flexConfigs/redis-cmd-raw-example.yml"
```

- Running without any flags, will default to run all configs within ./flexConfigs

```bash
./nri-flex
```

- Additional Logging

```bash
./nri-flex -verbose
```

## Compile from source

### Requirements

- Make
- Go 1.13 or later
- [dep](https://github.com/golang/dep) - Dependency management tool (if not using go mod, which we advise you to use)
- [golangci-lint v1.22.2](https://github.com/golangci/golangci-lint)
- Docker Compose (Integration tests)

### Setup

_Note:_ This assumes that you have a functional Go environment.

```bash
go get github.com/newrelic/nri-flex

cd ${GOPATH}/src/github.com/newrelic/nri-flex

# Ensure a clean start
make clean

# Download all required libraries
make dep
```

### Build

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

### Cross-Compiling

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

### Packaging

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

### Docker Related

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

### Other Utility Commands

```bash
# Use godocdown to create Markdown documentation for all commands and packages
# this is run by default.
make document
```

## Releasing

The build process sets the package version based on the latest git tag. After
all changes have been made for the lastest release, make a new tag with NO
commits after, and then `make package-all` to create the artifacts.

This process should be automated someday.

Finally, upload the artifacts on Github to the tag release.

## Docker

- Set your configs, modify Dockerfile if need be
- Build & Run Image

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

## Features & Support

-   Run any HTTP/S request, read file, shell command, consume from any Prometheus Exporter, Database Query, or JMX Query. (Java 7+ is required for JMX to work)
-   Service / Container Discovery
-   Attempt to cleverly flatten to samples
-   Use environment variables anywhere within config files (eg. using double dollar signs -> \$\$MY_ENV_VAR)
-   Detect and flatten dimensional data from Prometheus style payloads (vector, matrix, targets supported)
-   Merge different samples and outputs together
-   Key Remover & Replacer
-   Metric Parser for RATE & DELTA support (has capability to auto set rates and deltas)
-   Define multiple APIs / commands or mix
-   event_type autoset or override
-   Define custom attributes (more granular control, compared to NR infra agent)
-   Command allows horizontal split (useful for table style data) (use only once per command set)
-   snake_case to CamelCase conversion
-   Percentage to Decimal conversion
-   ToLower conversion
-   SubParse functionality (see redis config for an example)
-   LookUp Store - save attributes from previously generated samples to use in requests later (see rabbit example)
-   LazyFlatten - for arrays
-   Inbuilt data caching - useful for processing samples at different points
-   [+more here](https://github.com/newrelic/nri-flex/wiki/Functions)

## Integrations

For all see within the examples directory as there are many more.

-   All Prometheus Exporters
-   Consul
-   Vault (shows merge functionality)
-   Bamboo
-   Teamcity
-   CircleCI
-   RabbitMQ (shows metric parser, and lookup store functionality)
-   Elasticsearch (shows inbuilt URL cache functionality)
-   Traefik
-   Kong
-   etcd (shows custom sample keys functionality)
-   Varnish
-   Redis (more metrics, multi instance support, multi db support) (shows snake to camel, perc to decimal, replace keys, rename keys & sub parse functionality)
-   Zookeeper
-   OpsGenie
-   VictorOps
-   PagerDuty (shows lazy_flatten functionality)
-   AlertOps (shows lazy_flatten functionality)
-   New Relic Alert Ingestion (provides similar output to nri-alerts-pipe)
-   New Relic App Status Health Ingestion (appSample to present your app health, language, and aggregated summary)
-   http/s testing & request performance via curl
-   Postgres Custom Querying
-   MySQL Custom Querying
-   MariaDB Custom Querying
-   Percona Server, Google CloudSQL or Sphinx (2.2.3+) Custom Querying
-   MS SQL Server Custom Querying
-   JMX via nrjmx // (nrjmx is targetted to work with Java 7+, see cassandra and tomcat examples)
-   cassandra - via jmx
-   tomcat - via jmx
-   bind9
-   df display disk & inode info (shows horizontal split functionality)

## Disclaimer

New Relic has open-sourced this integration to enable monitoring of various technologies. This integration is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to this integration via GitHub. Support for this integration is available with an Expert Services subscription.

<!-- vscode-markdown-toc -->
1. [Before you start](#Beforeyoustart)
2. [Installation](#Installation)
3. [Standard configuration](#Standardconfiguration)
4. [New Relic Insights Event API](#NewRelicInsightsEventAPI)
5. [Serverless](#Serverless)
6. [Kubernetes](#Kubernetes)
7. [Testing](#Testing)
8. [Compile from source](#Compilefromsource)
	8.1. [Requirements](#Requirements)
	8.2. [Setup](#Setup)
	8.3. [Build](#Build)
	8.4. [Cross-compiling](#Cross-compiling)
	8.5. [Packaging](#Packaging)
	8.6. [Docker Related](#DockerRelated)
	8.7. [Other Utility Commands](#OtherUtilityCommands)
9. [Releasing](#Releasing)
10. [Docker](#Docker)
11. [Features & Support](#FeaturesSupport)
12. [Integrations](#Integrations)
13. [Disclaimer](#Disclaimer)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc --># Development mode

Flex can run in isolation while you're developing and testing your configuration files; these instructions below apply only when running in this mode. For instructions on running Flex with the New Relic Infrastructure agent, see the [docs](./docs/wiki/README.md).



##  1. <a name='Beforeyoustart'></a>Before you start

- See inside `examples/flexConfigs` for configuration examples that you can reuse.
- Flex runs everything by default in the `flexConfigs/` folder, next to the binary file.
- Flex provides two options for ingesting your events, via the New Relic Infrastructure agent and the New Relic Insights Event API.
- If you want the see the results of running Flex against you configuration file, Flex outputs to the terminal/console, so you don't need to send the data to New Relic during debug.
- Flex has a built-in Event Limiter - the `event_limit` argument is available to ensure you don't spam heaps of events unknowingly; default is `500` per execution/run, which can be increased if required.

##  2. <a name='Installation'></a>Installation

1. Go to [releases](https://github.com/newrelic/nri-flex/releases) and download the latest release for your development platform.
2. Unpack the archive.

Run `./nri-flex -help` to see all available flags.

##  3. <a name='Standardconfiguration'></a>Standard configuration

Flex looks for configuration files in a folder named */flexConfigs* by default.

If you have your configuration files somewhere else, you can use the following flags to instruct Flex to read configuration files from somewhere else than the defautl folder:

    -config_dir `string` Specifies a directory of configurations files

    -config_file (or -config_path) `string` Specifies a single config file

##  4. <a name='NewRelicInsightsEventAPI'></a>New Relic Insights Event API

If you want to see your data in New Relic while developing, use the New Relic Insights Event API in Flex. Note that it is not officially supported and we may remove this at any point in the future.

To use the New Relic Insight Event API from within Flex, you need to:

- Create an Insert API Key here: <https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/manage/api_keys>
- Use the below flags to configure Flex:
  
    -`insights_api_key` (`string`): Insights API key
    -`insights_url` (`string`): Insights API endpoint
    (for example, <https://insights-collector.newrelic.com/v1/accounts/YOUR_ACCOUNT_ID/events>)
    -`insights_output` (`bool`): whether events generated go to standard out

##  5. <a name='Serverless'></a>Serverless

See the [Serverless README](examples/lambdaExample/README.md).

##  6. <a name='Kubernetes'></a>Kubernetes

- Build your Docker Image, and deploy as a daemonset, see an example in [examples/nri-flex-k8s.yml](examples/nri-flex-k8s.yml)
- See the docs for more information on how to use [service discovery](https://github.com/newrelic/nri-flex/wiki/Service-Discovery)

##  7. <a name='Testing'></a>Testing

For additional logging, use the `-verbose` flag.

Testing a single config:

```bash
./nri-flex -config_file "examples/flexConfigs/redis-cmd-raw-example.yml"
```

Running without any flags, defaults to running all configs within `./flexConfigs`:

```bash
./nri-flex
```

For additional logging:

```bash
./nri-flex -verbose
```

##  8. <a name='Compilefromsource'></a>Compile from source

###  8.1. <a name='Requirements'></a>Requirements

- Make
- Go 1.13 or higher
- [dep](https://github.com/golang/dep) - Dependency management tool (if not using `go mod`, which we advise you to use)
- [golangci-lint v1.22.2](https://github.com/golangci/golangci-lint)
- Docker Compose (for integration tests)

###  8.2. <a name='Setup'></a>Setup

This assumes that you have a functional Go environment:

```bash
go get github.com/newrelic/nri-flex

cd ${GOPATH}/src/github.com/newrelic/nri-flex

# Ensure a clean start
make clean

# Download all required libraries
make dep
```

###  8.3. <a name='Build'></a>Build

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

###  8.4. <a name='Cross-compiling'></a>Cross-compiling

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

###  8.5. <a name='Packaging'></a>Packaging

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

###  8.6. <a name='DockerRelated'></a>Docker Related

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

###  8.7. <a name='OtherUtilityCommands'></a>Other Utility Commands

```bash
# Use godocdown to create Markdown documentation for all commands and packages
# this is run by default.
make document
```

##  9. <a name='Releasing'></a>Releasing

The build process sets the package version based on the latest git tag. 

After all changes have been made for the latest release, make a new tag with no commits after, and then `make package-all` to create the artifacts. 

Finally, upload the artifacts on Github to the tag release.

##  10. <a name='Docker'></a>Docker

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

##  11. <a name='FeaturesSupport'></a>Features & Support

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

##  12. <a name='Integrations'></a>Integrations

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

##  13. <a name='Disclaimer'></a>Disclaimer

New Relic has open-sourced this integration to enable monitoring of various technologies. This integration is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to this integration via GitHub. Support for this integration is available with an Expert Services subscription.

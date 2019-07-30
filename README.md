# New Relic - Flex
[![Build Status](https://badge.buildkite.com/fe011ab8474a98a28ac3255e2141ec3887e9accda7d3c31196.svg?branch=master)](https://buildkite.com/kav91/nri-flex)

- Flex is an agnostic AIO New Relic Integration, that can:
  - Abstract the need for end users to write any code other then to define a configuration yml, allowing you to consume metrics from practically anywhere!
  - Run any HTTP/S request, read file, shell command, consume from any Prometheus Exporter, Database Query, or JMX Query. (Java 7+ is required for JMX to work)
  - Can generate New Relic metric samples automatically from almost [any endpoint for almost any payload, useful helper functions exist to tidy up your output neatly](#features--support).
  - Simplies deployment and configuration as a single Flex integration can be running multiple integrations which would be just config files.
  - Provides over [200+ Integrations](#integrations)
  - Has agnostic [Service / Container Discovery](https://github.com/newrelic/nri-flex/wiki/Service-Discovery) built-in
  - As updates and upgrades are made, all Flex Integrations reap the benefits.
  - Can send data via the New Relic Infrastructure Agent, or the New Relic Insights Event API

## Disclaimer
New Relic has open-sourced this integration to enable monitoring of various technologies. This integration is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to this integration via GitHub. Support for this integration is available with an Expert Services subscription.

## Requirements
- Linux
- Windows
- New Relic Infrastructure

## Usage
- [Download the latest compiled release under the Releases section](https://github.com/newrelic/nri-flex/releases)
- [Config Examples](https://github.com/newrelic/nri-flex/tree/master/cmd/flex/examples)
- [Testing](#testing)
- [Standard Config Layout](#standard-configuration)
- [Installation](#installation)
- [Development](#development)
- [Releasing](#releasing)
- [Contributing](#contributing)

## Further Documentation
- [Features & Support](#features--support)
- [Existing Integrations](#integrations)
- [Wiki](https://github.com/newrelic/nri-flex/wiki)
- [Available Functions](https://github.com/newrelic/nri-flex/wiki/Functions)
- [Using Service Discovery](https://github.com/newrelic/nri-flex/wiki/Service-Discovery)
- [Using Prometheus-Integrations-(Exporters)](https://github.com/newrelic/nri-flex/wiki/Prometheus-Integrations-(Exporters))
- [Creating your own Flex configuration(s)](https://github.com/newrelic/nri-flex/wiki/Creating-Flex-Configs)

---

## Testing
- [Compiled Releases](https://github.com/newrelic/nri-flex/releases)
- For additional logging use the `-verbose` flag
```
Testing a single config
./nri-flex -config_file "examples/flexConfigs/redis-cmd-raw-example.yml"

Running without any flags, will default to run all configs within ./flexConfigs
./nri-flex 

Additional Logging
./nri-flex -verbose
```

## Standard Configuration
- Default configuration looks for Flex config files in /flexConfigs.
- Run ./nri-flex -help for all available flags.
- Flex has an Event Limiter built in - the event_limit argument is available and there to ensure you don't spam heaps of events unknowingly, the default is 500 per execution/run, which can be dialled up if required.

``` 
The below two flags you could specific a single Flex Config, or another config directory.

-config_dir string
        Set directory of config files (default "flexConfigs/")
-config_file string
        Set a specific config file

With these flags, you could also define multiple instances with different configurations of Flex within "nri-flex-config.yml" 
```

## Installation

- Setup your configuration(s) see inside examples/flexConfigs for examples
- Flex will run everything by default in the default flexConfigs/ folder (so keep what you want before deploy)
- Flex provides two options for ingesting your events, via the New Relic Infrastructure Agent, & the New Relic Insights Event API

### New Relic Infrastructure Agent
- Review the commented out portions in the install_linux.sh and/or Dockerfile depending on your config setup
- Run scripts/install_linux.sh or build the docker image
- Alternatively use the scripts/install_linux.sh as a guide for setting up (or scripts/install_win.bat)

#### Typical file/directory structure layout:
```
/etc/newrelic-infra/integrations.d/nri-flex-config.yml <- config 
(/examples/nri-flex-config.yml)

/var/db/newrelic-infra/custom-integrations/nri-flex-def-nix.yml <- definition 
(/examples/nri-flex-def-nix.yml)

/var/db/newrelic-infra/custom-integrations/nri-flex <- binary 
(compiled binary)

/var/db/newrelic-infra/custom-integrations/flexConfigs/ <- standard flexConfigs (refer to examples here: /examples/flexConfigs)

/var/db/newrelic-infra/custom-integrations/flexContainerDiscovery/ <- if using container discovery 
(refer to examples here: /examples/flexContainerDiscovery)
```

### New Relic Insights Event API

- Able to execute the binary wherever and however you like
- The Flex specific config folders will remain the same
- To use this method, create an Insert API Key from here: https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/manage/api_keys
- Use the below flags to configure
```
  -insights_api_key string
        Set Insights API key - from link above
  -insights_url string
        Set Insights URL eg. "https://insights-collector.newrelic.com/v1/accounts/YOUR_ACCOUNT_ID/events"
  -insights_output bool
        Output the events generated to standard out true/false
```
- Run ./nri-flex -help for all available options

#### Typical file/directory structure layout:
```
From any location:

nri-flex <- binary
# below folders in the same location as the binary unless you've specific a different location
flexConfigs/ <- folder
flexContainerDiscovery/ <- folder
```

### Serverless

[Serverless README](examples/lambdaExample/README.md)

## Development

### Requirements

* Make
* Go 1.10 or later
* [dep](https://github.com/golang/dep) - Dependency management tool
* Docker Compose (Integration tests)

### Setup

*Note:* This assumes that you have a functional Go environment.

```
go get github.com/newrelic/nri-flex

cd ${GOPATH}/src/github.com/newrelic/nri-flex

# Ensure a clean start
make clean

# Download all required libraries
make dep
```

### Build

```
# Default command runs clean, linter, unit test, and compiles for the local OS
make

# run all tests + linter
make test

# run integration tests
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

```
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

```
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

```
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

```
# Use godocdown to create Markdown documentation for all commands and packages
# this is run by default.
make documentation
```

## Releasing

The build process sets the package version based on the latest git tag. After
all changes have been made for the lastest release, make a new tag with NO
commits after, and then `make package-all` to create the artifacts.

This process should be automated someday.

Finally, upload the artifacts on Github to the tag release.


## Contributing
- Submit a pull request for review.
- If it's a code change make sure you run "make test" to confirm it's all good.
- Since Flex has a lot of moving parts, its best to write tests so all contributors don't impact each others work.

## Docker
- Set your configs, modify Dockerfile if need be
- Build & Run Image

```
Example:

BUILD
docker build -t nri-flex .

RUN - standard
docker run -d --name nri-flex --network=host --cap-add=SYS_PTRACE -v "/:/host:ro" -v "/var/run/docker.sock:/var/run/docker.sock" -e NRIA_LICENSE_KEY="yourInfraLicenseKey" nri-flex:latest

RUN - with container discovery reverse lookup (ensure -container_discovery is set to true nri-flex-config.yml)
docker run -d --name nri-flex --network=host --cap-add=SYS_PTRACE -l flexDiscoveryRedis="t=redis,c=redis,tt=img,tm=contains,r=true"  -v "/:/host:ro" -v "/var/run/docker.sock:/var/run/docker.sock" -e NRIA_LICENSE_KEY="yourInfraLicenseKey" nri-flex:latest

Example: Run Redis with a flex discovery label
docker run -it -p 9696:6379 --label flexDiscoveryRedis="t=redis,c=redis,tt=img,tm=contains" --name redis-svr -d redis
```

## Features & Support
- Run any HTTP/S request, read file, shell command, consume from any Prometheus Exporter, Database Query, or JMX Query. (Java 7+ is required for JMX to work)
- Service / Container Discovery
- Attempt to cleverly flatten to samples
- Use environment variables anywhere within config files (eg. using double dollar signs -> $$MY_ENV_VAR)
- Detect and flatten dimensional data from Prometheus style payloads (vector, matrix, targets supported)
- Merge different samples and outputs together
- Key Remover & Replacer
- Metric Parser for RATE & DELTA support (has capability to auto set rates and deltas)
- Define multiple APIs / commands or mix
- event_type autoset or override
- Define custom attributes (more granular control, compared to NR infra agent)
- Command allows horizontal split (useful for table style data) (use only once per command set)
- snake_case to CamelCase conversion
- Percentage to Decimal conversion
- ToLower conversion
- SubParse functionality (see redis config for an example)
- LookUp Store - save attributes from previously generated samples to use in requests later (see rabbit example)
- LazyFlatten - for arrays
- Inbuilt data caching - useful for processing samples at different points
- [+more here](https://github.com/newrelic/nri-flex/wiki/Functions)

## Integrations 
For all see within the examples directory.

- All Prometheus Exporters
- Consul
- Vault (shows merge functionality)
- Bamboo
- Teamcity
- CircleCI
- RabbitMQ (shows metric parser, and lookup store functionality)
- Elasticsearch (shows inbuilt URL cache functionality)
- Traefik
- Kong
- etcd (shows custom sample keys functionality)
- Varnish
- Redis (more metrics, multi instance support, multi db support) (shows snake to camel, perc to decimal, replace keys, rename keys & sub parse functionality)
- Zookeeper
- OpsGenie
- VictorOps
- PagerDuty (shows lazy_flatten functionality)
- AlertOps (shows lazy_flatten functionality)
- New Relic Alert Ingestion (provides similar output to nri-alerts-pipe)
- New Relic App Status Health Ingestion (appSample to present your app health, language, and aggregated summary)
- http/s testing & request performance via curl
- Postgres Custom Querying
- MySQL Custom Querying
- MariaDB Custom Querying
- Percona Server, Google CloudSQL or Sphinx (2.2.3+) Custom Querying
- MS SQL Server Custom Querying
- JMX via nrjmx // (nrjmx is targetted to work with Java 7+, see cassandra and tomcat examples)
- cassandra - via jmx
- tomcat - via jmx
- bind9
- df display disk & inode info (shows horizontal split functionality)

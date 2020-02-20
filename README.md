# New Relic - Flex

[![Build Status](https://travis-ci.org/newrelic/nri-flex.svg?branch=master)](https://travis-ci.com/newrelic/nri-flex)

- [New Relic - Flex](#new-relic---flex)
  - [Requirements](#requirements)
  - [Tutorials](#tutorials)
  - [Visualizing & Managing with the Flex Manager UI](#visualizing--managing-with-the-flex-manager-ui)
  - [Further Documentation](#further-documentation)
  - [Installation](#installation)
  - [Local development](#local-development)
  - [Example integrations](#example-integrations)
  - [Experimental features](#experimental-features)
  - [Disclaimer](#disclaimer)

- Flex is an agnostic AIO New Relic Integration, that can:
  - Abstract the need for end users to write any code other then to define a configuration yml, allowing you to consume metrics from a large variety of services.
  - Run any HTTP/S request, shell command.
  - Can generate New Relic metric samples automatically from almost any HTTP(S) endpoint for almost any payload with useful helper functions to parse and tidy up the output. See <https://github.com/newrelic/nri-flex/tree/master/docs/wiki/apis/functions.md>
  - Provides examples for over [200+ Integrations](#integrations)
  - As updates and upgrades are made, all Flex Integrations reap the benefits.
  - Official support for Linux hosts only at the moment

## Requirements

- Linux
- Windows (experimental support)
- New Relic Infrastructure

## Tutorials

- [Basic step-by-step tutorial](./docs/basic-tutorial.md)

## Visualizing & Managing with the Flex Manager UI

- [Flex Manager](https://github.com/newrelic/nr1-flex-manager)

## Further Documentation

- [Wiki](https://github.com/newrelic/nri-flex/tree/master/docs/wiki)
- [Example Integrations](#exampleintegrations)
- [Available Functions](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/apis/functions.md)
- [Creating your own Flex configuration(s)](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/basics/creating_configs.md)

---

## Installation

Flex is now being integrated by default in the New Relic Infrastructure agent.
See <https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure> for more information

## Local development

If you are in the process of setting up Flex configurations you can use Flex in isolatin mode, ie. without using the New Relic Infrastructure agent.
See [development](./development.md) for more information.

## Example integrations

For all examples see <https://github.com/newrelic/nri-flex/tree/master/examples>

Note: some of these examples may use features that are experimental (not officially supported )and/or are deprecated and may me removed in the future.

See [experimental](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/experimental) and [deprecated](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/experimental))

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
- Cassandra - via jmx
- Tomcat - via jmx
- bind9
- Linux disk usage & inode info (shows horizontal split functionality)

## Experimental features

Flex implements other features apart from those that we officially support. While you can use them if you need to solve your specific use-case, we don't guarantee at this point that they will work as expected and so we offer no official support.

They include, but are not limited to:

- Consume metrics from any Prometheus Exporter
- Consume metrics from database queries
- Consume metrics metrics from JMX queries (Java 7+ is required for JMX to work)
- Consume metrics from csv and json files
- Windows Hosts, Kubernetes, ECS, Fargate, and other container based platforms.
- [Service / Container Discovery](https://github.com/newrelic/nri-flex/wiki/Service-Discovery) built-in.
  This feature is deprecated in favour of the New Relic Infrastructure agent discovery support.
  It is kept for backwards compatibility and may be removed in the future.
- Send data via New Relic Insights Event API. Should be used only in development mode.

## Disclaimer

New Relic has open-sourced this integration to enable monitoring of various technologies. This integration is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to this integration via GitHub. Support for this integration is available with an Expert Services subscription.

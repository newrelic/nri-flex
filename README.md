# New Relic - Flex

[![Build Status](https://travis-ci.org/newrelic/nri-flex.svg?branch=master)](https://travis-ci.com/newrelic/nri-flex)

- [New Relic Flex](#new-relic---flex)
  - [Requirements](#requirements)
  - [Tutorial](#tutorial)
  - [Installation](#installation)
  - [Visualizing and managing Flex](#visualize-and-manage-flex)
  - [Local development](#local-development)
  - [Example integrations](#example-integrations)
  - [Experimental features](#experimental-features)
  - [Further documentation](#further-documentation)
  - [Disclaimer](#disclaimer)


Flex is an agnostic AIO New Relic Integration, that allows you to consume metrics from a wide variety of services, removing the need for end users to write any code other then to define a configuration YAML.

Among other things, Flex can:
  - Run any HTTP request or shell command.
  - Generate New Relic metric samples automatically via [useful helper functions](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/apis/functions.md) that parse and tidy up the output.

Only Linux is officially supported at the moment. As updates and upgrades are made, all Flex Integrations reap the benefits. See the examples for over [200+ Integrations](#integrations)!

## Requirements

- Linux
- Windows (experimental support)
- New Relic Infrastructure

## Tutorial

- [Flex step-by-step tutorial](./docs/basic-tutorial.md)

## Installation

Flex now comes bundled with the New Relic Infrastructure agent. See [Install, configure, and manage Infrastructure](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure).

## Visualize and manage Flex

Use the [Flex manager](https://github.com/newrelic/nr1-flex-manager) in New Relic One to visualize Flex data and manage the Flex integration.

## Local development

If you are setting up Flex configurations, you can use Flex in isolation mode, that is, without using the New Relic Infrastructure agent. See [Development](./development.md) for more information.

## Example integrations

All examples are located in <https://github.com/newrelic/nri-flex/tree/master/examples>.

Some of these examples may use features that are [experimental](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/experimental) (not officially supported) or [deprecated](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/experimental), and may me removed in the future.

- AlertOps (shows `lazy_flatten` functionality)
- All Prometheus Exporters
- Bamboo
- bind9
- Cassandra (via JMX)
- CircleCI
- Consul
- Elasticsearch (shows built-in URL cache functionality)
- etcd (shows custom sample keys functionality)
- HTTP/s testing & request performance via curl
- JMX via nrjmx (nrjmx is targetted to work with Java 7+, see Cassandra and Tomcat examples)
- Kong
- Linux disk usage and inode info (shows horizontal split functionality)
- MariaDB Custom Querying
- MSSQL Server Custom Querying
- MySQL Custom Querying
- New Relic Alert Ingestion (provides similar output to `nri-alerts-pipe`)
- New Relic App Status Health Ingestion (`appSample` to present your app health, language, and aggregated summary)
- OpsGenie
- PagerDuty (shows `lazy_flatten` functionality)
- Percona Server, Google CloudSQL or Sphinx (2.2.3+) Custom Querying
- Postgres Custom Querying
- RabbitMQ (shows metric parser, and lookup store functionality)
- Redis (more metrics, multi instance support, multi db support) (shows snake to camel, perc to decimal, replace keys, rename keys and sub parse functionality)
- Teamcity
- Tomcat - via JMX
- Traefik
- Varnish
- Vault (shows merge functionality)
- VictorOps
- Zookeeper

## Experimental features

Flex implements other features apart from those that we officially support. While you can use them if you need to solve your specific use cases, we can't guarantee that they'll work as expected, and we don't support them officially.

Experimental features include, but are not limited to:

- Consume metrics from any Prometheus Exporter.
- Consume metrics from database queries.
- Consume metrics metrics from JMX queries (Java 7+ is required for JMX to work).
- Consume metrics from csv and json files.
- Windows Hosts, Kubernetes, ECS, Fargate, and other container based platforms.
- [Service / Container Discovery](https://github.com/newrelic/nri-flex/wiki/Service-Discovery) built-in. This feature is deprecated in favour of the New Relic Infrastructure agent discovery support, and kept for backwards compatibility and may be removed in the future.
- Send data via New Relic Insights Event API. Should be used only in development mode.

## Further Documentation

- [NRI-Flex Docs](https://github.com/newrelic/nri-flex/tree/master/docs/wiki)
- [Example Integrations](#exampleintegrations)
- [Available Functions](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/apis/functions.md)
- [Create your own Flex configurations](https://github.com/newrelic/nri-flex/tree/master/docs/wiki/basics/creating_configs.md)

## Disclaimer

New Relic has open-sourced this integration to enable monitoring of various technologies. This integration is provided AS-IS WITHOUT WARRANTY OR SUPPORT, although you can report issues and contribute to this integration via GitHub. Support for this integration is available with an Expert Services subscription.

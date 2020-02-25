# New Relic Flex

[![Build Status](https://travis-ci.org/newrelic/nri-flex.svg?branch=master)](https://travis-ci.com/newrelic/nri-flex)

- [New Relic Flex](#new-relic---flex)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Getting started](#getting-started)
  - [Example integrations](#example-integrations)
  - [Experimental features](#experimental-features)
  - [Development](#development)
  - [Documentation](#documentation)
  - [Support](#support)
  - [License](#license)


Flex is an application-agnostic, all-in-one [New Relic integration](https://docs.newrelic.com/docs/integrations) that allows you to collect metric data from a wide variety of services. You can instrument any app that outputs to the terminal: you create a [config file](/docs/basics/creating_configs.md), start the Infrastructure agent, and data starts pouring into New Relic - see the [200+ example integrations](#example-integrations)!

Flex works in two steps:
  1. It runs any HTTP request or shell command, with or without parameters.
  2. It generates metric samples using [functions](/docs/apis/functions.md) that parse and tidy up the output from the commands/requests.

Only Linux is officially supported at the moment. As updates and upgrades are made, all Flex Integrations reap the benefits.

## Requirements

- Linux
- Windows (experimental support)
- New Relic [Infrastructure Pro](https://newrelic.com/infrastructure/pricing) subscription or trial

## Installation

Flex now comes bundled with the New Relic Infrastructure agent. To install the Infrastructure agent, see [Install, configure, and manage Infrastructure](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure).

## Getting started

The Flex [step-by-step tutorial](./docs/basic-tutorial.md) is a great starting point.

## Example integrations

All examples are located in [/examples](https://github.com/newrelic/nri-flex/tree/master/examples).

Note that some examples may use features that are [experimental](https://github.com/newrelic/nri-flex/tree/master/docs/experimental) (not officially supported) or [deprecated](https://github.com/newrelic/nri-flex/tree/master/docs/experimental), and may be removed in the future.

## Experimental features

Flex implements more features than those that New Relic officially supports. While you can use them to solve specific use cases, we can't guarantee that they'll work as expected, and New Relic does not support them officially.

Experimental features include, but are not limited to:

- Collect metrics from any Prometheus Exporter.
- Collect metrics from database queries.
- Collect metrics metrics from JMX queries (Java 7+ required).
- Collect metrics from CSV and JSON files.
- Support for Windows Hosts, Kubernetes, ECS, Fargate, and other container based platforms.
- Send data via New Relic Insights Event API. Should be used only in development mode.

### Flex Manager

Use the [Flex manager](https://github.com/newrelic/nr1-flex-manager) in New Relic One to visualize Flex data and manage the Flex integration.

## Development

If you are setting up Flex configurations, you can use Flex in isolation mode, that is, without using the New Relic Infrastructure agent. For more information, see [Development](./development.md).

## Documentation

- [Flex documentation - Main page](https://github.com/newrelic/nri-flex/tree/master/docs/readme.md)
- [File and directory structure](https://github.com/newrelic/nri-flex/tree/master/docs/basics/file_layout.md)
- [Create your own Flex configurations](https://github.com/newrelic/nri-flex/tree/master/docs/basics/creating_configs.md)
- [Supported functions](https://github.com/newrelic/nri-flex/tree/master/docs/apis/functions.md)
- [Experimental functions](https://github.com/newrelic/nri-flex/tree/master/docs/experimental/functions.md)

## Support

You can find more detailed documentation [on our website](http://newrelic.com/docs).

If you can't find what you're looking for there, reach out to us on our [support site](http://support.newrelic.com/) or our [community forum](http://forum.newrelic.com) and we'll be happy to help you.

Found a bug? Contact us at [support.newrelic.com](http://support.newrelic.com/), or email support@newrelic.com.

### Community

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub. You can find this project's topic/threads here:

https://discuss.newrelic.com/c/support-products-agents/new-relic-infrastructure

### Issues / Enhancement Requests

Issues and enhancement requests can be submitted in the [Issues tab of this repository](../../issues). Please search for and review the existing open issues before submitting a new issue.


## License

The project is released under version 2.0 of the [Apache license](http://www.apache.org/licenses/LICENSE-2.0).
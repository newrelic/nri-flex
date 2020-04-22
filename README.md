# New Relic Flex

[![Build Status](https://travis-ci.org/newrelic/nri-flex.svg?branch=master)](https://travis-ci.org/newrelic/nri-flex)

Flex is an application-agnostic, all-in-one [New Relic integration](https://docs.newrelic.com/docs/integrations) with which you can instrument any app that exposes metrics over a standard protocol (HTTP, file, shell) in a standard format (for example, JSON or plain text): you create a [config file](/docs/basics/configure.md), start the Infrastructure agent, and data starts pouring into New Relic.

Flex can take any input using [data source APIs](/docs/apis/README.md), process it through [functions](/docs/basics/functions.md), and send metric samples to New Relic as if they came from an integration:

![Flex diagram](https://newrelic-wpengine.netdna-ssl.com/wp-content/uploads/flex_diagram.jpg)

For a quick introduction on Flex, [read our blog post](https://blog.newrelic.com/product-news/how-to-use-new-relic-flex/). You can also have a look at the [200+ example integrations](#example-integrations)!

  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Getting started](#getting-started)
  - [Example integrations](#example-integrations)
  - [Development](#development)
  - [Documentation](#documentation)
  - [Support](#support)
  - [License](#license)

## Compatibility and requirements

Flex requires a New Relic [Infrastructure Pro](https://newrelic.com/infrastructure/pricing) subscription or trial and is compatible with the following operating systems/platforms:

- Linux
- Kubernetes
- Windows (Experimental)

For more information on compatible distros and versions, see the [Infrastructure agent compatibility page](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/getting-started/compatibility-requirements-new-relic-infrastructure).

## Installation

Flex comes bundled with the New Relic Infrastructure agent. To install the Infrastructure agent, see [Install, configure, and manage Infrastructure](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure).

If you're using Kubernetes take a look at [How to monitor services in Kubernetes](https://docs.newrelic.com/docs/integrations/kubernetes-integration/link-apps-services/monitor-services-running-kubernetes).

## Getting started

The [Flex step-by-step tutorial](./docs/basic-tutorial.md) is a great starting point.

## Example integrations

All examples are located in [/examples](https://github.com/newrelic/nri-flex/tree/master/examples).

> Note that some examples may use features that are [experimental](https://github.com/newrelic/nri-flex/tree/master/docs/experimental) (not officially supported) or [deprecated](https://github.com/newrelic/nri-flex/tree/master/docs/experimental).

## Development

While developing your own Flex integrations, you can use Flex without the New Relic Infrastructure agent for debugging. For more information, see [Development](/docs/development.md).

## Documentation

- [Flex documentation - Main page](docs/README.md)
- [Configure Flex](/docs/basics/configure.md)
- [Configure Flex in Kubernetes](/docs/basics/k8s_configure.md)
- [Data sources / APIS](/docs/apis/README.md)
- [Data transformation functions](docs/basics/functions.md)
- [Experimental functions](docs/experimental/functions.md)

### Flex Manager

Use the [Flex manager](https://github.com/newrelic/nr1-flex-manager) in New Relic One to visualize Flex data and manage the Flex integration.

## Support

Need help? See our [troubleshooting page](troubleshooting.md). You can find more detailed documentation [on the New Relic docs site](http://newrelic.com/docs).

If you can't find what you're looking for there, reach out to us on our [support site](http://support.newrelic.com/) or our [community forum](http://forum.newrelic.com) and we'll be happy to help you.

Found a bug? Contact us at [support.newrelic.com](http://support.newrelic.com/)

### Community

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. Like all official New Relic open source projects, there's a related Community topic in the New Relic Explorers Hub. You can find this project's topic/threads here:

https://discuss.newrelic.com/c/support-products-agents/new-relic-infrastructure

### Issues / Enhancement Requests

Issues and enhancement requests can be submitted in the [Issues tab of this repository](../../issues). Please search for and review the existing open issues before submitting a new issue.

## License

The project is released under version 2.0 of the [Apache license](http://www.apache.org/licenses/LICENSE-2.0).

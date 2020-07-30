[![Community Project header](https://github.com/newrelic/open-source-office/raw/master/examples/categories/images/Community_Project.png)](https://github.com/newrelic/open-source-office/blob/master/examples/categories/index.md#category-community-project)

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
- Windows
- Kubernetes

For more information on compatible distros and versions, see the [Infrastructure agent compatibility page](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/getting-started/compatibility-requirements-new-relic-infrastructure).

## Install

Flex comes bundled with the New Relic Infrastructure agent. To install the Infrastructure agent, see [Install, configure, and manage Infrastructure](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure).

If you're using Kubernetes, see [Configure Flex in Kubernetes](https://github.com/newrelic/nri-flex/blob/master/docs/basics/k8s_configure.md).

## Getting started

The [Flex step-by-step tutorial](./docs/basic-tutorial.md) is a great starting point.

## Example integrations

All examples are located in [/examples](https://github.com/newrelic/nri-flex/tree/master/examples).

> Note that some examples may use features that are [experimental](https://github.com/newrelic/nri-flex/tree/master/docs/experimental) (not officially supported) or [deprecated](https://github.com/newrelic/nri-flex/tree/master/docs/experimental).

## Development

While developing your own Flex integrations, you can use Flex without the New Relic Infrastructure agent for debugging. For more information, see [Development](/docs/development.md).

## Documentation

- [Flex usage documentation](docs/README.md)
- [How to configure Flex](/docs/basics/configure.md)
- [Flex under Kubernetes](/docs/basics/k8s_configure.md)
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

## Privacy

At New Relic we take your privacy and the security of your information seriously, and are committed to protecting your information. We must emphasize the importance of not sharing personal data in public forums, and ask all users to scrub logs and diagnostic information for sensitive information, whether personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or identifiable individual, including, for example, your name, phone number, post code or zip code, Device ID, IP address and email address.

Please review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy) for more information.

## Contributing

We encourage your contributions to improve Flex! Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.

If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company, drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](/SECURITY.md), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, please review [these guidelines](./CONTRIBUTING.md).

To all contributors, we thank you!  Without your contribution, this project would not be what it is today.

## License

The project is released under version 2.0 of the [Apache license](http://www.apache.org/licenses/LICENSE-2.0).

# Flex documentation

## Get started

- [Basic tutorial](basic-tutorial)
- [Configure Flex](configure)
- [Functions for data manipulation](functions)

## Data sources

With Flex you can acquire data from multiple sources for processing:

- [commands](./apis/commands): standard output from command-line tools
- [url](./apis/url): JSON output from HTTP/HTTPS endpoints

More data sources will be added in future updates. 

## Experimental features

Flex implements the following experimental features. 'Experimental' here means that New Relic does not yet provides support for them.

- [Experimental functions](./experimental/functions)
- [Database queries](./experimental/db)
- [Net dial](./experimental/dial)
- [Git configuration synchronization](./experimental/git_sync)

## Deprecated features

The following functionalities are still provided by Flex for backwards compatibility, but its use is discouraged and unsupported because New Relic provides more convenient implementations of such functionalities.

For each deprecated functionality, please consider migrating to the New Relic supported equivalent, as linked in the right column of the following table. 

| Deprecated feature | New Relic supported equivalent |
|---|---|
| [Discovery](./deprecated/discovery) | [Container auto-discovery for On-Host Integrations](https://docs.newrelic.com/docs/integrations/host-integrations/installation/container-auto-discovery) |
| [JMX](./deprecated/jmx) | [New Relic JMX On-Host Integration](http://github.com/newrelic/nri-jmx) |
| [Prometheus](./deprecated/prometheus) | [New Relic Prometheus OpenMetrics integration for Docker and Kubernetes](https://docs.newrelic.com/docs/integrations/prometheus-integrations) |
| [Secrets management](./deprecated/secrets) | [Secrets management for On-Host Integrations](https://docs.newrelic.com/docs/integrations/host-integrations/installation/secrets-management) |
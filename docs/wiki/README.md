# NRI-Flex documentation

## Basics

- [Flex files layout](basics/file_layout.md)
- [Creating configs](basics/creating_configs.md)
- [Order of operations](basics/order_of_operations.md)

## APIs

The Flex APIs provide means to acquire data from multiple sources, for its later
[manipulation](apis/functions.md)

- [`commands`](apis/commands.md)
- [`url`](apis/url.md)

Please refer to the [Functions for data manipulation](apis/functions.md) section for
a reference of all the functions available to manipulate the data that is acquired by
the APIs.

## Experimental functionalities

Flex implements the following functionalities, but they are still experimental. This means
that New Relic does not (yet) provides customer support for them.

- [Git configuration synchronization](experimental/git_sync.md)

## Deprecated functionalities

The following functionalities are still provided by Flex for backwards compatibility, but
its use is discouraged and unsupported because New Relic provides more convenient implementations
of such functionalities.

For each deprecated functionality, please consider migrating to the New Relic supported equivalent,
as linked in the right column of the following table. 

| Deprecated functionality | New Relic supported equivalent |
|---|---|
| [Discovery](deprecated/discovery.md) | [Container auto-discovery for On-Host Integrations](https://docs.newrelic.com/docs/integrations/host-integrations/installation/container-auto-discovery) |
| [Secrets management](deprecated/secrets.md) | [Secrets management for On-Host Integrations](https://docs.newrelic.com/docs/integrations/host-integrations/installation/secrets-management) |
| [Prometheus](deprecated/prometheus.md) | [New Relic Prometheus OpenMetrics integration for Docker and Kubernetes](https://docs.newrelic.com/docs/integrations/prometheus-integrations) |



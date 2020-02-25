#  File and directory structure

The following file/directory structure applies when using Flex with the New Relic Infrastructure Agent for Linux.

Note that within the OHI Config File, you can spawn multiple instances of Flex with different configurations if needed.

## Flex executable

New Relic Infrastructure version 1.10.0 or higher already bundles Flex with the agent in the same package.

The Flex executable is in `/var/db/newrelic-infra/newrelic-integrations/bin/nri-flex`.

## Flex and OHI - Joint configuration

Flex is configured in the same folder as the rest of [on-host integrations](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180): `/etc/newrelic-infra/integrations.d`.

The main section of the configuration YAML file is used by the Infrastructure agent to execute Flex like any other on-host integration; under `config`, you can add your Flex configuration, which is only used by Flex.

For example, a `/etc/newrelic-infra/integrations.d/my-flex-config.yml` file could contain:

```yaml
integrations:
  - name: nri-flex
    interval: 60s
    timeout: 5s
    config:
      name: example
      apis:
        - event_type: ExampleSample
          url: https://my-host:8443/admin/metrics.json
```

* **On-Host Integration Configuration**: the first 5 lines are read by the agent to execute the `nri-flex` binary every 60 seconds, canceling the execution if it lasts more than 5 seconds. Refer to [Integration configuration file specifications](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180) for more details about the contents of the on-host integrations configuration file.
* **Flex configuration**: the agent writes the four last lines into a temporary YAML file that is read by the `nri-flex` command (the file being set in the `CONFIG_PATH` environment variable). It contains the actions to be taken. Refer to the [APIs](../README.md#apis) reference for examples and details.

## Flex and OHI - Separate configuration

The agent allows linking the Flex file path from `/etc/newrelic-infra/integrations.d` by replacing the `config` contents of the `config_template_path` property, which contains the path of the Flex config file.

This way, the equivalent of `/etc/newrelic-infra/integrations.d/my-flex-config.yml` from the previous section would contain:

```yaml
integrations:
  - name: nri-flex
    interval: 60s
    timeout: 5s
    config_template_path: /path/to/flex/integration.yml
```

And the `/path/to/flex/integration.yml` file would contain the contents that previously were inside the `config`
section:

```yaml
name: example
apis:
  - event_type: ExampleSample
    url: https://my-host:8443/admin/metrics.json
```

> We recommend splitting the configuration of Flex as an on-host integration (that is, how the Infrastructure agent must execute Flex) and the actual Flex configuration in two separate files.

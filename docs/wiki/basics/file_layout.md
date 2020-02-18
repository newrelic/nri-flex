# Files Layout

When using Flex with the New Relic Infrastructure Agent, this is the file/directory structure that
should appear for Linux.

Within the OHI Config File, we could spawn multiple instances of Flex with different configurations if needed.

## Flex executable

The Flex executable will be available in the path `/var/db/newrelic-infra/newrelic-integrations/bin/nri-flex`.

New Relic Infrastructure agent version 1.10.0 and higher bundles Flex within the agent package, so you don't
need to perform any extra step for its installation.

## Flex & OHI joint configuration

Flex can be configured in the same folder as the rest of [on-host integrations](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180):
`/etc/newrelic-infra/`.

You need to put there a YAML file that may have two parts: the main document is used by the Infrastructure agent to
execute Flex as any other on-host integration (user, interval, label decoration, etc...); in the `config` section
of the main document, you can write your Flex configuration contents, which is exclusively used by Flex for its
proper operation.

For example, let a `/etc/newrelic-infra/my-flex-config.yml` text file contain:

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

* **On-Host Integration Configuration**: the first 5 lines are read by the Agent to know that it must execute the `nri-flex` executable every 60 seconds,
  canceling the execution if it lasts more than 5 seconds.
    - Please refer to the [on-host integrations configuration documentation](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180)
      for more details about the contents of the on-host integrations' configuration file.
* **Flex configuration**: the Agent writes the 4 last lines into a temporary YAML file that is read by the `nri-flex`
  command (receiving the file path via the `CONFIG_PATH` environment variable), containing the actions
  that need to be taken.
    - Please refer to the [APIs](../README.md#apis) reference for more examples and details of the contents of the
      a Flex configuration.

## Flex & OHI separate configuration

You may feel more comfortable keeping on one side the configuration of Flex as an on-host integration (this is,
how the agent must execute flex), and on the other side the configuration of what Flex has to do, in two separate
files.

The agent allows linking the Flex file path from the configuration YML in `/etc/newrelic-infra`, by replacing
the `config` property contents by the `config_template_path` property containing the path of the Flex file.

This way, the equivalent `/etc/newrelic-infra/my-flex-config.yml` from the previous section would contain:

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


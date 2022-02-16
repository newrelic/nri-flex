# Configure Flex

Flex integrations are code-less: They only require that you write configuration files for each integration as YAML files.

Configurations can [go embedded](#AddyourFlexconfigurationtointegrations.d) in the Infrastructure agent configuration file, or live as stand-alone config files that you can test separately and [link](#Linktoaseparateconfigurationfile) from the main integrations config file. What approach to follow is up to you.

- [Add your Flex configuration to `integrations.d`](#AddyourFlexconfigurationtointegrations.d)
- [Link to a separate configuration file](#Linktoaseparateconfigurationfile)
- [Configuration schema](#Configurationschema)
- [Configuration example](#Configurationexample)

## <a name='AddyourFlexconfigurationtointegrations.d'></a>Add your configuration to `integrations.d`

Since it comes bundled with the Infrastructure agent, Flex's configuration must be stored as YAML in the same folder as the rest of [on-host integrations](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180):

- Linux: `/etc/newrelic-infra/integrations.d`
- Windows: `C:\Program Files\New Relic\newrelic-infra\integrations.d\`

The main section of the integrations config file is used by the Infrastructure agent to execute Flex like any other integration; you can add your Flex configuration under `config`.

For example, `/etc/newrelic-infra/integrations.d/my-flex-config.yml` could contain the following:

```yaml
integrations:
  - name: nri-flex # We're telling the Infra agent to run Flex
    interval: 60s
    timeout: 5s
    config: # Flex configuration starts here!
      name: example
      apis:
        - event_type: ExampleSample
          url: https://my-host:8443/admin/metrics.json
```

- **On-Host Integration Configuration**: the first five lines are read by the agent to execute the `nri-flex` binary every 60 seconds, canceling the execution if it lasts more than 5 seconds. Refer to [Integration configuration file specifications](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180) for more details about the contents of the on-host integrations configuration file.
- **Flex configuration**: contains the actions to be taken using data sources APIs, such as the [url](../apis/url.md) API.

> To get a quick, first picture of a Flex configuration file, you can start following our [basic, step-by-step tutorial](../../basic-tutorial.md) or check existing config files under [/examples](../../examples).

## <a name='Linktoaseparateconfigurationfile'></a>Link to a separate configuration file

You can store the Flex configuration in a separate YAML file (for example, after [developing and testing a config file](../development.md)) and reference it by replacing `config` with `config_template_path` property, which contains the path of the Flex config file.

This way, the equivalent of `/etc/newrelic-infra/integrations.d/my-flex-config.yml` from the previous section would contain:

```yaml
integrations:
  - name: nri-flex
    interval: 60s
    timeout: 5s
    config_template_path: /path/to/flex/integration.yml # Reference to a separate Flex config file
```

`/path/to/flex/integration.yml` would contain the contents that previously were inside the `config` section:

```yaml
name: example
apis:
  - event_type: ExampleSample
    url: https://my-host:8443/admin/metrics.json
```

## <a name='Configurationschema'></a>Configuration schema

The following schema describes the overall structure of a Flex configuration.

```
+----------------------+
| name                 |
| global?              |
| +--------------+     |
| | <properties> |     |   Suffixes:
| +--------------|     |       ? optional
| custom_attributes?   |       * multiple repetitions
| +----------------+   |
| | <key>: <value> | * |
| +----------------+   |
| apis                 |
| +---------------+    |
| |  name?        |    |
| |  event_type?  | *  |
| |  <api>        |    |
| |  <functions>* |    |
| +---------------+    |
+----------------------+
```

### <a name='name'></a>name

The name of the Flex configuration. It should be something short and meaningful.

### <a name='global'></a>global

Set of global properties that apply to the overall file. The aim of this section is to avoid repeating some values (like URLs or user credentials).

**Example**:

```yaml
global:
  base_url: http://localhost:9200/
  user: elastic
  pass: 3l4st1c
  headers:
    accept: application/json
```

These are all the possible `global` properties:

| Property       | Description                                                                                                                    |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `base_url`     | Base URL. See [specifying a common base URL](../apis/url.md#SpecifyacommonbaseURL)                                             |
| `user`         | Username for APIs that require user and password authentication                                                                |
| `pass`     | Password for APIs that require user and password authentication                                                                    |
| `pass_phrase`  | Pass phrase for encrypte `password` properties                                                                                 |
| `proxy`        | Proxy URL for APIs whose connections require it                                                                                |
| `timeout`      | Timeout for the API connections, in milliseconds                                                                               |
| `headers`      | Key-value map of headers for the HTTP/HTTPS connections                                                                        |
| `tls_config`   | TLS configuration. See [configuring your HTTPS connections](../apis/url.md#ConfigureyourHTTPSconnections)                      |
| `ssh_pem_file` | Path to PEM file to enable SSH authentication                                                                                  |
| `JMX`          | See [JMX](../experimental/jmx.md) (experimental)                                                                               |

### <a name='apis'></a>apis

The `apis` section allows you to define multiple entries for data acquisition and processing. Each entry must have a `name` or `event_type`, which is used to name the event type in New Relic:

- `event_type` provides a name for each sample, which is used as table name for querying the metrics
  in the New Relic UI. `event_type` usually have names like `MySQLSample`, `MyRemoteSample`, `FolderSample`, etc.
- If `event_type` is not defined and `name` is, the submitted event type is `name`
  with the `Sample` prefix concatenated.
  - For example, `name: FolderSize` would make Flex to create events named `event_type: FolderSizeSample`.
- Each `event_type` is automatically decorated by the infrastructure agent with a `timestamp` attribute.
  Flex can overwrite this attribute, but it's not recommended as it can impact the alerts configured for this `event_type`.

In addition to the fields that define the name of the sample, each `apis` entry requires the type of API to parse data from, and, optionally, a list of [functions](../basics/functions.md) for processing the data coming from the API.

### <a name='Cache'></a>Cache

Flex by default stores the result of an API execution in it's internal cache. You can then use this cache as input to another API for further processing.

For example, consider a service that returns the following payload

```json
{
  "id": "eca0338f4ea31566",
  "leaderInfo": {
    "leader": "8a69d5f6b7814500",
    "startTime": "2014-10-24T13:15:51.186620747-07:00",
    "uptime": "10m59.322358947s",
    "abc": {
      "def": 123,
      "hij": 234
    }
  },
  "name": "node3"
}
```

as a result of executing the following `url` API:

```yaml
name: example
apis:
  - name: someService
    url: http://some-service.com/status
```

As we want to process it in another API, we use the `cache` function. Note that the cache `key` is the URL because it's an `url` API:

```yaml
name: example
apis:
  - name: status
    url: http://some-service.com/status
  - name: otherStatus
    cache: http://some-service.com/status
    strip_keys:
      - id
      - name
```

With a `commands` API, you should use the name of the API instead:

```yaml
name: example
apis:
  - name: status
    commands:
      # assume that this file contains the same json payload showed above the beginning
      - run: cat /var/some/file
  - name: otherStatus
    cache: status
    strip_keys:
      - id
      - name
```

### <a name='Customattributes'></a>Custom attributes

With Flex you can add your own custom attributes to samples. Add any custom attribute using key-value pairs under the `global` directive, and at the API level by declaring an array named `custom_attributes`.

```yaml
custom_attributes:
  greeting: hello
```

Custom attributes can be defined nearly anywhere in your configuration. For example, under `global`, or `api`, or further nested under each command.
Attributes defined at the lowest level take precedence.

Custom attributes defined at the `global` level are added to all samples, while custom attributes defined at the API level are added only at the level of the API where they are defined.

### <a name='Environmentvariables'></a>Environment variables

You can inject values for environment variables anywhere in a Flex config file. To inject the value for an environment variable, use a double dollar sign before the name of the variable (for example `$$MY_ENVIRONMENT_VAR`).

## <a name='Configurationexample'></a>Configuration example

Here's an example of a Flex configuration embedded in the main integrations configuration file:

```yaml
integrations: # OHI configuration starts here
  - name: nri-flex # OHI to be executed by the Agent
    config: # OHI configuration to be parsed by Flex
      # Actual Flex configuration starts here
      name: linuxDirectorySize # Flex configuration name
      apis:
        - name: DirectorySize # Event type will be DirectorySizeSample
          commands: # Selecting the API `commands`
            - run: du -c $$DIR # Running a shell command
              split: horizontal # Post-processing function: split horizontally
              set_header: [dirSizeBytes, dirName] # Names for the headers of the table resulting from split
              regex_match: true # Splits horizontally matching a regular expression
              split_by: (\d+)\s+(.*) # Captures the regexpes between parentheses as the headers above
```

### <a name='specialenvs'></a>Special environment variables

- `FLEX_META` - setting this environment variable with a flat JSON payload will unpack this into global custom attributes.
- `FLEX_CMD_PREPEND` - automatically prepend to commands being run (Requires AllowEnvCommands enabled).
- `FLEX_CMD_APPEND` - automatically append to a commands being run (Requires AllowEnvCommands enabled).
- `FLEX_CMD_WRAP` - automatically wrap the command in quotes (Requires AllowEnvCommands enabled).

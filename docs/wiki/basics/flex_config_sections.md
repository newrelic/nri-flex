# Anatomy of a Flex configuration file

This document describes the sections that compose a Flex configuration file. To get a quick, first picture of
a Flex configuration file, you can start following our [basic, step-by-step tutorial](../../basic-tutorial.md).

Flex configurations are all defined by a YAML file. As described in the [File Layout](./file_layout.md) page,
it can be written in two different ways:
 
1. As part of an [on-host integration (OHI) configuration file](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180).
   E.g. the `/etc/newrelic-infra/integrations.d/my-bundled-config.yml` would contain:
   ```yaml
   integrations:
     - name: nri-flex
       config:
         <Flex configuration YAML>
   ```
2. Referenced from the OHI configuration file, by means of the `config_template_path` option. E.g.
   the `/etc/newrelic-infra/integrations.d/my-bundled-config.yml` would contain:
   ```yaml
   integrations:
     - name: nri-flex
       config_template_path: /path/to/flex-config.yml
   ```
   And the YAML file in `/path/to/flex-config.yml` would contain the actual Flex configuration file.

This document page focus on the Flex configuration YAML sections. For the OHI configuration options, please
read the [OHI configuration file specification](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180).

The following schema depicts the overall structure of a Flex configuration (the one that should go inside the `config`
OHI configuration, or the file referenced from the `config_template_path` property).

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

The rest of this document describes the above sections of the Flex configuration.

## name

The name of the Flex configuration. It should be something short and meaningful.

## global

Set of global properties that would apply to the overall file. The aim of this section
is to avoid repeating some values (e.g. URLs, user credentials...) when they need to be
used from multiple places.

Example:

```yaml
global:
  base_url: http://localhost:9200/
  user: elastic
  pass: 3l4st1c
  headers:
    accept: application/json
```

The following table enumerated all the possible global properties:

| Property | Description |
|---|---|
| `base_url` | See [specifying a common base URL](../apis/url.md#specifying-a-common-base-url-with-base_url) |
| `user` | If you are using an API that requires user and password authentication, the user name |
| `password` | If you are using an API that requires user and password authentication, the password |
| `proxy` | If your chosen API requires connecting through a proxy, the URL of the actual proxy |
| `timeout` | Timeout for the API connections, in milliseconds |
| `headers` | Key/value map of headers for the HTTP/HTTPS connections. |
| `tls_config` | See [configuring your HTTPS connections](../apis/url.md#configuring-your-https-connections-with-tls_config) |
| `JMX` | See [JMX](../experimental/jmx.md) (experimental) |
| `pass_phrase` | If the above `password` property is ciphered, a pass phrase to decipher it |
| `ssh_pem_file` | Path to a PEM file to enable SSH authentication. |  

## custom_attributes

The `custom_attributes` accepts any key/values map, and allows decorating the samples with the
contained values. Example:

```yaml
custom_attributes:
  environment: production
  role: database
```

Custom attributes can be defined nearly anywhere in your configuration. E.g. under `global`, or `api`,
or further nested under each command. The lowest level defined attribute will take precedence.

## apis

The `apis` section allows you defining multiple entries for data acquisition and processing. Each enty needs to have
a `name` or `event_type` entry, which will be used to provide the name of the event type in infrastructure:

* `event_type` provides a name for each sample, which will be used as table name for querying the metrics
  in the New Relic UI. `event_type` would usually have names like `MySQLSample`, `MyRemoteSample`, `FolderSample`...
* If `event_type` is not defined and `name` is, the submitted event type will be the `name`
  with the `Sample` prefix concatenated.
    - E.g. `name: FolderSize` would make Flex creating events named with `event_type: FolderSizeSample`

In addition to the fields that define the name of the sample, each `apis` entry will require the type of API to
parse data from, and optionally a list of [functions](../apis/functions.md) for processing the data from the API.

Currently supported APIs are:

* [`commands`](../apis/commands.md) to execute a shell command and use its standard output as source
  of metrics (usually to be processed by a list of [functions](../apis/functions.md). 
* [`url`](../apis/url.md) to retrieve data from an HTTP or HTTPS endpoint.

## Example

An example of Flex configuration file (embedded in the OHI configuration) would be like:

```yaml
integrations:                                    # OHI configuration starts here  
  - name: nri-flex                               # OHI to be executed by the Agent
    config:                                      # OHI configuration to be parsed by Flex
      # Actual Flex configuration starts here
      name: linuxDirectorySize                   # Flex configuration name
      apis:                                       
        - name: DirectorySize                    # Event type will be DirectorySizeSample
          commands:                              # Selecting the API `commands`
            - run: du -c $$DIR                   # Running a shell command
              split: horizontal                  # Post-processing function: split horizontally
              set_header: [dirSizeBytes,dirName] # Names for the headers of the table resulting from split
              regex_match: true                  # Split horizontally matching a regular expression
              split_by: (\d+)\s+(.*)             # Capture the regexpes between parentheses as the headers above   
```

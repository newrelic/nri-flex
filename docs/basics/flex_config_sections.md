# Structure of a Flex configuration file

> To get a quick, first picture of a Flex configuration file, you can start following our [basic, step-by-step tutorial](../../basic-tutorial.md).

Flex configurations are written in YAML. They can be created in two ways, as described in the [File Layout](./file_layout.md) page:
 
1. As part of an [on-host integration (OHI) configuration file](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180).
   For example,  the `/etc/newrelic-infra/integrations.d/my-bundled-config.yml` would contain:
    
   ```yaml
   integrations:
     - name: nri-flex
       config:
         <Flex configuration YAML>
   ```
2. Referenced from the OHI configuration file, by means of the `config_template_path` option. 
   For example, the `/etc/newrelic-infra/integrations.d/my-bundled-config.yml` would contain:
   ```yaml
   integrations:
     - name: nri-flex
       config_template_path: /path/to/flex-config.yml
   ```
   While `/path/to/flex-config.yml` would contain the actual Flex configuration file.

Here we focus on creating the Flex configuration YAML file. For OHI configuration settings, see the [OHI configuration file specification](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180).

## Overview

The following schema describes the overall structure of a Flex configuration file (the one that should go inside the `config` OHI configuration, or the file referenced in `config_template_path`).

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

## name

The name of the Flex configuration. It should be something short and meaningful.

## global

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

| Property | Description |
|---|---|
| `base_url` | Base URL. See [specifying a common base URL](../apis/url.md#specifying-a-common-base-url-with-base_url) |
| `user` | Username for APIs that require user and password authentication |
| `password` | Password for APIs that require user and password authentication |
| `pass_phrase` | Pass phrase for encrypte `password` properties  |
| `proxy` | Proxy URL for APIs whose connections require it |
| `timeout` | Timeout for the API connections, in milliseconds |
| `headers` | Key-value map of headers for the HTTP/HTTPS connections |
| `tls_config` | TLS configuration. See [configuring your HTTPS connections](../apis/url.md#configuring-your-https-connections-with-tls_config) |
| `ssh_pem_file` | Path to PEM file to enable SSH authentication |  
| `JMX` | See [JMX](../experimental/jmx.md) (experimental) |

## apis

The `apis` section allows you to define multiple entries for data acquisition and processing. Each entry must have a `name` or `event_type`, which is used to name the event type in New Relic:

* `event_type` provides a name for each sample, which is used as table name for querying the metrics
  in the New Relic UI. `event_type` usually have names like `MySQLSample`, `MyRemoteSample`, `FolderSample`, etc.
* If `event_type` is not defined and `name` is, the submitted event type is `name`
  with the `Sample` prefix concatenated.
    - For example, `name: FolderSize` would make Flex to create events named `event_type: FolderSizeSample`.

In addition to the fields that define the name of the sample, each `apis` entry requires the type of API to parse data from, and, optionally, a list of [functions](../apis/functions.md) for processing the data coming from the API.

For a list of currently supported APIs, see [`Officially supported APIs`](creating_configs.md#OfficiallysupportedAPIs).

## Cache

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
## Example

An example of a Flex configuration file (embedded in the OHI configuration):

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
              regex_match: true                  # Splits horizontally matching a regular expression
              split_by: (\d+)\s+(.*)             # Captures the regexpes between parentheses as the headers above   
```
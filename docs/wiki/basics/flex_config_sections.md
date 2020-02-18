# Sections of a Flex configuration file

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

An example of Flex configuration file (embedded in the OHI configuration) would be like:

```yaml
integrations: # OHI configuration starts here  
  - name: nri-flex
    config:
      # Flex configuration starts here
      name: linuxDirectorySize
      apis:
        - name: linuxDirectorySize
          commands:
            - run: du -c $$DIR
              split: horizontal
              set_header: [dirSizeBytes,dirName]
              regex_match: true
              split_by: (\d+)\s+(.*)
```

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

**TODO**: describe each entry or link to the proper documents.

```go
type Global struct {
	BaseURL    string `yaml:"base_url"`
	User, Pass string
	Proxy      string
	Timeout    int
	Headers    map[string]string `yaml:"headers"`
	Jmx        JMX               `yaml:"jmx"`          // Don't explain here. Link to the JMX doc
	TLSConfig  TLSConfig         `yaml:"tls_config"`   // Not explaining here. Link to the url api doc
	Passphrase string            `yaml:"pass_phrase"`
	SSHPEMFile string            `yaml:"ssh_pem_file"`
}
```

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

The `apis` section allows you defining multiple entries for data acquisition and processing. Each



### Options
- [commands](#commands) Run any standard commands
- [net dial](#net-dial) Can be used for port testing or sending messages and processing the response
- [http](#http) General http requests
- [database queries](#database-queries)
- [using prometheus exporters](https://github.com/newrelic/nri-flex/wiki/Prometheus-Integrations-(Exporters))

### Further Configuration

#### [Functions available for things like pagination, manipulating the output, secret mgmt etc.](https://github.com/newrelic/nri-flex/wiki/Functions)
#### [Metric Parser for Rate & Delta Support](https://github.com/newrelic/nri-flex/wiki/Functions#metric_parser)
#### [Global Config](#global-config-that-is-passed-down)
#### [Setting Custom Attributes](#custom-attributes)
#### Environment variables can be used throughout any Flex config files by simply using a double dollar sign eg. $$MY_ENVIRONMENT_VAR.

***


### Commands

With the below example, we can create a redis integration in 6 lines, by simply running a command and parsing it.

```
---
name: redisFlex
apis: 
  - name: redis
    commands: 
      - run: (printf "info\r\n"; sleep 1) | nc -q0 127.0.0.1 6379 ### remove -q0 if testing on mac
        split_by: ":"
```


#### Run Command Specific Options
```
"shell"             // command shell
"run"               // command to run
"split"             // default "vertical", can be "horizontal" useful for outputs that look like a table
"split_by"          // character to split by
"set_header"        // manually set header column names (used when split is is set to horizontal)
"group_by"          // group by character
"regex"             // process SplitBy as regex (true/false)
"line_limit"        // stop processing at this line number
"row_header"        // start the row header at a different line (integer, used when split is horizontal)
"row_start"         // start creating samples from this line number, to be used with SplitBy
"ignore_output"     // ignore command output - useful chaining commands together
"custom_attributes" // set additional custom attributes
"line_end"          // stop processing at this line number
"timeout"           // when to timeout command in milliseconds (default 10s)
"dial"              // address to dial
"network"           // network to use (default tcp) (currently only used for dial)

```
See the redis example for a typical split, and look at the "df" command example for a horizontal split by example.

***

### Net Dial

Dial is a parameter used under commands.

port test eg.
```
name: portTestFlex
apis: 
  - timeout: 1000 ### default 1000 ms increase if you'd like
    commands:
    - dial: "google.com:80"
```

sending a message and processing the output eg.
```
---
name: redisFlex
apis: 
  - name: redis
    commands: 
      - dial: 127.0.0.1:6379
        run: "info\r\n"
        split_by: ":"
```

#### Global Config that is passed down
```
base_url
user
pass
proxy
timeout
headers:
 headerX: valueX
jmx:
* domain
* user
* pass
* host
* port
* key_store
* key_store_pass
* trust_store
* trust_store_pass
```

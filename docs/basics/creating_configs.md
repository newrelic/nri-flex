# Create Flex configurations

Flex allows you to define multiple APIs or data sources so you can monitor multiple services with just one Flex configuration file.

Flex configuration files are in YAML format. The easiest way to start learning how to configure Flex is by checking existing config files under [/examples](https://github.com/newrelic/nri-flex/tree/master/examples).

* [Officially supported APIs](#OfficiallysupportedAPIs)
* [Alpha APIs](#AlphaAPIs)
* [Data parsing and transformation functions](#Dataparsingandtransformationfunctions)
* [Custom attributes](#Customattributes)
* [Environment variables](#Environmentvariables)
* [Global configuration options](#Globalconfigurationoptions)

##  <a name='OfficiallysupportedAPIs'></a>Officially supported APIs

Althoug Flex supports a variety of APIs, the following APIs are officially supported by New Relic.

- [Shell commands](../apis/commands.md): Run any standard shell command or application (Linux or Windows).
- [HTTP requests](../apis/url.md): Query any standard HTTP/HTTPS endpoint returning JSON or XML.

More APIs will be added in the future.

## <a name='AlphaAPIs'></a>Alpha APIs

The following APIs are in alpha status, and while they may work for your usecase, New Relic does not yet support them. 

- [Net dial](#net-dial): Can be used for port testing or for sending messages and processing responses.
- [Database queries](#database-queries)
- [JMX queries](#jmx-queries): Uses nrjmx to send JMX requests to be processed.

##  <a name='Dataparsingandtransformationfunctions'></a>Data parsing and transformation functions

For data parsing and transformation functions like key renaming, key removal, output splitting, and others, see [Functions](../apis/functions.md).

##  <a name='Customattributes'></a>Custom attributes

With Flex you can add your own custom attributes to samples. Add any custom attribute using key-value pairs under the `global` directive, and at the API level by declaring an array named `custom_attributes`.

```yaml
custom_attributes:
  greeting: hello
```

Custom attributes defined at the config level are added to all samples, while custom attributes defined at the API level are added only at the level of the API where they are defined. Custom attributes can also be added for each command declared under `commands`.

##  <a name='Environmentvariables'></a>Environment variables

You can inject values for environment variables anywhere in a Flex config file. To inject the value for an environment variable, use a double dollar sign before the name of the variable (for example `$$MY_ENVIRONMENT_VAR`).

##  <a name='Globalconfigurationoptions'></a>Global configuration options

Flex allows you to configure some properties at global level. These properties overwrites any default (if applicable) and some can also be overwritten at API level.

| Name | Type | Applies to | Description |
|---:|:---:|:---:|:---|
|`base_url`| string | `url` | Base URL for HTTP requests. Allows to define URLs using just path segments |
|`user`| string | `url` | Default user for Basic authentication when doing HTTP requests |
|`pass`| string | `url` | Default password for Basic authentication when doing HTTP requests |
|`proxy`| string | `url` | Default proxy to use when doing HTTP requests |
|`timeout`| string | `url` | Default timeout for HTTP requests |
|`headers`| string to string map | `url` | Default headers to send when doing HTTP requests |
|`tls_config`| structure | `url` | TLS configuration to use whe sending HTTP request. See [http requests](../apis/url.md) for more information |

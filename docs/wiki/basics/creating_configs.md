# Creating Flex Configs

Flex configurations are all defined by a yaml file.
The easiest way to kick start learning is to check the existing configs under [/examples](https://github.com/newrelic/nri-flex/tree/master/examples).

Flex allows you to define multiple APIs or data sources so you can monitor multiple services with just one Flex configuration file.

Each API outputs it's resulting data into a single sample. A sample is a data structure that contains a name, some metadata and metric attributes. This data structure is then serialized and sent to New Relic.

Some APIs, like [commands](#commands), allow you to define multiple sequential entries that work in conjuction to process data.

## Officially supported APIs

Althoug Flex supports a variety of APIs, the following APIs are considered officially supported by New Relic. More APIs will be gaining GA status in the not so distant future.

- [shell commands](../apis/commands.md) Run any standard shell command or application
- [http requests](../apis/url.md) Query any standard http(s) endpoint returning JSON or XML

## Alpha APIs

The following APIs are considered Alpha status and while they may work for your use-case, New Relic does not at the moment offer official support for them.

- [net dial](#net-dial) Can be used for port testing or sending messages and processing the response.
- [database queries](#database-queries)
- [jmx queries](#jmx-queries) Uses nrjmx to send JMX requests which will be processed further.

## Data parsing and transformation functions

Data parsing and transformation functions available for features like `key renaming`, `key removal`, `output splitting` and others, see [Functions](../apis/functions.md).

## Custom Attributes

Flex allows you to add your own custom attributes to the resulting sample of an API. You can add custom attributes under the `Global` directive and at the API level, by declaring an array named `custom_attributes`.

Custom attributes defined at the Config level will be added to all samples, and custom attributes defined at the API level will be added only in the specific API where they are defined.

Custom attributes can also be added for the in the `commands` API in each declared command.

A standard key:pair structure is used.

```yaml
custom_attributes:
  greeting: hello
```

## Environment variables

Flex allows you to inject environment variable values throughout any Flex config file. Simply use a double dollar sign before the name of the environment value you want injected (ex: $$MY_ENVIRONMENT_VAR)

## Global configuration options

Flex allows you to configure some properties at the a global level. These properties  will overwrite the defaults (if applicable) and some can also be overwritten at the API level.

| Name | Type | Default | Applies to | Description |
|---:|:---:|:---:|:---:|:---|
|`base_url`| string | n.a. | `url` | Set it to define the base `url` used for  http requests. If using a `base_url` you can then use just the path of the endpoint when doing HTTP requests |
|`user`| string | n.a. | `url` | Set it do define the default `user` to use for Basic authentication when doing HTTP requests |
|`pass`| string | n.a. | `url` | Set it do define the default `password` to use for Basic authentication when doing HTTP requests |
|`proxy`| string | n.a. | `url` | Set it do define the default `proxy` to use when doing HTTP requests |
|`timeout`| string | n.a. | `url` | Set it do define the default `timeout` to use when doing HTTP requests |
|`headers`| map of string to string | n.a. | `url` | Set it do define the default `headers` to send when doing HTTP requests |
|`tls_config`| structure | n.a. | `url` | Set it do define the TLS configuration to use whe nding HTTP request. See [http requests](../apis/url.md) for more information |

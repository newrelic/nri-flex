## Flex OHI Config Options

> ⚠️ **Notice** ⚠️: the following documents may contain some deprecated functionalities that
are still supported by New Relic for backwards compatibility. However, an updated version of this
document is in progress. 

# update status sample examples

* Note this config file `nri-flex-config.yml` is a config file required by the New Relic Infrastructure Agent, to run integrations. 
    * This is different from the Flex config files, that are used to write and run Flex built integrations, that you would see within the examples/flexConfigs folder.

* Running `./nri-flex -help` will show that there are a fair amount of options available, the important ones are covered further below.
    * These options can be set in the arguments section of your `nri-flex-config.yml` file
    * Or can be used on the command line, for testing manually:
        * `./nri-flex -config_file "abc.yml" -verbose -pretty`

### Config Options

```
integration_name: com.newrelic.nri-flex

instances:
  - name: nri-flex
    command: metrics
    arguments:
      # container_discovery: true ## default false
      # container_discovery_dir: "anotherDir" default "flexContainerDiscovery" 
      # config_file: "../myConfigFile.yml" ## default "" - run only a single specific config file
      # config_dir: "anotherConfigDir/" ## default "flexConfigs/"
      # event_limit: 500 ## default 500
      # insights_api_key: abc
      # insights_url: https://insights...
      # insights_output: output the payload to stdout
    # labels:
      # owner: cloud
```

### Running Multiple Instances of Flex

```
integration_name: com.newrelic.nri-flex

instances:
  - name: nri-flex-a
    command: metrics
    arguments:
      config_dir: "databaseConfigs/"
  - name: nri-flex-b
    command: metrics
    arguments:
      config_dir: "shellConfigs/"
```
---
## Event Limiter

* There is a built in event limiter, as Flex can enable anyone to develop and run integrations in a trivial amount of time, you can also create a stack of events in the process. To safe guard NR systems and costs for the account, a configurable event limiter is built - `event_limit`. It defaults to a 500 event limit per execution, which can be bumped up if needed.
    * Flex also generates a flexStatusSample event, each run which indicates the volume of events being created and configs processed.
    * If any events get dropped due to the event limiter, they can be identified here.
```
eg. SELECT * FROM flexStatusSample
{
  "event_type": "flexStatusSample",
  "flex.Hostname": "Z09W63Z5DTY0",
  "flex.IntegrationVersion": "Unknown-SNAPSHOT",
  "flex.counter.ConfigsProcessed": 1,
  "flex.counter.EventCount": 0,
  "flex.counter.EventDropCount": 0
}
```
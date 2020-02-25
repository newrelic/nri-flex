# Configure Flex

> ⚠️ **Notice**: the following documents may contain deprecated functionalities that are still provided for backwards compatibility. 

`nri-flex-config.yml` is a config file required by the New Relic Infrastructure agent to run integrations. This is different from the Flex configuration files, which are used to write and run Flex integrations.

Running `./nri-flex -help` will show that there are a fair amount of options available. Options can be set in the arguments section of your `nri-flex-config.yml` file or can be used on the command line, for manual testing:

```bash
./nri-flex -config_file "abc.yml" -verbose -pretty
```

## Settings

```yaml
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

## Run multiple instances of Flex

```yaml
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

## Event limiter

Flex includes an event limiter, as anyone can develop and run integrations in a trivial amount of time, and generate a huge stack of events in the process. 

To avoid flooding your New Relic account with events, use `event_limit`. It defaults to a 500 event limit per execution, which can be increased if needed.

Flex also generates a `flexStatusSample` event, each run of which indicates the volume of events being created and configs processed.

Use this query to identify events dropped due to the event limiter:
```sql
SELECT * FROM flexStatusSample
{
  "event_type": "flexStatusSample",
  "flex.Hostname": "Z09W63Z5DTY0",
  "flex.IntegrationVersion": "Unknown-SNAPSHOT",
  "flex.counter.ConfigsProcessed": 1,
  "flex.counter.EventCount": 0,
  "flex.counter.EventDropCount": 0
}
```
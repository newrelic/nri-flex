# Troubleshooting Flex 

## Command Line arguments

The Infrastructure agent installs bundled binaries together in the following directory: 

```shell
# Linux
/var/db/newrelic-infra/newrelic-integrations/bin/

# Windows
C:\Program Files\New Relic\newrelic-infra\newrelic-integrations\
```

And expects the associated configuration files to be placed in the following directory:

```shell
# Linux
/etc/newrelic-infra/integrations.d/

# Windows
C:\Program Files\New Relic\newrelic-infra\integrations.d\
```

Flex itself comes with several command line arguments that can assist you during development of your configurations as you shape your metric outputs. You can run the following command to see all of the options available to you:

```shell
# Linux
/var/db/newrelic-infra/newrelic-integrations/bin/nri-flex -help

# Windows
C:\Program Files\New Relic\newrelic-infra\newrelic-integrations\nri-flex.exe -help
```

As you can see, there are well over 40 unique command line arguments that are available. The majority are for use when Flex is run without the Infrastructure agent, or are related to [experimental features](https://github.com/newrelic/nri-flex/tree/master/docs/experimental). But there are three that you should be aware of for troubleshooting purposes:
  * `config_path` - Set a specific config file to execute
  * `pretty` - Print pretty formatted JSON
  * `verbose` - Print more information to logs

These arguments allow you to see in real-time how Flex is processing the config you send it and the associated data output from the API(s) in the config.

### Testing a config

You can manually test a config file to ensure the output meets your expectations by running a command like this, replacing `<FILE_NAME>` with the name of your config file: 

```shell
# Linux
./nri-flex -verbose -pretty -config_path /etc/newrelic-infra/integrations.d/<FILE_NAME>

# Windows
.\nri-flex.exe -verbose -pretty -config_path "C:\Program Files\New Relic\newrelic-infra\integrations.d\<FILE_NAME>"
```

This will give you an output similar to this: 

```shell
INFO[0000] com.newrelic.nri-flex                         GOARCH=amd64 GOOS=linux version=1.3.5
DEBU[0000] config: git sync configuration not set       
WARN[0000] config: testing agent config, agent features will not be available 
DEBU[0000] config: running async                         name=pingTest
DEBU[0000] config: processing apis                       apis=1 name=pingTest
DEBU[0000] fetch: collect data                           name=pingTest
DEBU[0000] commands: executing                           count=1 name=pingTest
DEBU[0004] processor-data: running data handler          name=pingTest
DEBU[0004] config: finished variable processing apis     apis=1 name=pingTest
INFO[0004] flex: completed processing configs            configs=1
{
	"name": "com.newrelic.nri-flex",
	"protocol_version": "3",
	"integration_version": "1.3.5",
	"data": [
		{
			"metrics": [
				{
					"avg": 17.491,
					"event_type": "pingTest",
					"flex.commandTimeMs": 4029,
					"integration_name": "com.newrelic.nri-flex",
					"integration_version": "1.3.5",
					"max": 17.55,
					"min": 17.425,
					"packetLoss": 0,
					"packetsReceived": 5,
					"packetsTransmitted": 5,
					"stddev": 0.151,
					"timeMs": 4006,
					"url": "google.com"
				},
				{
					"event_type": "flexStatusSample",
					"flex.Hostname": "<REDACTED>.us-east-2.compute.internal",
					"flex.IntegrationVersion": "1.3.5",
					"flex.counter.ConfigsProcessed": 1,
					"flex.counter.EventCount": 1,
					"flex.counter.Event
					
					
				DropCount": 0,
					"flex.counter.pingTest": 1,
					"flex.time.elaspedMs": 4031,
					"flex.time.endMs": 1609978782729,
					"flex.time.startMs": 1609978778698
				}
			],
			"inventory": {},
			"events": []
		}
	]
}
```

This has 4 major sections: 

#### Verbose logging

```shell
INFO[0000] com.newrelic.nri-flex                         GOARCH=amd64 GOOS=linux version=1.3.5
DEBU[0000] config: git sync configuration not set       
WARN[0000] config: testing agent config, agent features will not be available 
DEBU[0000] config: running async                         name=pingTest
DEBU[0000] config: processing apis                       apis=1 name=pingTest
DEBU[0000] fetch: collect data                           name=pingTest
DEBU[0000] commands: executing                           count=1 name=pingTest
DEBU[0004] processor-data: running data handler          name=pingTest
DEBU[0004] config: finished variable processing apis     apis=1 name=pingTest
INFO[0004] flex: completed processing configs            configs=1
```

The intro section shows a log of how `nri-flex` executed your config file. If you have a malformed YAML file, you'll oftentimes see the last line show a zero entry like this: 
`INFO[0004] flex: completed processing configs            configs=0`

#### Versioning Info

```json
{
	"name": "com.newrelic.nri-flex",
	"protocol_version": "3",
	"integration_version": "1.3.5",
	"data": [
		{
			"metrics": [
```

This next section leads into the more interesting part of the payload, but you can quickly see the version of `nri-flex` you're executing with here if needed.

#### Config Execution Output

```json
{
					"avg": 17.491,
					"event_type": "pingTest",
					"flex.commandTimeMs": 4029,
					"integration_name": "com.newrelic.nri-flex",
					"integration_version": "1.3.5",
					"max": 17.55,
					"min": 17.425,
					"packetLoss": 0,
					"packetsReceived": 5,
					"packetsTransmitted": 5,
					"stddev": 0.151,
					"timeMs": 4006,
					"url": "google.com"
				},
```

The next JSON object(s) will represent the output of `nri-flex` executing the supplied configuration file. This is where you will see error output, if applicable. It's a representation of the raw data that will be shipped into New Relic One's Telemetry Data Platform. 

#### Flex Status Sample Output

```json
{
					"event_type": "flexStatusSample",
					"flex.Hostname": "<REDACTED>.us-east-2.compute.internal",
					"flex.IntegrationVersion": "1.3.5",
					"flex.counter.ConfigsProcessed": 1,
					"flex.counter.EventCount": 1,
					"flex.counter.Event
					
					
				DropCount": 0,
					"flex.counter.pingTest": 1,
					"flex.time.elaspedMs": 4031,
					"flex.time.endMs": 1609978782729,
					"flex.time.startMs": 1609978778698
				}
```

The last major section shows the `flexStatusSample` event. This is a heartbeat event that is sent along with every successful execution of `nri-flex` and can be used to evaluate whether a problem lies in your config, or with the Flex binary itself.

## Common issues

Flex is pretty forgiving, but there may be times that the data you aimed at capturing won't show up in New Relic. There may be several reasons to this. Here are the most common, by category.

* [Configuration issues](#Configurationissues)
	* [Wrong location of the data source](#Wronglocationofthedatasource)
	* [Empty or non-responsive data source](#Emptyornon-responsivedatasource)
	* [You asked Flex to perform something impossible](#YouaskedFlextoperformsomethingimpossible)
	* [Bad indentation](#Badindentation)
* [Agent issues](#Agentissues)
	* [Stand-alone configuration files](#Stand-aloneconfigurationfiles)
	* [Bad license key or invalid license](#Badlicensekeyorinvalidlicense)
* [System issues](#Systemissues)
	* [The Infrastructure agent is not running](#TheInfrastructureagentisnotrunning)
	* [Your system cannot run the Infrastructure agent](#YoursystemcannotruntheInfrastructureagent)
	* [There's a communication problem between the agent and New Relic](#TheresacommunicationproblembetweentheagentandNewRelic)

### <a name='Configurationissues'></a>Configuration issues

#### <a name='Wronglocationofthedatasource'></a>Wrong location of the data source

Check that the path, URI, or resource locator is correct. Flex can't guess if there's a typo in your configuration.

#### <a name='Emptyornon-responsivedatasource'></a>Empty or non-responsive data source

While the path of the data source may be correct, it may not contain the data you expected. If it's an HTTP endpoint, try using `curl` to verify it's responsive (that is, you get a `200 OK`). Commands may not always be available or you may be lacking permissions to run them. Files can be empty or protected. 

Always explore the availability of your sources prior to using Flex.

#### <a name='YouaskedFlextoperformsomethingimpossible'></a>You asked Flex to perform something impossible

While Flex is pretty flexible about things, you may have instructed it to perform something impossible, such as splitting records by colons when the file only contains comma-separated values, or apply a regex filter that returns zero matches. 

Check the [functions](/basics/functions.md) documentation for instructions on how to best use each Flex feature.

#### <a name='Badindentation'></a>Bad indentation

YAML files are notoriously picky with indentation. Lack of indentation, or using tabs instead of spaces, can cause Flex to reject any configuration file. We recommend that you two-space indentation for each level. Do not use tabs or other characters -- they're forbidden by the YAML specification. Use a YAML linter or checker to verify that your YAML conforms to the standard.

### <a name='Agentissues'></a>Agent issues

#### <a name='Stand-aloneconfigurationfiles'></a>Stand-alone configuration files

The Infrastructure agent needs to know that your Flex configurations must be run using Flex. Add the following lines at the top of your Flex configuration file and indent the rest accordingly...

```yaml
integrations:
  - name: nri-flex
    config:
      #Your Flex config goes here
```

Or reference the Flex stand-alone configuration file from the main Infrastructure agent configuration file:

```yaml
integrations:
  - name: nri-flex
    config_template_path: /path/to/flex/integration.yml # It can be any folder
```

#### <a name='Badlicensekeyorinvalidlicense'></a>Bad license key or invalid license

Check that your [New Relic license key](/docs/accounts/install-new-relic/account-setup/license-key) is correct and active. You won't be able to send data to New Relic otherwise, though you can still use Flex in [development mode](development.md).

### <a name='Systemissues'></a>System issues

#### <a name='TheInfrastructureagentisnotrunning'></a>The Infrastructure agent is not running

Flex will not send data to New Relic if the Infrastructure agent is not running. Check whether it's running or not by running these commands:

* Linux: `systemctl status newrelic-infra` or `pgrep newrelic-infra`
* Windows: `net status newrelic-infra`

If the agent is not running, learn how to [start, stop, restart, or check the agent status](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/configuration/start-stop-restart-check-infrastructure-agent-status).

#### <a name='YoursystemcannotruntheInfrastructureagent'></a>Your system cannot run the Infrastructure agent

The New Relic Infrastructure agent is a highly optimized piece of software, but still requires that the host meets certain requirements and is compatible. Check the [Infra agent requirements](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/getting-started/compatibility-requirements-new-relic-infrastructure).

Flex only works when the agent is [running in root mode](https://docs.newrelic.com/docs/infrastructure/install-configure-infrastructure/linux-installation/linux-agent-running-modes).

#### <a name='TheresacommunicationproblembetweentheagentandNewRelic'></a>There's a communication problem between the agent and New Relic

Verify that your network connection is up and that nothing is preventing the Infrastructure agent from communicating with the New Relic platform. 

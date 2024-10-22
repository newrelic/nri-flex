# Standalone mode

> **Disclaimer**: this function is bundled as alpha. That means that it is not yet supported by New Relic.

Flex can run without the infrastructure agent reporting telemetry directly to New Relic [Event API](https://docs.newrelic.com/docs/data-apis/ingest-apis/event-api/introduction-event-api/). 

## Steps to setup Flex standalone
1. Get the Flex executable on your target host. If the infrastructure agent is installed, use the default installation paths:
  * Linux: `/var/db/newrelic-infra/newrelic-integrations/bin/nri-flex`
  * Windows: `C:/Program Files/New Relic/newrelic-infra/newrelic-integrations/nri-flex.exe`
  * macOS: Flex is not bundled with the infrastructure agent on macOS.

    You can also download the lastest package from [Github releases](https://github.com/newrelic/nri-flex/releases) and extract it in your target host.

2. Prepare the execution command with the expected parameters as [explained below](#running-flex-without-the-infrastructure-agent). 

3. Determine when and how Flex should be executed. For example, define a cronjob to execute it on a fixed interval. 

4. Check that the custom events are reported as expected in the New Relic UI.

Refer to [troubleshooting](https://docs.newrelic.com/docs/infrastructure/host-integrations/host-integrations-list/flex-integration-tool-build-your-own-integration/#troubleshooting) documentation for help. 

## Running Flex without the infrastructure agent

Command example:
```
  /path/to/nri-flex -insights_api_key YOUR_LICENSE_KEY -insights_url https://insights-collector.newrelic.com/v1/accounts/YOUR_ACCOUNT_ID/events -config_path /path/to/YOUR_CONFIG.yml
```

* Adapt the path to point to your local executable for `nri-flex`
* Replace `YOUR_LICENSE_KEY` with your ingest key.
* Replace `YOUR_ACCOUNT_ID` with your New Relic account. 
* For EU accounts, use `insights-collector.eu01.nr-data.net`.
* Adapt the path to point to your local Flex configuration / integration.


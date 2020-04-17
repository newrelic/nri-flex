# Troubleshooting Flex issues

Flex is pretty forgiving, but there may be times that the data you aimed at capturing won't show up in New Relic. There may be several reasons to this. Here are the most common, by category.

## Configuration issues

### Wrong location of the data source

Check that the path, URI, or resource locator is correct. Flex can't guess if there's a typo in your configuration.

### Empty or non-responsive data source

While the path of the data source may be correct, it may not contain the data you expected. If it's an HTTP endpoint, try using `curl` to verify it's responsive (that is, you get a `200 OK`). Commands may not always be available or you may be lacking permissions to run them. Files can be empty or protected. 

Always explore the availability of your sources prior to using Flex.

### You asked Flex to perform something impossible

While Flex is pretty flexible about things, you may have instructed it to perform something impossible, such as splitting records by colons when the file only contains comma-separated values, or apply a regex filter that returns zero matches. 

Check the [functions](/basics/functions.md) documentation for instructions on how to best use each Flex feature.

### Bad indentation

YAML files are notoriously picky with indentation. Lack of indentation, or using tabs instead of spaces, can cause Flex to reject any configuration file. We recommend that you two-space indentation for each level. Do not use tabs or other characters -- they're forbidden by the YAML specification. Use a YAML linter or checker to verify that your YAML conforms to the standard.

## Agent issues

### Stand-alone configuration files

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

### Bad license key or invalid license

Check that your [New Relic license key](/docs/accounts/install-new-relic/account-setup/license-key) is correct and active. You won't be able to send data to New Relic otherwise, though you can still use Flex in [development mode](development.md).

## System issues

### The Infrastructure agent is not running

Flex will not send data to New Relic if the Infrastructure agent is not running. Check whether it's running or not by running these commands:

* Linux: `systemctl status newrelic-infra` or `pgrep newrelic-infra`
* Windows: `net status newrelic-infra`

If the agent is not running, learn how to [start, stop, restart, or check the agent status](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/configuration/start-stop-restart-check-infrastructure-agent-status).

### Your system cannot run the Infrastructure agent

The New Relic Infrastructure agent is a highly optimized piece of software, but still requires that the host meets certain requirements and is compatible. Check the [Infra agent requirements](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/getting-started/compatibility-requirements-new-relic-infrastructure).

Flex only works when the agent is [running in root mode](https://docs.newrelic.com/docs/infrastructure/install-configure-infrastructure/linux-installation/linux-agent-running-modes).

### There's a communication problem between the agent and New Relic

Verify that your network connection is up and that nothing is preventing the Infrastructure agent from communicating with the New Relic platform. 

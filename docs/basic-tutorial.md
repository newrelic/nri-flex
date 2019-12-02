# FLEX step-by-step basic tutorial

TODO: windows options

Requirements:
* Infrastructure Agent 1.8.0 (older versions also work with Flex, but the new integrations configuration
  system makes Flex easier to configure)
* Run as root mode

## Install

### Automatic

Say that the new Flex version will come bundled with the agent, and Agent 1.8.0 works better with Flex
for configuration.

From [Compiled releases](https://github.com/newrelic/nri-flex/releases)

```
$ wget https://github.com/newrelic/nri-flex/releases/download/v0.8.3-pre/nri-flex-linux-0.8.3-pre.tar.gz
$ tar xvf nri-flex-linux-0.8.3-pre.tar.gz
```

Folder structure:

```
nri-flex-linux-0.8.3-pre
|-- examples
|   |-- flexConfigs
|   |-- flexContainerDiscovery
|   |-- fullConfigExamples
|   |   |-- containerDiscovery
|   |   |   `-- flexContainerDiscovery
|   |   `-- standard
|   |       `-- flexConfigs
|   `-- lambdaExample
|       `-- pkg
|           `-- flexConfigs
`-- nrjmx
```

Run the ./install_linux.sh (or `install_windows.bat`) script.

### Manual
1 - copy nri-flex to /var/db/newrelic-infra/newrelic-integrations/

## Configure infra-agent

From [Compiled releases](https://github.com/newrelic/nri-flex/releases)

```
$ wget https://github.com/newrelic/nri-flex/releases/download/v0.8.3-pre/nri-flex-linux-0.8.3-pre.tar.gz
$ tar xvf nri-flex-linux-0.8.3-pre.tar.gz
```

cp
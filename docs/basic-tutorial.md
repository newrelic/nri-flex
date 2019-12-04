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

Run the ./install_linux.sh (or `install_windows.bat`) script.

Just copy the executable to newrelic-integrations
```
$ wget https://github.com/newrelic/nri-flex/releases/download/v0.8.4-pre/nri-flex-linux-0.8.4-pre.tar.gz
$ tar xvf nri-flex-linux-0.8.4-pre.tar.gz
$ cd nri-flex-linux-0.8.4-pre
$ sudo ./install_linux.sh
```



```sql
from flexStatusSample  select * where entityName = '#9033841498900681752'
```


El PROBLEMA con embedded config: FLEX no reconoce archivos que no estén como YAML
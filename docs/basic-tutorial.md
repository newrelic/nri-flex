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

**TODO** create a install_standalone.sh and install_standalone.ps1 scripts.

## Simple case: your mounted filesystems 


```yaml
integrations:
  - name: nri-flex
    interval: 15s
    config:
      name: my-flex-integration
```

```sql
from flexStatusSample  select * where entityName = '#9033841498900681752' LIMIT 1
```

|  |  |
|-----------|----------------|
| agentName | Infrastructure |
| agentVersion | 1.8.2 |
| coreCount | 1 |
| criticalViolationCount | 0 |
| entityGuid | MTAxNTQyOXxJTkZSQXxOQXw5MDMzODQxNDk4OTAwNjgxNzUy |
| entityId | 9033841498900681752 |
| entityKey | mmacias-vagrant |
| entityName | #9033841498900681752 |
| event_type | flexStatusSample |
| flex.Hostname | localhost.localdomain |
| flex.IntegrationVersion | 0.8.4-pre-1-g032d7a0 |
| flex.counter.ConfigsProcessed | 1  |
| flex.counter.EventCount | 2 |
| flex.counter.EventDropCount | 0 |
| flex.counter.FileSystemSample | 2 |
| flex.pd.1 | {\"name\":\"systemd\", "\"cmd\":\"/usr/lib/systemd/systemd ...  |
| flex.pd.1021 | {\"name\":\"master\", "\"cmd\":\"/usr/libexec/postfix/master ...  |
| flex.pd.3386 | {\"name\":\"sshd\", "\"cmd\":\"sshd: vagrant [priv] ... |
| flex.pd.6865 | {\"name\":\"newrelic-infra\", "\"cmd\":\"/usr/bin/newrelic-infra ... |
| flex.pd.7000 | {\"name\":\"sshd\", "\"cmd\":\"sshd: vagrant [priv] ... |
| flex.pd.871 | {\"name\":\"sshd\", "\"cmd\":\"/usr/sbin/sshd -D -u0 ... |
| flex.time.elaspedMs | 33 |
| flex.time.endMs | 1575553487305 |
| flex.time.startMs | 1575553487272 |
| instanceType | unknown |
| kernelVersion | 3.10.0-1062.7.1.el7.x86_64 |
| linuxDistribution | CentOS Linux 7 (Core) |
| nr.entityType | HOST |
| nr.ingestTimeMs | 1575553487000 |
| operatingSystem | linux |
| processorCount | 1 |
| systemMemoryBytes | 510697472 |
| timestamp | 1575553487000 |
| warningViolationCount | 0 |

```yaml
integrations:
  - name: nri-flex
    interval: 15s
    config:
      name: linuxFileSystem
      apis:
        - name: FileSystem
          commands:
            - run: 'df -PT -B1 -x tmpfs -x xfs -x vxfs -x btrfs -x ext -x ext2 -x ext3 -x ext4'
              split: horizontal
              set_header: [fs,fsType,capacityBytes,usedBytes,availableBytes,usedPerc,mountedOn]
              regex_match: true
              split_by: (\S+.\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(.*)
          perc_to_decimal: true
```

```
SELECT average(availableBytes / 1000000) AS AvailableMB From FileSystemSample TIMESERIES FACET mountedOn
```
  
## Using HTTP services

E.g. elastic search:

```
docker run -d --name elasticsearch -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.5.0
```

```yaml
integrations:
  - name: nri-flex
    interval: 15s
    config:
      name: elasticsearch-monitor
      global:
          base_url: http://localhost:9200/
          headers:
            accept: application/json
      apis:
        - event_type: elasticsearchNodeSample
          url: _nodes/stats
          sample_keys:
            nodes: nodes>node.id # there are objects within objects distinguished by key, we can create samples like so
          rename_keys:
            _nodes: parentNodes # bring the parent node attributes within the node sample
          remove_keys:
            - ingest.pipelines.xpack
            - roleSampleSamples
            - fs. ### no need for file system usage stats if infra collects it already
```

```
FROM elasticsearchNodeSample SELECT *
```

## Paso 2: añadir remote entity

```
integrations:
  - name: nri-flex
    interval: 15s
    env:
      LOCAL: false
      ENTITY: elastic-container
    config:
      name: elasticsearch-monitor
      global:
          base_url: http://localhost:9200/
          headers:
            accept: application/json
      apis:
        - event_type: elasticsearchNodeSample
          url: _nodes/stats
          sample_keys:
            nodes: nodes>node.id # there are objects within objects distinguished by key, we can create samples like so
          rename_keys:
            _nodes: parentNodes # bring the parent node attributes within the node sample
          remove_keys:
            - ingest.pipelines.xpack
            - roleSampleSamples
            - fs. ### no need for file system usage stats if infra collects it already
```

```
FROM elasticsearchNodeSample SELECT average(jvm.mem.heap_used_percent) AS HeapPercent  TIMESERIES FACET entityName 
```

## Paso 3: añadir discovery



## Moving your configuration outside

* Automatic: drop into `/var/db/newrelic-infra/custom-integrations/flexConfigs`
  - By now, only usable for Flex internal discovery mechanism
* Explicit: one `nri-flex` entry with a `config_template_path` per file
  - Usable From Agent discovery

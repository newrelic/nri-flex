# Flex step-by-step tutorial

Follow this tutorial to get started with Flex!

## Requirements

* Infrastructure agent version 1.8 or higher.
  - Flex can also work with older versions, but this tutorial relies on the latest integrations engine which has been added in the [version 1.8.0](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180).  
* Run the Infrastructure agent in root/administrator mode. 
  - Current versions of Flex require administrator permissions for the management of temporary files.
* Flex 0.8.5 or higher.
  - This version is prepared to work with the new integrations system introduced by the Infrastructure agent 1.8.0, which is used in this tutorial.
  - Previous versions of Flex also work, but require few extra configuration steps that are not addressed by this tutorial.

## Installation

Starting from New Relic Infrastructure agent version 1.10.0, Flex is bundled with the agent in the same package, so you don't need to perform any extra step for its installation.

## Check that Flex is up and running

1. Create a file named `my-flex-configs.yml` (or similar) in this folder:
    * Linux: `/etc/newrelic-infra/integrations.d`
    * Windows: `C:\Program Files\New Relic\newrelic-infra\integrations.d`
2. Edit the file and add this snippet:
   ```yaml
   integrations:
     - name: nri-flex
       config:
         name: just-testing
   ```
3. Go to New Relic and run the following [NRQL query](https://docs.newrelic.com/docs/query-data/nrql-new-relic-query-language):

```sql
FROM flexStatusSample SELECT * LIMIT 1
```

The query should produce a table similar to this:

![](./img/basic-table.png)

### What happened behind the scenes

1. The Infrastructure agent detected that a new integration, `nri-flex`, has been added.
2. The agent looks for an executable named `nri-flex` in `/var/db/newrelic-infra/newrelic-integrations/`.
3. A temporary configuration file is created with this content:
   ```yaml
   name: just-testing
   ```
4. `nri-flex` is executed and gets the path of the config file via the `CONFIG_PATH` environment
   variable.
5. Flex recognizes a configuration named `just-testing`, but since it does not provide extra information
   it just returns a `flexStatusSample` with some internal status of the Flex integration.

## Your first Flex integration

This example shows how to collect disk metrics from file systems not natively supported by New Relic using the `df` command in a Linux environment. A similar result could be achieved in Windows with just a few changes.

The goal of Flex is to process the output of the `df` command, showing the file system and 1-byte blocks, while excluding file systems already supported by the agent. If unsupported file system are not mounted, remove the `-x` arguments.

```bash
$ df -PT -B1 -x tmpfs -x xfs -x vxfs -x btrfs -x ext -x ext2 -x ext3 -x ext4
Filesystem     Type         1-blocks         Used    Available Capacity Mounted on
devtmpfs       devtmpfs    246296576            0    246296576       0% /dev
go_src         vboxsf   499963170816 361339486208 138623684608      73% /go/src
``` 

We want Flex to convert the above tabular text output into a set of equivalent JSON samples with this format. Notice that the agent decorates each sample with extra fields:

```json
{
  "event": {
    "event_type": "FileSystemSample",
    "fs": "go_src",
    "fsType": "vboxsf",
    "capacityBytes": 499963170816,
    "usedBytes": 361345331200,
    "availableBytes": 138617839616,
    "usedPerc": 73,
    "mountedOn": "/go/src"
  }
}
```

First, you need to tell Flex how to perform the above "table text to JSON" transformation by specifying the following:

- Name of the metric: `FileSystem`. Flex appends the `Sample` suffix, which results in  `FileSystemSample`.
- Which command to run: `df -PT -B1 ...`
- How to split the output table from `df` and how to assign the values to given metric names.

This is achieved by placing the content below in the YAML configuration file:

```yaml
integrations:
  - name: nri-flex
    config:
      name: linuxFileSystemIntegration
      apis:
        - name: FileSystem
          commands:
            - run: 'df -PT -B1 -x tmpfs -x xfs -x vxfs -x btrfs -x ext -x ext2 -x ext3 -x ext4'
              split: horizontal
              split_by: \s+
              set_header: [fs,fsType,capacityBytes,usedBytes,availableBytes,usedPerc,mountedOn]
          perc_to_decimal: true
```

Some sections from the above YAML are worth mentioning: 

- The `apis` section is an array of entries for each sample. Each entry sets a name for the sample, as well as the commands/procedures to get it and how to process it.
- The first entry in `apis` is named `FileSystem`: it is used to build the `FileSystemSample` event name.
- In the `commands` section, we specify how to get the information:
    - `run: 'df -PT -B1...` specifies the command to run.
    - `split: horizontal` states that each output line may return a metric.
    - `split_by` explains how to split each line in different fields. In this case, we use the `\s+` regular expression, which tells Flex that any sequence of one or more white spaces is a separator.
    - `set_header` specifies, in order, a matching name for each value of the aforementioned array.
    - `perc_to_decimal: true` aims at converting any percentage string into a decimal value, removing the trailing `%` symbol, if it exists.

**Once the Flex config is created, the Infrastructure agent auto-detects the new config and begins collecting data.**

To check that our new integration is working, try executing the following NRQL query:

```sql
FROM FileSystemSample SELECT mountedOn, fs, usedBytes, capacityBytes, usedBytes
```

![](./img/basic-filesystem.png)

## For more examples

Check the [flex configs examples](../examples/flexConfigs) folder for more working examples of Flex.

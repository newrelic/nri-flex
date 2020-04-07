# Flex step-by-step tutorial

Follow this tutorial to get started with Flex!

1. [Install the Infrastructure agent](#InstalltheInfrastructureagent)
2. [Check that Flex is up and running](#CheckthatFlexisupandrunning)
3. [Your first Flex integration](#YourfirstFlexintegration)
4. [How to add more integrations](#Howtoaddmoreintegrations)
5. [What's next?](#Whatsnext)

##  1. <a name='InstalltheInfrastructureagent'></a>Install the Infrastructure agent

Starting from New Relic Infrastructure agent version 1.10.7, Flex comes bundled with the agent. To install the Infrastructure agent, see:

- [Install Infrastructure for Linux using the package manager](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure/linux-installation/install-infrastructure-linux-using-package-manager)
- [Install Infrastructure for Windows Server using the MSI installer](https://docs.newrelic.com/docs/infrastructure/install-configure-manage-infrastructure/windows-installation/install-infrastructure-windows-server-using-msi-installer)

You can [start, stop, restart, and check](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/configuration/start-stop-restart-check-infrastructure-agent-status) the Infrastructure agent from the command line. The agent must run in [root/administrator mode](https://docs.newrelic.com/docs/infrastructure/install-configure-infrastructure/linux-installation/linux-agent-running-modes).

##  2. <a name='CheckthatFlexisupandrunning'></a>Check that Flex is up and running

1. Navigate to the integrations folder of the Infrastructure agent:
    * Linux: `/etc/newrelic-infra/integrations.d`
    * Windows: `C:\Program Files\New Relic\newrelic-infra\integrations.d\`
2. Create the integration configuration file (for example, `integrations.yml`) if it doesn't exist.
2. Add the Flex configuration to the file.

   ```yaml
   integrations:
     - name: nri-flex
       config:
         name: just-testing
   ```
   If you already have an integrations section in the file, add `nri-flex` to it.

After a few minutes, go to New Relic and run the following [NRQL query](https://docs.newrelic.com/docs/query-data/nrql-new-relic-query-language):

```sql 
FROM flexStatusSample SELECT * LIMIT 1
```

The query should produce a table similar to this:

![](./img/basic-table.png)

##  3. <a name='YourfirstFlexintegration'></a>Your first Flex integration

This example shows how to collect disk metrics from file systems not natively supported by New Relic using the `df` command in Linux.

The goal of Flex is to process the output of the `df` command, showing the file system and 1-byte blocks, while excluding file systems already supported by the agent. If unsupported file systems are not mounted, remove the `-x` arguments.

```bash
$ df -PT -B1 -x tmpfs -x xfs -x vxfs -x btrfs -x ext -x ext2 -x ext3 -x ext4
Filesystem     Type         1-blocks         Used    Available Capacity Mounted on
devtmpfs       devtmpfs    246296576            0    246296576       0% /dev
go_src         vboxsf   499963170816 361339486208 138623684608      73% /go/src
``` 

We want Flex to convert the above tabular text output into a set of equivalent JSON samples with the following format. Notice that the agent decorates each sample with extra fields:

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

- Name of the metric: `FileSystem`
- Which command to run: `df -PT -B1 ...`
- How to split the output table from `df`
- How to assign the values to given metric names

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
              row_start: 1
              set_header: [fs,fsType,capacityBytes,usedBytes,availableBytes,usedPerc,mountedOn]
          perc_to_decimal: true
```

- `apis` is an array of entries for each sample. Each entry sets a `name`for the sample, as well as the commands/procedures to get and process the sample. The first entry in the example is named `FileSystem`, which is used to name the `FileSystemSample` event.
- `commands` specifies how to get the information from CLI applications:
    - `run: 'df -PT -B1...` specifies the command to run.
    - `split: horizontal` states that each output line may return a metric.
    - `split_by` explains how to split each line in different fields. In this case, we use the `\s+` regular expression, which tells Flex that any sequence of one or more white spaces is a separator.
    - `row_start` specifies that data starts right after the first row (which is `0`).
    - `set_header` specifies, in order, a matching name for each value of the aforementioned array.
    - `perc_to_decimal: true` indicates to convert any percentage string into a decimal value, removing the trailing `%` symbol.

**Once the Flex config is created, the Infrastructure agent autodetects the new config and begins collecting data.**

To check that your new integration is working, execute the following [NRQL query](https://docs.newrelic.com/docs/query-data/nrql-new-relic-query-language):

```sql
FROM FileSystemSample SELECT mountedOn, fs, usedBytes, capacityBytes, usedBytes
```

The query should now produce a table similar to this:

![](./img/basic-filesystem.png)

##  4. <a name='Howtoaddmoreintegrations'></a>How to add more Flex integrations

Stand-alone Flex configurations, like most of our examples, start with the name of the integration and the [apis](/apis/readme.md). For example:

```yaml
name: linuxOpenFD
apis:
  - name: linuxOpenFD
    commands:
      - run: cat /proc/sys/fs/file-nr | awk '{print $1-$2,$3}'
        split: horizontal
        set_header: [openFD,maxFD]
        regex_match: true
        split_by: (\d+)\s+(.*)
```
These stand-alone configurations can be tested by invoking Flex from the command line; this is useful when [developing Flex integrations](../development.md), since invoking Flex directly doesn't send data to the New Relic platform:

```bash
sudo /var/db/newrelic-infra/newrelic-integrations/bin/nri-flex --verbose --pretty --config_file ./myconfig.yml
```

To use Flex configurations files with the Infrastructure agent, you need to add some lines at the beginning. For example, if we add the example above to our `integrations.d` file, we would get the following (notice the indentation):

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
              row_start: 1
              set_header: [fs,fsType,capacityBytes,usedBytes,availableBytes,usedPerc,mountedOn]
          perc_to_decimal: true
      name: linuxOpenFD
      apis:
        - name: linuxOpenFD
          commands:
            - run: cat /proc/sys/fs/file-nr | awk '{print $1-$2,$3}'
              split: horizontal
              set_header: [openFD,maxFD]
              regex_match: true
              split_by: (\d+)\s+(.*)
```
To insert multiple Flex configurations to the `integrations.d` config file, you can add multiple `nri-flex` blocks, each with an embedded Flex config:

```yaml
integrations:
 - name: nri-flex
   config:
     name: flexName_1
     # Flex config goes here
 - name: nri-flex
   config:
     name: flexName_2
     # Flex config goes here
 - name: nri-flex
   config:
     name: flexName_3
     # Flex config goes here
```
To minimize indentation issues, you can link to stand-alone Flex configuration files using the `config_template_path` directive:
```yaml
integrations:
  - name: nri-flex
    config_template_path: /path/to/flex/integration.yml
```

In the Flex repo you can find more than [200 config examples](../examples/flexConfigs) of custom integrations. Remember to add them under `config` in your integrations config file, or link to them using `config_template_path` statements.

>We strongly recommend that you use a YAML linter in your code editor to check for indentation issues in your config files. Most of the times, Flex rejects badly indented configurations.

##  5. <a name='Whatsnext'></a>What's next?

- Learn more about the Flex configuration schema in [Configure Flex](/basics/configure.md).
- Read about the [url](/apis/url.md) and [commands](/apis/command.md) APIs and how to create Flex integrations with them.
- See the [list of supported functions](/basics/functions.md) to understand what Flex is capable of.
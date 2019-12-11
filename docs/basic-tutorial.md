# FLEX step-by-step basic tutorial

Requirements:

* Infrastructure Agent version 1.8.x or higher
  - Flex can also work with older versions, but this tutorial relies on the
    latest integrations' configuration engine that has been added in the
    version 1.8.0.
  - [Please check the documentation to know more about its
    advantages](https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-configuration-file-specifications-agent-v180)       
* Run the Infrastructure Agent in Root/Administrator mode 
  - Current versions of Flex require administrator permissions for the
    management of temporary files.
  - It is the user that executes the Agent by default.  
* Flex 0.8.5 or higher.
  - To work with the new integrations system of the Infrastructure Agent 1.8.0
    and higher, it is recommended to use Flex 0.8.5.
  - Previous versions of Flex work correctly with the old integrations engine.

## Install

You can download Flex for Windows and Linux from the [Flex releases page](https://github.com/newrelic/nri-flex/releases).

The Flex package comes with an installer script (`install_linux.sh` or `install_windows.bat`).
It is aimed at installing all the required files for older versions of the Agent, and also
start/stopping the Agent service so the changes take effect.

From the Agent version 1.8.0, you can just manually copy the `nri-flex` executable from the
tarball into the `/var/db/newrelic-infra/newrelic-integrations/` folder:

Steps (from a command-line):
```
$ wget https://github.com/newrelic/nri-flex/releases/download/v0.8.5/nri-flex-linux-0.8.5.tar.gz
$ tar xzf nri-flex-linux-0.8.5.tar.gz
$ sudo cp nri-flex-linux-0.8.5/nri-flex /var/db/newrelic-infra/newrelic-integrations/
```

> Windows users: copy the `nri-flex.exe` file into `C:\Program Files\New Relic\newrelic-infra\newrelic-integrations`. 

Flex is already installed and ready to work with the agent.

## Checking that Flex is up and running

1. Create a file named `my-flex-configs.yml` (or any other name of your choose) into the
   `/etc/newrelic-infra/integrations.d` folder. 
    - Windows users: `C:\Program Files\New Relic\newrelic-infra\integrations.d` folder.
2. Set the following contents for the previously created file:
   ```yaml
   integrations:
     - name: nri-flex
       interval: 30s
       config:
         name: just-testing
   ```
3. Go to Insights and run the following query:

```sql
from flexStatusSample select * LIMIT 1
```

The query should show a table similar to the following:

![](./img/basic-table.png)

### What happened behind the scenes

1. The Infrastructure Agent detected that a new integration named `nri-flex` has been added.
2. The Agent looks for an executable named `nri-flex` in `/var/db/newrelic-infra/newrelic-integrations/`.
3. A temporary configuration file is created with the following contents:
   ```yaml
   name: just-testing
   ```
4. `nri-flex` is executed receiving the path of the above YAML file via the `CONFIG_PATH` environment
   variable.
5. Flex recognizes a configuration named `just-testing`, but since it does not provide extra information
   it just returns a `flexStatusSample` with some internal status of the Flex integration.

## Our first Flex integration

For this example, you will need a linux-based operating system, as it depends on Unix commands
that won't work in windows.

This example is about reporting disk metrics from file systems not natively supported by
New Relic using the `df` (Disk Free) command.

The objective of flex is to convert the text output of this command (disk free showing
file system, blocks and excluding the file systems already supported by the agent):

```
$ df -PT -B1 -x tmpfs -x xfs -x vxfs -x btrfs -x ext -x ext2 -x ext3 -x ext4
Filesystem     Type         1-blocks         Used    Available Capacity Mounted on
devtmpfs       devtmpfs    246296576            0    246296576       0% /dev
go_src         vboxsf   499963170816 361339486208 138623684608      73% /go/src
``` 

Into a set of equivalent JSON samples with this format:

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
    "mountedOn": "/go/src",
  }
}
```

> Note that the Agent will decorate each sample with extra fields

You need to tell Flex how to perform the above "table text to JSON" transformation,
concretely:

- Name of the metric (`FileSystem`). Flex will append the `Sample` suffix, resulting into
  `FileSystemSample`.
- Which command to run `df -PT -B1 ...`.
- How to split the output table from `df` and how to assign those values to given metric
  names.

The YAML configuration providing all the above information would be as following:
This is achieved placing the content below in the YAML configuration file from the previous
section:

```yaml
integrations:
  - name: nri-flex
    interval: 15s
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

Sections from the above YAML worth mentioning: 

- The `apis` section is an array of entries for each sample. Each entry sets a name for the
  sample, as well as the commands/procedures to get it and how to process it.
- First `apis` entry is named `FileSystem` (it will be used to construct the `FileSystemSample`
  event name).
- In the `commands` section, we specify how to get the information, basically:
    - `run: 'df -PT -B1...` specifies the command to run.
    - `split: horizontal` states that each output line may return a metric.
    - `split_by` explains how to split each line in different fields. In this case, we use
      a regular expression `\s+` telling that any sequence of 1 or more white spaces should
      be considered a separator. E.g. it would divide a line like:
      ```
      devtmpfs       devtmpfs    246296576            0    246296576       0% /dev
      ```
      Into an array containing `["devtmpfs", "devtmpfs", "246296576", "0", "246296576", "0%", "/dev"]`
    - `set_header` specifies a matching name for each value of the aforementioned array.
    - `perc_to_decimal: true` aims for converting any percentage string into a decimal value
      (this is, removing the trailing `%` symbol, if exists).

To check that our new integration is working, you can try executing the following query
in Insights:

```sql
FROM FileSystemSample SELECT mountedOn, fs, usedBytes, capacityBytes, usedBytes
```

![](./img/basic-filesystem.png)

## For more examples

You can check the [flex configs examples](../examples/flexConfigs) folder for more
working examples of Flex.
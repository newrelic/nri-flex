# `commands`

The `commands` API allows you to retrieve information from any application or shell command. It can accept multiple commands executed in sequence, raw text, and JSON data, and is platform-agnostic. Keep in mind that platform specific applications differ between different operating systems, and some may not be available (for example, Kubernetes).

- [Basic usage](#Basicusage)
- [Configuration properties](#Configurationproperties)
- [Advanced usage](#Advancedusage)

## <a name='Basicusage'></a>Basic usage

### Linux example

```yaml
name: example
apis:
  - name: linuxDirectorySize
    commands:
      - run: du -c /somedir
        split: horizontal
        set_header: [dirSizeBytes, dirName]
        regex_match: true
        split_by: (\d+)\s+(.*)
```

This Linux configuration retrieves the raw output provided by the command defined in `run`, which outputs a pair of values: the directory size, and the directory name. It also informs Flex that the output is horizontally formatted and has two columns as defined in `set_header`. Finally, it extracts the values using the regex expression defined in `split_by`, and assigns to each of the columns set in `set_header`.

### Windows example

A Windows configuration would follow the exact same structure:

```yaml
name: winSvc
apis:
  - name: winSvc
    commands:
      - run: wmic service get State,Name,DisplayName /format:csv
        split: horizontal
        regex_match: false
        split_by: \,
        row_header: 1
```

## <a name='Configurationproperties'></a>Configuration properties

The following table describes the properties of the `commands` API, which accepts a list of commands, each requiring a `run` directive.


| Name | Type | Default | Description|  
| ------ | ------ | ------ | ------ |
|                `run` |      string      |                                   | Command or application that you want to run. It accepts any valid shell command. You can also use environment variables with the format `$$ENV_VAR_NAME`                                                                                                                                                                       |
|              `shell` |      string      | `/bin/sh` (Linux) `cmd` (Windows) | Shell to use when executing the command defined in `run`. All native Linux shells, Windows CMD, and Windows PowerShell v1-5.x (`powershell`) and v6+ (`pwsh`) are supported.                                                                                                                                                   |
|              `split` |      string      |            `vertical`             | Mode of processing of the command output, either vertical with one value per line, or horizontal with more than one value per line (table format). Only used when `ignore_output` is false                                                                                                                                     |
|           `split_by` |      string      |                                   | Regular expression used to split metric data. It can accept a list of expressions when `sub_parse` is enabled                                                                                                                                                                                                                  |
|        `regex_match` |      string      |                                   | Whether the regular expression defined in `split_by` should be interpreted as a match expression (`true`) or as a split expression (`false`)                                                                                                                                                                                   |
|         `row_header` |       int        |                `0`                | Line that contains the header. Only applies if the value is not equal to `row_header` and is greater than or equal to `1`                                                                                                                                                                                                      |
|          `row_start` |       int        |                `0`                | Line number where Flex starts processing metric data. If `split` is set to `horizontal`, `row_start` is only used if `row_start` is not equal to `row_header` and `row_start` is greater than or equal to `1`                                                                                                                  |
|         `line_start` |       int        |                `0`                | Line number where Flex will start processing data. If `split` is set to `horizontal` and `row_start` is defined, `line_start` will only be used if `line_start` is not equal to `row_header` and `line_start` is greater than or equal to `1`. If both `row_start` and `line_start` are defined, `line_start` takes precedence |
|           `line_end` |       int        |                `0`                | Line number (exclusive) at which Flex stops processing data. Only applies if `split` is not equal to `horizontal`                                                                                                                                                                                                              |
|         `set_header` | array of strings |               `[]`                | Name and number of columns Flex should extract data from. Only applies if `split` is equal to `horizontal`                                                                                                                                                                                                                     |
| `header_regex_match` |       bool       |              `false`              | Whether the regular expression in `header_split_by` should be interpreted as a match expression (`true`) or as a split expression (`false`). Applies only if `split` is equal to `horizontal`                                                                                                                                  |
|    `header_split_by` |      string      |                                   | Regular expression applied to the header line. Applies only if `split` is equal to `horizontal`                                                                                                                                                                                                                                |
|       `split_output` |      string      |                                   | Regular expression used to split the output into blocks of data                                                                                                                                                                                                                                                                |
|            `timeout` |       int        |              `10000`              | Time to wait, in milliseconds, for the command to execute. If the command takes longer than `timeout`, Flex ignores the output and returns an error. Note that Flex waits for the command to stop by itself                                                                                                                    |     |
|             `assert` |       map        |                                   | [Check if command output matches or not matches your assertion string](#Assert-output-exists-before-processing)                                                                                                                                                                                                                |

## <a name='Advancedusage'></a>Advanced usage

The `commands` API accepts additional format directives to better parse the output.

### <a name='Rawdataparsing'></a>Raw data parsing

In the example below, the output from the `df` command is treated as a table with a header. We instruct Flex to extract values using a regex split expression; in this particular case, the regex expression tells Flex that the values are separated by spaces. Value are then assigned in order of appearance to the corresponding keys.

```yaml
---
name: example
apis:
  - name: diskFree
    commands:
      - run: df -T
        split: horizontal
        set_header:
          [fs, fsType, blocks, usedBytes, availableBytes, usedPerc, mountedOn]
        row_start: 1
        split_by: \s+
```

In this next example, we don't specify the keys to which the values will be assigned and instead use a regex expression to extract the keys from the header. Since we are not specifying the header row number neither the values start row, Flex assumes the first row is the header and the next lines are the value rows.

```yaml
---
name: example
apis:
  - name: diskFree
    commands:
      - run: df -T
        header_split_by: \s+
        split: horizontal
        regex_match: true
        split_by: (\S+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)%\s+(\d+)\s+(\d+)\s+(\d+)%\s+(.*)
```

To extract the values we use a regex expression in `split_by`. Note that in this case we extract the names of the metric attributes from raw data, so we must be sure that those are correct.

### <a name='Specifytheshell'></a>Specify the shell

All commands are executed using `/bin/sh` (Linux) or `cmd` (Windows). If you want to use a different shell, you can specify it at API level for all commands, or at command level, which overrides values set at the API level.

In the example below, `/bin/zsh` applies only to `some_other_command`, since `shell: /bin/bash` overrides the API level statement for `df -T`.

```yaml
---
name: example
apis:
  - name: diskFree
    shell: /bin/zsh
    commands:
      - run: df -T
        shell: /bin/bash
        header_split_by: \s+
        split: horizontal
        regex_match: true
        split_by: (\S+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)%\s+(\d+)\s+(\d+)\s+(\d+)%\s+(.*)
      - run: some_other_command
        split_by: \s+
```

As stated before, in Windows, you can also override the default shell, `cmd`, with, for example, PowerShell.

```yaml
# Used to get list of windows services
name: windowsServiceList
apis:
  - name: windowsServiceList
    commands:
      - run: Get-Service | Format-Table -Property Status, Name, DisplayName -Autosize
        shell: powershell
        event_type: WindowsCommand
        split: horizontal
        set_header: [status, name, displayname]
        row_start: 1
        regex_match: true
        split_by: (\w+)\s+(\w+)\s+(\w+)
```

In this example we are executing a command, `Get-Service`, using PowerShell as the command shell.

### <a name='Specifyatimeout'></a>Specify a timeout

Flex defines a 10 second timeout for each command by default. If the command does not complete within the timeout period, Flex stops processing the current command and moves to the next. You can change the timeout at both API and command levels. Timeout values are specified in milliseconds (for example, 15 seconds are specified as `15000`).

```yaml
name: example
apis:
  - name: linuxTimeout
    timeout: 5000
    shell: /bin/sh
    commands:
      - run: sleep 10
        name: timesOut
      - run: sleep 3
        name: doesNotTimeout
      - run: sleep 10
        name: alsoDoesNotTimeout
        timeout: 15000
```

In the example above we declare three commands. The timeout value of `5000` declared at the API level applies to all commands except `alsoDoesNotTimeout`, which has its own `timeout` statement. Note that Flex will return an error for the first command:

```
command: timed out err="context deadline exceeded" exec="sleep 10"`)
```

### <a name='Splittheoutput'></a>Split the output

In case the output of the command is not a sequential list of lines/values, use `split_output` to separate results into blocks that are to be processed sequentially. The directive accepts a regex expression for splitting the output into blocks. You can either use a list of regex expressions to extract data from each block, or try raw processing using `split_by`:

```yaml
name: example
apis:
  - name: splitOutput
    commands:
      - run: echo "key:value" && echo "---" && echo "other_key:otherValue"
        split_output: ---
        regex_matches:
          - expression: \S*key:(\w+)
            keys: [value]
```

If you run the command defined in the `run` directive, you get the following result:

```text
key:value
---
other_key:otherValue
```

The `split_output` command as defined above gets two blocks of data to which Flex then applies the `regex_matches` expressions for extracting the values. This example results in two metric samples:

```json
[
  {
    "event_type": "splitOutputSample",
    "integration_name": "com.newrelic.nri-flex",
    "value": "value"
  },
  {
    "event_type": "splitOutputSample",
    "integration_name": "com.newrelic.nri-flex",
    "value": "otherValue"
  }
]
```

### <a name='Manuallyspecifyblocksofdatatoprocess'></a>Manually specify blocks of data to process

If you know at which line the relevant data starts and where it ends, you can use `line_start` and `line_end` (optional) to limit the data processing to a specific number of lines from the output.

```yml
name: example
apis:
  - name: lineStart
    commands:
      - run: echo "this is noise" && echo "key:value"
        line_start: 1
        split_by: ":"
```

In the example above, we only want to process data after the first (`0`) line, so we set `line_start` to `1`. If we know that after a specific line the data is not useful, we can limit the set of lines Flex processes by adding `line_end`:

```yml
name: example
apis:
  - name: lineStart
    commands:
      - run: echo "this is noise" && echo "key:value" && echo "otherKey:otherValue" && echo "more noise"
        line_start: 1
        line_end: 3
        split_by: ":"
```

Note that `line_end` is exclusive, meaning that you have to add `1` (`0` indexed) to the actual line you want to stop being processed.

### Assert output exists before processing

Add an assert block with a nested variable of `match` and/or `not_match`.
Note both options use regex.

```yml
#### Command output must contain the string "hi" - this will returned in the payload
integrations:
  - name: nri-flex
    config:
      name: SomeIntegration
      apis:
        - event_type: SomeSample
          commands:
            - run: "echo hi:bye"
              split_by: ":"
              assert:
                match: hi #### <--------
```

```yml
#### Command output must contain the string "foo" - this will be discarded and not added to the payload
integrations:
  - name: nri-flex
    config:
      name: SomeIntegration
      apis:
        - event_type: SomeSample
          commands:
            - run: "echo hi:bye"
              split_by: ":"
              assert:
                match: foo #### <--------
```

```yml
#### Command output must contain the string "hi" and not contain the string "foo - this will be added to the payload
integrations:
  - name: nri-flex
    config:
      name: SomeIntegration
      apis:
        - event_type: SomeSample
          commands:
            - run: "echo hi:bye"
              split_by: ":"
              assert:
                match: hi ##### <------------
                not_match: foo #### <--------
```

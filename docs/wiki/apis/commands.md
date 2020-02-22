# `commands`

The `commands` API allows retrieving information from any application or shell command executed. It can accept multiple commands that are executed sequentially.

## Basic usage

```yaml
---
name: example
apis:
  - name: linuxDirectorySize
    commands:
      - run: du -c /somedir
        split: horizontal
        set_header: [dirSizeBytes,dirName]
        regex_match: true
        split_by: (\d+)\s+(.*)
```

The above Flex configuration retrieves the raw output provided by the command defined in the `run` directive, which outputs a pair of values, the directory size and the directory name. 
It also informs Flex that the output is horizontally formatted and has 2 columns as defined by the `set_header` directive.
Finally it extract the values via a regex match, using the regex expression defined in the `split_by` directive, and assignes them in order to each of the columns defined in the `set_header` directive.

## Configuration properties

The following table defines the properties of the `commands` API. The API accepts a list of **command**, each requiring the `run` directive, and accepting one of more directives described in the following table.

| Name | Type | Default | Description |
|---:|:---:|:---:|---|
| `run` | string | n.a. | Set it to the shell command (or application) that you want to run. It accepts any valid shell command. You can also use environment variables with the format $$ENV_VAR_NAME |
| `shell` | string | "/bin/sh" (linux)  "cmd"(windows) | Specifies the shell to use when executing the command defined in `run` |
| `split`| string | vertical | Defines the way the result of the command should be processed, either vertically with 1 value per line, or horizontally with more than 1 value per line (table like format). Only used when `ignore_output` is false (default)|
| `split_by` | string | | Defines the regular expression used to process "metric" data |
| `regex_match` | string | | Defines if the regular expression defined in `split_by` should be interpreted as match expression (true) or a split expression (false) |  
| `row_header` | int | 0 | Specifies which line contains the header. Only applies if != `row_header` and >= 1 otherwise ignored. |
| `row_start` | int | 0 | Specifies where the line number where Flex will start processing "metric" data. If `split` is set to `horizontal`, `row_start` will only be used if `row_start` != `row_header` and `row_start` >= 1 |
| `line_start` | int | 0 | Specifies the line number where Flex will start processing data. If `split` is set to `horizontal` and `row_start` is defined, `line_start` will only be used if `line_start` != `row_header` and `line_start` >= 1. If both `row_start` and `line_start` are defined, `line_start` has precedence |
| `line_end` | int | 0 | Specified the line number (exclusive) at which Flex will stop processing data. Only applies if `split` != `horizontal`|
| `set_header` | array of string | [] | Defines the name (and number) of columns Flex should extract data from. Only applies if `split` = `horizontal` |
| `header_regex_match` | bool | false | Defines whether the regular expression in `header_split_by` should be interpreted as a match expression (true) or a split expression (false). Applies only if `split` = `horizontal` |
| `header_split_by` | string | | Defines the regular expression that will be applied to the header line. Applies only if `split` = `horizontal` |
| `split_output` | string | | Defines the regular expression that is used to split the output into blocks of data. |
| `timeout` | int | 10 | Defines the time to wait for the command to execute. If the command takes longer Flex will ignore the output and return an error. Note that Flex will not forcefully stop the command, after the timeout, it will still wait for it to stop by itself| 

## Advanced usage

The `commands` api allows for more format directives to be defined to help Flex in parsing the output and different ways to achieve the same result.

### Raw data parsing

In the example below, the command being executed, `df`, outputs in a table-like format that includes a header defined by the directive `set_header`, and so the values start at `row_start`, which tells Flex the values start at row 1. 
We also inform Flex to extract the values via simple regex split expression. In this particular case, the expression tells Flex that the values are separated by spaces. Each value is assigned in order each key defined by the `set_header` directive.

```yaml
---
name: example
apis:
  - name: diskFree
    commands:
      - run: df -T
        split: horizontal
        set_header: [fs,fsType,blocks,usedBytes,availableBytes,usedPerc,mountedOn]
        row_start: 1
        split_by: \s+
```

In this next example, we do not specify the keys to which the values will be assigned and instead use a simple regex expression to extract the keys from the header. Since we are not specifying the header row number neither the values start row,
Flex assumes the first row is the header and the next lines are the values rows.
To extract the values we use a regex expression by setting the directive `regex_match` to true and declaring the expression in the `split_by` directive.
Note that in this case the names of the metric attributes will be extracted from the raw data so be sure that those are correct.

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

### Specifying the shell

By default, all commands are run with the shell `/bin/sh` (or `cmd` under windows). If for some reason you want to use a different shell, you case do so either by specifying it at the API level, which will apply equally to each command or you can specify it at the command level, which will apply to that command only, and as expected overrides any value set at the API level.

In the example below the shell `/bin/zsh` will apply only to the second command, even if it's declared at the API level, becuase the first command overwrites it, with `/bin/bash`.

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

## Specifying a timeout

By default, Flex defines a timeout of 10 seconds for each command. If the command does not complete within the timeout Flex will stop the processing the current command and move the next one.
You can, however, change timeout at the API level and at the command level. The timeout values are specified at the milisecond level, so for exmaple if you want to specify a value of 15 seconds you should use `15000`.

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

In the example above we declare 3 commands. The timeout value of `5000` declared at the API level will apply to all commands that do not override it locally, which in this configuration would the first (`timesOut`) and second (`doesNotTimeout`). 
Running Flex with this configuration will result in the first command returning an error (`command: timed out err="context deadline exceeded" exec="sleep 10"`), the second command running successfully and the third command also running successfully, since it overrides the timeout at the API level.

## Splitting the output

In case the output of the command is not a sequential list of lines/values you can use `split_output` to separate the results into blocks, that are then processed sequentially.

The directive accepts a regex expression that it uses to split the output into blocks.
Then it can either use a list of regex expressions to extract data from each block, or it can try and process it "raw" with a simple `split_by` directive.

```yml
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

If you run the command defined in the `run` directive you get the following result:

```text
key:value
---
other_key:otherValue
```

Using the `split_output` command as defined above, you will get 2 blocks of data that Flex will then apply the `regex_matches` expressions to extract the values.

The above example results in 2 metric samples similar to the example below.

```json
[{
  "event_type": "splitOutputSample",
  "integration_name": "com.newrelic.nri-flex",
  "value": "value"
},
{
  "event_type": "splitOutputSample",
  "integration_name": "com.newrelic.nri-flex",
  "value": "otherValue"
}]

```

### Manually specifying block of data to process

If you know at which line the relevant data starts and optionally where it end you can use the directives `line_start` and `line_end` to limit the data processing to a specific number of lines of output.

```yml
name: example
apis:
  - name: lineStart
    commands:
      - run: echo "this is noise" && echo "key:value"
        line_start: 1
        split_by: ":"
```

In the example above, we ony want to process data after the first (0) line, so we set `line_start` to 1. 
If we know that after some line, the data is not useful, we can limit the set of lines Flex will process, ass exemplified in the next example.

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

Note that `line_end` is exclusive, meaning that you have to use +1 (0 indexed) on the actual line you want to stop being processed.

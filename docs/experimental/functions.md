# Experimental functions

Experimental functions are available for use but currently not recommend for use in production environments (unless critical for your use case). They are not officially supported by New Relic.

- [Experimental functions](#experimental-functions)
  - [inherit_attributes](#inheritattributes)
  - [metric_parser](#metricparser)
  - [pagination](#pagination)
  - [pluck_numbers](#plucknumbers)
  - [rename_samples](#renamesamples)
  - [sample_include_match_all_filter](#sampleincludematchallfilter)
  - [sample_keys](#samplekeys)
  - [save_output](#saveoutput)
  - [store_variables](#storevariables)
  - [sub_parse](#subparse)

## inherit_attributes

When you use `start_key` to start processing a nested payload, you may want to inherit the attributes above it as well. That's what `inherit_attributes`is for. **Note**:  only supported when `start_key` is used.

Consider a service that returns the following payload:

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc":{
            "def":123,
            "hij":234
        }
    },
    "name": "node3"
}
```

You could include the top-level attributes in the sample:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      start_key: 
        - leaderInfo>abc
```

Which would give you a result similar to:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 1,
  "leaderInfo.abc.hij": 2,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node3"
}]
```

## metric_parser

`metric_parser` enables setting rates and deltas. It expects an exact match to avoid any conflicts, though the `mode` attribute can be set as shown below to automatically match multiple keys. `mode` supports `regex`, `suffix`, `prefix` and `contains` for matching.

Flex automatically attempts to set a namespace as required for your attributes; else you can namespace based on existing attributes or a custom attributes.

See the `metric_parser` example below:

```yaml
name: redisFlex
apis:
  - name: redis
    commands:
      - run: (printf "info\r\n"; sleep 1) | nc 127.0.0.1 6379 # or even ### run: "redis-cli -h 127.0.0.1 -p 6379 info" ### (depends on redis-cli)
        split_by: ":"
    remove_keys: # remove any keys that contain any of the following strings
      - human
    snake_to_camel: true
    perc_to_decimal: true
    sub_parse:
      - type: prefix
        key: db
        split_by:
          - ","
          - "="
    custom_attributes:
      myCustomAttr: theValue
    metric_parser:
      metrics:
        totalNetInputBytes: RATE
        rate$: RATE
      namespace: # you can create a namespace with a custom attribute, or chain together existing attributes, else it will default
        # custom_attr: "mySpecialRedisServer"
        existing_attr:
          - redisVersion
          - tcpPort
      # mode: regex ### switches metric parser to use a defined mode rather then exact match, options include "regex" ,"suffix", "prefix" & "contains"

```

## pagination

See the inline comments on how to use pagination.

```yaml
---
name: paginationTest
apis:
  - event_type: paginationTest
    url: https://reqres.in/api/users?page=${page}&per_page=2
    # url: https://reqres.in/api/users?page=${page}&per_page=${limit}
    pagination:
      page_start: 1 ### select the page to start from
      # increment: 10 ### number to increment by // default 1
      # page_limit: 2 ### can be used as a page offset place ${limit} into the url
      # page_limit_key: per_page ### select a key in the payload to set the page limit / offset
      # page_next_key: next_page ### select a key in the payload to set the next page to walk too
      # max_pages: 3 ### set max number of pages to walk
      # max_pages_key: total_pages ### select a key in the payload to set the max pages to walk
      # next_cursor_key: nextCursor ### if using cursor pagination look for this key instead, will get substituted into ${page}
      ############################# you will need to also set a ?flex=${page} query parameter for tracking eg. https://reqres.in/api/users?flex=${page}
      # next_link_key: nextLink ### look for specified key to navigate to next
      ############################# you will need to also set a ?flex=${page} query parameter for tracking eg. https://reqres.in/api/users?flex=${page}
      payload_key: data ### select a key in the payload to check if there is still content being returned
```
## pluck_numbers

Retrieves any attribute with a number value and assigns it to another attribute. Any value that contains numbers is automatically plucked out. If no number is found, the value is left as is.

Consider the following payload:

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc": {
            "def": 123,
            "hij": 234
        }
    },
    "name": "node3"
}
```

You could retrieve the values from `leaderInfo.abc.def` and `leaderInfo.abc.hij`:

```yaml
name: squidFlex
apis:
  - name: squidMgrUtilization
    commands:
      - run: squidclient -v mgr:utilization
        split_by: " = "
        line_limit: 88 # stop processing at this line as we only care about last 5 minute metrics
    pluck_numbers: true # find any numbers within a string and pluck them out
    value_parser:
      time: "[0-9]+" # use regex to find any time values, and pluck the first found integer out with the value_parser
```

## rename_samples

Uses a regular expression to rename a sample (`event_type` attribute) if the current sample has a key that matches. In the example below, if the `db` key is found, it's renamed `redisDbSample`; if `cmd` is found, rename to `redisCmdSample`.

```yaml
---
name: redis
apis:
  - name: redis
    url: http://127.0.0.1:8887/metrics
    prometheus:
      enable: true
    rename_samples:
      db: redisDbSample
      cmd: redisCmdSample
```
## sample_include_match_all_filter

Similar to the supported smaple_include_filter but will create samples only when ALL the specified filter keys and values are present in the sample.

Consider a service that returns the following payload:

```json
{
    "usageInfo": [
        {
            "incidentCode": 77,
            "serviceId": "compute"
        },
        {
            "incidentCode": 143,
            "serviceId": "compute"
        }
    ]
}
```

You may only want to have `"serviceId": "compute"` in the output sample when `"incidentCode": 77`:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/samples
      sample_include_match_all_filter:
          - serviceId: "compute"
          - incidentCode: 77
```

Which would return the following:

```json
"metrics": [
  {
    "api.StatusCode": 200,
    "incidentCode": 77,
    "serviceId": "compute",
    "event_type": "usageInfoSample",
    "integration_name": "com.newrelic.nri-flex",
    "integration_version": "Unknown-SNAPSHOT",
    "quantities": 10
  },
]
```

## sample_keys

Creates different samples based on a key from a larger object. There can be cases where you can receive a payload where you have subobjects identified by a key (like a map), and you want to extract them as a different sample. You can target a nested key and split them out into samples.

Consider the following payload:

```json 
{
    "followers": {
        "6e3bd23ae5f1eae0": {
            "counts": {
                "fail": 0,
                "success": 745
            },
            "latency": {
                "average": 0.017039507382550306,
                "current": 0.000138,
                "maximum": 1.007649,
                "minimum": 0,
                "standardDeviation": 0.05289178277920594
            }
        },
        "a8266ecf031671f3": {
            "counts": {
                "fail": 0,
                "success": 735
            },
            "latency": {
                "average": 0.012124141496598642,
                "current": 0.000559,
                "maximum": 0.791547,
                "minimum": 0,
                "standardDeviation": 0.04187900156583733
            }
        }
    },
    "leader": "924e2e83e93f2560"
}
```

You could create samples based on each of the `followers`:

```yaml
name: example
apis:
  - name: startKey
    url: http//some-service.com/status
    sample_keys:
      # create samples distinguished by the follower id
      followerSample: followers>follower.id
```


## save_output

| Applies to  | Description |
| :---------- | :---------- |
| API | Saves sample output to a .JSON file specified by the user, any directories in the path must exist prior. |

**Example**

Consider a service that returns the following payload:

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc": {
            "def": 123,
            "hij": 234
        }
    },
    "name": "node3"
}
```
You could save the output in a file called *results.json* in the flexConfigs folder:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      save_output: "./flexConfigs/results.json"
```
Files are saved with 0644 UNIX file permissions.


## store_variables

Stores variables from any API result that can be accessed anywhere in any subsequent API.

Consider the following payload:

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc":{
            "def":123,
            "hij":234
        }
    },
    "name": "node3"
}
```

You could store the value of key `id` to be used in the next API:

```yaml
name: example
apis:
  - name: storeVariables
    url: http://some-service.com/status
    store_variables:
      nodeId: id
  - name: useVariables
    url: http://some-service.com/${var:nodeId}/status
```

```yaml
---
name: dummyFlex
apis:
  - name: todo
    url: https://jsonplaceholder.typicode.com/todos/2
    store_variables:
      storedId: userId ### store the userId from this response into storedId
  - name: user
    url: https://jsonplaceholder.typicode.com/users/${var:storedId}  ### query the user route with the previously stored userId which is storedId
```

## sub_parse

Splits nested values out from one line. For example, `db0:keys=1,expires=0,avg_ttl=0` to `db0.keys = 1, db0.expires = 0, db0.avg_ttl = 0`.

```yaml
apis:
  - name: redis
    commands:
      - run: (printf "info\r\n"; sleep 1) | nc -q0 127.0.0.1 6379
        split_by: ":"
    sub_parse:
      - type: prefix
        key: db
        split_by:
          - ","
          - "="
```

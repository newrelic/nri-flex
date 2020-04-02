# Data parsing and transformation functions

Flex has many useful functions, which can be combined in different ways to help you manipulate and tidy up your data.

- [Data parsing and transformation functions](#data-parsing-and-transformation-functions)
  - [Function precedence order](#function-precedence-order)
  - [Flex supported functions](#flex-supported-functions)
    - [add_attribute](#addattribute)
    - [convert_space](#convertspace)
    - [ignore_output](#ignoreoutput)
    - [keep_keys](#keepkeys)
    - [lazy_flatten](#lazyflatten)
    - [lookup_file](#lookupfile)
    - [math](#math)
    - [perc_to_decimal](#perctodecimal)
    - [remove_keys](#removekeys)
    - [rename_keys / replace_keys](#renamekeys--replacekeys)
    - [sample_filter](#samplefilter)
    - [save_output](#saveoutput)
    - [sample_include_filter](#sampleincludefilter)
    - [sample_exclude_filter](#sampleexcludefilter)
    - [snake_to_camel](#snaketocamel)
    - [split_array](#splitarray)
    - [split_objects](#splitobjects)
    - [start_key](#startkey)
    - [store_lookups](#storelookups)
    - [strip_keys](#stripkeys)
    - [timestamp](#timestamp)
    - [to_lower](#tolower)
    - [value_parser](#valueparser)
    - [value_transformer](#valuetransformer)

## Function precedence order

Flex applies data parsing and transformation functions in a specific order, regardless of where in the configuration files you declare them. Keep the functions precedence order in mind to avoid unexpected or empty results.

- [Data parsing and transformation functions](#data-parsing-and-transformation-functions)
  - [Function precedence order](#function-precedence-order)
  - [Flex supported functions](#flex-supported-functions)
    - [add_attribute](#addattribute)
    - [convert_space](#convertspace)
    - [ignore_output](#ignoreoutput)
    - [keep_keys](#keepkeys)
    - [lazy_flatten](#lazyflatten)
    - [lookup_file](#lookupfile)
    - [math](#math)
    - [perc_to_decimal](#perctodecimal)
    - [remove_keys](#removekeys)
    - [rename_keys / replace_keys](#renamekeys--replacekeys)
    - [sample_filter](#samplefilter)
    - [save_output](#saveoutput)
    - [sample_include_filter](#sampleincludefilter)
    - [sample_exclude_filter](#sampleexcludefilter)
    - [snake_to_camel](#snaketocamel)
    - [split_array](#splitarray)
    - [split_objects](#splitobjects)
    - [start_key](#startkey)
    - [store_lookups](#storelookups)
    - [strip_keys](#stripkeys)
    - [timestamp](#timestamp)
    - [to_lower](#tolower)
    - [value_parser](#valueparser)
    - [value_transformer](#valuetransformer)

\*Happens before attribute modification and autoflattening, which is useful to get rid of unwanted data and arrays early on.

## Flex supported functions

Here is a list of supported functions. Be aware that while all the examples use JSON payloads for convenience, source data can be in a variety of different formats.

### add_attribute

| Applies to | Description                                                                                                      |
| :--------- | :--------------------------------------------------------------------------------------------------------------- |
| API        | Adds extra attributes to the resulting sample. Can use attributes from the result to create the extra attribute. |

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

You could use the payload to generate a link that can be added as an extra attribute to the resulting sample:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      add_attribute:
          # use the 'id' attribute of the service output
          link: https://some-other-service/nodes/${id}
```

Which would return the following:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 1,
  "leaderInfo.abc.hij": 2,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "link": "https://some-other-service/nodes/eca0338f4ea31566",
  "name": "node3"
}]
```

### convert_space

| Applies to | Description                                         |
| :--------- | :-------------------------------------------------- |
| API        | Replaces spaces in key names with other characters. |

**Example**

Consider a service that returns the following payload:

```json
{
    "id": "eca0338f4ea31566",
    "leader info": {
        "leader": "8a69d5f6b7814500",
        "start time": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc": {
            "def": 123,
            "hij": 234
        }
    },
    "name": "node3"
}
```

You could convert the spaces in `leader info` and `start time` to, for example, underscores:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      convert_space: '_'
```

Which would return the following:

```json
"metrics": [{
  "id": "eca0338f4ea31566",
  "leader_info.abc.def": 123,
  "leader_info.leader": "8a69d5f6b7814500",
  "leader_info.abc.hij": 234,
  "leader_info.start_time": "2014-10-24T13:15:51.186620747-07:00",
  "leader_info.uptime": "10m59.322358947s",
  "name": "node3"
}]
```

### ignore_output

| Applies to | Description                                                                                                                                                                     |
| :--------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| API        | Ignores the output of some API, that is, it does not create a sample for the result, but still caches it. This is useful when creating lookups/cache for other APIs executions. |

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

You could use it as source of values for other APIs:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      store_lookups:
          # store the 'id' into a lookup key named 'nodeId'
          nodeId: id
      ignore_output: true
    - name: useLookup
      # use the 'nodeId' stored in the previous APIs to execute this one
      url: http://some-other-service.com/${lookup:nodeId}/status
```

### keep_keys

| Applies to | Description                                                                                             |
| :--------- | :------------------------------------------------------------------------------------------------------ |
| API        | Keeps only the keys matching the regular expressions. This is useful for keeping just some key metrics. |

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

You could keep just the `id` and `name` fields by using this configuration:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      keep_keys:
          - id
          - name
```

### lazy_flatten

| Applies to | Description                                                                                                                                                                                                                                         |
| :--------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| API        | Performs a lazy flattening operation. The result differs depending on the object that's flattened. By default, Flex always performs data flattening; depending on the type of payload it either creates one sample or many, all with the same name. |

**Example**

Consider a service that returns the following json payload:

```json
{
    "contacts": [
        {
            "name": "batman",
            "number": 911
        },
        {
            "name": "robin",
            "number": 112
        }
    ]
}
```

Flex flattens the structure and creates two samples if not asked to perform any transformation.

For example, using the following configuration:

```yaml
name: example
apis:
    - name: status
      url: http://some-service.com/status
```

Will give you a result similar to:

```json

"metrics": [{
    "event_type": "statusSample",
    "name": "batman",
    "number": 911
  },
  {
    "event_type": "statusSample",
    "name": "robin",
    "number": 112
  }
]
```

If you want to have all the data in the same sample, you could perform a `lazy_flatten`:

```yaml
name: example
apis:
    - name: status
      url: http//some-service.com/status
    lazy_flatten:
      - contacts
```

Which would return something similar to the following:

```json
"metrics": [{
    "contacts.flat.0.name": "batman",
    "contacts.flat.0.number": 911,
    "contacts.flat.1.name": "robin",
    "contacts.flat.1.number": 112,
  }
]
```

On the other hand, with a payload like the following:

```json
{
    "contacts": {
        "first": {
            "name": "batman",
            "number": 911
        },
        "second": {
            "name": "robin",
            "number": 112
        }
    }
}
```

The same configuration gives the following Which would return the following:

```json
"metrics": [{
    "contacts.flat.first.name": "batman",
    "contacts.flat.first.number": 911,
    "contacts.flat.second.name": "robin",
    "contacts.flat.second.number": 112,
  }
]
```

### lookup_file

| Applies to | Description                                                                                     |
| :--------- | :---------------------------------------------------------------------------------------------- |
| API        | Dynamically inject values into configurations using a JSON file containing an array of objects. |

**Example**

Uses a lookup file to dynamically generate separate configuration files for each object within the array, and substitutes the variables in the configuration using the expression `${lf:var-name}`.

Consider a file with the following content:

```json
[
    {
        "name": "some-service",
        "addr": "some-service.com:80"
    },
    {
        "name": "another-service",
        "addr": "another-service.com:80"
    }
]
```

Assuming each service returns a payload similar to:

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

You could use `lookup_file` to generate multiple API executions and therefore samples:

```yaml
name: example
lookup_file: addresses.json
apis:
    - name: someService
      url: http://${lf:addr}/status
```

Which would return the following:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node3",
  },{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node3",
}
```

### math

| Applies to | Description                                                                                                           |
| :--------- | :-------------------------------------------------------------------------------------------------------------------- |
| API        | Performs math operations with the values of the attributes specified in the expression and/or other explicit numbers. |

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

You could create another attribute that is, for example, the **sum** of attributes `leaderInfo.abc.def` and `leaderInfo.abc.hij`:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      math:
          sum: ${leaderInfo.abc.def} + ${leaderInfo.abc.hij} + 1
```

Which would return the following:

```json
"metrics": [{
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "a2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": 10,
  "name": "node3",
  "sum": 358
}]
```

### perc_to_decimal

| Applies to | Description                                                              |
| :--------- | :----------------------------------------------------------------------- |
| API        | Converts any percentage formatted value into its decimal representation. |

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
            "def": "123%",
            "hij": "234%"
        }
    },
    "name": "node3"
}
```

You could convert the percentage formatted values in `leaderInfo.abc.def` and `leaderInfo.abc.hij` to their decimal representations.

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      perc_to_decimal: true
```

Which would return the following:

```json
"metrics": [{
  "id": "eca0338f4ea31566",
  "leader_info.abc.def": 123,
  "leader_info.abc.hij": 234,
  "leader_info.leader": "8a69d5f6b7814500",
  "leader_info.start_time": "2014-10-24T13:15:51.186620747-07:00",
  "leader_info.uptime": "10m59.322358947s",
  "name": "node3"
}]
```

### remove_keys

| Applies to | Description                                                                    |
| :--------- | :----------------------------------------------------------------------------- |
| API        | Uses a regular expression to remove selected keys (attributes) from your data: |

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

You could remove some of the keys using `remove_keys`:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      remove_keys:
          - time
```

Which would return something similar to the following:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "8a69d5f6b7814500",
  "name": "node3"
}]
```

Be aware that the value of `remove_keys` matches at any level, meaning that it could remove complete objects if any part of the name matches the regular expression.

### rename_keys / replace_keys

| Applies to | Description                          |
| :--------- | :----------------------------------- |
| API        | Uses a regex to find and rename keys |

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

You could rename the key `id` to `identifier`, and `name` to `nodeName`:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      # replace_keys for backcompat
      rename_keys:
          id: identifier
          name: nodeName
```

Which would return the following:

```json
"metrics": [{
  "identifier": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "8a69d5f6b7814500",
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "nodeName": "node3"
}]
```

### sample_filter

| Applies to | Description                                                              |
| :--------- | :----------------------------------------------------------------------- |
| API        | Skips creating the sample if both a key and value is found in the sample |
|            | if `sample_exclude_filter` is present, both filters will be applied.     |

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

### save_output

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
      sample_filter:
          - name: node3
```

Which would return the following:

```json
"metrics": []
```
Files are saved with 0644 UNIX file permissions.

### sample_include_filter

| Applies to | Description                                                                                                      |
| :--------- | :--------------------------------------------------------------------------------------------------------------- |
| API        | If a sample is included using sample_include_filter, Flex evaluates sample_filter and sample_exclude_filter next |

**Example**

Consider a service that returns the following payload:

```json
{
    "usageInfo": [
        {
            "quantities": 10,
            "customerId": "abc"
        },
        {
            "quantities": 20,
            "customerId": "xyz"
        }
    ]
}
```

You may only want to have `"customerId": "abc"` in the ouput sample:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/samples
      sample_include_filter:
          - customerId: abc
```

Which would return the following:

```json
"metrics": [
  {
    "api.StatusCode": 200,
    "customerId": "abc",
    "event_type": "usageInfoSample",
    "integration_name": "com.newrelic.nri-flex",
    "integration_version": "Unknown-SNAPSHOT",
    "quantities": 10
  },
]
```

### sample_exclude_filter

| Applies to | Description                                                                                  |
| :--------- | :------------------------------------------------------------------------------------------- |
| API        | Skips creating the sample if both a key and value is found in the sample.                    |
|            | If `sample_filter` is present, both `sample_filter` and `sample_exclude_filter` are applied. |

**Example**

Consider a service that returns the following payload:

```json
{
    "usageInfo": [
        {
            "quantities": 10,
            "customerId": "abc"
        },
        {
            "quantities": 20,
            "customerId": "xyz"
        }
    ]
}
```

You may want to exclude `"customerId": "abc"` from the output sample:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/samples
      sample_exclude_filter:
          - customerId: abc
```

Which would return the following:

```json
"metrics": [
  {
    "api.StatusCode": 200,
    "customerId": "xyz",
    "event_type": "usageInfoSample",
    "integration_name": "com.newrelic.nri-flex",
    "integration_version": "Unknown-SNAPSHOT",
    "quantities": 20
  },
]
```

### snake_to_camel

| Applies to | Description                                                          |
| :--------- | :------------------------------------------------------------------- |
| API        | Converts all snake-cased attributes into camelCased formatted names. |

**Example**

Consider a service that returns the following payload:

```json
{
    "id": "eca0338f4ea31566",
    "leader_info": {
        "leader": "8a69d5f6b7814500",
        "start_time": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc": {
            "def": 123,
            "hij": 234
        }
    },
    "name": "node3"
}
```

You could convert `leader_info` and `start_time` to camelCase for increased consistency:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      snake_to_camel: true
```

Which would return the following:

```json
"metrics": [{
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.leader": "8a69d5f6b7814500",
  "leaderInfo.abc.hij": 234,
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node3"
}]
```

### split_array

| Applies to | Description                           |
| :--------- | :------------------------------------ |
| API        | Split an array that has nested arrays |

**Example**

Consider a service that returns the following payload:

```json
{
    "status": 1,
    "appstatus": -128,
    "statusstring": null,
    "appstatusstring": null,
    "results": [
        {
            "status": -128,
            "schema": [
                {
                    "name": "TIMESTAMP",
                    "type": 6
                },
                {
                    "name": "HOST_ID",
                    "type": 5
                },
                {
                    "name": "HOSTNAME",
                    "type": 9
                },
                {
                    "name": "PERCENT_USED",
                    "type": 6
                }
            ],
            "data": [
                [1582159853733, 0, "7605f6bec898", 0],
                [1582159853733, 2, "067ea6fc4c22", 0],
                [1582159853733, 1, "62a10d3f45e3", 0]
            ]
        }
    ]
}
```

You could split the configuration:

```yaml
name: example
apis:
    - name: voltdb_cpu
      event_type: voltdb
      url: http://some-service.com/status
      split_array: true
      set_header: [TIMESTAMP, HOST_ID, HOSTNAME, PERCENT_USED]
      start_key:
          - results>data
```

Which would return the something like following:

```json
"metrics": [{
   "HOSTNAME": "7605f6bec898",
   "HOST_ID": 0,
   "PERCENT_USED": 4,
   "TIMESTAMP": 1582161013979,
   "event_type": "voltdb",
   },{
   "HOSTNAME": "067ea6fc4c22",
   "HOST_ID": 2,
   "PERCENT_USED": 4,
   "TIMESTAMP": 1582161013978,
   "event_type": "voltdb",
   },{
   "HOSTNAME": "62a10d3f45e3",
   "HOST_ID": 1,
   "PERCENT_USED": 4,
   "TIMESTAMP": 1582161013980,
   "event_type": "voltdb",
 }]
```

### split_objects

| Applies to | Description                                             |
| :--------- | :------------------------------------------------------ |
| API        | Splits an object that has nested objects into an array. |

**Example**

Consider a service that return the following payload:

```json
{
    "first": {
        "id": "eca0338f4ea31566",
        "leaderInfo": {
            "uptime": "10m59.322358947s",
            "abc": {
                "def": 123,
                "hij": 234
            }
        },
        "name": "node1"
    },
    "second": {
        "id": "eca0338f4ea31566",
        "leaderInfo": {
            "uptime": "10m59.322358947s",
            "abc": {
                "def": 123,
                "hij": 234
            }
        },
        "name": "node2"
    }
}
```

You could split the single object into two separate objects:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      split_objects: true
```

Which would return something similar to the following:

```json
"metrics": [{
  "event_type": "Sample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node1",
  "split.id": "first"
},
{
  "event_type": "Sample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node2",
  "split.id": "second"
}]
```

### start_key

| Applies to | Description                                                  |
| :--------- | :----------------------------------------------------------- |
| API        | Starts processing data at a different point in your payload. |

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

You could tell Flex to start processing the payload from `leaderInfo`:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      start_key:
          - leaderInfo
```

This would mean processing only the following data:

```json
{
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc": {
            "def": 123,
            "hij": 234
        }
    }
}
```

Which would return something similar to

```json
"metrics": [{
  "abc.def": 123,
  "abc.hij": 234,
  "event_type": "someServiceSample",
  "leader": "8a69d5f6b7814500",
  "startTime": "2014-10-24T13:15:51.186620747-07:00",
  "uptime": "10m59.322358947s"
}]
```

Or further down:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      start_key:
          - leaderInfo
          - abc
```

Which would mean processing only this data:

```json
{
    "abc": {
        "def": 123,
        "hij": 234
    }
}
```

Which would return something similar to

```json
"metrics": [{
  "def": 123,
  "event_type": "someServiceSample",
  "hij": 234
}
```

### store_lookups

| Applies to | Description                                                          |
| :--------- | :------------------------------------------------------------------- |
| API        | Stores attributes from a API that you could use in a subsequent API. |

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

You could store the `id` attribute to be used in a subsequent API.

```yaml
name: example
apis:
    - name: storeLookups
      url: http://some-service.com/status
      store_lookups:
          # store the 'id' into a lookup key named 'nodeId'
          nodeId: id
    - name: useLookup
      url: http://some-other-service.com/${lookup:nodeId}/status
```

### strip_keys

| Applies to | Description                                     |
| :--------- | :---------------------------------------------- |
| API        | Removes entire keys or objects from the output. |

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

You could completely remove the `leaderInfo` object:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      strip_keys:
          - leaderInfo
```

This would return something similar to:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "name": "node3"
}]
```

You could also remove nested keys, for example `leader` and `startTime` under the `leaderInfo` object:

```yaml
name: example
apis:
    - name: stripKeys
      url: http://some-service.com/status
      strip_keys:
          - leaderInfo>leader
          - leaderInfo>startTime
```

Which would return something similar to:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node3"
}]
```

Note that Flex strips all keys that match the payload. This means that if the payload has multiple objects that match the `strip_keys` value, all are be removed.

### timestamp

| Applies to | Description                                                                                     |
| :--------- | :---------------------------------------------------------------------------------------------- |
| Anywhere   | Injects timestamps anywhere in your config and also performs additions or subtractions on them. |

You can use the following expressions to inject a timestamp formatted in various ways:

```
${timestamp:[ms|ns|s|date|datetime|datetimetz|dateutc|datetimeutc|datetimeutctz][+|-][Number][ms|milli|millisecond|ns|nano|nanosecond|s|sec|second|m|min|minute|h|hr|hour]}
```

-   `ms` - milliseconds
-   `s` - seconds
-   `ns` - nanoseconds
-   `date` - current date
-   `datetime` - current datetime
-   `datetimetz` - current datetime with timezone
-   `dateutc` - current utc date
-   `datetimeutc` - current utc datetime
-   `datetimeutctz` - current utc datetime with timezone

For example:

```
${timestamp:ms} - current timestamp in milliseconds
${timestamp:date} - date in local timezone: 2006-01-02
${timestamp:datetime} - date and time in local timezone : 2006-01-02T03:04
${timestamp:datetimetz} - date and time in local timezone, with timezone : 2006-01-02T15:04:05Z07:00
${timestamp:dateutc} - date in utc timezone: 2006-01-02
${timestamp:datetimeutc} - date and time in  utc timezone: 2006-01-02T03:04
${timestamp:datetimeutctz} - date and time in utc timezone, with timezone: 2006-01-02T15:04:05Z07:00
```

To perform calculations, you can use any of the following expressions (or similar):

```
${timestamp:ms-5000} subtract 5000 from current timestamp in milliseconds
${timestamp:ms+10000}" add 10000 to current timestamp in milliseconds

${timestamp:datetime-1hr} subtract 1 hour from current datetime, return datetime
${timestamp:datetime+60min} add 60 minutes to current datetime, return datetime
```

### to_lower

| Applies to | Description                     |
| :--------- | :------------------------------ |
| API        | Converts all keys to lowercase. |

**Example**

Consider a service that returns the following payload:

```json
{
    "Id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc": {
            "def": 123,
            "hij": 234
        }
    },
    "Name": "node3"
}
```

You could rename all keys to lowercase:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      to_lower: true
```

The result would be similar to the following (notice all keys are lowercase, including keys that would be camelCased):

```json
"metrics": [{
  "id": "eca0338f4ea31566",
  "leaderinfo.abc.def": 123,
  "leaderinfo.leader": "8a69d5f6b7814500",
  "leaderinfo.abc.hij": 234,
  "leaderinfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderinfo.uptime": "10m59.322358947s",
  "name": "node3"
}]
```

### value_parser

| Applies to | Description                                                                                                   |
| :--------- | :------------------------------------------------------------------------------------------------------------ |
| API        | Finds keys using a regular expression and applies another regular expresion to extract the first value found. |

**Example**

Consider a service that returns the following payload:

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "a8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc": {
            "def1": "a:123",
            "def2": "a:234"
        }
    },
    "name": "node3"
}
```

You could use `value_parser` to extract/transform the numbers on keys `leaderInfo.abc.def1` and `leaderInfo.abc.def2`, and replace them in the result with the transformed values:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      value_parser:
          def: '[0-9]+'
```

Which would return the following:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "a2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": 10,
  "name": "node3",
}]
```

### value_transformer

| Applies to | Description                                                      |
| :--------- | :--------------------------------------------------------------- |
| API        | Uses a regular expression to find a key and transforms its value |

**Example**

Consider a service that returns the following payload:

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "uptime": "10m59.322358947s",
        "abc": {
            "def": 123,
            "hij": 234
        }
    },
    "name": "node3"
}
```

Without declaring any other transformation you would get a result similar to:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderTime.abc.def": 123,
  "leaderTime.abc.hij": 234,
  "leaderTime.leader": "8a69d5f6b7814500",
  "leaderTime.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderTime.uptime": "10m59.322358947s",
  "name": "node3"
}
```

If you want to transform the value of key `name` into a format like for example **<node/name>**, you could use `value_transformer` like this:

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      value_transformer:
          name: node_${value}
```

Which would return something similar to:

```json
"metrics": [{
  "event_type": "someServiceSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "8a69d5f6b7814500",
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node/node3"
}]
```

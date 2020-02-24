# Data parsing and transformation functions

-   [Data parsing and transformation functions](#data-parsing-and-transformation-functions)
    -   [Order of application](#order-of-application)
    -   [Supported functions](#supported-functions)
        -   [snake_to_camel](#snaketocamel)
        -   [to_lower](#tolower)
        -   [convert_space](#convertspace)
        -   [perc_to_decimal](#perctodecimal)
        -   [start_key](#startkey)
        -   [strip_keys](#stripkeys)
        -   [rename_keys / replace_keys](#renamekeys--replacekeys)
        -   [remove_keys](#removekeys)
        -   [keep_keys](#keepkeys)
        -   [store_lookups](#storelookups)
        -   [cache](#cache)
        -   [lazy_flatten](#lazyflatten)
        -   [value_transformer](#valuetransformer)
        -   [split_objects](#splitobjects)
        -   [split_array](#split_array)
        -   [value_parser](#valueparser)
        -   [math](#math)
        -   [lookup_file](#lookupfile)
        -   [ignore_output](#ignoreoutput)
        -   [add_attribute](#addattribute)
        -   [custom_attributes](#customattributes)
        -   [sample_filter](#samplefilter)
        -   [timestamp](#timestamp)

Flex has many useful functions to help you manipulate and tidy up your data output in many way, by combining these functions in different ways.

## Order of application

While you can declare these functions in the configuration file(s) in mostly any order, Flex internally applies them in a specific order. You have to be aware of this order so that you don't get unexpected results or is osme cases no result at all.

[**Note**: cleanup this]

-   Find Start Key
-   Strip Keys - Happens before attribute modifiction and auto flattening, useful to get rid of unneeded data and arrays early
-   Lazy Flatten
-   Standard Flatten
-   Remove Keys
-   Strip Keys (second round)
-   Merge (if used)
-   ToLower Case
-   Convert Space
-   snake_case to camelCase
-   Value Parser
-   Pluck numbers
-   Sub parse
-   Value Transformer
-   Rename Key // uses regex to find keys to replace
-   Store Lookups
-   Keep Keys // keeps only keys you want to keep, and removes the rest

## Supported functions

The following functions are currently officially supported.
These function are applied by Flex after the result of the API has been extracted, so while we show the results using JSON representation dues to it's ease of visualization, it is not exactly how it works internally.

### snake_to_camel

| Valid at | Description                                                                         |
| :------- | :---------------------------------------------------------------------------------- |
| API      | Converts all attributes with "snake" formatted names into camelCase formatted names |

Given a service that returns the following payload

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

you can convert `leader_info` and `start_time` to camelCase in order to make the result more consistent

```yaml
name: example
apis:
    - name: startKey
      url: http://some-service.com/status
      snake_to_camel: true
```

which should give a similar result to

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

### to_lower

| Valid at | Description                    |
| :------- | :----------------------------- |
| API      | Converts all keys to lowercase |

Given a service that returns the following payload

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

you can rename all key to lowercase

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      to_lower: true
```

which would give you a result similar to (notice all keys are lowercase, including keys that would be camelCased)

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

### convert_space

| Valid at | Description                                             |
| :------- | :------------------------------------------------------ |
| API      | Converts spaces in keys names into some other character |

Given a service that returns the following payload

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

you can convert the spaces in `leader info` and `start time` to, for example, underscores

```yaml
name: example
apis:
    - name: startKey
      url: http://some-service.com/status
      convert_space: '_'
```

which should give you a similar result to

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

### perc_to_decimal

| Valid at | Description                                                             |
| :------- | :---------------------------------------------------------------------- |
| API      | Converts any percentage formatted value into its decimal representation |

Given a service that returns the following payload

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

you can convert the percentage formatted values in `leaderInfo.abc.def` and `leaderInfo.abc.hij` to their decimal representations.

```yaml
name: example
apis:
    - name: startKey
      url: http://some-service.com/status
      perc_to_decimal: true
```

which should give you a similar result to

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

### start_key

| Valid at | Description                                                 |
| :------- | :---------------------------------------------------------- |
| API      | Starts processing data at a different point in your payload |

Given a service that returns the following payload

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

you can tell Flex to start processing the payload from `leaderInfo`

```yaml
name: example
apis:
    - name: startKey
      url: http://some-service.com/status
      start_key:
          - leaderInfo
```

which would mean processing only the following data

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

of further down,

```yaml
name: example
apis:
    - name: startKey
      url: http://some-service.com/status
      start_key:
          - leaderInfo
          - abc
```

which would mean processing only the following data

```json
{
    "abc": {
        "def": 123,
        "hij": 234
    }
}
```

### strip_keys

| Valid at | Description                                     |
| :------- | :---------------------------------------------- |
| API      | Remove entire keys or "objects" from the output |

Given the following payload:

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

you can completely remove the `abc` object,

```yaml
name: example
apis:
    - name: stripKeys
      url: http://some-service.com/status
      strip_keys:
          - abc
```

giving you the following resulting payload,

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s"
    },
    "name": "node3"
}
```

or remove nested key(s), for example `leader` and `startTime` under the `leaderInfo` object,

```yaml
name: example
apis:
    - name: stripKeys
      url: http://some-service.com/status
      strip_keys:
          - leaderInfo>leader
          - leaderInfo>startTime
```

giving you the following resulting payload,

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

Note that Flex will strip all keys that match from the payload. This means that if the payload has multiple objects that match the `strip_keys` value, all of them will be removed.

### rename_keys / replace_keys

| Valid at | Description                          |
| :------- | :----------------------------------- |
| API      | Uses a regex to find and rename keys |

Given a service that returns the following payload

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

you can rename the key `id` to `identifier` and `name` to `nodeName`

```yaml
name: example
apis:
    - name: startKey
      url: http://some-service.com/status
      # replace_keys for backcompat
      rename_keys:
          - id: identifier
          - node: nodeName
```

giving you the following result

```json
{
    "identifier": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": "2014-10-24T13:15:51.186620747-07:00",
        "uptime": "10m59.322358947s",
        "abc": {
            "def": 123,
            "hij": 234
        }
    },
    "nodeName": "node3"
}
```

### remove_keys

| Valid at | Description                                                                      |
| :------- | :------------------------------------------------------------------------------- |
| API      | Uses regular expression to select keys (attributes) to be removed from your data |

Given a service that returns the following payload

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

you can remove some of the keys

```yaml
name: example
apis:
    - name: removeKeys
      url: http://some-service.com/status
      remove_keys:
          - time
```

which would result in the following

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "abc": {
            "def": 123,
            "hij": 234
        }
    },
    "name": "node3"
}
```

Be aware that `remove_keys` will match at any level, meanig that it can remove complete "objects" if any part of the name matches the regular expression.

### keep_keys

| Valid at | Description                                                                                            |
| :------- | :----------------------------------------------------------------------------------------------------- |
| API      | Keeps only the keys matching the regular expressions. This is useful for keeping just some key metrics |

Given a service that returns the following payload

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

you can keep just the `id` and `name` by using the following configuration

```yaml
name: example
apis:
    - name: removeKeys
      url: http://some-service.com/status
      keep_keys:
          - id
          - name
```

### store_lookups

| Valid at | Description                                                       |
| :------- | :---------------------------------------------------------------- |
| API      | Stores attributes from a API that you can use in a subsequent API |

Given a service that returns the following payload

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

you can store the `id` attribute to be used in a subsequent API.

```yaml
name: example
apis:
    - name: storeLookups
      url: http://some-service.com/status
      store_lookup:
          # store the 'id' into a lookup key named 'nodeId'
          - nodeId: id
    - name: useLookup
      url: http://some-other-service.com/${lookup:nodeId}/status
```

### cache

| Valid at     | Description                                                                                                                                                                                                                                                                                              |
| :----------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| API, Command | Allows you to reuse the output from a previous API as the result of current API. Flex keeps the results of each API (`url` or `commands`) in a "cache" that can be used as the input of a subsequent API. For `url` APIs the cache **key** is the URL and for the `commands` APIs is the name of the API |

Given a service that returns the following payload

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

resulting from executing the following `url` API

```yaml
name: example
apis:
    - name: stripKeys
      url: http://some-service.com/status
```

we may want to process it in another API (for example because we want to use a different sample for the result).
Note here that the cache `key` is the URL because it's an `url` API.

```yaml
name: example
apis:
    - name: status
      url: http://some-service.com/status
    - name: otherStatus
      cache: http://some-service.com/status
      strip_keys:
          - id
          - name
```

With a `commands` API you should use the name of the API,

```yaml
name: example
apis:
    - name: status
      commands:
          # assume that this file contains the same json payload showed above the beginning
          - run: cat /var/some/file
    - name: otherStatus
      cache: status
      strip_keys:
          - id
          - name
```

### lazy_flatten

| Valid at           | Description                                                                                                                                                                                                                                                       |
| :----------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| API                | Performs a lazy flatten. Depending on the type of "object" it is flattening , array of objects or map, the result will differ. Flex by default always performs data flattening and depending on the type of payload it will either create one (1) sample or many, |
| with the same name |

Given a service that returns the following json payload,

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

Flex will flatten the structure and create 2 samples, if you don't ask it to perform any transformation.
For example, using the following configuration

```yaml
name: example
apis:
    - name: status
      url: http//some-service.com/status
```

will give you a result similar to

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

If you want to have all the data in the same sample you can perform a `lazy_flatten`

```yaml
name: example
apis:
    - name: status
      url: http//some-service.com/status
    lazy_flatten:
      - contacts
```

which will give a result similar to

```json
"metrics": [{
    "contacts.flat.0.name": "batman",
    "contacts.flat.0.number": 911,
    "contacts.flat.1.name": "robin",
    "contacts.flat.1.number": 112,
  }
]
```

On the other hand, with a payload like the following

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

the same configuration will give a result like

```json
"metrics": [{
    "contacts.flat.first.name": "batman",
    "contacts.flat.first.number": 911,
    "contacts.flat.second.name": "robin",
    "contacts.flat.second.number": 112,
  }
]
```

### value_transformer

| Valid at | Description                                                      |
| :------- | :--------------------------------------------------------------- |
| API      | Uses a regular expression to find a key and transforms its value |

Given a service that returns the following payload

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

without declaring any other transformation you will get a result similar to

```json
"metrics": [{
  "id": "eca0338f4ea31566",
  "leaderTime.abc.def": 123,
  "leaderTime.abc.hij": 234,
  "leaderTime.leader": "8a69d5f6b7814500",
  "leaderTime.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderTime.uptime": "10m59.322358947s",
  "name": "node3"
}
```

If you want to transform the value of key `name` into a format like for example **<node/name>** you can use `value_transformer` like this

```yaml
name: example
apis:
    - name: removeKeys
      url: http://some-service.com/status
      value_transformer:
          name: node_${value}
```

### split_objects

| Valid at | Description                                            |
| :------- | :----------------------------------------------------- |
| API      | Splits an object that has nested objects into an array |

Given a service that return the following payload

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

You can split the single "object" into 2 different ones

```yaml
name: example
apis:
    - name: removeKeys
      url: http://some-service.com/status
      split_objects: true
```

which will give you a similar result to the following

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

### value_parser

| Valid at | Description                                                                                                |
| :------- | :--------------------------------------------------------------------------------------------------------- |
| API      | Finds keys using a regular expression and another regular expresion again to extract the first value found |

Given a service that returns the following payload

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

you can use `value_parser` to extract/transform the numbers on keys `leaderInfo.abc.def1` and `leaderInfo.abc.def2` and replace them in the result with the values transformed

```yaml
name: example
apis:
    - name: removeKeys
      url: http://some-service.com/status
      value_parser:
          def: '[0-9]+'
```

which will give you a result similar to

```json
"metrics": [{
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "a2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": 10,
  "name": "node3",
}]
```

### math

| Valid at | Description                                                                                                          |
| :------- | :------------------------------------------------------------------------------------------------------------------- |
| API      | Performs math operations with the values of the attributes specified in the expression and/or other explicit numbers |

Given a service that returns the following payload

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

you can create another attribute that is the, for example **sum**, of attributes `leaderInfo.abc.def` and `leaderInfo.abc.hij`

```yaml
name: example
apis:
    - name: removeKeys
      url: http://some-service.com/status
      math:
          sum: ${leaderInfo.abc.def} + ${leaderInfo.abc.hij} + 1
```

which will give you a reuslt similar to

```json
"metrics": [{
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 100,
  "leaderInfo.abc.hij": 100,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "a2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": 10,
  "name": "node3",
  "sum": 201
}]
```

### lookup_file

| Valid at | Description                                                                                       |
| :------- | :------------------------------------------------------------------------------------------------ |
| API      | Uses a json file containing an array of objects, to dynamically inject values into configurations |

Using a lookup file will generate a separate config file dynamically for each object within the array, and substitute the variables in the configuration using the expression **\${lf:var-name}**

Given a file with the following contents

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

and assuming each service returns a payload similar to:

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

you can use to generate multiple API executions and therefore samples

```yaml
name: example
lookup_file: addresses.json
apis:
    - name: lookupFile
      url: http://${lf.addr}/status
      math:
          sum: ${leaderInfo.abc.def} + ${leaderInfo.abc.hij}
```

which would give you a result similar to

```json
"metrics": [{
  "event_type": "lookupFileSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node3",
  "sum": 357
  },{
  "event_type": "lookupFileSample",
  "id": "eca0338f4ea31566",
  "leaderInfo.abc.def": 123,
  "leaderInfo.abc.hij": 234,
  "leaderInfo.leader": "a8a69d5f6b7814500",
  "leaderInfo.startTime": "2014-10-24T13:15:51.186620747-07:00",
  "leaderInfo.uptime": "10m59.322358947s",
  "name": "node3",
  "sum": 357
}
```

### ignore_output

| Valid at | Description                                                                                                                                                                 |
| :------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| API      | Ignore the output of some API, ie, does not create a sample for the result, but still caches the result. It is useful when creating lookups/cache for other APIs executions |

Given an service that returns the following payload

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

we can use it as source of values for other APIs

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      store_lookup:
          # store the 'id' into a lookup key named 'nodeId'
          - nodeId: id
      ignore_output: true
    - name: useLookup
      # use the 'nodeId' stored in the previous APIs to execute this one
      url: http://some-other-service.com/${lookup:nodeId}/status
```

### add_attribute

| Valid at | Description                                                                                                     |
| :------- | :-------------------------------------------------------------------------------------------------------------- |
| API      | Adds extra attributes to the resulting sample. Can use attributes from the result to create the extra attribute |

Given a service that returns the following payload

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

we can use it to generate a link that can be added as an extra attribute to the resulting sample

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      add_attribute:
          # use the 'id' attribute of the service output
          link: https://some-other-service/nodes/${id}
```

which would give you a result similat to

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
}
```

### custom_attributes

| Valid at    | Description                                   |
| :---------- | :-------------------------------------------- |
| Config, API | Adds extra attributes to the resulting sample |

Given a service that returns the following payload

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

you can add extra attributes at the resulting sample

```yaml
name: example
custom_attributes:
    global_attr: global_value
apis:
    - name: someService
      url: http://some-service.com/status
      custom_attributes:
          api_attr: api_value
```

which would give you a result similar to

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
  "global_attr": "global_value",
  "api_attr": "api_value",
}
```

### sample_filter

| Valid at | Description                                                              |
| :------- | :----------------------------------------------------------------------- |
| API      | Skips creating the sample if both a key and value is found in the sample |

Given a service that returns the following payload

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

you can completeley skip creating the output sample

```yaml
name: example
apis:
    - name: someService
      url: http://some-service.com/status
      sample_filter:
          name: node3
```

which would give you a result similar to

```json
"metrics": []
```

### timestamp

| Valid at | Description                                                                                 |
| :------- | :------------------------------------------------------------------------------------------ |
| Anywhere | Injects timestamps anywhere in your config and also perform addition or subtraction on them |

Injects timestamps anywhere in your config and also performs addition or subtraction on them.
You can use the following expressions to inject a timestamp formatted in various ways:

```raw
${timestamp:[ms|ns|s|date|datetime|datetimetz|dateutc|datetimeutc|datetimeutctz][+|-][Number][ms|milli|millisecond|ns|nano|nanosecond|s|sec|second|m|min|minute|h|hr|hour]}
```

-   "ms" - milliseconds
-   "s" - seconds
-   "ns" - nanoseconds
-   "date" - current date
-   "datetime" - current datetime
-   "datetimetz" - current datetime with timezone
-   "dateutc" - current utc date
-   "datetimeutc" - current utc datetime
-   "datetimeutctz" - current utc datetime with timezone

For example:

```raw
${timestamp:ms} - current timestamp in milliseconds
${timestamp:date} - date in local timezone: 2006-01-02
${timestamp:datetime} - date and time in local timezone : 2006-01-02T03:04
${timestamp:datetimetz} - date and time in local timezone, with timezone : 2006-01-02T15:04:05Z07:00
${timestamp:dateutc} - date in utc timezone: 2006-01-02
${timestamp:datetimeutc} - date and time in  utc timezone: 2006-01-02T03:04
${timestamp:datetimeutctz} - date and time in utc timezone, with timezone: 2006-01-02T15:04:05Z07:00
```

To perform calculations you can use any of the following expressions (or similar):

```raw
${timestamp:ms-5000} subtract 5000 from current timestamp in milliseconds
${timestamp:ms+10000}" add 10000 to current timestamp in milliseconds

${timestamp:datetime-1hr} subtract 1 hour from current datetime, return datetime
${timestamp:datetime+60min} add 60 minutes to current datetime, return datetime
```

### split_array

| Valid at | Description                           |
| :------- | :------------------------------------ |
| API      | Split an array that has nested arrays |

Split an array that has nested arrays.
eg. You receive a payload that looks like below

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

The following config can split these neatly for you.

```yaml
apis:
    - name: voltdb_cpu
      event_type: voltdb
      # url: <voltdb json api - CPU URL > e.g. http://127.0.0.1:32952/api/1.0/?Procedure=@Statistics&Parameters=["CPU"]
      url: http://127.0.0.1:32952/api/1.0/?Procedure=@Statistics&Parameters=["CPU"]
      split_array: true
      set_header: [TIMESTAMP, HOST_ID, HOSTNAME, PERCENT_USED]
      start_key:
          - results>data
```

Output:

```json
{
    "name": "com.newrelic.nri-flex",
    "protocol_version": "3",
    "integration_version": "Unknown-SNAPSHOT",
    "data": [
        {
            "metrics": [
                {
                    "HOSTNAME": "7605f6bec898",
                    "HOST_ID": 0,
                    "PERCENT_USED": 4,
                    "TIMESTAMP": 1582161013979,
                    "event_type": "voltdb",
                    "integration_name": "com.newrelic.nri-flex",
                    "integration_version": "Unknown-SNAPSHOT"
                },
                {
                    "HOSTNAME": "067ea6fc4c22",
                    "HOST_ID": 2,
                    "PERCENT_USED": 4,
                    "TIMESTAMP": 1582161013978,
                    "event_type": "voltdb",
                    "integration_name": "com.newrelic.nri-flex",
                    "integration_version": "Unknown-SNAPSHOT"
                },
                {
                    "HOSTNAME": "62a10d3f45e3",
                    "HOST_ID": 1,
                    "PERCENT_USED": 4,
                    "TIMESTAMP": 1582161013980,
                    "event_type": "voltdb",
                    "integration_name": "com.newrelic.nri-flex",
                    "integration_version": "Unknown-SNAPSHOT"
                },
                {
                    "event_type": "flexStatusSample",
                    "flex.Hostname": "C02W60KWHTD8",
                    "flex.IntegrationVersion": "Unknown-SNAPSHOT",
                    "flex.counter.ConfigsProcessed": 1,
                    "flex.counter.EventCount": 3,
                    "flex.counter.EventDropCount": 0,
                    "flex.counter.HttpRequests": 1,
                    "flex.counter.voltdb": 3,
                    "flex.time.elaspedMs": 15,
                    "flex.time.endMs": 1582161013969,
                    "flex.time.startMs": 1582161013954
                }
            ],
            "inventory": {},
            "events": []
        }
    ]
}
```

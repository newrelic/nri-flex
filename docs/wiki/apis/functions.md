Flex has many useful functions to help you manipulate and tidy up your data output in many ways, and can also be used in various combinations.

The below functions are available under each defined api:
* [start_key](#start_key) Start processing data at a different point in your payload
* [cache](#cache) Reuse output from something previously run eg. a command or url
* [store_lookups](#store_lookups) Store some values to be used as a lookup in a subsequent request
* [store_variables](#store_variables) Store variables to reuse in any subsequent url, cmd, jmx query etc.
* [strip_keys](#strip_keys) Strip keys from a payload
* [lazy_flatten](#lazy_flatten) Lazy flatten a payload
* [sample_keys](#sample_keys) Create samples out of nested objects
* [replace_keys](#replace_keys) Replace entire keys or parts of a key using regex
* [rename_samples](#rename_samples) Using regex to find a key, rename the sample if found
* [rename_keys](#rename_keys) Rename keys using a contains match
* [remove_keys](#remove_keys) Remove keys entirely
* [keep_keys](#keep_keys) Remove all keys and keep particular keys
* [to_lower](#to_lower) Converts the key to all lower case
* [snake_to_camel](#snake_to_camel) Converts snake_case to camelCase, eg. super_hero -> superHero
* [perc_to_decimal](#perc_to_decimal) Converts percentages to decimals
* [sub_parse](#sub_parse) Sub parse a nested value
* [metric_parser](#metric_parser) For setting metrics to be RATEs or DELTAs
* [value_parser](#value_parser) Use regex to find a key, and use regex again to pluck a value
* [value_transformer](#value_transformer) Use regex to find a key, and then transform the value
* [pluck_numbers](#pluck_numbers) Pluck all numbers out, if not found leave value as is
* [math](#math) Perform math operations with or without the current data set
* [timestamp](#timestamp) Apply timestamps anywhere in configs, and add or subtract to timestamps as well
* [split_objects](#split_objects) Split an object that has nested objects
* [lookup_file](#lookup_file) Supply a json file containing an array of objects, to dynamically substitute into config(s)
* [pagination](#pagination) Handle Pagination for HTTP URLs
* [inherit_attributes](#inherit_attributes) Inherits attributes from parent of payload
* [add_attribute](#add_attribute) Construct new attributes based on data from current sample
* [ignore_output](#ignore_output) Completely ignore the output of a sample

```yaml
---
name: myFlexIntegration
apis: 
  - event_type: alertSample
    url: https://someapi.com/stats
    | <- here!
    start_key: # <- example 
```

#### start_key
Can be used to process data from a different point in your payload.

eg. 
```yaml
you could start processing the below payload from leaderInfo
start_key:
 - leaderInfo

or

start processing the payload from a further nested point
start_key:
 - leaderInfo
 - abc

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

#### cache

Allows you to reuse the output from a previously run command or http response by the name or url used.

```yaml
name: elasticsearchFlex
global:
    base_url: http://localhost:9200/
    # user: elastic
    # pass: elastic
    headers:
      accept: application/json
apis: 
  ### here we are calling the _stats endpoint, but splitting the stats sample to several samples to avoid overloading a single sample
  ### since Flex has inbuilt caching we can retrieve URL data previously fetched and split the sample out
  - event_type: elasticsearchTotalSample
    url: _stats
    start_key:
      - _all
      - total
  - event_type: elasticsearchPrimarySample
    cache: _stats ### we can use previously cached calls saved in memory
    start_key:
      - _all
      - total
```

With commands
```yaml
---
name: commandCache
apis: 
  - name: getHost
    commands: 
      - run: echo "zHost:$(hostname)"
        split_by: ":"
  - name: abc
    cache: getHost
```

From URL output, to commands

```yaml
### see the full example in the repo: examples/flexConfigs/nginx-opensource-stub-example.yml

name: nginxFlex
apis:
  - name: nginxStub
    url: http://127.0.0.1/nginx_status
  - name: nginx
    merge: NginxSample
    commands:
      - cache: http://127.0.0.1/nginx_status
        split_by: ": "
        line_end: 1
```

#### store_lookups

You can have a situation, where a previous API call made, may have attributes you would like to use in a subsequent call elsewhere. Flex allows you to store these attributes, and use them in a subsequent call or command.

```yaml
- event_type: rabbitVHostSample
    url: api/vhosts
    snake_to_camel: true
    store_lookups: ### we store all vhost "name" attributes found into the lookupStore to be used later
      vhosts: name 
  - event_type: rabbitAliveTestSample         ### as we use the special lookup key, we fetch the vhosts found, and execute the next request(s)
    url: api/aliveness-test/${lookup:vhosts}  ### if multiple vhosts were stored, this will issue multiple requests
```

#### strip_keys

Strip keys can remove entire keys and objects from processing.

```yaml
Remove the entire "incidents" object

apis: 
  - event_type: incidentSample
    url: incidents
    strip_keys:
      - incidents

Remove nested key(s) under the "incidents" object

We remove the "transitions", and "pagedPolicies" objects/keys under the "incidents" object.
apis: 
  - event_type: incidentSample
    url: incidents
    strip_keys:
      - incidents>transitions
      - incidents>pagedPolicies
```

#### lazy_flatten

Performs a lazy flatten. Uses no smarts to flatten a payload, eg. you could have an array of contact numbers, or some sort of dimensional data.

```
eg.

receive:
{ contacts:[ {name:batman,number:911}, {name:robin,number:000} ] }
after lazy_flatten:
contacts.0.name = batman
contacts.0.number = 911
contacts.1.name = robin
contacts.1.number = 000

apis: 
  - event_type: incidentSample
    url: incidents
    lazy_flatten:
      - numbers
      - incidents>pagedUsers # we can also do it another level nested within
      - incidents>pagedTeams
```

#### sample_keys

You can receive payloads where they are not in an array, but in a larger object split by key.
We can easily target a nested key and split them out into samples.

eg.
``` 
receive:
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

config:
name: etcdFlex
global:
    base_url: http://127.0.0.1:2379/v2/
    headers:
      accept: application/json
apis: 
  - event_type: etcdLeaderSample
    url: stats/leader
    sample_keys:
      followerSample: followers>follower.id # there are objects within objects distinguished by key, we can create samples like so
```

#### replace_keys

Uses regex to find keys to replace. 

```yaml
# replace the "os" key in the payload, to "operatingSystem"
apis: 
  - name: redis
    commands: 
      - run: (printf "info\r\n"; sleep 1) | nc -q0 127.0.0.1 6379 # or even ### run: "redis-cli -h 127.0.0.1 -p 6379 info" ### (depends on redis-cli)
        split_by: ":"
    replace_keys:
      os: operatingSystem # replaces os > operatingSystem
```

#### rename_samples

Using regex to find a key, if found rename to have a different sample name. As below eg. if "db" key is found, rename to redisDbSample, if "cmd" is found, rename to redisCmdSample.

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

#### rename_keys

Renames any part of a key using regex to match, with supplied text.
eg. you could rename "superhero" to "hero" by
 - renaming, super to nothing "" , or 
 - rename superhero to "hero"

```yaml
example:
  - event_type: elasticsearchNodeSample
    url: _nodes/stats
    sample_keys:
      nodes: nodes>node.id # there are objects within objects distinguished by key, we can create samples like so
    rename_keys: # <- here
      _nodes: parentNodes # rename nodes to parentNodes
```

#### remove_keys

Uses regex to remove any key(s) (attributes) from your data.

```yaml
apis: 
  - name: redis
    commands: 
      - run: (printf "info\r\n"; sleep 1) | nc -q0 127.0.0.1 6379 # or even ### run: "redis-cli -h 127.0.0.1 -p 6379 info" ### (depends on redis-cli)
        split_by: ":"
    remove_keys: # remove any keys that contain any of the following strings
      - human # <- here
```

#### keep_keys

Delete all other keys, and keep only the ones you define using regex.

```yaml
apis: 
  - name: tomcatThreads
    event_type: tomcatThreadSample
    ### note "keep_keys" will do the inverse, if you want all metrics remove the keep keys blocks completely
    ### otherwise tailor specific keys you would like to keep, this uses regex for filtering
    ### this is useful for keeping key metrics
    keep_keys: ###
      - bean
      - maxThreads
      - connectionCount
```

#### sub_parse
Splits nested values out from one line 
eg. db0:keys=1,expires=0,avg_ttl=0
to db0.keys = 1, db0.expires = 0, db0.avg_ttl = 0
```yaml
apis: 
  - name: redis
    commands: 
      - run: (printf "info\r\n"; sleep 1) | nc -q0 127.0.0.1 6379 # remove -q0 if testing on mac
        split_by: ":"
    sub_parse:
      - type: prefix
        key: db
        split_by:
          - ","
          - "="
```

#### metric_parser

Setting rates and deltas can be performed by the metric_parser. 

By default it expects an exact match to avoid any conflicts, however the "mode" attribute can be set as shown below to automatically match many keys. The "mode" option supports "regex", "suffix", "prefix" & "contains" for matching.

Flex will automatically attempt to set a namespace as required for your attributes, else you can namespace based on existing attributes or a custom attributes.

eg. see metric_parser options further below
set either RATE or DELTA

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

#### value_parser
Find keys using regex, and again use regex to pluck the first found value out

eg. find "time" and pluck the first found integer out.
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

#### value_transformer
eg. Transform the value of a key, called "key" where the value is "world" to "hello-world"

key: world -> key: hello-world

```yaml
---
name: dummy
apis: 
  - event_type: myEventSample
    url: http://127.0.0.1:8887/test
    value_transformer:
      key: hello-${value} 
## value is already world, so we can substitute it back in, which would now equal "key": "hello-world"
```

#### pluck_numbers
Any values that contain numbers, are automatically plucked out. If no number found, the value is left as is.

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

#### store_variables
Store variables from any execution point, and reuse in any subsequent url, command, query etc.

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
#### math

Perform math operations with or without current data set.
[See full example here](https://github.com/newrelic/nri-flex/blob/master/examples/flexConfigs/ping-url.yml)

```yaml
name: nginxFlex
apis:
  - name: nginxStub
    url: http://127.0.0.1/nginx_status
    rename_keys:
      "Active connections": "net.connectionsActive"
    pluck_numbers: true
    math: ### <- here, the key will become the new metric key, the attributes wrapped in ${attribute} can be from any existing attributes, or you can use your own numbers too
      net.connectionsDroppedPerSecond: ${net.connectionsAcceptedPerSecond} - ${net.handledPerSecond}
..........
```

#### timestamp

Apply timestamps anywhere in your config, and subtract or add to them.
eg. apply the below anywhere in your config.
```
"${timestamp:ms}" current timestamp in milliseconds
"${timestamp:ms-5000}" subtract 5000 from current timestamp in milliseconds
"${timestamp:ms+10000}" add 10000 to current timestamp in milliseconds

"${timestamp:date}" - date in date format local timezone: 2006-01-02
"${timestamp:datetime}" - datetime in date and time format local timezone : 2006-01-02T03:04
"${timestamp:datetimetz}" - datetime in date and time format local timezone : 2006-01-02T15:04:05Z07:00
"${timestamp:dateutc}" - date in date format utc timezone: 2006-01-02
"${timestamp:datetimeutc}" - datetime in date and time format utc timezone: 2006-01-02T03:04
"${timestamp:datetimeutctz}" - datetime in date and time format utc timezone: 2006-01-02T15:04:05Z07:00

"${timestamp:datetime-1hr}" subtract 1 hour from current datetime, return datetime
"${timestamp:datetime+60min}" add 60 minutes to current datetime, return datetime

* ${timestamp:[ms|ns|s|date|datetime|datetimetz|dateutc|datetimeutc|datetimeutctz][+|-][Number][ms|milli|millisecond|ns|nano|nanosecond|s|sec|second|m|min|minute|h|hr|hour]}

```
Supports:
* "ms" - milliseconds
* "s" - seconds
* "ns" - nanoseconds
* "date" - current date
* "datetime" - current datetime
* "datetimetz" - current datetime with timezone
* "dateutc" - current utc date
* "datetimeutc" - current utc datetime
* "datetimeutctz" - current utc datetime with timezone


#### split_objects

Split an object that has nested objects.
eg. You receive a payload that looks like below
```json
{
    "hg.nginx.org": {
        "processing": 0,
        "requests": 204,
        "responses": {
            "1xx": 0,
            "2xx": 191,
            "3xx": 12,
            "4xx": 1,
            "5xx": 0,
            "total": 204
        },
        "discarded": 0,
        "received": 45310,
        "sent": 2913986
    },
    "trac.nginx.org": {
        "processing": 0,
        "requests": 278,
        "responses": {
            "1xx": 0,
            "2xx": 185,
            "3xx": 84,
            "4xx": 2,
            "5xx": 6,
            "total": 277
        },
        "discarded": 1,
        "received": 65422,
        "sent": 2825682
    }
}
```
The following config can split these neatly for you.

```yaml
name: splitObjectExample
apis:
  - event_type: NginxEndpointSample
    url: http://demo.nginx.com/api/3/http/server_zones
    split_objects: true
```

Output:

```json
{
  "name": "com.newrelic.nri-flex",
  "protocol_version": "2",
  "integration_version": "0.6.0-pre",
  "data": [
    {
      "metrics": [
        {
          "discarded": 0,
          "event_type": "NginxEndpointSample",
          "integration_name": "com.newrelic.nri-flex",
          "integration_version": "0.6.0-pre",
          "processing": 0,
          "received": 54808,
          "requests": 250,
          "responses.1xx": 0,
          "responses.2xx": 236,
          "responses.3xx": 13,
          "responses.4xx": 1,
          "responses.5xx": 0,
          "responses.total": 250,
          "sent": 3357038,
          "split.id": "hg.nginx.org"
        },
        {
          "discarded": 1,
          "event_type": "NginxEndpointSample",
          "integration_name": "com.newrelic.nri-flex",
          "integration_version": "0.6.0-pre",
          "processing": 0,
          "received": 71475,
          "requests": 324,
          "responses.1xx": 0,
          "responses.2xx": 213,
          "responses.3xx": 99,
          "responses.4xx": 2,
          "responses.5xx": 9,
          "responses.total": 323,
          "sent": 3292360,
          "split.id": "trac.nginx.org"
        },
        {
          "event_type": "flexStatusSample",
          "flex.ConfigsProcessed": 1,
          "flex.EventCount": 2,
          "flex.EventDropCount": 0,
          "flex.NginxEndpointSample": 2
        }
      ],
      "inventory": {},
      "events": []
    }
  ]
}
```

#### lookup_file

Supply a json file containing an array of objects, to dynamically substitute into config(s)
This will generate a separate config file dynamically for each object within the array, and substitute the variables in with the below helper substitutions.

eg. ${lf:addr} expects there to be an "addr" attribute in the object

config example:
```yaml
name: portTestWithLookup
lookup_file: testLookup.json ### location of lookup file
apis: 
  - name: portTest
    timeout: 1000 ### default 1000 ms increase if you'd like
    commands:
      - dial: ${lf:addr}
    custom_attributes: ### add some additional attributes if you'd like
      name: ${lf:name}
      abc: ${lf:abc}
```

testLookup.json could be for eg.
```json
 [
     {
         "name":"google.com",
         "addr":"google.com:80",
         "abc":"def"
     },
     {
         "name":"yahoo",
         "addr":"yahoo.com:80",
         "abc":"zyx"
     },
     {
         "name":"redis",
         "addr":"localhost:6379",
         "abc":"efg"
     }
 ]
```

#### pagination

See below inline comments on how to use pagination.

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

#### inherit_attributes
When you process a payload which is further nested, you may want to inherit the attributes above it as well. 
Note: currently only supported when start_key is used.
```yaml
---
name: example
apis: 
  - event_type: exampleAbc
    url: http://127.0.0.1:8887/azure-insights-vm.json
    inherit_attributes: true
    start_key:
      - value>timeseries
    lazy_flatten:
      - timeseries
```

#### add_attribute

Under add_attribute, we create a `linkToIncident` variable which will use the value from an existing key that is available in the current payload. The attribute that would be seen in this sample returned would be `"links.incident_id"` and it is dynamically populated into the placeholder `${links.incident_id}`.

```yaml
  - name: allAlerts
    event_type: alertSample
    url: ${var:alert_API_URL}&start_date=${timestamp:datetimeutc-8hr}
    headers:
      X-Api-Key: ${lf:rest_API_key}
      Content-Type: application/json
    start_key:
      - violations
    store_lookups:
      storeInc: links.incident_id
    value_parser:
      links.incident_id: "[0-9]+"
    add_attribute:
      linkToIncident: https://alerts.newrelic.com/accounts/12345/incidents/${links.incident_id}/violations
```

#### ignore_output

Supply ignore_output: true to an api, to completed ignore its output, this is useful when creating lookups.

```yaml
---
name: myFlexConfig
apis:
  - event_type: getConsumers
    url: http://127.0.0.1:8887/consumers.json
    store_lookups:
      consumers: consumers
    ignore_output: true
  - event_type: getClusters
    url: http://127.0.0.1:8887/clusters.json
    store_lookups:
      clusters: clusters
    ignore_output: true
  - event_type: SomeEvent
    url: http://127.0.0.1:8887/${lookup:consumers}/${lookup:clusters}.json
```

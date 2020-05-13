# `file`

The `file` API lets you retrieve information from any `JSON` or `CSV` file.

* [Basic usage](#Basicusage)
* [Configuration properties](#Configurationproperties)
* [Advanced usage](#Advancedusage)

##  <a name='Basicusage'></a>Basic usage

```yaml
name: linuxExample
apis:
  - name: jsonFileExample
    file: /tmp/payload.json
```

On Windows, make sure `file` uses backslashes "\\" to separate directories in the file path.

```yaml
name: windowsExample
apis:
  - name: jsonFileExample
    file: C:\Program Files\My App\tmp\payload.json
```

Other than that, there are no differences between Linux and Windows features.

`file` accepts a path to any JSON or CSV file. If the file does not have an extension, it's processed as JSON by default. 
To process CSV files, the `.csv` extension is required.

##  <a name='Configurationproperties'></a>Configuration properties

The following table describes the properties of the `file` API.

| Name | Type | Default | Description |
|---:|:---:|:---:|---|
| `set_header` | array of strings | `[]` | Name and number of columns Flex should extract data from. Only applies to CSV files. If this property is not set, the first row of data is used as the header.

##  <a name='Advancedusage'></a>Advanced usage

The `file` API can be used alongside other Flex functions. In the following example, we use some Flex data processing [functions](../basics/functions.md).

```yaml
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: /tmp/etcdSelf.json
    start_key:
      - leaderInfo
    rename_keys:
      startTime: timestamp
    custom_attributes:
      env: production
```

Given the following `/tmp/etcdSelf.json` file:

```json
{
    "id": "eca0338f4ea31566",
    "leaderInfo": {
        "leader": "8a69d5f6b7814500",
        "startTime": 1588232295,
        "uptime": 3600
    },
    "name": "node3",
    "recvAppendRequestCnt": 5944,
    "recvBandwidthRate": 570.6254930219969,
    "recvPkgRate": 9.00892789741075,
    "sendAppendRequestCnt": 0,
    "state": "StateFollower"
}
``` 

The generated sample contains the following attributes:

```
"leader": "8a69d5f6b7814500",
"timestamp": 1588232295,
"uptime": 3600
"env": "production"
```

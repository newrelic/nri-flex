# Gemini Enterprise Flex Integration

This flex integration provides an example of collecting Gemini Enterprise
metrics for a given GCP **project** using the
[GCP Cloud Monitoring API v3](https://docs.cloud.google.com/monitoring/api/v3).

## Getting Started

To get started with this Flex integration, you will need the following:

* A host machine running the
  [Infrastructure Agent](https://docs.newrelic.com/docs/infrastructure/infrastructure-agent/linux-installation/package-manager-install/)
* A
  [GCP Service Account](https://developers.google.com/identity/protocols/oauth2/service-account#creatinganaccount)
* A
  [service account key](https://docs.cloud.google.com/iam/docs/keys-create-delete#creating)
  JSON file for the GCP service account

**NOTE:** The service account must have the permissions below or be assigned a
role that has these permissions:

* monitoring.timeSeries.list

## Installation

1. Install the
   [GCP gcloud CLI](https://docs.cloud.google.com/sdk/docs/install-sdk) on the
   host running the infrastructure agent.
1. Copy the service account key JSON file to the host running the integration.
1. Copy the `gemini-enterprise.yml` Flex configuration file to
   `/etc/newrelic-infra/integrations.d` on the host running the infrastructure
   agent.
1. Run Flex manually **one time** with the `gcp-auth` block uncommented. This
   will authenticate for all subsequent executions of the `gcloud` CLI via Flex
   on the host machine.

   ```bash
   [sudo] /opt/newrelic-infra/newrelic-integrations/bin/nri-flex --verbose --pretty --config_file gemini-enterprise.yml
   ```

   Comment out the `gcp-auth` block after this is done successfully.
1. [Restart the infrastructure agent](https://docs.newrelic.com/docs/infrastructure/infrastructure-agent/manage-your-agent/start-stop-restart-infrastructure-agent/).

## Configuration

The `gemini-enterprise.yml` Flex configuration requires the following
configuration parameters:

| Parameter Name          | Description                                                                                    |
| ----------------------- | ---------------------------------------------------------------------------------------------- |
| `service_account_email` | The email address of the service account used for authentication                               |
| `service_account_key`   | The path on the local file system to the service account key JSON file used for authentication |
| `project_id`            | The GCP project ID containing the Gemini Enterprise application(s) to monitor                  |

These parameters are substituted in to the shell commands used to authenticate
with GCP and the API URLs used to collect metrics as appropriate.

Additionally, the `INSIGHTS*` environment variables in the `env:` section of the
configuration can be used to remove all infrastructure agent metadata added to
each payload forwarded to New Relic. To do this, perform the following steps:

1. Uncomment the line containing the `INSIGHTS_API_KEY` environment variable
   definition and replace `<ingest_key>` with your New Relic license key.
1. Uncomment the line containing the `INSIGHTS_URL` environment variable
   definition and replace `<account_id>` with your New Relic account ID.

## Data Collected

This example Flex integration collects the following metrics.

| Metric Type                                                              | Filter                                                     | Kind  | Type  | Unit | Ingest Delay |
| ------------------------------------------------------------------------ | ---------------------------------------------------------- | ----- | ----- | ---- | ------------ |
| `serviceruntime.googleapis.com/api/request_count`                        | `resource.labels.service="discoveryengine.googleapis.com"` | DELTA | INT64 | 1    | 300s         |
| `discoveryengine.googleapis.com/data_stores_regional`                    |                                                            | GAUGE | INT64 | 1    | 120s         |
| `discoveryengine.googleapis.com/engines_regional`                        |                                                            | GAUGE | INT64 | 1    | 120s         |
| `discoveryengine.googleapis.com/quota/search_requests_regional/exceeded` |                                                            | DELTA | INT64 | 1    | 150s         |
| `discoveryengine.googleapis.com/quota/search_requests_regional/limit`    |                                                            | GAUGE | INT64 | 1    | 150s         |
| `discoveryengine.googleapis.com/quota/search_requests_regional/usage`    |                                                            | DELTA | INT64 | 1    | 150s         |
| `discoveryengine.googleapis.com/quota/data_stores_regional/exceeded`     |                                                            | DELTA | INT64 | 1    | 150s         |
| `discoveryengine.googleapis.com/quota/data_stores_regional/limit`        |                                                            | GAUGE | INT64 | 1    | 150s         |
| `discoveryengine.googleapis.com/quota/data_stores_regional/usage`        |                                                            | GAUGE | INT64 | 1    | 150s         |

**NOTE:** See
["Value types and metric kinds"](https://docs.cloud.google.com/monitoring/api/v3/kinds-and-types)
for a description of the Kind, Type, and Unit values listed above.

### Using the Data

The metrics above are sent to New Relic as infrastructure
[event data](https://docs.newrelic.com/docs/data-apis/understand-data/new-relic-data-types/#event-data).
The provided events and attributes are listed below.

**NOTE:** Some of the text below is sourced from the
[GCP Cloud Monitoring API v3 documentation](https://docs.cloud.google.com/monitoring/api/v3).

#### `GeminiEnterpriseSample`

In this example Flex configuration, the `GeminiEnterpriseSample` event is used
to report the following metrics:

* The `serviceruntime.googleapis.com/api/request_count` metric type for the
  `discoveryengine.googleapis.com` service
* The `discoveryengine.googleapis.com/data_stores_regional` metric type
* The `discoveryengine.googleapis.com/engines_regional` metric type.

The table below shows what attributes are included in each event based on the
[metric type](https://docs.cloud.google.com/monitoring/api/v3/metric-model#metric_types)
(metadata attributes added by the infrastructure agent are not included below).

| Attribute Name       | Data Type | Metric Type(s)                                    | Description                                     |
| -------------------- | --------- | ------------------------------------------------- | ----------------------------------------------- |
| `metricType`         | string    | all                                               | The GCP Cloud Monitoring v3 metric type         |
| `projectId`          | string    | all                                               | The GCP project ID of the monitored resource    |
| `resourceType`       | string    | all                                               | The monitored resource type                     |
| `resourceLocation`   | string    | all                                               | The location of the monitored resource          |
| `regionalLocation`   | string    | `discoveryengine.googleapis.com/*`                | The multi region identifier                     |
| `value`              | number    | all                                               | The numeric measured value                      |
| `interval.startTime` | number    | all                                               | The measurement interval start time, in seconds |
| `interval.endTime`   | number    | all                                               | The measurement interval end time, in seconds   |

**NOTE:** For events reported for the
`serviceruntime.googleapis.com/api/request_count` metric type, the
`resourceLocation` value represents the service specific notion of location.
This can be a name of a zone or region. If a service does not have any notion of
zones then 'global' can be used.

#### `GeminiEnterpriseQuotaSample`

In this example Flex configuration, the `GeminiEnterpriseQuotaSample` event is
used to report the metrics for all `discoveryengine.googleapis.com/quota/*`
metric types.

The table below shows what attributes are included in each event based on the
metric type (metadata attributes added by the infrastructure agent are not
included below).

| Attribute Name       | Data Type | Metric Type(s)                                                      | Description                                                |
| -------------------- | --------- | ------------------------------------------------------------------- | ---------------------------------------------------------- |
| `metricType`         | string    | all                                                                 | The GCP Cloud Monitoring v3 metric type                    |
| `projectId`          | string    | all                                                                 | The GCP project ID of the monitored `Location` resource    |
| `resourceType`       | string    | all                                                                 | The monitored resource type                                |
| `resourceLocation`   | string    | all                                                                 | The location of the monitored resource                     |
| `limitName`          | string    | all                                                                 | The limit name                                             |
| `regionalLocation`   | string    | all                                                                 | The multi region identifier                                |
| `method`             | string    | discoveryengine.googleapis.com/quota/search_requests_regional/usage | The method                                                 |
| `value`              | number    | all                                                                 | The numeric measured value                                 |
| `interval.startTime` | number    | all                                                                 | The measurement interval start time, in seconds            |
| `interval.endTime`   | number    | all                                                                 | The measurement interval end time, in seconds              |

**NOTE:** The `resourceType` value for all `GeminiEnterpriseQuotaSample` events
is `discoveryengine.googleapis.com/Location`.

### Collecting Additional Metrics

This Flex configuration can be extended to collect additional Gemini Enterprise
metrics.

#### Supported value types and metric kinds

The example API calls in this Flex configuration can be used to collect all
[supported value type and metric kind combinations](https://docs.cloud.google.com/monitoring/api/v3/kinds-and-types#kind-type-combos)
except the `DISTRIBUTION` value type combinations.

**NOTE:** Use caution when visualizing or alerting on `CUMULATIVE` metrics as
the value of cumulative metrics monotonically increase over time.

#### Adding a new metric

To add a new metric to the Flex configuration, first inspect the response from
the GCP Cloud Monitoring API v3 call for the metric. You can do this using the
[APIs explorer](https://docs.cloud.google.com/monitoring/api/apis-explorer),
an API client tool, or a CLI tool like `curl` or `wget`. In the example below,
we will use `curl` for the API call and the
[GCP `gcloud` CLI](https://docs.cloud.google.com/sdk/docs/install-sdk) for
authentication.

1. Compose the
   [monitoring filter](https://docs.cloud.google.com/monitoring/api/v3/filters)
   for the desired metric.
1. Activate your service account credentials using `gcloud` by running the
   following command:

   ```bash
   gcloud auth activate-service-account MY_SERVICE_ACCOUNT_EMAIL --key-file=MY_PATH_TO_SERVICE_ACCOUNT_KEY_FILE
   ```

   Replace `MY_SERVICE_ACCOUNT_EMAIL` with the email address of your service
   account and `MY_PATH_TO_SERVICE_ACCOUNT_KEY_FILE` with the file system path
   to your service account key file.
1. Generate an authentication token by running the following command:

   ```bash
   TOKEN=$(gcloud auth print-access-token)
   ```
1. Run the following `curl` command:

   ```bash
   curl -s -H "Authorization: Bearer $TOKEN" \
   "https://monitoring.googleapis.com/v3/projects/MY_PROJECT_ID/timeSeries?interval.startTime=MY_START_TIME&interval.endTime=MY_END_TIME&filter=MY_MONITORING_FILTER"
   ```

   Replace the placeholders above as follows.

   | Placeholder            | Value                                                                                       |
   | ---------------------- | ------------------------------------------------------------------------------------------- |
   | `MY_PROJECT_ID`        | Your GCP project ID                                                                         |
   | `MY_START_TIME`        | The start of the time interval to query using RFC3339 format: `YYYY-MM-DDTHH:MM:SSZ`        |
   | `MY_END_TIME`          | The end of the time interval to query using RFC3339 format: `YYYY-MM-DDTHH:MM:SSZ`          |
   | `MY_MONITORING_FILTER` | The monitoring filter with `"` characters escaped using `\"` and spaces escaped with `%20`  |

   **NOTE:**
   * You may need to escape other special shell characters in
     `MY_MONITORING_FILTER` using the `\` character.
   * To format the JSON output in the API response, pipe the result of the
     `curl` command into the `jq` utility (if available) using `| jq .`.

Sample output for the
`discoveryengine.googleapis.com/quota/search_requests_regional/usage` metric
type is shown below.

```json
{
  "timeSeries": [
    {
      "metric": {
        "labels": {
          "method": "",
          "regional_location": "us",
          "limit_name": "SearchRequestsPerMinutePerProjectPerRegion"
        },
        "type": "discoveryengine.googleapis.com/quota/search_requests_regional/usage"
      },
      "resource": {
        "type": "discoveryengine.googleapis.com/Location",
        "labels": {
          "project_id": "PROJECT_ID",
          "location": "global"
        }
      },
      "metricKind": "DELTA",
      "valueType": "INT64",
      "points": [
        {
          "interval": {
            "startTime": "2026-03-12T12:53:10.001Z",
            "endTime": "2026-03-12T12:54:10Z"
          },
          "value": {
            "int64Value": "1"
          }
        }
      ]
    }
  ],
  "unit": "1"
}
```

Using the response from the GCP Cloud Monitoring API v3 call, generate a Flex
`url` API configuration as follows:

1. Start with the following template:

   ```yml
    - name: MY_API_CONFIG_NAME
      event_type: MY_EVENT_NAME
      url: MY_API_URL
      headers:
        Authorization: Bearer ${var:access_token}
      start_key:
      - timeSeries>points
      inherit_attributes: true
      timestamp_conversion:
        interval.startTime: TIMESTAMP::RFC3339
        interval.endTime: TIMESTAMP::RFC3339
      rename_keys:
        parent.0.metric.type: metricType
        parent.0.metric.labels.regional_location: regionalLocation
        parent.0.resource.type: resourceType
        value.MY_TYPED_VALUE: value
      sample_include_filter:
      - metricType: MY_METRIC_TYPE
      add_attribute:
        timestamp: ${interval.endTime}
      remove_keys:
        - parent.0.api.StatusCode
        - parent.0.metricKind
        - parent.0.valueType
        - parent.0.unit
   ```
1. Replace `MY_API_CONFIG_NAME` with a name for the configuration. For example,
   `discovery-engine-quota-search-requests-regional-usage`.
1. Replace `MY_EVENT_NAME` with a name for the event generated for the metric.
   The two event names used in this example are `GeminiEnterpriseSample` and
   `GeminiEnterpriseQuotaSample` but using these names is not required.
1. Replace `MY_API_URL` with the API URL used in the `curl` command above with
   the following changes:

   * Replace the project ID in the URL path with the text `${var:project_id}`.
   * Replace the timestamp for the `startTime` query parameter with the text
     `${timestamp:datetimeutctz-MY_START_OFFSET}`.
   * Replace the timestamp for the `endTime` query parameter with the text
     `${timestamp:datetimeutctz-MY_END_OFFSET}`.
   * Leave the space characters in the monitoring filter escaped using '%20'.
   * Remove the escaping for the other special characters in the monitoring
     filter.

   Set the `MY_START_OFFSET` and `MY_END_OFFSET` values as follows:

   1. Refer to the metric description for the
      [`discoveryengine` metric](https://docs.cloud.google.com/monitoring/api/metrics_gcp_d_h#gcp-discoveryengine)
      being collected. It should indicate the time interval after which the
      sampled data is visible. For example, the description for the
      `quota/search_requests_regional/usage` metric indicates the following:

      > After sampling, data is not visible for up to 150 seconds.

      Use this value for `MY_END_OFFSET` using the suffix `s` to indicate
      seconds. For example, `150s`.
   1. Add `60` to the end offset and use this value for `MY_START_OFFSET`. For
      example, `210s`.

   This ensures that data points are not missed due to requesting data too
   early.

   **NOTE:**
   * This means that for a given minute of data, that data will not be
     visible in New Relic until at least `MY_END_OFFSET` seconds after the given
     minute. In the table in the ["Data Collected"](#data-collected) section
     above, this value is listed as the "Ingest Delay".
   * This method works for data sampled at 30 and 60 second intervals. Data
     sampled at different intervals may require separate Flex configurations
     with a different `interval` value. Make sure to check the sampling interval
     in the metric type description prior to adding a new metric.

1. Replace `MY_METRIC_TYPE` with the
   [metric type](https://docs.cloud.google.com/monitoring/api/v3/metric-model#metric_types)
   used in your monitoring filter.
1. Replace `MY_TYPED_VALUE` with one of the following as appropriate for the
   value type of the metric type.

   * `boolValue`
   * `int64Value`
   * `doubleValue`
   * `stringValue`

1. For each additional property in the JSON response, determine whether you want
   to keep the property or remove the property and add the appropriate property
   expression to the `rename_keys` map or the `remove_keys` list, respectively.

   For example, the response shown above for the
   `discoveryengine.googleapis.com/quota/search_requests_regional/usage` metric
   type has the following additional properties, identified by the nested key
   expressions below.

   * `parent.0.metric.limit_name`
   * `parent.0.metric.method`
   * `parent.0.resource.location`
   * `parent.0.resource.project_id`

   In this example, we want to keep all the additional properties.

Following the above steps for the response shown above for the
`discoveryengine.googleapis.com/quota/search_requests_regional/usage` metric,
the final Flex `url` API configuration would be the following:

```yml
- name: discovery-engine-quota-search-requests-regional-usage
  event_type: GeminiEnterpriseQuotaSample
  url: https://monitoring.googleapis.com/v3/projects/${var:project_id}/timeSeries?interval.startTime=${timestamp:datetimeutctz-210s}&interval.endTime=${timestamp:datetimeutctz-150s}&filter=metric.type="discoveryengine.googleapis.com/quota/search_requests_regional/usage"
  headers:
    Authorization: Bearer ${var:access_token}
  start_key:
  - timeSeries>points
  inherit_attributes: true
  timestamp_conversion:
    interval.startTime: TIMESTAMP::RFC3339
    interval.endTime: TIMESTAMP::RFC3339
  rename_keys:
    parent.0.metric.type: metricType
    parent.0.metric.labels.limit_name: limitName
    parent.0.metric.labels.method: method
    parent.0.metric.labels.regional_location: regionalLocation
    parent.0.resource.type: resourceType
    parent.0.resource.labels.location: resourceLocation
    parent.0.resource.labels.project_id: projectId
    value.int64Value: value
  sample_include_filter:
  - metricType: discoveryengine.googleapis.com/quota/search_requests_regional/usage
  add_attribute:
    timestamp: ${interval.endTime}
  remove_keys:
    - parent.0.api.StatusCode
    - parent.0.metricKind
    - parent.0.valueType
    - parent.0.unit
```

## Known Limitations

Please note the following known limitations when using this Flex configuration.

* Pagination of GCP Cloud Monitoring API v3 calls is not supported. This means
  that for queries which result in more than 100,000 data points, only the first
  100,000 data points will be collected.
* As indicated above, this Flex configuration works for data sampled at 30 and
  60 second intervals. Data sampled at different intervals may require separate
  Flex configurations with a different `interval` value. Make sure to check the
  sampling interval in the metric type description prior to adding a new metric.

## Troubleshooting

To troubleshoot issues with the integration refer to the
[Flex troubleshooting documentation](https://github.com/sdewitt-newrelic/nri-flex/blob/master/docs/troubleshooting.md).

# Big Query Flex Integration

This flex integration runs bq queries via GCP CLI and reports the output to New Relic. Most useful information is found under the [INFORMATION_SCHEMA](https://cloud.google.com/bigquery/docs/information-schema-intro) view.


## Pre-requirements
* [Infrastructure Agent](https://docs.newrelic.com/docs/infrastructure/infrastructure-agent/linux-installation/package-manager-install/) installed
* [GCP CLI](https://cloud.google.com/sdk/docs/install#linux) installed & configured on host running infrastructure agent
* [GCP Service Account](https://developers.google.com/identity/protocols/oauth2/service-account#creatinganaccount) created with json key file downloaded

Service account must have the following permissions:
* Bigquery.tables.get
* Bigquery.tables.list
* Bigquery.routines.get
* Bigquery.routines.list
* Bigquery.jobs.listAll



## Installation
1. Copy service account json key file to host running the integration
2. Copy `bq.yml` under `/etc/newrelic-infra/integrations.d`
3. Authenticate the CLI with the service account key file:

```bash
gcloud auth login --cred-file=/path/to/key.json
```

4. Run Flex manually **one time** with the `bq-auth` block uncommented. This will authenticate for all subsequent executions of the CLI via Flex.

```bash
[sudo] /opt/newrelic-infra/newrelic-integrations/bin/nri-flex --verbose --pretty --config_file bq.yml
```

Comment out the `bq-auth` block after this is done successfully.

5. [Restart the infrastructure agent](https://docs.newrelic.com/docs/infrastructure/infrastructure-agent/manage-your-agent/start-stop-restart-infrastructure-agent/)



## Configuration
The bq configuration requires the service account email, GCP project id, and region. These values are substituted dynamically into each bq CLI command ran, so any additional queries added can follow the same format as the examples provided.

Additionally, the polling interval can be set at the top (in seconds), and the `INSIGHTS*` environment variables can be used to remove all infrastructure agent metadata tacked onto each bq payload forwarded to New Relic. These are configured with an ingest key and an account id within the URL variable.

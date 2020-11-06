# Flex running as a GCP Function

Golang GCP Functions require the source tree be zipped and deployed to GCP.

## Installation

### Quick start

- Clone nri-flex
- cd into the clone directory
- Create `flexConfigs/` directory 
- Add your Flex configuration files to `flexConfigs`
- Create a ZIP file to upload to GCP 
- Upload ZIP file to GCP Bucket
- Define GCP Function
- Define GCP Scheduler Job to run Function

```bash
git clone https://github.com/newrelic/nri-flex.git
cd nri-flex
mkdir flexConfigs
cp <YOUR_CONFIG_FILES> flexConfigs/
zip -r ../nri-flex.zip * -x vendor/\* bin/\*
```

### Details
- Clone Flex
  - `git clone https://github.com/newrelic/nri-flex.git`
- Change into the cloned Flex's directory
  - `cd nri-flex`
- Create a `flexConfigs` subdirectory
  - `mkdir flexConfigs`
- Add your Flex configuration files to `flexConfigs`
  - `cp <YOUR_CONFIG_FILES> flexConfigs/`
- Create the ZIP file distribution
  - `zip -r ../nri-flex.zip * -x vendor/\* bin/\*`
- Create a GCP Bucket if you don't already have one
- Upload the ZIP file to your GCP Bucket
- [Create a GCP Function in the GCP Console](https://cloud.google.com/functions/docs/deploying/console)
  - Configuration
    - Function name: your choice
    - Region: your choice
    - Trigger
      - Trigger type: `Cloud Pub/Sub`
      - Select a Cloud Pub/Sub topic: choose `CREATE A TOPIC` and follow the directions
      - Click `Save`
    - Expand `VARIABLES, NETWORKING, AND ADVANCED SETTINGS` and choose the `ENVIRONMENT VARIABLES` tab
      - Create two `Runtime environment variables`
        - Name: INSIGHTS_API_KEY    Value: YOUR_INSIGHTS_API_KEY
        - Name: INSIGHTS_URL        Value: https://insights-collector.newrelic.com/v1/accounts/284929/events
        - (Optional) Name: VERBOSE  Value: true
  - Code
    - Runtime: `Go 1.13`
    - Entry point: `FlexPubSub`
    - Source code: `ZIP from Cloud Storage`
    - Cloud Storage location: click `Browser` and browser to the previously uploaded ZIP file
    - Click `DEPLOY`
- [Create a new Cloud Scheduler job in the GCP Console](https://cloud.google.com/scheduler/docs/creating)
  - Name: your choice
  - [Frequency](https://cloud.google.com/scheduler/docs/configuring/cron-job-schedules): 
  - Timezone: your choice
  - Target: `Pub/Sub`
    - Topic: The Topic name you created above
    - Payload: `{}` (empty curly braces, an empty JSON Object)
  

## Trouble  shooting
You can test your setup and configuration locally if you have Go installed on your machine using

    go run test/serverless/gcp/function.go

To see the available command line parameters use

    go run test/serverless/gcp/function.go -help
    
The Test tab for the GCP Function in the GCP Console will run Flex once and show you the log, set the `VERBOSE` environment variable to get more detail.

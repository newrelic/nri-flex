# Running Flex as a Serverless Function

The following documentation explains how to configure the New Relic FLEX integration to run as a serverless function. Using the `nri-flex` as a stand-alone ingest tool can be beneficial to monitor database performance in AWS RDS, check on external status for 3rd-party services, and many more functions.

## Note

Due to the nature of serverless functions, some, possibly all, FLEX command `run` functions many not be available in the function's context. APIs that are native to FLEX, like `database` and `url`, should function as expected.

## Deploying with Serverless Framework

- Ensure serverless framework is installed - https://serverless.com/framework/docs/providers/aws/guide/installation/
- Within serverless.yml update any parameters as required
- Add the latest linux flex binary into the /pkg folder
- Add your Flex configs into /pkg/flexConfigs/
- Deploy with `sls deploy -v`
- Alternatively create a package and deploy with your own method

## Deploying with AWS native tools

The following provides an example of deploying the `nri-flex` Golang executable as an Amazon Web Services (AWS) Lambda function using native AWS tools. It is important to first make sure that you have configured the AWS command line interface (CLI) locally [link](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)
for the required account or are know how to use the AWS Lambda web interface to upload a `.zip` file and configure the handle value. The following will describe how to do all this via the CLI.

### Create the AWS Lambda

The following steps assume the AWS lambda function already exists, has the correct permissions to execute, and that you have the appropriate rights to manage the lambda.

#### Lambda Requirements

- Set to use Go 1.x (or appropriate language version of the time)
- Set the `Handler` to be `nri-flex` or follow the CLI steps below to do the same
- Environment Variables
  - Make sure that the New Relic insights insert key maps to the correct New Relic account ID
  - INSIGHTS_API_KEY: [NR-INSIGHTS-INSERT-KEY](https://docs.newrelic.com/docs/apis/get-started/intro-apis/types-new-relic-api-keys#insights-insert-key)
  - INSIGHTS_URL: https://insights-collector.newrelic.com/v1/accounts/NR-ACCOUNT-ID/events

### Build Lambda Layout

```bash
cd ~/
mkdir flex-lambda
cd flex-lambda
mkdir deploy

# download the latest linux x86_64 tagged release of nri-flex
# https://github.com/newrelic/nri-flex/releases
curl -L https://github.com/newrelic/nri-flex/releases/download/v<version>/nri-flex_<version>_Linux_x86_64.tar.gz -o nri-flex.tar.gz

# extract just the nri-flex binary, unless you want to review the examples and other files provided
# browse the files using 'tar ztf nri-flex.tar.gz'
tar -xf nri-flex.tar.gz nri-flex

# move the nri-flex binary into the root of the deploy folder
mv nri-flex deploy/nri-flex

cd deploy
mkdir pkg
mkdir pkg/flexConfigs

# use your favorite text editor to create a flex YAML file
nano pkg/flexConfigs/test.yaml
```

Add the following YAML contents to the file and save it.

```yaml
name: httpExample
custom_attributes: # applies to all apis
  myCustAttr: myCustVal
apis:
  - event_type: LambdaFlexSample
    url: https://jsonplaceholder.typicode.com/todos/1
    custom_attributes:
      nestedCustAttr: nestedCustVal # nested custom attributes specific to each api
```

### Create the Lambda `.zip` File

The following is used to create the `.zip` of the folder structure. On Windows, the `.zip` file creation is recommend to be done per the AWS documentation which references a `build-lambda-zip.exe` tool [link](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html). It appears that the native Windows `.zip` format may not work for this, but it is unclear at this time. If possible, use the Linux subsystem on Windows to create the `.zip` file.

```bash
# while still in ~/flex-lambda/deploy/
zip ../nri-flex.zip -r .
# the command should out the the following if all is correct
  adding: nri-flex (deflated 64%)
  adding: pkg/ (stored 0%)
  adding: pkg/flexConfigs/ (stored 0%)
  adding: pkg/flexConfigs/test.yaml (deflated 38%)
```

### Update the Lambda's Handle

Note, this step is only needed if it wasn't done during the creation of the Lambda function.

```bash
aws lambda update-function-configuration --function-name <function-name> --handler nri-flex
```

### Update the Lambda's Code

```bash
aws lambda update-function-code --function-name <function-name> --zip-file fileb://../nri-flex.zip
```

### Set CloudWatch to Execute

Go to CloudWatch and configure an Event (schedule rule) which calls the Lambda function as needed by your requirements.

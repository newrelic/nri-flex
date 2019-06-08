# Running Flex as a Lambda

### Deploying with Serverless Framework
- Ensure serverless framework is installed - https://serverless.com/framework/docs/providers/aws/guide/installation/
- Within serverless.yml update any parameters as required
- Add the latest linux flex binary into the /pkg folder
- Add your Flex configs into /pkg/flexConfigs/
- Deploy with `sls deploy -v`
- Alternatively create a package and deploy with your own method
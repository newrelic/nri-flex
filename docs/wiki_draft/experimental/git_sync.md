# Config Sync with Git (Experimental)

**Disclaimer**: this function is bundled in Alpha status. That means that it is not yet supported by New Relic.

To dynamically sync integrations with github several methods are available.

- CLI Flags
```
./nri-flex -verbose -git_user myUser -git_token 13nasdasj13jadf -git_repo https://github.com/myUser/my-config-repo
```
- Environment Variables
```
GIT_REPO=https://github.com/myUser/my-config-repo
GIT_USER=myUser
GIT_TOKEN=13nasdasj13jadf
```
- Setting in nri-flex-config.yml
```yaml
### /etc/newrelic-infra/integrations.d/nri-flex-config.yml
integration_name: com.newrelic.nri-flex

instances:
  - name: nri-flex
    command: metrics
    arguments:
      git_repo: https://github.com/userName/repoName
      git_user: userName
      git_token: abcd
      # fargate: true ## default false
      # container_discovery: true ## default false
      # container_discovery_dir: "anotherDir" default "flexContainerDiscovery" 
      # config_file: "../myConfigFile.yml" ## default "" - run only a single specific config file
      # config_dir: "anotherConfigDir/" ## default "flexConfigs/"
      # event_limit: 500 ## default 500
      # insights_api_key: abc
      # insights_url: https://insights...
      # insights_output: output the payload to stdout
    # labels:
      # owner: cloud
```
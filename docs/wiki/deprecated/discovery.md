### Flex Auto Container Discovery 

> ⚠️ **Notice** ⚠️: this document contains a deprecated functionality that is still
> supported by New Relic for backwards compatibility. However, we encourage you to
> use the improved, fully-supported [container auto-discovery mechanism for On-Host Integrations](https://docs.newrelic.com/docs/integrations/host-integrations/installation/container-auto-discovery). 

- Flex has the capability to auto discovery containers in your surrounding environment, and dynamically monitor them regardless of changing IP addresses and ports
- Requires access to /var/run/docker.sock (same as the New Relic Infrastructure Agent, so it is convenient to bake Flex into the newrelic/infrastructure image or your k8s daemon set)

#### Configuration
```
  -container_discovery
        Enable container auto discovery
  -container_discovery_dir string
        Set directory of auto discovery config files (default "flexContainerDiscovery/")
  -docker_api_version string
        Force Docker client API version
```

### Versions
- [V2 Container Discovery](https://github.com/newrelic/nri-flex/wiki/Service-Discovery#V2-Container-Discovery)
- [V1 Container Discovery](https://github.com/newrelic/nri-flex/wiki/Service-Discovery#V1-Container-Discovery)
---

### V2 Container Discovery
- V2 simplifies discovery by allowing everything to be defined by your Flex configuration file
- You are able to place config files within your standard flexConfigs folder, the container_discovery flag still needs to be enabled to true.
- These configs are prefixed with "cd-" eg. "flexConfigs/`cd-redis.yml`"
- Use ${auto:ip}, ${auto:host} and ${auto:port} as your substitution flags, see below example.

Example: consider `container_discovery` is enabled and a yml file called cd-redis.yml is placed within your flexConfigs folder.
```
name: redis
apis: 
  - name: redis
    entity: ${auto:ip}
    commands: 
      - dial: ${auto:ip}:${auto:port}
        run: "info\r\n"
        split_by: ":"
```

Flex by default will look for containers with the name that contains "redis", the name is taken from the text after "cd-".
To override this or further configure how you would like to target containers consider the below example.
`the below are all optional parameters`
```
name: redis
container_discovery:
  mode: contains # default: contains, other options: prefix, suffix, contains, regex
  type: container # default: container, other options: image
  target: some-other-container-name # default: text after the "cd-" in the config file name
  ip_mode: private # default: automatically set, other options: public, private
apis: 
  - name: redis
    entity: ${auto:ip}
    commands: 
      - dial: ${auto:ip}:${auto:port}
        run: "info\r\n"
        split_by: ":"
```
---
### V1 Container Discovery

Also note as an alternative in a Kubernetes environment, service environment variables are passed to every container, so your flex container will have access. Therefore you could even simply use environment variable substitution in a normal config file without using this container discovery feature. Environment variables can be passed into config files by simply using a double dollar sign eg. $$MY_ENVIRONMENT_VAR

- Add a label or annotation that contains the keyword - "flexDiscovery"
- For Kubernetes add it as an environment variable
- To that same label or annotation add a flex discovery configuration eg. "t=redis,c=redis,tt=img,tm=contains"
- Complete example                                              flexDiscoveryRedis="t=redis,c=redis,tt=img,tm=contains"
- You could have varying configs on one container as well like flexDiscoveryRedis1, flexDiscoveryZookeeper etc.
- Flex Container Discovery Configs are placed within "flexContainerDiscovery/" directory 
- For an example see "flexContainerDiscovery/redis.yml" 
- Use `${auto:host}` and `${auto:port}` anywhere in your config, this will dynamically be substituted per container discovered
- This makes it possible to have multiple containers re-use the same config with different ip/port configurations
- You can also instead opt to apply the same environment variable(s) to the Flex container itself, this allows you to monitor other containers without modifying your existing environment
- If you have configs that overlap between something set on a container, and something set on the Flex container itself, the Flex container will take precedence.
- Flex uses the environment variables that are presented to the Docker Engine, if you have modified the discovery environment variable(s) after deployment for whatever reason, the Docker Engine will be unaware of these changes.

#### Flex Discovery Configuration Parameters
- tt=targetType - are we targeting an img = image or cname = containerName? (default "img")
- t=target - the keyword to target based on the targetType eg. "redis"
- tm=targetMode - contains, prefix or regex to match the target (default "contains")
- c=config - which config file will we use to create the dynamic configs from eg. "redis" .yml (defaults to the "target value")
- p=port - force set a chosen target port
- ip=ipMode - default private can be set to public
- If config is nil, use the target (t), as the yaml file to look up, eg. if target (t) = redis, lookup the config (c) redis.yml if config not set
### Flex Auto Container Discovery 

> ⚠️ **Note**: this is deprecated functionality that is still provided for backwards compatibility. We encourage you to use the improved, fully-supported [container auto-discovery mechanism for On-Host Integrations](https://docs.newrelic.com/docs/integrations/host-integrations/installation/container-auto-discovery). 

Flex has the capability to auto discovery containers in your surrounding environment, and dynamically monitor them regardless of changing IP addresses and ports.

This requires access to `/var/run/docker.sock`, same as the New Relic Infrastructure Agent, so it is convenient to include Flex in the `newrelic/infrastructure` image or your k8s daemon set.

## Configuration
```
  -container_discovery
        Enable container auto discovery
  -container_discovery_dir string
        Set directory of auto discovery config files (default "flexContainerDiscovery/")
  -docker_api_version string
        Force Docker client API version
```

## Versions
- [V2 Container Discovery](#V2-Container-Discovery)
- [V1 Container Discovery](V1-Container-Discovery)


### V2 Container Discovery

V2 simplifies discovery by allowing everything to be defined by your Flex configuration file: you can place config files within your standard `flexConfigs` folder. The `container_discovery` flag still needs to be enabled to true.

These configs are prefixed with `cd-` (for example, `flexConfigs/cd-redis.yml`). Use `${auto:ip}`, `${auto:host}` and `${auto:port}` as your substitution flags.

In this example, `container_discovery` is enabled and a YAML file called `cd-redis.yml` is placed in your `flexConfigs` folder:
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

This tells Flex to look for containers with the name that contains `redis` by default, as the name is taken from the text after `cd-`.

To override this behavior or further configure how you would like to target containers, consider the example below (all are optional parameters):

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
### V1 Container Discovery

In a Kubernetes environment, service environment variables can be passed to every container, so your Flex container can access them. You can use environment variable substitution in a normal config file without using the container discovery feature. Environment variables can be passed to config files by using a double dollar sign (for example, `$$MY_ENVIRONMENT_VAR`).

1. Add a label or annotation that contains the keyword - `flexDiscovery`.
2. For Kubernetes, add it as an environment variable.
3. Add a flex discovery configuration to that same label or annotation, for example, `flexDiscoveryRedis="t=redis,c=redis,tt=img,tm=contains`

You could have varying configs on one container, like `flexDiscoveryRedis1`, `flexDiscoveryZookeeper`, etc.

Flex Container Discovery Configs are placed within "flexContainerDiscovery/" directory. For an example see `flexContainerDiscovery/redis.yml`.

Use `${auto:host}` and `${auto:port}` anywhere in your config: these are dynamically substituted per container discovered, and makes it possible to have multiple containers re-use the same config with different ip/port configurations.

You can also opt to apply the same environment variable(s) to the Flex container itself, which allows you to monitor other containers without modifying your existing environment. If you have configs that overlap between something set on a container, and something set on the Flex container itself, the Flex container will take precedence.

Flex uses the environment variables that are presented to the Docker Engine; if you have modified the discovery environment variable(s) after deployment for whatever reason, the Docker Engine will be unaware of these changes.

#### Flex Discovery Configuration Parameters
- `tt=targetType` - Target type (default is `img`)
- `t=target` - Keyword to target based on the targetType (for example, `redis`)
- `tm=targetMode` - Prefix or regex to match the target (default is `contains`)
- `c=config` - Config file used to create the dynamic configs (defaults to the `target` value)
- `p=port` - Target port
- `ip=ipMode` - IP mode. Default is `private` - it can be set to `public`

If a config is `nil`, use the target (`t`), as the YAML file to look up. For example, if target (`t`) = `redis`, look up the config (`c`) `redis.yml` if the config is not set.
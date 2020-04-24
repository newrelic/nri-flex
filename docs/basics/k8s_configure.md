# Configure Flex in Kubernetes

## New Relic container images

If you're running services in Kubernetes we recommend you use the official container images from New Relic.

There are three container images you can use depending on how you plan on using Flex: 

- To run Flex only to monitor services running in Kubernetes, use the [newrelic/infrastructure](https://hub.docker.com/r/newrelic/infrastructure) container image. This image only contains the Infrastructure agent, and the Docker and Flex integrations. You will not be able to perform service discovery or use other New Relic integrations.
- To run Flex alongside other New Relic integrations, use the [newrelic/infrastructure-bundle](https://hub.docker.com/r/newrelic/infrastructure-bundle) container image. This adds all the other For more information, see [New Relic integrations](https://docs.newrelic.com/docs/integrations/kubernetes-integration/link-apps-services/monitor-services-running-kubernetes).
- If you also want to monitor your Kubernetes cluster, use the [newrelic/infrastructure-k8s](https://hub.docker.com/r/newrelic/infrastructure-k8s) container image. This image adds all the integrations, including the [Kubernetes integration](https://docs.newrelic.com/docs/integrations/kubernetes-integration/get-started/introduction-kubernetes-integration). 

## Configure Flex in Kubernetes

The preferred way to configure Flex (and other New Relic integrations) in Kubernetes is through a Config Map.

Configurations can [go embedded](#AddyourFlexconfigurationtointegrations.d) in the Infrastructure agent config file, or live as standalone config files that you can test separately and [link](#Linktoaseparateconfigurationfile) from the main integrations config file. 

What approach to follow is up to you.

* [Add your Flex configuration to configMap](#AddyourFlexconfigurationtoconfigMap)
* [Link to a separate configuration file](#Linktoaseparateconfigurationfile)
* [Add multiple Flex configurations](#Addmultipleflexconfigurations)

### <a name='AddyourFlexconfigurationtoconfigMap'></a>Add your configuration to `configMap`

Assuming that you have installed the Kubernetes integration following the [official instructions](https://docs.newrelic.com/docs/integrations/kubernetes-integration/installation/kubernetes-installation-configuration#install), you should have the Infrastructure agent running in your cluster, as well as two configMaps:

- `nri-default-integration-cfg` : config map used to enable the New Relic Kubernetes integration. Can be removed if you do not want to use this integration.
- `nri-integration-cfg-example` : config map used to enable other integrations. You should this config map to enable Flex and other New Relic integrations.

To enable Flex, create a `data` section in the config map and add an Infrastructure agent integration configuration under `data`:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nri-integration-cfg-example
  namespace: default
data:
  nri-flex.yml: |
    integrations:
      - name: nri-flex
        config:
          name: example
          apis:
            - event_type: ExampleSample
              url: https://my-host:8443/admin/metrics.json
``` 

> To get a quick, first picture of a Flex configuration file, follow our [basic step-by-step tutorial](../../basic-tutorial.md) or check existing config files under [/examples](../../examples).

### <a name='Linktoaseparateconfigurationfile'></a>Link to a separate configuration file

You can store the Flex configuration files in a separate YAML file (for example, after [developing and testing a config file](../development.md)) and reference it by replacing `config` with the `config_template_path` property, which contains the path to a Flex configuration file.

The linking equivalent of the previous config map would contain:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nri-integration-cfg-example
  namespace: default
data:
  nri-flex.yml: |
    integrations:
      - name: nri-flex
        config_template_path: /path/to/flex/integration.yml # Reference to a separate Flex config file
``` 
While this works, you would either have to create a new config map with the contents of the external file and mount it into the container, or modify the original container image to include the file.

### <a name='#Addmultipleflexconfigurations'></a>Add multiple Flex configurations

You can also add multiple Flex configurations to the config map. There are two ways to do this:

1. Add multiple Flex entries to the same configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nri-integration-cfg-example
  namespace: default
data:
  nri-flex.yml: |
    integrations:
      - name: nri-flex-1
        config:
          name: example
          apis:
            - event_type: ExampleSample
              url: https://my-host:8443/admin/metrics.json
      - name: nri-flex-2
        config:
          name: other-example
          apis:
            - event_type: OtherExampleSample
              url: https://my-other-host:8080/metrics.json
```

2. Add distinct config map entries, one per configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nri-integration-cfg-example
  namespace: default
data:
  nri-flex-1.yml: |
    integrations:
      - name: nri-flex-1
        config:
          name: example
          apis:
            - event_type: ExampleSample
              url: https://my-host:8443/admin/metrics.json
  nri-flex-2.yml: |
    integrations:
      - name: nri-flex-2
        config:
          name: other-example
          apis:
            - event_type: OtherExampleSample
              url: https://my-other-host:8080/metrics.json
```

There are no distinctive advantages of using one over the other: choose according to your needs and preferences.

For an example of a deployment manifest, see [the k8s example](../../examples/nri-flex-k8s.yml).
# Configure Flex under Kubernetes

If you're running under Kubernetes we recommend you use the official container images from NewRelic.

There are 3 container images you can use depending on you use-case. 
- If you plan on using only Flex to monitor service in Kubernetes you can use the container image [newrelic/infrastructure](https://hub.docker.com/r/newrelic/infrastructure).
  
  This image only contains the Infrastructure agent, and the Docker and Flex integrations. You will not be able to do service discovery and use other NewRelic integrations.
- If you plan on using Flex and perhaps other NewRelic integration you can use the container image [newrelic/infrastructure-bundle](https://hub.docker.com/r/newrelic/infrastructure-bundle)

   This image builds on top on the [newrelic/infrastructure](https://hub.docker.com/r/newrelic/infrastructure) image and adds all the other NewRelic support integrations. See this [page](https://docs.newrelic.com/docs/integrations/kubernetes-integration/link-apps-services/monitor-services-running-kubernetes) for more information.

- If you also want to monitor your Kubernetes cluster you can use the container image [newrelic/infrastructure-k8s](https://hub.docker.com/r/newrelic/infrastructure-k8s).
  
  This image builds on top of [newrelic/infrastructure-bundle](https://hub.docker.com/r/newrelic/infrastructure-bundle) and add the Kubernetes integration. See this [page](https://docs.newrelic.com/docs/integrations/kubernetes-integration/get-started/introduction-kubernetes-integration) for more information.

The preferred way to configure Flex (and other NewRelic integrations) in Kubernetes is to use a ConfigMap.

Configurations can [go embedded](#AddyourFlexconfigurationtointegrations.d) in the Infrastructure agent co file, or live as stand-alone config files that you can test separately and [link](#Linktoaseparateconfigurationfile) from the main integrations config file. What approach to follow is up to you.

* [Add your Flex configuration to configMap](#AddyourFlexconfigurationtoconfigMap)
* [Link to a separate configuration file](#Linktoaseparateconfigurationfile)
* [Add multiple Flex configurations](#Addmultipleflexconfigurations)

## <a name='AddyourFlexconfigurationtoconfigMap'></a>Add your configuration to `configMap`

Assuming you have installed the Infrastructure K8s integration (which installs the Infra agent) using the [instructions](https://docs.newrelic.com/docs/integrations/kubernetes-integration/installation/kubernetes-installation-configuration#install) (you can skip step 1 if you don't want to monitor you cluster) 
you should have the Infra agent running in your cluster and 2 configMaps:

- `nri-default-integration-cfg` : configMap used to enable the New Relic Kubernetes integration. Can be removed if you do not want to use this integration.
- `nri-integration-cfg` : configMap used to enable other integrations. You should this configMap to enable Flex and other NewRelic integrations.

To enable Flex, create a `data` section in the configMap and under data add an Infrastructure agent integration configuration "file":
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nri-integration-cfg # aimed to be safely overridden by users
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

> To get a quick, first picture of a Flex configuration file, you can start following our [basic, step-by-step tutorial](../../basic-tutorial.md) or check existing config files under [/examples](../../examples).

## <a name='Linktoaseparateconfigurationfile'></a>Link to a separate configuration file

You can store the Flex configuration in a separate YAML file (for example, after [developing and testing a config file](../development.md)) and reference it by replacing `config` with `config_template_path` property, which contains the path of the Flex config file.

This way, the equivalent of previous configMap from the previous section would contain:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nri-integration-cfg # aimed to be safely overridden by users
  namespace: default
data:
  nri-flex.yml: |
    integrations:
      - name: nri-flex
        config_template_path: /path/to/flex/integration.yml # Reference to a separate Flex config file
``` 
While this would work, you would either have to create a new configMap with the contents of the external file and mount it into the container, 
or you would have to modify the original container image to include the file. For this reason we recommend you embedded you configuration as explained in the previous section. 

## <a name='#Addmultipleflexconfigurations'></a>Add multiple Flex configurations

You can also add multiple Flex configurations to the configMap. There are two ways to do this:

- add multiple Flex entries to the same "file"

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nri-integration-cfg # aimed to be safely overridden by users
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

- add distinct configMap entries, one per "file"

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nri-integration-cfg # aimed to be safely overridden by users
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

Which one you use is up to you. There are no distinctive advantages of using one over the other.

For an example of a deployment manifest see [example](../../examples/nri-flex-k8s.yml).
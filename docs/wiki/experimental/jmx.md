# JMX Queries

> ⚠️ **Notice** ⚠️: this document contains a deprecated functionality that is still
> provided for backwards compatibility. However, we encourage you to
> use the improved, fully-supported [nri-jmx On-Host Integration](http://github.com/newrelic/nri-jmx). 

Flex can uses [nrjmx](http://github.com/newrelic/nrjmx) to send JMX requests that can be processed later.

### JMX Queries

- Flex is able to utilise [nrjmx](https://github.com/newrelic/nrjmx) to run any JMX queries you need.
- A copy of the nrjmx jar file is kept [here](https://github.com/newrelic/nri-flex/tree/master/nrjmx)

#### Install JMX On Linux

Below are the steps to install the nrjmx file and .jar manually to enable support for JMX. It is important to substitute the `{desired version}` below with the version of `nri-flex` you want to use.

##### Steps
1. Download the package manually
   * Find the [release](https://github.com/newrelic/nri-flex/releases) needed
   * `curl -L https://github.com/newrelic/nri-flex/releases/download/{desired version}/nri-flex-linux-{desired version}.tar.gz -o nri-flex-linux-{desired version}.tar.gz`
2. Extract the compressed file
   * `tar -xvf nri-flex-linux-{desired version}.tar.gz`
3. Use the extracted directory as the working directory
   * `cd nri-flex-linux-{desired version}`
4. Run `install_linux.sh --jmx` with privileges, typically `sudo`, to create files in `/var/db/newrelic-infra`
5. (optional) Copy over a JMX configuration to begin monitoring
   * In the `nri-flex-linux-{desired version}/examples/flexConfigs/` folder there are different JMX examples which can be used as a reference
   * After copying an example config to `/var/db/newrelic-infra/custom-integrations/flexConfigs/`, edit the file according to your monitoring needs and save the file
   * Restart the infrastructure agent, `sudo systemctl restart newrelic-infra`

#### JMX options available are:
```
domain
user
pass
host
port
key_store
key_store_pass
trust_store
trust_store_pass
```

Below is a simple tomcat example.

```
name: tomcatFlex
global:
  jmx:
      host: "127.0.0.1"
      port: "9001"
apis: 
  - name: tomcatThreads
    event_type: tomcatThreadSample
    ### note "keep_keys" will do the inverse, if you want all metrics remove the keep keys blocks completely
    ### otherwise tailor specific keys you would like to keep, this uses regex for filtering
    ### this is useful for keeping key metrics
    keep_keys: ###
      - bean
      - maxThreads
      - connectionCount
    commands: 
      - run: Catalina:type=ThreadPool,name=*
        output: jmx
  - name: tomcatRequest
    event_type: tomcatRequestSample
    keep_keys:
      - bean
      - bytesSent
      - bytesReceived
      - errorCount
      - requestCount
    commands: 
      - run: Catalina:type=GlobalRequestProcessor,name=*
        output: jmx
  - name: manager
    event_type: tomcatManagerSample
    keep_keys:
      - bean
      - errorCount
      - requestCount
    commands: 
      - run: Catalina:type=GlobalRequestProcessor,name=*
        output: jmx
  - name: datasource
    event_type: tomcatDatasourceSample
    keep_keys:
      - bean
      - numActive
      - numIdle
```

***


#### Global Config that is passed down
```
base_url
user
pass
proxy
timeout
headers:
 headerX: valueX
jmx:
* domain
* user
* pass
* host
* port
* key_store
* key_store_pass
* trust_store
* trust_store_pass
```

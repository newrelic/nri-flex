# `url`

The `url` API allows retrieving information from an HTTP endpoint. 

## Basic usage

```yaml
---
name: example
apis:
  - event_type: ExampleSample
    url: http://my-host:8080/admin/metrics.json
    headers:
      accept: application/json
```

The above Flex configuration retrieves a JSON from the provided URL, containing a set of metrics.
Please notice that the `url` API may be followed by a `headers` section, which allows specifying
the HTTP headers.

## `POST` / `PUT` HTTP methods

You can use the `method` and `payload` properties to specify a `POST` or `PUT` request with its
body. 

```
---
name: httpPostExample 
apis: 
  - name: httpPost
    url: https://jsonplaceholder.typicode.com/posts
    method: POST
    payload: > 
      {"title": "foo","body": "bar","userId": 1}
```

## Configuring your HTTPS connections with `tls_config`

If you are using TLS endpoints with self-signed certificates, you may need to specify a `tls_config`
section with any of the following items:

| Name | Type | Default | Description |
|---:|:---:|:---:|---|
| `enable` | Bool | `false` | Set it to `true` to enable a custom TLS configuration for your HTTPS connection. It is used in conjunction with the rest of properties in this table |
| `insecure_skip_verify` | Bool | false | Set to `true` to skip the verification of TLS certificates for your HTTPS endpoint |
| `ca` | string | _empty_ | Provide the Certificate Authority PEM certificate, in case your HTTPS endpoint has self-signed certificates. |  

Example:

```yaml
---
name: example
apis:
  - event_type: ExampleSample
    url: https://my-host:8443/admin/metrics.json
    headers:
      accept: application/json
    tls_config:
      enable: true
      ca: /etc/bundles/my-ca-cert.pem
```

## Specifying a common base URL with `base_url`

If you have to query the same host multiple times, you may want to set up a common
`base_url` global field. Then you only have to provide the URL path in the rest of
`url` sections.

Example

```yaml
name: consulFlex
global:
    base_url: http://consul-host/v1/
    headers:
      X-Consul-Token: my-root-consul-token
apis: 
  - event_type: ConsulHealthSample
    url: health/service/consul
  - event_type: ConsulCheckSample
    url: health/state/any
  - event_type: ConsulMemberSample
    url: agent/members
```

## `url` with `cache` for later processing

The URL invocations are cached, so you can process them later without having to query it
repeatedly.

In the following example, the NGINX status endpoint is invoked, and its output is
retrieved from the cache for its later extraction of meaningful data:

```yaml
name: nginxFlex
apis:
  - name: nginxStub
    url: http://127.0.0.1/nginx_status
  - name: nginx
    event_type: NginxSample
    commands:
      - cache: http://127.0.0.1/nginx_status
        split_output: Active
        regex_matches:
          - expression: Active connections:\s(\S+)
            keys: [net.connectionsActive]
          - expression: \s?(\d+)\s(\d+)\s(\d+)
            keys: [net.connectionsAcceptedPerSecond, net.handledPerSecond, net.requestsPerSecond]
          - expression: Reading:\s(\d+)\s\S+\s(\d+)\s\S+\s(\d+)
            keys: [net.connectionsReading, net.connectionsWriting, net.connectionsWaiting]
    math:
      net.connectionsDroppedPerSecond: ${net.connectionsAcceptedPerSecond} - ${net.handledPerSecond} 
```
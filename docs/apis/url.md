# url API

The `url` API allows you to retrieve information from an HTTP endpoint.

- [Basic usage](#Basicusage)
- [Use POST/PUT methods with a body](#UsePOSTPUTmethodswithabody)
- [Configure your HTTPS connections](#ConfigureyourHTTPSconnections)
- [Specify a common base URL](#SpecifyacommonbaseURL)
- [URL with cache for later processing](#URLwithcacheforlaterprocessing)
- [Include response headers on sample](#ReturnResponseHeaders)

## <a name='Basicusage'></a>Basic usage

```yaml
name: example
apis:
  - event_type: ExampleSample
    url: http://my-host:8080/admin/metrics.json
    headers:
      accept: application/json
```

The above Flex configuration retrieves a JSON file containing a set of metrics from the provided URL. Note that the `url` key can be followed by a `headers` section, which allows specifying HTTP headers.

## <a name='UsePOSTPUTmethodswithabody'></a>Use POST/PUT methods with a body

To specify a `POST` or `PUT` request with a body, use the `method` and `payload` properties.

```yaml
name: httpPostExample
apis:
  - name: httpPost
    url: https://jsonplaceholder.typicode.com/posts
    method: POST
    payload: >
      {"title": "foo","body": "bar","userId": 1}
```

## <a name='ConfigureyourHTTPSconnections'></a>Configure your HTTPS connections

When using TLS endpoints with self-signed certificates, define a `tls_config` section with any of the following items:

|                   Name |  Type  | Default | Description                                                                                                  |
| ---------------------: | :----: | :-----: | ------------------------------------------------------------------------------------------------------------ |
|               `enable` |  bool  | `false` | Set to `true` to enable custom TLS configuration. Requires `ca` to be defined if enabled.                    |
| `insecure_skip_verify` |  bool  | `false` | Set to `true` to skip the verification of TLS certificates.                                                  |
|                   `ca` | string | _Empty_ | The Certificate Authority PEM certificate, in case your HTTPS endpoint has self-signed certificates.         |
|                 `cert` | string | _Empty_ | PEM encoded certificate (must be used with `key`), in case your HTTPS endpoint has self-signed certificates. |
|                  `key` | string | _Empty_ | PEM encoded key (must be used with `cert`), in case your HTTPS endpoint has self-signed certificates.        |

### TLS configuration example:

```yaml
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

## <a name='SpecifyacommonbaseURL'></a>Specify a common base URL

When you have to query several different URLs, specifying a `base_url` under `global` can be quite helpful, as it allows you to provide URL path segment in `url` fields instead of full URLs.

### Base URL example

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

## <a name='URLwithcacheforlaterprocessing'></a>URL with cache for later processing

URL invocations are cached to avoid having to query them repeatedly. Use `cache` under `command` to read cached data.

In this example, the NGINX status endpoint is invoked, and the output is retrieved from the cache for later processing:

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
            keys:
              [
                net.connectionsAcceptedPerSecond,
                net.handledPerSecond,
                net.requestsPerSecond,
              ]
          - expression: Reading:\s(\d+)\s\S+\s(\d+)\s\S+\s(\d+)
            keys:
              [
                net.connectionsReading,
                net.connectionsWriting,
                net.connectionsWaiting,
              ]
    math:
      net.connectionsDroppedPerSecond: ${net.connectionsAcceptedPerSecond} - ${net.handledPerSecond}
```

## <a name='ReturnResponseHeaders'></a>Include response headers on sample

To include response headers on the metric sample set `return_headers` attribute to true.

### Return headers example

```yaml
name: example
apis:
  - name: ExampleSample
    url: https://my-host:8443/admin/metrics/1
    return_headers: true
```

Given the following output for each metric:

```json
{
  "event_type": "ExampleSample",
  "id": 1,
  "completed": "true",
  "api.StatusCode": 200,
  "api.header.Access-Control-Allow-Credentials": "[true]",
  "api.header.Age": "[4459]",
  "api.header.Content-Type": "[application/json; charset=utf-8]",
  "api.header.Date": "[Mon, 25 May 2020 16:23:53 GMT]",
  "api.header.Expires": "[-1]",
  "api.header.Retry-Count": "[0]"
}
```

# Prometheus Integrations (Exporters)

> ⚠️ **Notice** ⚠️: this document contains a deprecated functionality that is still
> supported by New Relic for backwards compatibility. However, we encourage you to
> use the improved, fully-supported [New Relic Prometheus OpenMetrics integration for Docker and Kubernetes](https://docs.newrelic.com/docs/integrations/prometheus-integrations). 

Prometheus Support - [Exporters](https://prometheus.io/docs/instrumenting/exporters/)

- Supports all Prometheus exporters
- Flex will attempt to flatten all Prometheus metrics for you to save on events being generated, however you may need to do some minor additional configuration (below) to get the best output
- With the automatically flattened event, histogram & summary, count & sum values are retained
- If you would like the full qauntiles and buckets, consider flagging on histogram, and/or summary to true
- Target the /metrics endpoint and set your desired configuration, see further below for options
- To quickly find out what metrics may need to be in their own samples or merged into the main sample, set -force_log and view the /metrics endpoint you are targetting

```
# This is a configuration example explaining all the possible options on setting up prometheus exporter metrics to be ingested

---
name: redisFlex # https://github.com/oliver006/redis_exporter
apis: 
  - name: redis
    url: http://localhost:9121/metrics 
    prometheus: 
      enable: true
      # flattened_event: "redisCustomSample" ## as the api name is "redis" the default event type created will be redisSample to override the flattened event type use this parameter 
      # histogram: true  ## default false - enable histogram metrics, as the api name is "redis" the event type created will be redisHistogramSample
      # histogram_event: nginxMyCustomHistogramEvent ## override the auto event type for histogram metrics
      # summary: true ## default false - enable summary metrics, as the api name is "redis" the event type created will be redisSummarySample
      # summary_event: nginxMyCustomHistogramEvent ## override the auto event type for summary metrics
      # go_metrics: false ## default false - the exporters internal go metrics
      # unflatten: true ## use with caution as it can generate a lot of events, this will generate an event per prometheus metric
    # sample_filter:
      # - .*: GAUGE ## remove all gauge metrics
      # - .*: COUNTER ## remove all counter metrics
```
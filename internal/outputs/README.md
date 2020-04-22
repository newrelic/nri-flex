# outputs
--
    import "github.com/newrelic/nri-flex/internal/outputs"


## Usage

#### func  GetMetricBatches

```go
func GetMetricBatches() [][]*metric.Set
```
GetMetricBatches batch metrics by entity with a maximum batch size defined by
'InsightBatchSize' config.

#### func  InfraIntegration

```go
func InfraIntegration() error
```
InfraIntegration Creates Infrastructure SDK Integration

#### func  SendBatchToInsights

```go
func SendBatchToInsights(metrics []*metric.Set) error
```
SendBatchToInsights - Send processed events to insights.

#### func  SendToMetricAPI

```go
func SendToMetricAPI() error
```
SendToMetricAPI - Send processed events to insights

#### func  StatusSample

```go
func StatusSample()
```
StatusSample creates flexStatusSample

#### func  StoreJSON

```go
func StoreJSON(samples []interface{}, path string)
```
function to store samples as a JSON object at specified path

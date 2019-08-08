# outputs
--
    import "github.com/newrelic/nri-flex/internal/outputs"


## Usage

#### func  InfraIntegration

```go
func InfraIntegration()
```
InfraIntegration Creates Infrastructure SDK Integration

#### func  InfraRemoteEntity

```go
func InfraRemoteEntity()
```
InfraRemoteEntity Creates Infrastructure Remote Entity

#### func  SendToInsights

```go
func SendToInsights()
```
SendToInsights - Send processed events to insights loop through integration
entities as there could be multiple that have been set when posted they are
batched by entity

#### func  SendToMetricAPI

```go
func SendToMetricAPI()
```
SendToMetricAPI - Send processed events to insights

#### func  StatusSample

```go
func StatusSample()
```
StatusSample creates flexStatusSample

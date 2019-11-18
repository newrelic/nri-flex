# inputs
--
    import "github.com/newrelic/nri-flex/internal/inputs"


## Usage

#### func  NetDialWithTimeout

```go
func NetDialWithTimeout(dataStore *[]interface{}, command load.Command, dataSample *map[string]interface{}, api load.API, processType *string)
```
NetDialWithTimeout performs network dial without timeout

#### func  ParseJMX

```go
func ParseJMX(dataStore *[]interface{}, dataInterface interface{}, command load.Command, dataSample *map[string]interface{})
```
ParseJMX Processes JMX Data

#### func  ParseReader

```go
func ParseReader(in io.Reader, ch chan<- *dto.MetricFamily) error
```
ParseReader consumes an io.Reader and pushes it to the MetricFamily channel. It
returns when all MetricFamilies are parsed and put on the channel.

#### func  ProcessQueries

```go
func ProcessQueries(dataStore *[]interface{}, yml *load.Config, apiNo int)
```
ProcessQueries processes database queries

#### func  Prometheus

```go
func Prometheus(dataStore *[]interface{}, input io.Reader, cfg *load.Config, api *load.API)
```
Prometheus from http io

#### func  RunCommands

```go
func RunCommands(dataStore *[]interface{}, yml *load.Config, apiNo int)
```
RunCommands executes the given commands to create one merged sampled

#### func  RunFile

```go
func RunFile(dataStore *[]interface{}, cfg *load.Config, apiNo int)
```
RunFile runs file read data collection

#### func  RunHTTP

```go
func RunHTTP(dataStore *[]interface{}, doLoop *bool, yml *load.Config, api load.API, reqURL *string)
```
RunHTTP Executes HTTP Requests

#### func  RunScpWithTimeout

```go
func RunScpWithTimeout(dataStore *[]interface{}, cfg *load.Config, api load.API, apiNo int)
```
RunScpWithTimeout performs scp with timeout

#### func  SetJMXCommand

```go
func SetJMXCommand(dataStore *[]interface{}, runCommand *string, command load.Command, api load.API, config *load.Config)
```
SetJMXCommand Add parameters to JMX call

#### type Family

```go
type Family struct {
	//Time    time.Time
	Name    string                         `json:"name"`
	Help    string                         `json:"help"`
	Type    string                         `json:"type"`
	Metrics map[int]map[string]interface{} `json:"metrics,omitempty"` // Either metric or summary.
}
```

Family mirrors the MetricFamily proto message.

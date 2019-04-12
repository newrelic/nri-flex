# parser
--
    import "github.com/newrelic/nri-flex/internal/parsers"


## Usage

#### func  ParseJMX

```go
func ParseJMX(dataInterface interface{}, dataStore *[]interface{}, command load.Command, dataSample *map[string]interface{})
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
func ProcessQueries(api load.API, dataStore *[]interface{})
```
ProcessQueries processes database queries

#### func  Prometheus

```go
func Prometheus(input io.Reader, dataStore *[]interface{}, api *load.API)
```
Prometheus from http io

#### func  RunCommands

```go
func RunCommands(yml *load.Config, api load.API, dataStore *[]interface{})
```
RunCommands executes the given commands to create one merged sampled

#### func  RunHTTP

```go
func RunHTTP(doLoop *bool, yml *load.Config, api load.API, reqURL *string, dataStore *[]interface{})
```
RunHTTP Executes HTTP Requests

#### func  SetJMXCommand

```go
func SetJMXCommand(runCommand *string, command load.Command, api load.API, config *load.Config)
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

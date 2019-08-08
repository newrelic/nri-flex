# integration
--
    import "github.com/newrelic/nri-flex/internal/integration"


## Usage

#### func  HandleRequest

```go
func HandleRequest(ctx context.Context, event interface{}) (string, error)
```
HandleRequest Handles incoming lambda request

#### func  Lambda

```go
func Lambda()
```
Lambda handles lambda invocation

#### func  LambdaCheck

```go
func LambdaCheck() bool
```
LambdaCheck check if Flex is running within a Lambda and insights url and api
key has been supplied

#### func  RunFlex

```go
func RunFlex(mode string)
```
RunFlex runs flex if mode is "" run in default mode

#### func  SetDefaults

```go
func SetDefaults()
```
SetDefaults set flex defaults

#### func  SetEnvs

```go
func SetEnvs()
```
SetEnvs set environment variable argument overrides

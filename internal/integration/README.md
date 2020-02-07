# integration
--
    import "github.com/newrelic/nri-flex/internal/integration"


## Usage

#### func  HandleLambda

```go
func HandleLambda()
```
HandleLambda handles lambda invocation

#### func  HandleRequest

```go
func HandleRequest(ctx context.Context, event interface{}) (string, error)
```
HandleRequest Handles incoming lambda request

#### func  IsLambda

```go
func IsLambda() bool
```
IsLambda check if Flex is running within a Lambda.

#### func  RunFlex

```go
func RunFlex(runMode FlexRunMode)
```
RunFlex runs flex.

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

#### func  ValidateLambdaConfig

```go
func ValidateLambdaConfig() error
```
ValidateLambdaConfig: while running within a Lambda insights url and api key are
required.

#### type FlexRunMode

```go
type FlexRunMode int
```

FlexRunMode is used to switch the mode of running flex.

```go
const (
	// FlexModeDefault is the usual way of running flex.
	FlexModeDefault FlexRunMode = iota
	// FlexModeLambda is used when flex is running within a lambda.
	FlexModeLambda
	// FlexModeTest is used when running tests.
	FlexModeTest
)
```

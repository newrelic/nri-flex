# runtime
--
    import "."


## Usage

#### func  CommonPostInit

```go
func CommonPostInit()
```
Post-initialization common to all runtime types here

#### func  CommonPreInit

```go
func CommonPreInit()
```
Pre-initialization common to all runtime types here

#### func  RunFlex

```go
func RunFlex(instance Instance) error
```
Common run (once) function

#### type Default

```go
type Default struct {
}
```


#### func (*Default) SetConfigDir

```go
func (i *Default) SetConfigDir(s string)
```

#### type Function

```go
type Function struct {
}
```

GCP Function runtime

#### func (*Function) SetConfigDir

```go
func (i *Function) SetConfigDir(s string)
```

#### type Instance

```go
type Instance interface {
	SetConfigDir(string)
	// contains filtered or unexported methods
}
```

Serverless runtimes must implement this

#### func  GetDefaultRuntime

```go
func GetDefaultRuntime() Instance
```
Get a server-based, default, runtime

#### func  GetFlexRuntime

```go
func GetFlexRuntime() Instance
```
Get the first available runtime type, defaults to the server-based (Linux |
Windows) Default type

#### func  GetTestRuntime

```go
func GetTestRuntime() Instance
```
Get the test runtime

#### type Lambda

```go
type Lambda struct {
}
```

AWS Lambda runtime

#### func (*Lambda) FlexAsALambdaHandler

```go
func (i *Lambda) FlexAsALambdaHandler(ctx context.Context, event interface{}) (string, error)
```
FlexAsALambdaHandler receives the incoming lambda request, from the AWS
perspective this is the entry point

#### func (*Lambda) SetConfigDir

```go
func (i *Lambda) SetConfigDir(s string)
```

#### type Test

```go
type Test struct {
}
```


#### func (*Test) SetConfigDir

```go
func (i *Test) SetConfigDir(s string)
```

# discovery
--
    import "github.com/newrelic/nri-flex/internal/discovery"


## Usage

#### func  CreateDynamicContainerConfigs

```go
func CreateDynamicContainerConfigs(containers []types.Container, files []os.FileInfo, path string, ymls *[]load.Config)
```
CreateDynamicContainerConfigs Creates dynamic configs for each container

#### func  ExecContainerCommand

```go
func ExecContainerCommand(containerID string, command []string) (string, error)
```
ExecContainerCommand execute command against a container

#### func  FindFlexContainerID

```go
func FindFlexContainerID(read string)
```
FindFlexContainerID detects if Flex is running within a container and sets the
ID

#### func  Processes

```go
func Processes()
```
Processes loops through tcp connections and returns the corresponding process
and connection information

#### func  Readln

```go
func Readln(r *bufio.Reader) (string, error)
```
Readln from bufioReader

#### func  Run

```go
func Run(configs *[]load.Config)
```
Run discover containers

#### type ProcessNetworkStat

```go
type ProcessNetworkStat struct {
	Name string
	Data string
}
```

ProcessNetworkStat x

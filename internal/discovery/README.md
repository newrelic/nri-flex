# discovery
--
    import "github.com/newrelic/nri-flex/internal/discovery"


## Usage

#### func  CreateDynamicContainerConfigs

```go
func CreateDynamicContainerConfigs(containers []types.Container, files []os.FileInfo, path string, ymls *[]load.Config)
```
CreateDynamicContainerConfigs Creates dynamic configs for each container

#### func  FindFlexContainerID

```go
func FindFlexContainerID()
```
FindFlexContainerID detects if Flex is running within a container and sets the
ID

#### func  Readln

```go
func Readln(r *bufio.Reader) (string, error)
```
Readln from bufioReader

#### func  Run

```go
func Run(cfg *[]load.Config)
```
Run discover containers

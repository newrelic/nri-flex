# logger
--
    import "github.com/newrelic/nri-flex/internal/logger"


## Usage

#### func  Flex

```go
func Flex(logType string, err error, message string, createEvent bool)
```
Flex generic log handler to support force logging and creating additional events
for insights debugging

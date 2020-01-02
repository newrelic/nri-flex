# config
--
    import "github.com/newrelic/nri-flex/internal/config"


## Usage

#### func  FetchData

```go
func FetchData(apiNo int, yml *load.Config) []interface{}
```
FetchData fetches data from various inputs Also handles paginated responses for
HTTP requests (tested against NR APIs)

#### func  FetchLookups

```go
func FetchLookups(cfg *load.Config, i int) bool
```
FetchLookups x

#### func  GitCheckout

```go
func GitCheckout(w *git.Worktree)
```
GitCheckout git checkout

#### func  GitClone

```go
func GitClone(dir string, u *url.URL) error
```
GitClone git clone

#### func  GitPull

```go
func GitPull(dir string) error
```
GitPull git pull

#### func  LoadFile

```go
func LoadFile(configs *[]load.Config, f os.FileInfo, path string) error
```
LoadFile loads a single Flex config file

#### func  LoadFiles

```go
func LoadFiles(configs *[]load.Config, files []os.FileInfo, path string)
```
LoadFiles Loads Flex config files

#### func  ReadYML

```go
func ReadYML(yml string) (load.Config, error)
```
ReadYML Unmarshals yml files

#### func  Run

```go
func Run(yml load.Config)
```
Run Action each config file

#### func  RunFiles

```go
func RunFiles(configs *[]load.Config)
```
RunFiles Processes yml files

#### func  RunVariableProcessor

```go
func RunVariableProcessor(i int, cfg *load.Config)
```
RunVariableProcessor substitute store variables into specific parts of config
files

#### func  SubEnvVariables

```go
func SubEnvVariables(strConf *string)
```
SubEnvVariables substitutes environment variables into config Use a double
dollar sign eg. $$MY_ENV_VAR to subsitute that environment variable into the
config file Can be useful with kubernetes service environment variables

#### func  SubLookupFileData

```go
func SubLookupFileData(configs *[]load.Config, config load.Config)
```
SubLookupFileData substitutes data from lookup files into config

#### func  SubTimestamps

```go
func SubTimestamps(strConf *string)
```
SubTimestamps - return timestamp/date/datetime of current date/time with
optional adjustment in various format

#### func  SyncGitConfigs

```go
func SyncGitConfigs(customDir string) bool
```
SyncGitConfigs Clone git repo if already exists, else pull latest version

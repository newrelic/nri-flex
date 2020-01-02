# load
--
    import "github.com/newrelic/nri-flex/internal/load"


## Usage

```go
const (
	DefaultSplitBy     = ":"                      // unused currently
	DefaultTimeout     = 10000 * time.Millisecond // 10 seconds, used for raw commands
	DefaultDialTimeout = 1000                     // 1 seconds, used for dial
	DefaultPingTimeout = 5000                     // 5 seconds
	DefaultHANA        = "hdb"
	DefaultPostgres    = "postgres"
	DefaultMSSQLServer = "sqlserver"
	DefaultMySQL       = "mysql"
	DefaultOracle      = "ora"
	DefaultVertica     = "vertica"
	DefaultJmxPath     = "./nrjmx/"
	DefaultJmxHost     = "127.0.0.1"
	DefaultJmxPort     = "9999"
	DefaultJmxUser     = "admin"
	DefaultJmxPass     = "admin"
	DefaultShell       = "/bin/sh"
	DefaultLineLimit   = 255
	Public             = "public"
	Private            = "private"
	Jmx                = "jmx"
	Img                = "img"
	Image              = "image"
	TypeContainer      = "container"
	TypeJSON           = "json"
	TypeColumns        = "columns"
	Contains           = "contains"
)
```

```go
var AWSExecutionEnv string
```
AWSExecutionEnv AWS execution environment

```go
var ContainerID string
```
ContainerID current container id

```go
var DiscoveredProcesses map[string]string
```
DiscoveredProcesses discovered processes

```go
var Entity *integration.Entity
```
Entity Infrastructure SDK Entity

```go
var FlexStatusCounter = struct {
	sync.RWMutex
	M map[string]int
}{M: make(map[string]int)}
```
FlexStatusCounter count internal metrics

```go
var Hostname string
```
Hostname current host

```go
var IngestData interface{}
```
IngestData store ingested data

```go
var Integration *integration.Integration
```
Integration Infrastructure SDK Integration

```go
var IntegrationName = "com.newrelic.nri-flex" // IntegrationName Name

```

```go
var IntegrationNameShort = "nri-flex" // IntegrationNameShort Short Name

```

```go
var IntegrationVersion = "Unknown-SNAPSHOT" // IntegrationVersion Version

```

```go
var IsFargate bool
```
IsFargate basic check if running on fargate

```go
var IsKubernetes bool
```
IsKubernetes basic check if running on k8s

```go
var LambdaName string
```
LambdaName if running on lambda add name from AWS_LAMBDA_FUNCTION_NAME

```go
var Logrus = logrus.New()
```
Logrus create instance of the logger

```go
var MetricsStore = struct {
	sync.RWMutex
	Data []Metrics
}{}
```
MetricsStore for Dimensional Metrics to store data and lock and unlock when
needed

```go
var StartTime int64
```
StartTime time Flex starts in Nanoseconds

#### func  MakeTimestamp

```go
func MakeTimestamp() int64
```
MakeTimestamp creates timestamp in milliseconds

#### func  MetricsStoreAppend

```go
func MetricsStoreAppend(metrics Metrics)
```
MetricsStoreAppend Append data to store

#### func  MetricsStoreEmpty

```go
func MetricsStoreEmpty()
```
MetricsStoreEmpty empties stored data

#### func  Refresh

```go
func Refresh()
```
Refresh Helper function used for testing

#### func  StatusCounterIncrement

```go
func StatusCounterIncrement(key string)
```
StatusCounterIncrement increment the status counter for a particular key

#### func  StatusCounterRead

```go
func StatusCounterRead(key string) int
```
StatusCounterRead the status counter for a particular key

#### func  TimestampMs

```go
func TimestampMs() int64
```
TimestampMs create a timestamp in milliseconds

#### type API

```go
type API struct {
	Name              string            `yaml:"name"`
	EventType         string            `yaml:"event_type"`  // override eventType
	Entity            string            `yaml:"entity"`      // define a custom entity name
	EntityType        string            `yaml:"entity_type"` // define a custom entity type (namespace)
	Ingest            bool              `yaml:"ingest"`
	Inventory         map[string]string `yaml:"inventory"`      // set as inventory
	InventoryOnly     bool              `yaml:"inventory_only"` // only generate inventory data
	Events            map[string]string `yaml:"events"`         // set as events
	EventsOnly        bool              `yaml:"events_only"`    // only generate events
	Merge             string            `yaml:"merge"`          // merge into another eventType
	Prefix            string            `yaml:"prefix"`         // prefix attribute keys
	File              string            `yaml:"file"`
	URL               string            `yaml:"url"`
	Pagination        Pagination        `yaml:"pagination"`
	EscapeURL         bool              `yaml:"escape_url"`
	Prometheus        Prometheus        `yaml:"prometheus"`
	Cache             string            `yaml:"cache"` // read data from datastore
	Database          string            `yaml:"database"`
	DbDriver          string            `yaml:"db_driver"`
	DbConn            string            `yaml:"db_conn"`
	Shell             string            `yaml:"shell"`
	CommandsAsync     bool              `yaml:"commands_async"` // run commands async
	Commands          []Command         `yaml:"commands"`
	DbQueries         []Command         `yaml:"db_queries"`
	DbAsync           bool              `yaml:"db_async"` // perform db queries async
	Jmx               JMX               `yaml:"jmx"`
	IgnoreLines       []int             // not implemented - idea is to ignore particular lines starting from 0 of the command output
	User, Pass        string
	Proxy             string
	TLSConfig         TLSConfig `yaml:"tls_config"`
	Timeout           int
	Method            string
	Payload           string
	Headers           map[string]string `yaml:"headers"`
	DisableParentAttr bool              `yaml:"disable_parent_attr"`
	StartKey          []string          `yaml:"start_key"` // start from a different section of the payload
	StoreLookups      map[string]string `yaml:"store_lookups"`
	StoreVariables    map[string]string `yaml:"store_variables"`
	LazyFlatten       []string          `yaml:"lazy_flatten"`
	SampleKeys        map[string]string `yaml:"sample_keys"`
	RenameSamples     map[string]string `yaml:"rename_samples"`     // using regex if sample has a key that matches, make that a different sample
	SkipProcessing    []string          `yaml:"skip_processing"`    // skip processing particular keys using an array of regex strings
	InheritAttributes bool              `yaml:"inherit_attributes"` // attempts to inherit attributes were possible
	CustomAttributes  map[string]string `yaml:"custom_attributes"`  // set additional custom attributes
	SplitObjects      bool              `yaml:"split_objects"`      // convert object with nested objects to array

	// Key manipulation
	ToLower      bool              `yaml:"to_lower"`       // convert all unicode letters mapped to their lower case.
	ConvertSpace string            `yaml:"convert_space"`  // convert spaces to another char
	SnakeToCamel bool              `yaml:"snake_to_camel"` // snake_case to camelCase
	ReplaceKeys  map[string]string `yaml:"replace_keys"`   // uses rename_keys functionality
	RenameKeys   map[string]string `yaml:"rename_keys"`    // use regex to find keys, then replace value
	AddAttribute map[string]string `yaml:"add_attribute"`  // add attribute // needs description

	// Value manipulation
	PercToDecimal    bool              `yaml:"perc_to_decimal"` // will check strings, and perform a trimRight for the %
	PluckNumbers     bool              `yaml:"pluck_numbers"`   // plucks numbers out of the value
	Math             map[string]string `yaml:"math"`            // perform match across processed metrics
	SubParse         []Parse           `yaml:"sub_parse"`
	ValueParser      map[string]string `yaml:"value_parser"`      // find keys with regex, and parse the value with regex
	ValueTransformer map[string]string `yaml:"value_transformer"` // find key(s) with regex, and modify the value
	MetricParser     MetricParser      `yaml:"metric_parser"`     // to use the MetricParser for setting deltas and gauges a namespace needs to be set

	// Command based options
	Split     string   `yaml:"split"`      // default vertical, can be set to horizontal (column) useful for tabular outputs
	SplitBy   string   `yaml:"split_by"`   // character to split by
	SetHeader []string `yaml:"set_header"` // manually set header column names
	Regex     bool     `yaml:"regex"`      // process SplitBy as regex
	RowHeader int      `yaml:"row_header"` // set the row header, to be used with SplitBy
	RowStart  int      `yaml:"row_start"`  // start from this line, to be used with SplitBy

	// Filtering Options
	EventFilter  []Filter            `yaml:"event_filter"` // filters events in/out
	KeyFilter    []Filter            `yaml:"key_filter"`   // filters keys in/out
	StripKeys    []string            `yaml:"strip_keys"`
	RemoveKeys   []string            `yaml:"remove_keys"`
	KeepKeys     []string            `yaml:"keep_keys"`     // inverse of removing keys
	SampleFilter []map[string]string `yaml:"sample_filter"` // sample filter key pair values with regex

	// Debug Options
	Debug   bool `yaml:"debug"` // logs out additional data, should not be enabled for production use!
	Logging struct {
		Open bool `yaml:"open"` // log open related errors
	}
}
```

API YAML Struct

#### type ArgumentList

```go
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	ForceLogEvent         bool   `default:"false" help:"Force create an event for everything - useful for testing"`
	OverrideIPMode        string `default:"" help:"Force override ipMode used for container discovery set as private or public - useful for testing"`
	Local                 bool   `default:"true" help:"Collect local entity info"`
	ConfigPath            string `default:"" help:"Set a specific config file."`
	ConfigFile            string `default:"" help:"(deprecated) Set a specific config file. Alias for config_path"`
	ConfigDir             string `default:"flexConfigs/" help:"Set directory of config files"`
	ContainerDiscoveryDir string `default:"flexContainerDiscovery/" help:"Set directory of auto discovery config files"`
	ContainerDiscovery    bool   `default:"false" help:"Enable container auto discovery"`
	Fargate               bool   `default:"false" help:"Enable Fargate discovery"`
	DockerAPIVersion      string `default:"" help:"Force Docker client API version"`
	EventLimit            int    `default:"500" help:"Event limiter - max amount of events per execution"`
	Entity                string `default:"" help:"Manually set a remote entity name"`
	InsightsURL           string `default:"" help:"Set Insights URL"`
	InsightsAPIKey        string `default:"" help:"Set Insights API key"`
	InsightsOutput        bool   `default:"false" help:"Output the events generated to standard out"`
	MetricAPIUrl          string `default:"https://metric-api.newrelic.com/metric/v1" help:"Set Metric API URL"`
	MetricAPIKey          string `default:"" help:"Set Metric API key"`
	GitFlexDir            string `default:"flexGitConfigs/" help:"Set directory to store configs from git repository"`
	GitService            string `default:"github" help:"Set git service"`
	GitToken              string `default:"" help:"Set git token"`
	GitUser               string `default:"" help:"Set git user"`
	GitRepo               string `default:"" help:"Set git repository to sync"`
	GitURL                string `default:"" help:"Set alternate git url"`
	GitBranch             string `default:"master" help:"Checkout to specified git branch"`
	GitCommit             string `default:"" help:"Checkout to specified git commit, if set will not use branch"`
	ProcessConfigsSync    bool   `default:"false" help:"Process configs synchronously rather then async"`
}
```

ArgumentList Available Arguments

```go
var Args ArgumentList
```
Args Infrastructure SDK Arguments List

#### type BlkioStatEntry

```go
type BlkioStatEntry struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Op    string `json:"op"`
	Value uint64 `json:"value"`
}
```

BlkioStatEntry is one small entity to store a piece of Blkio stats Not used on
Windows.

#### type BlkioStats

```go
type BlkioStats struct {
	// number of bytes transferred to and from the block device
	IoServiceBytesRecursive []BlkioStatEntry `json:"io_service_bytes_recursive"`
	IoServicedRecursive     []BlkioStatEntry `json:"io_serviced_recursive"`
	IoQueuedRecursive       []BlkioStatEntry `json:"io_queue_recursive"`
	IoServiceTimeRecursive  []BlkioStatEntry `json:"io_service_time_recursive"`
	IoWaitTimeRecursive     []BlkioStatEntry `json:"io_wait_time_recursive"`
	IoMergedRecursive       []BlkioStatEntry `json:"io_merged_recursive"`
	IoTimeRecursive         []BlkioStatEntry `json:"io_time_recursive"`
	SectorsRecursive        []BlkioStatEntry `json:"sectors_recursive"`
}
```

BlkioStats stores All IO service stats for data read and write. This is a Linux
specific structure as the differences between expressing block I/O on Windows
and Linux are sufficiently significant to make little sense attempting to morph
into a combined structure.

#### type CPUStats

```go
type CPUStats struct {
	// CPU Usage. Linux and Windows.
	CPUUsage CPUUsage `json:"cpu_usage"`

	// System Usage. Linux only.
	SystemUsage uint64 `json:"system_cpu_usage,omitempty"`

	// Online CPUs. Linux only.
	OnlineCPUs uint32 `json:"online_cpus,omitempty"`

	// Throttling Data. Linux only.
	ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
}
```

CPUStats aggregates and wraps all CPU related info of container

#### type CPUUsage

```go
type CPUUsage struct {
	// Total CPU time consumed.
	// Units: nanoseconds (Linux)
	// Units: 100's of nanoseconds (Windows)
	TotalUsage uint64 `json:"total_usage"`

	// Total CPU time consumed per core (Linux). Not used on Windows.
	// Units: nanoseconds.
	PercpuUsage []uint64 `json:"percpu_usage,omitempty"`

	// Time spent by tasks of the cgroup in kernel mode (Linux).
	// Time spent by all container processes in kernel mode (Windows).
	// Units: nanoseconds (Linux).
	// Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers.
	UsageInKernelmode uint64 `json:"usage_in_kernelmode"`

	// Time spent by tasks of the cgroup in user mode (Linux).
	// Time spent by all container processes in user mode (Windows).
	// Units: nanoseconds (Linux).
	// Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers
	UsageInUsermode uint64 `json:"usage_in_usermode"`
}
```

CPUUsage stores All CPU stats aggregated since container inception.

#### type Command

```go
type Command struct {
	Name             string            `yaml:"name"`              // required for database use
	EventType        string            `yaml:"event_type"`        // override eventType (currently used for db only)
	Shell            string            `yaml:"shell"`             // command shell
	Cache            string            `yaml:"cache"`             // use content from cache instead of a run command
	Run              string            `yaml:"run"`               // runs commands, but if database is set, then this is used to run queries
	ContainerExec    string            `yaml:"container_exec"`    // execute a command against a container
	Jmx              JMX               `yaml:"jmx"`               // if wanting to run different jmx endpoints to merge
	CompressBean     bool              `yaml:"compress_bean"`     // compress bean name //unused
	IgnoreOutput     bool              `yaml:"ignore_output"`     // can be useful for chaining commands together
	MetricParser     MetricParser      `yaml:"metric_parser"`     // not used yet
	CustomAttributes map[string]string `yaml:"custom_attributes"` // set additional custom attributes
	Output           string            `yaml:"output"`            // jmx, raw, json
	LineEnd          int               `yaml:"line_end"`          // stop processing command output after a certain amount of lines
	LineStart        int               `yaml:"line_start"`        // start from this line
	Timeout          int               `yaml:"timeout"`           // command timeout
	Dial             string            `yaml:"dial"`              // eg. google.com:80
	Network          string            `yaml:"network"`           // default tcp

	// Parsing Options - Body
	Split       string `yaml:"split"`        // default vertical, can be set to horizontal (column) useful for outputs that look like a table
	SplitBy     string `yaml:"split_by"`     // character/match to split by
	SplitOutput string `yaml:"split_output"` // split output by found regex
	RegexMatch  bool   `yaml:"regex_match"`  // process SplitBy as a regex match
	GroupBy     string `yaml:"group_by"`     // group by character
	RowHeader   int    `yaml:"row_header"`   // set the row header, to be used with SplitBy
	RowStart    int    `yaml:"row_start"`    // start from this line, to be used with SplitBy

	// Parsing Options - Header
	SetHeader        []string `yaml:"set_header"`         // manually set header column names (used when split is is set to horizontal)
	HeaderSplitBy    string   `yaml:"header_split_by"`    // character/match to split header by
	HeaderRegexMatch bool     `yaml:"header_regex_match"` // process HeaderSplitBy as a regex match

	// RegexMatches
	RegexMatches []RegMatch `yaml:"regex_matches"`
}
```

Command Struct

#### type Config

```go
type Config struct {
	FileName           string             `yaml:"file_name"`           // set when file is read
	FilePath           string             `yaml:"file_path"`           // set when file is read
	ContainerDiscovery ContainerDiscovery `yaml:"container_discovery"` // provide container discovery parameter at config level
	Name               string
	Global             Global
	APIs               []API
	Datastore          map[string][]interface{} `yaml:"datastore"`
	LookupStore        map[string][]string      `yaml:"lookup_store"`
	LookupFile         string                   `yaml:"lookup_file"`
	VariableStore      map[string]string        `yaml:"variable_store"`
	Secrets            map[string]Secret        `yaml:"secrets"`
	CustomAttributes   map[string]string        `yaml:"custom_attributes"` // set additional custom attributes
	MetricAPI          bool                     `yaml:"metric_api"`        // enable use of the dimensional data models metric api
}
```

Config YAML Struct

#### type Container

```go
type Container struct {

	// The Docker ID for the container.
	DockerID string `json:"DockerID"`

	// The name of the container as specified in the task definition.
	Name string `json:"Name"`

	// The name of the container supplied to Docker.
	// The Amazon ECS container agent generates a unique name for the container to avoid name collisions when multiple copies of the same task definition are run on a single instance.
	DockerName string `json:"DockerName"`

	// The image for the container.
	Image string `json:"Image"`

	//The SHA-256 digest for the image.
	ImageID string `json:"ImageID,omitempty"`

	// Any ports exposed for the container. This parameter is omitted if there are no exposed ports.
	Ports string `json:"Ports"`

	// Any labels applied to the container. This parameter is omitted if there are no labels applied.
	Labels map[string]string `json:"Labels"`

	// Any labels applied to the container. This parameter is omitted if there are no labels applied.
	Limits map[string]uint64 `json:"Limits"`

	// The desired status for the container from Amazon ECS.
	DesiredStatus string `json:"DesiredStatus"`

	// The known status for the container from Amazon ECS.
	KnownStatus string `json:"KnownStatus"`

	// The exit code for the container. This parameter is omitted if the container has not exited.
	ExitCode string `json:"ExitCode"`

	// The time stamp for when the container was created. This parameter is omitted if the container has not been created yet.
	CreatedAt string `json:"CreatedAt"`

	// The time stamp for when the container started. This parameter is omitted if the container has not started yet.
	StartedAt string `json:"StartedAt"` // 2017-11-17T17:14:07.781711848Z

	// The time stamp for when the container stopped. This parameter is omitted if the container has not stopped yet.
	FinishedAt string `json:"FinishedAt"`

	// The type of the container. Containers that are specified in your task definition are of type NORMAL.
	// You can ignore other container types, which are used for internal task resource provisioning by the Amazon ECS container agent.
	Type string `json:"Type"`

	// The network information for the container, such as the network mode and IP address. This parameter is omitted if no network information is defined.
	Networks []Network `json:"Networks"`
}
```

Container as defined by the ECS metadata API

#### type ContainerDiscovery

```go
type ContainerDiscovery struct {
	Target          string `yaml:"target"`  // string of container or image to target
	Type            string `yaml:"type"`    // container or image
	Mode            string `yaml:"mode"`    // contains, prefix, exact
	Port            int    `yaml:"port"`    // port
	IPMode          string `yaml:"ip_mode"` // public / private
	FileName        string `yaml:"file_name"`
	ReplaceComplete bool   `yaml:"replace_complete"`
}
```

ContainerDiscovery struct

#### type Filter

```go
type Filter struct {
	Key     string `yaml:"key"`
	Value   string `yaml:"value"`
	Mode    string `yaml:"mode"`    // default regex, other options contains, prefix, suffix
	Inverse bool   `yaml:"inverse"` // inverse only works when being used for keys currently (setting to true is like using keep keys)
}
```

Filter struct

#### type Global

```go
type Global struct {
	BaseURL    string `yaml:"base_url"`
	User, Pass string
	Proxy      string
	Timeout    int
	Headers    map[string]string `yaml:"headers"`
	Jmx        JMX               `yaml:"jmx"`
	TLSConfig  TLSConfig         `yaml:"tls_config"`
}
```

Global struct

#### type JMX

```go
type JMX struct {
	Domain         string `yaml:"domain"`
	User           string `yaml:"user"`
	Pass           string `yaml:"pass"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	URIPath        string `yaml:"uri_path"`
	KeyStore       string `yaml:"key_store"`
	KeyStorePass   string `yaml:"key_store_pass"`
	TrustStore     string `yaml:"trust_store"`
	TrustStorePass string `yaml:"trust_store_pass"`
}
```

JMX struct

#### type MemoryStats

```go
type MemoryStats struct {

	// current res_counter usage for memory
	Usage uint64 `json:"usage,omitempty"`
	// maximum usage ever recorded.
	MaxUsage uint64 `json:"max_usage,omitempty"`
	// TODO(vishh): Export these as stronger types.
	// all the stats exported via memory.stat.
	Stats map[string]uint64 `json:"stats,omitempty"`
	// number of times memory usage hits limits.
	Failcnt uint64 `json:"failcnt,omitempty"`
	Limit   uint64 `json:"limit,omitempty"`

	// committed bytes
	Commit uint64 `json:"commitbytes,omitempty"`
	// peak committed bytes
	CommitPeak uint64 `json:"commitpeakbytes,omitempty"`
	// private working set
	PrivateWorkingSet uint64 `json:"privateworkingset,omitempty"`
}
```

MemoryStats aggregates all memory stats since container inception on Linux.
Windows returns stats for commit and private working set only.

#### type MetricParser

```go
type MetricParser struct {
	Namespace Namespace                         `yaml:"namespace"`
	Metrics   map[string]string                 `yaml:"metrics"`  // inputBytesPerSecond: RATE
	Mode      string                            `yaml:"mode"`     // options regex, prefix, suffix, contains
	AutoSet   bool                              `yaml:"auto_set"` // if set to true, will attempt to do a contains instead of a direct key match, this is useful for setting multiple metrics
	Counts    map[string]int64                  `yaml:"counts"`
	Summaries map[string]map[string]interface{} `yaml:"summaries"`
}
```

MetricParser Struct

#### type Metrics

```go
type Metrics struct {
	TimestampMs      int64                    `json:"timestamp.ms,omitempty"` // required for every metric at root or nested
	IntervalMs       int64                    `json:"interval.ms,omitempty"`  // required for count & summary
	CommonAttributes map[string]interface{}   `json:"commonAttributes,omitempty"`
	Metrics          []map[string]interface{} `json:"metrics"` // summaries have a different value structure then gauges or counters
}
```

Metrics struct

#### type Namespace

```go
type Namespace struct {
	// if neither of the below are set and the MetricParser is used, the namespace will default to the "Name" attribute
	CustomAttr   string   `yaml:"custom_attr"`   // set your own custom namespace attribute
	ExistingAttr []string `yaml:"existing_attr"` // utilise existing attributes and chain together to create a custom namespace
}
```

Namespace Struct

#### type NetStats

```go
type NetStats struct {
	RxBytes   uint64 `json:"rx_bytes"`
	RxPackets uint64 `json:"rx_packets"`
	TxBytes   uint64 `json:"tx_bytes"`
	TxPackets uint64 `json:"tx_packets"`
}
```

NetStats ECS container network usage

#### type Network

```go
type Network struct {

	// NetworkMode currently only supported mode is awsvpc
	NetworkMode string `json:"NetworkMode"`

	// IPv4 Addresses supplied in a single element list
	IPv4Addresses []string `json:"IPv4Addresses"`
}
```

Network information of the container

#### type NetworkStats

```go
type NetworkStats struct {
	// Bytes received. Windows and Linux.
	RxBytes uint64 `json:"rx_bytes"`
	// Packets received. Windows and Linux.
	RxPackets uint64 `json:"rx_packets"`
	// Received errors. Not used on Windows. Note that we don't `omitempty` this
	// field as it is expected in the >=v1.21 API stats structure.
	RxErrors uint64 `json:"rx_errors"`
	// Incoming packets dropped. Windows and Linux.
	RxDropped uint64 `json:"rx_dropped"`
	// Bytes sent. Windows and Linux.
	TxBytes uint64 `json:"tx_bytes"`
	// Packets sent. Windows and Linux.
	TxPackets uint64 `json:"tx_packets"`
	// Sent errors. Not used on Windows. Note that we don't `omitempty` this
	// field as it is expected in the >=v1.21 API stats structure.
	TxErrors uint64 `json:"tx_errors"`
	// Outgoing packets dropped. Windows and Linux.
	TxDropped uint64 `json:"tx_dropped"`
	// Endpoint ID. Not used on Linux.
	EndpointID string `json:"endpoint_id,omitempty"`
	// Instance ID. Not used on Linux.
	InstanceID string `json:"instance_id,omitempty"`
}
```

NetworkStats aggregates the network stats of one container

#### type Pagination

```go
type Pagination struct {
	// internal attribute use
	OriginalURL  string `yaml:"original_url"`  // internal use (not intended for user use)
	NextLink     string `yaml:"next_link"`     // internal use (not intended for user use)
	NoPages      int    `yaml:"no_pages"`      // used to track how many pages walked (not intended for user use)
	PageMarker   int    `yaml:"page_marker"`   // used as a page marker (not intended for user use)
	CursorMarker string `yaml:"cursor_marker"` // used as a marker currently for cursors (not intended for user use)

	// attributes for use
	Increment   int    `yaml:"increment"`     // number to increment by
	MaxPages    int    `yaml:"max_pages"`     // set the max number of pages to walk (needs to be set or payload_key)
	MaxPagesKey string `yaml:"max_pages_key"` // set the max number of pages to walk (needs to be set or payload_key)

	PageStart    int    `yaml:"page_start"`     // page to start walking from
	PageNextKey  string `yaml:"page_next_key"`  // set a key to look for the next page to walk too - regex eg. "next":.(\d+)
	PageLimit    int    `yaml:"page_limit"`     // manually set the page_limit to use
	PageLimitKey string `yaml:"page_limit_key"` // set a key to look for the limit / page size / offset to use - regex eg. "limit":.(\d+)

	PayloadKey string `yaml:"payload_key"` // set a key to watch if data exists at a particular attribute (needs to be set or max_pages) regex eg. "someKey":(\[(.*?)\]|\{(.*?)\})

	NextCursorKey string `yaml:"next_cursor_key"` // watch for next cursor to query next
	MaxCursorKey  string `yaml:"max_cursor_key"`  // watch for max cursor to stop at

	NextLinkKey string `yaml:"next_link_key"` // look for a next link key to browse too
}
```

Pagination handles request pagination

#### type Parse

```go
type Parse struct {
	Type    string   `yaml:"type"` // perform a contains, match, hasPrefix or regex for specified key
	Key     string   `yaml:"key"`
	SplitBy []string `yaml:"split_by"`
}
```

Parse struct

#### type PidsStats

```go
type PidsStats struct {
	// Current is the number of pids in the cgroup
	Current uint64 `json:"current,omitempty"`
	// Limit is the hard limit on the number of pids in the cgroup.
	// A "Limit" of 0 means that there is no limit.
	Limit uint64 `json:"limit,omitempty"`
}
```

PidsStats contains the stats of a container's pids

#### type Prometheus

```go
type Prometheus struct {
	Enable           bool              `yaml:"enable"`
	Raw              bool              `yaml:"raw"`             // creates an event per prometheus metric retaining all metadata
	Unflatten        bool              `yaml:"unflatten"`       // unflattens all counters and gauges into separate metric samples retaining all their metadata // make this map[string]string
	FlattenedEvent   string            `yaml:"flattened_event"` // name of the flattenedEvent
	KeyMerge         []string          `yaml:"key_merge"`       // list of keys to merge into the key name when flattening, not usable when unflatten set to true
	KeepLabels       bool              `yaml:"keep_labels"`     // not usable when unflatten set to true
	KeepHelp         bool              `yaml:"keep_help"`       // not usable when unflatten set to true
	CustomAttributes map[string]string `yaml:"custom_attributes"`
	SampleKeys       map[string]string `yaml:"sample_keys"`
	Histogram        bool              `yaml:"histogram"`       // if flattening by default, create a full histogram sample
	HistogramEvent   string            `yaml:"histogram_event"` // override histogram event type
	Summary          bool              `yaml:"summary"`         // if flattening by default, create a full summary sample
	SummaryEvent     string            `yaml:"summaryevent"`    // override summary event type
	GoMetrics        bool              `yaml:"go_metrics"`      // enable go metrics

}
```

Prometheus struct

#### type RegMatch

```go
type RegMatch struct {
	Expression string   `yaml:"expression"`
	Keys       []string `yaml:"keys"`
	KeysMulti  []string `yaml:"keys_multi"`
}
```

RegMatch support for regex matches

#### type SampleMerge

```go
type SampleMerge struct {
	EventType string   `yaml:"event_type"` // new event_type name for the sample
	Samples   []string `yaml:"samples"`    // list of samples to be merged
}
```

SampleMerge merge multiple samples into one (will remove previous samples)

#### type Secret

```go
type Secret struct {
	Kind           string                 `yaml:"kind"` // eg. aws, vault
	Key            string                 `yaml:"key"`
	Token          string                 `yaml:"token"`
	CredentialFile string                 `yaml:"credential_file"`
	ConfigFile     string                 `yaml:"config_file"`
	File           string                 `yaml:"file"`
	Data           string                 `yaml:"data"`
	HTTP           API                    `yaml:"http"`
	Region         string                 `yaml:"region"`
	Base64Decode   bool                   `yaml:"base64_decode"`
	Type           string                 `yaml:"type"` // basic, equal, json
	Values         map[string]interface{} `yaml:"values"`
}
```

Secret Struct

#### type Stats

```go
type Stats struct {
	// Common stats
	Read    time.Time `json:"read"`
	PreRead time.Time `json:"preread"`

	// Linux specific stats, not populated on Windows.
	PidsStats  PidsStats  `json:"pids_stats,omitempty"`
	BlkioStats BlkioStats `json:"blkio_stats,omitempty"`

	// Windows specific stats, not populated on Linux.
	NumProcs     uint32       `json:"num_procs"`
	StorageStats StorageStats `json:"storage_stats,omitempty"`

	// Shared stats
	CPUStats    CPUStats    `json:"cpu_stats,omitempty"`
	PreCPUStats CPUStats    `json:"precpu_stats,omitempty"` // "Pre"="Previous"
	MemoryStats MemoryStats `json:"memory_stats,omitempty"`

	// Network from AWS Stats API
	Network NetStats `json:"network"`

	// Networks request version >=1.21
	Networks map[string]NetworkStats `json:"networks,omitempty"`
}
```

Stats is Ultimate struct aggregating all types of stats of one container

#### type StorageStats

```go
type StorageStats struct {
	ReadCountNormalized  uint64 `json:"read_count_normalized,omitempty"`
	ReadSizeBytes        uint64 `json:"read_size_bytes,omitempty"`
	WriteCountNormalized uint64 `json:"write_count_normalized,omitempty"`
	WriteSizeBytes       uint64 `json:"write_size_bytes,omitempty"`
}
```

StorageStats is the disk I/O stats for read/write on Windows.

#### type TLSConfig

```go
type TLSConfig struct {
	Enable             bool   `yaml:"enable"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	MinVersion         uint16 `yaml:"min_version"`
	MaxVersion         uint16 `yaml:"max_version"`
	Ca                 string `yaml:"ca"` // path to ca to read
}
```

TLSConfig struct

#### type TaskMetadata

```go
type TaskMetadata struct {

	// The name of the cluster that hosts the task.
	Cluster string `json:"Cluster"`

	// The Amazon Resource Name (ARN) of the task.
	TaskArn string `locationName:"taskArn" type:"string"`

	// The family of the Amazon ECS task definition for the task.
	Family string `locationName:"family" type:"string"`

	// The revision of the Amazon ECS task definition for the task.
	Revision string `locationName:"revision" type:"string"`

	// The desired status of the task. For more information, see Task Lifecycle
	// (https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_life_cycle.html).
	DesiredStatus string `locationName:"desiredStatus" type:"string"`

	// The known status for the task from Amazon ECS.
	KnownStatus string `locationName:"knownStatus" type:"string"`

	// A list of container metadata for each container associated with the task.
	Containers []Container `json:"Containers"`

	// The resource limits specified at the task level (such as CPU and memory). This parameter is omitted if no resource limits are defined.
	Limits map[string]float64 `json:"Limits"`

	// The Unix timestamp for when the container image pull began.
	PullStartedAt *time.Time `locationName:"pullStartedAt" type:"timestamp"`

	// The Unix timestamp for when the container image pull completed.
	PullStoppedAt *time.Time `locationName:"pullStoppedAt" type:"timestamp"`

	// The Unix timestamp for when the task execution stopped.
	ExecutionStoppedAt *time.Time `locationName:"executionStoppedAt" type:"timestamp"`
}
```

TaskMetadata
https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v2.html#task-metadata-endpoint-v2-response

#### type ThrottlingData

```go
type ThrottlingData struct {
	// Number of periods with throttling active
	Periods uint64 `json:"periods"`
	// Number of periods when the container hits its throttling limit.
	ThrottledPeriods uint64 `json:"throttled_periods"`
	// Aggregate time the container was throttled for in nanoseconds.
	ThrottledTime uint64 `json:"throttled_time"`
}
```

ThrottlingData stores CPU throttling stats of one running container. Not used on
Windows.

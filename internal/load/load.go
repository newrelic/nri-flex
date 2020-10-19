/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package load

import (
	"sync"
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	logrus "github.com/sirupsen/logrus"
)

// ArgumentList Available Arguments
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	ForceLogEvent           bool   `default:"false" help:"Force create an event for everything - useful for testing"`
	OverrideIPMode          string `default:"" help:"Force override ipMode used for container discovery set as private or public - useful for testing"`
	Local                   bool   `default:"true" help:"Collect local entity info"`
	ConfigPath              string `default:"" help:"Set a specific config file."`
	ConfigFile              string `default:"" help:"(deprecated) Set a specific config file. Alias for config_path"`
	ConfigDir               string `default:"flexConfigs/" help:"Set directory of config files"`
	ContainerDiscoveryDir   string `default:"flexContainerDiscovery/" help:"Set directory of auto discovery config files"`
	ContainerDiscovery      bool   `default:"false" help:"Enable container auto discovery"`
	ContainerDiscoveryMulti bool   `default:"false" help:"Allow a container to be matched multiple times"`
	ContainerDump           bool   `default:"false" help:"Dump all containers, useful for debugging"`
	Fargate                 bool   `default:"false" help:"Enable Fargate discovery"`
	DockerAPIVersion        string `default:"" help:"Force Docker client API version"`
	EventLimit              int    `default:"0" help:"Event limiter - limit events per execution, 0 to disable"`
	Entity                  string `default:"" help:"Manually set a remote entity name"`
	InsightsURL             string `default:"" help:"Set Insights URL"`
	InsightsAPIKey          string `default:"" help:"Set Insights API key"`
	InsightsOutput          bool   `default:"false" help:"Output the events generated to standard out"`
	InsightBatchSize        int    `default:"5000" help:"Batch Size - number of metrics per post call to Insight endpoint"`
	MetricAPIUrl            string `default:"https://metric-api.newrelic.com/metric/v1" help:"Set Metric API URL"`
	MetricAPIKey            string `default:"" help:"Set Metric API key"`
	GitFlexDir              string `default:"flexGitConfigs/" help:"Set directory to store configs from git repository"`
	GitService              string `default:"github" help:"Set git service"`
	GitToken                string `default:"" help:"Set git token"`
	GitUser                 string `default:"" help:"Set git user"`
	GitRepo                 string `default:"" help:"Set git repository to sync"`
	GitURL                  string `default:"" help:"Set alternate git url"`
	GitBranch               string `default:"master" help:"Checkout to specified git branch"`
	GitCommit               string `default:"" help:"Checkout to specified git commit, if set will not use branch"`
	ProcessConfigsSync      bool   `default:"false" help:"Process configs synchronously rather then async"`
	// ProcessDiscovery      bool   `default:"true" help:"Enable process discovery"`
	EncryptPass          string `default:"" help:"Pass to be encypted"`
	PassPhrase           string `default:"N3wR3lic!" help:"PassPhrase used to de/encrypt"`
	DiscoverProcessWin   bool   `default:"false" help:"Discover Process info on Windows OS"`
	DiscoverProcessLinux bool   `default:"false" help:"Discover Process info on Linux OS"`
	NRJMXToolPath        string `default:"/usr/lib/nrjmx/" help:"Set a custom path for nrjmx tool"`
	StructuredLogs       bool   `default:"false" help:"output logs in Json structure format for external tool parsing"`
}

// Args Infrastructure SDK Arguments List
var Args ArgumentList

// StartTime time Flex starts in Nanoseconds
var StartTime int64

// Integration Infrastructure SDK Integration
var Integration *integration.Integration

// IgnoredIntegrationData this is used for lookups with ignored output
var IgnoredIntegrationData []map[string]interface{}

// Entity Infrastructure SDK Entity
var Entity *integration.Entity

// Hostname current host
var Hostname string

// ContainerID current container id
var ContainerID string

// IsKubernetes basic check if running on k8s
var IsKubernetes bool

// IsFargate basic check if running on fargate
var IsFargate bool

// LambdaName if running on lambda add name from AWS_LAMBDA_FUNCTION_NAME
var LambdaName string

// AWSExecutionEnv AWS execution environment
var AWSExecutionEnv string

// DiscoveredProcesses discovered processes
var DiscoveredProcesses map[string]string

// IngestData store ingested data
var IngestData interface{}

// Logrus create instance of the logger
var Logrus = logrus.New()

var IntegrationName = "com.newrelic.nri-flex" // IntegrationName Name
var IntegrationNameShort = "nri-flex"         // IntegrationNameShort Short Name
var IntegrationVersion = "Unknown-SNAPSHOT"   // IntegrationVersion Version

const (
	DefaultSplitBy     = ":"                      // unused currently
	DefaultTimeout     = 10000 * time.Millisecond // 10 seconds, used for raw commands
	DefaultDialTimeout = 1000                     // 1 seconds, used for dial
	DefaultPingTimeout = 5000                     // 5 seconds
	DefaultHANA        = "hdb"
	DefaultDB2         = "go_ibm_db"
	DefaultPostgres    = "postgres"
	DefaultMSSQLServer = "sqlserver"
	DefaultMySQL       = "mysql"
	DefaultOracle      = "godror"
	DefaultSybase      = "ase"
	DefaultVertica     = "vertica"
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
	TypeCname          = "cname"
	TypeJSON           = "json"
	TypeXML            = "xml"
	TypeColumns        = "columns"
	Contains           = "contains"
)

// MetricsStore for Dimensional Metrics to store data and lock and unlock when needed
var MetricsStore = struct {
	sync.RWMutex
	Data []Metrics
}{}

// MetricsStoreAppend Append data to store
func MetricsStoreAppend(metrics Metrics) {
	MetricsStore.Lock()
	MetricsStore.Data = append(MetricsStore.Data, metrics)
	MetricsStore.Unlock()
}

// MetricsStoreEmpty empties stored data
func MetricsStoreEmpty() {
	MetricsStore.Lock()
	MetricsStore.Data = []Metrics{}
	MetricsStore.Unlock()
}

// Metrics struct
type Metrics struct {
	TimestampMs      int64                    `json:"timestamp.ms,omitempty"` // required for every metric at root or nested
	IntervalMs       int64                    `json:"interval.ms,omitempty"`  // required for count & summary
	CommonAttributes map[string]interface{}   `json:"commonAttributes,omitempty"`
	Metrics          []map[string]interface{} `json:"metrics"` // summaries have a different value structure then gauges or counters
}

// AgentConfig stores the information from a single V4 integrations file
// This has been added so that Flex can understand the V4 agent format when users are using the config_file parameter
type AgentConfig struct {
	Integrations []ConfigEntry `yaml:"integrations"`
}

// ConfigEntry holds an integrations YAML configuration entry. It may define multiple types of tasks
type ConfigEntry struct {
	Name   string `yaml:"name"`
	Config Config `yaml:"config"`
}

// Config YAML Struct
type Config struct {
	FileName           string             `yaml:"file_name"`           // set when file is read
	FilePath           string             `yaml:"file_path"`           // set when file is read
	ContainerDiscovery ContainerDiscovery `yaml:"container_discovery"` // provide container discovery parameter at config level
	Name               string
	Global             Global
	APIs               []API
	Datastore          map[string][]interface{}       `yaml:"datastore"`
	LookupStore        map[string]map[string]struct{} `yaml:"lookup_store"` // ensures uniqueness vs a slice
	LookupFile         string                         `yaml:"lookup_file"`
	VariableStore      map[string]string              `yaml:"variable_store"`
	Secrets            map[string]Secret              `yaml:"secrets"`
	CustomAttributes   map[string]string              `yaml:"custom_attributes"` // set additional custom attributes
	MetricAPI          bool                           `yaml:"metric_api"`        // enable use of the dimensional data models metric api
}

// Secret Struct
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

// ContainerDiscovery struct
type ContainerDiscovery struct {
	Target          string `yaml:"target"`  // string of container or image to target
	Type            string `yaml:"type"`    // container or image
	Mode            string `yaml:"mode"`    // contains, prefix, exact
	Port            int    `yaml:"port"`    // port
	IPMode          string `yaml:"ip_mode"` // public / private
	FileName        string `yaml:"file_name"`
	ReplaceComplete bool   `yaml:"replace_complete"`
}

// Global struct
type Global struct {
	BaseURL    string `yaml:"base_url"`
	User, Pass string
	Proxy      string
	Timeout    int
	Headers    map[string]string `yaml:"headers"`
	Jmx        JMX               `yaml:"jmx"`
	TLSConfig  TLSConfig         `yaml:"tls_config"`
	Passphrase string            `yaml:"pass_phrase"`
	SSHPEMFile string            `yaml:"ssh_pem_file"`
}

// TLSConfig struct
type TLSConfig struct {
	Enable             bool   `yaml:"enable"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	MinVersion         uint16 `yaml:"min_version"`
	MaxVersion         uint16 `yaml:"max_version"`
	Ca                 string `yaml:"ca"` // path to ca to read
}

// SampleMerge merge multiple samples into one (will remove previous samples)
type SampleMerge struct {
	EventType string   `yaml:"event_type"` // new event_type name for the sample
	Samples   []string `yaml:"samples"`    // list of samples to be merged
}

// API YAML Struct
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
	RunAsync          bool              `yaml:"run_async" `     // API block to run in Async mode when using with lookupstore
	JoinKey           string            `yaml:"join_key"`       // merge into another eventType
	Prefix            string            `yaml:"prefix"`         // prefix attribute keys
	File              string            `yaml:"file"`
	URL               string            `yaml:"url"`
	Pagination        Pagination        `yaml:"pagination"`
	EscapeURL         bool              `yaml:"escape_url"`
	Prometheus        Prometheus        `yaml:"prometheus"`
	Cache             string            `yaml:"cache"` // read data from datastore
	Database          string            `yaml:"database"`
	DBDriver          string            `yaml:"db_driver"`
	DBConn            string            `yaml:"db_conn"`
	Shell             string            `yaml:"shell"`
	CommandsAsync     bool              `yaml:"commands_async"` // run commands async
	Commands          []Command         `yaml:"commands"`
	DBQueries         []Command         `yaml:"db_queries"`
	DBAsync           bool              `yaml:"db_async"`   // perform db queries async
	Jq                string            `yaml:"jq"`         // parse data using jq
	ParseHTML         bool              `yaml:"parse_html"` // parse text/html content type table element to JSON
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
	DedupeLookups     []string          `yaml:"dedupe_lookups"`
	StoreVariables    map[string]string `yaml:"store_variables"`
	LazyFlatten       []string          `yaml:"lazy_flatten"`
	SampleKeys        map[string]string `yaml:"sample_keys"`
	RenameSamples     map[string]string `yaml:"rename_samples"`     // using regex if sample has a key that matches, make that a different sample
	SkipProcessing    []string          `yaml:"skip_processing"`    // skip processing particular keys using an array of regex strings
	InheritAttributes bool              `yaml:"inherit_attributes"` // attempts to inherit attributes were possible
	CustomAttributes  map[string]string `yaml:"custom_attributes"`  // set additional custom attributes
	SplitObjects      bool              `yaml:"split_objects"`      // convert object with nested objects to array
	SplitArray        bool              `yaml:"split_array"`        // convert array to samples, use SetHeader to set attribute name
	LeafArray         bool              `yaml:"leaf_array"`         // convert array element to samples when SplitArray, use SetHeader to set attribute name
	Scp               SCP               `yaml:"scp"`
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

	ValueMapper         map[string][]string `yaml:"value_mapper"`         // Map the value of the key based on regex pattern,  "*.?\s(Service Status)=>$1-Good"
	TimestampConversion map[string]string   `yaml:"timestamp_conversion"` // find keys with regex, convert date<=>timestamp

	// Command based options
	Split     string   `yaml:"split"`      // default vertical, can be set to horizontal (column) useful for tabular outputs
	SplitBy   string   `yaml:"split_by"`   // character to split by
	SetHeader []string `yaml:"set_header"` // manually set header column names
	Regex     bool     `yaml:"regex"`      // process SplitBy as regex
	RowHeader int      `yaml:"row_header"` // set the row header, to be used with SplitBy
	RowStart  int      `yaml:"row_start"`  // start from this line, to be used with SplitBy

	// Filtering Options
	EventFilter                 []Filter            `yaml:"event_filter"` // filters events in/out
	KeyFilter                   []Filter            `yaml:"key_filter"`   // filters keys in/out
	StripKeys                   []string            `yaml:"strip_keys"`
	RemoveKeys                  []string            `yaml:"remove_keys"`
	KeepKeys                    []string            `yaml:"keep_keys"`                       // inverse of removing keys
	SampleFilter                []map[string]string `yaml:"sample_filter"`                   // exclude sample filter key pair values with regex === sample_exclude_filter
	SampleIncludeFilter         []map[string]string `yaml:"sample_include_filter"`           // include sample filter key pair values with regex
	SampleExcludeFilter         []map[string]string `yaml:"sample_exclude_filter"`           // exclude sample filter key pair values with regex
	SampleIncludeMatchAllFilter []map[string]string `yaml:"sample_include_match_all_filter"` //include samples where multiple keys match the specified
	IgnoreOutput                bool                `yaml:"ignore_output"`                   // ignore the output completely, useful when creating lookups

	SaveOutput string `yaml:"save_output"` // Save output samples to a file

	// Debug Options
	Debug   bool     `yaml:"debug"` // logs out additional data, should not be enabled for production use!
	Logging struct { // log to insights
		Open bool `yaml:"open"` // log open related errors
	}

	ReturnHeaders bool `yaml:"return_headers"`
}

// Filter struct
type Filter struct {
	Key     string `yaml:"key"`
	Value   string `yaml:"value"`
	Mode    string `yaml:"mode"`    // default regex, other options contains, prefix, suffix
	Inverse bool   `yaml:"inverse"` // inverse only works when being used for keys currently (setting to true is like using keep keys)
}

// Command Struct
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
	Output           string            `yaml:"output"`            // jmx, raw, json,xml
	LineEnd          int               `yaml:"line_end"`          // stop processing command output after a certain amount of lines
	LineStart        int               `yaml:"line_start"`        // start from this line
	Timeout          int               `yaml:"timeout"`           // command timeout
	Dial             string            `yaml:"dial"`              // eg. google.com:80
	Network          string            `yaml:"network"`           // default tcp
	OS               string            `yaml:"os"`                // default empty for any operating system, if set will check if the OS matches else will skip execution

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

// Pagination handles request pagination
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

	NextLinkKey  string `yaml:"next_link_key"`  // look for a next link key to browse too
	NextLinkHost string `yaml:"next_link_host"` // set next link host - useful when next_link_key returns a partial URL, e.g "/mynextlinkABC", the next link will be {next_link_host}/mynextlinkABC
}

// RegMatch support for regex matches
type RegMatch struct {
	Expression string   `yaml:"expression"`
	Keys       []string `yaml:"keys"`
	KeysMulti  []string `yaml:"keys_multi"`
}

// Prometheus struct
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
	// Metrics          []string          `yaml:"metrics"`         // filter metrics
}

// JMX struct
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

// SCP struct
type SCP struct {
	User       string `yaml:"user"`
	Pass       string `yaml:"pass"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	RemoteFile string `yaml:"remote_file"`
	Passphrase string `yaml:"pass_phrase"`
	SSHPEMFile string `yaml:"ssh_pem_file"`
}

// Parse struct
type Parse struct {
	Type    string   `yaml:"type"` // perform a contains, match, hasPrefix or regex for specified key
	Key     string   `yaml:"key"`
	SplitBy []string `yaml:"split_by"`
}

// MetricParser Struct
type MetricParser struct {
	Namespace Namespace                         `yaml:"namespace"`
	Metrics   map[string]string                 `yaml:"metrics"`  // inputBytesPerSecond: RATE
	Mode      string                            `yaml:"mode"`     // options regex, prefix, suffix, contains
	AutoSet   bool                              `yaml:"auto_set"` // if set to true, will attempt to do a contains instead of a direct key match, this is useful for setting multiple metrics
	Counts    map[string]int64                  `yaml:"counts"`
	Summaries map[string]map[string]interface{} `yaml:"summaries"`
}

// Namespace Struct
type Namespace struct {
	// if neither of the below are set and the MetricParser is used, the namespace will default to the "Name" attribute
	CustomAttr   string   `yaml:"custom_attr"`   // set your own custom namespace attribute
	ExistingAttr []string `yaml:"existing_attr"` // utilise existing attributes and chain together to create a custom namespace
}

// MakeTimestamp creates timestamp in milliseconds
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// SamplesToMerge keep merge sapmles
type SamplesToMerge struct {
	sync.RWMutex
	Data map[string][]interface{}
}

// SampleAppend append sample with locking
func (s *SamplesToMerge) SampleAppend(key string, sample interface{}) {
	s.Lock()
	defer s.Unlock()
	(s.Data)[key] = append((s.Data)[key], sample)
}

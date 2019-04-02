package load

import (
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

// ArgumentList Available Arguments
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	ForceLogEvent         bool   `default:"false" help:"Force create an event for everything - useful for testing"`
	OverrideIPMode        string `default:"" help:"Force override ipMode used for container discovery set as private or public - useful for testing"`
	Local                 bool   `default:"true" help:"Collect local entity info"`
	ConfigFile            string `default:"" help:"Set a specific config file - not usable for container discovery"`
	ConfigDir             string `default:"flexConfigs/" help:"Set directory of config files"`
	ContainerDiscoveryDir string `default:"flexContainerDiscovery/" help:"Set directory of auto discovery config files"`
	ContainerDiscovery    bool   `default:"false" help:"Enable container auto discovery"`
	DockerAPIVersion      string `default:"" help:"Force Docker client API version"`
	EventLimit            int    `default:"500" help:"Event limiter - max amount of events per execution"`
	Entity                string `default:"" help:"Manually set a remote entity name"`
	InsightsURL           string `default:"" help:"Set Insights URL"`
	InsightsAPIKey        string `default:"" help:"Set Insights API key"`
	InsightsInterval      int    `default:"0" help:"Run Insights mode periodically at this set interval"`
	InsightsOutput        bool   `default:"false" help:"Output the events generated to standard out"`

	// not implemented yet
	// ClusterModeKey string `default:"" help:"Set key used for cluster mode identification"`
	// ClusterModeExp string `default:"60s" help:"Set cluster mode key identifier expiration"`
}

// Args Infrastructure SDK Arguments List
var Args ArgumentList

// Integration Infrastructure SDK Integration
var Integration *integration.Integration

// Entity Infrastructure SDK Entity
var Entity *integration.Entity

// Hostname current host
var Hostname string

// ContainerID current container id
var ContainerID string

// EventDropCount current number of events dropped due to limiter
var EventDropCount = 0

// EventCount number of events processed
var EventCount = 0

// EventDistribution number of events distributed per sample
var EventDistribution = map[string]int{}

// ConfigsProcessed number of configs processed
var ConfigsProcessed = 0

const (
	IntegrationName      = "com.kav91.nri-flex"     // IntegrationName Name
	IntegrationNameShort = "nri-flex"               // IntegrationNameShort Short Name
	IntegrationVersion   = "0.5.0-pre"              // IntegrationVersion Version
	DefaultSplitBy       = ":"                      // unused currently
	DefaultTimeout       = 10000 * time.Millisecond // 10 seconds, used for raw commands
	DefaultPingTimeout   = 5000                     // 5 seconds
	DefaultPostgres      = "postgres"
	DefaultMSSQLServer   = "sqlserver"
	DefaultMySQL         = "mysql"
	DefaultOracle        = "ora"
	DefaultJmxPath       = "./nrjmx/"
	DefaultJmxHost       = "127.0.0.1"
	DefaultJmxPort       = "9999"
	DefaultJmxUser       = "admin"
	DefaultJmxPass       = "admin"
	DefaultIPMode        = "private"
	DefaultShell         = "/bin/sh"
	DefaultLineLimit     = 255
	Public               = "public"
	Private              = "private"
	Jmx                  = "jmx"
	Img                  = "img"
	TypeJSON             = "json"
	TypeColumns          = "columns"
)

// Config YAML Struct
type Config struct {
	FileName         string // this will be set when files are read
	Name             string
	Global           Global
	APIs             []API
	Datastore        map[string][]interface{} `yaml:"datastore"`
	LookupStore      map[string][]string      `yaml:"lookup_store"`
	VariableStore    map[string]string        `yaml:"variable_store"`
	CustomAttributes map[string]string        `yaml:"custom_attributes"` // set additional custom attributes
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
}

// TLSConfig struct
type TLSConfig struct {
	Enable             bool   `yaml:"enable"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	MinVersion         uint16 `yaml:"min_version"`
	MaxVersion         uint16 `yaml:"max_version"`
}

// SampleMerge merge multiple samples into one (will remove previous samples)
type SampleMerge struct {
	EventType string   `yaml:"event_type"` // new event_type name for the sample
	Samples   []string `yaml:"samples"`    // list of samples to be merged
}

// API YAML Struct
type API struct {
	EventType         string     `yaml:"event_type"` // override eventType
	Merge             string     `yaml:"merge"`      // merge into another eventType
	Prefix            string     `yaml:"prefix"`     // prefix attribute keys
	Name              string     `yaml:"name"`
	File              string     `yaml:"file"`
	URL               string     `yaml:"url"`
	Prometheus        Prometheus `yaml:"prometheus"`
	Cache             string     `yaml:"cache"` // read data from datastore
	Database          string     `yaml:"database"`
	DbDriver          string     `yaml:"db_driver"`
	DbConn            string     `yaml:"db_conn"`
	Shell             string     `yaml:"shell"`
	Commands          []Command  `yaml:"commands"`
	DbQueries         []Command  `yaml:"db_queries"`
	Jmx               JMX        `yaml:"jmx"`
	IgnoreLines       []int      // not implemented - idea is to ignore particular lines starting from 0 of the command output
	User, Pass        string
	Proxy             string
	TLSConfig         TLSConfig `yaml:"tls_config"`
	Timeout           int
	Method            string
	Payload           string
	DisableParentAttr bool                `yaml:"disable_parent_attr"`
	Headers           map[string]string   `yaml:"headers"`
	StartKey          []string            `yaml:"start_key"`
	StoreLookups      map[string]string   `yaml:"store_lookups"`
	StoreVariables    map[string]string   `yaml:"store_variables"`
	StripKeys         []string            `yaml:"strip_keys"`
	LazyFlatten       []string            `yaml:"lazy_flatten"`
	SampleKeys        map[string]string   `yaml:"sample_keys"`
	ReplaceKeys       map[string]string   `yaml:"replace_keys"`   // uses rename_keys functionality
	RenameKeys        map[string]string   `yaml:"rename_keys"`    // use regex to find keys, then replace value
	RenameSamples     map[string]string   `yaml:"rename_samples"` // using regex if sample has a key that matches, make that a different sample
	RemoveKeys        []string            `yaml:"remove_keys"`
	KeepKeys          []string            `yaml:"keep_keys"`     // inverse of removing keys
	ToLower           bool                `yaml:"to_lower"`      // convert all unicode letters mapped to their lower case.
	ConvertSpace      string              `yaml:"convert_space"` // convert spaces to another char
	SnakeToCamel      bool                `yaml:"snake_to_camel"`
	PercToDecimal     bool                `yaml:"perc_to_decimal"` // will check strings, and perform a trimRight for the %
	PluckNumbers      bool                `yaml:"pluck_numbers"`   // plucks numbers out of the value
	Math              map[string]string   `yaml:"math"`            // perform match across processed metrics
	SubParse          []Parse             `yaml:"sub_parse"`
	CustomAttributes  map[string]string   `yaml:"custom_attributes"` // set additional custom attributes
	ValueParser       map[string]string   `yaml:"value_parser"`      // find keys with regex, and parse the value with regex
	ValueTransformer  map[string]string   `yaml:"value_transformer"` // find key(s) with regex, and modify the value
	MetricParser      MetricParser        `yaml:"metric_parser"`     // to use the MetricParser for setting deltas and gauges a namespace needs to be set
	SampleFilter      []map[string]string `yaml:"sample_filter"`     // sample filter key pair values with regex
	Split             string              `yaml:"split"`             // default vertical, can be set to horizontal (column) useful for tabular outputs
	SplitBy           string              `yaml:"split_by"`          // character to split by
	SetHeader         []string            `yaml:"set_header"`        // manually set header column names
	Regex             bool                `yaml:"regex"`             // process SplitBy as regex
	RowHeader         int                 `yaml:"row_header"`        // set the row header, to be used with SplitBy
	RowStart          int                 `yaml:"row_start"`         // start from this line, to be used with SplitBy
	Logging           struct {            // log to insights
		Open bool `yaml:"open"` // log open related errors
	}
}

// Command Struct
type Command struct {
	Name             string            `yaml:"name"`              // required for database use
	EventType        string            `yaml:"event_type"`        // override eventType (currently used for db only)
	Shell            string            `yaml:"shell"`             // command shell
	Cache            string            `yaml:"cache"`             // use content from cache instead of a run command
	Run              string            `yaml:"run"`               // runs commands, but if database is set, then this is used to run queries
	Jmx              JMX               `yaml:"jmx"`               // if wanting to run different jmx endpoints to merge
	CompressBean     bool              `yaml:"compress_bean"`     // compress bean name //unused
	IgnoreOutput     bool              `yaml:"ignore_output"`     // can be useful for chaining commands together
	MetricParser     MetricParser      `yaml:"metric_parser"`     // not used yet
	CustomAttributes map[string]string `yaml:"custom_attributes"` // set additional custom attributes
	Output           string            `yaml:"output"`            // jmx, raw, json
	LineEnd          int               `yaml:"line_end"`          // stop processing command output after a certain amount of lines
	LineStart        int               `yaml:"line_start"`        // start from this line
	Timeout          int               `yaml:"timeout"`           // command timeout

	// Parsing Options - Body
	Split      string `yaml:"split"`       // default vertical, can be set to horizontal (column) useful for outputs that look like a table
	SplitBy    string `yaml:"split_by"`    // character/match to split by
	RegexMatch bool   `yaml:"regex_match"` // process SplitBy as a regex match
	GroupBy    string `yaml:"group_by"`    // group by character
	RowHeader  int    `yaml:"row_header"`  // set the row header, to be used with SplitBy
	RowStart   int    `yaml:"row_start"`   // start from this line, to be used with SplitBy

	// Parsing Options - Header
	SetHeader        []string `yaml:"set_header"`         // manually set header column names (used when split is is set to horizontal)
	HeaderSplitBy    string   `yaml:"header_split_by"`    // character/match to split header by
	HeaderRegexMatch bool     `yaml:"header_regex_match"` // process HeaderSplitBy as a regex match
}

// Prometheus struct
type Prometheus struct {
	Enable           bool              `yaml:"enable"`
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

// JMX struct
type JMX struct {
	Domain         string `yaml:"domain"`
	User           string `yaml:"user"`
	Pass           string `yaml:"pass"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	KeyStore       string `yaml:"key_store"`
	KeyStorePass   string `yaml:"key_store_pass"`
	TrustStore     string `yaml:"trust_store"`
	TrustStorePass string `yaml:"trust_store_pass"`
}

// Parse struct
type Parse struct {
	Type    string   `yaml:"type"` // perform a contains, match, hasPrefix or regex for specified key
	Key     string   `yaml:"key"`
	SplitBy []string `yaml:"split_by"`
}

// MetricParser Struct
type MetricParser struct {
	Namespace Namespace         `yaml:"namespace"`
	Metrics   map[string]string `yaml:"metrics"`  // inputBytesPerSecond: RATE
	AutoSet   bool              `yaml:"auto_set"` // if set to true, will attempt to do a contains instead of a direct key match, this is useful for setting multiple metrics
}

// Namespace Struct
type Namespace struct {
	// if neither of the below are set and the MetricParser is used, the namespace will default to the "Name" attribute
	CustomAttr   string   `yaml:"custom_attr"`   // set your own custom namespace attribute
	ExistingAttr []string `yaml:"existing_attr"` // utilise existing attributes and chain together to create a custom namespace
}

// Refresh Helper function used for testing
func Refresh() {
	EventCount = 0
	EventDropCount = 0
	ConfigsProcessed = 0
	Args.ConfigDir = ""
	Args.ConfigFile = ""
	Args.ContainerDiscovery = false
	Args.ContainerDiscoveryDir = ""
}

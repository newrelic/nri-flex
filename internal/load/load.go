package load

import (
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

// ArgumentList Available Arguments
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	ForceLog              bool   `default:"false" help:"Force log everything to standard out - useful for testing"`
	ForceLogEvent         bool   `default:"false" help:"Force create an event for everything - useful for testing"`
	OverrideIPMode        string `default:"" help:"Force override ipMode used for container discovery set as private or public - useful for testing"`
	Local                 bool   `default:"true" help:"Collect local entity info"`
	ConfigFile            string `default:"" help:"Set a specific config file - not usable for container discovery"`
	ConfigDir             string `default:"flexConfigs/" help:"Set directory of config files"`
	ContainerDiscoveryDir string `default:"flexContainerDiscovery/" help:"Set directory of auto discovery config files"`
	ContainerDiscovery    bool   `default:"false" help:"Enable container auto discovery"`
	DockerAPIVersion      string `default:"" help:"Force Docker client API version"`
	EventLimit            int    `default:"100" help:"Event limiter - max amount of events per execution"`
	Entity                string `default:"" help:"Manually set a remote entity name"`

	// not implemented yet
	ClusterModeKey string `default:"" help:"Set key used for cluster mode identification"`
	ClusterModeExp string `default:"60s" help:"Set cluster mode key identifier expiration"`
}

// Args Infrastructure SDK Arguments List
var Args ArgumentList

// Integration Infrastructure SDK Integration
var Integration *integration.Integration

// Entity Infrastructure SDK Entity
var Entity *integration.Entity

// Hostname current host
var Hostname string

// EventDropCount current number of events dropped due to limiter
var EventDropCount = 0

// EventCount number of events processed
var EventCount = 0

// EventDistribution number of events distributed per sample
var EventDistribution = map[string]int{}

// ConfigsProcessed number of configs processed
var ConfigsProcessed = 0

const (
	IntegrationName    = "com.kav91.nri-flex"     // IntegrationName Name
	IntegrationVersion = "0.3.3-pre"              // IntegrationVersion Version
	DefaultSplitBy     = ":"                      // unused currently
	DefaultTimeout     = 10000 * time.Millisecond // 10 seconds, used for raw commands
	DefaultPingTimeout = 5000                     // 5 seconds
	DefaultPostgres    = "postgres"
	DefaultMSSQLServer = "sqlserver"
	DefaultMySQL       = "mysql"
	DefaultOracle      = "ora"
	DefaultJmxPath     = "./nrjmx/"
	DefaultJmxHost     = "127.0.0.1"
	DefaultJmxPort     = "9999"
	DefaultJmxUser     = "admin"
	DefaultJmxPass     = "admin"
	DefaultIPMode      = "private"
	DefaultShell       = "/bin/sh"
	Public             = "public"
	Private            = "private"
	Jmx                = "jmx"
)

// Config YAML Struct
type Config struct {
	FileName         string // this will be set when files are read
	Name             string
	Global           Global
	APIs             []API
	Datastore        map[string][]interface{} `yaml:"datastore"`
	LookupStore      map[string][]string      `yaml:"lookup_store"`
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
	Timeout           int
	Method            string
	Payload           string
	DisableParentAttr bool              `yaml:"disable_parent_attr"`
	Headers           map[string]string `yaml:"headers"`
	StartKey          []string          `yaml:"start_key"`
	StoreLookups      map[string]string `yaml:"store_lookups"`
	StripKeys         []string          `yaml:"strip_keys"`
	LazyFlatten       []string          `yaml:"lazy_flatten"`
	SampleKeys        map[string]string `yaml:"sample_keys"`
	ReplaceKeys       map[string]string `yaml:"replace_keys"`
	RenameKeys        map[string]string `yaml:"rename_keys"`
	RemoveKeys        []string          `yaml:"remove_keys"`
	KeepKeys          []string          `yaml:"keep_keys"`     // inverse of removing keys
	ToLower           bool              `yaml:"to_lower"`      // convert all unicode letters mapped to their lower case.
	ConvertSpace      string            `yaml:"convert_space"` // convert spaces to another char
	SnakeToCamel      bool              `yaml:"snake_to_camel"`
	PercToDecimal     bool              `yaml:"perc_to_decimal"` // will check strings, and perform a trimRight for the %
	SubParse          []Parse           `yaml:"sub_parse"`
	CustomAttributes  map[string]string `yaml:"custom_attributes"` // set additional custom attributes
	MetricParser      MetricParser      `yaml:"metric_parser"`     // to use the MetricParser for setting deltas and gauges a namespace needs to be set
	SampleFilters     []string          `yaml:"sample_filters"`    // sample filter key pair values with regex
	Split             string            `yaml:"split"`             // default vertical, can be set to horizontal (column) useful for tabular outputs
	SplitBy           string            `yaml:"split_by"`          // character to split by
	SetHeader         []string          `yaml:"set_header"`        // manually set header column names
	Regex             bool              `yaml:"regex"`             // process SplitBy as regex
	RowHeader         int               `yaml:"row_header"`        // set the row header, to be used with SplitBy
	RowStart          int               `yaml:"row_start"`         // start from this line, to be used with SplitBy
	Logging           struct {          // log to insights
		Open bool `yaml:"open"` // log open related errors
	}
}

// Command Struct
type Command struct {
	Name             string            `yaml:"name"`              // required for database use
	Shell            string            `yaml:"shell"`             // command shell
	Run              string            `yaml:"run"`               // runs commands, but if database is set, then this is used to run queries
	Jmx              JMX               `yaml:"jmx"`               // if wanting to run different jmx endpoints to merge
	CompressBean     bool              `yaml:"compress_bean"`     // compress bean name //unused
	KeyFilters       map[string]string `yaml:"key_filters"`       // filter keys out with regex
	Output           string            `yaml:"output"`            // jmx, raw, json
	Split            string            `yaml:"split"`             // default vertical, can be set to horizontal (column) useful for outputs that look like a table
	SplitBy          string            `yaml:"split_by"`          // character to split by
	SetHeader        []string          `yaml:"set_header"`        // manually set header column names (used when split is is set to horizontal)
	GroupBy          string            `yaml:"group_by"`          // group by character
	Regex            bool              `yaml:"regex"`             // process SplitBy as regex
	RowHeader        int               `yaml:"row_header"`        // set the row header, to be used with SplitBy
	RowStart         int               `yaml:"row_start"`         // start from this line, to be used with SplitBy
	IgnoreOutput     bool              `yaml:"ignore_output"`     // can be useful for chaining commands together
	MetricParser     MetricParser      `yaml:"metric_parser"`     // not used yet
	CustomAttributes map[string]string `yaml:"custom_attributes"` // set additional custom attributes
	// SplitBy      string `yaml:"split_by"`      // performs horizontal split by eg. split_by ":" and splits  "myMetric:myValue"
	// ColSplitBy   string `yaml:"col_split_by"`  // performs vertical split, if data output is in a table type output, first row is taken as the headers, and following as the samples
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
	Histogram        bool              `yaml:"histogram"` // if flattening by default, create a full histogram sample
	Summary          bool              `yaml:"summary"`   // if flattening by default, create a full summary sample
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

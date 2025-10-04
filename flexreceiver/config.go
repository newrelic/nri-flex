package flexreceiver

import (
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
)

// Config defines the configuration for the Flex receiver.
type Config struct {
	// The directory containing Flex configuration files.
	ConfigDir string `mapstructure:"config_dir"`
	// A single Flex configuration file.
	ConfigFile string `mapstructure:"config_file"`
	// The interval at which to collect metrics.
	CollectionInterval time.Duration `mapstructure:"collection_interval"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the receiver configuration is valid.
func (cfg *Config) Validate() error {
	if cfg.CollectionInterval <= 0 {
		return errors.New("collection_interval must be a positive duration")
	}
	if cfg.ConfigDir == "" && cfg.ConfigFile == "" {
		return errors.New("either 'config_dir' or 'config_file' must be provided")
	}
	return nil
}

package outputs

import (
	"os"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestConfigPath_Override(t *testing.T) {
	const (
		configPath   = "CONFIG_PATH"
		expectedPath = "/usr/local/foo"
	)

	oldPath, ok := os.LookupEnv(configPath)
	os.Setenv(configPath, expectedPath)
	defer func() {
		if ok {
			os.Setenv(configPath, oldPath)
		} else {
			os.Unsetenv(configPath)
		}
	}()

	InfraIntegration()

	if load.Args.ConfigPath != expectedPath ||
		load.Args.ConfigFile != expectedPath {
		t.Errorf("config_path and config_file should equal %s (actual values: %s, %s",
			expectedPath, load.Args.ConfigPath, load.Args.ConfigFile)
	}
}

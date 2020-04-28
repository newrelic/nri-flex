// +build windows

package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"

	//"github.com/newrelic/nri-flex/integration-test/integration"
	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/load"

	"github.com/stretchr/testify/require"
)

var flexMain = filepath.Join("..", "cmd", "nri-flex", "nri-flex.go")

func Test_WindowsCommands_ReturnsData(t *testing.T) {
	// given
	load.Refresh()

	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// Load a single config file
	var configs []load.Config
	configFile, _ := os.Stat(filepath.Join("configs", "windows-cmd-test.yml"))
	err := config.LoadFile(&configs, configFile, "configs")
	require.NoError(t, err)

	// when
	errs := config.RunFiles(&configs)
	require.Empty(t, errs)

	metricsSet := load.Entity.Metrics
	require.NotEmpty(t, metricsSet)

	for _, ms := range metricsSet {
		require.NotNil(t, ms.Metrics, "status")
		require.NotNil(t, ms.Metrics, "name")
		require.NotNil(t, ms.Metrics, "displayname")
	}

	// check for a specific service, because Flex ingests everything, even output "garbage"
	// any Windows should always have the Themes service, so check for that
	var found bool
	for _, ms := range metricsSet {
		if ms.Metrics["name"] == "Themes" {
			found = true
		}
	}

	require.Truef(t, found, "didn't find the 'Themes' service. check that the configuration is correct")
}

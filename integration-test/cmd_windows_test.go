// +build integration
// +build windows

package integration

import (
	"path/filepath"
	"testing"

	sdk "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/require"
)

func Test_WindowsCommands_ReturnsData(t *testing.T) {
	configDirPath := filepath.Join("configs", "windows")

	load.Refresh()

	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set file to load
	load.Args.ConfigFile = filepath.Join(configDirPath, "windows-cmd-test.yml")

	// when
	err := integration.RunFlex(integration.FlexModeDefault)
	require.NoError(t, err)

	//then
	metricsSet := load.Entity.Metrics
	require.NotEmpty(t, metricsSet)

	for _, ms := range metricsSet {
		if ms.Metrics["event_type"] == "flexStatusSample" {
			continue
		}
		require.NotNil(t, ms.Metrics["status"], "status")
		require.NotNil(t, ms.Metrics["name"], "name")
		require.NotNil(t, ms.Metrics["displayname"], "displayname")
	}

	// check for a specific service, because Flex ingests everything, even output "garbage"
	// any Windows version should always have the Themes service, so check for that
	var found bool
	for _, ms := range metricsSet {
		if ms.Metrics["name"] == "Themes" {
			found = true
		}
	}

	require.Truef(t, found, "didn't find the 'Themes' service. check that the configuration is correct")
}

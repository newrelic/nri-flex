// +build windows

package integration

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-flex/integration-test/integration"

	"github.com/newrelic/nri-flex/integration-test/gofile"

	"github.com/stretchr/testify/require"
)

var flexMain = filepath.Join("..", "cmd", "nri-flex", "nri-flex.go")

func Test_WindowsCommands_ReturnsData(t *testing.T) {

	configFilePath := filepath.Join("configs", "windows-cmd-test.yml")
	flexMainPath, _ := filepath.Abs(flexMain)

	// WHEN executing nri flex with the provided configuration
	stdout, err := gofile.Run(flexMainPath, "-verbose", "-config_path="+configFilePath)
	require.NoError(t, err)
	payload := integration.JSON{}
	require.NoError(t, json.Unmarshal(stdout, &payload))

	// THEN samples are received with the metrics properly extracted
	require.NotEmpty(t, payload.Data)
	found := false
	for _, data := range payload.Data {
		if data.Entity == nil {
			continue
		}
		found = true
		require.Len(t, data.Metrics, 1)
		m := data.Metrics[0]
		require.Equal(t, "windowsServiceListSample", m["event_type"])
		require.Equal(t, "com.newrelic.nri-flex", m["integration_name"])
		require.Contains(t, m, "integration_version")
		require.Contains(t, m, "status")
		require.Contains(t, m, "name")
		require.Contains(t, m, "displayname")
	}
	require.Truef(t, found, "did not find any result in the integration %s", string(stdout))
}

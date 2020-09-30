// +build integration
// +build linux

package integration_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-flex/integration-test/integration"

	"github.com/newrelic/nri-flex/integration-test/gofile"

	"github.com/stretchr/testify/require"
)

var flexMain = filepath.Join("..", "cmd", "nri-flex", "nri-flex.go")

func TestHTTP(t *testing.T) {
	// WHEN executing nri flex with the provided http-test.yml configuration
	stdout, err := gofile.Run(flexMain, "-verbose", "-config_path="+filepath.Join("configs", "http-test.yml"))
	require.NoError(t, err)
	payload := integration.JSON{}
	require.NoError(t, json.Unmarshal(stdout, &payload))

	// THEN an NginxSample is received with the metrics properly extracted
	require.NotEmpty(t, payload.Data)
	found := false
	for _, data := range payload.Data {
		if data.Entity == nil {
			continue
		}
		found = true
		require.Equal(t, "https-server", data.Entity.Name)
		require.Equal(t, "instance", data.Entity.Type)
		require.Len(t, data.Metrics, 1)
		m := data.Metrics[0]
		require.Equal(t, "NginxSample", m["event_type"])
		require.Contains(t, m, "flex.commandTimeMs")
		require.Equal(t, "com.newrelic.nri-flex", m["integration_name"])
		require.Contains(t, m, "integration_version")
		require.InDelta(t, m["net.connectionsActive"], 43, 0.1)
		require.InDelta(t, m["net.connectionsReading"], 0, 0.1)
		require.InDelta(t, m["net.connectionsWaiting"], 38, 0.1)
		require.InDelta(t, m["net.connectionsWriting"], 5, 0.1)
		require.InDelta(t, m["net.connectionsAcceptedPerSecond"], 8000, 0.1)
		require.InDelta(t, m["net.handledPerSecond"], 7368, 0.1)
		require.InDelta(t, m["net.requestsPerSecond"], 10993, 0.1)
		require.InDelta(t, m["net.connectionsDroppedPerSecond"], 8000-7368, 0.1)
	}
	require.Truef(t, found, "did not find any 'http-server' nginx instance in the integration %s", string(stdout))
}

func TestHTTPS(t *testing.T) {
	// GIVEN a Nginx HTTPS server (already running in the Dockerfile_https container)

	// WHEN executing Flex with the provided http-test.yml configuration
	stdout, err := gofile.Run(flexMain, "-verbose", "-config_path="+filepath.Join("configs", "https-test.yml"))
	require.NoError(t, err)
	payload := integration.JSON{}
	require.NoError(t, json.Unmarshal(stdout, &payload))

	// THEN an NginxSample is received with the metrics properly extracted
	require.NotEmpty(t, payload.Data)
	found := false
	for _, data := range payload.Data {
		if data.Entity == nil {
			continue
		}
		found = true
		require.Equal(t, "https-server", data.Entity.Name)
		require.Equal(t, "instance", data.Entity.Type)
		require.Len(t, data.Metrics, 1)
		m := data.Metrics[0]
		require.Equal(t, "NginxSample", m["event_type"])
		require.Contains(t, m, "flex.commandTimeMs")
		require.Equal(t, "com.newrelic.nri-flex", m["integration_name"])
		require.Contains(t, m, "integration_version")
		require.InDelta(t, m["net.connectionsActive"], 43, 0.1)
		require.InDelta(t, m["net.connectionsReading"], 0, 0.1)
		require.InDelta(t, m["net.connectionsWaiting"], 38, 0.1)
		require.InDelta(t, m["net.connectionsWriting"], 5, 0.1)
		require.InDelta(t, m["net.connectionsAcceptedPerSecond"], 8000, 0.1)
		require.InDelta(t, m["net.handledPerSecond"], 7368, 0.1)
		require.InDelta(t, m["net.requestsPerSecond"], 10993, 0.1)
		require.InDelta(t, m["net.connectionsDroppedPerSecond"], 8000-7368, 0.1)
	}
	require.Truef(t, found, "did not find any 'https-server' nginx instance in the integration %s", string(stdout))
}

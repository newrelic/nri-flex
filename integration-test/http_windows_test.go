// +build integration
// +build windows

package integration

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	sdk "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/integration-test/gofile"
	"github.com/newrelic/nri-flex/internal/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/require"
)

func Test_WindowsHttp_ReturnsData(t *testing.T) {
	go startServer("http")

	load.Refresh()
	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set file to load
	configDirPath := filepath.Join("configs", "windows")
	load.Args.ConfigFile = filepath.Join(configDirPath, "windows-http-test.yml")

	// when
	err := integration.RunFlex(integration.FlexModeDefault)
	require.NoError(t, err)

	//then
	metricsSet := load.Entity.Metrics
	checkOutput(t, metricsSet, 1)
}

func Test_WindowsHttps_ReturnsData(t *testing.T) {
	go startServer("https")

	load.Refresh()
	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set file to load
	configDirPath := filepath.Join("configs", "windows")
	load.Args.ConfigFile = filepath.Join(configDirPath, "windows-https-test.yml")

	// when
	err := integration.RunFlex(integration.FlexModeDefault)
	require.NoError(t, err)

	//then
	metricsSet := load.Entity.Metrics
	checkOutput(t, metricsSet, 1)
}

func Test_WindowsHttps_ConfigFolder_ReturnsData(t *testing.T) {
	go startServer("http")
	go startServer("https")

	load.Refresh()
	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set file to load
	configDirPath := filepath.Join("configs", "windows")
	load.Args.ConfigDir = configDirPath
	load.Args.Verbose = true

	// when
	err := integration.RunFlex(integration.FlexModeDefault)
	require.NoError(t, err)

	//then
	metricsSet := load.Entity.Metrics
	checkOutput(t, metricsSet, 2)
}

func checkOutput(t *testing.T, metrics []*metric.Set, expectedCount int) {
	require.NotEmpty(t, metrics)

	var actualCount int
	for _, ms := range metrics {
		if ms.Metrics["event_type"] != "WindowsHttpSample" {
			continue
		}
		require.Equal(t, "WindowsHttpSample", ms.Metrics["event_type"])
		require.Equal(t, 10.0, ms.Metrics["cpu"], "cpu")
		require.Equal(t, float64(3500), ms.Metrics["memory"], "memory")
		require.Equal(t, float64(500), ms.Metrics["disk"], "disk")
		actualCount++
	}

	require.Equal(t, expectedCount, actualCount)
}

func startServer(server string) {
	_, err := gofile.Run(filepath.Join("https-server", server, "server.go"))
	if err != nil {
		fmt.Println(err)
	}
}

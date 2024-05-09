//go:build integration && windows
// +build integration,windows

package integration_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	sdk "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
	"github.com/stretchr/testify/require"
)

func Test_WindowsHttp_ReturnsData(t *testing.T) {
	err := os.Setenv("TEST_HTTP_SERVER_PORT", startServer(t, false))
	require.NoError(t, err)

	load.Refresh()
	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set file to load
	configDirPath := filepath.Join("configs", "windows")
	load.Args.ConfigFile = filepath.Join(configDirPath, "windows-http-test.yml")

	// when
	r := runtime.GetDefaultRuntime()
	err = runtime.RunFlex(r)
	require.NoError(t, err)

	//then
	metricsSet := load.Entity.Metrics
	checkOutput(t, metricsSet, 1)
}

func Test_WindowsHttps_ReturnsData(t *testing.T) {
	err := os.Setenv("TEST_HTTPS_SERVER_PORT", startServer(t, true))
	require.NoError(t, err)

	load.Refresh()
	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set file to load
	configDirPath := filepath.Join("configs", "windows")
	load.Args.ConfigFile = filepath.Join(configDirPath, "windows-https-test.yml")

	// when
	r := runtime.GetDefaultRuntime()
	err = runtime.RunFlex(r)
	require.NoError(t, err)

	//then
	metricsSet := load.Entity.Metrics
	checkOutput(t, metricsSet, 1)
}

func Test_WindowsHttps_ConfigFolder_ReturnsData(t *testing.T) {
	err := os.Setenv("TEST_HTTP_SERVER_PORT", startServer(t, false))
	require.NoError(t, err)
	err = os.Setenv("TEST_HTTPS_SERVER_PORT", startServer(t, true))
	require.NoError(t, err)

	load.Refresh()
	i, _ := sdk.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")

	// set file to load
	configDirPath := filepath.Join("configs", "windows")
	load.Args.ConfigDir = configDirPath
	load.Args.Verbose = true

	// when
	r := runtime.GetDefaultRuntime()
	err = runtime.RunFlex(r)
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

func startServer(t *testing.T, tls bool) (port string) {
	t.Helper()

	srv := &httptest.Server{}
	if tls {
		srv = httptest.NewTLSServer(http.HandlerFunc(serveJSON))
	} else {
		srv = httptest.NewServer(http.HandlerFunc(serveJSON))
	}

	url, err := url.Parse(srv.URL)
	require.NoError(t, err)
	return url.Port()
}

func serveJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	_, _ = w.Write([]byte(`
	{
		"metrics": [
			{
			 "cpu": 10.0,
			 "memory": 3500,
			 "disk": 500
			} 
		]
	}
	`))
}

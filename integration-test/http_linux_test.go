package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-flex/integration-test/ijson"

	"github.com/newrelic/nri-flex/integration-test/gofile"

	"github.com/stretchr/testify/require"
)

var flexMain = filepath.Join("..", "cmd", "nri-flex", "nri-flex.go")

func TestHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`Active connections: 43
server accepts handled requests
8000 7368 10993
Reading: 0 Writing: 5 Waiting: 38
`))
	}))
	defer server.Close()

	configPath, err := replaceDiscoveryPort(filepath.Join("configs", "http-test.yml"), server)
	require.NoError(t, err)

	stdout, err := gofile.Run(flexMain, "-config_path="+configPath)
	require.NoError(t, err)

	payload := ijson.Payload{}
	require.NoError(t, json.Unmarshal(stdout, &payload))

	require.NotEmpty(t, payload.Data)
	found := false
	for _, data := range payload.Data {
		if data.Entity == nil {
			continue
		}
		found = true
		require.Equal(t, "127.0.0.1", data.Entity.Name)
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
	require.Truef(t, found, "did not find any 127.0.0.1 nginx instance in the integration %s", string(stdout))
}

// returns the path of a temporary config file loaded from the provided file path, replacing the
// ${discovery.port} by the port of the server
func replaceDiscoveryPort(configFilePath string, server *httptest.Server) (string, error) {
	configTemplate, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	surl, err := url.Parse(server.URL)
	if err != nil {
		return "", err
	}
	replaced := bytes.ReplaceAll(configTemplate, []byte("${discovery.port}"), []byte(surl.Port()))
	cfg, err := ioutil.TempFile("", "httptest")
	if err != nil {
		return "", err
	}
	defer cfg.Close()
	_, err = cfg.Write(replaced)
	if err != nil {
		return "", err
	}
	return cfg.Name(), nil
}

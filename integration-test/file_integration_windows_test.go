// +build integration
// +build windows

// The tests in this file are supposed to be run in the CI using docker-compose
// You can run then from the IDE or manually but you'll need some setup first
package integration_test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestConfig_Read_Files(t *testing.T) {
	tests := map[string]struct {
		configFile string
		want       Metrics
	}{
		"jsonEtcdSelf": {"configs/windows-json-test.yml", wantJsonMetrics},
		"csvTest":      {"configs/windows-csv-test.yml", wantCsvMetrics},
	}

	// copy payloads to the current user's TMP folder plus a folder with spaces
	assert.NoError(t, exec.Command("cmd", "/C", "xcopy", "/y", "..\\test\\payloads\\etcdSelf.json", filepath.Join("%TMP%", "test payloads", "etcdSelf.json*")).Run())
	assert.NoError(t, exec.Command("cmd", "/C", "xcopy", "/y", "..\\test\\payloads\\test.csv", filepath.Join("%TMP%", "test payloads", "test.csv*")).Run())

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			flexOutput := runConfigFile(t, tc.configFile)

			for i, metrics := range tc.want {
				assert.NoError(t, isMapSubset(metrics, flexOutput[i].Metrics))
			}
		})
	}
}

// +build integration
// +build windows

// The tests in this file are supposed to be run in the CI using docker-compose
// You can run then from the IDE or manually but you'll need some setup first
package integration_test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestConfig_Read_Files(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test payloads")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tmpDir))
	}()

	assert.NoError(t, os.Setenv("MY_TEMP_FOLDER", tmpDir))
	defer func() {
		assert.NoError(t, os.Unsetenv("MY_TEMP_FOLDER"))
	}()

	tests := map[string]struct {
		configFile  string
		payloadFile string
		want        Metrics
	}{
		"jsonEtcdSelf": {"configs/windows-json-test.yml", "..\\test\\payloads\\etcdSelf.json", wantJsonMetrics},
		"csvTest":      {"configs/windows-csv-test.yml", "..\\test\\payloads\\test.csv", wantCsvMetrics},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// copy the payload file to the temporary folder
			assert.NoError(t, exec.Command("cmd", "/C", "xcopy", "/y", tc.payloadFile, tmpDir).Run())

			flexOutput := runConfigFile(t, tc.configFile)

			for i, metrics := range tc.want {
				assert.NoError(t, isMapSubset(metrics, flexOutput[i].Metrics))
			}
		})
	}
}

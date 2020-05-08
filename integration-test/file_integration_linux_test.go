// +build integration
// +build linux

// The tests in this file are supposed to be run in the CI using docker-compose
// You can run then from the IDE or manually but you'll need some setup first
package integration_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_Read_Files(t *testing.T) {
	tests := map[string]struct {
		configFile string
		want       Metrics
	}{
		"jsonEtcdSelf": {"configs/json-test.yml", wantJsonMetrics},
		"csvTest":      {"configs/csv-test.yml", wantCsvMetrics},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			flexOutput := runConfigFile(t, tc.configFile)

			for i, metrics := range tc.want {
				assert.NoError(t, isMapSubset(metrics, flexOutput[i].Metrics))
			}
		})
	}
}

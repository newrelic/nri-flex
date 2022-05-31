package testbed

import (
	"encoding/json"
	"fmt"
	Integration "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	errOutputsDoNotMatch = fmt.Errorf("Outputs do not match")
	errDifferentMetrics  = fmt.Errorf("Outputs have different metrics")
)

type ExecutionValidator interface {
	// Validate given stdout and stderr logs
	Validate(*testing.T, string, string) error
}

type MetricValidator struct {
	matchStdout Integration.Integration
	matchStderr string
}

func NewMetricValidator(stdout, stderr string) (*MetricValidator, error) {
	m := &MetricValidator{matchStderr: stderr}
	err := json.Unmarshal([]byte(stdout), &m.matchStdout)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (e *MetricValidator) Validate(t *testing.T, stdout string, stderr string) error {
	var actualDataset Integration.Integration
	err := json.Unmarshal([]byte(stdout), &actualDataset)
	if err != nil {
		return err
	}

	// checks that the number of entities is the same
	assert.Equal(t, len(e.matchStdout.Entities), len(actualDataset.Entities))

	for i, expectedSet := range e.matchStdout.Entities {
		// TODO: check nil cases
		actualMetrics := actualDataset.Entities[i].Metrics
		assert.Equal(t, len(expectedSet.Metrics), len(actualMetrics))
		for j, expectedMetrics := range expectedSet.Metrics {
			assert.Equal(t, len(expectedMetrics.Metrics), len(actualMetrics[j].Metrics))
		}
	}

	return nil
}

package testbed

import (
	"encoding/json"
	Integration "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
	"testing"
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

	totalExpectedMetrics := 0
	totalActualMetrics := 0

	// two loops because the order of the metrics can be different depending on the execution
	for _, expectedSet := range e.matchStdout.Entities {
		// TODO: check nil cases
		for _, expectedMetrics := range expectedSet.Metrics {
			totalExpectedMetrics += len(expectedMetrics.Metrics)
		}
	}
	for _, actualSet := range actualDataset.Entities {
		// TODO: check nil cases
		for _, actualMetrics := range actualSet.Metrics {
			totalActualMetrics += len(actualMetrics.Metrics)
		}
	}

	assert.Equal(t, totalExpectedMetrics, totalActualMetrics)

	return nil
}

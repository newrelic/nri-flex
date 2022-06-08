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

	assert.Equal(t, e.matchStdout.Name, actualDataset.Name)
	assert.Equal(t, e.matchStdout.ProtocolVersion, actualDataset.ProtocolVersion)

	// checks that the number of entities is the same
	assert.Equal(t, len(e.matchStdout.Entities), len(actualDataset.Entities))

	totalExpectedMetrics, totalExpectedEvents := 0, 0
	totalActualMetrics, totalActualEvents := 0, 0

	for i, expectedSet := range e.matchStdout.Entities {
		// two loops because the order of the metrics can be different depending on the execution
		for _, expectedMetrics := range expectedSet.Metrics {
			totalExpectedMetrics += len(expectedMetrics.Metrics)
		}
		for _, actualMetrics := range actualDataset.Entities[i].Metrics {
			totalActualMetrics += len(actualMetrics.Metrics)
		}
		totalExpectedEvents += len(expectedSet.Events)
		totalActualEvents += len(actualDataset.Entities[i].Events)
	}

	assert.Equal(t, totalExpectedMetrics, totalActualMetrics)
	assert.Equal(t, totalActualEvents, totalActualEvents)

	return nil
}

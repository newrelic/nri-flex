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

// IntegrationValidator validates an integration stdout
type IntegrationValidator struct {
	expectedIntegration Integration.Integration
}

// NewIntegrationValidator returns a new IntegrationValidator instance
func NewIntegrationValidator(stdout, stderr string) (*IntegrationValidator, error) {
	m := &IntegrationValidator{}
	err := json.Unmarshal([]byte(stdout), &m.expectedIntegration)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Validate verifies the number of entities, events and metrics of a given integration output, stderr is not validated
func (e *IntegrationValidator) Validate(t *testing.T, stdout string, stderr string) error {
	var actualDataset Integration.Integration

	err := json.Unmarshal([]byte(stdout), &actualDataset)
	if err != nil {
		return err
	}

	assert.Equal(t, e.expectedIntegration.Name, actualDataset.Name)
	assert.Equal(t, e.expectedIntegration.ProtocolVersion, actualDataset.ProtocolVersion)

	// checks that the number of entities is the same
	assert.Equal(t, len(e.expectedIntegration.Entities), len(actualDataset.Entities))

	totalExpectedMetrics, totalExpectedEvents := 0, 0
	expectedMetricsKeys := make(map[string]interface{})
	totalActualMetrics, totalActualEvents := 0, 0
	actualMetricsKeys := make(map[string]interface{})

	for i, expectedSet := range e.expectedIntegration.Entities {
		// two loops because the order of the metrics can be different depending on the execution
		for _, expectedMetrics := range expectedSet.Metrics {
			totalExpectedMetrics += len(expectedMetrics.Metrics)
			for key, value := range expectedMetrics.Metrics {
				expectedMetricsKeys[key] = value
			}
		}
		for _, actualMetrics := range actualDataset.Entities[i].Metrics {
			totalActualMetrics += len(actualMetrics.Metrics)
			for key, value := range actualMetrics.Metrics {
				actualMetricsKeys[key] = value
			}
		}
		totalExpectedEvents += len(expectedSet.Events)
		totalActualEvents += len(actualDataset.Entities[i].Events)
	}

	assert.Equal(t, totalExpectedMetrics, totalActualMetrics)
	assert.Equal(t, totalActualEvents, totalActualEvents)
	// verify that all keys exists in both maps and types match
	for expectedKey, expectedValue := range expectedMetricsKeys {
		assert.Contains(t, actualMetricsKeys, expectedKey)
		assert.IsType(t, expectedValue, actualMetricsKeys[expectedKey])
	}

	return nil
}

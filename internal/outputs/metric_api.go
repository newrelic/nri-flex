package outputs

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

// SendToMetricAPI - Send processed events to insights
func SendToMetricAPI() {
	key := load.Args.InsightsAPIKey
	if load.Args.MetricAPIKey != "" {
		key = load.Args.MetricAPIKey
	}
	jsonData, err := json.Marshal(load.MetricsPayload) // may need to implement some sort of chunking or batching
	if err != nil {
		logger.Flex("error", err, "failed to marshal", false)
	} else {
		postRequest(load.Args.MetricAPIUrl, key, jsonData)
		if load.Args.InsightsOutput {
			fmt.Println(string(jsonData))
		}
	}
}

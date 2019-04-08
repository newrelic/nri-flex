package outputs

import (
	"encoding/json"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// SendToInsights - Send processed events to insights
func SendToInsights() {
	modifyEventType()
	jsonData, err := json.Marshal(load.Entity.Metrics)
	if err != nil {
		logger.Flex("debug", err, "failed to marshal", false)
	} else {
		postRequest(load.Args.InsightsURL, load.Args.InsightsAPIKey, jsonData)
		// empty the infrastructure entity metrics by default
		if !load.Args.InsightsOutput {
			load.Entity.Metrics = []*metric.Set{}
		}
	}
}

// modifyEventType insights uses eventType key in camel case whereas infrastructure uses event_type
func modifyEventType() {
	for _, event := range load.Entity.Metrics {
		event.Metrics["eventType"] = event.Metrics["event_type"].(string)
		delete(event.Metrics, "event_type")
	}
}

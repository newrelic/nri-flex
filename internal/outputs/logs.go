/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
)

// GetLogMetricBatches batch metrics by entity with a maximum batch size defined by 'LogBatchSize' config.
func GetLogMetricBatches() [][]*metric.Set {
	var result [][]*metric.Set
	for _, entity := range load.Integration.Entities {
		// split the payload into smaller batches so that the payload size does not exceed the Log endpoint size limit
		batchSize := load.Args.InsightBatchSize
		sizeMetrics := len(entity.Metrics)
		batches := sizeMetrics/batchSize + 1

		modifyLogEventType(entity)

		for i := 0; i < batches; i++ {
			start := i * batchSize
			end := start + batchSize
			if end > sizeMetrics {
				end = sizeMetrics
			}
			result = append(result, (entity.Metrics)[start:end])
		}
		// empty the infrastructure entity metrics by default
		entity.Metrics = []*metric.Set{}
	}
	return result
}

// SendBatchToLogApi - Send processed events to log api.
func SendBatchToLogApi(metrics []*metric.Set) error {
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("log api: failed to marshal json, %v", err)
	}

	load.Logrus.Debugf("posting %d events to log api", len(metrics))

	if load.Args.LogOutput {
		fmt.Println(string(jsonData))
	}

	err = postRequest(load.Args.LogApiURL, load.Args.LogApiKey, jsonData)
	if err != nil {
		return fmt.Errorf("log api: sending events failed, %v", err)
	}

	return nil
}

// modifyEventType insights uses eventType key in camel case whereas infrastructure uses event_type
func modifyLogEventType(entity *integration.Entity) {
	for _, event := range entity.Metrics {
		if event.Metrics["event_type"] != nil {
			event.Metrics["flexEventType"] = event.Metrics["event_type"].(string)
		}
		delete(event.Metrics, "event_type")
	}
}

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/nri-flex/internal/load"
)

// SendToMetricAPI - Send processed events to insights
func SendToMetricAPI() error {
	key := load.Args.InsightsAPIKey
	if load.Args.MetricAPIKey != "" {
		key = load.Args.MetricAPIKey
	}
	jsonData, err := json.Marshal(load.MetricsStore.Data) // may need to implement some sort of chunking or batching
	if err != nil {
		return fmt.Errorf("metrics api: failed to marshal json, %v", err)
	}

	if load.Args.InsightsOutput {
		fmt.Println(string(jsonData))
	}

	err = postRequest(load.Args.MetricAPIUrl, key, jsonData)
	if err != nil {
		return fmt.Errorf("insights: sending events failed, %v", err)
	}
	return nil
}

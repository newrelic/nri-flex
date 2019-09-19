/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

// SendToMetricAPI - Send processed events to insights
func SendToMetricAPI() {
	key := load.Args.InsightsAPIKey
	if load.Args.MetricAPIKey != "" {
		key = load.Args.MetricAPIKey
	}
	jsonData, err := json.Marshal(load.MetricsStore.Data) // may need to implement some sort of chunking or batching
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("metrics: failed to marshal")
	} else {
		postRequest(load.Args.MetricAPIUrl, key, jsonData)
		if load.Args.InsightsOutput {
			fmt.Println(string(jsonData))
		}
	}
}

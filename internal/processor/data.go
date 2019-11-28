/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

// RunDataHandler handles the data received for processing
func RunDataHandler(dataSets []interface{}, samplesToMerge *map[string][]interface{}, i int, cfg *load.Config) {
	load.Logrus.WithFields(logrus.Fields{
		"name": cfg.Name,
	}).Debug("processor: running data handler")
	for _, dataSet := range dataSets {
		switch dataSet := dataSet.(type) {
		case map[string]interface{}:
			ds := dataSet
			processDataSet(&ds, samplesToMerge, i, cfg)
		case []interface{}:
			nextDataSets := dataSet
			RunDataHandler(nextDataSets, samplesToMerge, i, cfg)
		default:
			load.Logrus.WithFields(logrus.Fields{
				"name": cfg.Name,
			}).Debug("processor: not sure what to do with this?!")
		}
	}
}

// processDataSet performs the core flattening on the map[string]interface then executes createMetricSets finally
func processDataSet(dataSet *map[string]interface{}, samplesToMerge *map[string][]interface{}, i int, cfg *load.Config) {
	ds := (*dataSet)

	FindStartKey(&ds, cfg.APIs[i].StartKey, cfg.APIs[i].InheritAttributes) // start at a later part in the received data
	StripKeys(&ds, cfg.APIs[i].StripKeys)                                  // remove before flattening
	RunLazyFlatten(&ds, cfg, i)                                            // perform lazy flatten if needed
	flattenedData := FlattenData(ds, map[string]interface{}{}, "", cfg.APIs[i].SampleKeys, &cfg.APIs[i])

	// also strip from flattened data
	for _, stripKey := range cfg.APIs[i].StripKeys {
		delete(flattenedData, stripKey)
		delete(flattenedData, stripKey+"Samples")
	}

	mergedData := FinalMerge(flattenedData)

	// hren: moved this portion to CreateMetricSets
	// if cfg.APIs[i].Merge != "" {
	// 	for _, mergeItem := range mergedData {
	// 		mergeItem.(map[string]interface{})["_sampleNo"] = i
	// 		// hren overwrite event_type if it is merge operation
	// 		mergeItem.(map[string]interface{})["event_type"] = cfg.APIs[i].Merge
	// 		(*samplesToMerge)[cfg.APIs[i].Merge] = append((*samplesToMerge)[cfg.APIs[i].Merge], mergeItem)
	// 	}
	// 	mergedSample = true
	// }

	if cfg.APIs[i].Merge == "" {
		CreateMetricSets(mergedData, cfg, i, false, nil)
	} else {
		CreateMetricSets(mergedData, cfg, i, true, samplesToMerge)
	}
}

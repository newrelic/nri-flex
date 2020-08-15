/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

// RunDataHandler handles the data received for processing
// The originalAPINo is to track the original API sequential No. in the Flex config file. This is to diffentiate the new API Seq No. created by StoreLookup.
// The originalAPINo is used for Merge and Join operation
func RunDataHandler(dataSets []interface{}, samplesToMerge *load.SamplesToMerge, i int, cfg *load.Config, originalAPINo int) {
	load.Logrus.WithFields(logrus.Fields{
		"name": cfg.Name,
	}).Debug("processor-data: running data handler")

	if cfg.APIs[originalAPINo].Jq != "" {
		dataSets = runJq(dataSets, cfg.APIs[originalAPINo])
	}

	for _, dataSet := range dataSets {
		switch dataSet := dataSet.(type) {
		case map[string]interface{}:
			ds := dataSet
			processDataSet(&ds, samplesToMerge, i, cfg, originalAPINo)
		case []interface{}:
			nextDataSets := dataSet
			RunDataHandler(nextDataSets, samplesToMerge, i, cfg, originalAPINo)
		default:
			load.Logrus.WithFields(logrus.Fields{
				"name": cfg.Name,
			}).Debugf("processor-data: unsupported data type %T", dataSet)
		}
	}
}

// processDataSet performs the core flattening on the map[string]interface then executes createMetricSets finally
func processDataSet(dataSet *map[string]interface{}, samplesToMerge *load.SamplesToMerge, i int, cfg *load.Config, originalAPINo int) {
	ds := (*dataSet)

	if cfg.LookupStore == nil {
		cfg.LookupStore = map[string]map[string]struct{}{}
	}

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

	if cfg.APIs[i].Merge == "" {
		CreateMetricSets(mergedData, cfg, i, false, nil, originalAPINo)
	} else {
		CreateMetricSets(mergedData, cfg, i, true, samplesToMerge, originalAPINo)
	}
}

func runJq(dataSets interface{}, api load.API) []interface{} {
	// due to the current design, dataSets is always a slice
	// so an api that returns a map will still be a slice returned for processing
	// for better usability if there is no array prefix access the first item in the slice automatically for the user
	if !strings.HasPrefix(api.Jq, ".[") {
		api.Jq = fmt.Sprintf(".[0]%v", api.Jq)
	}

	query, err := gojq.Parse(api.Jq)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"api": api.Name,
		}).WithError(err).Error("jq: failed to parse")
		return []interface{}{}
	}

	iter := query.Run(dataSets)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		switch value := v.(type) {
		case []interface{}:
			return value
		case map[string]interface{}:
			return []interface{}{value}
		case error:
			load.Logrus.WithFields(logrus.Fields{
				"api": api.Name,
				"jq":  api.Jq,
			}).WithError(err).Debug("jq: failed to process")
			return []interface{}{}
		}

	}

	return []interface{}{}
}

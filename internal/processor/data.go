package processor

import (
	"fmt"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

// RunDataHandler handles the data received for processing
func RunDataHandler(dataSets []interface{}, samplesToMerge *map[string][]interface{}, i int, cfg *load.Config) {
	logger.Flex("debug", nil, fmt.Sprintf("running data handler for %v", cfg.Name), false)
	for _, dataSet := range dataSets {
		switch dataSet := dataSet.(type) {
		case map[string]interface{}:
			ds := dataSet
			processDataSet(&ds, samplesToMerge, i, cfg)
		case []interface{}:
			nextDataSets := dataSet
			RunDataHandler(nextDataSets, samplesToMerge, i, cfg)
		default:
			logger.Flex("debug", nil, "not sure what to do with "+cfg.Name, false)
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
	mergedSample := false

	if len(mergedData) == 1 {
		if cfg.APIs[i].Merge != "" {
			mergedData[0].(map[string]interface{})["_sampleNo"] = i
			(*samplesToMerge)[cfg.APIs[i].Merge] = append((*samplesToMerge)[cfg.APIs[i].Merge], mergedData[0])
			mergedSample = true
		}
	}

	if !mergedSample {
		CreateMetricSets(mergedData, cfg, i)
	}
}

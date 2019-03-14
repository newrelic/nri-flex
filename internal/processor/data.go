package processor

import (
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
)

// runDataHandler handles the data received for processing
func runDataHandler(dataSets []interface{}, samplesToMerge *map[string][]interface{}, i int, yml *load.Config) {
	for _, dataSet := range dataSets {
		switch dataSet := dataSet.(type) {
		case map[string]interface{}:
			ds := dataSet
			processDataSet(&ds, samplesToMerge, i, yml)
		case []interface{}:
			nextDataSets := dataSet
			runDataHandler(nextDataSets, samplesToMerge, i, yml)
		default:
			logger.Flex("debug", nil, "not sure what to do with "+yml.Name, false)
		}
	}
}

// processDataSet performs the core flattening on the map[string]interface then executes createMetricSets finally
func processDataSet(dataSet *map[string]interface{}, samplesToMerge *map[string][]interface{}, i int, yml *load.Config) {
	ds := (*dataSet)

	FindStartKey(yml.APIs[i].StartKey, &ds)      // start at a later part in the received data
	StripKeys(yml.APIs[i].StripKeys, &ds)        // remove before flattening
	RunLazyFlatten(yml.APIs[i].LazyFlatten, &ds) // perform lazy flatten if needed
	flattenedData := FlattenData(ds, map[string]interface{}{}, "", yml.APIs[i].SampleKeys)

	// also strip from flattened data
	for _, stripKey := range yml.APIs[i].StripKeys {
		delete(flattenedData, stripKey)
		delete(flattenedData, stripKey+"Samples")
	}

	mergedData := FinalMerge(flattenedData)
	mergedSample := false

	if len(mergedData) == 1 {
		if yml.APIs[i].Merge != "" {
			mergedData[0].(map[string]interface{})["_sampleNo"] = i
			(*samplesToMerge)[yml.APIs[i].Merge] = append((*samplesToMerge)[yml.APIs[i].Merge], mergedData[0])
			mergedSample = true
		}
	}

	if !mergedSample {
		CreateMetricSets(mergedData, yml, i)
	}
}

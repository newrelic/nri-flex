/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"fmt"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

// FlattenData flatten an interface
func FlattenData(unknown interface{}, data map[string]interface{}, key string, sampleKeys map[string]string, api *load.API) map[string]interface{} {
	switch unknown := unknown.(type) {
	case []interface{}:
		dataSamples := []interface{}{}
		dataSamples = append(dataSamples, unknown...)
		key = checkPluralSlice(key)
		data[key+"FlexSamples"] = dataSamples
	case map[string]interface{}:
		if api.SplitObjects { // split objects can only be used once
			dataSamples := splitObjects(&unknown, api)
			FlattenData(dataSamples, data, key, sampleKeys, api)
		} else {
			for loopKey := range unknown {
				finalKey := loopKey
				if key != "" {
					finalKey = key + "." + loopKey
				}

				// Check sample keys and convert to samples using ">" as the sample key splitter if defined
				// Knowing the sampleKey itself isn't really needed as these get turned into samples
				for _, sampleVal := range sampleKeys {
					keys := strings.Split(sampleVal, ">")
					flexSamples := []interface{}{}
					// if one of the keys == the loopKey we know to create samples
					if len(keys) > 0 && keys[0] == loopKey {
						switch unknown[loopKey].(type) {
						case map[string]interface{}:
							dataSamples := unknown[loopKey].(map[string]interface{})
							for dataSampleKey, dataSample := range dataSamples {
								newSample := dataSample.(map[string]interface{})
								newSample[keys[1]] = dataSampleKey
								flexSamples = append(flexSamples, FlattenData(newSample, map[string]interface{}{}, "", sampleKeys, api))
							}
							unknown[loopKey] = flexSamples
						}
					}
				}

				FlattenData(unknown[loopKey], data, finalKey, sampleKeys, api)
			}
		}
	default:
		data[key] = unknown
	}

	for dataKey := range data {
		// separately flatten the flex samples, adding them back into the slice with a new key
		// & removing the old from data thus a replace
		if strings.Contains(dataKey, "FlexSamples") {
			strippedDataKey, newSamples := processFlexSamples(dataKey, data[dataKey].([]interface{}), sampleKeys, api)
			data[strippedDataKey] = newSamples
			delete(data, dataKey)
		}
	}

	return data
}

// FinalMerge Perform final data merging
// Separates detected samples and already flattened attributes
func FinalMerge(data map[string]interface{}) []interface{} {
	finalAttributes := map[string]interface{}{}
	finalSampleSets := map[string]interface{}{}

	// store all flat attributes in final attributes
	// store detected samples in SampleSets
	for key := range data {
		if !strings.Contains(key, "Samples") {
			finalAttributes[key] = data[key]
		} else {
			finalSampleSets[key] = data[key]
		}
	}

	finalMergedSamples := []interface{}{}
	for sampleSet := range finalSampleSets {
		switch ss := finalSampleSets[sampleSet].(type) {
		case []interface{}:
			for _, sample := range ss {
				switch sample := sample.(type) {
				case map[string]interface{}:
					newSample := sample
					newSample["event_type"] = sampleSet
					for attribute := range finalAttributes {
						newSample[attribute] = finalAttributes[attribute]
					}
					finalMergedSamples = append(finalMergedSamples, newSample)
				default:
					logger.Flex("debug", nil, fmt.Sprintf("%v not sure what to do with this?", sample), false)
				}
			}
		case map[string]interface{}:
			newSample := ss
			newSample["event_type"] = sampleSet
			for attribute := range finalAttributes {
				newSample[attribute] = finalAttributes[attribute]
			}
			finalMergedSamples = append(finalMergedSamples, newSample)
		}
	}

	if len(finalMergedSamples) > 0 {
		return finalMergedSamples
	}

	finalMergedSamples = append(finalMergedSamples, finalAttributes)
	return finalMergedSamples
}

// ProcessSamplesToMerge used to merge multiple samples together
func ProcessSamplesToMerge(samplesToMerge *map[string][]interface{}, yml *load.Config) {
	for eventType, sampleSet := range *samplesToMerge {
		newSample := map[string]interface{}{}
		newSample["event_type"] = eventType
		for _, sample := range sampleSet {
			prefix := yml.APIs[sample.(map[string]interface{})["_sampleNo"].(int)].Prefix
			for k, v := range sample.(map[string]interface{}) {
				if k != "_sampleNo" {
					newSample[prefix+k] = v
				}
			}
		}
		CreateMetricSets([]interface{}{newSample}, yml, 0)
	}
}

// processFlexSamples Processes Flex detected samples
func processFlexSamples(dataKey string, dataSamples []interface{}, sampleKeys map[string]string, api *load.API) (string, []interface{}) {
	newSamples := []interface{}{}
	for _, sample := range dataSamples {
		sampleFlatten := FlattenData(sample, map[string]interface{}{}, "", sampleKeys, api)
		if sampleFlatten["valuesPrometheusSamples"] != nil {
			for _, prometheusSample := range sampleFlatten["valuesPrometheusSamples"].([]interface{}) {
				// this could be optimized
				newSample := FlattenData(sample, map[string]interface{}{}, "", sampleKeys, api)
				newSample["timestamp"] = int(prometheusSample.([]interface{})[0].(float64))
				newSample["value"] = prometheusSample.([]interface{})[1]
				delete(newSample, "valuesPrometheusSamples")
				newSamples = append(newSamples, newSample)
			}
		} else if sampleFlatten["valuePrometheusSamples"] != nil {
			newSample := FlattenData(sample, map[string]interface{}{}, "", sampleKeys, api)
			newSample["timestamp"] = int(sampleFlatten["valuePrometheusSamples"].([]interface{})[0].(float64))
			newSample["value"] = sampleFlatten["valuePrometheusSamples"].([]interface{})[1]
			delete(newSample, "valuePrometheusSamples")
			newSamples = append(newSamples, newSample)
		} else {
			newSamples = append(newSamples, sampleFlatten)
		}
	}
	strippedDataKey := strings.Replace(dataKey, "Flex", "", -1)
	return strippedDataKey, newSamples
}

// checkPluralSlice Checks if a key is plural to use for FlexSample naming
// An assumption is made that the slice key is plural
func checkPluralSlice(key string) string {
	if len(key) > 0 {
		if key[len(key)-1:] == "s" {
			return key[:len(key)-1]
		}
	}
	return key
}

// splitObjects splits a map string interface / object with nested objects
// this will drop and ignore and slices/arrays that could exist
func splitObjects(unknown *map[string]interface{}, api *load.API) []interface{} {
	dataSamples := []interface{}{}
	for loopKey := range *unknown {
		switch data := (*unknown)[loopKey].(type) {
		case map[string]interface{}:
			logger.Flex("debug", nil, fmt.Sprintf("splitting object %v", loopKey), false)
			data["split.id"] = loopKey
			dataSamples = append(dataSamples, data)
		}
	}
	(*api).SplitObjects = false // only allow this to be run once
	return dataSamples
}

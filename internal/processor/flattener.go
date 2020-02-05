/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
)

const _originalAPINo = "_originalAPINo"

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
					load.Logrus.Debug("processor: flattener - not sure what to do with this?")
					load.Logrus.Debug(sample)
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
// hren: this function is replaced by ProcessSamplesMergeJoin
// func ProcessSamplesToMerge(samplesToMerge *map[string][]interface{}, yml *load.Config) {
// 	for eventType, sampleSet := range *samplesToMerge {
// 		newSample := map[string]interface{}{}
// 		newSample["event_type"] = eventType
// 		for _, sample := range sampleSet {
// 			prefix := yml.APIs[sample.(map[string]interface{})[_sampleNo].(int)].Prefix
// 			for k, v := range sample.(map[string]interface{}) {
// 				if k != _sampleNo {
// 					newSample[prefix+k] = v
// 				}
// 			}
// 		}
// 		CreateMetricSets([]interface{}{newSample}, yml, 0)
// 	}
// }

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
			load.Logrus.Debug(fmt.Sprintf("processor: splitting object %v", loopKey))
			data["split.id"] = loopKey
			dataSamples = append(dataSamples, data)
		}
	}
	(*api).SplitObjects = false // only allow this to be run once
	return dataSamples
}

// ProcessSamplesMergeJoin used to merge/join multiple samples together hren
func ProcessSamplesMergeJoin(samplesToMerge *load.SamplesToMerge, yml *load.Config) {
	for eventType, sampleSet := range samplesToMerge.Data {
		newSample := map[string]interface{}{}
		newSample["event_type"] = eventType
		primaryEvent := -1
		mergeEvent := 0
		primaryEventJoinKey := []string{""}
		primaryEventJoinWith := map[int]string{}
		primaryEventSamples := map[int]interface{}{}
		secondaryEventSamples := map[string]interface{}{}
		mergeOperation := false

		for rownum, sample := range sampleSet {
			originalAPINo := sample.(map[string]interface{})[_originalAPINo].(int)
			joinKey := yml.APIs[sample.(map[string]interface{})[_originalAPINo].(int)].JoinKey
			// merge or join operation
			if joinKey != "" {
				if primaryEvent == -1 {
					primaryEvent = originalAPINo
					primaryEventJoinKey = strings.Split(joinKey, ",")
					sort.Strings(primaryEventJoinKey)
				}
				if primaryEvent == originalAPINo {
					// key track of samples for primary event
					primaryEventSamples[rownum] = sample
				} else {
					// keep track of what events and the joinkey to be joint
					primaryEventJoinWith[originalAPINo] = joinKey
					// keep joining events in a map using sampleNo+joinkey value as key for later lookup
					for k, v := range sample.(map[string]interface{}) {
						if k == joinKey {
							value := ""
							switch (v).(type) {
							case float32, float64:
								// For float numbers, use decimal point format instead of scientific notation (e.g. 2026112.000000 vs 2.026112e+06 )
								// to allow the parser to process the original float number 2026112.000000 rather than 2.026112e+06
								value = fmt.Sprintf("%f", v)
							default:
								value = fmt.Sprintf("%v", v)
							}
							secondaryEventSamples[strconv.Itoa(originalAPINo)+value] = sample
						}
					}
				}
			} else {
				// doing merge processing if joinkey is empty
				mergeOperation = true
				mergeEvent = originalAPINo
				prefix := yml.APIs[sample.(map[string]interface{})[_originalAPINo].(int)].Prefix
				for k, v := range sample.(map[string]interface{}) {
					if k != _originalAPINo {
						newSample[prefix+k] = v
					}
				}
			}

		}

		if mergeOperation {
			CreateMetricSets([]interface{}{newSample}, yml, mergeEvent, false, nil, mergeEvent)
		}

		// if primary join event has samples
		for _, sample := range primaryEventSamples {
			newJoinSample := map[string]interface{}{}
			newJoinSample["event_type"] = eventType
			prefix := yml.APIs[sample.(map[string]interface{})[_originalAPINo].(int)].Prefix
			for k, v := range sample.(map[string]interface{}) {
				if k != _originalAPINo {
					newJoinSample[prefix+k] = v
					value := ""
					switch (v).(type) {
					case float32, float64:
						// For float numbers, use decimal point format instead of scientific notation (e.g. 2026112.000000 vs 2.026112e+06 )
						// to allow the parser to process the original float number 2026112.000000 rather than 2.026112e+06
						value = fmt.Sprintf("%f", v)
					default:
						value = fmt.Sprintf("%v", v)
					}
					if contains(primaryEventJoinKey, k) {
						// if the key is in the the primary join key list
						// then will lookup the secondary events for the key, if found, then join the keys from secondary events
						for kk := range primaryEventJoinWith {
							// for kk, vv := range primaryEventJoinWith {
							// we built the map with secondary event "_sampleNo" and join key earlier.
							// the map for the secondary event sample assumes unique lookup key per event
							if val, found := secondaryEventSamples[strconv.Itoa(kk)+value]; found {
								// use secondary's prefix to overwrite primary's
								prefix := yml.APIs[val.(map[string]interface{})[_originalAPINo].(int)].Prefix
								for kkk, vvv := range val.(map[string]interface{}) {
									if kkk != _originalAPINo {
										newJoinSample[prefix+kkk] = vvv
									}
								}

							}
						}
					}
				}
			}
			CreateMetricSets([]interface{}{newJoinSample}, yml, primaryEvent, false, nil, primaryEvent)
		}

	}
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

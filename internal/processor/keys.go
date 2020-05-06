/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
)

// RunKeyConversion handles to lower and snake to camel case for keys
func RunKeyConversion(key *string, api load.API, v interface{}, SkipProcessing *[]string) {
	if api.ToLower {
		*key = strings.ToLower(*key)
	}
	if api.ConvertSpace != "" {
		*key = strings.Replace(*key, " ", api.ConvertSpace, -1)
	}
	if api.SnakeToCamel {
		formatter.SnakeCaseToCamelCase(key)
	}
}

// RunKeepKeys will remove the key if is not defined in keep_keys.
func RunKeepKeys(keepKeys []string, key *string, currentSample *map[string]interface{}) {
	if len(keepKeys) > 0 {
		foundKey := false
		for _, keepKey := range keepKeys {
			if formatter.KvFinder("regex", *key, keepKey) {
				foundKey = true
				break
			}
		}
		if !foundKey {
			delete(*currentSample, *key)
		}
	}
}

// RunKeyRemover Remove unwanted keys with regex
func RunKeyRemover(currentSample *map[string]interface{}, removeKeys []string) {
	for _, removeKey := range removeKeys {
		for key := range *currentSample {
			// ignore case of key to remove
			if formatter.KvFinder("regex", key, "(?i)"+removeKey) {
				delete(*currentSample, key)
			}
		}
	}
}

// RunKeyRenamer find keys with regex, and replace the value
func RunKeyRenamer(renameKeys map[string]string, key *string) {
	for renameKey, renameVal := range renameKeys {
		// TODO: Should this first try matching as a plain string and after that try compile it as regex?
		validateKey := regexp.MustCompile(renameKey)
		*key = validateKey.ReplaceAllString(*key, renameVal)
	}
}

// StripKeys strip defined keys out
func StripKeys(ds *map[string]interface{}, stripKeys []string) {
	for _, stripKey := range stripKeys {
		delete(*ds, stripKey)
		if strings.Contains(stripKey, ">") {
			stripSplit := strings.Split(stripKey, ">")
			if len(stripSplit) == 2 {
				if (*ds)[stripSplit[0]] != nil {
					switch (*ds)[stripSplit[0]].(type) {
					case map[string]interface{}:
						delete((*ds)[stripSplit[0]].(map[string]interface{}), stripSplit[1])
					case []interface{}:
						for i := range (*ds)[stripSplit[0]].([]interface{}) {
							switch (*ds)[stripSplit[0]].([]interface{})[i].(type) {
							case map[string]interface{}:
								delete((*ds)[stripSplit[0]].([]interface{})[i].(map[string]interface{}), stripSplit[1])
							}
						}
					}
				}
			}
		}
	}
}

// FindStartKey start at a different section of a payload
func FindStartKey(mainDataset *map[string]interface{}, startKeys []string, inheritAttributes bool) {
	parentAttributes := map[string]interface{}{}
	for level, startKey := range startKeys {
		startSplit := strings.Split(startKey, ">")
		if len(startSplit) == 2 {
			if (*mainDataset)[startSplit[0]] != nil {
				storeParentAttributes(*mainDataset, parentAttributes, startKey, level, inheritAttributes)
				switch mainDs := (*mainDataset)[startSplit[0]].(type) {
				case []interface{}:
					var nestedSlices []interface{}
					for i, nested := range mainDs {
						switch sample := nested.(type) {
						case map[string]interface{}:
							if sample[startSplit[1]] != nil {
								switch nestedSample := sample[startSplit[1]].(type) {
								case map[string]interface{}:
									nestedSample["flexSliceIndex"] = i
									nestedSlices = append(nestedSlices, nestedSample)
								case []interface{}:
									nestedSlices = append(nestedSlices, applyIndexes(i, nestedSample)...)
								}
							}
						}
					}
					applyParentAttributes(nil, nestedSlices, parentAttributes)
					*mainDataset = map[string]interface{}{startSplit[1]: nestedSlices}
				}
			}
		} else if len(startSplit) == 1 {
			if (*mainDataset)[startKey] != nil {
				storeParentAttributes(*mainDataset, parentAttributes, startKey, level, inheritAttributes)
				switch mainDs := (*mainDataset)[startKey].(type) {
				case map[string]interface{}:
					*mainDataset = mainDs
					if len(startKeys)-1 == level {
						applyParentAttributes(*mainDataset, nil, parentAttributes)
					}
				case []interface{}:
					applyParentAttributes(nil, mainDs, parentAttributes)
					*mainDataset = map[string]interface{}{startKey: mainDs}
				}
			}
		}
	}
}

func applyIndexes(index int, slices []interface{}) []interface{} {
	newSlices := []interface{}{}
	for _, sample := range slices {
		switch sampleData := sample.(type) {
		case map[string]interface{}:
			sampleData["flexSliceIndex"] = index
			newSlices = append(newSlices, sampleData)
		default:
			newSlices = append(newSlices, sample)
		}
	}
	return newSlices
}

func storeParentAttributes(mainDataset map[string]interface{}, parentAttributes map[string]interface{}, startKey string, level int, inheritAttributes bool) {
	if inheritAttributes {
		startSplit := strings.Split(startKey, ">")
		for key, val := range mainDataset {
			if key != startKey {
				switch valueData := val.(type) {
				case map[string]interface{}, []interface{}:
					if len(startSplit) == 2 && key == startSplit[0] {
						// lazy flatten the slices and maps from the highest level
						for mapKey, mapVal := range flattenSlicesAndMaps(val) {
							switch data := mapVal.(type) {
							case map[string]interface{}:
								for innerMapKey, innerMapVal := range data {
									// avoid duplicates from the flattened data
									if !strings.Contains(innerMapKey, startSplit[1]) {
										parentAttributes[fmt.Sprintf("parent.%d.%v.%v", level, mapKey, innerMapKey)] = fmt.Sprintf("%v", innerMapVal)
									}
								}
							default:
								// avoid duplicates from the flattened data
								if !strings.Contains(mapKey, startSplit[1]) {
									parentAttributes[fmt.Sprintf("parent.%d.%v", level, mapKey)] = fmt.Sprintf("%v", mapVal)
								}
							}
						}
					}
				default:
					value := fmt.Sprintf("%v", valueData)
					parentAttributes[fmt.Sprintf("parent.%d.%v", level, key)] = value
				}
			}
		}
	}
}

func applyParentAttributes(mainDataset map[string]interface{}, datasets []interface{}, parentAttributes map[string]interface{}) {
	if mainDataset != nil {
		for key, val := range parentAttributes {
			mainDataset[key] = val
		}
	} else if len(datasets) > 0 {
		for _, dataset := range datasets {
			switch switchDs := dataset.(type) {
			case map[string]interface{}:
				for key, val := range parentAttributes {
					// check if this is a nested parent, and only apply if the index matches
					matches := formatter.RegMatch(key, "parent\\.(\\d+)\\.(\\d+)\\.(.+)")
					if len(matches) > 1 {
						sliceIndex := -1
						if switchDs["flexSliceIndex"] != nil {
							sliceIndex = switchDs["flexSliceIndex"].(int)
						}
						matchIndex2, err := strconv.Atoi(matches[1])
						if err == nil {
							// no need to add the second index into the key as we've unpacked at the matched level
							if sliceIndex == matchIndex2 {
								switchDs[fmt.Sprintf("parent.%v.%v", matches[0], matches[2])] = val
							}
						}
					} else {
						switchDs[key] = val
					}
				}
				delete(switchDs, "flexSliceIndex")
			}
		}
	}
}

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"fmt"
	"regexp"
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

// RunKeepKeys Removes all other keys/attributes and keep only those defined in keep_keys
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
			if formatter.KvFinder("regex", key, removeKey) {
				delete(*currentSample, key)
				break
			}
		}
	}
}

// RunKeyRenamer find keys with regex, and replace the value
func RunKeyRenamer(renameKeys map[string]string, key *string, originalKey string) {
	for renameKey, renameVal := range renameKeys {
		validateKey := regexp.MustCompile(renameKey)
		matches := validateKey.FindAllString(*key, -1)
		for _, match := range matches {
			*key = strings.Replace(*key, match, renameVal, -1)
		}
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
					nestedSlices := []interface{}{}
					for _, nested := range mainDs {
						switch sample := nested.(type) {
						case map[string]interface{}:
							if sample[startSplit[1]] != nil {
								switch nestedSample := sample[startSplit[1]].(type) {
								case map[string]interface{}:
									nestedSlices = append(nestedSlices, nestedSample)
								case []interface{}:
									nestedSlices = append(nestedSlices, nestedSample[0])
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

func storeParentAttributes(mainDataset map[string]interface{}, parentAttributes map[string]interface{}, startKey string, level int, inheritAttributes bool) {
	if inheritAttributes {
		for key, val := range mainDataset {
			if key != startKey {
				value := fmt.Sprintf("%v", val)
				if !strings.Contains(value, "map[") {
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
					switchDs[key] = val
				}
			}
		}
	}
}

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/newrelic/nri-flex/internal/inputs"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// FetchData fetches data from various inputs
// Also handles paginated responses for HTTP requests (tested against NR APIs)
func FetchData(apiNo int, yml *load.Config, samplesToMerge *load.SamplesToMerge) []interface{} {
	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
	}).Debug("fetch: collect data")

	api := yml.APIs[apiNo]
	file := yml.APIs[apiNo].File
	reqURL := api.URL

	doLoop := true
	dataStore := []interface{}{}

	continueProcessing := FetchLookups(yml, apiNo, samplesToMerge)

	if continueProcessing {
		if file != "" {
			inputs.RunFile(&dataStore, yml, apiNo)
		} else if api.Cache != "" {
			if yml.Datastore[api.Cache] != nil {
				dataStore = yml.Datastore[api.Cache]
			}
		} else if api.Ingest {
			if yml.Datastore["IngestData"] != nil {
				dataStore = yml.Datastore["IngestData"]
			}
		} else if len(api.Commands) > 0 && api.Database == "" && api.DbConn == "" {
			inputs.RunCommands(&dataStore, yml, apiNo)
		} else if reqURL != "" {
			inputs.RunHTTP(&dataStore, &doLoop, yml, api, &reqURL)
		} else if api.Database != "" && api.DbConn != "" {
			inputs.ProcessQueries(&dataStore, yml, apiNo)
		} else if api.Scp.Host != "" {
			inputs.RunScpWithTimeout(&dataStore, yml, api, apiNo)
		}
	}

	// cache output into datastore for later use
	if len(dataStore) > 0 {
		if api.URL != "" {
			if yml.Datastore == nil {
				yml.Datastore = map[string][]interface{}{}
			}
			yml.Datastore[api.URL] = dataStore
		} else if len(api.Commands) > 0 && api.Database == "" && api.DbConn == "" && api.Name != "" {
			if yml.Datastore == nil {
				yml.Datastore = map[string][]interface{}{}
			}
			yml.Datastore[api.Name] = dataStore
		} else if api.File != "" {
			if yml.Datastore == nil {
				yml.Datastore = map[string][]interface{}{}
			}
			yml.Datastore[api.File] = dataStore
		}
	}

	return dataStore
}

// FetchLookups x
func FetchLookups(cfg *load.Config, apiNo int, samplesToMerge *load.SamplesToMerge) bool {
	tmpCfgBytes, err := yaml.Marshal(&cfg.APIs[apiNo])

	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"name": cfg.Name,
			"err":  err,
		}).Error("fetch: lookup processor marshal failed")
	} else {
		tmpCfgStr := string(tmpCfgBytes)
		lookupsFound := regexp.MustCompile(`\${lookup:.*?}`).FindAllString(tmpCfgStr, -1)

		// if no lookups, do not continue running the processor
		if len(lookupsFound) == 0 {
			return true
		}

		// determine each unique lookup found
		lookupDimensions := []string{}
		for _, lookupVar := range lookupsFound {
			lookupDimensionFound := false
			// eg. ${lookup:consumers} -> consumers
			currentLookupDimension := strings.TrimSuffix(strings.Split(lookupVar, "${lookup:")[1], "}")

			for _, lookupDimension := range lookupDimensions {
				if currentLookupDimension == lookupDimension {
					lookupDimensionFound = true
					break
				}
			}

			// only if not found then append to ensure the slice is unique
			if !lookupDimensionFound {
				lookupDimensions = append(lookupDimensions, currentLookupDimension)
			}
		}

		load.Logrus.WithFields(logrus.Fields{
			"name": cfg.Name,
		}).Debug(fmt.Sprintf("fetch: unique lookups found in api %v", (lookupDimensions)))

		sliceIndexes := []int{}
		sliceKeys := []string{}
		sliceLookups := [][]string{}

		// init lookups
		for key, values := range cfg.LookupStore {
			// only create lookups for the found dimensions
			for _, dimKey := range lookupDimensions {
				if key == dimKey {
					sliceIndexes = append(sliceIndexes, 0)
					sliceKeys = append(sliceKeys, key)
					valueArray := []string{}
					for a := range values {
						valueArray = append(valueArray, a)
					}
					sliceLookups = append(sliceLookups, valueArray)
					break
				}
			}
		}

		loopNo := -1
		combinations := [][]string{}
		if len(sliceLookups) > 0 {
			loopLookups(loopNo, sliceIndexes, sliceLookups, &combinations)
		}

		load.Logrus.WithFields(logrus.Fields{
			"name": cfg.Name,
		}).Debug(fmt.Sprintf("fetch: combinations found %v", (combinations)))

		newAPIs := []string{}
		for _, combo := range combinations {
			tmpConfigWithLookupReplace := tmpCfgStr
			if len(combo) == len(sliceKeys) {
				for i, key := range sliceKeys {
					tmpConfigWithLookupReplace = strings.ReplaceAll(tmpConfigWithLookupReplace, fmt.Sprintf("${lookup:%v}", key), combo[i])
				}
				newAPIs = append(newAPIs, tmpConfigWithLookupReplace)
			} else {
				load.Logrus.WithFields(logrus.Fields{
					"name": cfg.Name,
				}).Debug("fetch: invalid lookup, missing a replace")
			}
		}

		lookupConfig := load.Config{
			Name:             cfg.Name,
			Global:           cfg.Global,
			FileName:         cfg.FileName,
			Datastore:        cfg.Datastore,
			LookupStore:      cfg.LookupStore,
			VariableStore:    cfg.VariableStore,
			CustomAttributes: cfg.CustomAttributes,
		}

		for _, newAPI := range newAPIs {
			API := load.API{}
			err := yaml.Unmarshal([]byte(newAPI), &API)
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"name": cfg.Name,
					"err":  err,
				}).Error("fetch: failed to unmarshal lookup config")
			} else {
				lookupConfig.APIs = append(lookupConfig.APIs, API)
			}
		}
		// Run(lookupConfig)
		// Please note:
		//          When in RunAsync/run_async mode, we will disable StoreLookups and VariableLookups due to potential concurrent map write.
		//          We will address this in the future if required. These two functions are probably not necessary for this use case.
		if cfg.APIs[apiNo].RunAsync {
			RunAsync(lookupConfig, samplesToMerge, apiNo)
		} else {
			RunSync(lookupConfig, samplesToMerge, apiNo)
		}
		return false
	}

	return true
}

func loopLookups(loopNo int, sliceIndexes []int, sliceLookups [][]string, combinations *[][]string) {
	loopNo++
	for i := range sliceLookups[loopNo] {
		// track the index of each loop
		(sliceIndexes)[loopNo] = i

		// this indicates we are in the inner most loop, else do another loop
		if loopNo+1 == len(sliceLookups) {
			keys := []string{}
			for x := 0; x <= loopNo; x++ {
				keys = append(keys, sliceLookups[x][sliceIndexes[x]])
			}
			*combinations = append(*combinations, keys)
		} else {
			loopLookups(loopNo, sliceIndexes, sliceLookups, combinations)
		}
	}
}

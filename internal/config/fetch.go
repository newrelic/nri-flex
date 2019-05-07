package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/newrelic/nri-flex/internal/inputs"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	yaml "gopkg.in/yaml.v2"
)

// FetchData fetches data from various inputs
// Also handles paginated responses for HTTP requests (tested against NR APIs)
func FetchData(apiNo int, yml *load.Config) []interface{} {
	logger.Flex("debug", nil, fmt.Sprintf("fetching data for %v", yml.Name), false)

	api := yml.APIs[apiNo]
	file := yml.APIs[apiNo].File
	reqURL := api.URL

	doLoop := true
	dataStore := []interface{}{}

	continueProcessing := FetchLookups(yml, apiNo)

	if continueProcessing {
		if file != "" {
			fileData, err := ioutil.ReadFile(file)
			if err != nil {
				logger.Flex("error", err, "unable to read file: "+file, false)
			} else {
				newBody := strings.Replace(string(fileData), " ", "", -1)
				var f interface{}
				err := json.Unmarshal([]byte(newBody), &f)
				if err != nil {
					logger.Flex("error", err, "failed to unmarshal", false)
				} else {
					dataStore = append(dataStore, f)
				}
			}
		} else if api.Cache != "" {
			if yml.Datastore[api.Cache] != nil {
				// load.StoreAppend(yml.Datastore[api.Cache])
				dataStore = yml.Datastore[api.Cache]
			}
		} else if len(api.Commands) > 0 && api.Database == "" && api.DbConn == "" {
			inputs.RunCommands(yml, apiNo, &dataStore)
		} else if reqURL != "" {
			inputs.RunHTTP(&dataStore, &doLoop, yml, api, &reqURL)
		} else if api.Database != "" && api.DbConn != "" {
			inputs.ProcessQueries(&dataStore, api)
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
func FetchLookups(cfg *load.Config, i int) bool {
	tmpCfgBytes, err := yaml.Marshal(&cfg.APIs[i])

	if err != nil {
		logger.Flex("error", err, "lookup processor marshal failed", false)
	} else {
		tmpCfgStr := string(tmpCfgBytes)

		// if no lookups, do not continue running the processor
		if !strings.Contains(tmpCfgStr, "${lookup:") {
			return true
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

		replaceOccured := false
		newAPIs := []string{}
		lookupIndex := 0
		for lookup, lookupKeys := range cfg.LookupStore {
			for z, key := range lookupKeys {
				if lookupIndex == 0 {
					newAPIs = append(newAPIs, tmpCfgStr)
				}
				newAPIs[z] = strings.Replace(newAPIs[z], ("${lookup:" + lookup + "}"), key, -1)
				replaceOccured = true
			}
			lookupIndex++
		}

		if replaceOccured {
			for _, newAPI := range newAPIs {
				API := load.API{}
				err := yaml.Unmarshal([]byte(newAPI), &API)
				if err != nil {
					logger.Flex("error", err, "failed to unmarshal lookup config", false)
				} else {
					lookupConfig.APIs = append(lookupConfig.APIs, API)
				}

			}
			Run(lookupConfig)
			return false
		}
	}

	return true
}

package processor

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	parser "github.com/newrelic/nri-flex/internal/parsers"
)

// fetchData fetches data and handles paginated responses
func fetchData(i int, yml *load.Config) []interface{} {
	continueProcessing := true

	api := yml.APIs[i]
	file := yml.APIs[i].File
	reqURL := api.URL

	if strings.Contains(reqURL, "${lookup:") {
		runLookupProcessor(reqURL, yml, i)
		// stop processing normally, and restart from fetchData with newly loaded urls to fetch
		continueProcessing = false
	}

	dataStore := []interface{}{}
	doLoop := true

	if continueProcessing {
		if file != "" {
			fileData, err := ioutil.ReadFile(file)
			if err != nil {
				logger.Flex("debug", err, "unable to read file: "+file, false)
			} else {
				newBody := strings.Replace(string(fileData), " ", "", -1)
				var f interface{}
				err := json.Unmarshal([]byte(newBody), &f)
				if err != nil {
					logger.Flex("debug", err, "failed to unmarshal", false)
				} else {
					dataStore = append(dataStore, f)
				}
			}
		} else if api.Cache != "" {
			if yml.Datastore[api.Cache] != nil {
				dataStore = yml.Datastore[api.Cache]
			}
		} else if len(api.Commands) > 0 && api.Database == "" && api.DbConn == "" {
			parser.RunCommands(*yml, api, &dataStore)
		} else if reqURL != "" {
			parser.RunHTTP(&doLoop, yml, api, &reqURL, &dataStore)
		} else if api.Database != "" && api.DbConn != "" {
			parser.ProcessQueries(api, &dataStore)
		}
	}

	// cache output into datastore for later use
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

	return dataStore
}

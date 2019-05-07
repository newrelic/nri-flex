package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	"github.com/newrelic/nri-flex/internal/processor"

	yaml "gopkg.in/yaml.v2"
)

// LoadFiles Loads Flex config files
func LoadFiles(configs *[]load.Config, files []os.FileInfo, path string) {
	for _, f := range files {
		b, err := ioutil.ReadFile(path + f.Name())
		if err != nil {
			logger.Flex("debug", err, "unable to readfile", false)
			continue
		}
		if !strings.Contains(f.Name(), "yml") && !strings.Contains(f.Name(), "yaml") {
			continue
		}
		ymlStr := string(b)
		SubEnvVariables(&ymlStr)
		SubTimestamps(&ymlStr)
		config, err := ReadYML(ymlStr)
		config.FileName = f.Name()
		if err != nil {
			logger.Flex("error", err, "unable to read yml", false)
			continue
		}
		if config.Name == "" {
			logger.Flex("error", fmt.Errorf("config file %v requires a name", f.Name()), "", false)
			continue
		}

		// if lookup files exist we need to potentially create multiple config files
		if config.LookupFile != "" {
			SubLookupFileData(configs, config)
		} else {
			*configs = append(*configs, config)
		}

	}
}

// ReadYML Unmarshals yml files
func ReadYML(yml string) (load.Config, error) {
	c := load.Config{}
	err := yaml.Unmarshal([]byte(yml), &c)
	if err != nil {
		return load.Config{}, err
	}
	return c, nil
}

// Run Action each config file
func Run(yml load.Config) {
	samplesToMerge := map[string][]interface{}{}
	for i := range yml.APIs {
		RunVariableProcessor(i, &yml)
		FetchData(i, &yml)
		// dataSets := FetchData(i, &yml)
		processor.RunDataHandler(load.Store.Data, &samplesToMerge, i, &yml)
		// processor.RunDataHandler(dataSets, &samplesToMerge, i, &yml)
	}
	processor.ProcessSamplesToMerge(&samplesToMerge, &yml)
}

// RunFiles Processes yml files
func RunFiles(configs *[]load.Config) {
	logger.Flex("debug", nil, fmt.Sprintf("starting to process %d configs", len(*configs)), false)
	var wg sync.WaitGroup
	wg.Add(len(*configs))
	for _, cfg := range *configs {
		go func(cfg load.Config) {
			defer wg.Done()
			logger.Flex("debug", nil, fmt.Sprintf("running config: %v", cfg.Name), false)
			Run(cfg)
			load.StatusCounterIncrement("ConfigsProcessed")
		}(cfg)
	}
	wg.Wait()
	logger.Flex("debug", nil, fmt.Sprintf("completed processing %d configs", len(*configs)), false)
}

// RunVariableProcessor substitute store variables into specific parts of config files
func RunVariableProcessor(i int, cfg *load.Config) {
	// don't use variable processor if nothing exists in variable store
	if len((*cfg).VariableStore) > 0 {
		logger.Flex("debug", nil, fmt.Sprintf("running variable processor %d items in store", len((*cfg).VariableStore)), false)
		// to simplify replacement, convert to string, and convert back later
		tmpCfgBytes, err := yaml.Marshal(&cfg)
		if err != nil {
			logger.Flex("error", err, "variable processor marshal failed", false)
		} else {
			tmpCfgStr := string(tmpCfgBytes)
			variableReplaces := regexp.MustCompile(`\${var:.*?}`).FindAllString(tmpCfgStr, -1)
			replaceOccured := false
			for _, variableReplace := range variableReplaces {
				variableKey := strings.TrimSuffix(strings.Split(variableReplace, "${var:")[1], "}") // eg. "channel"
				if cfg.VariableStore[variableKey] != "" {
					tmpCfgStr = strings.Replace(tmpCfgStr, variableReplace, cfg.VariableStore[variableKey], -1)
					replaceOccured = true
				}
			}
			// if replace occurred convert string to config yaml and reload
			if replaceOccured {
				newCfg, err := ReadYML(tmpCfgStr)
				if err != nil {
					logger.Flex("error", err, "variable processor unmarshal failed", false)
				} else {
					*cfg = newCfg
				}
			}
		}
	}
}

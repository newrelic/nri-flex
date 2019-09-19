/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/processor"
	"github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

// LoadFiles Loads Flex config files
func LoadFiles(configs *[]load.Config, files []os.FileInfo, path string) {
	for _, f := range files {
		filePath := path + f.Name()
		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			if strings.Contains(err.Error(), "is a directory") { // if it is a directory then recurse
				if !strings.Contains(filePath, ".git") && !strings.Contains(filePath, "nr-integrations") { // do not recurse through .git or nr-integrations folder
					load.Logrus.WithFields(logrus.Fields{
						"path": filePath,
					}).Debug("config: checking nested configs")
					nextPath := filePath + "/"
					files, err = ioutil.ReadDir(nextPath)
					if err != nil {
						load.Logrus.WithFields(logrus.Fields{
							"path": nextPath,
						}).Debug("config: failed to read")
					} else {
						LoadFiles(configs, files, nextPath)
					}
				}
			} else {
				load.Logrus.WithFields(logrus.Fields{
					"file": filePath,
					"err":  err,
				}).Debug("config: failed to read")
			}
			continue
		}
		// not done earlier as there could be a directory first
		if !strings.HasSuffix(f.Name(), "yml") && !strings.HasSuffix(f.Name(), "yaml") {
			continue
		}
		ymlStr := string(b)
		SubEnvVariables(&ymlStr)
		SubTimestamps(&ymlStr)
		config, err := ReadYML(ymlStr)
		config.FileName = f.Name()
		config.FilePath = path

		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"file": filePath,
				"err":  err,
			}).Error("config: failed to unmarshal yaml")
			continue
		}
		if config.Name == "" {
			load.Logrus.WithFields(logrus.Fields{
				"file": filePath,
			}).Error("config: flexConfig requires name")
			continue
		}

		checkIngestConfigs(&config)

		// if lookup files exist we need to potentially create multiple config files
		if config.LookupFile != "" {
			SubLookupFileData(configs, config)
		} else {
			*configs = append(*configs, config)
		}

	}
}

func checkIngestConfigs(config *load.Config) {
	if (*config).FileName == "flex-lambda-ingest.yml" && load.LambdaName != "" {
		if load.IngestData != nil {
			(*config).Datastore = map[string][]interface{}{}
			(*config).Datastore["IngestData"] = []interface{}{load.IngestData}
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
	load.Logrus.Debug(fmt.Sprintf("config: processing %d apis in %v", len(yml.APIs), yml.Name))

	// load secrets
	loadSecrets(&yml)

	// intentionally handled synchronously
	for i := range yml.APIs {
		RunVariableProcessor(i, &yml)
		dataSets := FetchData(i, &yml)
		processor.RunDataHandler(dataSets, &samplesToMerge, i, &yml)
	}

	load.Logrus.Debug(fmt.Sprintf("config: finished processing %d apis in %v", len(yml.APIs), yml.Name))
	processor.ProcessSamplesToMerge(&samplesToMerge, &yml)
}

// RunFiles Processes yml files
func RunFiles(configs *[]load.Config) {
	if load.Args.ProcessConfigsSync {
		for _, cfg := range *configs {
			if verifyConfig(cfg) {
				load.Logrus.WithFields(logrus.Fields{"name": cfg.Name}).Debug("config: running")
				Run(cfg)
				load.StatusCounterIncrement("ConfigsProcessed")
			}
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(len(*configs))
		for _, cfg := range *configs {
			go func(cfg load.Config) {
				defer wg.Done()
				if verifyConfig(cfg) {
					load.Logrus.WithFields(logrus.Fields{"name": cfg.Name}).Debug("config: running")
					Run(cfg)
					load.StatusCounterIncrement("ConfigsProcessed")
				}
			}(cfg)
		}
		wg.Wait()
	}
	load.Logrus.Info(fmt.Sprintf("flex: completed processing %d config(s)", load.StatusCounterRead("ConfigsProcessed")))
}

// verifyConfig ensure the config file doesn't have anything it should not run
func verifyConfig(cfg load.Config) bool {
	if strings.HasPrefix(cfg.FileName, "cd-") && !cfg.ContainerDiscovery.ReplaceComplete {
		return false
	}
	ymlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return false
	}
	ymlStr := string(ymlBytes)
	if strings.Contains(ymlStr, "${auto:host}") || strings.Contains(ymlStr, "${auto:port}") {
		return false
	}
	return true
}

// RunVariableProcessor substitute store variables into specific parts of config files
func RunVariableProcessor(i int, cfg *load.Config) {
	// don't use variable processor if nothing exists in variable store
	if len((*cfg).VariableStore) > 0 {
		load.Logrus.Debug(fmt.Sprintf("running variable processor %d items in store", len((*cfg).VariableStore)))
		// to simplify replacement, convert to string, and convert back later
		tmpCfgBytes, err := yaml.Marshal(&cfg)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{"err": err, "name": cfg.Name}).Error("config: variable processor marshal failed")
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
					load.Logrus.WithFields(logrus.Fields{"err": err, "name": cfg.Name}).Error("config: variable processor unmarshal failed")
				} else {
					*cfg = newCfg
				}
			}
		}
	}
}

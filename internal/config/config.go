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
		if f.IsDir() {
			recurseDirectory(filePath, configs)
			continue
		}
		// ignoring non-yaml files
		if !strings.HasSuffix(f.Name(), "yml") && !strings.HasSuffix(f.Name(), "yaml") {
			continue
		}
		if err := LoadFile(configs, f, path); err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"file": filePath,
				"err":  err,
			}).Error("can't load file")
		}
	}
}

// LoadFile loads a single Flex config file
func LoadFile(configs *[]load.Config, f os.FileInfo, path string) error {
	filePath := path + f.Name()

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	ymlStr := string(b)
	SubEnvVariables(&ymlStr)
	SubTimestamps(&ymlStr)
	config, err := ReadYML(ymlStr)
	if err != nil {
		return fmt.Errorf("config: failed to load config file, %v", err)
	}
	config.FileName = f.Name()
	config.FilePath = path

	if config.Name == "" {
		return fmt.Errorf("config: flexConfig requires name: %s", config.FilePath)
	}

	checkIngestConfigs(&config)

	// if lookup files exist we need to potentially create multiple config files
	if config.LookupFile != "" {
		err := SubLookupFileData(configs, config)
		if err != nil {
			load.Logrus.WithError(err).Error("config: failed to sub lookup file data")
		}
	} else {
		*configs = append(*configs, config)
	}
	return nil
}

func recurseDirectory(filePath string, configs *[]load.Config) {
	// do not recurse through .git or nr-integrations folder.
	if strings.Contains(filePath, ".git") ||
		strings.Contains(filePath, "nr-integrations") {
		return
	}

	load.Logrus.WithFields(logrus.Fields{
		"path": filePath,
	}).Debug("config: checking nested configs")

	nextPath := filePath + "/"
	files, err := ioutil.ReadDir(nextPath)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"path": nextPath,
		}).Debug("config: failed to read")
		return
	}
	LoadFiles(configs, files, nextPath)
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
	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: processing apis")

	// load secrets
	loadSecrets(&yml)

	// intentionally handled synchronously
	for i := range yml.APIs {
		if err := runVariableProcessor(&yml); err != nil {
			load.Logrus.WithError(err).Error("config: variable processor error")
		}

		dataSets := FetchData(i, &yml)
		processor.RunDataHandler(dataSets, &samplesToMerge, i, &yml)
	}

	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: finished processing apis")

	// processor.ProcessSamplesToMerge(&samplesToMerge, &yml)
	// hren joinAndMerge processing - replacing processor.ProcessSamplesToMerge
	processor.ProcessSamplesMergeJoin(&samplesToMerge, &yml)
}

// RunFiles Processes yml files
func RunFiles(configs *[]load.Config) {
	if load.Args.ProcessConfigsSync {
		for _, cfg := range *configs {
			if verifyConfig(cfg) {
				load.Logrus.WithFields(logrus.Fields{"name": cfg.Name}).Debug("config: running sync")
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
					load.Logrus.WithFields(logrus.Fields{"name": cfg.Name}).Debug("config: running async")
					Run(cfg)
					load.StatusCounterIncrement("ConfigsProcessed")
				}
			}(cfg)
		}
		wg.Wait()
	}

	load.Logrus.WithFields(logrus.Fields{
		"configs": load.StatusCounterRead("ConfigsProcessed"),
	}).Info("flex: completed processing")
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

// runVariableProcessor substitute store variables into specific parts of config files
func runVariableProcessor(cfg *load.Config) error {
	// don't use variable processor if nothing exists in variable store
	if len((*cfg).VariableStore) == 0 {
		return nil
	}

	load.Logrus.Debugf("running variable processor %d items in store", len((*cfg).VariableStore))
	// to simplify replacement, convert to string, and convert back later
	tmpCfgBytes, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("config %s: variable processor marshal failed, error: %v", cfg.Name, err)
	}

	tmpCfgStr := string(tmpCfgBytes)
	variableReplaces := regexp.MustCompile(`\${var:.*?}`).FindAllString(tmpCfgStr, -1)
	replaceOccurred := false
	for _, variableReplace := range variableReplaces {
		variableKey := strings.TrimSuffix(strings.Split(variableReplace, "${var:")[1], "}") // eg. "channel"
		if cfg.VariableStore[variableKey] != "" {
			tmpCfgStr = strings.Replace(tmpCfgStr, variableReplace, cfg.VariableStore[variableKey], -1)
			replaceOccurred = true
		}
	}
	// if replace occurred convert string to config yaml and reload
	if replaceOccurred {
		newCfg, err := ReadYML(tmpCfgStr)
		if err != nil {
			return fmt.Errorf("config %s: variable processor read yml failed, error: %v", cfg.Name, err)
		}
		*cfg = newCfg
	}
	return nil
}

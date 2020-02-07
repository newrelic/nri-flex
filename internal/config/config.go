/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"errors"
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
		return err
	}
	config.FileName = f.Name()
	config.FilePath = path

	if err != nil {
		return err
	}
	if config.Name == "" {
		return errors.New("config: flexConfig requires name")
	}

	checkIngestConfigs(&config)

	// if lookup files exist we need to potentially create multiple config files
	if config.LookupFile != "" {
		SubLookupFileData(configs, config)
	} else {
		*configs = append(*configs, config)
	}
	return nil
}

func recurseDirectory(filePath string, configs *[]load.Config) {
	if !strings.Contains(filePath, ".git") && !strings.Contains(filePath, "nr-integrations") { // do not recurse through .git or nr-integrations folder
		load.Logrus.WithFields(logrus.Fields{
			"path": filePath,
		}).Debug("config: checking nested configs")
		nextPath := filePath + "/"
		files, err := ioutil.ReadDir(nextPath)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"path": nextPath,
			}).Debug("config: failed to read")
		} else {
			LoadFiles(configs, files, nextPath)
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
	// samplesToMerge := map[string][]interface{}{}
	var samplesToMerge load.SamplesToMerge
	samplesToMerge.Data = map[string][]interface{}{}
	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: processing apis")

	// load secrets
	loadSecrets(&yml)

	// intentionally handled synchronously
	for i := range yml.APIs {
		RunVariableProcessor(i, &yml)
		dataSets := FetchData(i, &yml, &samplesToMerge)
		processor.RunDataHandler(dataSets, &samplesToMerge, i, &yml, i)
	}

	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: finished processing apis")

	// processor.ProcessSamplesToMerge(&samplesToMerge, &yml)
	// hren MergeAndJoin processing - replacing processor.ProcessSamplesToMerge
	processor.ProcessSamplesMergeJoin(&samplesToMerge, &yml)
}

// RunAsync API in Async mode after lookup
func RunAsync(yml load.Config, samplesToMerge *load.SamplesToMerge, originalAPINo int) {
	load.Logrus.WithFields(logrus.Fields{
		"name":     yml.Name,
		"apis":     len(yml.APIs),
	}).Debug("config: processing apis: ASYNC mode. Will skip StoreLookups VariableLookups for: ")

	// load secrets
	loadSecrets(&yml)
	var wgapi sync.WaitGroup
	wgapi.Add(len(yml.APIs))

	for i := range yml.APIs {
		go func(originalAPINo int, i int) {
			defer wgapi.Done()
			dataSets := FetchData(i, &yml, samplesToMerge)
			processor.RunDataHandler(dataSets, samplesToMerge, i, &yml, originalAPINo)
		}(originalAPINo, i)
	}
	wgapi.Wait()

	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: finished processing apis: ASYNC Mode")

	// processor.ProcessSamplesToMerge(&samplesToMerge, &yml)
	// hren joinAndMerge processing - replacing processor.ProcessSamplesToMerge
	// ProcessSamplesMergeJoin will be processed in the run() function for the whole config
	// processor.ProcessSamplesMergeJoin(&samplesToMerge, &yml)
}

// RunSync API in Sync mode after lookup
func RunSync(yml load.Config, samplesToMerge *load.SamplesToMerge, originalAPINo int) {
	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: processing apis: Sync Mode")

	// load secrets
	loadSecrets(&yml)

	for i := range yml.APIs {
		dataSets := FetchData(i, &yml, samplesToMerge)
		processor.RunDataHandler(dataSets, samplesToMerge, i, &yml, originalAPINo)
	}

	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: finished processing apis: Sync Mode")

	// processor.ProcessSamplesToMerge(&samplesToMerge, &yml)
	// hren joinAndMerge processing - replacing processor.ProcessSamplesToMerge
	// ProcessSamplesMergeJoin will be processed in the run() function for the whole config
	// processor.ProcessSamplesMergeJoin(&samplesToMerge, &yml)
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

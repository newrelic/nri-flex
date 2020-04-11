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
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/processor"
	"github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

// LoadFiles Loads Flex config files
func LoadFiles(configs *[]load.Config, files []os.FileInfo, path string) []error {
	var errors []error
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

		// since we're reading many files, continue if one of more fails
		err := LoadFile(configs, f, path)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// LoadV4IntegrationConfig Loads Agent/Flex config files
func LoadV4IntegrationConfig(v4Str string, configs *[]load.Config, fileName string, filePath string) error {
	c := load.AgentConfig{}
	err := yaml.Unmarshal([]byte(v4Str), &c)

	if err != nil {
		load.Logrus.WithError(err).Error("config: failed to unmarshal v4 config file")
		return err
	}

	load.Logrus.Warn("config: testing agent config, agent features will not be available")

	for _, integration := range c.Integrations {
		// ensure it is a flex based integration
		if integration.Name == "nri-flex" {
			newConfig := integration.Config
			newConfig.FileName = fileName
			newConfig.FilePath = filePath

			if newConfig.Name == "" {
				load.Logrus.WithFields(logrus.Fields{
					"file": filePath,
				}).WithError(err).Error("config: flexConfig requires name")
				return err
			}

			// if lookup files exist we need to potentially create multiple config files
			if newConfig.LookupFile != "" {
				err = SubLookupFileData(configs, newConfig)
				if err != nil {
					load.Logrus.WithFields(logrus.Fields{
						"file": filePath,
					}).WithError(err).Error("config: failed to sub lookup file data")
				}
			} else {
				*configs = append(*configs, newConfig)
			}
		}
	}

	return nil
}

// LoadFile loads a single Flex config file
func LoadFile(configs *[]load.Config, f os.FileInfo, path string) error {
	filePath := path + f.Name()

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"file": filePath,
			"err":  err,
		}).Error("config: failed to read config file")
		return err
	}

	ymlStr := string(b)
	SubEnvVariables(&ymlStr)
	SubTimestamps(&ymlStr, time.Now())

	// Check if V4 Agent configuration
	if strings.HasPrefix(ymlStr, "integrations:") {
		err := LoadV4IntegrationConfig(ymlStr, configs, f.Name(), path)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"file": filePath,
			}).WithError(err).Error("config: failed to load v4 config file")
			return err
		}
	} else {
		config, err := ReadYML(ymlStr)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"file": filePath,
			}).WithError(err).Error("config: failed to load config file")
			return err
		}

		config.FileName = f.Name()
		config.FilePath = path
		if config.Name == "" {
			load.Logrus.WithFields(logrus.Fields{
				"file": filePath,
			}).WithError(err).Error("config: flexConfig requires name")
			return err
		}

		checkIngestConfigs(&config)

		// if lookup files exist we need to potentially create multiple config files
		if config.LookupFile != "" {
			err = SubLookupFileData(configs, config)
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"file": filePath,
				}).WithError(err).Error("config: failed to sub lookup file data")
			}
		} else {
			*configs = append(*configs, config)
		}
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
		}).WithError(err).Debug("config: failed to read")
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
	// samplesToMerge := map[string][]interface{}{}
	var samplesToMerge load.SamplesToMerge
	samplesToMerge.Data = map[string][]interface{}{}
	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: processing apis")

	// load secrets
	_ = loadSecrets(&yml)

	// intentionally handled synchronously
	for i := range yml.APIs {
		if err := runVariableProcessor(&yml); err != nil {
			load.Logrus.WithError(err).Error("config: variable processor error")
		}
		dataSets := FetchData(i, &yml, &samplesToMerge)
		processor.RunDataHandler(dataSets, &samplesToMerge, i, &yml, i)
	}

	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: finished variable processing apis")

	// processor.ProcessSamplesToMerge(&samplesToMerge, &yml)
	// hren MergeAndJoin processing - replacing processor.ProcessSamplesToMerge
	processor.ProcessSamplesMergeJoin(&samplesToMerge, &yml)
}

// RunAsync API in Async mode after lookup
func RunAsync(yml load.Config, samplesToMerge *load.SamplesToMerge, originalAPINo int) {
	load.Logrus.WithFields(logrus.Fields{
		"name": yml.Name,
		"apis": len(yml.APIs),
	}).Debug("config: processing apis: ASYNC mode. Will skip StoreLookups VariableLookups for: ")

	// load secrets
	_ = loadSecrets(&yml)
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
	_ = loadSecrets(&yml)

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
func RunFiles(configs *[]load.Config) []error {
	var errors []error
	if load.Args.ProcessConfigsSync {
		for _, cfg := range *configs {
			err := verifyConfig(cfg)
			if err != nil {
				errors = append(errors, err)
			} else {
				load.Logrus.WithFields(logrus.Fields{"name": cfg.Name}).Debug("config: running sync")
				Run(cfg)
				load.StatusCounterIncrement("ConfigsProcessed")
			}
		}
	} else {
		errorChannel := make(chan error)
		// listen for errors coming from the verification of the configs and store them for later
		go func() {
			for err := range errorChannel {
				if err != nil {
					errors = append(errors, err)
				}
			}
		}()

		var wg sync.WaitGroup
		wg.Add(len(*configs))
		for _, cfg := range *configs {
			go func(cfg load.Config) {
				defer wg.Done()
				err := verifyConfig(cfg)
				if err != nil {
					errorChannel <- err
				} else {
					load.Logrus.WithFields(logrus.Fields{"name": cfg.Name}).Debug("config: running async")
					Run(cfg)
					load.StatusCounterIncrement("ConfigsProcessed")
				}
			}(cfg)
		}

		wg.Wait()
		close(errorChannel)
	}

	load.Logrus.WithFields(logrus.Fields{
		"configs": load.StatusCounterRead("ConfigsProcessed"),
	}).Info("flex: completed processing configs")

	return errors
}

// verifyConfig ensure the config file doesn't have anything it should not run
func verifyConfig(cfg load.Config) error {
	if strings.HasPrefix(cfg.FileName, "cd-") && !cfg.ContainerDiscovery.ReplaceComplete {
		return fmt.Errorf("config: failed to apply discovery to config: '%s'", cfg.Name)
	}
	ymlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	ymlStr := string(ymlBytes)
	if strings.Contains(ymlStr, "${auto:host}") || strings.Contains(ymlStr, "${auto:port}") {
		return fmt.Errorf("config: cannot have 'auto' token replacements: '%s'", cfg.Name)
	}
	return nil
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

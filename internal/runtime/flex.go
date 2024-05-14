/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/outputs"
	"github.com/newrelic/nri-flex/internal/utils"
	"github.com/sirupsen/logrus"
)

// Serverless runtimes must implement this
type Instance interface {
	isAvailable() bool
	loadConfigs(*[]load.Config) error
	SetConfigDir(string)
}

// Add new  runtime types to this  list. Test & Default don't go here
var runtimeTypes = [2]Instance{new(Lambda), new(Function)}
var log = load.Logrus

// Get the first available runtime type, defaults to the server-based (Linux | Windows) Default type
func GetFlexRuntime() Instance {
	for _, r := range runtimeTypes {
		if r.isAvailable() {
			return r
		}
	}
	return GetDefaultRuntime()
}

// Get a server-based, default, runtime
func GetDefaultRuntime() Instance {
	return new(Default)
}

// Get the test runtime
func GetTestRuntime() Instance {
	return new(Test)
}

// Post-initialization common to all runtime types here
func CommonPostInit() {
	err := load.Integration.Publish()
	if err != nil {
		load.Logrus.WithError(err).Fatal("runtime.CommonPostInit: failed to publish")
	}
}

// CommonPreInit Pre-initialization common to all runtime types here
func CommonPreInit() {
	load.SetupLogger()
	load.StartTime = load.MakeTimestamp()
	setEnvs()

	err := outputs.InfraIntegration()
	if err != nil {
		load.Logrus.WithError(err).Fatal("runtime.CommonPreInit: failed to initialize runtime")
	}
}

// RunFlex Common run (once) function
func RunFlex(instance Instance) error {
	setStatusCounters()

	log.WithFields(logrus.Fields{
		"version": load.IntegrationVersion,
		"GOOS":    runtime.GOOS,
		"GOARCH":  runtime.GOARCH,
	}).Info(load.IntegrationName)

	var configs []load.Config

	// runtime instance specific run
	err := instance.loadConfigs(&configs)
	if err != nil {
		return err
	}

	errors := config.RunFiles(&configs)
	if len(errors) > 0 {
		return fmt.Errorf("runtime.RunFlex: failed to run configuration files")
	}

	outputs.StatusSample()

	if load.Args.InsightsURL != "" && load.Args.InsightsAPIKey != "" {
		for _, batch := range outputs.GetMetricBatches() {
			if err := outputs.SendBatchToInsights(batch); err != nil {
				log.WithError(err).Error("runtime.RunFlex: failed to send batch to insights")
			}
		}
	} else if load.Args.LogApiURL != "" && load.Args.LogApiKey != "" {
		for _, batch := range outputs.GetLogMetricBatches() {
			if err := outputs.SendBatchToLogApi(batch); err != nil {
				log.WithError(err).Error("runtime.RunFlex: failed to send batch to log api")
			}
		}
	} else if load.Args.MetricAPIUrl != "" && (load.Args.InsightsAPIKey != "" || load.Args.MetricAPIKey != "") && len(load.MetricsStore.Data) > 0 {
		if err := outputs.SendToMetricAPI(); err != nil {
			log.WithError(err).Error("runtime.RunFlex: failed to send metrics")
		}
	} else if len(load.MetricsStore.Data) > 0 && (load.Args.MetricAPIUrl == "" || (load.Args.InsightsAPIKey == "" || load.Args.MetricAPIKey == "")) {
		log.Debug("runtime.RunFlex: metric_api is being used, but metric url and/or key has not been set")
	}
	return nil
}

func addSingleConfigFile(configFile string, configs *[]load.Config) error {
	file, err := os.Stat(configFile)
	if err != nil {
		log.WithFields(logrus.Fields{
			"err":  err,
			"file": configFile,
		}).Fatal("config: failed to read")
		return err
	}
	path := strings.Replace(filepath.FromSlash(configFile), file.Name(), "", -1)
	err = config.LoadFile(configs, file, path)
	return err
}

func addConfigsFromPath(path string, configs *[]load.Config) []error {
	configPath := filepath.FromSlash(path)
	files, err := ioutil.ReadDir(configPath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"dir": path,
		}).WithError(err).Fatal("config: failed to read configuration folder")
		return []error{err}
	}

	return config.LoadFiles(configs, files, configPath)
}

func logEncryptPass() error {
	log.Info("*****Encryption Result*****")
	cipherText, err := utils.Encrypt([]byte(load.Args.EncryptPass), load.Args.PassPhrase)
	if err != nil {
		log.WithError(err).Error("EncryptPass: Failed to encrypt")
		return err
	}
	cleartext, err := utils.Decrypt(cipherText, load.Args.PassPhrase)
	if err != nil {
		log.WithError(err).Error("EncryptPass: Failed to Decrypt")
		return err
	}
	log.Infof("   encrypt_pass: %s", cleartext)
	log.Infof("    pass_phrase: %s", load.Args.PassPhrase)
	log.Infof(" encrypted pass: %x", cipherText)

	return nil
}

func setStatusCounters() {
	log.Out = os.Stderr
	load.FlexStatusCounter.M = make(map[string]int)
	load.FlexStatusCounter.M["EventCount"] = 0
	load.FlexStatusCounter.M["EventDropCount"] = 0
	load.FlexStatusCounter.M["ConfigsProcessed"] = 0
}

// setEnvs set environment variable argument overrides
func setEnvs() {
	gitService := os.Getenv("GIT_SERVICE")
	if gitService != "" {
		load.Args.GitService = gitService
	}
	gitRepo := os.Getenv("GIT_REPO")
	if gitRepo != "" {
		load.Args.GitRepo = gitRepo
		load.Args.GitToken = os.Getenv("GIT_TOKEN")
		load.Args.GitUser = os.Getenv("GIT_USER")
	}
	insightsAPIKey := os.Getenv("INSIGHTS_API_KEY")
	if insightsAPIKey != "" {
		load.Args.InsightsAPIKey = insightsAPIKey
		load.Args.InsightsURL = os.Getenv("INSIGHTS_URL")
	}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		load.IsKubernetes = true
	}
	eventLimit, err := strconv.Atoi(os.Getenv("EVENT_LIMIT"))
	if err == nil && eventLimit > 0 {
		load.Args.EventLimit = eventLimit
	}
	configSync, err := strconv.ParseBool(os.Getenv("PROCESS_CONFIGS_SYNC"))
	if err == nil && configSync {
		load.Args.ProcessConfigsSync = configSync
	}
	asyncRate, err := strconv.Atoi(os.Getenv("ASYNC_RATE"))
	if err == nil && asyncRate > 0 {
		load.Args.AsyncRate = asyncRate
	}

	fargate, err := strconv.ParseBool(os.Getenv("FARGATE"))
	if err == nil && fargate {
		load.Args.Fargate = fargate
	}
	cd, err := strconv.ParseBool(os.Getenv("CONTAINER_DISCOVERY"))
	if err == nil && cd {
		load.Args.ContainerDiscovery = cd
	}
	load.Args.MetricAPIUrl = os.Getenv("METRIC_API_URL")
	load.Args.MetricAPIKey = os.Getenv("METRIC_API_KEY")
}

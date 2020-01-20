/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/discovery"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/outputs"
	"github.com/newrelic/nri-flex/internal/utils"
	"github.com/sirupsen/logrus"
)

// RunFlex runs flex
// if mode is "" run in default mode
func RunFlex(mode string) {
	verboseLogging := os.Getenv("VERBOSE")
	if load.Args.Verbose || verboseLogging == "true" || verboseLogging == "1" {
		load.Logrus.SetLevel(logrus.TraceLevel)
	}

	load.Logrus.WithFields(logrus.Fields{
		"version": load.IntegrationVersion,
		"GOOS":    runtime.GOOS,
		"GOARCH":  runtime.GOARCH,
	}).Info(load.IntegrationName)

	// store config ymls
	configs := []load.Config{}

	switch mode {
	case "lambda":
		addConfigsFromPath("/var/task/pkg/flexConfigs/", &configs)
		if config.SyncGitConfigs("/tmp/") {
			addConfigsFromPath("/tmp/", &configs)
		}
	default:
		// running as default
		if load.Args.EncryptPass != "" {
			load.Logrus.Info(fmt.Sprintf("*****Encryption Result*****"))
			ciphertext, err := utils.Encrypt([]byte(load.Args.EncryptPass), load.Args.PassPhrase)
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{"err": err}).Error("EncryptPass: Failed to encrypt")
			}
			cleartext, err := utils.Decrypt(ciphertext, load.Args.PassPhrase)
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{"err": err}).Error("EncryptPass: Failed to Decrypt")
			}
			load.Logrus.Info(fmt.Sprintf("   encrypt_pass: %s", cleartext))
			load.Logrus.Info(fmt.Sprintf("    pass_phrase: %s", load.Args.PassPhrase))
			load.Logrus.Info(fmt.Sprintf(" encrypted pass: %x", ciphertext))
			os.Exit(0)

		} else {
			config.SyncGitConfigs("")
			if load.Args.ConfigFile != "" {
				addSingleConfigFile(load.Args.ConfigFile, &configs)
			} else {
				addConfigsFromPath(load.Args.ConfigDir, &configs)
			}
			if load.Args.ContainerDiscovery || load.Args.Fargate {
				discovery.Run(&configs)
			}
		}

	}

	if load.ContainerID == "" && mode != "test" && mode != "lambda" && runtime.GOOS != "darwin" {
		discovery.Processes()
	}

	load.Logrus.Info(fmt.Sprintf("flex: config files loaded %d", len(configs)))

	config.RunFiles(&configs)
	outputs.StatusSample()

	if load.Args.InsightsURL != "" && load.Args.InsightsAPIKey != "" {
		outputs.SendToInsights()
	} else if load.Args.MetricAPIUrl != "" && (load.Args.InsightsAPIKey != "" || load.Args.MetricAPIKey != "") && len(load.MetricsStore.Data) > 0 {
		outputs.SendToMetricAPI()
	} else if len(load.MetricsStore.Data) > 0 && (load.Args.MetricAPIUrl == "" || (load.Args.InsightsAPIKey == "" || load.Args.MetricAPIKey == "")) {
		load.Logrus.Debug("flex: metric_api is being used, but metric url and/or key has not been set", len(configs))
	}
}

func addSingleConfigFile(configFile string, configs *[]load.Config) {
	file, err := os.Stat(configFile)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err":  err,
			"file": configFile,
		}).Fatal("config: failed to read")
	}
	path := strings.Replace(filepath.FromSlash(configFile), file.Name(), "", -1)
	if err := config.LoadFile(configs, file, path); err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err":  err,
			"file": configFile,
		}).Error("config: failed to load file")
	}
}

func addConfigsFromPath(path string, configs *[]load.Config) {
	configPath := filepath.FromSlash(path)
	files, err := ioutil.ReadDir(configPath)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
			"dir": path,
		}).Fatal("config: failed to read")
	}
	config.LoadFiles(configs, files, configPath)
}

// SetDefaults set flex defaults
func SetDefaults() {
	load.Logrus.Out = os.Stderr
	load.FlexStatusCounter.M = make(map[string]int)
	load.FlexStatusCounter.M["EventCount"] = 0
	load.FlexStatusCounter.M["EventDropCount"] = 0
	load.FlexStatusCounter.M["ConfigsProcessed"] = 0
}

// SetEnvs set environment variable argument overrides
func SetEnvs() {
	load.AWSExecutionEnv = os.Getenv("AWS_EXECUTION_ENV")
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

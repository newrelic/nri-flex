/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package integration

import (
	"fmt"
	"github.com/newrelic/nri-flex/internal/utils"
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
	"github.com/sirupsen/logrus"
)

// FlexRunMode is used to switch the mode of running flex.
type FlexRunMode int

const (
	// FlexModeDefault is the usual way of running flex.
	FlexModeDefault FlexRunMode = iota
	// FlexModeLambda is used when flex is running within a lambda.
	FlexModeLambda
	// FlexModeTest is used when running tests.
	FlexModeTest
)

// RunFlex runs flex.
func RunFlex(runMode FlexRunMode) error {
	setupLogger()

	load.Logrus.WithFields(logrus.Fields{
		"version": load.IntegrationVersion,
		"GOOS":    runtime.GOOS,
		"GOARCH":  runtime.GOARCH,
	}).Info(load.IntegrationName)

	// store config yml
	var configs []load.Config

	switch runMode {
	case FlexModeLambda:
		errors := addConfigsFromPath("/var/task/pkg/flexConfigs/", &configs)
		if len(errors) > 0 {
			return fmt.Errorf("flex: failed to read some configuration files. please review them before continuing")
		}

		isSyncGitConfigured, err := config.SyncGitConfigs("/tmp/")
		if err != nil {
			logrus.WithError(err).Warn("flex: failed to sync git configs")
		} else if isSyncGitConfigured {
			errors = addConfigsFromPath("/tmp/", &configs)
			if len(errors) > 0 {
				logrus.WithError(err).Warn("flex: failed to load git sync configuration files. ignoring and continuing")
			}
		}

	default:
		if load.Args.EncryptPass != "" {
			err := logEncryptPass()
			if err != nil {
				return err
			}
		}

		_, err := config.SyncGitConfigs("")
		if err != nil {
			logrus.WithError(err).Warn("flex: failed to sync git configs")
		}

		var errors []error
		if load.Args.ConfigFile != "" {
			err = addSingleConfigFile(load.Args.ConfigFile, &configs)
			if err != nil {
				return err
			}
		} else {
			errors = addConfigsFromPath(load.Args.ConfigDir, &configs)
			if len(errors) > 0 {
				return fmt.Errorf("flex: failed to load configurations files")
			}
		}

		// should we stop if this fails??
		if load.Args.ContainerDiscovery || load.Args.Fargate {
			discovery.Run(&configs)
		}
	}

	if load.ContainerID == "" && runMode == FlexModeDefault {
		switch runtime.GOOS {
		case "windows":
			if load.Args.DiscoverProcessWin {
				discovery.Processes()
			}
		case "linux":
			if load.Args.DiscoverProcessLinux {
				discovery.Processes()
			}
		}
	}
	errors := config.RunFiles(&configs)
	if len(errors) > 0 {
		return fmt.Errorf("flex: failed to run configuration files")
	}

	outputs.StatusSample()

	if load.Args.InsightsURL != "" && load.Args.InsightsAPIKey != "" {
		for _, batch := range outputs.GetMetricBatches() {
			if err := outputs.SendBatchToInsights(batch); err != nil {
				load.Logrus.WithError(err).Error("flex: failed to send batch to insights")
			}
		}
	} else if load.Args.MetricAPIUrl != "" && (load.Args.InsightsAPIKey != "" || load.Args.MetricAPIKey != "") && len(load.MetricsStore.Data) > 0 {
		if err := outputs.SendToMetricAPI(); err != nil {
			load.Logrus.WithError(err).Error("flex: failed to send metrics")
		}
	} else if len(load.MetricsStore.Data) > 0 && (load.Args.MetricAPIUrl == "" || (load.Args.InsightsAPIKey == "" || load.Args.MetricAPIKey == "")) {
		load.Logrus.Debug("flex: metric_api is being used, but metric url and/or key has not been set")
	}
	return nil
}

func setupLogger() {
	verboseLogging := os.Getenv("VERBOSE")
	if load.Args.Verbose || verboseLogging == "true" || verboseLogging == "1" {
		load.Logrus.SetLevel(logrus.TraceLevel)
	}

	if load.Args.StructuredLogs {
		load.Logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func addSingleConfigFile(configFile string, configs *[]load.Config) error {
	file, err := os.Stat(configFile)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
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
		load.Logrus.WithFields(logrus.Fields{
			"dir": path,
		}).WithError(err).Fatal("config: failed to read confguration folder")
		return []error{err}
	}

	return config.LoadFiles(configs, files, configPath)
}

func logEncryptPass() error {
	load.Logrus.Info("*****Encryption Result*****")
	cipherText, err := utils.Encrypt([]byte(load.Args.EncryptPass), load.Args.PassPhrase)
	if err != nil {
		load.Logrus.WithError(err).Error("EncryptPass: Failed to encrypt")
		return err
	}
	cleartext, err := utils.Decrypt(cipherText, load.Args.PassPhrase)
	if err != nil {
		load.Logrus.WithError(err).Error("EncryptPass: Failed to Decrypt")
		return err
	}
	load.Logrus.Infof("   encrypt_pass: %s", cleartext)
	load.Logrus.Infof("    pass_phrase: %s", load.Args.PassPhrase)
	load.Logrus.Infof(" encrypted pass: %x", cipherText)

	return nil
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

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package runtime

import (
	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
	"os"
)

// GCP Function runtime
type Function struct {
	configDir string
}

// Test to see if we're running as a GCP Function
func (i *Function) isAvailable() bool {
	log.Debugf("Function.isAvailable: enter")
	status := false
	if os.Getenv("FUNCTION_TARGET") != "" {
		err := i.init()
		if err != nil {
			load.Logrus.WithError(err).Fatal("Function.isAvailable: failed to validate Function required config")
		}
		status = true
	}
	log.Debugf("Function.isAvailable: exit status: %t", status)
	return status
}

func (i *Function) loadConfigs(configs *[]load.Config) error {
	log.Debugf("Function.loadConfigs: enter")

	// Get the configs
	errors := addConfigsFromPath(i.configDir, configs)
	if len(errors) > 0 {
		log.Error("Function.loadConfigs: failed to read some configuration files, please review them")
	}

	// Sync to Git if required. Does Function need this?
	isSyncGitConfigured, err := config.SyncGitConfigs("/tmp/")
	if err != nil {
		log.WithError(err).Warn("Function.loadConfigs: failed to sync git configs")
	} else if isSyncGitConfigured {
		errors = addConfigsFromPath("/tmp/", configs)
		if len(errors) > 0 {
			log.Error("Function.loadConfigs: failed to load git sync configuration files, ignoring and continuing")
		}
	}
	log.Debugf("Function.loadConfigs: exit")
	return nil
}

func (i *Function) init() error {
	log.Debugf("Function.init: enter")
	i.SetConfigDir("./serverless_function_source_code/flexConfigs/")

	load.ServerlessName = os.Getenv("X_GOOGLE_FUNCTION_NAME")
	load.ServerlessExecutionEnv = os.Getenv("GOOGLE_RUNTIME") + " " + os.Getenv("GOOGLE_RUNTIME_VERSION")
	load.Logrus.SetFormatter(&logrus.JSONFormatter{})
	log.Debugf("Function.init: exit")

	return nil
}

func (i *Function) SetConfigDir(s string) {
	log.Debugf("SetConfigDir: enter: %s", i.configDir)
	i.configDir = s
	log.Debugf("SetConfigDir: exit: %s", i.configDir)
}

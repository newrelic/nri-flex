/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package runtime

import (
	"context"
	"fmt"
	"github.com/newrelic/nri-flex/internal/config"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

// AWS Lambda runtime
type Lambda struct {
	configDir string
}

// Test to see if we're running as an AWS Lambda
func (i *Lambda) isAvailable() bool {
	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		err := i.init()
		if err != nil {
			load.Logrus.WithError(err).Fatal("Lambda.isAvailable: failed to validate lambda required config")
		}
		return true
	}
	return false
}

// Run Flex in the Lambda
func (i *Lambda) loadConfigs(configs *[]load.Config) error {
	load.Logrus.Info("Lambda.loadConfigs: running as Lambda")
	errors := addConfigsFromPath(i.configDir, configs)
	if len(errors) > 0 {
		log.Error("Lambda.loadConfigs: failed to read some configuration files, please review them")
	}

	isSyncGitConfigured, err := config.SyncGitConfigs("/tmp/")
	if err != nil {
		log.WithError(err).Warn("Lambda.loadConfigs: failed to sync git configs")
	} else if isSyncGitConfigured {
		errors = addConfigsFromPath("/tmp/", configs)
		if len(errors) > 0 {
			log.Error("Lambda.loadConfigs: failed to load git sync configuration files, ignoring and continuing")
		}
	}
	return nil
}

// init: while running within a Lambda insights url and api key are required.
func (i *Lambda) init() error {
	if i.configDir == "" {
		i.configDir = "/var/task/pkg/flexConfigs/"
	}
	if load.Args.InsightsURL == "" || load.Args.InsightsAPIKey == "" {
		return fmt.Errorf(" Lambda.init: missing insights URL and/or API key")
	}
	load.ServerlessName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	load.ServerlessExecutionEnv = os.Getenv("AWS_EXECUTION_ENV")
	load.Logrus.SetFormatter(&logrus.JSONFormatter{})

	// Register the Lambda entrypoint with  AWS
	lambda.Start(i.FlexAsALambdaHandler)
	return nil
}

// FlexAsALambdaHandler receives the incoming lambda request, from the AWS perspective this is the entry point
func (i *Lambda) FlexAsALambdaHandler(ctx context.Context, event interface{}) (string, error) {
	if event != nil {
		load.IngestData = event
	}

	_ = RunFlex(i)
	return "Flex Lambda Complete", nil
}

func (i *Lambda) SetConfigDir(s string) {
	i.configDir = s
}

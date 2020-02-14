/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package integration

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

// IsLambda check if Flex is running within a Lambda.
func IsLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// ValidateLambdaConfig: while running within a Lambda insights url and api key are required.
func ValidateLambdaConfig() error {
	if load.Args.InsightsURL == "" || load.Args.InsightsAPIKey == "" {
		return fmt.Errorf("lambda: missing insights URL and/or API key")
	}
	return nil
}

// HandleLambda handles lambda invocation
func HandleLambda() {
	load.LambdaName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	load.Logrus.SetFormatter(&logrus.JSONFormatter{})

	lambda.Start(HandleRequest)
}

// HandleRequest Handles incoming lambda request
func HandleRequest(ctx context.Context, event interface{}) (string, error) {
	load.Logrus.Info("flex: running as lambda")
	SetDefaults()

	if event != nil {
		load.IngestData = event
	}

	_ = RunFlex(FlexModeLambda)

	return "Flex Lambda Complete", nil
}

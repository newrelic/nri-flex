/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
)

// Test the AWS Lambda path without having to run on AWS
// Use: go run test/serverless/aws/lambda.go
//      use -help to see the command line params
func main() {
	runtime.CommonPreInit()

	i := new(runtime.Lambda)

	i.SetConfigDir("./flexConfigs")
	err := runtime.RunFlex(i)
	if err != nil {
		load.Logrus.WithError(err).Fatal("Lambda: failed to run runtime")
	}

	runtime.CommonPostInit()
	_, err = i.FlexAsALambdaHandler(context.TODO(), nil)
	if err != nil {
		load.Logrus.WithError(err).Fatal("Lambda: failed to run handler")
	}

}

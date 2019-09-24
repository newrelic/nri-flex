/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/newrelic/nri-flex/internal/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/outputs"
	"github.com/sirupsen/logrus"
)

func main() {
	load.StartTime = load.MakeTimestamp()
	integration.SetEnvs()
	outputs.InfraIntegration()

	if integration.LambdaCheck() {
		integration.Lambda()
	} else {
		// default process
		integration.SetDefaults()
		integration.RunFlex("")
	}

	err := load.Integration.Publish()
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{"err": err}).Fatal("flex: unable to publish")
	}
}

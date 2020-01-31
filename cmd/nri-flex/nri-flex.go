/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/newrelic/nri-flex/internal/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/outputs"
)

func main() {
	load.StartTime = load.MakeTimestamp()
	integration.SetEnvs()
	outputs.InfraIntegration()

	if integration.IsLambda() {
		if err := integration.ValidateLambdaConfig(); err != nil {
			load.Logrus.WithError(err).Fatal("flex: failed to validate lambda required config")
		}
		integration.HandleLambda()
	} else {
		// default process
		integration.SetDefaults()
		integration.RunFlex(integration.FlexModeDefault)
	}

	if err := load.Integration.Publish(); err != nil {
		load.Logrus.WithError(err).Fatal("flex: unable to publish")
	}
}

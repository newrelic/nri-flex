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

	err := outputs.InfraIntegration()
	if err != nil {
		load.Logrus.WithError(err).Fatal("flex: failed to initialize integration")
	}

	if integration.IsLambda() {
		err = integration.ValidateLambdaConfig()
		if err != nil {
			load.Logrus.WithError(err).Fatal("flex: failed to validate lambda required config")
		}
		integration.HandleLambda()
	} else {
		// default process
		integration.SetDefaults()
		err = integration.RunFlex(integration.FlexModeDefault)
		if err != nil {
			load.Logrus.WithError(err).Fatal("flex: failed to run integration")
		}
	}

	err = load.Integration.Publish()
	if err != nil {
		load.Logrus.WithError(err).Fatal("flex: failed to publish")
	}
}

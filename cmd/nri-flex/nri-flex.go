/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"time"

	"github.com/newrelic/nri-flex/internal/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	"github.com/newrelic/nri-flex/internal/outputs"
)

func main() {
	load.StartTime = time.Now().UnixNano()
	integration.SetEnvs()
	outputs.InfraIntegration()

	if integration.LambdaCheck() {
		integration.Lambda()
	} else {
		// default process
		integration.SetDefaults()
		integration.RunFlex("")
	}

	logger.Flex("fatal", load.Integration.Publish(), "unable to publish", false)
}

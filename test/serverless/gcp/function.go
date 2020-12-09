/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/newrelic/nri-flex"
	"github.com/newrelic/nri-flex/internal/load"
	"time"
)

// Test the GCP Function path without having to run on GCP
// Use: go run test/serverless/gcp/function.go
//      use -help to see the command line params
func main() {
	load.Logrus.Infof("test.function.main: enter")
	for {
		nriflex.FlexPubSub(nil, nriflex.PubSubMessage{})
		time.Sleep(60 * time.Second)
		load.Logrus.Info("")
		load.Logrus.Info("")
		load.Logrus.Info("")
	}
	load.Logrus.Infof("test.function.main: exit")
}

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
	"sync"
)

var log = load.Logrus

// Test the GCP Function path without having to run on GCP
// Use: go run test/serverless/gcp/function.go
//      use -help to see the command line params
func main() {
	log.Infof("main: enter")

	var r runtime.Instance

	var once sync.Once
	once.Do(func() {
		log.Infof("main: once.Do: enter")
		// Generate the Function runtime singleton
		r = runtime.GetFlexRuntime()
		log.Infof("main: once.Do: exit")
	})

	r.SetConfigDir("./flexConfigs/")
	runtime.CommonPreInit()
	err := runtime.RunFlex(r)
	if err != nil {
		load.Logrus.WithError(err).Fatal("main: failed to run runtime")
	}
	runtime.CommonPostInit()
	log.Infof("main: exit")
}

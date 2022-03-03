/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
)

func main() {
	runtime.CommonPreInit()

	i := runtime.GetFlexRuntime()
	err := runtime.RunFlex(i)
	if err != nil {
		load.Logrus.WithError(err).Fatal("flex: failed to run runtime")
	}
	load.Logrus.Info("test")
	runtime.CommonPostInit()
}

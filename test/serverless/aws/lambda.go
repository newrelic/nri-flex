/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"
	"fmt"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
	rt "runtime"
	"time"
)

// Test the AWS Lambda path without having to run on AWS
// Use: go run test/serverless/aws/lambda.go
//      use -help to see the command line params
func main() {
	for {
		runtime.CommonPreInit()

		i := new(runtime.Lambda)

		i.SetConfigDir("./flexConfigs")

		start := time.Now()
		err := runtime.RunFlex(i)
		load.Logrus.Infof("Test: elapsed time: %s", time.Since(start))
		printMemUsage()
		time.Sleep(60 * time.Second)
		load.Logrus.Info("")
		load.Logrus.Info("")
		load.Logrus.Info("")

		if err != nil {
			load.Logrus.WithError(err).Fatal("Lambda: failed to run runtime")
		}

		runtime.CommonPostInit()
		_, err = i.FlexAsALambdaHandler(context.TODO(), nil)
		if err != nil {
			load.Logrus.WithError(err).Fatal("Lambda: failed to run handler")
		}
	}

}
func printMemUsage() {
	var m rt.MemStats
	rt.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

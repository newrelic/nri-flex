/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"
	"fmt"
	"github.com/newrelic/nri-flex"
	"github.com/newrelic/nri-flex/internal/load"
	"runtime"
	"time"
)

// Test the GCP Function path without having to run on GCP
// Use: go run test/serverless/gcp/function.go
//      use -help to see the command line params
func main() {
	load.Logrus.Infof("test.function.main: enter")
	for {
		start := time.Now()
		_ = nriflex.FlexPubSub(context.TODO(), nriflex.PubSubMessage{})
		load.Logrus.Infof("Test: elapsed time: %s", time.Since(start))
		printMemUsage()
		time.Sleep(60 * time.Second)
		load.Logrus.Info("")
		load.Logrus.Info("")
		load.Logrus.Info("")
	}
}
func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

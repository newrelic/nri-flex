/*
* Copyright 2020 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package nriflex

import (
	"context"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
	"net/http"
	"sync"
)

// nri-flex.main is not called by GCP, so everything has to happen here

// Singleton, we only want to set this once
var r runtime.Instance
var log = load.Logrus

// Target Function endpoint if you want to use HTTP  with the Cloud Scheduler
func FlexHTTP(w http.ResponseWriter, req *http.Request) {
	log.Debugf("FlexHTTP: enter")
	run()
	log.Debugf("FlexHTTP: exit")
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Target Function endpoint if you want to use Pub/Sub  with the Cloud Scheduler
func FlexPubSub(ctx context.Context, m PubSubMessage) error {
	log.Debugf("FlexPubSub: enter")
	run()
	log.Debugf("FlexPubSub: exit")
	return nil
}

// One time only init
var once sync.Once

// Do the actual, common, work here
func run() {
	log.Debugf("nriflex.run: enter")
	once.Do(func() {
		log.Debugf("nriflex.run: once.Do: enter")
		// Generate the Function runtime singleton
		r = runtime.GetFlexRuntime()
		log.Debugf("nriflex.run: once.Do: exit")
	})
	runtime.CommonPreInit()
	err := runtime.RunFlex(r)
	if err != nil {
		load.Logrus.WithError(err).Fatal("flex: failed to run runtime")
	}
	runtime.CommonPostInit()
	log.Debugf("nriflex.run: exit")
}

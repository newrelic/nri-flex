/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	fintegration "github.com/newrelic/nri-flex/internal/integration"
	"github.com/sirupsen/logrus"

	"github.com/newrelic/nri-flex/internal/load"
)

// testSamples as samples could be generated in different orders, so we test per sample
func testSamples(expectedSamples []string, entityMetrics []*metric.Set, t *testing.T) {
	if len(entityMetrics) != len(expectedSamples) {
		t.Errorf("Missing samples, got: %v, want: %v.", (entityMetrics), (expectedSamples))
	}

	for _, expectedSample := range expectedSamples {
		matchedSample := false
		for _, sample := range entityMetrics {
			delete(sample.Metrics, "flex.time.startNs")
			delete(sample.Metrics, "flex.time.endNs")
			delete(sample.Metrics, "flex.time.elaspedNs")
			out, err := sample.MarshalJSON()
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"err": err,
				}).Debug("failed to marshal")
			} else if expectedSample == string(out) {
				matchedSample = true
				break
			}
		}
		if !matchedSample {
			completeMetrics, _ := json.Marshal(entityMetrics)
			t.Errorf("Unable to find expected payload, received: %v, want: %v.", string(completeMetrics), expectedSample)
		}
	}
}

func TestConfigDir(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestReadJsonCmdDir", "nri-flex")
	load.Args.ConfigDir = "../../test/configs/"
	fintegration.RunFlex("test")
	expectedSamples := []string{
		`{"event_type":"flexStatusSample","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.commandJsonOutSample":1}`,
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.newrelic.nri-flex",` +
			`"integration_version":"Unknown-SNAPSHOT","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestConfigFile(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestReadJsonCmd", "nri-flex")
	load.Args.ConfigFile = "../../test/configs/json-read-cmd-example.yml"
	fintegration.RunFlex("test")
	expectedSamples := []string{
		`{"event_type":"flexStatusSample","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.commandJsonOutSample":1}`,
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.newrelic.nri-flex",` +
			`"integration_version":"Unknown-SNAPSHOT","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

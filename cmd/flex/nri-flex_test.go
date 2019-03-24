package main

import (
	"encoding/json"
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

// testSamples as samples could be generated in different orders, so we test per sample
func testSamples(expectedSamples []string, entityMetrics []*metric.Set, t *testing.T) {
	if len(entityMetrics) != len(expectedSamples) {
		t.Errorf("Missing samples, got: %v, want: %v.", (entityMetrics), (expectedSamples))

		// t.Errorf("Missing samples, got: %v, want: %v.", len(entityMetrics), len(expectedSamples))
	}
	for _, expectedSample := range expectedSamples {
		matchedSample := false
		for _, sample := range entityMetrics {
			out, err := sample.MarshalJSON()
			if err != nil {
				logger.Flex("debug", err, "failed to marshal", false)
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
	runIntegration()
	expectedSamples := []string{
		`{"commandJsonOutSample_count":1,"configsProcessed":1,"eventCount":1,"eventDropCount":0,"event_type":"flexStatusSample"}`,
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.kav91.nri-flex",` +
			`"integration_version":"0.4.3-pre","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestConfigFile(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestReadJsonCmd", "nri-flex")
	load.Args.ConfigFile = "../../test/configs/json-read-cmd-example.yml"
	runIntegration()
	expectedSamples := []string{
		`{"commandJsonOutSample_count":2,"configsProcessed":1,"eventCount":1,"eventDropCount":0,"event_type":"flexStatusSample"}`,
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.kav91.nri-flex",` +
			`"integration_version":"0.4.3-pre","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

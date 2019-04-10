package main

import (
	"encoding/json"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

// testSamples as samples could be generated in different orders, so we test per sample
func testSamples(expectedSamples []string, entityMetrics []*metric.Set, t *testing.T) {
	if len(entityMetrics) != len(expectedSamples) {
		t.Errorf("Missing samples, got: %v, want: %v.", (entityMetrics), (expectedSamples))
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
		`{"event_type":"flexStatusSample","flex.ConfigsProcessed":1,"flex.EventCount":1,"flex.EventDropCount":0,"flex.commandJsonOutSample":1}`,
<<<<<<< HEAD:cmd/nri-flex/nri-flex_test.go
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.newrelic.nri-flex",` +
			`"integration_version":"Unknown-SNAPSHOT","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
=======
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.kav91.nri-flex",` +
			`"integration_version":"0.5.3-pre","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
>>>>>>> master:cmd/flex/nri-flex_test.go
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestConfigFile(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestReadJsonCmd", "nri-flex")
	load.Args.ConfigFile = "../../test/configs/json-read-cmd-example.yml"
	runIntegration()
	expectedSamples := []string{
		`{"event_type":"flexStatusSample","flex.ConfigsProcessed":1,"flex.EventCount":1,"flex.EventDropCount":0,"flex.commandJsonOutSample":1}`,
<<<<<<< HEAD:cmd/nri-flex/nri-flex_test.go
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.newrelic.nri-flex",` +
			`"integration_version":"Unknown-SNAPSHOT","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
=======
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.kav91.nri-flex",` +
			`"integration_version":"0.5.3-pre","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
>>>>>>> master:cmd/flex/nri-flex_test.go
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

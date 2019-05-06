package logger

import (
	"fmt"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"

	"github.com/newrelic/infra-integrations-sdk/integration"
)

func TestLogger(t *testing.T) {
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestLogger", "nri-flex")
	Flex("debug", nil, "hi", true)
	if len(load.Entity.Metrics) != 1 {
		t.Errorf("Incorrect number of events created %d", len(load.Entity.Metrics))
	} else {
		if load.Entity.Metrics[0].Metrics["event_type"] != "flexDebug" {
			t.Errorf("incorrect event type want: flexDebug, got: %v", fmt.Sprintf("%v", load.Entity.Metrics[0].Metrics["event_type"]))
		}
	}

	load.Refresh()
	load.Entity, _ = i.Entity("TestLogger2", "nri-flex")
	Flex("debug", fmt.Errorf("testing"), "123", true)
	if len(load.Entity.Metrics) != 1 {
		t.Errorf("Incorrect number of events created %d", len(load.Entity.Metrics))
	} else {
		if load.Entity.Metrics[0].Metrics["event_type"] != "flexDebug" {
			t.Errorf("incorrect event type want: flexDebug, got: %v", fmt.Sprintf("%v", load.Entity.Metrics[0].Metrics["event_type"]))
		}
	}

	load.Refresh()
	load.Entity, _ = i.Entity("TestLogger3", "nri-flex")
	Flex("debug", fmt.Errorf("testing"), "123", true)
	if len(load.Entity.Metrics) != 1 {
		t.Errorf("Incorrect number of events created %d", len(load.Entity.Metrics))
	} else {
		if load.Entity.Metrics[0].Metrics["event_type"] != "flexInfo" {
			t.Errorf("incorrect event type want: flexInfo, got: %v", fmt.Sprintf("%v", load.Entity.Metrics[0].Metrics["event_type"]))
		}
	}
}

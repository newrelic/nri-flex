package outputs

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"net/http"
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// SendToInsights - Send processed events to insights
func SendToInsights() {
	modifyEventType()
	postRequest()
	// empty the infrastructure entity metrics by default
	if !load.Args.InsightsOutput {
		load.Entity.Metrics = []*metric.Set{}
	}
}

// modifyEventType insights uses eventType key in camel case whereas infrastructure uses event_type
func modifyEventType() {
	for _, event := range load.Entity.Metrics {
		event.Metrics["eventType"] = event.Metrics["event_type"].(string)
		delete(event.Metrics, "event_type")
	}
}

// postRequest wraps request and attaches needed headers and zlib compression
func postRequest() {
	jsonData, err := json.Marshal(load.Entity.Metrics)
	if err != nil {
		logger.Flex("debug", err, "failed to marshal", false)
	} else {
		var zlibCompressedPayload bytes.Buffer
		w := zlib.NewWriter(&zlibCompressedPayload)
		_, err := w.Write(jsonData)
		logger.Flex("debug", err, "unable to write zlib compressed form", false)
		logger.Flex("debug", w.Close(), "unable to close zlib writer", false)
		if err != nil {
			logger.Flex("debug", fmt.Errorf("failed to compress payload"), "", false)
		} else {
			tr := &http.Transport{IdleConnTimeout: 15 * time.Second}
			client := &http.Client{Transport: tr}
			req, err := http.NewRequest("POST", load.Args.InsightsURL, bytes.NewBuffer(zlibCompressedPayload.Bytes()))
			logger.Flex("info", nil, fmt.Sprintf("insights: bytes %d events %d", len(zlibCompressedPayload.Bytes()), len(load.Entity.Metrics)), false)

			if err != nil {
				logger.Flex("debug", err, "unable to create http.Request", false)
			} else {
				req.Header.Set("Content-Encoding", "deflate")
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Insert-Key", load.Args.InsightsAPIKey)
				_, err := client.Do(req)
				logger.Flex("debug", err, "unable to send", false)
			}
		}
	}
}

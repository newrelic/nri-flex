/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

func TestPrometheusRedis(t *testing.T) {
	load.Refresh()

	// create a listener with desired port
	l, err := net.Listen("tcp", "127.0.0.1:9122")
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fileData, err := ioutil.ReadFile("../../test/payloads/prometheusRedis.out")
		if err != nil {
			t.Fatal(err)
		}
		_, err = rw.Write(fileData)
		if err != nil {
			t.Fatal(err)
		}
		logger.Flex("debug", err, "failed to write", false)
	}))

	// NewUnstartedServer creates a listener. Close listener and replace with the one we created.
	ts.Listener.Close()
	ts.Listener = l
	// Start the server.
	ts.Start()

	config := load.Config{
		APIs: []load.API{
			{
				Name: "redis",
				URL:  "http://localhost:9122",
				Prometheus: load.Prometheus{
					Enable: true,
					CustomAttributes: map[string]string{
						"abc": "def",
					},
				},
			},
		},
	}

	var jsonOut interface{}
	expectedOutput, err := ioutil.ReadFile("../../test/payloadsExpected/promRedis.json")
	if err != nil {
		t.Fatal(err)
	}
	json.Unmarshal(expectedOutput, &jsonOut)
	expectedDatastore := jsonOut.([]interface{})

	doLoop := true
	dataStore := []interface{}{}
	RunHTTP(&dataStore, &doLoop, &config, config.APIs[0], &config.APIs[0].URL)

	if len(expectedDatastore) != len(dataStore) {
		t.Errorf("Incorrect number of samples generated expected: %d, got: %d", len(expectedDatastore), len(dataStore))
		t.Errorf("%v", (dataStore))
	}

	for _, sample := range expectedDatastore {
		switch sample := sample.(type) {
		case map[string]interface{}:
			for _, rSample := range dataStore {
				switch recSample := rSample.(type) {
				case map[string]interface{}:

					if fmt.Sprintf("%v", recSample["name"]) == "main" && fmt.Sprintf("%v", sample["name"]) == "main" {
						for key := range sample {
							if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
								t.Errorf("%v want %v, got %v", key, sample[key], recSample[key])
							}
						}
					}

					if fmt.Sprintf("%v", recSample["db"]) != "<nil>" && fmt.Sprintf("%v", recSample["db"]) == fmt.Sprintf("%v", sample["db"]) {
						for key := range sample {
							if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
								t.Errorf("dbSample %v want %v, got %v", key, sample[key], recSample[key])
							}
						}
					}
				}
			}
		}
	}
}

func TestPrometheusNginx(t *testing.T) {
	load.Refresh()

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fileData, _ := ioutil.ReadFile("../../test/payloads/prometheusNginx.out")
		_, err := rw.Write(fileData)
		logger.Flex("debug", err, "failed to write", false)
	}))
	// Close the server when test finishes
	defer server.Close()

	config := load.Config{
		APIs: []load.API{
			{
				Name: "nginxIngress",
				URL:  server.URL,
				Prometheus: load.Prometheus{
					Enable: true,
				},
			},
		},
	}

	var jsonOut interface{}
	expectedOutput, _ := ioutil.ReadFile("../../test/payloadsExpected/promNginxFull1.json")
	json.Unmarshal(expectedOutput, &jsonOut)
	expectedDatastore := jsonOut.([]interface{})

	doLoop := true
	dataStore := []interface{}{}
	RunHTTP(&dataStore, &doLoop, &config, config.APIs[0], &config.APIs[0].URL)

	if len(expectedDatastore) != len(dataStore) {
		t.Errorf("Incorrect number of samples generated expected: %d, got: %d", len(expectedDatastore), len(dataStore))
		t.Errorf("%v", (dataStore))
	}

	for _, sample := range expectedDatastore {
		switch sample := sample.(type) {
		case map[string]interface{}:
			for _, rSample := range dataStore {
				switch recSample := rSample.(type) {
				case map[string]interface{}:
					if fmt.Sprintf("%v", recSample["name"]) == "main" && fmt.Sprintf("%v", sample["name"]) == "main" {
						for key := range sample {
							if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
								t.Errorf("%v want %v, got %v", key, sample[key], recSample[key])
							}
						}
					}
				}
			}
		}
	}
}

func TestPrometheusNginx2(t *testing.T) {
	load.Refresh()

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fileData, _ := ioutil.ReadFile("../../test/payloads/prometheusNginx.out")
		_, err := rw.Write(fileData)
		logger.Flex("debug", err, "failed to write", false)
	}))
	// Close the server when test finishes
	defer server.Close()

	config := load.Config{
		APIs: []load.API{
			{
				Name: "nginxIngress",
				URL:  server.URL,
				Prometheus: load.Prometheus{
					Enable:    true,
					Histogram: true,
					Summary:   true,
				},
			},
		},
	}

	var jsonOut interface{}
	expectedOutput, _ := ioutil.ReadFile("../../test/payloadsExpected/promNginxFull2.json")
	json.Unmarshal(expectedOutput, &jsonOut)
	expectedDatastore := jsonOut.([]interface{})

	doLoop := true
	dataStore := []interface{}{}
	RunHTTP(&dataStore, &doLoop, &config, config.APIs[0], &config.APIs[0].URL)

	if len(expectedDatastore) != len(dataStore) {
		t.Errorf("Incorrect number of samples generated expected: %d, got: %d", len(expectedDatastore), len(dataStore))
		t.Errorf("%v", (dataStore))
	}

	for _, sample := range expectedDatastore {
		switch sample := sample.(type) {
		case map[string]interface{}:
			for _, rSample := range dataStore {
				switch recSample := rSample.(type) {
				case map[string]interface{}:

					if fmt.Sprintf("%v", recSample["name"]) == "main" && fmt.Sprintf("%v", sample["name"]) == "main" {
						for key := range sample {
							if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
								t.Errorf("%v want %v, got %v", key, sample[key], recSample[key])
							}
						}
					}
				}
			}
		}
	}
}

func TestPrometheusNginx3(t *testing.T) {
	load.Refresh()
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fileData, _ := ioutil.ReadFile("../../test/payloads/prometheusNginx.out")
		_, err := rw.Write(fileData)
		logger.Flex("debug", err, "failed to write", false)
	}))
	// Close the server when test finishes
	defer server.Close()

	config := load.Config{
		APIs: []load.API{
			{
				Name: "nginxIngress",
				URL:  server.URL,
				Prometheus: load.Prometheus{
					Enable:    true,
					Unflatten: true,
				},
			},
		},
	}

	var jsonOut interface{}
	expectedOutput, _ := ioutil.ReadFile("../../test/payloadsExpected/promNginxFull3.json")
	json.Unmarshal(expectedOutput, &jsonOut)
	expectedDatastore := jsonOut.([]interface{})

	doLoop := true
	dataStore := []interface{}{}
	RunHTTP(&dataStore, &doLoop, &config, config.APIs[0], &config.APIs[0].URL)

	if len(expectedDatastore) != len(dataStore) {
		t.Errorf("Incorrect number of samples generated expected: %d, got: %d", len(expectedDatastore), len(dataStore))
		t.Errorf("%v", (dataStore))
	}
}

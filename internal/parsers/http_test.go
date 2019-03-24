package parser

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
	"testing"
)

func TestRunHTTP(t *testing.T) {
	// create a listener with desired port
	l, _ := net.Listen("tcp", "127.0.0.1:9123")
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		fileData, _ := ioutil.ReadFile("../../test/payloadsExpected/httpTest.json")
		_, err := rw.Write(fileData)
		logger.Flex("debug", err, "failed to write", false)
	}))
	// NewUnstartedServer creates a listener. Close listener and replace with the one we created.
	ts.Listener.Close()
	ts.Listener = l
	// Start the server.
	ts.Start()

	doLoop := true
	dataStore := []interface{}{}
	config := load.Config{
		Name: "httpExample",
		Global: load.Global{
			BaseURL: "http://127.0.0.1:9123",
			Timeout: 5000,
			User:    "batman",
			Pass:    "robin",
			Headers: map[string]string{
				"test": "abc",
			},
		},
		APIs: []load.API{
			{
				EventType: "httpExample",
				URL:       "/",
				Timeout:   5100,
				User:      "batman",
				Pass:      "robin",
				Headers: map[string]string{
					"test2": "abc2",
				},
			},
			{
				EventType: "httpExample2",
				URL:       "todos/abc",
				Timeout:   5100,
				User:      "batman",
				Pass:      "robin",
				Headers: map[string]string{
					"test2": "abc2",
				},
			},
			{
				EventType: "httpExample3",
				URL:       "users",
				Timeout:   5100,
				User:      "batman",
				Pass:      "robin",
				Headers: map[string]string{
					"test2": "abc2",
				},
			},
		},
	}

	expectedSamples := []interface{}{
		map[string]interface{}{
			"userId":         1,
			"id":             1,
			"title":          "delectus aut autem",
			"completed":      "false",
			"api.StatusCode": 200,
		},
	}
	RunHTTP(&doLoop, &config, config.APIs[0], &config.APIs[0].URL, &dataStore)

	if len(dataStore) != len(expectedSamples) {
		t.Errorf("received sample count %d does not match expected %d", len(dataStore), len(expectedSamples))
		t.Errorf("%v", dataStore)
	}

	for key := range dataStore[0].(map[string]interface{}) {
		if fmt.Sprintf("%v", dataStore[0].(map[string]interface{})[key]) != fmt.Sprintf("%v", expectedSamples[0].(map[string]interface{})[key]) {
			t.Errorf(fmt.Sprintf("doesnt match %v : %v - %v", key, dataStore[0].(map[string]interface{})[key], expectedSamples[0].(map[string]interface{})[key]))
		}
	}
}

package parser

import (
	"fmt"
	"nri-flex/internal/load"
	"testing"
)

func TestRunHTTP(t *testing.T) {
	doLoop := true
	dataStore := []interface{}{}
	config := load.Config{
		Name: "httpExample",
		Global: load.Global{
			BaseURL: "https://jsonplaceholder.typicode.com/",
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
				URL:       "todos/1",
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
	}

	for key := range dataStore[0].(map[string]interface{}) {
		if fmt.Sprintf("%v", dataStore[0].(map[string]interface{})[key]) != fmt.Sprintf("%v", expectedSamples[0].(map[string]interface{})[key]) {
			t.Errorf(fmt.Sprintf("doesnt match %v : %v - %v", key, dataStore[0].(map[string]interface{})[key], expectedSamples[0].(map[string]interface{})[key]))
		}
	}
}

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestRunHttp(t *testing.T) {
	tests := map[string]struct {
		address          string
		config           load.Config
		expected         []interface{}
		expectedFilePath string
	}{
		"base-sample": {
			"127.0.0.1:9123",
			load.Config{
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
			},
			[]interface{}{
				map[string]interface{}{
					"userId":         float64(1),
					"id":             float64(1),
					"title":          "delectus aut autem",
					"completed":      "false",
					"api.StatusCode": 200,
				},
			},
			"../../test/payloadsExpected/http_response-single_object.json",
		},
		"sample-with-headers-single-response-object": {
			"127.0.0.1:9124",
			load.Config{
				Name: "return-headers-example",
				Global: load.Global{
					BaseURL: "http://127.0.0.1:9124",
				},
				APIs: []load.API{
					{
						EventType: "return-headers-example",
						URL:       "/",
						Timeout:   5100,
						ReturnHeaders: true,
					},
				},
			},
			[]interface{}{
				map[string]interface{}{
					"userId":                    float64(1),
					"id":                        float64(1),
					"title":                     "delectus aut autem",
					"completed":                 "false",
					"api.StatusCode":            200,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"127"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			"../../test/payloadsExpected/http_response-single_object.json",
		},
		"sample-with-headers-multiple-response-object": {
			"127.0.0.1:9125",
			load.Config{
				Name: "return-headers-example",
				Global: load.Global{
					BaseURL: "http://127.0.0.1:9125",
				},
				APIs: []load.API{
					{
						EventType:     "return-headers-example",
						URL:           "/",
						Timeout:       5100,
						ReturnHeaders: true,
					},
				},
			},
			[]interface{}{
				map[string]interface{}{
					"userId":                    float64(1),
					"id":                        float64(1),
					"title":                     "delectus aut autem",
					"completed":                 "false",
					"api.StatusCode":            200,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"216"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
				map[string]interface{}{
					"userId":                    float64(1),
					"id":                        float64(2),
					"title":                     "quis ut nam facilis et officia qui",
					"completed":                 "false",
					"api.StatusCode":            200,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"216"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			"../../test/payloadsExpected/http_response-multiple_objects.json",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			load.Refresh()
			doLoop := true

			err := mockHttpServer(tc.address, tc.expectedFilePath)
			assert.NoError(t, err)

			var dataStore []interface{}
			RunHTTP(&dataStore, &doLoop, &tc.config, tc.config.APIs[0], &tc.config.APIs[0].URL)
			assertElementsMatch(t, dataStore, tc)
		})
	}
}

func assertElementsMatch(t *testing.T, dataStore []interface{}, tc struct {
	address          string
	config           load.Config
	expected         []interface{}
	expectedFilePath string
}) {

	for index, result := range dataStore {
		for key := range result.(map[string]interface{}) {
			a := result.(map[string]interface{})[key]
			e := tc.expected[index].(map[string]interface{})[key]

			if fmt.Sprintf("%v", a) != fmt.Sprintf("%v", e) || reflect.TypeOf(a) != reflect.TypeOf(e) {
				t.Errorf(fmt.Sprintf("mismatch in '%v' key: expected value %v(%v) - actual value %v(%v)", key, e, reflect.TypeOf(e).String(), a,reflect.TypeOf(a).String()))
			}
		}
	}
	assert.ElementsMatch(t, dataStore, tc.expected)
}

func mockHttpServer(url string, filePath string) error {
	l, err := net.Listen("tcp", url)
	if err != nil {
		load.Logrus.WithError(err).Error("http: failed to create listener")
		return err
	}

	mockHttpHandler := mockHttpHandler{
		filePath: filePath,
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(mockHttpHandler.ServeHTTP))

	// NewUnstartedServer creates a listener. Close listener and replace with the one we created.
	ts.Listener.Close()
	ts.Listener = l
	// Start the server.
	ts.Start()
	return nil
}

type mockHttpHandler struct {
	filePath string
}

func (h *mockHttpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Date", "Mon, 18 May 2020 09:38:35 GMT")
	fileData, _ := ioutil.ReadFile(h.filePath)
	_, err := rw.Write(fileData)
	if err != nil {
		load.Logrus.WithError(err).Error("http: failed to write")
	}
}

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"github.com/stretchr/testify/assert"

	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
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
						ReturnHeaders: false,
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
						ReturnHeaders: false,
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
						ReturnHeaders: false,
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
			"../../test/payloadsExpected/http-response_single-object.json",
		},
		"sample-with-headers": {
			"127.0.0.1:9124",
			load.Config{
				Name: "return-headers-example",
				Global: load.Global{
					BaseURL: "http://127.0.0.1:9124",
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
					"api.header.Content-Length": []string{"154"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			"../../test/payloadsExpected/http-response_single-object.json",
		},
		"sample-without-headers": {
			"127.0.0.1:9125",
			load.Config{
				Name: "return-headers-example",
				Global: load.Global{
					BaseURL: "http://127.0.0.1:9125",
				},
				APIs: []load.API{
					{
						EventType: "return-headers-example",
						URL:       "/",
						Timeout:   5100,
						//returnHeaders: true
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
			"../../test/payloadsExpected/http-response_single-object.json",
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
			assert.ElementsMatch(t, dataStore, tc.expected)
		})
	}
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
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("http: failed to write")
	}
}

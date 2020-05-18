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
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

func TestRunHttp(t *testing.T) {
	portNumber := "9123"
	initializeHttpListener(portNumber)

	tests := map[string]struct {
		config   load.Config
		expected []interface{}
	}{
		"base-sample": {
			load.Config{
				Name: "httpExample",
				Global: load.Global{
					BaseURL: fmt.Sprintf("http://127.0.0.1:%v", portNumber),
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
		},
		"sample-with-headers": {
			load.Config{
				Name: "return-headers-example",
				Global: load.Global{
					BaseURL: fmt.Sprintf("http://127.0.0.1:%v", portNumber),
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
		},
		"sample-without-headers": {
			load.Config{
				Name: "return-headers-example",
				Global: load.Global{
					BaseURL: fmt.Sprintf("http://127.0.0.1:%v", portNumber),
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
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			load.Refresh()
			doLoop := true

			var dataStore []interface{}
			RunHTTP(&dataStore, &doLoop, &tc.config, tc.config.APIs[0], &tc.config.APIs[0].URL)
			assert.ElementsMatch(t, dataStore, tc.expected)
		})
	}
}

func initializeHttpListener(port string) {
	// create a listener with desired port
	address := fmt.Sprintf("127.0.0.1:%v", port)
	l, _ := net.Listen("tcp", address)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Date", "Mon, 18 May 2020 09:38:35 GMT")
		fileData, _ := ioutil.ReadFile("../../test/payloadsExpected/httpTest.json")
		_, err := rw.Write(fileData)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("http: failed to write")
		}
	}))
	// NewUnstartedServer creates a listener. Close listener and replace with the one we created.
	ts.Listener.Close()
	ts.Listener = l
	// Start the server.
	ts.Start()
}

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
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
			"../../test/payloadsExpected/http_response/http_response-single_object.json",
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
					"api.header.Content-Length": []string{"127"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			"../../test/payloadsExpected/http_response/http_response-single_object.json",
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
			"../../test/payloadsExpected/http_response/http_response-multiple_objects.json",
		},
		"sample-with-headers-string-response": {
			"127.0.0.1:9126",
			load.Config{
				Name: "return-headers-example",
				Global: load.Global{
					BaseURL: "http://127.0.0.1:9126",
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
					"output":                    "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua",
					"api.StatusCode":            200,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"136"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			"../../test/payloadsExpected/http_response/http_response-string_response.json",
		},
		"sample-error-with-headers": {
			"127.0.0.1:9127",
			load.Config{
				Name: "return-headers-example",
				Global: load.Global{
					BaseURL: "http://127.0.0.1:9127",
				},
				APIs: []load.API{
					{
						EventType:     "return-headers-example",
						URL:           "/",
						Timeout:       5100,
						ReturnHeaders: true,
					},
					{
						EventType:     "return-headers-example",
						URL:           "/todos",
						Timeout:       5100,
						ReturnHeaders: true,
					},
				},
			},
			[]interface{}{
				map[string]interface{}{
					"error":                     "Missing Required Parameters",
					"api.StatusCode":            200,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"52"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			"../../test/payloadsExpected/http_response/http_response-error_message.json",
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

func TestHttp_handleJSON_unmarshalError(t *testing.T) {
	// Given a test logger
	load.Logrus.SetOutput(ioutil.Discard)  // discard logs so not to break race tests
	defer log.SetOutput(os.Stderr) // return back to default
	hook := new(test.Hook)
	load.Logrus.AddHook(hook)

	// AND a http response with invalid json format
	url := ""
	doLoop := false
	var resp gorequest.Response
	var sample []interface{}

	body := []byte("invalidJSON")

	// When get body from response
	handleJSON(&sample, body, &resp, &doLoop, &url, "", false)

	// THEN one error entry were found
	require.NotEmpty(t, hook.Entries)
	entry := hook.LastEntry()
	assert.Equal(t, logrus.ErrorLevel, entry.Level)
	assert.Equal(t, "http: failed to unmarshal json", entry.Message)
	assert.EqualError(t, entry.Data["error"].(error), "invalid character 'i' looking for beginning of value")

	// AND no sample was built
	assert.Equal(t, 0, len(sample))
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
				t.Errorf(fmt.Sprintf("mismatch in '%v' key: expected value %v(%v) - actual value %v(%v)", key, e, reflect.TypeOf(e).String(), a, reflect.TypeOf(a).String()))
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

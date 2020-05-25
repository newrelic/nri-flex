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
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestRunHttp(t *testing.T) {
	successStatusCode := 200
	internalServerErrorStatusCode := 500

	tests := map[string]struct {
		config             load.Config
		expected           []interface{}
		expectedFilePath   string
		expectedStatusCode int
	}{
		"base-sample": {
			load.Config{
				Name: "httpExample",
				Global: load.Global{
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
					"api.StatusCode": successStatusCode,
				},
			},
			path.Join("..", "..", "test", "payloadsExpected", "http_response", "http_response-single_object.json"),
			successStatusCode,
		},
		"sample-with-headers-single-response-object": {
			load.Config{
				Name: "return-headers-example",
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
					"api.StatusCode":            successStatusCode,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"999"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			path.Join("..", "..", "test", "payloadsExpected", "http_response", "http_response-single_object.json"),
			successStatusCode,
		},
		"sample-with-headers-multiple-response-object": {
			load.Config{
				Name: "return-headers-example",
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
					"api.StatusCode":            successStatusCode,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"999"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
				map[string]interface{}{
					"userId":                    float64(1),
					"id":                        float64(2),
					"title":                     "quis ut nam facilis et officia qui",
					"completed":                 "false",
					"api.StatusCode":            successStatusCode,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"999"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			path.Join("..", "..", "test", "payloadsExpected", "http_response", "http_response-multiple_objects.json"),
			successStatusCode,
		},
		"sample-with-headers-string-response": {
			load.Config{
				Name: "return-headers-example",
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
					"api.StatusCode":            successStatusCode,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"999"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			path.Join("..", "..", "test", "payloadsExpected", "http_response", "http_response-string_response.json"),
			successStatusCode,
		},
		"sample-error-with-headers": {
			load.Config{
				Name: "return-headers-example",
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
					"api.StatusCode":            internalServerErrorStatusCode,
					"api.header.Content-Type":   []string{"application/json"},
					"api.header.Content-Length": []string{"999"},
					"api.header.Date":           []string{"Mon, 18 May 2020 09:38:35 GMT"},
					"api.header.Retry-Count":    []string{"0"},
				},
			},
			path.Join("..", "..", "test", "payloadsExpected", "http_response", "http_response-error_message.json"),
			internalServerErrorStatusCode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			load.Refresh()
			doLoop := true

			ts := mockHttpServer(tc.expectedFilePath, tc.expectedStatusCode)
			defer ts.Close()

			tc.config.Global.BaseURL = ts.URL

			var dataStore []interface{}
			RunHTTP(&dataStore, &doLoop, &tc.config, tc.config.APIs[0], &tc.config.APIs[0].URL)
			assertElementsMatch(t, dataStore, tc.expected)
		})
	}
}

func TestHttp_handleJSON_unmarshalError(t *testing.T) {
	// Given a test logger
	load.Logrus.SetOutput(ioutil.Discard) // discard logs so not to break race tests
	defer log.SetOutput(os.Stderr)        // return back to default
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

func assertElementsMatch(t *testing.T, actual []interface{}, expected []interface{}) {
	for index, result := range actual {
		for key := range result.(map[string]interface{}) {
			a := result.(map[string]interface{})[key]
			e := expected[index].(map[string]interface{})[key]

			if fmt.Sprintf("%v(%T)", a, a) != fmt.Sprintf("%v(%T)", e, e) {
				t.Errorf(fmt.Sprintf("mismatch in '%v' key: expected value %v(%T) - actual value %v(%T)", key, e, e, a, a))
			}
		}
	}
	assert.ElementsMatch(t, actual, expected)
}

func mockHttpServer(filePath string, statusCode int) *httptest.Server {
	mockHttpHandler := mockHttpHandler{
		filePath:   filePath,
		statusCode: statusCode,
	}

	return httptest.NewServer(http.HandlerFunc(mockHttpHandler.ServeHTTP))
}

type mockHttpHandler struct {
	filePath   string
	statusCode int
}

func (h *mockHttpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Content-Length", "999")
	rw.Header().Set("Date", "Mon, 18 May 2020 09:38:35 GMT")
	rw.WriteHeader(h.statusCode)
	fileData, _ := ioutil.ReadFile(h.filePath)
	_, err := rw.Write(fileData)
	if err != nil {
		load.Logrus.WithError(err).Error("http: failed to write")
	}
}

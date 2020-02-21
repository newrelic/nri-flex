/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */
package processor

import (
	"encoding/json"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/assert"
)

func TestSplitObjects(t *testing.T) {

	getAPI := func() *load.API {
		return &load.API{
			SplitObjects: true,
		}
	}

	testCases := map[string]struct {
		data           map[string]interface{}
		expectedData   string // SplitObjects changes also the input map.
		expectedResult string
	}{
		"SimpleMap": {
			data: map[string]interface{}{
				"abc": 1,
			},
			expectedData:   `{"abc":1}`,
			expectedResult: `null`,
		},
		"NestedMaps": {
			data: map[string]interface{}{
				"abc": map[string]interface{}{
					"def": 1,
				},
			},
			expectedData:   `{"abc":{"def":1,"split.id":"abc"}}`,
			expectedResult: `[{"def":1,"split.id":"abc"}]`,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			api := getAPI()
			result := splitObjects(&testCase.data, api)

			got, _ := json.Marshal(result)
			assert.Equal(t, testCase.expectedResult, string(got))
			// Running splitObjects should update SplitObjects from config.
			assert.False(t, api.SplitObjects)

			got2, _ := json.Marshal(testCase.data)
			assert.Equal(t, testCase.expectedData, string(got2))
		})
	}
}

func TestSplitArrays(t *testing.T) {
	getAPI := func() *load.API {
		return &load.API{
			SplitObjects: true,
			SetHeader: []string{
				"TIMESTAMP", "HOST_ID", "HOSTNAME", "PERCENT_USED",
			},
		}
	}

	inputArray := []interface{}{
		[]interface{}{1582159853733, 1, "host1", 10},
		[]interface{}{1582159853733, 2, "host2", 20},
		[]interface{}{1582159853733, 3, "host3", 30},
	}

	expectedResult := `[{"HOSTNAME":"host1","HOST_ID":"1","PERCENT_USED":"10","TIMESTAMP":"1582159853733"},{"HOSTNAME":"host2","HOST_ID":"2","PERCENT_USED":"20","TIMESTAMP":"1582159853733"},{"HOSTNAME":"host3","HOST_ID":"3","PERCENT_USED":"30","TIMESTAMP":"1582159853733"}]`

	api := getAPI()
	result := splitArrays(&inputArray, map[string]interface{}{}, "", api, &[]interface{}{})
	got, _ := json.Marshal(result)
	assert.Equal(t, expectedResult, string(got))

}

func TestRunLazyFlatten(t *testing.T) {
	getConfig := func(lazyFlatten []string) *load.Config {
		return &load.Config{
			APIs: []load.API{
				0: {
					LazyFlatten: lazyFlatten,
				},
			},
		}
	}

	testCases := map[string]struct {
		sample      map[string]interface{}
		lazyFlatten []string
		expected    string
	}{
		"SimpleMatchWithArray": {
			sample: map[string]interface{}{
				"abc": []interface{}{
					0: map[string]interface{}{
						"def": map[string]interface{}{
							"x": 1,
						},
					},
				},
				"xyz": []interface{}{
					0: map[string]interface{}{
						"def": map[string]interface{}{
							"x": 1,
						},
					},
				},
			},
			lazyFlatten: []string{
				"abc",
			},
			expected: `{"abc":{"flat.0.def.x":1},"xyz":[{"def":{"x":1}}]}`,
		},
		"SimpleMatchWithMap": {
			sample: map[string]interface{}{
				"abc": map[string]interface{}{
					"def": map[string]interface{}{
						"x": 1,
					},
				},
				"xyz": map[string]interface{}{
					"def": map[string]interface{}{
						"x": 1,
					},
				},
			},
			lazyFlatten: []string{
				"abc",
			},
			expected: `{"abc":{"flat.def.x":1},"xyz":{"def":{"x":1}}}`,
		},
		"CascadeMatchWithArray": {
			sample: map[string]interface{}{
				"abc": []interface{}{
					0: map[string]interface{}{
						"def": map[string]interface{}{
							"x": 1,
						},
					},
				},
				"xyz": []interface{}{
					0: map[string]interface{}{
						"def": map[string]interface{}{
							"x": 1,
						},
					},
				},
			},
			lazyFlatten: []string{
				"abc>def",
			},
			expected: `{"abc":[{"def.x":1}],"xyz":[{"def":{"x":1}}]}`,
		},
		"CascadeMatchWithMap": {
			sample: map[string]interface{}{
				"abc": map[string]interface{}{
					"def": map[string]interface{}{
						"x": 1,
					},
				},
				"xyz": map[string]interface{}{
					"def": map[string]interface{}{
						"x": 1,
					},
				},
			},
			lazyFlatten: []string{
				"abc>def",
			},
			expected: `{"abc":{"def":{"def.x":1}},"xyz":{"def":{"x":1}}}`,
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunLazyFlatten(&testCase.sample, getConfig(testCase.lazyFlatten), 0)

			got, _ := json.Marshal(testCase.sample)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

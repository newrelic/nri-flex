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
	// test SplitArray
	getAPI := func() *load.API {
		return &load.API{
			SetHeader: []string{
				"TIMESTAMP", "VALUE", "AA",
			},
		}
	}

	inputArray := []interface{}{[]interface{}{
		map[string]interface{}{
			"samples": []interface{}{
				[]interface{}{uint64(1582513500000), 303.6965867733333},
				[]interface{}{uint64(1582513500000), 404.6965867733333},
			},
			"type":    "B_BITS",
			"subtype": "Out",
			"unit": map[string]interface{}{
				"name":     "aaaMbps",
				"fullName": "aaaMegabits per second",
			},
		},
	},
		[]interface{}{
			map[string]interface{}{
				"samples": []interface{}{
					[]interface{}{uint64(1582513500000), 1000.0},
					[]interface{}{uint64(1582513500000), 2000.0},
				},
				"type":    "B_BITS",
				"subtype": "Configured speed",
				"unit": map[string]interface{}{
					"name":     "bbbMbps",
					"fullName": "bbbMegabits per second",
				},
			},
		},
		[]interface{}{
			map[string]interface{}{
				"samples": []interface{}{
					[]interface{}{uint64(1582513500000), 260.34207970666665},
					[]interface{}{uint64(1582513500000), 370.34207970666665},
				},
				"type":    "B_BITS",
				"subtype": "In",
				"unit": map[string]interface{}{
					"name":     "cccMbps",
					"fullName": "cccMegabits per second",
				},
			},
		},
	}

	expectedResult := `[{"TIMESTAMP":"1582513500000","VALUE":"303.696587","data-subtype":"Out","data-type":"B_BITS","data-unit-fullName":"aaaMegabits per second","data-unit-name":"aaaMbps"},{"TIMESTAMP":"1582513500000","VALUE":"404.696587","data-subtype":"Out","data-type":"B_BITS","data-unit-fullName":"aaaMegabits per second","data-unit-name":"aaaMbps"},{"TIMESTAMP":"1582513500000","VALUE":"1000.000000","data-subtype":"Configured speed","data-type":"B_BITS","data-unit-fullName":"bbbMegabits per second","data-unit-name":"bbbMbps"},{"TIMESTAMP":"1582513500000","VALUE":"2000.000000","data-subtype":"Configured speed","data-type":"B_BITS","data-unit-fullName":"bbbMegabits per second","data-unit-name":"bbbMbps"},{"TIMESTAMP":"1582513500000","VALUE":"260.342080","data-subtype":"In","data-type":"B_BITS","data-unit-fullName":"cccMegabits per second","data-unit-name":"cccMbps"},{"TIMESTAMP":"1582513500000","VALUE":"370.342080","data-subtype":"In","data-type":"B_BITS","data-unit-fullName":"cccMegabits per second","data-unit-name":"cccMbps"}]`

	api := getAPI()
	result := splitArrays(&inputArray, map[string]interface{}{}, "data", api, &[]interface{}{}, map[string]interface{}{})

	got, _ := json.Marshal(result)
	assert.Equal(t, expectedResult, string(got))

	// test SplitArray with leaf_array

	getAPI2 := func() *load.API {
		return &load.API{
			SetHeader: []string{
				"TIMESTAMP", "VALUE", "AA",
			},
			LeafArray: true,
		}
	}

	inputArray2 := []interface{}{[]interface{}{
		map[string]interface{}{
			"timestamps": []interface{}{
				uint64(1585662957000), uint64(1585662958000), uint64(1585662959000),
			},
			"type": "time_series",
		},
	},

		[]interface{}{},
	}

	expectedResult2 := `[{"TIMESTAMP":"1585662957000","concurrent_plays-type":"time_series","index":0},{"TIMESTAMP":"1585662958000","concurrent_plays-type":"time_series","index":1},{"TIMESTAMP":"1585662959000","concurrent_plays-type":"time_series","index":2}]`

	api2 := getAPI2()
	result2 := splitArrays(&inputArray2, map[string]interface{}{}, "concurrent_plays", api2, &[]interface{}{}, map[string]interface{}{})

	got2, _ := json.Marshal(result2)
	assert.Equal(t, expectedResult2, string(got2))

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

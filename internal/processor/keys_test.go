/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestKeyConversionToLower(t *testing.T) {
	api := load.API{
		Name:    "testKeyConversion",
		ToLower: true,
	}
	key := "Hi ThEre!"
	expected := "hi there!"
	v := 10

	RunKeyConversion(&key, api, v, &api.SkipProcessing)
	if key != expected {
		t.Errorf("want: %v got: %v", expected, key)
	}
}

func TestKeyConversionSnakeToCamel(t *testing.T) {
	api := load.API{
		Name:         "testKeyConversion",
		ToLower:      true,
		SnakeToCamel: true,
	}

	key := "hi_there_mario"
	expected := "hiThereMario"
	v := 10

	RunKeyConversion(&key, api, v, &api.SkipProcessing)
	if key != expected {
		t.Errorf("want: %v got: %v", expected, key)
	}
}

func TestKeyConversionConvertSpace(t *testing.T) {
	api := load.API{
		Name:         "testKeyConversion",
		ConvertSpace: "-",
	}

	key := "hi there mario"
	expected := "hi-there-mario"
	v := 10

	RunKeyConversion(&key, api, v, &api.SkipProcessing)
	if key != expected {
		t.Errorf("want: %v got: %v", expected, key)
	}
}

func TestKeepKeys(t *testing.T) {
	getSample := func() map[string]interface{} {
		return map[string]interface{}{
			"abc":   1,
			"xyz":   2,
			"dfadf": "dafsdfa",
			"def":   "initely busted",
		}
	}

	testCases := map[string]struct {
		keepKeys []string
		key      string
		sample   map[string]interface{}
		expected string
	}{
		"WithStrings": {
			keepKeys: []string{"abc", "xyz"},
			key:      "def",
			sample:   getSample(),
			expected: `{"abc":1,"dfadf":"dafsdfa","xyz":2}`,
		},
		"WithRegex": {
			keepKeys: []string{".*?", "xyz"}, // This should capture any key and avoid deleting it.
			key:      "abc",
			sample:   getSample(),
			expected: `{"abc":1,"def":"initely busted","dfadf":"dafsdfa","xyz":2}`,
		},
		"WithEmptyKeepKeysArray": {
			keepKeys: []string{},
			key:      "def",
			sample:   getSample(),
			expected: `{"abc":1,"def":"initely busted","dfadf":"dafsdfa","xyz":2}`,
		},
		"WithEmptyKey": {
			keepKeys: []string{},
			key:      "",
			sample:   getSample(),
			expected: `{"abc":1,"def":"initely busted","dfadf":"dafsdfa","xyz":2}`,
		},
		"WithNilSample": {
			keepKeys: []string{},
			key:      "",
			sample:   nil,
			expected: `null`,
		},
	}

	for caseName, testCase := range testCases {
		t.Run(caseName, func(t *testing.T) {
			RunKeepKeys(testCase.keepKeys, &testCase.key, &testCase.sample)
			actual, _ := json.Marshal(testCase.sample)
			assert.Equal(t, testCase.expected, string(actual))
		})
	}
}

func TestKeyRemover(t *testing.T) {

	testCases := map[string]struct {
		removeKeys []string
		sample     map[string]interface{}
		expected   string
	}{
		"LowerCase": {
			removeKeys: []string{"abc", "xyz"},
			sample: map[string]interface{}{
				"abc":   1,
				"xyz":   2,
				"dfadf": "dafsdfa",
			},
			expected: `{"dfadf":"dafsdfa"}`,
		},
		"CaseInsensitive": {
			removeKeys: []string{"ABc", "xYz"},
			sample: map[string]interface{}{
				"abc":   1,
				"ABC":   2, // Same content of key should also be removed.
				"xyz":   3,
				"dfadf": "dafsdfa",
			},
			expected: `{"dfadf":"dafsdfa"}`,
		},
		"UnknownKey": {
			removeKeys: []string{"aaaaa"},
			sample: map[string]interface{}{
				"abc":   1,
				"xyz":   3,
				"dfadf": "dafsdfa",
			},
			expected: `{"abc":1,"dfadf":"dafsdfa","xyz":3}`,
		},
		"NilRemoveKeys": {
			removeKeys: nil,
			sample: map[string]interface{}{
				"abc":   1,
				"xyz":   3,
				"dfadf": "dafsdfa",
			},
			expected: `{"abc":1,"dfadf":"dafsdfa","xyz":3}`,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunKeyRemover(&testCase.sample, testCase.removeKeys)
			got, _ := json.Marshal(testCase.sample)

			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

func TestKeyRenamer(t *testing.T) {
	testCases := map[string]struct {
		renameKeys map[string]string
		key        string
		expected   string
	}{
		"RegexWithSuffix": {
			renameKeys: map[string]string{"abc$": "xyz"},
			key:        "aaaaabc",
			expected:   `aaaaxyz`,
		},
		"MultipleMatches": { // TODO: is this a valid case?
			renameKeys: map[string]string{"ab": "xyz"},
			key:        "abab",
			expected:   `xyzxyz`,
		},
		"StringMatch": {
			//TODO: with current RunKeyRenamer implementation this will also match parentX1.uptime. Should this be treated as a string?
			renameKeys: map[string]string{"parent.1": "aws"},
			key:        "parent.1.uptime",
			expected:   `aws.uptime`,
		},
		"NoMatch": {
			renameKeys: map[string]string{"asd": "xyz"},
			key:        "aaaaabc",
			expected:   `aaaaabc`,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunKeyRenamer(testCase.renameKeys, &testCase.key)
			assert.Equal(t, testCase.expected, testCase.key)
		})
	}
}

func TestStripKeys(t *testing.T) {
	testCases := map[string]struct {
		sample    map[string]interface{}
		stripKeys []string
		expected  string
	}{
		"NoChangeIfNoMatch": {
			sample: map[string]interface{}{
				"abc": 1,
				"def": 2,
			},
			stripKeys: []string{"xyz"},
			expected:  `{"abc":1,"def":2}`,
		},
		"AreChangesOnSimpleMatch": {
			sample: map[string]interface{}{
				"abc": 1,
				"def": 2,
			},
			stripKeys: []string{"abc"},
			expected:  `{"def":2}`,
		},
		"ChangesOnSubMapSimpleMatch": {
			sample: map[string]interface{}{
				"abc": map[string]interface{}{
					"def": 1,
				},
				"xyz": map[string]interface{}{
					"zyx": 2,
				},
				"aaa": nil,
			},
			stripKeys: []string{"abc"},
			expected:  `{"aaa":null,"xyz":{"zyx":2}}`,
		},
		"ChangesOnSubMapCascadeMatch": {
			sample: map[string]interface{}{
				"abc": map[string]interface{}{
					"def": 1,
				},
				"xyz": map[string]interface{}{
					"zyx": 2,
				},
				"aaa": nil,
			},
			stripKeys: []string{"abc>def"},
			expected:  `{"aaa":null,"abc":{},"xyz":{"zyx":2}}`,
		},
		"ChangesOnSubMapUnknownKey": {
			sample: map[string]interface{}{
				"abc": map[string]interface{}{
				},
				"xyz": map[string]interface{}{
					"zyx": 2,
				},
				"aaa": nil,
			},
			stripKeys: []string{"abc>def"},
			expected:  `{"aaa":null,"abc":{},"xyz":{"zyx":2}}`,
		},
		"ChangesOnSubArraySimpleMatch": {
			sample: map[string]interface{}{
				"abc": []interface{}{
					0: map[string]interface{}{
						"def": map[string]interface{}{
							"x": 1,
						},
					},
				},
				"aaa": nil,
			},
			stripKeys: []string{"abc"},
			expected:  `{"aaa":null}`,
		},
		"ChangesOnSubArrayCascadeMatch": {
			sample: map[string]interface{}{
				"abc": []interface{}{
					0: map[string]interface{}{
						"def": map[string]interface{}{
							"x": 1,
						},
						"def2": map[string]interface{}{
							"x2": 2,
						},
					},
					1: map[string]interface{}{
						"def": map[string]interface{}{
							"x": 1,
						},
						"def2": map[string]interface{}{
							"x2": 2,
						},
					},
				},
				"aaa": nil,
			},
			stripKeys: []string{"abc>def"},
			expected:  `{"aaa":null,"abc":[{"def2":{"x2":2}},{"def2":{"x2":2}}]}`,
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			StripKeys(&testCase.sample, testCase.stripKeys)

			got, _ := json.Marshal(testCase.sample)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

func TestStartKeys(t *testing.T) {
	testCases := map[string]struct {
		dataSet      map[string]interface{}
		startKeys    []string
		expected     string
		inheritAttrs bool
	}{
		"SimpleKeyWithNestedMap": {
			dataSet: map[string]interface{}{
				"abc": 1,
				"def": map[string]interface{}{
					"def2": map[string]interface{}{
						"def3": 3,
					},
					"xyz": "test",
				},
			},
			startKeys: []string{"def", "def2"},
			expected:  `{"def3":3}`,
		},
		"SimpleKeyWithNestedMap_InheritAttrs": {
			dataSet: map[string]interface{}{
				"abc": 1,
				"def": map[string]interface{}{
					"def2": map[string]interface{}{
						"def3": 3,
					},
					"xyz": "test",
				},
			},
			startKeys:    []string{"def", "def2"},
			expected:     `{"def3":3,"parent.0.abc":"1","parent.1.xyz":"test"}`,
			inheritAttrs: true,
		},
		"SimpleKeyWithArray_InheritAttrs": {
			dataSet: map[string]interface{}{
				"abc": 1,
				"def": []interface{}{
					map[string]interface{}{
						"def2": map[string]interface{}{
							"def3": 3,
						},
						"xyz": "test",
					},
					map[string]interface{}{
						"def2": map[string]interface{}{
							"def3": 3,
						},
						"xyz": "test2",
					},
				},
			},
			startKeys:    []string{"def"},
			expected:     `{"def":[{"def2":{"def3":3},"parent.0.abc":"1","xyz":"test"},{"def2":{"def3":3},"parent.0.abc":"1","xyz":"test2"}]}`,
			inheritAttrs: true,
		},
		"NestedKeyWithMap_InheritAttrs": {
			dataSet: map[string]interface{}{
				"abc": 1,
				"def": []interface{}{
					map[string]interface{}{
						"def2": map[string]interface{}{
							"def3": map[string]interface{}{
								"def4": 4,
							},
						},
						"xyz": "test",
					},
				},
			},
			startKeys: []string{"def>def2"},
			// TODO: shouldn't def2 be considered as another key and store also xyz as parent attribute?
			expected:     `{"def2":[{"def3":{"def4":4},"parent.0.abc":"1"}]}`,
			inheritAttrs: true,
		},
		"NestedKeyWithMapAndArray_InheritAttrs": {
			dataSet: map[string]interface{}{
				"abc": 1,
				"def": []interface{}{
					map[string]interface{}{
						"def2": []interface{}{
							"def3",
							"def4",
						},
						"xyz": "test",
					},
				},
			},
			startKeys: []string{"def>def2"},
			// TODO: shouldn't here contain also def4?
			expected:     `{"def2":["def3"]}`,
			inheritAttrs: true,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			FindStartKey(&testCase.dataSet, testCase.startKeys, testCase.inheritAttrs)

			got, _ := json.Marshal(testCase.dataSet)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

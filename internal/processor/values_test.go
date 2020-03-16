package processor

import (
	"encoding/json"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/assert"
)

func TestRunSubParse(t *testing.T) {
	getConfig := func() []load.Parse {
		return []load.Parse{
			{
				Type: "prefix",
				Key:  "db",
				SplitBy: []string{
					",",
					"=",
				},
			},
		}
	}

	testCases := map[string]struct {
		sample   map[string]interface{}
		value    interface{}
		parseCfg []load.Parse
		expected string
	}{
		"ValidParseConfig": {
			sample:   map[string]interface{}{},
			parseCfg: getConfig(),
			value:    `db0:keys=10237963,expires=224098,avg_ttl=0`,
			expected: `{"db.avg_ttl":"0","db.db0:keys":"10237963","db.expires":"224098"}`,
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunSubParse(testCase.parseCfg, &testCase.sample, "db", testCase.value)

			got, _ := json.Marshal(testCase.sample)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

func TestRunValueParser(t *testing.T) {
	getConfig := func(valueParser map[string]string) load.API {
		return load.API{
			ValueParser: valueParser,
		}
	}

	testCases := map[string]struct {
		parseCfg load.API
		value    interface{}
		key      string
		expected string
	}{
		"ParseString": {
			parseCfg: getConfig(
				map[string]string{
					"time[0-9]": "[0-9]+:[0-9]+:[0-9]+",
				}),
			key:      `time8`,
			value:    `12:31:32PM`,
			expected: `"12:31:32"`,
		},
		"ParseFloat": {
			parseCfg: getConfig(
				map[string]string{
					"int": "[0-9]+",
				}),
			key:      `int`,
			value:    13.8,
			expected: `"13"`,
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunValueParser(&testCase.value, testCase.parseCfg, &testCase.key)

			got, _ := json.Marshal(testCase.value)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

func TestRunValueTransformer(t *testing.T) {

	getConfig := func(valueTransformer map[string]string) load.API {
		return load.API{
			ValueTransformer: valueTransformer,
		}
	}

	testCases := map[string]struct {
		transformerCfg load.API
		value          interface{}
		key            string
		expected       string
	}{

		"TransformString": {
			transformerCfg: getConfig(
				map[string]string{
					"test.": "hello-${value}",
				}),
			key:      `test0`,
			value:    `world`,
			expected: `"hello-world"`,
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunValueTransformer(&testCase.value, testCase.transformerCfg, &testCase.key)

			got, _ := json.Marshal(testCase.value)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

func TestRunMathCalculations(t *testing.T) {
	testCases := map[string]struct {
		formulas map[string]string
		sample   map[string]interface{}
		key      string
		expected string
	}{
		"SimpleFormula": {
			formulas: map[string]string{
				"net.connectionsDroppedPerSecond": `${net.connectionsAcceptedPerSecond} - ${net.handledPerSecond}`,
			},
			sample: map[string]interface{}{
				"net.connectionsAcceptedPerSecond": 4,
				"net.handledPerSecond":             3,
			},
			expected: `{"net.connectionsAcceptedPerSecond":4,"net.connectionsDroppedPerSecond":1,"net.handledPerSecond":3}`,
		},
	}
	mathDefault := ""
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunMathCalculations(&testCase.formulas, &mathDefault, &testCase.sample)

			got, _ := json.Marshal(testCase.sample)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

func TestRunMathCalculationsWithDefault(t *testing.T) {
	testCases := map[string]struct {
		formulas map[string]string
		sample   map[string]interface{}
		key      string
		expected string
	}{
		"SimpleFormula": {
			formulas: map[string]string{
				"net.connectionsDroppedPerSecond": `${NonExistentAttribute} - ${net.handledPerSecond}`,
			},
			sample: map[string]interface{}{
				"net.connectionsAcceptedPerSecond": 4,
				"net.handledPerSecond":             3,
			},
			expected: `{"net.connectionsAcceptedPerSecond":4,"net.connectionsDroppedPerSecond":97,"net.handledPerSecond":3}`,
		},
	}
	mathDefault := "100"
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunMathCalculations(&testCase.formulas, &mathDefault, &testCase.sample)

			got, _ := json.Marshal(testCase.sample)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

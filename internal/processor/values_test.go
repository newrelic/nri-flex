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
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunMathCalculations(&testCase.formulas, &testCase.sample)

			got, _ := json.Marshal(testCase.sample)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

func TestRunValueMapper(t *testing.T) {
	// 	name: testValueMapper
	// 	apis:
	//   - name: getSomeData
	//     commands:
	//       - run: echo load:1.04, 2.20, 3.01
	//         split_by: ":"
	//     value_mapper:
	//       load=>load1:
	//         - (.+), (.+), (.+)=>$1
	//       load=>load2:
	//         - (.+), (.+), (.+)=>$2
	//       load=>load3:
	//         - (.+), (.+), (.+)=>$3

	getConfig := func(valueMapper map[string][]string) load.API {
		return load.API{
			ValueMapper: valueMapper,
		}
	}

	testCases := map[string]struct {
		valueMapperCfg load.API
		value          interface{}
		key            string
		expected       string
		sample         map[string]interface{}
	}{

		"TransformString": {
			valueMapperCfg: getConfig(
				map[string][]string{
					"load=>load1": []string{"(.+), (.+), (.+)=>$1"},
					"load=>load2": []string{"(.+), (.+), (.+)=>$2"},
					"load=>load3": []string{"(.+), (.+), (.+)=>$3"},
				}),
			key:    `load`,
			value:  `1.04, 2.20, 3.01`,
			sample: map[string]interface{}{},
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			RunValueMapper(testCase.valueMapperCfg.ValueMapper, &testCase.sample, testCase.key, &testCase.value)
			assert.Equal(t, "1.04", testCase.sample["load1"])
			assert.Equal(t, "2.20", testCase.sample["load2"])
			assert.Equal(t, "3.01", testCase.sample["load3"])
		})
	}
}

func TestRunTimestampConversion(t *testing.T) {
	getConfig := func(TimestampConversion map[string]string) load.API {
		return load.API{
			TimestampConversion: TimestampConversion,
		}
	}

	testCases := map[string]struct {
		parseCfg load.API
		value    interface{}
		key      string
		expected string
	}{
		"DATE2TIMESTAMP_Predefined_Date_Format": {
			parseCfg: getConfig(
				map[string]string{
					"started_at": "DATE2TIMESTAMP::RFC3339",
				}),
			key:      `started_at`,
			value:    `2020-07-20T14:34:05Z`,
			expected: `"1595255645"`,
		},
		"DATE2TIMESTAMP_Custom_Date_Format": {
			parseCfg: getConfig(
				map[string]string{
					"started_at": "DATE2TIMESTAMP::2006-01-02",
				}),
			key:      `started_at`,
			value:    `2020-07-20`,
			expected: `"1595203200"`,
		},

		"TIMESTAMP2DATE_Custom_Date_Format": {
			parseCfg: getConfig(
				map[string]string{
					"endtime": "TIMESTAMP2DATE::2006-01-02",
				}),
			key:      `endtime`,
			value:    1595598897,
			expected: `"2020-07-24"`,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {

			RunTimestampConversion(&testCase.value, testCase.parseCfg, &testCase.key)

			got, _ := json.Marshal(testCase.value)
			assert.Equal(t, testCase.expected, string(got))
		})
	}
}

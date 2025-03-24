/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */
package config

import (
	"os"
	"testing"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/assert"
)

func TestSubTimestamps(t *testing.T) {

	fileContent := `"${timestamp:ms}"
"${timestamp:ns}"
"${timestamp:s}"
"${timestamp:date}"
"${timestamp:datetime}"
"${timestamp:datetimetz}"
"${timestamp:dateutc}"
"${timestamp:datetimeutc}"
"${timestamp:datetimeutctz}"
"${timestamp:year}"
"${timestamp:month}"
"${timestamp:day}"
"${timestamp:hour}"
"${timestamp:minute}"
"${timestamp:second}"
"${timestamp:utcyear}"
"${timestamp:utcmonth}"
"${timestamp:utcday}"
"${timestamp:utchour}"
"${timestamp:utcminute}"
"${timestamp:utcsecond}"
"${timestamp:ms+10}"
"${timestamp:ns-10s}"
"${timestamp:ns-[Digits&NonDigits]}"`

	expected := `"138157323000"
"138157323000000004"
"138157323"
"1974-05-19"
"1974-05-19T01:02:03"
"1974-05-19T01:02:03Z"
"1974-05-19"
"1974-05-19T01:02:03"
"1974-05-19T01:02:03Z"
"1974"
"5"
"19"
"1"
"2"
"3"
"1974"
"5"
"19"
"1"
"2"
"3"
"138157323010"
"138157313000000004"
"138157323000"`

	fixedDate := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

	SubTimestamps(&fileContent, fixedDate)

	assert.Equal(t, expected, fileContent)
}

func Test_parseLookupFileAndUpdateConfig(t *testing.T) {

	testConfig := `
custom_attributes:
  replace_id_float: ${lf:id_float}
  replace_id_int: ${lf:id_int}
  replace_name: ${lf:name}
`

	item := map[string]interface{}{
		"id_float": 2456853.0,
		"id_int":   2456854,
		"name":     "AMP_eov_ntet-np",
	}

	actual, err := fillTemplateConfigWithValues(item, testConfig)

	assert.NoError(t, err)
	assert.NotNil(t, actual)

	expected := load.Config{
		CustomAttributes: map[string]string{
			"replace_id_float": "2456853.000000",
			"replace_id_int":   "2456854",
			"replace_name":     "AMP_eov_ntet-np",
		},
	}

	assert.ObjectsAreEqual(expected.CustomAttributes, (*actual).CustomAttributes)
}

func Test_toString(t *testing.T) {
	valueInt := 2456853
	assert.Equal(t, "2456853", toString(valueInt))
	valueFloat32 := float32(2456853)
	assert.Equal(t, "2456853.000000", toString(valueFloat32))
	valueFloat64 := float64(2456853)
	assert.Equal(t, "2456853.000000", toString(valueFloat64))
	valueString := "2456853"
	assert.Equal(t, "2456853", toString(valueString))
	valueMap := map[string]interface{}{"foo": "baz"}
	assert.Equal(t, "map[foo:baz]", toString(valueMap))
}

func TestSubEnvVariablescheck(t *testing.T) {
	// Set up environment variables
	os.Setenv("TEST_ENV_VAR", "test_value")
	os.Setenv("ANOTHER_ENV_VAR", "another_value")
	os.Setenv("FARGATE", "true")
	os.Setenv("FARGATE_TASK", "something")
	defer os.Unsetenv("TEST_ENV_VAR")
	defer os.Unsetenv("ANOTHER_ENV_VAR")
	defer os.Unsetenv("FARGATE")
	defer os.Unsetenv("FARGATE_TASK")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single substitution",
			input:    "Value is $$TEST_ENV_VAR",
			expected: "Value is test_value",
		},
		{
			name:     "Multiple substitutions",
			input:    "Values are $$TEST_ENV_VAR and $$ANOTHER_ENV_VAR",
			expected: "Values are test_value and another_value",
		},
		{
			name:     "No substitution",
			input:    "No env vars here",
			expected: "No env vars here",
		},
		{
			name:     "Partial substitution",
			input:    "Value is $$TEST_ENV_VAR and $$MISSING_ENV_VAR",
			expected: "Value is test_value and $$MISSING_ENV_VAR",
		},
		{
			name:     "Substitution with substr of env var",
			input:    "Value is $$FARGATE_TASK",
			expected: "Value is something",
		},
		{
			name:     "Substitution for only env var",
			input:    "FARGATE value is: $$FARGATE",
			expected: "FARGATE value is: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input
			SubEnvVariables(&input)
			assert.Equal(t, tt.expected, input)
		})
	}
}

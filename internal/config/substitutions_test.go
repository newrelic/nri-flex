/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */
package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

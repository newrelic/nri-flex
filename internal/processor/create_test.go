/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */
package processor

import (
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/assert"
)

func TestRunSampleFilter(t *testing.T) {
	getAPI := func() *load.API {
		return &load.API{
			SampleFilter: []map[string]string{
				{"customerId": "xyz"},
				{"customerId": "abc"},
				{"secretKey": "alpha"}
			},
		}
	}
	api := getAPI()

	createSample := true
	currentSample := map[string]interface{}{
		"customerId": "abc",
		"quantities": 10,
	}
	expectedResult := false
	RunSampleFilter(currentSample, api.SampleFilter, &createSample)
	assert.Equal(t, expectedResult, createSample)

	createSample = true
	currentSample = map[string]interface{}{
		"customerId": "xyz",
		"quantities": 20,
	}
	expectedResult = false
	RunSampleFilter(currentSample, api.SampleFilter, &createSample)
	assert.Equal(t, expectedResult, createSample)

	createSample = true
	currentSample = map[string]interface{}{
		"customerId": "aaaa",
		"quantities": 20,
	}
	expectedResult = true
	RunSampleFilter(currentSample, api.SampleFilter, &createSample)
	assert.Equal(t, expectedResult, createSample)

	createSample = true
	currentSample = map[string]interface{}{
		"customerId": "abc",
		"secretKey": "oof",
	}
	expectedResult = false
	RunSampleFilterMatchAll(currentSample, api.SampleFilter, &createSample)
	assert.Equal(t, expectedResult, createSample)

	createSample = true
	currentSample = map[string]interface{}{
		"customerId": "abc",
		"secretKey": "alpha",
	}
	expectedResult = true
	RunSampleFilterMatchAll(currentSample, api.SampleFilter, &createSample)
	assert.Equal(t, expectedResult, createSample)

}

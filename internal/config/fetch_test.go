package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDimensionalLookup(t *testing.T) {
	lookupStore := map[string][]string{}
	lookupStore["animal"] = []string{
		"dog",
		"cat",
		"cow",
	}
	lookupStore["gender"] = []string{
		"m",
		"f",
	}
	lookupStore["numbers"] = []string{
		"10",
		"11",
	}

	lookupDimensions := []string{
		"animal",
		"gender",
		"numbers",
	}

	var sliceIndexes []int
	var sliceKeys []string
	var sliceLookups [][]string

	// init lookups
	for _, dimKey := range lookupDimensions {
		if value, found := lookupStore[dimKey]; found {
			sliceIndexes = append(sliceIndexes, 0)
			sliceKeys = append(sliceKeys, dimKey)
			sliceLookups = append(sliceLookups, value)
		}
	}

	loopNo := -1
	var combinations [][]string
	loopLookups(loopNo, sliceIndexes, sliceLookups, &combinations)

	expected := make([][]string, 12)
	expected[0] = []string{"dog", "m", "10"}
	expected[1] = []string{"dog", "m", "11"}
	expected[2] = []string{"dog", "f", "10"}
	expected[3] = []string{"dog", "f", "11"}
	expected[4] = []string{"cat", "m", "10"}
	expected[5] = []string{"cat", "m", "11"}
	expected[6] = []string{"cat", "f", "10"}
	expected[7] = []string{"cat", "f", "11"}
	expected[8] = []string{"cow", "m", "10"}
	expected[9] = []string{"cow", "m", "11"}
	expected[10] = []string{"cow", "f", "10"}
	expected[11] = []string{"cow", "f", "11"}

	if len(expected) != len(combinations) {
		t.Errorf("want: %v combinations, got: %v", len(expected), len(combinations))
	}

	for i, combo := range combinations {
		for z, comboKey := range combo {
			if expected[i][z] != comboKey {
				t.Errorf("want: %v, got: %v", expected[i][z], comboKey)
			}
		}
	}
}

func Test_findLookups(t *testing.T) {
	testCases := []struct {
		name            string
		content         string
		expectedLookups [][]string
	}{
		{
			name:            "no lookups",
			content:         "",
			expectedLookups: nil,
		},
		{
			name:    "one lookup",
			content: "http://some-other-service.com/users/${lookup.postSample:userId}",
			expectedLookups: [][]string{
				{
					"${lookup.postSample:userId}", "postSample", "userId",
				},
			},
		},
		{
			name:    "multiple lookup",
			content: "http://some-other-service.com/users/${lookup.postSample:userId}_${lookup.postSample:name}__${lookup.anotherSample:anotherValue}",
			expectedLookups: [][]string{
				{
					"${lookup.postSample:userId}", "postSample", "userId",
				},
				{
					"${lookup.postSample:name}", "postSample", "name",
				},
				{
					"${lookup.anotherSample:anotherValue}", "anotherSample", "anotherValue",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedLookups, findLookups(tc.content))
		})
	}
}

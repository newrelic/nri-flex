package config

import (
	"testing"
)

func TestDimensionalLookup(t *testing.T) {
	lookupStore := map[string]map[string]struct{}{}
	lookupStore["animal"] = map[string]struct{}{
		"dog": struct{}{},
		"cat": struct{}{},
		"cow": struct{}{}}
	lookupStore["gender"] = map[string]struct{}{
		"m": struct{}{},
		"f": struct{}{}}
	lookupStore["numbers"] = map[string]struct{}{
		"10": struct{}{},
		"11": struct{}{}}

	lookupDimensions := []string{}
	for lookupVar := range lookupStore {
		lookupDimensions = append(lookupDimensions, lookupVar)
	}

	sliceIndexes := []int{}
	sliceKeys := []string{}
	sliceLookups := [][]string{}

	// init lookups
	for key, values := range lookupStore {
		// only create lookups for the found dimensions
		for _, dimKey := range lookupDimensions {
			if key == dimKey {
				sliceIndexes = append(sliceIndexes, 0)
				sliceKeys = append(sliceKeys, key)
				valueArray := []string{}
				for a := range values {
					valueArray = append(valueArray, a)
				}
				sliceLookups = append(sliceLookups, valueArray)
				break
			}
		}
	}

	loopNo := -1
	combinations := [][]string{}
	loopLookups(loopNo, sliceIndexes, sliceLookups, &combinations)

	expected := map[int][]string{}
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

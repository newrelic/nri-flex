/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"fmt"
	"testing"
)

func TestStoreLookups(t *testing.T) {

	storeLookups := map[string]string{
		"blah": "abc",
	}
	key := "abc"
	lookupStore := map[string]map[string]struct{}{}
	var v interface{}
	v = "myStoredValue"

	StoreLookups(storeLookups, &lookupStore, key, v)
	valueArray := []string{}
	for a := range lookupStore["blah"] {
		valueArray = append(valueArray, a)
	}

	if fmt.Sprintf("%v", valueArray[0]) != v {
		t.Errorf("want: %v got: %v", v, valueArray[0])
	}
}

func TestVariableLookups(t *testing.T) {
	storeLookups := map[string]string{
		"blah": "abc",
	}
	key := "abc"
	variableStore := map[string]string{}
	var v interface{}
	v = "myStoredValue"

	VariableLookups(storeLookups, &variableStore, key, v)

	if fmt.Sprintf("%v", variableStore["blah"]) != v {
		t.Errorf("want: %v got: %v", v, variableStore["blah"])
	}
}

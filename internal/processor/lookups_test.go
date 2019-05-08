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
	lookupStore := map[string][]string{}
	var v interface{}
	v = "myStoredValue"

	StoreLookups(storeLookups, &key, &lookupStore, &v)

	if fmt.Sprintf("%v", lookupStore["blah"][0]) != v {
		t.Errorf("want: %v got: %v", v, lookupStore["blah"][0])
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

	VariableLookups(storeLookups, &key, &variableStore, &v)

	if fmt.Sprintf("%v", variableStore["blah"]) != v {
		t.Errorf("want: %v got: %v", v, variableStore["blah"])
	}
}

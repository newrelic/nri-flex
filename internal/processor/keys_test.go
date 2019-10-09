/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"encoding/json"
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
	keepKeys := []string{"abc", "xyz"}
	currentSample := map[string]interface{}{
		"abc":   1,
		"xyz":   2,
		"dfadf": "dafsdfa",
		"def":   "initely busted",
	}

	key := "def"

	RunKeepKeys(keepKeys, &key, &currentSample)
	output, _ := json.Marshal(currentSample)
	expected := `{"abc":1,"dfadf":"dafsdfa","xyz":2}`

	if string(output) != expected {
		t.Errorf("want: %v got: %v", expected, string(output))
	}
}

func TestKeyRemover(t *testing.T) {
	removeKeys := []string{"abc", "xyz"}
	currentSample := map[string]interface{}{
		"abc":   1,
		"xyz":   2,
		"dfadf": "dafsdfa",
	}

	RunKeyRemover(&currentSample, removeKeys)
	output, _ := json.Marshal(currentSample)
	expected := `{"dfadf":"dafsdfa"}`

	if string(output) != expected {
		t.Errorf("want: %v got: %v", expected, string(output))
	}
}

func TestKeyRenamer(t *testing.T) {
	renameKeys := map[string]string{"abc$": "xyz"}
	key := "aaaaabc"
	originalKey := "abc"

	RunKeyRenamer(renameKeys, &key, originalKey)
	if key != "aaaaxyz" {
		t.Errorf("want: xyz got: %v", key)
	}
}

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
	}

	key := "def"
	k := "def"

	RunKeepKeys(keepKeys, &key, &currentSample, &k)
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
	key := "abc"
	progress := false

	RunKeyRemover(removeKeys, &key, &progress, &currentSample)
	output, _ := json.Marshal(currentSample)
	expected := `{"dfadf":"dafsdfa","xyz":2}`

	if string(output) != expected {
		t.Errorf("want: %v got: %v", expected, string(output))
	}
}

func TestKeyRenamer(t *testing.T) {
	renameKeys := map[string]string{"abc": "xyz"}
	key := "abc"

	RunKeyRenamer(renameKeys, &key)
	if key != "xyz" {
		t.Errorf("want: xyz got: %v", key)
	}
}

// // RunKeyRenamer find key with regex, and replace the value
// func RunKeyRenamer(renameKeys map[string]string, key *string) {
// 	for renameKey, renameVal := range renameKeys {
// 		if formatter.KvFinder("regex", *key, renameKey) {
// 			*key = strings.Replace(*key, renameKey, renameVal, -1)
// 			break
// 		}
// 	}
// }

package processor

import (
	"strings"

	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
)

// RunKeyConversion handles to lower and snake to camel case for keys
func RunKeyConversion(key *string, api load.API, v interface{}, SkipProcessing *[]string) {
	if api.ToLower {
		*key = strings.ToLower(*key)
	}
	if api.ConvertSpace != "" {
		*key = strings.Replace(*key, " ", api.ConvertSpace, -1)
	}
	if api.SnakeToCamel {
		formatter.SnakeCaseToCamelCase(key)
	}
}

// RunKeepKeys Removes all other keys/attributes and keep only those defined in keep_keys
func RunKeepKeys(keepKeys []string, key *string, currentSample *map[string]interface{}, k *string) {
	if len(keepKeys) > 0 {
		foundKey := false
		for _, keepKey := range keepKeys {
			if formatter.KvFinder("regex", *key, keepKey) {
				foundKey = true
				break
			}
		}
		if !foundKey {
			delete(*currentSample, *k)
		}
	}
}

// RunKeyRemover Remove unwanted keys
func RunKeyRemover(removeKeys []string, key *string, progress *bool, currentSample *map[string]interface{}) {
	for _, removeKey := range removeKeys {
		if formatter.KvFinder("regex", *key, removeKey) {
			*progress = false
			delete(*currentSample, *key)
			break
		}
	}
}

// RunKeyRenamer find key with regex, and replace the value
func RunKeyRenamer(renameKeys map[string]string, key *string) {
	for renameKey, renameVal := range renameKeys {
		if formatter.KvFinder("regex", *key, renameKey) {
			*key = strings.Replace(*key, renameKey, renameVal, -1)
			break
		}
	}
}

// StripKeys strip defined keys out
func StripKeys(ds *map[string]interface{}, stripKeys []string) {
	for _, stripKey := range stripKeys {
		delete(*ds, stripKey)
		if strings.Contains(stripKey, ">") {
			stripSplit := strings.Split(stripKey, ">")
			if len(stripSplit) == 2 {
				if (*ds)[stripSplit[0]] != nil {
					switch (*ds)[stripSplit[0]].(type) {
					case map[string]interface{}:
						delete((*ds)[stripSplit[0]].(map[string]interface{}), stripSplit[1])
					case []interface{}:
						for i := range (*ds)[stripSplit[0]].([]interface{}) {
							switch (*ds)[stripSplit[0]].([]interface{})[i].(type) {
							case map[string]interface{}:
								delete((*ds)[stripSplit[0]].([]interface{})[i].(map[string]interface{}), stripSplit[1])
							}
						}
					}
				}
			}
		}
	}
}

// FindStartKey start at a different section of a payload
func FindStartKey(mainDataset *map[string]interface{}, startKeys []string) {
	for _, startKey := range startKeys {
		// fmt.Println(i, startKey)
		// fmt.Println()

		startSplit := strings.Split(startKey, ">")
		if len(startSplit) == 2 {
			// fmt.Println(startSplit)
			if (*mainDataset)[startSplit[0]] != nil {
				switch mainDs := (*mainDataset)[startSplit[0]].(type) {
				case []interface{}:
					nestedSlices := []interface{}{}
					for _, nested := range mainDs {
						switch sample := nested.(type) {
						case map[string]interface{}:
							if sample[startSplit[1]] != nil {
								switch nestedSample := sample[startSplit[1]].(type) {
								case map[string]interface{}:
									nestedSlices = append(nestedSlices, nestedSample)
								case []interface{}:
									nestedSlices = append(nestedSlices, nestedSample[0])
								}
							}
						}
					}
					*mainDataset = map[string]interface{}{startSplit[1]: nestedSlices}
				}
			}
		} else if len(startSplit) == 1 {
			if (*mainDataset)[startKey] != nil {
				switch mainDs := (*mainDataset)[startKey].(type) {
				case map[string]interface{}:
					*mainDataset = mainDs
				case []interface{}:
					*mainDataset = map[string]interface{}{startKey: mainDs}
					// x, _ := json.Marshal(*mainDataset)
					// fmt.Println(string(x))
				}
			} else {
				// fmt.Println("didn't find it:", startKey)
			}
		}
	}
}

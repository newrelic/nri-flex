package outputs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//function to store samples as a JSON object at specified path
func StoreJson(samples []interface{}, path string) error {
	if samples == nil {
		return nil
	}
	bytes, _ := json.Marshal(samples)
	err := ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return fmt.Errorf("file output: failed to write ")
	}
	return nil
}
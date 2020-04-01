package outputs

import (
	"encoding/json"
	"io/ioutil"
	"errors"
)

//function to store samples as a JSON object at specified path
func StoreJson(samples []interface{}, path string) {
	bytes, _ := json.Marshal(samples)
	err := ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		panic(errors.New("failed to write file: %v", err))
	}
}
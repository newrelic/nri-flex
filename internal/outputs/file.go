package outputs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

//function to store samples as a JSON object at specified path
func StoreJSON(samples []interface{}, path string) {
	bytes, _ := json.Marshal(samples)
	err := ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		panic(errors.New("failed to write file: %v"))
	}
}

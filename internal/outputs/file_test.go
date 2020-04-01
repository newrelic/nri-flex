package outputs

import (
	"testing"	
	"encoding/json"
	"io/ioutil"
	"fmt"
	"string"
	"github.com/stretchr/testify/assert"
	"github.com/newrelic/nri-flex/internal/load"
)

const testData = `[{"abd":"def"},{"123":456}]`

func testDataProvider() []interface{} {
	jsonData := []byte(testData)
	var out []interface{} 
	if err := json.Unmarshal(jsonData, &out); err != nil {
		panic(err)
	}
	return out
}

func readJsonFile(file string, output *string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("file input: failed to read file: %v", err)
	}
	*output = string(b)
}

// testing function for output writes
func testOutputWrite(t *testing.T) {
	fname := "writeTest.json"
	var out []interface{} = testDataProvider()
	storeJson(out, fname)
	var readVal string
	readJsonFile(fname, &readVal)
	assert.Equals(testData, readVal, "Values should be the same");
}
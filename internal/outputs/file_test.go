package outputs

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
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

func readJSONFile(file string, output *string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic("reading JSON file failed")
	}
	*output = string(b)
}

// testing function for output writes
func testOutputWrite(t *testing.T) {
	fname := "writeTest.json"
	var out []interface{} = testDataProvider()
	storeJSON(out, fname)
	var readVal string
	readJSONFile(fname, &readVal)
	assert.Equals(testData, readVal, "Values should be the same")
}

package outputs

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"os"
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

func removeTestFile(path string) {
	err := os.Remove(path)
	if err != null {
		panic("unable to remove writeTest.json")
	}
}

// testing function for output writes
func TestOutputWrite(t *testing.T) {
	fname := "writeTest.json"
	var out []interface{} = testDataProvider()
	StoreJSON(out, fname)
	var readVal string
	readJSONFile(fname, &readVal)
	assert.Equal(t, testData, readVal)
	removeTestFile(fname)
}

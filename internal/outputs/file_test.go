package outputs

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const testData = `[{"abd":"def"},{"123":456}]`

func testDataProvider(t *testing.T) []interface{} {
	jsonData := []byte(testData)
	var out []interface{}
	err := json.Unmarshal(jsonData, &out)
	require.NoError(t, err)
	return out
}

func readJSONFile(t *testing.T, file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		require.NoError(t, err, "reading JSON file failed")
	}
	return string(b)
}

// testing function for output writes
func TestOutputWrite(t *testing.T) {
	dir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(dir) }()

	fname := path.Join(dir, "writeTest.json")
	var out = testDataProvider(t)
	StoreJSON(out, fname)
	readVal := readJSONFile(t, fname)
	assert.Equal(t, testData, readVal)
}

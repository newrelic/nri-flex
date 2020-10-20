package inputs

import (
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/assert"
)

func TestParseToJSON(t *testing.T) {
	getConfig := func(parseHTML bool) load.API {
		return load.API{
			ParseHTML: parseHTML,
		}
	}

	testCases := map[string]struct {
		parseCfg load.API
		value    string
		key      string
		expected string
	}{
		"SingleTable": {
			parseCfg: getConfig(
				true),
			value: `<html><body>
			<table source="myTestPage1">
				<tr><th>Heading 1</th><th>Heading 11</th><th>Heading 12</th><th>Heading 13</th><th>Heading 14</th></tr>
				<tr><td>Data 11</td><td>Data 12</td></tr>
				<tr><td>Data 21</td><td>Data 22</td></tr>
				<tr><td>Data 31</td><td>Data 32</td></tr>
				<tr><td>Data 41</td><td>Data 42</td></tr>
			</table>
			</html>`,
			expected: `[{"table":[{ "Heading 1": "Data 11", "Heading 11": "Data 12", "Heading 12": "", "Heading 13": "", "Heading 14": ""},{ "Heading 1": "Data 21", "Heading 11": "Data 22", "Heading 12": "", "Heading 13": "", "Heading 14": ""},{ "Heading 1": "Data 31", "Heading 11": "Data 32", "Heading 12": "", "Heading 13": "", "Heading 14": ""},{ "Heading 1": "Data 41", "Heading 11": "Data 42", "Heading 12": "", "Heading 13": "", "Heading 14": ""}], "source": "myTestPage1","Index":0 }]`,
		},
		"TwoTableWithAttribute": {
			parseCfg: getConfig(
				true),
			value: `<html><body>
			<table source="myTestTable1">
			  <tr><th>Heading 1</th><th>Heading 11</th><th>Heading 12</th><th>Heading 13</th><th>Heading 14</th></tr>
			  <tr><td>Data 11</td><td>Data 12</td></tr>
			  <tr><td>Data 21</td><td>Data 22</td></tr>
			  <tr><td>Data 31</td><td>Data 32</td></tr>
			  <tr><td>Data 41</td><td>Data 42</td></tr>
			</table>
			<p>Stuff in here</p>
			<table source="myTestTable2">
			  <tr><th>Heading 21</th><th>Heading 22</th></tr>
			  <tr><td>Data 211</td><td>Data 212</td></tr>
			  <tr><td>Data 221</td><td>Data 222</td></tr>
			  <tr><td>Data 231</td><td><span></span><span><a href="">Data 232</a></span></td></tr>
			  <tr><td>Data 241</td><td>Data 242</td></tr>
			</table>
			</body>
			</html>`,
			expected: `[{"table":[{ "Heading 1": "Data 11", "Heading 11": "Data 12", "Heading 12": "", "Heading 13": "", "Heading 14": ""},{ "Heading 1": "Data 21", "Heading 11": "Data 22", "Heading 12": "", "Heading 13": "", "Heading 14": ""},{ "Heading 1": "Data 31", "Heading 11": "Data 32", "Heading 12": "", "Heading 13": "", "Heading 14": ""},{ "Heading 1": "Data 41", "Heading 11": "Data 42", "Heading 12": "", "Heading 13": "", "Heading 14": ""}], "source": "myTestTable1","Index":0 },{"table":[{ "Heading 21": "Data 211", "Heading 22": "Data 212"},{ "Heading 21": "Data 221", "Heading 22": "Data 222"},{ "Heading 21": "Data 231", "Heading 22": "Data 232"},{ "Heading 21": "Data 241", "Heading 22": "Data 242"}], "source": "myTestTable2","Index":1 }]`,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {

			result, _ := ParseToJSON([]byte(testCase.value))
			assert.Equal(t, testCase.expected, string(result))
		})
	}
}

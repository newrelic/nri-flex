package processor

import (
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/assert"
)

func TestRunJq(t *testing.T) {

	input := []interface{}{map[string]interface{}{
		"foo": map[string]interface{}{
			"data": map[string]interface{}{
				"abc": 1,
				"def": 2,
			},
		},
	}}

	api := load.API{
		Jq: ".foo.data",
	}

	var dataSets = runJq(input, api)
	var expectedResult = []interface{}{map[string]interface{}{"abc": 1, "def": 2}}

	if len(dataSets) != len(expectedResult) {
		t.Errorf("Missing samples, got: %v, want: %v.", dataSets, expectedResult)
	}

	abc1 := expectedResult[0].(map[string]interface{})["abc"]
	abc2 := dataSets[0].(map[string]interface{})["abc"]
	if abc1 != abc2 {
		t.Errorf("Wrong data for abc, got: %v, want: %v.", abc2, abc1)
	}

	def1 := expectedResult[0].(map[string]interface{})["def"]
	def2 := dataSets[0].(map[string]interface{})["def"]
	if def1 != def2 {
		t.Errorf("Wrong data for def, got: %v, want: %v.", def2, def1)
	}

}

func TestRunJqMultipleResults(t *testing.T) {

	input := []interface{}{map[string]interface{}{
		"foo": map[string]interface{}{
			"data": map[string]interface{}{
				"abc": 1,
				"def": 2,
			},
			"more_data": map[string]interface{}{
				"ghi": 3,
				"jkl": 4,
			},
		},
	}}

	api := load.API{
		Jq: ".foo.data,.[0].foo.more_data",
	}

	var dataSets = runJq(input, api)
	var expectedResult = []interface{}{
		map[string]interface{}{"abc": 1, "def": 2},
		map[string]interface{}{"ghi": 3, "jkl": 4},
	}

	assert.Equal(t, expectedResult, dataSets)
}

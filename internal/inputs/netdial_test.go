package inputs

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestNetDial(t *testing.T) {
	load.Refresh()

	config := load.Config{
		APIs: []load.API{
			{
				Name: "failure",
				Commands: []load.Command{
					load.Command{
						Dial: "fake12311290.com:9989",
					},
				},
			},
		},
	}

	var jsonOut interface{}
	expectedOutput, _ := ioutil.ReadFile("../../test/payloadsExpected/portTestSingle.json")
	json.Unmarshal(expectedOutput, &jsonOut)
	expectedDatastore := jsonOut.([]interface{})

	dataStore := []interface{}{}
	dataSample := map[string]interface{}{}
	processType := ""
	NetDialWithTimeout(&dataStore, config.APIs[0].Commands[0], &dataSample, config.APIs[0], &processType)

	if len(expectedDatastore) != len(dataStore) {
		t.Errorf("Incorrect number of samples generated expected: %d, got: %d", len(expectedDatastore), len(dataStore))
		t.Errorf("%v", (dataStore))
	}

	for _, sample := range expectedDatastore {
		switch sample := sample.(type) {
		case map[string]interface{}:
			for _, rSample := range dataStore {
				switch recSample := rSample.(type) {
				case map[string]interface{}:
					for key := range sample {
						if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
							t.Errorf("dbSample %v want %v, got %v", key, sample[key], recSample[key])
						}
					}
				}
			}
		}
	}
}

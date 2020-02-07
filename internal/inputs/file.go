package inputs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

// ProcessFile read and process the file into data collection.
func ProcessFile(dataStore *[]interface{}, cfg *load.Config, apiNo int) error {
	file := cfg.APIs[apiNo].File

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("file input: failed to read file: %v", err)
	}

	fileContent := string(b)

	if strings.HasSuffix(file, ".csv") {
		return processCsv(dataStore, cfg.Name, file, &fileContent, cfg.APIs[apiNo].SetHeader)
	}

	return processJSON(dataStore, fileContent)
}

func processCsv(dataStore *[]interface{}, cfgName, file string, data *string, header []string) error {
	load.Logrus.WithFields(logrus.Fields{
		"name": cfgName,
		"file": file,
	}).Debug("file input: reading csv")

	r := csv.NewReader(strings.NewReader(*data))
	var keys []string
	if len(header) > 0 {
		keys = header
	}

	index := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// TODO: Should we return error here?
			load.Logrus.WithFields(logrus.Fields{
				"err":  err,
				"file": file,
			}).Error("file input: failed to read csv line")
		}

		// establish header / keys
		if index == 0 && len(keys) == 0 {
			keys = append(keys, record...)
		} else {
			if len(record) != len(keys) {
				return fmt.Errorf("file input: csv header and record length mismatch: %d headerValues vs %d recordValues",
					len(keys), len(record))
			}
			newSample := map[string]interface{}{}
			for i, key := range keys {
				newSample[key] = record[i]
			}
			*dataStore = append(*dataStore, newSample)
		}
		index++
	}
	return nil
}

func processJSON(dataStore *[]interface{}, data string) error {
	newBody := strings.Replace(data, " ", "", -1)

	var f interface{}
	err := json.Unmarshal([]byte(newBody), &f)
	if err != nil {
		return fmt.Errorf("file input: failed to unmarshal JSON: %v", err)
	}

	*dataStore = append(*dataStore, f)
	return nil
}

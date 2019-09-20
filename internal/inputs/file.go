package inputs

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

// RunFile runs file read data collection
func RunFile(dataStore *[]interface{}, cfg *load.Config, apiNo int) {
	file := cfg.APIs[apiNo].File
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"name": cfg.Name,
			"file": file,
		}).Error("fetch: failed to read")
	} else {
		if strings.HasSuffix(file, ".csv") {
			load.Logrus.WithFields(logrus.Fields{
				"name": cfg.Name,
				"file": file,
			}).Debug("fetch: reading csv")
			processCsv(dataStore, file, string(fileData), cfg.APIs[apiNo].SetHeader)
		} else {
			newBody := strings.Replace(string(fileData), " ", "", -1)
			var f interface{}
			err := json.Unmarshal([]byte(newBody), &f)
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"name": cfg.Name,
					"file": file,
				}).Error("fetch: failed to unmarshal")
			} else {
				*dataStore = append(*dataStore, f)
			}
		}
	}
}

func processCsv(dataStore *[]interface{}, file string, data string, header []string) {
	r := csv.NewReader(strings.NewReader(data))

	keys := []string{}
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
			load.Logrus.WithFields(logrus.Fields{
				"err":  err,
				"file": file,
			}).Error("commands: failed to read csv line")
		}

		// establish header / keys
		if index == 0 && len(keys) == 0 {
			keys = append(keys, record...)
		} else {
			if len(record) == len(keys) {
				newSample := map[string]interface{}{}
				for i, key := range keys {
					newSample[key] = record[i]
				}
				*dataStore = append(*dataStore, newSample)
			} else {
				load.Logrus.WithFields(logrus.Fields{
					"headerValues": len(keys),
					"recordValues": len(record),
					"file":         file,
				}).Error("commands: csv header and record length mismatch")
				break
			}
		}
		index++
	}
}

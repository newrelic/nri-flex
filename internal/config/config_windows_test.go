//+build windows

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

// testSamples as samples could be generated in different orders, so we test per sample
func testSamples(expectedSamples []metric.Set, t *testing.T) {
	entityMetrics, _ := json.Marshal(load.Entity.Metrics)
	expectedMetrics, _ := json.Marshal(expectedSamples)
	if len(load.Entity.Metrics) != len(expectedSamples) {
		t.Errorf("Missing samples, got: %v, want: %v.", string(entityMetrics), string(expectedMetrics))
	}

	for _, sample := range load.Entity.Metrics {
		success := false
		for _, expectedSample := range expectedSamples {
			keyMatches := 0
			for k := range expectedSample.Metrics {
				if fmt.Sprintf("%v", expectedSample.Metrics[k]) == fmt.Sprintf("%v", sample.Metrics[k]) {
					keyMatches++
				}
			}
			if keyMatches == len(expectedSample.Metrics) && keyMatches != 0 {
				success = true
				t.Logf("matched keys: %d", keyMatches)
				break
			}
		}
		if !success {
			t.Errorf("Failed to match sample, got: %v, want: %v.", string(entityMetrics), string(expectedMetrics))
		}
	}
}

func TestConfigDir(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestReadJsonCmdDir", "nri-flex")
	load.Args.ConfigDir = "../../test/configs/windows/"

	var ymls []load.Config
	var files []os.FileInfo

	path := filepath.FromSlash(load.Args.ConfigDir)
	var err error
	files, err = ioutil.ReadDir(path)

	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("failed to read config dir: " + load.Args.ConfigDir)
	}

	LoadFiles(&ymls, files, path) // load standard configs if available
	RunFiles(&ymls)

	jsonFile, _ := ioutil.ReadFile("../../test/payloadsExpected/configDir.json")
	expectedOutput := []metric.Set{}
	err = json.Unmarshal(jsonFile, &expectedOutput)
	if err != nil {
		t.Error(err)
	}

	testSamples(expectedOutput, t)
}

func TestConfigFile(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestReadJsonCmd", "nri-flex")
	load.Args.ConfigFile = "../../test/configs/windows/json-read-cmd-example.yml"

	// Read a single config file
	var files []os.FileInfo
	var ymls []load.Config
	file, _ := os.Stat(load.Args.ConfigFile)
	path := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)

	LoadFiles(&ymls, files, path) // load standard configs if available
	RunFiles(&ymls)

	jsonFile, _ := ioutil.ReadFile("../../test/payloadsExpected/configFile.json")
	expectedOutput := []metric.Set{}
	err := json.Unmarshal(jsonFile, &expectedOutput)
	if err != nil {
		t.Error(err)
	}
	testSamples(expectedOutput, t)
}

func TestV4ConfigFile(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestV4Cmd", "nri-flex")
	load.Args.ConfigFile = "../../test/configs/windows/v4-integrations-example.yml"

	// Read a single config file
	var files []os.FileInfo
	var ymls []load.Config
	file, _ := os.Stat(load.Args.ConfigFile)
	path := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)

	LoadFiles(&ymls, files, path) // load standard configs if available
	RunFiles(&ymls)

	jsonFile, _ := ioutil.ReadFile("../../test/payloadsExpected/configFileV4.json")
	expectedOutput := []metric.Set{}
	err := json.Unmarshal(jsonFile, &expectedOutput)
	if err != nil {
		t.Error(err)
	}
	testSamples(expectedOutput, t)
}

func TestSubEnvVariables(t *testing.T) {
	str := " hi there $$PWD bye"
	SubEnvVariables(&str)
	if strings.Count(str, "$$") != 0 {
		t.Errorf("failed to sub all variables %v", str)
	}
}

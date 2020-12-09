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
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
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

	configsPath := path.Join("..", "..", "test", "configs")
	if runtime.GOOS == "windows" {
		configsPath = path.Join("..", "..", "test", "config_windows")
	}
	load.Args.ConfigDir = configsPath

	var ymls []load.Config
	var files []os.FileInfo
	var err error

	files, err = ioutil.ReadDir(load.Args.ConfigDir)
	require.NoError(t, err, "failed to read config dir: %s", load.Args.ConfigDir)

	errs := LoadFiles(&ymls, files, load.Args.ConfigDir) // load standard configs if available
	for _, err = range errs {
		assert.NoError(t, err)
	}
	require.Empty(t, errs)

	errs = RunFiles(&ymls)
	for _, err = range errs {
		assert.NoError(t, err)
	}
	require.Empty(t, errs)

	jsonFile, _ := ioutil.ReadFile(path.Join("..", "..", "test", "payloadsExpected", "configDir.json"))
	var expectedOutput []metric.Set
	err = json.Unmarshal(jsonFile, &expectedOutput)
	require.NoError(t, err)

	testSamples(expectedOutput, t)
}

func TestConfigFile(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestReadJsonCmd", "nri-flex")

	configsPath := path.Join("..", "..", "test", "configs")
	if runtime.GOOS == "windows" {
		configsPath = path.Join("..", "..", "test", "config_windows")
	}
	load.Args.ConfigFile = path.Join(configsPath, "json-read-cmd-example.yml")

	// Read a single config file
	var files []os.FileInfo
	var ymls []load.Config
	file, _ := os.Stat(load.Args.ConfigFile)
	filePath := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)

	errs := LoadFiles(&ymls, files, filePath) // load standard configs if available
	for _, err := range errs {
		assert.NoError(t, err)
	}
	require.Empty(t, errs)

	RunFiles(&ymls)

	jsonFile, _ := ioutil.ReadFile(path.Join("..", "..", "test", "payloadsExpected", "configFile.json"))
	var expectedOutput []metric.Set
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

	configsPath := path.Join("..", "..", "test", "configs")
	if runtime.GOOS == "windows" {
		configsPath = path.Join("..", "..", "test", "config_windows")
	}
	load.Args.ConfigFile = path.Join(configsPath, "v4-integrations-example.yml")

	// Read a single config file
	var files []os.FileInfo
	var ymls []load.Config
	file, _ := os.Stat(load.Args.ConfigFile)
	filePath := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)

	LoadFiles(&ymls, files, filePath) // load standard configs if available
	RunFiles(&ymls)

	jsonFile, _ := ioutil.ReadFile(path.Join("..", "..", "test", "payloadsExpected", "configFileV4.json"))
	var expectedOutput []metric.Set
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

func TestApplyFlexMeta(t *testing.T) {
	os.Setenv("FLEX_META", "{\"abc\":\"def\",\"hello\":123}")
	config := load.Config{}
	applyFlexMeta(&config)

	if config.CustomAttributes["abc"] != "def" {
		t.Errorf("failed to apply flex meta variable abc expected %v got %v", "def", config.CustomAttributes["hello"])
	} else if config.CustomAttributes["hello"] != "123" {
		t.Errorf("failed to apply flex meta variable hello expected %v got %v", "123", config.CustomAttributes["hello"])
	}
}

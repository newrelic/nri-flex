// +build integration

// The tests in this file are supposed to be run in the CI using docker-compose
// You can run then from the IDE or manually but you'll need some setup first
package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/stretchr/testify/assert"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
)

func TestConfig_cmd_LinuxDirUsage(t *testing.T) {
	// given
	load.Refresh()

	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")
	load.Args.ConfigFile = "configs/linux-directory-size.yml"

	// Read a single config file
	var files []os.FileInfo
	var configs []load.Config
	file, err := os.Stat(load.Args.ConfigFile)
	if err != nil {
		panic("config file not found: " + load.Args.ConfigFile)
	}
	path := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)
	config.LoadFiles(&configs, files, path)

	// when
	config.RunFiles(&configs)

	// 'du' return one line per dir + total UNLESS we use 'summary' flag, then it return 2 lines
	// - value, dirname
	// - value, total
	// we're only interest in the first (index 0)
	// then
	metricsSet := load.Entity.Metrics[0]
	assert.NotEmpty(t, metricsSet)

	// these were the names we gave the 2 'columns' of the command result
	assert.NotNil(t, metricsSet.Metrics["dirSizeBytes"])
	assert.NotNil(t, metricsSet.Metrics["dirName"])
}

func TestConfig_cmd_LinuxDiskFree(t *testing.T) {
	// given
	load.Refresh()

	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")
	load.Args.ConfigFile = "configs/linux-disk-free.yml"

	// Read a single config file
	var files []os.FileInfo
	var configs []load.Config
	file, err := os.Stat(load.Args.ConfigFile)
	if err != nil {
		panic("config file not found: " + load.Args.ConfigFile)
	}
	path := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)
	config.LoadFiles(&configs, files, path)

	// when
	config.RunFiles(&configs)

	// fs,fsType,usedBytes,availableBytes,usedPerc,mountedOn
	metricsSet := load.Entity.Metrics
	assert.NotEmpty(t, metricsSet)
	for _, ms := range metricsSet {
		assert.NotNil(t, ms.Metrics["fs"])
		assert.NotNil(t, ms.Metrics["fsType"])
		assert.NotNil(t, ms.Metrics["usedBytes"])
		assert.NotNil(t, ms.Metrics["availableBytes"])
		assert.NotNil(t, ms.Metrics["usedPerc"])
		assert.NotNil(t, ms.Metrics["mountedOn"])
	}
}

func TestConfig_cmd_OpenFDs(t *testing.T) {
	// given
	load.Refresh()

	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")
	load.Args.ConfigFile = "../examples/linux/linux-open-fds.yml"

	// Read a single config file
	var files []os.FileInfo
	var configs []load.Config
	file, err := os.Stat(load.Args.ConfigFile)
	if err != nil {
		panic("config file not found: " + load.Args.ConfigFile)
	}
	path := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)
	config.LoadFiles(&configs, files, path)

	// when
	config.RunFiles(&configs)

	// openFD,maxFD
	// 1 record only
	metricsSet := load.Entity.Metrics
	assert.NotEmpty(t, metricsSet)
	assert.Len(t, metricsSet, 1)
	for _, ms := range metricsSet {
		fmt.Println(fmt.Sprint(ms))
		assert.NotNil(t, ms.Metrics["openFD"])
		assert.GreaterOrEqual(t, ms.Metrics["openFD"].(float64), float64(0))
		assert.NotNil(t, ms.Metrics["maxFD"])
		assert.GreaterOrEqual(t, ms.Metrics["maxFD"].(float64), float64(0))
	}
}

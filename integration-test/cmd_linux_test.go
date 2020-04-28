// +build integration
// +build linux

// The tests in this file are supposed to be run in the CI using docker-compose
// You can run then from the IDE or manually but you'll need some setup first
package integration_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/stretchr/testify/assert"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
)

func tmpData() (data []byte) {
	data = make([]byte, 512)
	for i := range data {
		data[i] = 1
	}
	return data
}

func TestConfig_cmd_LinuxDirUsage(t *testing.T) {
	// given
	load.Refresh()

	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("IntegrationTest", "nri-flex")
	load.Args.ConfigFile = "configs/linux-directory-size.yml"

	tmpDir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	t.Logf("Tmp dir is: %s", tmpDir)
	tmpFile := path.Join(tmpDir, "test_file")
	err = ioutil.WriteFile(tmpFile, tmpData(), 0644)
	require.NoError(t, err)

	err = os.Setenv("TMP_TEST_DIR", tmpDir)
	require.NoError(t, err)
	defer os.Unsetenv("TMP_TEST_DIR")

	// Read a single config file
	var files []os.FileInfo
	var configs []load.Config
	file, err := os.Stat(load.Args.ConfigFile)
	require.NoError(t, err, "config file not found: %s", load.Args.ConfigFile)
	path := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)
	errs := config.LoadFiles(&configs, files, path)
	require.Empty(t, errs)

	// when
	errs = config.RunFiles(&configs)
	require.Empty(t, errs)

	// 'du' return one line per dir + total UNLESS we use 'summary' flag, then it return 2 lines
	// - value, dirname
	// - value, total
	// we're only interest in the first (index 0)
	// then
	metricsSet := load.Entity.Metrics[0]
	assert.NotEmpty(t, metricsSet)

	// these were the names we gave the 2 'columns' of the command result
	assert.Equal(t, float64(8), metricsSet.Metrics["dirSizeBytes"])
	assert.Equal(t, tmpDir, metricsSet.Metrics["dirName"])
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
	errs := config.LoadFiles(&configs, files, path)
	require.Empty(t, errs)

	// when
	errs = config.RunFiles(&configs)
	require.Empty(t, errs)

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
	load.Args.ConfigFile = "../examples/flexConfigs/linux-open-fds.yml"

	// Read a single config file
	var files []os.FileInfo
	var configs []load.Config
	file, err := os.Stat(load.Args.ConfigFile)
	if err != nil {
		panic("config file not found: " + load.Args.ConfigFile)
	}
	path := strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
	files = append(files, file)
	errs := config.LoadFiles(&configs, files, path)
	require.Empty(t, errs)

	// when
	errs = config.RunFiles(&configs)
	require.Empty(t, errs)

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

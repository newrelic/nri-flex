package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/newrelic/nri-flex/internal/discovery"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	"github.com/newrelic/nri-flex/internal/outputs"
	"github.com/newrelic/nri-flex/internal/processor"
)

func main() {
	load.FlexStatusCounter.M["EventCount"] = 0
	load.FlexStatusCounter.M["EventDropCount"] = 0
	load.FlexStatusCounter.M["ConfigsProcessed"] = 0

	outputs.CreateIntegration()
	logger.Flex("info", nil, fmt.Sprintf("%v: v%v", load.IntegrationName, load.IntegrationVersion), false)

	// todo: port cluster mode here
	runIntegration()
	logger.Flex("fatal", load.Integration.Publish(), "unable to publish", false)
}

func runIntegration() {
	// store config ymls
	var configs []load.Config

	// Check if specific config file has been specified
	// if not check flexConfigs dir and run container discovery if enabled
	var path string
	var files []os.FileInfo
	if load.Args.ConfigFile != "" {
		// Read a single config file
		file, err := os.Stat(load.Args.ConfigFile)
		logger.Flex("fatal", err, "failed to read specified config file: "+load.Args.ConfigFile, false)
		path = strings.Replace(filepath.FromSlash(load.Args.ConfigFile), file.Name(), "", -1)
		files = append(files, file)
	} else {
		// List config files in directory
		path = filepath.FromSlash(load.Args.ConfigDir)
		var err error
		files, err = ioutil.ReadDir(path)
		logger.Flex("fatal", err, "failed to read config dir: "+load.Args.ConfigDir, false)
		if load.Args.ContainerDiscovery {
			discovery.Run(&configs)
		}
	}

	processor.LoadConfigFiles(&configs, files, path)
	processor.RunConfigFiles(&configs)
	outputs.CreateStatusSample()
	if load.Args.InsightsURL != "" && load.Args.InsightsAPIKey != "" {
		outputs.SendToInsights()
	}
}

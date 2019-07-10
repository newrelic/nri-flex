package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/discovery"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	"github.com/newrelic/nri-flex/internal/outputs"
	"github.com/sirupsen/logrus"
)

func main() {
	load.Logrus.Out = os.Stdout
	load.FlexStatusCounter.M["EventCount"] = 0
	load.FlexStatusCounter.M["EventDropCount"] = 0
	load.FlexStatusCounter.M["ConfigsProcessed"] = 0

	discovery.Processes()
	outputs.InfraIntegration()
	outputs.LambdaCheck()
	runIntegration()

	if outputs.LambdaEnabled {
		outputs.LambdaFinish()
	}

	logger.Flex("fatal", load.Integration.Publish(), "unable to publish", false)
}

// runIntegration runs nri-flex
func runIntegration() {
	if load.Args.Verbose && os.Getenv("VERBOSE") != "true" && os.Getenv("VERBOSE") != "1" {
		load.Logrus.SetLevel(logrus.TraceLevel)
	}

	logger.Flex("debug", nil, fmt.Sprintf("%v: v%v %v:%v", load.IntegrationName, load.IntegrationVersion, runtime.GOOS, runtime.GOARCH), false)

	if load.Args.GitService != "" && load.Args.GitToken != "" && load.Args.GitUser != "" && load.Args.GitRepo != "" {
		config.SyncGitConfigs()
	}

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
		config.LoadFiles(&configs, files, path)
	} else {
		var err error
		// List config files in directory
		path = filepath.FromSlash(load.Args.ConfigDir)
		files, err = ioutil.ReadDir(path)
		config.LoadFiles(&configs, files, path)
		logger.Flex("fatal", err, fmt.Sprintf("failed to read config dir %v", load.Args.ConfigDir), false)
		if load.Args.ContainerDiscovery {
			discovery.Run(&configs)
		}
	}

	config.RunFiles(&configs)
	outputs.StatusSample()
	if load.Args.InsightsURL != "" && load.Args.InsightsAPIKey != "" {
		outputs.SendToInsights()
	} else if load.Args.MetricAPIUrl != "" && (load.Args.InsightsAPIKey != "" || load.Args.MetricAPIKey != "") && len(load.MetricsStore.Data) > 0 {
		outputs.SendToMetricAPI()
	} else if len(load.MetricsStore.Data) > 0 && (load.Args.MetricAPIUrl == "" || (load.Args.InsightsAPIKey == "" || load.Args.MetricAPIKey == "")) {
		logger.Flex("debug", nil, "metric_api is being used, but metric url and/or key has not been set", false)
	}
}

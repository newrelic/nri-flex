package main

import (
	"fmt"
	"io/ioutil"
	"nri-flex/internal/discovery"
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
	"nri-flex/internal/outputs"
	"nri-flex/internal/processor"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

func main() {
	outputs.CreateIntegration()
	if load.Args.Version {
		logger.Flex("info", nil, fmt.Sprintf("%v: v%v", load.IntegrationName, load.IntegrationVersion), false)
	}
	if load.Args.ForceLog {
		logger.Flex("fatal", fmt.Errorf("-force_log deprecated use -verbose instead"), "", false)
	}
	// todo: port cluster mode here
	runIntegration()
	err := load.Integration.Publish()
	logger.Flex("fatal", err, "unable to publish", false)
}

func runIntegration() {
	// set current containerID
	cpuset, err := ioutil.ReadFile("/proc/1/cpuset")
	// output eg. /kubepods/besteffort/podaa8aee52-49b6-11e9-95e2-080027000d3d/d49ee19ddec683e0cd80ca881a27d45a88105f8c439a4c9d5607b675341e394e
	if err == nil {
		strCpuset := strings.TrimSpace(string(cpuset))
		logger.Flex("info", nil, fmt.Sprintf("cpuset: %v", strCpuset), false)
		values := strings.Split(strCpuset, "/")
		if len(values) > 0 {
			if len(values[len(values)-1]) == 64 {
				load.ContainerID = values[len(values)-1]
				logger.Flex("debug", fmt.Errorf("flex container id: %v", load.ContainerID), "", false)
			}
		}
	} else {
		logger.Flex("debug", err, "potentially not running within a container", false)
	}

	// set hostname
	load.Hostname, _ = os.Hostname()

	// store config ymls
	var ymls []load.Config

	containerDiscoveryAvailable := false
	var containers []types.Container
	if load.Args.ContainerDiscovery {
		discovery.Run(&containerDiscoveryAvailable, &containers)
	}

	// Check if specific config file has been specified if not check flexConfigs directory
	var path, containerDiscoveryPath string
	var files, containerDiscoveryFiles []os.FileInfo
	if load.Args.ConfigFile != "" {
		// Read config file
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

		// List config files in containerDiscoveryPath directory
		if containerDiscoveryAvailable {
			containerDiscoveryPath = filepath.FromSlash(load.Args.ContainerDiscoveryDir)
			containerDiscoveryFiles, err = ioutil.ReadDir(containerDiscoveryPath)
			logger.Flex("debug", err, "failed to read config dir: "+load.Args.ContainerDiscoveryDir, false)

			if len(containerDiscoveryFiles) == 0 {
				containerDiscoveryAvailable = false
				logger.Flex("info", nil, "no configs found: "+load.Args.ContainerDiscoveryDir, false)
			} else if len(containerDiscoveryFiles) > 0 && err == nil {
				discovery.CreateDynamicContainerConfigs(containers, containerDiscoveryFiles, containerDiscoveryPath, &ymls)
			}
		}
	}

	processor.LoadConfigFiles(&ymls, files, path) // load standard configs if available
	processor.YmlFiles(&ymls)

	// flexStatusSample
	flexStatusSample := load.Entity.NewMetricSet("flexStatusSample")
	logger.Flex("debug", flexStatusSample.SetMetric("eventCount", load.EventCount, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("eventDropCount", load.EventDropCount, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("configsProcessed", load.ConfigsProcessed, metric.GAUGE), "", false)
	for sample, count := range load.EventDistribution {
		logger.Flex("debug", flexStatusSample.SetMetric(sample+"_count", count, metric.GAUGE), "", false)
	}

	// SendToInsights
	if load.Args.InsightsURL != "" && load.Args.InsightsAPIKey != "" {
		outputs.SendToInsights()
	}

}

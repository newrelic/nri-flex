/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// runFargateDiscovery check aws metadata endpoint for containers
func runFargateDiscovery(configs *[]load.Config) {
	load.Logrus.Debug("discovery: running fargate discovery")
	client := &http.Client{}
	resp, err := client.Get("http://169.254.170.2/v2/metadata")
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("discovery: fargate aws v2 metadata fetch failed")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		load.IsFargate = true
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("discovery: fargate response body read fail")
		} else {
			TaskMetadata := load.TaskMetadata{}
			err := json.Unmarshal(bodyBytes, &TaskMetadata)
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("discovery: fargate responde body unmarshal fail")
			} else {
				determineDynamicFargateConfigs(configs, TaskMetadata)
			}
		}
	}
}

// determineDynamicFargateConfigs
func determineDynamicFargateConfigs(configs *[]load.Config, TaskMetadata load.TaskMetadata) {
	for _, currentConfig := range *configs {
		if strings.HasPrefix(currentConfig.FileName, "cd-") {
			if currentConfig.ContainerDiscovery.Mode == "" {
				currentConfig.ContainerDiscovery.Mode = "contains"
			}
			if currentConfig.ContainerDiscovery.Type == "" {
				currentConfig.ContainerDiscovery.Type = load.TypeContainer
			}
			if currentConfig.ContainerDiscovery.Target == "" {
				target := strings.Replace(strings.TrimPrefix(currentConfig.FileName, "cd-"), ".yml", "", -1)
				currentConfig.ContainerDiscovery.Target = strings.Replace(target, ".yaml", "", -1)
			}

			// check all containers async
			var wg sync.WaitGroup
			wg.Add(len(TaskMetadata.Containers))
			for _, container := range TaskMetadata.Containers {
				go func(container load.Container) {
					defer wg.Done()
					// do not target the flex container
					if container.DockerID != load.ContainerID {
						if checkContainerMatch(container, currentConfig.ContainerDiscovery) {
							load.Logrus.Debug(fmt.Sprintf("discovery: fargate lookup matched %v - file %v", container.DockerID, currentConfig.FileName))
							if len(container.Networks) > 0 {
								if len(container.Networks[0].IPv4Addresses) > 0 {
									addDynamicFargateConfig(configs, currentConfig, container)
								} else {
									load.Logrus.Debug(fmt.Sprintf("discovery: fargate container %v file %v - does not have any IPv4 Addresses configured", container.DockerID, currentConfig.FileName))
								}
							} else {
								load.Logrus.Debug(fmt.Sprintf("discovery: fargate container %v file %v - does not have any networks configured", container.DockerID, currentConfig.FileName))
							}
						}
					}
				}(container)
			}
			wg.Wait()
		}
	}
}

func checkContainerMatch(container load.Container, containerDiscovery load.ContainerDiscovery) bool {
	switch containerDiscovery.Type {
	case "cname", load.TypeContainer:
		if formatter.KvFinder(containerDiscovery.Mode, container.Name, containerDiscovery.Target) {
			return true
		}
	case load.Img, load.Image:
		if formatter.KvFinder(containerDiscovery.Mode, container.Image, containerDiscovery.Target) {
			return true
		}
	default:
		load.Logrus.Debug("discovery: fargate targetType not set id: " + container.DockerID)
	}
	return false
}

func addDynamicFargateConfig(configs *[]load.Config, currentConfig load.Config, container load.Container) {
	tmpCfgBytes, err := yaml.Marshal(&currentConfig)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("discovery: fargate marshal config failed")
	} else {
		tmpCfgStr := string(tmpCfgBytes)
		fargateIP := container.Networks[0].IPv4Addresses[0]
		tmpCfgStr = strings.Replace(tmpCfgStr, "${auto:host}", fargateIP, -1)
		tmpCfgStr = strings.Replace(tmpCfgStr, "${auto:ip}", fargateIP, -1)
		newConfig, err := config.ReadYML(tmpCfgStr)

		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"name": newConfig.Name,
				"err":  err,
			}).Error("discovery: fargate unmarshal config failed")
		} else {
			newConfig.ContainerDiscovery.ReplaceComplete = true

			//add extra attributes
			newConfig.CustomAttributes["containerId"] = container.DockerID
			newConfig.CustomAttributes["containerName"] = container.Name
			newConfig.CustomAttributes["image"] = container.Image
			newConfig.CustomAttributes["imageId"] = container.ImageID
			for key, val := range container.Labels {
				newConfig.CustomAttributes[key] = val
			}

			*configs = append(*configs, newConfig)
		}
	}
}

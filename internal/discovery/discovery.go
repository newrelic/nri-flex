/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Run discover containers
func Run(configs *[]load.Config) {
	FindFlexContainerID("/proc/1/cpuset")
	if load.Args.Fargate {
		runFargateDiscovery(configs)
	} else {
		var err error
		cli, err = setDockerClient()
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("discovery: unable to set docker client")
		} else {
			ctx := context.Background()
			containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{})
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("discovery: unable to set perform container list")
			} else if len(containerList) > 0 {
				// List config files in containerDiscoveryPath directory
				containerDiscoveryPath := filepath.FromSlash(load.Args.ContainerDiscoveryDir)
				containerDiscoveryFiles, err := ioutil.ReadDir(containerDiscoveryPath)

				if err != nil {
					load.Logrus.WithFields(logrus.Fields{
						"err": err,
					}).Debug(fmt.Sprintf("discovery: %v directory unavailable or not used", load.Args.ContainerDiscoveryDir))
				}

				CreateDynamicContainerConfigs(containerList, containerDiscoveryFiles, containerDiscoveryPath, configs)

				if err == nil && len(containerDiscoveryFiles) == 0 {
					load.Logrus.WithFields(logrus.Fields{
						"message": "if you are using v2 discovery then ignore",
						"dir":     load.Args.ContainerDiscoveryDir,
					}).Debug("discovery: no configs found")
				}
			}
		}
	}
}

// FindFlexContainerID detects if Flex is running within a container and sets the ID
func FindFlexContainerID(read string) {
	// set current containerID
	cpuset, err := ioutil.ReadFile(read) // read = "/proc/1/cpuset"
	// output eg. /kubepods/besteffort/podaa8aee52-49b6-11e9-95e2-080027000d3d/d49ee19ddec683e0cd80ca881a27d45a88105f8c439a4c9d5607b675341e394e
	if err == nil {
		strCpuset := strings.TrimSpace(string(cpuset))

		load.Logrus.WithFields(logrus.Fields{
			"read":   read,
			"output": strCpuset,
		}).Debug("discovery: finding flex container id")

		if strings.Contains(strCpuset, "kube") {
			load.IsKubernetes = true
		}

		values := strings.Split(strCpuset, "/")
		if len(values) > 0 {
			if len(values[len(values)-1]) == 64 {
				load.ContainerID = values[len(values)-1]
				load.Logrus.WithFields(logrus.Fields{
					"containerId": load.ContainerID,
				}).Debug("discovery: flex container id found")
			}
		}
		// fallback on cgroup
		if load.ContainerID == "" && read != "/proc/self/cgroup" {
			FindFlexContainerID("/proc/self/cgroup")
		}
	} else {
		load.Logrus.Debug("discovery: flex potentially not running within a container")
	}
}

// CreateDynamicContainerConfigs Creates dynamic configs for each container
func CreateDynamicContainerConfigs(containers []types.Container, files []os.FileInfo, path string, ymls *[]load.Config) {
	var containerYmls []load.Config
	config.LoadFiles(&containerYmls, files, path)
	foundTargetContainerIds := []string{}

	// store inspected containers, so we do not re-inspect anything unnecessarily
	inspectedContainers := []types.ContainerJSON{}
	// find flex container id if not set
	if load.ContainerID == "" {
		fallbackFindFlexContainerID(&containers)
	}

	// filter containers out
	filteredContainers := []types.Container{}
	for _, container := range containers {
		// do not target self or nr images
		if container.ID == load.ContainerID || strings.HasPrefix(container.Image, "newrelic/") {
			continue
		}
		// do not target nr sidecar
		if len(container.Names) > 0 {
			if strings.HasPrefix(container.Names[0], "/k8s_newrelic-sidecar") {
				continue
			}
		}
		filteredContainers = append(filteredContainers, container)
	}

	load.Logrus.Debug(fmt.Sprintf("discovery: containers %d filtered containers %d", len(containers), len(filteredContainers)))

	if load.Args.ContainerDump {
		for _, container := range filteredContainers {
			load.Logrus.Debug(fmt.Sprintf("container: %v %v %v", container.ID, container.Image, container.Names))
		}
	}

	// flex config file container_discovery parameter -> container
	runConfigLookup(cli, &filteredContainers, &inspectedContainers, &foundTargetContainerIds, ymls)
	// flex envs/labels -> container
	runReverseLookup(cli, &filteredContainers, &inspectedContainers, &foundTargetContainerIds, &containerYmls, ymls, path)
	// container envs/labels -> flex
	runForwardLookup(cli, &filteredContainers, &inspectedContainers, &foundTargetContainerIds, &containerYmls, ymls, path)
}

// runForwardLookup container envs -> flex
func runForwardLookup(dockerClient *client.Client, containers *[]types.Container, inspectedContainers *[]types.ContainerJSON, foundTargetContainerIds *[]string, containerYmls *[]load.Config, ymls *[]load.Config, path string) {

	var wg sync.WaitGroup
	wg.Add(len(*containers))
	for _, container := range *containers {
		go func(container types.Container) {
			defer wg.Done()

			// do not target already targeted containers unless container discovery multi param is true
			targeted := false
			if !load.Args.ContainerDiscoveryMulti {
				for _, foundTargetContainerID := range *foundTargetContainerIds {
					if foundTargetContainerID == container.ID {
						targeted = true
					}
				}
			}

			// do not do a forward lookup against the flex container
			if load.ContainerID != container.ID && !targeted {

				discoveryLoop := map[string]string{}
				// add container labels to check for disc configs
				for key, val := range container.Labels {
					discoveryLoop[key] = val
				}

				// check env variables for disc configs, so use container inspect
				var containerInspect types.ContainerJSON

				// check if the container has already been inspected
				for _, inspectedContainer := range *inspectedContainers {
					if inspectedContainer.ID == container.ID {
						containerInspect = inspectedContainer
						break
					}
				}

				if containerInspect.Config == nil {
					ctx := context.Background()
					var err error
					containerInspect, err = dockerClient.ContainerInspect(ctx, container.ID)
					if err != nil {
						load.Logrus.WithFields(logrus.Fields{
							"err": err,
						}).Error("discovery: container inspect failed")
					} else {
						*inspectedContainers = append(*inspectedContainers, containerInspect)
					}
				}

				if containerInspect.Config != nil {
					// add container labels to check for disc configs
					for key, val := range containerInspect.Config.Labels {
						discoveryLoop[key] = val
					}

					// add env variables to check for disc configs
					for _, envVar := range containerInspect.Config.Env {
						environmentVar := strings.SplitN(envVar, "=", 2)
						if len(environmentVar) == 2 {
							discoveryLoop[environmentVar[0]] = environmentVar[1]
						}
					}

					// create discoveryConfigs - look for flex label and split
					for key, val := range discoveryLoop {
						if strings.Contains(key, "flexDiscovery") {

							load.Logrus.Debug(fmt.Sprintf("discovery: fwd lookup for %v", key))

							discoveryConfig := map[string]interface{}{}
							parseFlexDiscoveryLabel(discoveryConfig, key, val)
							// t = target, c = config, r = reverse, tt = target type, tm = target mode, ip = ip mode, p = port
							// check if we have a target to find, and config to run
							if discoveryConfig["t"] != nil {
								// if config is nil, use the <target> , as the yaml file to look up eg. if target (t) = redis, lookup the config (c) redis.yml
								if discoveryConfig["c"] == nil {
									discoveryConfig["c"] = discoveryConfig["t"]
								}
								// auto will mean that if set to true, it will loop through all other containers to find a match
								// if not set / set to false it will target the current container
								if discoveryConfig["r"] == nil {
									discoveryConfig["r"] = "false"
								}
								if discoveryConfig["tt"] == nil {
									discoveryConfig["tt"] = load.Img // cname == containerName , img = image
								}
								if discoveryConfig["tm"] == nil {
									discoveryConfig["tm"] = load.Contains
								}

								// addDynamicConfig will ensure the config file matches, so the above parameters are no longer enforced
								addDynamicConfig(containerYmls, discoveryConfig, ymls, container, containerInspect, path)
							}
						}
					}
				}
			}
		}(container)
	}
	wg.Wait()
}

// runConfigLookup flex config -> container
func runConfigLookup(dockerClient *client.Client, containers *[]types.Container, inspectedContainers *[]types.ContainerJSON, foundTargetContainerIds *[]string, ymls *[]load.Config) {
	discoveryLoop := []load.ContainerDiscovery{}

	// this can probably be handled differently
	for _, yml := range *ymls {
		if strings.HasPrefix(yml.FileName, "cd-") {
			target := strings.SplitAfter(yml.FileName, "cd-")[1]
			target = strings.Replace(target, ".yml", "", -1)
			target = strings.Replace(target, ".yaml", "", -1)
			targetType := load.TypeContainer
			mode := "contains"

			cd := load.ContainerDiscovery{
				FileName: yml.FileName,
				Target:   target,
				Type:     targetType,
				Mode:     mode,
			}

			if yml.ContainerDiscovery.Mode != "" {
				cd.Mode = yml.ContainerDiscovery.Mode
			}
			if yml.ContainerDiscovery.Type != "" {
				cd.Type = yml.ContainerDiscovery.Type
			}
			if yml.ContainerDiscovery.Target != "" {
				cd.Target = yml.ContainerDiscovery.Target
			}
			if yml.ContainerDiscovery.IPMode != "" {
				cd.IPMode = yml.ContainerDiscovery.IPMode
			}

			discoveryLoop = append(discoveryLoop, cd)
		}
	}

	// check our discovery loop continue reverse lookup
	if len(discoveryLoop) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(*containers))
		for _, container := range *containers {
			go func(container types.Container) {
				defer wg.Done()
				// create discoveryConfigs
				for _, cd := range discoveryLoop {
					discoveryConfig := map[string]interface{}{}
					discoveryConfig["t"] = cd.Target
					discoveryConfig["c"] = cd.FileName
					discoveryConfig["tt"] = cd.Type
					discoveryConfig["tm"] = cd.Mode

					ctx := context.Background()
					reverseContainerInspect, err := dockerClient.ContainerInspect(ctx, container.ID)
					if err != nil {
						load.Logrus.WithFields(logrus.Fields{
							"err": err,
						}).Error(fmt.Sprintf("discovery: %v cfg container inspect failed - %v", container.ID, cd.FileName))
					} else {
						if findContainerTarget(discoveryConfig, container, foundTargetContainerIds) {

							switch discoveryConfig["tt"].(string) {
							case load.TypeCname, load.TypeContainer:
								load.Logrus.Debug(fmt.Sprintf("discovery: %v cfg lookup matched %v %v - %v", container.ID, container.Names, cd.Target, cd.FileName))
							case load.Img, load.Image:
								load.Logrus.Debug(fmt.Sprintf("discovery: %v cfg lookup matched %v %v - %v", container.ID, container.Image, cd.Target, cd.FileName))
							}

							*inspectedContainers = append(*inspectedContainers, reverseContainerInspect)
							addDynamicConfig(ymls, discoveryConfig, ymls, container, reverseContainerInspect, "")
						}
					}
				}
			}(container)
		}
		wg.Wait()
	}

}

// runReverseLookup flex envs -> container
func runReverseLookup(dockerClient *client.Client, containers *[]types.Container, inspectedContainers *[]types.ContainerJSON, foundTargetContainerIds *[]string, containerYmls *[]load.Config, ymls *[]load.Config, path string) {
	var flexContainerInspect types.ContainerJSON
	ctx := context.Background()
	discoveryLoop := map[string]string{}

	if load.ContainerID != "" {
		var err error
		flexContainerInspect, err = dockerClient.ContainerInspect(ctx, load.ContainerID)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("discovery: container inspect failed")
		} else if flexContainerInspect.Config != nil {
			*inspectedContainers = append(*inspectedContainers, flexContainerInspect)

			// add container labels to check for disc configs
			for key, val := range flexContainerInspect.Config.Labels {
				discoveryLoop[key] = val
			}

			// add env variables to check for disc configs
			for _, envVar := range flexContainerInspect.Config.Env {
				environmentVar := strings.SplitN(envVar, "=", 2)
				if len(environmentVar) == 2 {
					discoveryLoop[environmentVar[0]] = environmentVar[1]
				}
			}
		}
	}

	// check our discovery loop continue reverse lookup
	if len(discoveryLoop) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(*containers))
		for _, container := range *containers {
			go func(container types.Container) {
				defer wg.Done()
				for key, val := range discoveryLoop {
					if strings.Contains(key, "flexDiscovery") {
						discoveryConfig := map[string]interface{}{}
						parseFlexDiscoveryLabel(discoveryConfig, key, val)
						// t = target, c = config, r = reverse, tt = target type, tm = target mode, ip = ip mode, p = port
						// check if we have a target to find, and config to run
						if discoveryConfig["t"] != nil {
							// if config is nil, use the <target> , as the yaml file to look up eg. if target (t) = redis, lookup the config (c) redis.yml
							if discoveryConfig["c"] == nil {
								discoveryConfig["c"] = discoveryConfig["t"]
							}
							if discoveryConfig["tt"] == nil {
								discoveryConfig["tt"] = load.Img // cname == containerName , img = image
							}
							if discoveryConfig["tm"] == nil {
								discoveryConfig["tm"] = load.Contains
							}
						}
						ctx := context.Background()
						reverseContainerInspect, err := dockerClient.ContainerInspect(ctx, container.ID)
						if err != nil {
							load.Logrus.WithFields(logrus.Fields{
								"err": err,
							}).Error(fmt.Sprintf("discovery: rev container inspect failed on cid:%v key:%v val:%v", container.ID, key, val))
						} else {
							if findContainerTarget(discoveryConfig, container, foundTargetContainerIds) {
								load.Logrus.Debug(fmt.Sprintf("discovery: rev lookup matched %v: %v - %v", container.ID, key, val))
								*inspectedContainers = append(*inspectedContainers, reverseContainerInspect)
								addDynamicConfig(containerYmls, discoveryConfig, ymls, container, reverseContainerInspect, path)
							}
						}
					}
				}
			}(container)
		}
		wg.Wait()
	}

}

func addDynamicConfig(containerYmls *[]load.Config, discoveryConfig map[string]interface{}, ymls *[]load.Config, targetContainer types.Container, targetContainerInspect types.ContainerJSON, path string) {
	for _, containerYml := range *containerYmls {
		configName := ""
		switch cfg := discoveryConfig["c"].(type) {
		case string:
			configName = cfg
			if !strings.HasSuffix(cfg, ".yml") {
				configName = configName + ".yml"
			}
		default:
			load.Logrus.Error(fmt.Sprintf("discovery: config file error %v", (discoveryConfig["c"])))
		}

		if containerYml.FileName == configName {
			load.Logrus.Debug(fmt.Sprintf("discovery: %v matched %v", targetContainer.ID, containerYml.FileName))
			if path == "" {
				path = containerYml.FilePath
			}
			b, err := ioutil.ReadFile(path + containerYml.FileName)
			if err != nil {
				load.Logrus.Error(fmt.Sprintf("discovery: unable to read flex config: " + path + containerYml.FileName))
			} else {
				ymlString := string(b)
				config.SubEnvVariables(&ymlString)
				config.SubTimestamps(&ymlString)
				discoveryIPAddress := "" // we require IP at least
				discoveryPort := ""      // we don't require port
				networkIPAddress := ""
				privatePort := ""
				publicIPAddress := ""
				publicPort := ""
				// use the first found IPAddress
				for _, network := range targetContainer.NetworkSettings.Networks {
					networkIPAddress = network.IPAddress
					break
				}
				// ability to override and select port
				if discoveryConfig["p"] != nil {
					discoveryPort = discoveryConfig["p"].(string)
				} else {
					// use the first found public port
					for _, port := range targetContainer.Ports {
						publicIPAddress = fmt.Sprintf("%v", port.IP)
						publicPort = fmt.Sprintf("%v", port.PublicPort)
						privatePort = fmt.Sprintf("%v", port.PrivatePort)
						break
					}
				}

				ipMode := ""
				if load.ContainerID == "" {
					ipMode = "public"
				} else {
					ipMode = "private"
				}

				if load.Args.OverrideIPMode != "" && (load.Args.OverrideIPMode == load.Public || load.Args.OverrideIPMode == load.Private) {
					ipMode = load.Args.OverrideIPMode
				} else if discoveryConfig["ip"] != nil {
					if discoveryConfig["ip"].(string) == load.Private || discoveryConfig["ip"].(string) == load.Public {
						ipMode = discoveryConfig["ip"].(string)
					}
				}

				switch ipMode {
				case load.Public:
					discoveryIPAddress = publicIPAddress
					discoveryPort = publicPort
				case load.Private:
					discoveryIPAddress = networkIPAddress
					discoveryPort = privatePort
				}

				if discoveryIPAddress != "" {
					ymlString = strings.Replace(ymlString, "${auto:ip}", discoveryIPAddress, -1)
					ymlString = strings.Replace(ymlString, "${auto:host}", discoveryIPAddress, -1)
				}

				// attempt low level ip fetch
				if discoveryIPAddress == "" {
					ipAddresses := lowLevelIpv4Fetch(&discoveryIPAddress, targetContainerInspect.State.Pid, targetContainer.ID)
					for i, address := range ipAddresses {
						if i == 0 {
							ymlString = strings.Replace(ymlString, "${auto:ip}", address, -1)
							ymlString = strings.Replace(ymlString, "${auto:host}", address, -1)
							discoveryIPAddress = address
						}
						ymlString = strings.Replace(ymlString, fmt.Sprintf("${auto:ip[%d]}", i), address, -1)
						ymlString = strings.Replace(ymlString, fmt.Sprintf("${auto:host[%d]}", i), address, -1)
					}
				}

				// attempt hostname fallback
				if discoveryIPAddress == "" {
					execHostnameFallback(&discoveryIPAddress, targetContainer.ID, &ymlString)
				}

				// substitute port into yml
				if discoveryPort != "" && discoveryPort != "0" {
					ymlString = strings.Replace(ymlString, "${auto:port}", discoveryPort, -1)
				} else {
					portFallback(targetContainer, targetContainerInspect, &ymlString, &discoveryPort)
				}

				load.Logrus.Debug(fmt.Sprintf("discovery: %v %v - %v - %v:%v", targetContainer.ID, containerYml.FileName, ipMode, discoveryIPAddress, discoveryPort))

				if strings.Contains(ymlString, "${auto:host}") || strings.Contains(ymlString, "${auto:ip}") || strings.Contains(ymlString, "${auto:port}") {
					containerName := ""
					if len(targetContainer.Names) > 0 {
						containerName = targetContainer.Names[0]
					}
					load.Logrus.Debug(fmt.Sprintf("discovery: %v %v couldn't build dynamic cfg", targetContainer.ID, containerName))
					load.Logrus.Debug(fmt.Sprintf("discovery: %v %v missing variable unable to create dynamic cfg ip:%v - port:%v", targetContainer.ID, containerName, discoveryIPAddress, discoveryPort))
				} else {
					yml, err := config.ReadYML(ymlString)
					if err != nil {
						load.Logrus.WithFields(logrus.Fields{
							"err":  err,
							"file": path + containerYml.FileName,
						}).Error("discovery: unable to unmarshal flexConfig")
						load.Logrus.Error(ymlString)
					} else {

						// decorate additional docker attributes
						if yml.CustomAttributes == nil {
							yml.CustomAttributes = map[string]string{}
						}
						for key, val := range targetContainer.Labels {
							yml.CustomAttributes[key] = val
						}
						// If Kubernetes, find the labels in metadata
						if load.IsKubernetes {
							// Get the pod name & namespace from labels above
							podName := targetContainer.Labels["io.kubernetes.pod.name"]
							podNamespace := targetContainer.Labels["io.kubernetes.pod.namespace"]
							load.Logrus.Debug(fmt.Sprintf("discovery: fetching k8 labels for  %v in namespace  %v ", podName, podNamespace))
							kubeLabels := getK8Labels(podName, podNamespace)
							for key, val := range kubeLabels {
								yml.CustomAttributes[key] = val
								//load.Logrus.Debug(fmt.Sprintf("discovery: adding kubeLabels:  %v : %v ", key, val))
							}
						}

						yml.CustomAttributes["containerId"] = targetContainer.ID
						yml.CustomAttributes["imageId"] = targetContainer.Image
						yml.CustomAttributes["IDShort"] = targetContainer.ID[0:12]
						yml.CustomAttributes["image"] = targetContainer.Image
						img := strings.Split(targetContainer.Image, "@")
						yml.CustomAttributes["imageShort"] = img[0]
						yml.CustomAttributes["container.duration"] = fmt.Sprintf("%d", MakeTimestamp()-(targetContainer.Created*1000))
						yml.CustomAttributes["container.state"] = targetContainer.State
						yml.CustomAttributes["container.status"] = targetContainer.Status
						yml.CustomAttributes["container.name"] = strings.TrimPrefix(targetContainerInspect.Name, "/")

						*ymls = append(*ymls, yml)
					}
				}

			}

		} else {
			load.Logrus.Debug(fmt.Sprintf("discovery: %v container containerFileName %v did not match configName %v", targetContainer.ID, containerYml.FileName, configName))
		}
	}
}

func parseFlexDiscoveryLabel(discoveryConfig map[string]interface{}, key string, val string) {
	if strings.Contains(val, "=") { // nicer for other setups
		labelValues := strings.Split(val, ",")
		for _, value := range labelValues {
			configKeyPair := strings.Split(value, "=")
			if len(configKeyPair) == 2 {
				discoveryConfig[configKeyPair[0]] = configKeyPair[1]
			}
		}
	} else if strings.Contains(val, ".") { // needed for kubernetes eg. flexDiscoveryRedis:"t_redis.c_redis.tt_img.tm_contains"
		labelValues := strings.Split(val, ".")
		for _, value := range labelValues {
			configKeyPair := strings.Split(value, "_")
			if len(configKeyPair) == 2 {
				discoveryConfig[configKeyPair[0]] = configKeyPair[1]
			}
		}
	}
}

func findContainerTarget(discoveryConfig map[string]interface{}, container types.Container, foundTargetContainerIds *[]string) bool {
	// do not add any dynamic configs for already targeted containers, unless container discovery multi param is true
	if !load.Args.ContainerDiscoveryMulti {
		for _, id := range *foundTargetContainerIds {
			if id == container.ID {
				load.Logrus.Debug("discovery: container id " + container.ID + " already targeted")
				return false
			}
		}
	}

	switch discoveryConfig["tt"].(type) {
	case string:
		switch discoveryConfig["tt"].(string) {
		case load.TypeCname, load.TypeContainer:
			for _, containerName := range container.Names {
				checkContainerName := strings.TrimPrefix(containerName, "/") // docker adds a / in front
				if checkContainerName != "" && formatter.KvFinder(discoveryConfig["tm"].(string), checkContainerName, discoveryConfig["t"].(string)) {
					*(foundTargetContainerIds) = append(*(foundTargetContainerIds), container.ID)
					return true
				}
				// kubernetes container name fallback via label
				for key, val := range container.Labels {
					if val != "" && key == "io.kubernetes.container.name" {
						if formatter.KvFinder(discoveryConfig["tm"].(string), val, discoveryConfig["t"].(string)) {
							*(foundTargetContainerIds) = append(*(foundTargetContainerIds), container.ID)
							return true
						}
					}
				}
			}
		case load.Img, load.Image:
			if formatter.KvFinder(discoveryConfig["tm"].(string), container.Image, discoveryConfig["t"].(string)) {
				*(foundTargetContainerIds) = append(*(foundTargetContainerIds), container.ID)
				return true
			}
		}
	case nil:
		load.Logrus.Debug("discovery: targetType not set id: " + container.ID)
	}

	return false
}

func lowLevelIpv4Fetch(discoveryIPAddress *string, pid int, containerID string) []string {
	if *discoveryIPAddress == "" {
		// targetContainerInspect.State.Pid
		// cat /host/proc/<pid>/net/fib_trie | awk '/32 host/ { print f } {f=$2}' | grep -v 127.0.0.1 | sort -u
		load.Logrus.Debug(fmt.Sprintf("discovery: %v attempting low level ip fetch", containerID))
		target := "/host/proc"
		if load.ContainerID == "" {
			target = "/proc"
		}
		fibTrie, err := ioutil.ReadFile(fmt.Sprintf("%v/%v/net/fib_trie", target, pid))
		if err == nil && len(fibTrie) > 0 {
			reg := regexp.MustCompile(`\b((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\.)){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))(\r\n|\r|\n)\s+\/32 host\b`)
			matches := reg.FindAllStringSubmatch(strings.Replace(string(fibTrie), "127.0.0.1", "", -1), -1)
			ipMatches := []string{}
			for _, ipmatch := range matches {
				if len(ipmatch) >= 2 && !sliceContains(ipMatches, ipmatch[1]) {
					ipMatches = append(ipMatches, ipmatch[1])
				}
			}
			if len(ipMatches) > 0 {
				load.Logrus.Debug(fmt.Sprintf("discovery: %v container ip addresses found %v", containerID, ipMatches))
				*discoveryIPAddress = ipMatches[0]
				return ipMatches
			}
			load.Logrus.Error(fmt.Sprintf("discovery: %v container low level ip fetch failed", containerID))
		} else {
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("discovery: unable to read fib_trie")
		}
	}
	return nil
}

func portFallback(container types.Container, containerInspect types.ContainerJSON, ymlString *string, discoveryPort *string) {
	// kubernetes port fallback
	for key, val := range container.Labels {
		if key == "annotation.io.kubernetes.container.ports" {
			var x []interface{}
			err := json.Unmarshal([]byte(val), &x)
			if err == nil {
				for _, kubePort := range x {
					if kubePort.(map[string]interface{})["containerPort"] != nil {
						*discoveryPort = fmt.Sprintf("%v", kubePort.(map[string]interface{})["containerPort"])
						*ymlString = strings.Replace(*ymlString, "${auto:port}", *discoveryPort, -1)
						break
					}
				}
			}
		}
	}
	// secondary inspect fallback
	if *discoveryPort == "" || *discoveryPort == "0" {
		if containerInspect.Config != nil {
			for port := range containerInspect.Config.ExposedPorts {
				*discoveryPort = strings.Split(port.Port(), "/")[0]
				*ymlString = strings.Replace(*ymlString, "${auto:port}", *discoveryPort, -1)
				break
			}
		}
	}
}

func execHostnameFallback(discoveryIPAddress *string, containerID string, ymlString *string) {
	// fall back if IP is not discovered
	// attempt to directly fetch IP from container
	load.Logrus.Debug(fmt.Sprintf("discovery: %v attempting hostname -i fallback", containerID))
	ip, err := ExecContainerCommand(containerID, []string{"hostname", "-i"})
	ipv4 := strings.TrimSpace(ip)
	re := regexp.MustCompile(`\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b`)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Debug(fmt.Sprintf("discovery: %v container secondary fetch container ip failed", containerID))
	} else if ip != "" && re.Match([]byte(ipv4)) && !strings.Contains(ip, "exec failed") {
		*discoveryIPAddress = ipv4
		*ymlString = strings.Replace(*ymlString, "${auto:host}", ipv4, -1)
		*ymlString = strings.Replace(*ymlString, "${auto:ip}", ipv4, -1)
	}
}

func fallbackFindFlexContainerID(containers *[]types.Container) {
	// fallback on looking for image name "nri-flex" if flex's container id was not found internally
	load.Logrus.Debug("discovery: flex container id has not been found internally, failing back... checking for flex image or container name")

	var wg sync.WaitGroup
	wg.Add(len(*containers))
	for _, container := range *containers {
		go func(container types.Container) {
			defer wg.Done()
			if strings.Contains(container.Image, load.IntegrationName) && load.ContainerID == "" {
				load.ContainerID = container.ID
			}

			// fallback - check standard container names
			if load.ContainerID == "" {
				for _, containerName := range container.Names {
					checkContainerName := strings.TrimPrefix(containerName, "/") // docker adds a / in front
					if strings.Contains(checkContainerName, load.IntegrationNameShort) {
						load.ContainerID = container.ID
					}
				}
			}

			// fallback - check kubernetes container name via label
			if load.ContainerID == "" {
				for key, val := range container.Labels {
					if key == "io.kubernetes.container.name" {
						if strings.Contains(val, load.IntegrationNameShort) {
							load.ContainerID = container.ID
						}
					}
				}
			}

		}(container)
	}
	wg.Wait()

	if load.ContainerID == "" {
		load.Logrus.Debug("discovery: unable to find flex container id")
	} else {
		load.Logrus.Debug(fmt.Sprintf("discovery: flex container id %v", load.ContainerID))
	}
}

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// MakeTimestamp struct
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

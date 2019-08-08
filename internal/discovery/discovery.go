package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

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
			logger.Flex("error", err, "unable to set docker client", false)
		} else {
			ctx := context.Background()
			containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{})
			if err != nil {
				logger.Flex("error", err, "unable to set perform container list", false)
			} else if len(containerList) > 0 {
				// List config files in containerDiscoveryPath directory
				containerDiscoveryPath := filepath.FromSlash(load.Args.ContainerDiscoveryDir)
				containerDiscoveryFiles, err := ioutil.ReadDir(containerDiscoveryPath)
				logger.Flex("error", err, "failed to read config dir: "+load.Args.ContainerDiscoveryDir, false)

				CreateDynamicContainerConfigs(containerList, containerDiscoveryFiles, containerDiscoveryPath, configs)
				if len(containerDiscoveryFiles) == 0 {
					logger.Flex("debug", nil, "no configs found: "+load.Args.ContainerDiscoveryDir, false)
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
		logger.Flex("debug", nil, fmt.Sprintf("cpuset: %v", strCpuset), false)
		if strings.Contains(strCpuset, "kube") {
			load.IsKubernetes = true
		}
		values := strings.Split(strCpuset, "/")
		if len(values) > 0 {
			if len(values[len(values)-1]) == 64 {
				load.ContainerID = values[len(values)-1]
				logger.Flex("debug", fmt.Errorf("flex container id: %v", load.ContainerID), "", false)
			}
		}
	} else {
		logger.Flex("debug", nil, "potentially not running within a container", false)
	}
}

// CreateDynamicContainerConfigs Creates dynamic configs for each container
func CreateDynamicContainerConfigs(containers []types.Container, files []os.FileInfo, path string, ymls *[]load.Config) {
	var containerYmls []load.Config
	config.LoadFiles(&containerYmls, files, path)
	foundTargetContainerIds := []string{}

	// store inspected containers, so we do not re-inspect anything unnecessarily
	inspectedContainers := []types.ContainerJSON{}
	logger.Flex("debug", fmt.Errorf("containers %d, containerDiscoveryConfigs %d", len(containers), len(containerYmls)), "", false)
	// find flex container id if not set
	if load.ContainerID == "" {
		fallbackFindFlexContainerID(&containers)
	}

	// flex config file container_discovery parameter -> container
	runConfigLookup(cli, &containers, &inspectedContainers, &foundTargetContainerIds, ymls)
	// flex envs/labels -> container
	runReverseLookup(cli, &containers, &inspectedContainers, &foundTargetContainerIds, &containerYmls, ymls, path)
	// container envs/labels -> flex
	runForwardLookup(cli, &containers, &inspectedContainers, &foundTargetContainerIds, &containerYmls, ymls, path)
}

// runForwardLookup container envs -> flex
func runForwardLookup(dockerClient *client.Client, containers *[]types.Container, inspectedContainers *[]types.ContainerJSON, foundTargetContainerIds *[]string, containerYmls *[]load.Config, ymls *[]load.Config, path string) {

	var wg sync.WaitGroup
	wg.Add(len(*containers))
	for _, container := range *containers {
		go func(container types.Container) {
			defer wg.Done()

			// do not target already targeted containers
			targeted := false
			for _, foundTargetContainerID := range *foundTargetContainerIds {
				if foundTargetContainerID == container.ID {
					targeted = true
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
						logger.Flex("debug", nil, "container inspect failed", false)
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
							logger.Flex("debug", fmt.Errorf("fwd lookup for %v", key), "", false)
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
				// do not target the flex container
				if container.ID != load.ContainerID {
					// create discoveryConfigs - look for flex label and split
					for _, cd := range discoveryLoop {
						discoveryConfig := map[string]interface{}{}
						discoveryConfig["t"] = cd.Target
						discoveryConfig["c"] = cd.FileName
						discoveryConfig["tt"] = cd.Type
						discoveryConfig["tm"] = cd.Mode

						ctx := context.Background()
						reverseContainerInspect, err := dockerClient.ContainerInspect(ctx, container.ID)
						if err != nil {
							logger.Flex("error", fmt.Errorf("cfg container inspect failed on cid:%v - %v", container.ID, cd.FileName), "", false)
						} else {
							if findContainerTarget(discoveryConfig, container, foundTargetContainerIds) {
								logger.Flex("debug", fmt.Errorf("cfg lookup matched %v: - %v", container.ID, cd.FileName), "", false)
								*inspectedContainers = append(*inspectedContainers, reverseContainerInspect)
								addDynamicConfig(ymls, discoveryConfig, ymls, container, reverseContainerInspect, "")
							}
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
			logger.Flex("debug", nil, "container inspect failed", false)
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
				// do not target the flex container
				if container.ID != load.ContainerID {
					// create discoveryConfigs - look for flex label and split
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
								logger.Flex("error", fmt.Errorf("rev container inspect failed on cid:%v key:%v val:%v", container.ID, key, val), "", false)
							} else {
								if findContainerTarget(discoveryConfig, container, foundTargetContainerIds) {
									logger.Flex("debug", fmt.Errorf("rev lookup matched %v: %v - %v", container.ID, key, val), "", false)
									*inspectedContainers = append(*inspectedContainers, reverseContainerInspect)
									addDynamicConfig(containerYmls, discoveryConfig, ymls, container, reverseContainerInspect, path)
								}
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
			logger.Flex("error", fmt.Errorf("container discovery config file error %v", (discoveryConfig["c"])), "", false)
		}
		if containerYml.FileName == configName {
			logger.Flex("debug", fmt.Errorf("container discovery %v matched %v", targetContainer.ID, containerYml.FileName), "", false)
			if path == "" {
				path = containerYml.FilePath
			}
			b, err := ioutil.ReadFile(path + containerYml.FileName)
			if err != nil {
				logger.Flex("error", err, "unable to read flex config: "+path+containerYml.FileName, false)
			} else {
				ymlString := string(b)
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

				// attempt low level ip fetch
				lowLevelIpv4Fetch(&discoveryIPAddress, targetContainerInspect.State.Pid)
				// attempt hostname fallback
				execHostnameFallback(&discoveryIPAddress, targetContainer.ID)

				if discoveryIPAddress != "" {
					// substitute ip into yml
					ymlString = strings.Replace(ymlString, "${auto:host}", discoveryIPAddress, -1)
					ymlString = strings.Replace(ymlString, "${auto:ip}", discoveryIPAddress, -1)
				}

				if discoveryPort != "" && discoveryPort != "0" {
					// substitute port into yml
					ymlString = strings.Replace(ymlString, "${auto:port}", discoveryPort, -1)
				} else {
					// kubernetes port fallback
					for key, val := range targetContainer.Labels {
						if key == "annotation.io.kubernetes.container.ports" {
							var x []interface{}
							err := json.Unmarshal([]byte(val), &x)
							if err == nil {
								for _, kubePort := range x {
									if kubePort.(map[string]interface{})["containerPort"] != nil {
										discoveryPort = fmt.Sprintf("%v", kubePort.(map[string]interface{})["containerPort"])
										ymlString = strings.Replace(ymlString, "${auto:port}", discoveryPort, -1)
										break
									}
								}
							}
						}
					}
					// secondary inspect fallback
					if discoveryPort == "" || discoveryPort == "0" {
						if targetContainerInspect.Config != nil {
							for port := range targetContainerInspect.Config.ExposedPorts {
								discoveryPort = strings.Split(port.Port(), "/")[0]
								ymlString = strings.Replace(ymlString, "${auto:port}", discoveryPort, -1)
								break
							}
						}
					}
				}

				logger.Flex("debug", nil, fmt.Sprintf("target: %v %v - %v - %v:%v", targetContainer.ID, containerYml.FileName, ipMode, discoveryIPAddress, discoveryPort), false)

				if strings.Contains(ymlString, "${auto:host}") || strings.Contains(ymlString, "${auto:ip}") || strings.Contains(ymlString, "${auto:port}") {
					logger.Flex("debug", nil, "couldn't build dynamic cfg for: "+targetContainer.Image+" - "+targetContainer.ID, false)
					logger.Flex("debug", nil, "missing variable unable to create dynamic cfg ip:<"+discoveryIPAddress+">-port:<"+discoveryPort+">", false)
				} else {
					yml, err := config.ReadYML(ymlString)
					if err != nil {
						logger.Flex("error", err, "unable to unmarshal yml config: "+path+containerYml.FileName, false)
						logger.Flex("error", fmt.Errorf(ymlString), "", false)
					} else {
						if yml.CustomAttributes == nil {
							yml.CustomAttributes = map[string]string{}
						}
						for key, val := range targetContainer.Labels {
							yml.CustomAttributes[key] = val
						}
						yml.CustomAttributes["containerId"] = targetContainer.ID
						yml.CustomAttributes["imageId"] = targetContainer.Image
						yml.CustomAttributes["IDShort"] = targetContainer.ID[0:12]
						*ymls = append(*ymls, yml)
					}
				}

			}

		} else {
			logger.Flex("debug", fmt.Errorf("container discovery %v: containerFileName %v did not match configName %v", targetContainer.ID, containerYml.FileName, configName), "", false)
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
	// do not add any dynamic configs for already targeted containers
	for _, id := range *foundTargetContainerIds {
		if id == container.ID {
			return false
		}
	}

	switch discoveryConfig["tt"].(type) {
	case string:
		switch discoveryConfig["tt"].(string) {
		case "cname", load.TypeContainer:
			for _, containerName := range container.Names {
				checkContainerName := strings.TrimPrefix(containerName, "/") // docker adds a / in front
				if formatter.KvFinder(discoveryConfig["tm"].(string), checkContainerName, discoveryConfig["t"].(string)) {
					*(foundTargetContainerIds) = append(*(foundTargetContainerIds), container.ID)
					return true
				}
				// kubernetes container name fallback via label
				for key, val := range container.Labels {
					if key == "io.kubernetes.container.name" {
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
		logger.Flex("debug", nil, "targetType not set id: "+container.ID, false)
	}

	return false
}

func lowLevelIpv4Fetch(discoveryIPAddress *string, pid int) {
	if *discoveryIPAddress == "" {
		// targetContainerInspect.State.Pid
		// cat /host/proc/<pid>/net/fib_trie | awk '/32 host/ { print f } {f=$2}' | grep -v 127.0.0.1 | sort -u

		logger.Flex("debug", nil, "attempting low level ip fetch", false)

		// Create a new context and add a timeout to it
		ctx, cancel := context.WithTimeout(context.Background(), load.DefaultTimeout)
		defer cancel() // The cancel should be deferred so resources are cleaned up

		target := "/host/proc/"
		if load.ContainerID == "" {
			target = "/proc/"
		}

		// Create the command with our context
		cmd := exec.CommandContext(ctx, "/bin/sh", "-c", fmt.Sprintf("cat %v/%v", target, pid)+
			`/net/fib_trie | awk '/32 host/ { print f } {f=$2}' | grep -v 127.0.0.1 | sort -u`)
		output, err := cmd.CombinedOutput()

		if err != nil {
			message := "command failed: " + ""
			if output != nil {
				message = message + " " + string(output)
			}
			logger.Flex("error", err, message, false)
		} else if ctx.Err() == context.DeadlineExceeded {
			logger.Flex("error", ctx.Err(), "command timed out", false)
		} else if ctx.Err() != nil {
			logger.Flex("debug", err, "command execution failed", false)
		} else {
			ipv4 := strings.TrimSpace(string(output))
			// ensure this is an ipv4 address
			re := regexp.MustCompile(`\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b`)
			if re.Match([]byte(ipv4)) {
				logger.Flex("debug", nil, fmt.Sprintf("fetched %v", ipv4), false)
				*discoveryIPAddress = ipv4
			} else {
				logger.Flex("debug", fmt.Errorf("low level fetch failed %v", ipv4), "", false)
			}
		}
	}
}

func execHostnameFallback(discoveryIPAddress *string, containerID string) {
	if *discoveryIPAddress == "" {
		// fall back if IP is not discovered
		// attempt to directly fetch IP from container
		ip, err := ExecContainerCommand(containerID, []string{"hostname", "-i"})
		ipv4 := strings.TrimSpace(ip)
		re := regexp.MustCompile(`\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b`)
		if err != nil {
			logger.Flex("debug", err, "secondary fetch container ip failed", false)
		} else if ip != "" && re.Match([]byte(ipv4)) && !strings.Contains(ip, "exec failed") {
			*discoveryIPAddress = ipv4
		}
	}
}

func fallbackFindFlexContainerID(containers *[]types.Container) {
	// fallback on looking for image name "nri-flex" if flex's container id was not found internally
	logger.Flex("debug", fmt.Errorf("flex container id has not been found internally"), "", false)
	logger.Flex("debug", fmt.Errorf("falling back - looking for 'nri-flex' image or container name"), "", false)

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
		logger.Flex("debug", fmt.Errorf("unable to find flex container id"), "", false)
	} else {
		logger.Flex("debug", fmt.Errorf("flex container id: %v", load.ContainerID), "", false)
	}
}

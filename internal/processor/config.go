package processor

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	yaml "gopkg.in/yaml.v2"
)

// LoadConfigFiles loads config files
func LoadConfigFiles(ymls *[]load.Config, files []os.FileInfo, path string) {
	for _, f := range files {
		b, err := ioutil.ReadFile(path + f.Name())
		if err != nil {
			logger.Flex("debug", err, "unable to readfile", false)
			continue
		}
		if !strings.Contains(f.Name(), "yml") && !strings.Contains(f.Name(), "yaml") {
			continue
		}
		ymlStr := string(b)
		SubEnvVariables(&ymlStr)
		SubTimestamps(&ymlStr)
		yml, err := ReadYML(ymlStr)
		yml.FileName = f.Name()
		if err != nil {
			logger.Flex("debug", err, "unable to read yml", false)
			continue
		}
		if yml.Name == "" {
			logger.Flex("debug", err, "please set a name on your config file", false)
			// fmt.Println("Please set a name on your config file", f.Name())
			continue
		}
		*ymls = append(*ymls, yml)
	}
}

// ReadYML Unmarshals yml files
func ReadYML(yml string) (load.Config, error) {
	c := load.Config{}
	err := yaml.Unmarshal([]byte(yml), &c)
	if err != nil {
		return load.Config{}, err
	}
	return c, nil
}

// RunConfig Action each config file
func RunConfig(yml load.Config) {
	samplesToMerge := map[string][]interface{}{}
	for i := range yml.APIs {
		runVariableProcessor(i, &yml)
		dataSets := fetchData(i, &yml)
		runDataHandler(dataSets, &samplesToMerge, i, &yml)
	}
	ProcessSamplesToMerge(&samplesToMerge, &yml)
}

// runVariableProcessor substitute store variables into specific parts of config files
func runVariableProcessor(i int, cfg *load.Config) {
	// don't use variable processor if nothing exists in variable store
	if len((*cfg).VariableStore) > 0 {
		// to simplify replacement, convert to string, and convert back later
		tmpCfgBytes, err := yaml.Marshal(&cfg)
		if err != nil {
			logger.Flex("debug", err, "variable processor marshal failed", false)
		} else {
			tmpCfgStr := string(tmpCfgBytes)
			variableReplaces := regexp.MustCompile(`\${var:.*?}`).FindAllString(tmpCfgStr, -1)
			replaceOccured := false
			for _, variableReplace := range variableReplaces {
				variableKey := strings.TrimSuffix(strings.Split(variableReplace, "${var:")[1], "}") // eg. "channel"
				if cfg.VariableStore[variableKey] != "" {
					tmpCfgStr = strings.Replace(tmpCfgStr, variableReplace, cfg.VariableStore[variableKey], -1)
					replaceOccured = true
				}
			}
			// if replace occurred convert string to config yaml and reload
			if replaceOccured {
				newCfg, err := ReadYML(tmpCfgStr)
				if err != nil {
					logger.Flex("debug", err, "variable processor unmarshal failed", false)
				} else {
					*cfg = newCfg
				}
			}
		}
	}
}

// runLookupProcessor
func runLookupProcessor(str string, cfg *load.Config, i int) {
	lookupReplaces := regexp.MustCompile(`\${lookup:.*?}`).FindAllString(str, -1)

	newConfig := load.Config{
		Name:   cfg.Name,
		Global: cfg.Global,
	}

	for _, lookupReplace := range lookupReplaces {
		// eg. lookupReplace == ${lookup:channels}
		lookupKey := strings.TrimSuffix(strings.Split(lookupReplace, "${lookup:")[1], "}") // eg. "channels"
		if cfg.LookupStore[lookupKey] != nil {
			for _, storedKey := range cfg.LookupStore[lookupKey] { // eg. list of channels
				// add into newConfig>API, and execute processConfig again
				newURL := strings.Replace(str, lookupReplace, url.QueryEscape(storedKey), -1)
				newAPI := load.API{
					Name:      cfg.APIs[i].EventType,
					URL:       newURL,
					EventType: cfg.APIs[i].EventType,
				}
				newConfig.APIs = append(newConfig.APIs, newAPI)
			}
		}
	}

	// re issue process config with newly built config
	RunConfig(newConfig)
}

// RunConfigFiles Processes yml files
func RunConfigFiles(ymls *[]load.Config) {
	var wg sync.WaitGroup
	wg.Add(len(*ymls))
	for _, yml := range *ymls {
		go func(yml load.Config) {
			defer wg.Done()
			RunConfig(yml)
			load.ConfigsProcessed++
		}(yml)
	}
	wg.Wait()
}

// SubTimestamps substitute timestamps into config
func SubTimestamps(strConf *string) {
	current := time.Now()
	currentNano := current.UnixNano()
	currentMs := currentNano / 1e+6
	currentSec := current.Unix()
	*strConf = strings.Replace(*strConf, "${timestamp:ms}", fmt.Sprint(currentMs), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:ns}", fmt.Sprint(currentNano), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:s}", fmt.Sprint(currentSec), -1)

	timestamps := regexp.MustCompile(`\${timestamp:.*?}`).FindAllString(*strConf, -1)
	for _, timestamp := range timestamps {
		newTimestamp := int64(0)
		matches := formatter.RegMatch(timestamp, `(\${timestamp:)(ms|ns|s)(-|\+)(\d*)`)
		if len(matches) == 4 {
			switch matches[1] {
			case "ms":
				newTimestamp = currentMs
			case "ns":
				newTimestamp = currentNano
			case "s":
				newTimestamp = currentSec
			default:
				break
			}
			value, err := strconv.ParseInt(matches[3], 10, 64)
			if err != nil {
				logger.Flex("debug", err, "failed to parse int", false)
			} else {
				switch matches[2] {
				case "+":
					newTimestamp += value
				case "-":
					newTimestamp -= value
				default:
					break
				}
				*strConf = strings.Replace(*strConf, timestamp, fmt.Sprint(newTimestamp), -1)
			}
		}
	}
}

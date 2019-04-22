package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

// SubEnvVariables substitutes environment variables into config
// Use a double dollar sign eg. $$MY_ENV_VAR to subsitute that environment variable into the config file
// Can be useful with kubernetes service environment variables
func SubEnvVariables(strConf *string) {
	subCount := strings.Count(*strConf, "$$")
	replaceCount := 0
	if subCount > 0 {
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			if len(pair) == 2 && pair[0] != "" {
				if strings.Contains(*strConf, "$$"+pair[0]) {
					*strConf = strings.Replace(*strConf, "$$"+pair[0], pair[1], -1)
					load.StatusCounterIncrement("environmentVariablesReplaced")
					replaceCount++
				}
			}
			if replaceCount >= subCount {
				break
			}
		}
	}
}

// SubTimestamps substitute timestamps into config
// Substitution keys:
// ${timestamp:ms} - timestamp in milliseconds
// ${timestamp:ns} - timestamp in nanoseconds
// ${timestamp:s} - timestamp in seconds
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
				load.StatusCounterIncrement("timestampsReplaced")
			}
		}
	}
}

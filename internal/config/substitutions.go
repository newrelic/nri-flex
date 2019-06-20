package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	yaml "gopkg.in/yaml.v2"
)

// SubLookupFileData substitutes data from lookup files into config
func SubLookupFileData(configs *[]load.Config, config load.Config) {
	logger.Flex("debug", nil, "running lookup files", false)

	tmpCfgBytes, err := yaml.Marshal(&config)
	if err != nil {
		logger.Flex("error", err, "sub lookup file data marshal failed", false)
	} else {

		b, err := ioutil.ReadFile(config.LookupFile)
		if err != nil {
			logger.Flex("error", err, "unable to readfile", false)
			return
		}

		jsonOut := []interface{}{}
		jsonErr := json.Unmarshal(b, &jsonOut)
		if jsonErr != nil {
			logger.Flex("error", jsonErr, config.LookupFile, false)
			return
		}

		// create a new config file per
		for _, item := range jsonOut {
			switch obj := item.(type) {
			case map[string]interface{}:
				tmpCfgStr := string(tmpCfgBytes)
				variableReplaces := regexp.MustCompile(`\${lf:.*?}`).FindAllString(tmpCfgStr, -1)
				replaceOccured := false
				for _, variableReplace := range variableReplaces {
					variableKey := strings.TrimSuffix(strings.Split(variableReplace, "${lf:")[1], "}") // eg. "channel"
					if obj[variableKey] != nil {
						tmpCfgStr = strings.Replace(tmpCfgStr, variableReplace, fmt.Sprintf("%v", obj[variableKey]), -1)
						replaceOccured = true
					}
				}
				// if replace occurred convert string to config yaml and reload
				if replaceOccured {
					newCfg, err := ReadYML(tmpCfgStr)
					if err != nil {
						logger.Flex("error", err, fmt.Sprintf("new lookup file unmarshal failed %v %v", config.Name, config.LookupFile), false)
						logger.Flex("error", fmt.Errorf("check for errors or run yaml lint against the below output:\n%v", tmpCfgStr), "", false)
					} else {
						*configs = append(*configs, newCfg)
					}
				}
			default:
				logger.Flex("debug", nil, "lookup file needs to contain an array of objects", false)
			}
		}
	}
}

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
// supported format
// ${timestamp:[ms|ns|s|date|datetime|datetimetz|dateutc|datetimeutc|datetimeutctz][+|-][Number][ms|milli|millisecond|ns|nano|nanosecond|s|sec|second|m|min|minute|h|hr|hour]}

// Substitution keys:
// ${timestamp:ms} - timestamp in milliseconds
// ${timestamp:ns} - timestamp in nanoseconds
// ${timestamp:s} - timestamp in seconds
// ${timestamp:date} - date in date format local timezone: 2006-01-02
// ${timestamp:datetime} - datetime in date and time format local timezone : 2006-01-02T03:04
// ${timestamp:datetimetz} - datetime in date and time format local timezone : 2006-01-02T15:04:05Z07:00
// ${timestamp:dateutc} - date in date format utc timezone: 2006-01-02
// ${timestamp:datetimeutc} - datetime in date and time format utc timezone: 2006-01-02T03:04
// ${timestamp:datetimeutctz} - datetime in date and time format utc timezone: 2006-01-02T15:04:05Z07:00

// SubTimestamps - return timestamp/date/datetime of current date/time with optional adjustment in various format
func SubTimestamps(strConf *string) {

	current := time.Now()
	currentUTC := time.Now().UTC()

	// date and datetime output format
	dateFormat := "2006-01-02"
	datetimeFormat := "2006-01-02T15:04:05"
	datetimeFormatTZ := "2006-01-02T15:04:05Z07:00"

	*strConf = strings.Replace(*strConf, "${timestamp:ms}", fmt.Sprint(current.UnixNano()/1e+6), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:ns}", fmt.Sprint(current.UnixNano()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:s}", fmt.Sprint(current.Unix()), -1)

	*strConf = strings.Replace(*strConf, "${timestamp:date}", fmt.Sprint(current.Format(dateFormat)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:datetime}", fmt.Sprint(current.Format(datetimeFormat)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:datetimetz}", fmt.Sprint(current.Format(datetimeFormatTZ)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:dateutc}", fmt.Sprint(currentUTC.Format(dateFormat)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:datetimeutc}", fmt.Sprint(currentUTC.Format(datetimeFormat)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:datetimeutctz}", fmt.Sprint(currentUTC.Format(datetimeFormatTZ)), -1)

	timestamps := regexp.MustCompile(`\${timestamp:.*?}`).FindAllString(*strConf, -1)
	for _, timestamp := range timestamps {

		durationType := time.Millisecond
		timestampCurrent := current
		timestampUTC := currentUTC
		timestampReturn := ""
		defaultTimestamp := "${timestamp:ms}"
		var err error

		matches := formatter.RegMatch(timestamp, `(\${timestamp:)(ms|ns|s|date|datetime|datetimetz|dateutc|datetimeutc|datetimeutctz)(-|\+)(\d+|\d+\D+)\}`)
		// matches patterns like {timestamp:ms+10} or {timestamp:ns-10s}, {timestamp:ns-[Digits&NonDigits]},etc
		if len(matches) == 4 {
			var duration int64

			matchDuration := formatter.RegMatch(matches[3], `(\d+)(\D+)`)
			if len(matchDuration) == 2 {
				//match case like {timestamp:ns-10s}
				duration, _ = strconv.ParseInt(matchDuration[0], 10, 64)

				switch strings.ToLower(matchDuration[1]) {
				case "ns", "nano", "nanosecond":
					durationType = time.Nanosecond
				case "ms", "milli", "millisecond":
					durationType = time.Millisecond
				case "s", "sec", "second":
					durationType = time.Second
				case "m", "min", "minute":
					durationType = time.Minute
				case "h", "hr", "hour":
					durationType = time.Hour
				default:
					logger.Flex("info", err, "unable to parse "+timestamp+", defaulting to "+defaultTimestamp, false)
				}

			} else {
				// match case like {timestamp:ns-10}, only digits are provided, use default durationType := time.Millisecond
				duration, _ = strconv.ParseInt(matches[3], 10, 64)
			}

			switch matches[2] {
			case "+":
			case "-":
				duration = -duration
			}

			// adjust the timestamp offset based on duration and durationType
			timestampCurrent = timestampCurrent.Add(time.Duration(duration) * durationType)
			timestampUTC = timestampUTC.Add(time.Duration(duration) * durationType)

			// prepare the timestamp return format
			switch matches[1] {
			case "ms":
				timestampReturn = fmt.Sprint(timestampCurrent.UnixNano() / 1e+6)
			case "ns":
				timestampReturn = fmt.Sprint(timestampCurrent.UnixNano())
			case "s":
				timestampReturn = fmt.Sprint(timestampCurrent.Unix())
			case "date":
				timestampReturn = fmt.Sprint(timestampCurrent.Format(dateFormat))
			case "datetime":
				timestampReturn = fmt.Sprint(timestampCurrent.Format(datetimeFormat))
			case "datetimetz":
				timestampReturn = fmt.Sprint(timestampCurrent.Format(datetimeFormatTZ))
			case "dateutc":
				timestampReturn = fmt.Sprint(timestampUTC.Format(dateFormat))
			case "datetimeutc":
				timestampReturn = fmt.Sprint(timestampUTC.Format(datetimeFormat))
			case "datetimeutctz":
				timestampReturn = fmt.Sprint(timestampUTC.Format(datetimeFormatTZ))

			default:
				// default to timestamp in unix milliseoncds
				logger.Flex("info", err, "unable to parse "+timestamp+", defaulting to "+defaultTimestamp, false)
				timestampReturn = fmt.Sprint(timestampCurrent.UnixNano() / 1e+6)
			}

		} else {

			// if the regex does not match,  default to the current timestamp in unix milliseoncds
			logger.Flex("info", err, "unable to parse "+timestamp+", defaulting to "+defaultTimestamp, false)
			timestampReturn = fmt.Sprint(timestampCurrent.UnixNano() / 1e+6)

		}
		*strConf = strings.Replace(*strConf, timestamp, timestampReturn, -1)
		load.StatusCounterIncrement("timestampsReplaced")

	}

}

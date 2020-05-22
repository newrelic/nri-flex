/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

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
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// SubLookupFileData substitutes data from lookup files into config
func SubLookupFileData(configs *[]load.Config, config load.Config) error {
	load.Logrus.WithFields(logrus.Fields{
		"name": config.Name,
	}).Debug("config: running lookup files")

	tmpCfgBytes, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("config name: %s: sub lookup file data marshal failed, error: %v", config.Name, err)
	}

	b, err := ioutil.ReadFile(config.LookupFile)
	if err != nil {
		return fmt.Errorf("config name: %s: failed to read lookup file, error: %v", config.LookupFile, err)
	}

	var jsonOut []interface{}
	jsonErr := json.Unmarshal(b, &jsonOut)
	if jsonErr != nil {
		return fmt.Errorf("config name: %s: failed to unmarshal lookup file, error: %v", config.LookupFile, err)
	}

	// create a new config file per
	for _, item := range jsonOut {
		tmpCfgStr := string(tmpCfgBytes)
		newCfg, err := fillTemplateConfigWithValues(item, tmpCfgStr)

		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"file":       config.LookupFile,
				"name":       config.Name,
				"suggestion": "check for errors or run yaml lint against the below output",
			}).WithError(err).Error("config: new lookup file unmarshal failed")
			load.Logrus.Error(tmpCfgStr)
		}

		if newCfg != nil {
			*configs = append(*configs, *newCfg)
		}
	}
	return nil
}

func fillTemplateConfigWithValues(values interface{}, configTemplate string) (*load.Config, error) {
	switch obj := values.(type) {
	case map[string]interface{}:
		variableReplaces := regexp.MustCompile(`\${lf:.*?}`).FindAllString(configTemplate, -1)
		replaceOccurred := false
		for _, variableReplace := range variableReplaces {
			variableKey := strings.TrimSuffix(strings.Split(variableReplace, "${lf:")[1], "}") // eg. "channel"
			if obj[variableKey] != nil {
				value := obj[variableKey]
				configTemplate = strings.Replace(
					configTemplate,
					variableReplace,
					toString(value),
					-1)
				replaceOccurred = true
			}
		}
		// if replace occurred convert string to values yaml and reload
		if replaceOccurred {
			newCfg, err := ReadYML(configTemplate)
			if err != nil {
				return nil, err
			}
			return &newCfg, nil
		}
	default:
		load.Logrus.Debug("config: lookup file needs to contain an array of objects")
	}
	return nil, nil
}

func toString(value interface{}) string {
	format := "%v"
	switch value.(type) {
	case int:
		format = "%d"
	case float64, float32:
		format = "%f"
	case string:
		format = "%s"
	}

	return fmt.Sprintf(format, value)
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
func SubTimestamps(strConf *string, currentTime time.Time) {
	currentUTC := currentTime.UTC()

	// date and datetime output format
	dateFormat := "2006-01-02"
	datetimeFormat := "2006-01-02T15:04:05"
	datetimeFormatTZ := "2006-01-02T15:04:05Z07:00"
	*strConf = strings.Replace(*strConf, "${timestamp:ms}", fmt.Sprint(currentTime.UnixNano()/1e+6), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:ns}", fmt.Sprint(currentTime.UnixNano()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:s}", fmt.Sprint(currentTime.Unix()), -1)

	*strConf = strings.Replace(*strConf, "${timestamp:date}", fmt.Sprint(currentTime.Format(dateFormat)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:datetime}", fmt.Sprint(currentTime.Format(datetimeFormat)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:datetimetz}", fmt.Sprint(currentTime.Format(datetimeFormatTZ)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:dateutc}", fmt.Sprint(currentUTC.Format(dateFormat)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:datetimeutc}", fmt.Sprint(currentUTC.Format(datetimeFormat)), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:datetimeutctz}", fmt.Sprint(currentUTC.Format(datetimeFormatTZ)), -1)

	*strConf = strings.Replace(*strConf, "${timestamp:year}", strconv.Itoa(currentTime.Year()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:month}", strconv.Itoa(int(currentTime.Month())), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:day}", strconv.Itoa(currentTime.Day()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:hour}", strconv.Itoa(currentTime.Hour()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:minute}", strconv.Itoa(currentTime.Minute()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:second}", strconv.Itoa(currentTime.Second()), -1)

	*strConf = strings.Replace(*strConf, "${timestamp:utcyear}", strconv.Itoa(currentUTC.Year()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:utcmonth}", strconv.Itoa(int(currentUTC.Month())), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:utcday}", strconv.Itoa(currentUTC.Day()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:utchour}", strconv.Itoa(currentUTC.Hour()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:utcminute}", strconv.Itoa(currentUTC.Minute()), -1)
	*strConf = strings.Replace(*strConf, "${timestamp:utcsecond}", strconv.Itoa(currentUTC.Second()), -1)

	timestamps := regexp.MustCompile(`\${timestamp:.*?}`).FindAllString(*strConf, -1)
	for _, timestamp := range timestamps {

		durationType := time.Millisecond
		timestampCurrent := currentTime
		timestampUTC := currentUTC
		timestampReturn := ""
		defaultTimestamp := "${timestamp:ms}"
		var err error

		matches := formatter.RegMatch(timestamp, `(\${timestamp:)(ms|ns|s|date|datetime|datetimetz|dateutc|datetimeutc|datetimeutctz|year|month|day|hour|second|utcyear|utcmonth|utcday|utchour|minute|utcminute|utcsecond)(-|\+)(\d+|\d+\D+)\}`)

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
					load.Logrus.WithError(err).
						Error("config: unable to parse " + timestamp + ", defaulting to " + defaultTimestamp)
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

			case "year":
				timestampReturn = strconv.Itoa(timestampCurrent.Year())
			case "month":
				timestampReturn = strconv.Itoa(int(timestampCurrent.Month()))
			case "day":
				timestampReturn = strconv.Itoa(timestampCurrent.Day())
			case "hour":
				timestampReturn = strconv.Itoa(timestampCurrent.Hour())
			case "minute":
				timestampReturn = strconv.Itoa(timestampCurrent.Minute())
			case "second":
				timestampReturn = strconv.Itoa(timestampCurrent.Second())

			case "utcyear":
				timestampReturn = strconv.Itoa(timestampUTC.Year())
			case "utcmonth":
				timestampReturn = strconv.Itoa(int(timestampCurrent.Month()))
			case "utcday":
				timestampReturn = strconv.Itoa(timestampUTC.Day())
			case "utchour":
				timestampReturn = strconv.Itoa(timestampUTC.Hour())
			case "utcminute":
				timestampReturn = strconv.Itoa(timestampUTC.Minute())

			case "utcsecond":
				timestampReturn = strconv.Itoa(timestampUTC.Second())

			default:
				// default to timestamp in unix milliseconds
				load.Logrus.WithFields(logrus.Fields{
					"err": err,
				}).Debug("config: unable to parse " + timestamp + ", defaulting to " + defaultTimestamp)

				timestampReturn = fmt.Sprint(timestampCurrent.UnixNano() / 1e+6)
			}

		} else {

			// if the regex does not match,  default to the current timestamp in unix milliseoncds
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Debug("config: unable to parse " + timestamp + ", defaulting to " + defaultTimestamp)

			timestampReturn = fmt.Sprint(timestampCurrent.UnixNano() / 1e+6)

		}
		*strConf = strings.Replace(*strConf, timestamp, timestampReturn, -1)
		load.StatusCounterIncrement("timestampsReplaced")

	}

}

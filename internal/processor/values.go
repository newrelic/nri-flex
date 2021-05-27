/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
)

// RunValConversion performs percentage to decimal & nano second to millisecond
func RunValConversion(v *interface{}, api load.API, key *string) {
	if api.PercToDecimal {
		formatter.PercToDecimal(v)
	}

	value := fmt.Sprintf("%v", *v)
	if strings.Contains(value, "µs") {
		valueSplit := strings.Split(value, "µs")
		newValue, _ := strconv.ParseFloat(valueSplit[0], 64)
		newValue /= 1000 // convert to ms
		*v = newValue
		*key += ".ms"
	}
}

// RunSubParse splits nested values out from one line eg. db0:keys=1,expires=0,avg_ttl=0
func RunSubParse(subParse []load.Parse, currentSample *map[string]interface{}, key string, v interface{}) {
	for _, parse := range subParse {
		if len(parse.SplitBy) == 2 {
			process := formatter.KvFinder(parse.Type, key, parse.Key)
			if process {
				values := strings.Split(fmt.Sprintf("%v", v), parse.SplitBy[0])
				for _, val := range values {
					nestedVal := strings.Split(val, parse.SplitBy[1])
					if len(nestedVal) == 2 {
						(*currentSample)[key+"."+nestedVal[0]] = nestedVal[1]
					}
				}
			}
		}
	}
}

// RunValueParser use regex to find a key, and pluck out its value by regex
func RunValueParser(v *interface{}, api load.API, key *string) {
	for regexKey, regexVal := range api.ValueParser {
		if formatter.KvFinder("regex", *key, regexKey) {
			value := ""
			switch (*v).(type) {
			case float32, float64:
				// For float numbers, use decimal point format instead of scientific notation (e.g. 2026112.000000 vs 2.026112e+06 )
				// to allow the parser to process the original float number 2026112.000000 rather than 2.026112e+06
				value = fmt.Sprintf("%f", *v)
			default:
				value = fmt.Sprintf("%v", *v)
			}
			*v = formatter.ValueParse(value, regexVal)
		}
	}
}

// RunValueTransformer use regex to find a key, and then transform the value
// eg. key: world
// key: hello-${value} == key: hello-world
func RunValueTransformer(v *interface{}, api load.API, key *string) {
	for regexKey, newValue := range api.ValueTransformer {
		if formatter.KvFinder("regex", *key, regexKey) {
			currentValue := fmt.Sprintf("%v", *v)
			*v = strings.Replace(newValue, "${value}", currentValue, -1)
		}
	}
}

// RunPluckNumbers pluck numbers out automatically with ValueParser
// eg. "sample_start_time = 1552864614.137869 (Sun, 17 Mar 2019 23:16:54 GMT)"
// returns 1552864614.137869
func RunPluckNumbers(v *interface{}, api load.API, key *string) {
	if api.PluckNumbers {
		value := fmt.Sprintf("%v", *v)
		*v = formatter.ValueParse(value, `[+-]?([0-9]*\.?[0-9]+|[0-9]+\.?[0-9]*)([eE][+-]?[0-9]+)?`)
	}
}

// RunMathCalculations performs math calculations
func RunMathCalculations(math *map[string]string, currentSample *map[string]interface{}) {
	for newMetric, formula := range *math {
		finalFormula := formula
		keys := regexp.MustCompile(`\${.*?}`).FindAllString(finalFormula, -1)
		for _, key := range keys {
			findKey := strings.TrimSuffix(strings.TrimPrefix(key, "${"), "}")
			if (*currentSample)[findKey] != nil {
				finalFormula = strings.Replace(finalFormula, key, fmt.Sprintf("%v", (*currentSample)[findKey]), -1)
			}
		}
		expression, err := govaluate.NewEvaluableExpression(finalFormula)
		if err != nil {
			load.Logrus.WithError(err).Errorf("processor-values: %v math exp failed %v", newMetric, finalFormula)
		} else {
			result, err := expression.Evaluate(nil)
			if err != nil {
				load.Logrus.WithError(err).Errorf("processor-values: %v math evaluate failed %v", newMetric, finalFormula)
			} else {
				(*currentSample)[newMetric] = result
			}
		}
	}
}

// RunValueMapper map the value using regex grouping for keys e.g.  "*.?\s(Service Status)=>$1-Good" -> "Service Status-Good"
func RunValueMapper(mapKeys map[string][]string, currentSample *map[string]interface{}, key string, v *interface{}) {
	for mapKey, mapVal := range mapKeys {
		keySplit := strings.Split(mapKey, "=>")
		if key == keySplit[0] {
			replacedValue := false
			for _, mapEntry := range mapVal {
				valueSplit := strings.Split(mapEntry, "=>")
				if len(valueSplit) == 2 {
					regexPattern := valueSplit[0]
					targetValue := valueSplit[1]
					r := regexp.MustCompile(regexPattern)
					res := r.FindStringSubmatch(toString(v))
					for i, value := range res {
						if i != 0 {
							targetValue = strings.ReplaceAll(targetValue, "$"+strconv.Itoa(i), value)
							replacedValue = true
						}
					}
					if replacedValue {
						if len(keySplit) == 2 {
							(*currentSample)[keySplit[1]] = targetValue
						} else {
							*v = targetValue
						}
						break
					}
				}
			}
		}
	}
}

// RunTimestampConversion find keys with regex, convert date<=>timestamp
func RunTimestampConversion(v *interface{}, api load.API, key *string) {
	for regexKey, regexVal := range api.TimestampConversion {
		if formatter.KvFinder("regex", *key, regexKey) {
			value := toString(v)
			convertDateStamp(regexVal, &value)
			*v = value
		}
	}
}

// convert value to string, float64 to string without decimal
func toString(v *interface{}) string {
	switch val := (*v).(type) {
	case int:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%0.f", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// convert value to string, float64 to string without decimal
// func toString(value interface{}) string {
// 	format := "%v"
// 	switch value.(type) {
// 	case int:
// 		format = "%d"
// 	case float64, float32:
// 		format = "%0.f"
// 	case string:
// 		format = "%s"
// 	}

// 	return fmt.Sprintf(format, value)
// }

func convertDateStamp(timestampTamplate string, targetValue *string) {

	DATEFORMAT := map[string]string{
		"ANSIC":       "Mon Jan _2 15:04:05 2006",
		"UnixDate":    "Mon Jan _2 15:04:05 MST 2006",
		"RubyDate":    "Mon Jan 02 15:04:05 -0700 2006",
		"ATOM":        "2006-01-02T15:04:05Z07:00",
		"COOKIE":      "Monday, 02-Jan-06 15:04:05 MST",
		"ISO8601":     "2006-01-02T15:04:05Z0700",
		"RFC822":      "Mon, 02 Jan 06 15:04:05 Z0700",
		"RFC850":      "Monday, 02-Jan-06 15:04:05 MST",
		"RFC1036":     "Mon, 02 Jan 06 15:04:05 Z0700",
		"RFC1123":     "Mon, 02 Jan 2006 15:04:05 Z0700",
		"RFC2822":     "Mon, 02 Jan 2006 15:04:05 Z0700",
		"RFC3339":     "2006-01-02T15:04:05Z07:00",
		"RFC3339Nano": "2006-01-02T15:04:05.999999999Z07:00",
		"RSS":         "Mon, 02 Jan 2006 15:04:05 Z0700",
		"W3C":         "2006-01-02T15:04:05Z07:00",
	}

	timestampFormat := strings.Split(timestampTamplate, "::")
	if len(timestampFormat) != 0 {
		if timestampFormat[0] == "TIMESTAMP" {
			srcDateformat := time.RFC3339
			if len(timestampFormat) == 2 {
				if val, ok := DATEFORMAT[timestampFormat[1]]; ok {
					srcDateformat = val
				} else {
					srcDateformat = timestampFormat[1]
				}
			}
			resTime, err := time.Parse(srcDateformat, *targetValue)
			if err != nil {
				load.Logrus.WithError(err).Errorf("processor-values: %v TimestampConversion failed %v", srcDateformat, *targetValue)
			} else {
				*targetValue = strconv.FormatInt(resTime.Unix(), 10)
			}
		}
		if timestampFormat[0] == "DATE" {
			dstDateformat := time.RFC3339
			if len(timestampFormat) == 2 {
				if val, ok := DATEFORMAT[timestampFormat[1]]; ok {
					dstDateformat = val
				} else {
					dstDateformat = timestampFormat[1]
				}
			}
			n, err := strconv.ParseInt(*targetValue, 10, 64)
			if err != nil {
				load.Logrus.WithError(err).Errorf("processor-values: %v TimestampConversion failed %v", dstDateformat, *targetValue)
			} else {
				unixTimeUTC := time.Unix(n, 0)
				*targetValue = unixTimeUTC.Format(dstDateformat)
			}
		}
	}
}

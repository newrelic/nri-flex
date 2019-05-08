package processor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/jeremywohl/flatten"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
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
			value := fmt.Sprintf("%v", *v)
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

// RunLazyFlatten lazy flattens the payload
func RunLazyFlatten(ds *map[string]interface{}, cfg *load.Config, api int) {
	// perform lazy flatten
	for _, flattenKey := range cfg.APIs[api].LazyFlatten {
		if strings.Contains(flattenKey, ">") {
			flatSplit := strings.Split(flattenKey, ">")
			if len(flatSplit) == 2 {
				if (*ds)[flatSplit[0]] != nil {
					switch (*ds)[flatSplit[0]].(type) {
					case map[string]interface{}:
						flat, err := flatten.Flatten((*ds)[flatSplit[0]].(map[string]interface{}), "", flatten.DotStyle)
						if err == nil {
							delete((*ds)[flatSplit[0]].(map[string]interface{}), flatSplit[1])
							(*ds)[flatSplit[0]].(map[string]interface{})[flatSplit[1]] = flat
						}
					case []interface{}:
						for i := range (*ds)[flatSplit[0]].([]interface{}) {
							switch (*ds)[flatSplit[0]].([]interface{})[i].(type) {
							case map[string]interface{}:
								// we need to flatten top level, then loop through and find the new keys and add back into the sample
								flat, err := flatten.Flatten((*ds)[flatSplit[0]].([]interface{})[i].(map[string]interface{}), "", flatten.DotStyle)
								if err == nil {
									// delete old data
									delete((*ds)[flatSplit[0]].([]interface{})[i].(map[string]interface{}), flatSplit[1])
									// depending if nested it may not be targeted correctly so auto set something in remove_keys - hacky workaround
									cfg.APIs[api].RemoveKeys = append(cfg.APIs[api].RemoveKeys, flatSplit[1]+"Samples")

									// add back into the datasample
									for k, v := range flat {
										if strings.Contains(k, flatSplit[1]) {
											(*ds)[flatSplit[0]].([]interface{})[i].(map[string]interface{})[k] = v
										}
									}
								}
							}
						}
					}
				}
			}
		} else {
			tmp := map[string]interface{}{"flat": (*ds)[flattenKey]}
			flat, err := flatten.Flatten(tmp, "", flatten.DotStyle)
			if err == nil {
				delete((*ds), flattenKey)
				(*ds)[flattenKey] = flat
			} else {
				logger.Flex("error", err, "unable to lazy_flatten", false)
			}
		}
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
			logger.Flex("error", err, fmt.Sprintf("%v math exp failed %v", newMetric, finalFormula), false)
		} else {
			result, err := expression.Evaluate(nil)
			if err != nil {
				logger.Flex("error", err, fmt.Sprintf("%v math evalute failed %v", newMetric, finalFormula), false)
			} else {
				(*currentSample)[newMetric] = result
			}
		}
	}
}

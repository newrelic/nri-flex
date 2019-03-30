package processor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/jeremywohl/flatten"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// CreateMetricSets creates metric sets
func CreateMetricSets(samples []interface{}, config *load.Config, i int) {
	api := config.APIs[i]
	// as it stands we know that this always receives map[string]interface{}'s
	for _, sample := range samples {

		// event limiter
		if (load.EventCount > load.Args.EventLimit) && load.Args.EventLimit != 0 {
			load.EventDropCount++
			if load.EventDropCount == 1 { // don't output the message more then once
				logger.Flex("debug",
					fmt.Errorf("event Limit %d has been reached, please increase if required", load.Args.EventLimit),
					"", false)
			}
			break
		}

		currentSample := sample.(map[string]interface{})
		eventType := "UnknownSample" // set an UnknownSample event name
		SetEventType(&currentSample, &eventType, api.EventType, api.Merge, api.Name)

		// modify existing sample before final processing
		createSample := true
		for k, v := range currentSample {
			key := k
			progress := true
			RunKeyRemover(api.RemoveKeys, &key, &progress, &currentSample)
			keyReplaced := false

			if progress {
				RunKeyConversion(&key, api, v)
				RunValConversion(&v, api, &key)
				RunValueParser(&v, api, &key)
				RunPluckNumbers(&v, api, &key)
				RunSubParse(api.SubParse, &currentSample, key, v) // subParse key pairs (see redis example)
				RunKeyReplace(api.ReplaceKeys, &keyReplaced, &key)
				RunKeyRenamer(api.RenameKeys, &keyReplaced, &key)                    // use key renamer if key replace hasn't occurred
				StoreLookups(api.StoreLookups, &key, &config.LookupStore, &v)        // store lookups
				VariableLookups(api.StoreVariables, &key, &config.VariableStore, &v) // store variable

				currentSample[key] = v
				if key != k {
					delete(currentSample, k)
				}

				// check if this contains any key pair values to filter out
				RunSampleFilter(api.SampleFilter, &createSample, key, v)
				// if keepkeys used will do inverse
				RunKeepKeys(api.KeepKeys, &key, &currentSample, &k)
				RunSampleRenamer(api.RenameSamples, &currentSample, key, &eventType)
			}
		}

		if createSample {
			RunMathCalculations(&api.Math, &currentSample)

			load.EventCount++
			load.EventDistribution[eventType]++

			// add custom attribute(s)
			// global
			for k, v := range config.CustomAttributes {
				currentSample[k] = v
			}
			// nested
			for k, v := range api.CustomAttributes {
				currentSample[k] = v
			}

			// inject some additional attributes if set
			if config.Global.BaseURL != "" {
				currentSample["baseUrl"] = config.Global.BaseURL
			}

			var metricSet *metric.Set
			// metricSet2 := configureNewMetricSet(&currentSample, api)
			// if metric parser is used, we need to namespace metrics for rate and delta support
			if len(api.MetricParser.Metrics) > 0 {
				useDefaultNamespace := false
				if api.MetricParser.Namespace.CustomAttr != "" {
					metricSet = load.Entity.NewMetricSet(eventType, metric.Attr("namespace", api.MetricParser.Namespace.CustomAttr))
				} else if len(api.MetricParser.Namespace.ExistingAttr) == 1 {
					nsKey := api.MetricParser.Namespace.ExistingAttr[0]
					nsVal := currentSample[nsKey]
					switch nsVal := nsVal.(type) {
					case string:
						metricSet = load.Entity.NewMetricSet(eventType, metric.Attr(nsKey, nsVal))
						delete(currentSample, nsKey) // can delete from sample as already set via namespace key
					default:
						useDefaultNamespace = true
					}
				} else if len(api.MetricParser.Namespace.ExistingAttr) > 1 {
					finalValue := ""
					for i, k := range api.MetricParser.Namespace.ExistingAttr {
						if currentSample[k] != nil {
							if i == 0 {
								finalValue = fmt.Sprintf("%v", currentSample[k])
							} else {
								finalValue = finalValue + "-" + fmt.Sprintf("%v", currentSample[k])
							}
						}
					}
					if finalValue != "" {
						metricSet = load.Entity.NewMetricSet(eventType, metric.Attr("namespace", finalValue))
					} else {
						useDefaultNamespace = true
					}
				}

				if useDefaultNamespace {
					logger.Flex("debug", fmt.Errorf("defaulting a namespace for:%v", api.Name), "", false)
					metricSet = load.Entity.NewMetricSet(eventType, metric.Attr("namespace", api.Name))
				}
			} else {
				metricSet = load.Entity.NewMetricSet(eventType)
			}

			// set default attribute(s)
			logger.Flex("debug", metricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE), "", false)
			logger.Flex("debug", metricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE), "", false)

			//add sample metrics
			for k, v := range currentSample {
				// key filter could be put here
				AutoSetMetric(k, v, metricSet, api.MetricParser.Metrics, api.MetricParser.AutoSet)
			}
		}

	}
}

// SetEventType sets the metricSet's eventType
func SetEventType(currentSample *map[string]interface{}, eventType *string, apiEventType string, apiMerge string, apiName string) {
	// if event_type is set use this, else attempt to autoset
	if (*currentSample)["event_type"] != nil && (*currentSample)["event_type"].(string) == "flexError" {
		*eventType = (*currentSample)["event_type"].(string)
		delete((*currentSample), "event_type")
	} else if apiEventType != "" && apiMerge == "" {
		*eventType = apiEventType
		delete((*currentSample), "event_type")
	} else {
		// pull out the event name, and remove if "Samples" is plural
		// if event_type not set, auto create via api name
		if (*currentSample)["event_type"] != nil {
			*eventType = (*currentSample)["event_type"].(string)
			if strings.Contains(*eventType, "Samples") {
				*eventType = strings.Replace(*eventType, "Samples", "Sample", -1)
			}
		} else {
			*eventType = apiName + "Sample"
		}
		delete((*currentSample), "event_type")
	}
}

// RunSampleRenamer using regex if sample has a key that matches, make that a different sample (event_type)
func RunSampleRenamer(renameSamples map[string]string, currentSample *map[string]interface{}, key string, eventType *string) {
	for regex, newEventType := range renameSamples {
		if formatter.KvFinder("regex", key, regex) {
			(*currentSample)["event_type"] = newEventType
			*eventType = newEventType
			break
		}
	}
}

// RunKeepKeys remove all other keys and keep these
func RunKeepKeys(keepKeys []string, key *string, currentSample *map[string]interface{}, k *string) {
	if len(keepKeys) > 0 {
		foundKey := false
		for _, keepKey := range keepKeys {
			if formatter.KvFinder("regex", *key, keepKey) {
				foundKey = true
				break
			}
		}
		if !foundKey {
			delete(*currentSample, *k)
		}
	}
}

// RunKeyRemover Remove unwanted keys
func RunKeyRemover(removeKeys []string, key *string, progress *bool, currentSample *map[string]interface{}) {
	for _, removeKey := range removeKeys {
		if formatter.KvFinder("regex", *key, removeKey) {
			*progress = false
			delete(*currentSample, *key)
			break
		}
	}
}

// RunKeyReplace String replace within a key
func RunKeyReplace(replaceKeys map[string]string, keyReplaced *bool, key *string) {
	for replaceKey, replaceVal := range replaceKeys {
		if formatter.KvFinder("regex", *key, replaceKey) {
			*key = strings.Replace(*key, replaceKey, replaceVal, -1)
			*keyReplaced = true
			break
		}
	}
}

// RunKeyRenamer Rename a key
func RunKeyRenamer(renameKeys map[string]string, keyReplaced *bool, key *string) {
	if !*keyReplaced {
		for renameKey, renameVal := range renameKeys {
			if strings.Contains(*key, renameKey) {
				*key = strings.Replace(*key, renameKey, renameVal, -1)
			}
		}
	}
}

// StoreLookups if key is found (using regex), store the values in the lookupStore as the defined lookupStoreKey for later use
func StoreLookups(storeLookups map[string]string, key *string, lookupStore *map[string][]string, v *interface{}) {
	for lookupStoreKey, lookupFindKey := range storeLookups {
		if formatter.KvFinder("regex", *key, lookupFindKey) {
			if *lookupStore == nil {
				*lookupStore = map[string][]string{}
			}
			(*lookupStore)[lookupStoreKey] = append((*lookupStore)[lookupStoreKey], fmt.Sprintf("%v", *v))
		}
	}
}

// VariableLookups if key is found (using regex), store the value in the variableStore, as the defined by the variableStoreKey for later use
func VariableLookups(variableLookups map[string]string, key *string, variableStore *map[string]string, v *interface{}) {
	for variableStoreKey, variableFindKey := range variableLookups {
		if *key == variableFindKey {
			if (*variableStore) == nil {
				(*variableStore) = map[string]string{}
			}
			(*variableStore)[variableStoreKey] = fmt.Sprintf("%v", *v)
		}
	}
}

// ProcessSamplesToMerge used to merge multiple samples together
func ProcessSamplesToMerge(samplesToMerge *map[string][]interface{}, yml *load.Config) {
	for eventType, sampleSet := range *samplesToMerge {
		newSample := map[string]interface{}{}
		newSample["event_type"] = eventType
		for _, sample := range sampleSet {
			prefix := yml.APIs[sample.(map[string]interface{})["_sampleNo"].(int)].Prefix
			for k, v := range sample.(map[string]interface{}) {
				if k != "_sampleNo" {
					newSample[prefix+k] = v
				}
			}
		}
		CreateMetricSets([]interface{}{newSample}, yml, 0)
	}
}

// FlattenData flatten an interface
func FlattenData(unknown interface{}, data map[string]interface{}, key string, sampleKeys map[string]string) map[string]interface{} {
	switch unknown := unknown.(type) {
	case []interface{}:
		dataSamples := []interface{}{}
		dataSamples = append(dataSamples, unknown...)
		// for _, loopVal := range unknown.([]interface{}) {
		// 	dataSamples = append(dataSamples, loopVal)
		// }

		// Check if Prometheus Style Metrics else process as normal (FlexSamples)
		if checkPrometheus(dataSamples) {
			data[key+"PrometheusSamples"] = dataSamples
		} else {
			key = checkPluralSlice(key)
			data[key+"FlexSamples"] = dataSamples
		}
	case map[string]interface{}:
		for loopKey := range unknown {
			finalKey := loopKey
			if key != "" {
				finalKey = key + "." + loopKey
			}

			// Check sample keys and convert to samples using ">" as the sample key splitter if defined
			// Knowing the sampleKey itself isn't really needed as these get turned into samples
			for _, sampleVal := range sampleKeys {
				keys := strings.Split(sampleVal, ">")
				flexSamples := []interface{}{}
				// if the one of the keys == the loopKey we know to create samples
				if len(keys) > 0 && keys[0] == loopKey {
					switch unknown[loopKey].(type) {
					case map[string]interface{}:
						dataSamples := unknown[loopKey].(map[string]interface{})
						for dataSampleKey, dataSample := range dataSamples {
							newSample := dataSample.(map[string]interface{})
							newSample[keys[1]] = dataSampleKey
							flexSamples = append(flexSamples, FlattenData(newSample, map[string]interface{}{}, "", sampleKeys))
						}
						unknown[loopKey] = flexSamples
					}

				}
			}

			FlattenData(unknown[loopKey], data, finalKey, sampleKeys)
		}
	default:
		data[key] = unknown
	}

	for dataKey := range data {
		// separately flatten the flex samples, adding them back into the slice with a new key
		// & removing the old from data thus a replace
		if strings.Contains(dataKey, "FlexSamples") {
			strippedDataKey, newSamples := processFlexSamples(dataKey, data[dataKey].([]interface{}), sampleKeys)
			data[strippedDataKey] = newSamples
			delete(data, dataKey)
		}
	}

	return data
}

// FinalMerge Perform final data merging
// Separates detected samples and already flattened attributes
func FinalMerge(data map[string]interface{}) []interface{} {
	finalAttributes := map[string]interface{}{}
	finalSampleSets := map[string]interface{}{}

	// store all flat attributes in final attributes
	// store detected samples in SampleSets
	for key := range data {
		if !strings.Contains(key, "Samples") {
			finalAttributes[key] = data[key]
		} else {
			finalSampleSets[key] = data[key]
		}
	}

	finalMergedSamples := []interface{}{}
	for sampleSet := range finalSampleSets {
		for _, sample := range finalSampleSets[sampleSet].([]interface{}) {
			switch sample := sample.(type) {
			case map[string]interface{}:
				newSample := sample
				newSample["event_type"] = sampleSet
				for attribute := range finalAttributes {
					newSample[attribute] = finalAttributes[attribute]
				}
				finalMergedSamples = append(finalMergedSamples, newSample)
			default:
				log.Debug("not sure what to do with this?")
				log.Debug(fmt.Sprintf("%v", sample))
			}
		}
	}

	if len(finalMergedSamples) > 0 {
		return finalMergedSamples
	}

	finalMergedSamples = append(finalMergedSamples, finalAttributes)
	return finalMergedSamples
}

// checkPrometheus Checks if the slice appears to be in a Prometheus style format
// some code duplication this can probably be cleaned up
func checkPrometheus(dataSamples []interface{}) bool {
	// needed when only 1 value set is returned from prometheus
	if len(dataSamples) == 2 {
		//check if the first value (timestamp) is a parse-able to a float
		value := fmt.Sprintf("%v", dataSamples[0])
		_, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return true
		}
	}

	for _, dataSample := range dataSamples {
		switch dataSample := dataSample.(type) {
		case []interface{}:
			//there should be 2 values a timestamp and value eg. [ 1435781430.781, "1" ]
			if len(dataSample) == 2 {
				//check if the first value (timestamp) is a parse-able to a float
				value := fmt.Sprintf("%v", dataSample[0])
				_, err := strconv.ParseFloat(value, 64)
				if err == nil {
					return true
				}
			}
		default:
			return false
		}
	}
	return false
}

// processFlexSamples Processes Flex detected samples
func processFlexSamples(dataKey string, dataSamples []interface{}, sampleKeys map[string]string) (string, []interface{}) {
	newSamples := []interface{}{}
	for _, sample := range dataSamples {
		sampleFlatten := FlattenData(sample, map[string]interface{}{}, "", sampleKeys)
		if sampleFlatten["valuesPrometheusSamples"] != nil {
			for _, prometheusSample := range sampleFlatten["valuesPrometheusSamples"].([]interface{}) {
				// this could be optimized
				newSample := FlattenData(sample, map[string]interface{}{}, "", sampleKeys)
				newSample["timestamp"] = int(prometheusSample.([]interface{})[0].(float64))
				newSample["value"] = prometheusSample.([]interface{})[1]
				delete(newSample, "valuesPrometheusSamples")
				newSamples = append(newSamples, newSample)
			}
		} else if sampleFlatten["valuePrometheusSamples"] != nil {
			newSample := FlattenData(sample, map[string]interface{}{}, "", sampleKeys)
			newSample["timestamp"] = int(sampleFlatten["valuePrometheusSamples"].([]interface{})[0].(float64))
			newSample["value"] = sampleFlatten["valuePrometheusSamples"].([]interface{})[1]
			delete(newSample, "valuePrometheusSamples")
			newSamples = append(newSamples, newSample)
		} else {
			newSamples = append(newSamples, sampleFlatten)
		}
	}
	strippedDataKey := strings.Replace(dataKey, "Flex", "", -1)
	return strippedDataKey, newSamples
}

// checkPluralSlice Checks if a key is plural to use for FlexSample naming
// An assumption is made that the slice key is plural
func checkPluralSlice(key string) string {
	if len(key) > 0 {
		if key[len(key)-1:] == "s" {
			return key[:len(key)-1]
		}
	}
	return key
}

// RunSampleFilter Filters samples generated
func RunSampleFilter(sampleFilters []map[string]string, createSample *bool, key string, v interface{}) {
	for _, sampleFilter := range sampleFilters {
		for regKey, regVal := range sampleFilter {
			regKeyFound := false
			regValFound := false
			if regKey != "" {
				validateKey := regexp.MustCompile(regKey)
				if validateKey.MatchString(key) {
					regKeyFound = true
				}
			}
			if regVal != "" {
				validateVal := regexp.MustCompile(regVal)
				if validateVal.MatchString(fmt.Sprintf("%v", v)) {
					regValFound = true
				}
			}
			if regKeyFound && regValFound {
				*createSample = false
			}
		}
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

// RunKeyConversion handles to lower and snake to camel case for keys
func RunKeyConversion(key *string, api load.API, v interface{}) {
	if api.ToLower {
		*key = strings.ToLower(*key)
	}
	if api.ConvertSpace != "" {
		*key = strings.Replace(*key, " ", api.ConvertSpace, -1)
	}
	if api.SnakeToCamel {
		formatter.SnakeCaseToCamelCase(key)
	}
}

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

// RunValueParser use regex to find a key, and pluck out its value by regex
func RunValueParser(v *interface{}, api load.API, key *string) {
	for regexKey, regexVal := range api.ValueParser {
		if formatter.KvFinder("regex", *key, regexKey) {
			value := fmt.Sprintf("%v", *v)
			*v = formatter.ValueParse(value, regexVal)
		}
	}
}

// RunPluckNumbers pluck numbers out automatically with ValueParser
func RunPluckNumbers(v *interface{}, api load.API, key *string) {
	//"sample_start_time = 1552864614.137869 (Sun, 17 Mar 2019 23:16:54 GMT)"
	// return 1552864614.137869
	if api.PluckNumbers {
		value := fmt.Sprintf("%v", *v)
		*v = formatter.ValueParse(value, `[+-]?([0-9]*\.?[0-9]+|[0-9]+\.?[0-9]*)([eE][+-]?[0-9]+)?`)
	}
}

// AutoSetMetric parse to number
func AutoSetMetric(k string, v interface{}, metricSet *metric.Set, metrics map[string]string, autoSet bool) {
	value := fmt.Sprintf("%v", v)
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || strings.EqualFold(value, "infinity") {
		logger.Flex("debug", metricSet.SetMetric(k, value, metric.ATTRIBUTE), "", false)
	} else {
		foundKey := false
		for metricKey, metricVal := range metrics {
			if (k == metricKey) || (autoSet && formatter.KvFinder("regex", k, metricKey)) {
				if metricVal == "RATE" {
					foundKey = true
					logger.Flex("debug", metricSet.SetMetric(k, parsed, metric.RATE), "", false)
					break
				} else if metricVal == "DELTA" {
					foundKey = true
					logger.Flex("debug", metricSet.SetMetric(k, parsed, metric.DELTA), "", false)
					break
				}
			}
		}
		if !foundKey {
			logger.Flex("debug", metricSet.SetMetric(k, parsed, metric.GAUGE), "", false)
		}
	}
}

// FindStartKey start at a different section of a payload
func FindStartKey(startKeys []string, mainDataset *map[string]interface{}) {
	for _, startKey := range startKeys {
		if (*mainDataset)[startKey] != nil {
			switch mainDs := (*mainDataset)[startKey].(type) {
			case map[string]interface{}:
				*mainDataset = mainDs
			case []interface{}:
				*mainDataset = map[string]interface{}{startKey: mainDs}
			}
		}
	}
}

// StripKeys strip defined keys out
func StripKeys(stripKeys []string, ds *map[string]interface{}) {
	for _, stripKey := range stripKeys {
		delete(*ds, stripKey)
		if strings.Contains(stripKey, ">") {
			stripSplit := strings.Split(stripKey, ">")
			if len(stripSplit) == 2 {
				if (*ds)[stripSplit[0]] != nil {
					switch (*ds)[stripSplit[0]].(type) {
					case map[string]interface{}:
						delete((*ds)[stripSplit[0]].(map[string]interface{}), stripSplit[1])
					case []interface{}:
						for i := range (*ds)[stripSplit[0]].([]interface{}) {
							switch (*ds)[stripSplit[0]].([]interface{})[i].(type) {
							case map[string]interface{}:
								// fmt.Println(ds[stripSplit[0]].([]interface{})[i].(map[string]interface{})[stripSplit[1]])
								delete((*ds)[stripSplit[0]].([]interface{})[i].(map[string]interface{}), stripSplit[1])
							}
						}
					}
				}
			}
		}
	}
}

// RunLazyFlatten lazy flattens the payload
func RunLazyFlatten(lazyFlatten []string, ds *map[string]interface{}) {
	// perform lazy flatten
	for _, flattenKey := range lazyFlatten {
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
			logger.Flex("debug", err, fmt.Sprintf("%v math exp failed %v", newMetric, finalFormula), false)
		} else {
			result, err := expression.Evaluate(nil)
			if err != nil {
				logger.Flex("debug", err, fmt.Sprintf("%v math evalute failed %v", newMetric, finalFormula), false)
			} else {
				(*currentSample)[newMetric] = result
			}
		}
	}
}

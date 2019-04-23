package processor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/newrelic/infra-integrations-sdk/data/event"

	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// CreateMetricSets creates metric sets
func CreateMetricSets(samples []interface{}, config *load.Config, i int) {
	api := config.APIs[i]
	// as it stands we know that this always receives map[string]interface{}'s
	for _, sample := range samples {
		// event limiter
		load.FlexStatusCounter.Lock()
		if (load.FlexStatusCounter.M["EventCount"] > load.Args.EventLimit) && load.Args.EventLimit != 0 {
			load.FlexStatusCounter.M["EventDropCount"]++
			if load.FlexStatusCounter.M["EventDropCount"] == 1 { // don't output the message more then once
				logger.Flex("debug",
					fmt.Errorf("event Limit %d has been reached, please increase if required", load.Args.EventLimit),
					"", false)
			}
			load.FlexStatusCounter.Unlock()
			break
		}
		load.FlexStatusCounter.Unlock()

		currentSample := sample.(map[string]interface{})
		eventType := "UnknownSample" // set an UnknownSample event name
		SetEventType(&currentSample, &eventType, api.EventType, api.Merge, api.Name)

		// modify existing sample before final processing
		createSample := true
		SkipProcessing := api.SkipProcessing
		for k, v := range currentSample { // k == original key
			key := k
			progress := true
			RunKeyRemover(api.RemoveKeys, &key, &progress, &currentSample)

			if progress {
				RunKeyConversion(&key, api, v, &SkipProcessing)
				RunValConversion(&v, api, &key)
				RunValueParser(&v, api, &key)
				RunPluckNumbers(&v, api, &key)
				RunSubParse(api.SubParse, &currentSample, key, v) // subParse key pairs (see redis example)
				RunValueTransformer(&v, api, &key)                // Needs to be run before KeyRenamer and KeyReplacer
				RunKeyRenamer(api.RenameKeys, &key)               // use key renamer if key replace hasn't occurred
				RunKeyRenamer(api.ReplaceKeys, &key)              // kept for backwards compatibility with replace_keys

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

			load.StatusCounterIncrement("EventCount")
			load.StatusCounterIncrement(eventType)

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

			workingEntity := setEntity(api.Entity, api.EntityType) // default type instance

			var metricSet *metric.Set
			// if metric parser is used, we need to namespace metrics for rate and delta support
			if len(api.MetricParser.Metrics) > 0 {
				useDefaultNamespace := false
				if api.MetricParser.Namespace.CustomAttr != "" {
					metricSet = workingEntity.NewMetricSet(eventType, metric.Attr("namespace", api.MetricParser.Namespace.CustomAttr))
				} else if len(api.MetricParser.Namespace.ExistingAttr) == 1 {
					nsKey := api.MetricParser.Namespace.ExistingAttr[0]
					switch nsVal := currentSample[nsKey].(type) {
					case string:
						metricSet = workingEntity.NewMetricSet(eventType, metric.Attr(nsKey, nsVal))
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
						metricSet = workingEntity.NewMetricSet(eventType, metric.Attr("namespace", finalValue))
					} else {
						useDefaultNamespace = true
					}
				}

				if useDefaultNamespace {
					logger.Flex("debug", fmt.Errorf("defaulting a namespace for:%v", api.Name), "", false)
					metricSet = workingEntity.NewMetricSet(eventType, metric.Attr("namespace", api.Name))
				}
			} else {
				metricSet = workingEntity.NewMetricSet(eventType)
			}

			// set default attribute(s)
			logger.Flex("debug", metricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE), "", false)
			logger.Flex("debug", metricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE), "", false)

			//add sample metrics
			for k, v := range currentSample {
				// add prefixing, prefixing for merged samples done elsewhere
				if api.Prefix != "" && api.Merge == "" {
					k = api.Prefix + k
				}

				StoreLookups(api.StoreLookups, &k, &config.LookupStore, &v)        // store lookups
				VariableLookups(api.StoreVariables, &k, &config.VariableStore, &v) // store variable

				if api.InventoryOnly {
					setInventory(workingEntity, api.Inventory, k, v)
				} else if api.EventsOnly {
					setEvents(workingEntity, api.Events, k, v)
				} else {
					// these can be set async
					var wg sync.WaitGroup
					wg.Add(3)
					go func() {
						defer wg.Done()
						setInventory(workingEntity, api.Inventory, k, v)
					}()
					go func() {
						defer wg.Done()
						setEvents(workingEntity, api.Events, k, v)
					}()
					go func() {
						defer wg.Done()
						AutoSetMetric(k, v, metricSet, api.MetricParser.Metrics, api.MetricParser.AutoSet)
					}()
					wg.Wait()
				}
			}
		}

	}
}

// setInventory sets infrastructure inventory metrics
func setInventory(entity *integration.Entity, inventory map[string]string, k string, v interface{}) {
	if inventory[k] != "" {
		if inventory[k] == "value" {
			logger.Flex("debug", entity.SetInventoryItem(k, "value", v), "", false)
		} else {
			logger.Flex("debug", entity.SetInventoryItem(inventory[k], k, v), "", false)
		}
	}
}

// setInventory sets infrastructure inventory metrics
func setEvents(entity *integration.Entity, inventory map[string]string, k string, v interface{}) {
	if inventory[k] != "" {
		value := fmt.Sprintf("%v", v)
		if inventory[k] != "default" {
			err := entity.AddEvent(&event.Event{
				Summary:  value,
				Category: inventory[k],
			})
			logger.Flex("debug", err, "", false)
		} else {
			err := entity.AddEvent(&event.Event{
				Summary:  value,
				Category: k,
			})
			logger.Flex("debug", err, "", false)
		}
	}
}

// setEntity sets the entity to be used for the configured API
// defaults the type aka namespace to instance
func setEntity(entity string, customNamespace string) *integration.Entity {
	if entity != "" {
		if customNamespace == "" {
			customNamespace = "instance"
		}
		workingEntity, err := load.Integration.Entity(entity, customNamespace)
		if err == nil {
			return workingEntity
		}
	}
	return load.Entity
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
		switch ss := finalSampleSets[sampleSet].(type) {
		case []interface{}:
			for _, sample := range ss {
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
		case map[string]interface{}:
			newSample := ss
			newSample["event_type"] = sampleSet
			for attribute := range finalAttributes {
				newSample[attribute] = finalAttributes[attribute]
			}
			finalMergedSamples = append(finalMergedSamples, newSample)
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
				} else if metricVal == "ATTRIBUTE" {
					foundKey = true
					logger.Flex("debug", metricSet.SetMetric(k, value, metric.ATTRIBUTE), "", false)
					break
				}
			}
		}
		if !foundKey {
			logger.Flex("debug", metricSet.SetMetric(k, parsed, metric.GAUGE), "", false)
		}
	}
}

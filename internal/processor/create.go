/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/event"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

const regex = "regex"

// CreateMetricSets creates metric sets
// hren added samplesToMerge parameter, moved merge operation to CreateMetricSets so that the "Run...." functions still apply before merge
func CreateMetricSets(samples []interface{}, config *load.Config, i int, mergeMetric bool, samplesToMerge *map[string][]interface{}) {
	api := config.APIs[i]
	// as it stands we know that this always receives map[string]interface{}'s
	for _, sample := range samples {
		currentSample := sample.(map[string]interface{})
		eventType := "UnknownSample" // set an UnknownSample event name
		SetEventType(&currentSample, &eventType, api.EventType, api.Merge, api.Name)

		// init lookup store
		if (&config.LookupStore) == nil {
			config.LookupStore = map[string]map[string]struct{}{}
		}

		// event limiter
		if (load.StatusCounterRead("EventCount") > load.Args.EventLimit) && load.Args.EventLimit != 0 {
			load.StatusCounterIncrement("EventDropCount")
			if load.StatusCounterRead("EventDropCount") == 1 { // don't output the message more then once
				load.Logrus.Errorf("flex: event limit %d has been reached, please increase if required", load.Args.EventLimit)
			}
			break
		}

		// modify existing sample before final processing
		SkipProcessing := api.SkipProcessing

		var modifiedKeys []string
		for k, v := range currentSample { // k == original key
			key := k
			RunKeyConversion(&key, api, v, &SkipProcessing)
			RunValConversion(&v, api, &key)
			RunValueParser(&v, api, &key)
			RunPluckNumbers(&v, api, &key)
			RunSubParse(api.SubParse, &currentSample, key, v) // subParse key pairs (see redis example)
			RunValueTransformer(&v, api, &key)                // Needs to be run before KeyRenamer and KeyReplacer

			RunValueMapper(api.ValueMapper, &currentSample, key, &v) // subParse key pairs (see redis example)
			// do not rename a key again, this is to avoid continuous replacement loops
			// eg. if you replace id with project.id
			// this could then again attempt to replace id within project.id to project.project.id
			if !sliceContains(modifiedKeys, k) {
				RunKeyRenamer(api.RenameKeys, &key, k)  // use key renamer if key replace hasn't occurred
				RunKeyRenamer(api.ReplaceKeys, &key, k) // kept for backwards compatibility with replace_keys
			}

			currentSample[key] = v
			if key != k {
				modifiedKeys = append(modifiedKeys, key)
				delete(currentSample, k)
			}

			StoreLookups(api.StoreLookups, &key, &config.LookupStore, &v)        // store lookups
			VariableLookups(api.StoreVariables, &key, &config.VariableStore, &v) // store variable

			// if keepkeys used will do inverse
			RunKeepKeys(api.KeepKeys, &key, &currentSample)
			RunSampleRenamer(api.RenameSamples, &currentSample, key, &eventType)
		}

		createSample := true
		// check if we should ignore this output completely
		// useful when requests are made to generate a lookup, but the data is not needed
		if api.IgnoreOutput {
			createSample = false
		} else {
			// check if this contains any key pair values to filter out
			RunSampleFilter(currentSample, api.SampleFilter, &createSample)
		}

		if createSample {
			// remove keys from sample
			RunKeyRemover(&currentSample, api.RemoveKeys)

			RunMathCalculations(&api.Math, &currentSample)

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

			addAttribute(currentSample, api.AddAttribute)

			// hren: if it is not mergeMetric, it will proceed to puslish metric
			if !mergeMetric {
				workingEntity := setEntity(api.Entity, api.EntityType) // default type instance
				if config.MetricAPI {
					AutoSetMetricAPI(&currentSample, &api)
				} else {
					AutoSetStandard(&currentSample, &api, workingEntity, eventType, config)
				}
			} else {
				// hren: it is mergeMetric, add the metric to mergeData, which will be published later
				currentSample["_sampleNo"] = i
				// hren overwrite event_type if it is merge operation
				currentSample["event_type"] = config.APIs[i].Merge
				(*samplesToMerge)[config.APIs[i].Merge] = append((*samplesToMerge)[config.APIs[i].Merge], currentSample)
			}

		}

	}
}

// setInventory sets infrastructure inventory metrics
func setInventory(entity *integration.Entity, inventory map[string]string, k string, v interface{}) {
	if inventory[k] != "" {
		if inventory[k] == "value" {
			checkError(entity.SetInventoryItem(k, "value", v))
		} else {
			checkError(entity.SetInventoryItem(inventory[k], k, v))
		}
	}
}

// setInventory sets infrastructure inventory metrics
func setEvents(entity *integration.Entity, inventory map[string]string, k string, v interface{}) {
	if inventory[k] != "" {
		value := cleanValue(&v)
		if inventory[k] != "default" {
			checkError(entity.AddEvent(&event.Event{
				Summary:  value,
				Category: inventory[k],
			}))
		} else {
			checkError(entity.AddEvent(&event.Event{
				Summary:  value,
				Category: k,
			}))
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

func cleanEvent(event string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	event = re.ReplaceAllLiteralString(event, "_")
	event = strings.TrimPrefix(event, "_")
	return event
}

// SetEventType sets the metricSet's eventType
func SetEventType(currentSample *map[string]interface{}, eventType *string, apiEventType string, apiMerge string, apiName string) {
	// if event_type is set use this, else attempt to autoset
	if (*currentSample)["event_type"] != nil && (*currentSample)["event_type"].(string) == "flexError" {
		*eventType = (*currentSample)["event_type"].(string)
		delete(*currentSample, "event_type")
	} else if apiEventType != "" && apiMerge == "" {
		*eventType = apiEventType
		delete(*currentSample, "event_type")
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
		delete(*currentSample, "event_type")
	}
	*eventType = cleanEvent(*eventType)
}

// RunSampleRenamer using regex if sample has a key that matches, make that a different sample (event_type)
func RunSampleRenamer(renameSamples map[string]string, currentSample *map[string]interface{}, key string, eventType *string) {
	for regexLocal, newEventType := range renameSamples {
		if formatter.KvFinder(regex, key, regexLocal) {
			(*currentSample)["event_type"] = newEventType
			*eventType = newEventType
			break
		}
	}
}

// RunSampleFilter Filters samples generated
func RunSampleFilter(currentSample map[string]interface{}, sampleFilters []map[string]string, createSample *bool) {
	for _, sampleFilter := range sampleFilters {
		for key, v := range currentSample {
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
					if validateVal.MatchString(cleanValue(&v)) {
						regValFound = true
					}
				}
				if regKeyFound && regValFound {
					*createSample = false
				}
			}
		}
	}
}

// RunEventFilter filters events generated
func RunEventFilter(filters []load.Filter, createEvent *bool, k string, v interface{}) {
	for _, filter := range filters {
		value := cleanValue(&v)
		filterMode := filter.Mode
		if filterMode == "" {
			filterMode = regex
		}
		filterValue := filter.Value
		if filterValue == "" && filter.Mode == regex {
			filterValue = ".*"
		}
		if formatter.KvFinder(filterMode, k, filter.Key) && formatter.KvFinder(filterMode, value, filterValue) {
			*createEvent = false
			break
		}
	}
}

// RunKeyFilter filters keys generated
func RunKeyFilter(filters []load.Filter, currentSample *map[string]interface{}, k string) {
	foundKey := false
	filterInverse := false

	for _, filter := range filters {
		filterMode := filter.Mode
		filterInverse = filter.Inverse
		if filterMode == "" {
			filterMode = regex
		}
		if formatter.KvFinder(filterMode, k, filter.Key) {
			if filterInverse {
				foundKey = true
				break
			}
		}
	}

	// delete the key if not found, and being used in inverse mode
	if filterInverse && !foundKey {
		delete(*currentSample, k)
	}
}

// AutoSetMetricAPI automatically set metrics for use with the metric api
func AutoSetMetricAPI(currentSample *map[string]interface{}, api *load.API) {
	// set current time
	currentTime := time.Now().UnixNano() / 1e+6
	// set common attributes
	commonAttributes := map[string]interface{}{
		"integration_version": load.IntegrationVersion,
		"integration_name":    load.IntegrationName,
	}

	// store numeric values, as metrics within Metrics
	var Metrics []map[string]interface{}
	SummaryMetrics := map[string]map[string]float64{}

	//add sample metrics
	for k, v := range *currentSample {
		// add prefixing, prefixing for merged samples done elsewhere
		if api.Prefix != "" && api.Merge == "" {
			k = api.Prefix + k
		}
		value := cleanValue(&v)
		parsed, err := strconv.ParseFloat(value, 64)
		// any non numeric values, are stored as common attributes
		if err != nil || strings.EqualFold(value, "infinity") {
			commonAttributes[k] = value
		} else {
			currentMetric := map[string]interface{}{
				"name":  k,
				"value": parsed,
				"type":  "",
			}

			// check if counter
			for metricKey, intervalMs := range (*api).MetricParser.Counts {
				if k == metricKey {
					currentMetric["type"] = "count"
					currentMetric["interval.ms"] = intervalMs
					load.StatusCounterIncrement("CounterMetrics")
					Metrics = append(Metrics, currentMetric)
					break
				}
			}

			// check if summary
			if currentMetric["type"] == "" {
				for rootSummary, metricTypes := range (*api).MetricParser.Summaries {
					for metric, keyName := range metricTypes {
						if metric == "min" || metric == "sum" || metric == "max" || metric == "count" {
							if keyName == k {
								if SummaryMetrics[rootSummary] != nil {
									SummaryMetrics[rootSummary][metric] = parsed
								} else {
									SummaryMetrics[rootSummary] = map[string]float64{
										metric: parsed,
									}
								}
								currentMetric["type"] = "summary" // setting just to avoid the gauge default
							}
						}
					}
				}
			}

			// if type still not set, default to gauge
			if currentMetric["type"] == "" {
				currentMetric["type"] = "gauge"
				load.StatusCounterIncrement("GaugeMetrics")
				Metrics = append(Metrics, currentMetric)
			}
		}
	}

	// add summary metrics into final metrics for MetricsStore
	for summaryName, metrics := range SummaryMetrics {
		v := (*api).MetricParser.Summaries[summaryName]["interval"]
		value := cleanValue(&v)
		intervalParsed, err := strconv.ParseFloat(value, 64)
		if err == nil && len(metrics) == 4 { // should be 4 for min/max/value/count
			currentMetric := map[string]interface{}{
				"name":        summaryName,
				"value":       metrics,
				"type":        "summary",
				"interval.ms": intervalParsed,
			}
			load.StatusCounterIncrement("SummaryMetrics")
			Metrics = append(Metrics, currentMetric)
		}
	}

	MetricsPayload := load.Metrics{
		CommonAttributes: commonAttributes,
		TimestampMs:      currentTime,
		Metrics:          Metrics,
	}

	load.MetricsStoreAppend(MetricsPayload)
}

// AutoSetStandard x
func AutoSetStandard(currentSample *map[string]interface{}, api *load.API, workingEntity *integration.Entity, eventType string, config *load.Config) {
	load.StatusCounterIncrement("EventCount")
	load.StatusCounterIncrement(eventType)

	var metricSet *metric.Set
	// if metric parser is used, we need to namespace metrics for rate and delta support
	if len(api.MetricParser.Metrics) > 0 {
		useDefaultNamespace := false
		if api.MetricParser.Namespace.CustomAttr != "" {
			metricSet = workingEntity.NewMetricSet(eventType, metric.Attr("namespace", api.MetricParser.Namespace.CustomAttr))
		} else if len(api.MetricParser.Namespace.ExistingAttr) == 1 {
			nsKey := api.MetricParser.Namespace.ExistingAttr[0]
			switch nsVal := (*currentSample)[nsKey].(type) {
			case string:
				metricSet = workingEntity.NewMetricSet(eventType, metric.Attr(nsKey, nsVal))
				delete((*currentSample), nsKey) // can delete from sample as already set via namespace key
			default:
				useDefaultNamespace = true
			}
		} else if len(api.MetricParser.Namespace.ExistingAttr) > 1 {
			finalValue := ""
			for i, k := range api.MetricParser.Namespace.ExistingAttr {
				if (*currentSample)[k] != nil {
					v := (*currentSample)[k]
					value := cleanValue(&v)
					if i == 0 {
						finalValue = value
					} else {
						finalValue = finalValue + "-" + value
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
			load.Logrus.Debugf("flex: defaulting a namespace for:%v", api.Name)
			metricSet = workingEntity.NewMetricSet(eventType, metric.Attr("namespace", api.Name))
		}
	} else {
		metricSet = workingEntity.NewMetricSet(eventType)
	}

	// set default attribute(s)
	checkError(metricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE))
	checkError(metricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE))

	//add sample metrics
	for k, v := range *currentSample {
		// add prefixing, prefixing for merged samples done elsewhere
		if api.Prefix != "" && api.Merge == "" {
			k = api.Prefix + k
		}

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
				AutoSetMetricInfra(k, v, metricSet, api.MetricParser.Metrics, api.MetricParser.AutoSet, api.MetricParser.Mode)
			}()
			wg.Wait()
		}
	}
}

// AutoSetMetricInfra parse to number
func AutoSetMetricInfra(k string, v interface{}, metricSet *metric.Set, metrics map[string]string, autoSet bool, mode string) {
	value := cleanValue(&v)
	parsed, err := strconv.ParseFloat(value, 64)

	if err != nil || strings.EqualFold(value, "infinity") || strings.EqualFold(value, "inf") {
		checkError(metricSet.SetMetric(k, value, metric.ATTRIBUTE))
	} else {
		foundKey := false
		for metricKey, metricVal := range metrics {
			if (k == metricKey) || (autoSet && formatter.KvFinder(regex, k, metricKey)) || (mode != "" && formatter.KvFinder(mode, k, metricKey)) {
				if metricVal == "RATE" {
					foundKey = true
					checkError(metricSet.SetMetric(k, parsed, metric.RATE))
					break
				} else if metricVal == "DELTA" {
					foundKey = true
					checkError(metricSet.SetMetric(k, parsed, metric.DELTA))
					break
				} else if metricVal == "ATTRIBUTE" {
					foundKey = true
					checkError(metricSet.SetMetric(k, value, metric.ATTRIBUTE))
					break
				}
			}
		}
		if !foundKey {
			checkError(metricSet.SetMetric(k, parsed, metric.GAUGE))
		}
	}
}

func addAttribute(currentSample map[string]interface{}, addAttribute map[string]string) {
	// add attribute, use attributes from current sample to create new attributes like http links
	for key, val := range addAttribute {
		newAttributeValue := val
		variableReplaceOccurred := false
		// in the value of each attribute find the keys that need replacing
		variableReplaces := regexp.MustCompile(`\${.*?}`).FindAllString(val, -1)
		for _, variableReplace := range variableReplaces {
			replaceKey := strings.TrimSuffix(strings.TrimPrefix(variableReplace, "${"), "}")

			if currentSample[replaceKey] != nil {
				value := currentSample[replaceKey]
				replacementValue := cleanValue(&value)
				newAttributeValue = strings.Replace(newAttributeValue, variableReplace, replacementValue, -1)

				// check if the replacement occurred
				// if this check is not in place there will be a unneeded templated sample generated
				if strings.Contains(newAttributeValue, replacementValue) {
					variableReplaceOccurred = true
				}
			}
		}
		if variableReplaceOccurred {
			currentSample[key] = newAttributeValue
		}
	}
}

// sliceContains check if slice contains an attribute
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func checkError(err error) {
	if err != nil {
		load.Logrus.WithError(err).Error("flex: failed to set")
	}
}

// deprecated
// // checkPrometheus Checks if the slice appears to be in a Prometheus style format
// // some code duplication this can probably be cleaned up
// func checkPrometheus(dataSamples []interface{}) bool {
// 	// needed when only 1 value set is returned from prometheus
// 	if len(dataSamples) == 2 {
// 		//check if the first value (timestamp) is a parse-able to a float
// 		value := fmt.Sprintf("%v", dataSamples[0])
// 		_, err := strconv.ParseFloat(value, 64)
// 		if err == nil {
// 			return true
// 		}
// 	}

// 	for _, dataSample := range dataSamples {
// 		switch dataSample := dataSample.(type) {
// 		case []interface{}:
// 			//there should be 2 values a timestamp and value eg. [ 1435781430.781, "1" ]
// 			if len(dataSample) == 2 {
// 				//check if the first value (timestamp) is a parse-able to a float
// 				value := fmt.Sprintf("%v", dataSample[0])
// 				_, err := strconv.ParseFloat(value, 64)
// 				if err == nil {
// 					return true
// 				}
// 			}
// 		default:
// 			return false
// 		}
// 	}
// 	return false
// }

package parser

import (
	"fmt"
	"io"
	"math"
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// Family mirrors the MetricFamily proto message.
type Family struct {
	//Time    time.Time
	Name    string                         `json:"name"`
	Help    string                         `json:"help"`
	Type    string                         `json:"type"`
	Metrics map[int]map[string]interface{} `json:"metrics,omitempty"` // Either metric or summary.
}

// Prometheus from http io
func Prometheus(input io.Reader, dataStore *[]interface{}, api *load.API) {
	mfChan := make(chan *dto.MetricFamily, 1024)
	go func() {
		if err := ParseReader(input, mfChan); err != nil {
			logger.Flex("debug", err, "prometheus parsing failure", false)
		}
	}()

	// store the flattened sample
	flattenedSample := map[string]interface{}{}
	if api.Prometheus.FlattenedEvent != "" {
		flattenedSample["event_type"] = api.Prometheus.FlattenedEvent
	} else {
		flattenedSample["event_type"] = api.Name + "Sample"
	}

	// initialize blank sampleKeys
	sampleKeys := map[string]map[string]interface{}{}

	// add standard metric families into datastore
	for mf := range mfChan {
		prometheusNewFamily(mf, dataStore, api, &flattenedSample, &sampleKeys)
	}
	// anything sampled add into datastore
	for sample := range sampleKeys {
		*dataStore = append(*dataStore, sampleKeys[sample])
	}
	// add flattened sample into datastore
	if len(flattenedSample) > 0 {
		applyCustomAttributes(&flattenedSample, &api.Prometheus.CustomAttributes)
		*dataStore = append(*dataStore, flattenedSample)
	}
}

// NewFamily consumes a MetricFamily and transforms it to a map[string]interface{}
func prometheusNewFamily(dtoMF *dto.MetricFamily, dataStore *[]interface{}, api *load.API, flattenedSample *map[string]interface{}, sampleKeys *map[string]map[string]interface{}) {

	for _, m := range dtoMF.Metric {
		// do not show go exporter metrics unless enabled
		if !api.Prometheus.GoMetrics && strings.Contains(dtoMF.GetName(), "go_") {
			break
		}

		// store counter/gauge samples with the same meta together
		useSecondary := false
		buildKey := ""
		if len(m.Label) > 0 && dtoMF.GetType() != dto.MetricType_SUMMARY && dtoMF.GetType() != dto.MetricType_HISTOGRAM {
			useSecondary = true
			sample := map[string]interface{}{}
			for _, label := range m.Label {
				sample[label.GetName()] = label.GetValue()
				buildKey += label.GetValue()
			}
			if (*sampleKeys)[buildKey] == nil {
				(*sampleKeys)[buildKey] = sample
			}
		}

		metric := map[string]interface{}{}
		metric["name"] = dtoMF.GetName()
		metric["help"] = dtoMF.GetHelp()
		metric["type"] = dtoMF.GetType().String()
		applyCustomAttributes(&metric, &api.Prometheus.CustomAttributes)
		prometheusMakeLabels(m, &metric)

		// create a custom sample to store the metrics associated with the particular key
		customSample := ""
		baseSample := ""

		for sample, key := range api.Prometheus.SampleKeys {
			if metric[key] != nil { // found key
				baseSample = sample
				if customSample == "" {
					customSample = baseSample + "." + fmt.Sprintf("%v", metric[key])
				} else {
					customSample = customSample + "." + fmt.Sprintf("%v", metric[key])
				}
			}
		}

		// when using a custom sample, we store into a workingSample as we can't target an address when nested
		workingSample := map[string]interface{}{}
		if customSample != "" {
			if (*sampleKeys)[customSample] == nil {
				(*sampleKeys)[customSample] = map[string]interface{}{}
			} else {
				workingSample = (*sampleKeys)[customSample]
			}
			workingSample["event_type"] = baseSample
			applyCustomAttributes(&workingSample, &api.Prometheus.CustomAttributes)
			prometheusMakeLabels(m, &workingSample) // possibility that a colision could occur from other samples
		}

		if dtoMF.GetType() == dto.MetricType_SUMMARY {
			if (*api).Prometheus.Unflatten {
				metric["count"] = fmt.Sprint(m.GetSummary().GetSampleCount())
				metric["sum"] = fmt.Sprint(m.GetSummary().GetSampleSum())
				prometheusMakeQuantiles(m, &metric, dtoMF, api.Prometheus.Unflatten)
				*dataStore = append(*dataStore, metric)
			} else if customSample != "" {
				workingSample[dtoMF.GetName()] = fmt.Sprint(getValue(m))
				workingSample[dtoMF.GetName()+".count"] = fmt.Sprint(m.GetSummary().GetSampleCount())
				workingSample[dtoMF.GetName()+".sum"] = fmt.Sprint(m.GetSummary().GetSampleSum())
				prometheusMakeQuantiles(m, &workingSample, dtoMF, api.Prometheus.Unflatten)
				(*sampleKeys)[customSample] = workingSample
			} else if api.Prometheus.Summary {
				metric["count"] = fmt.Sprint(m.GetSummary().GetSampleCount())
				metric["sum"] = fmt.Sprint(m.GetSummary().GetSampleSum())
				defaultEvent := api.Name
				if api.Prometheus.FlattenedEvent != "" {
					defaultEvent = api.Prometheus.FlattenedEvent
				}
				if strings.Contains(defaultEvent, "Sample") {
					defaultEvent = strings.Replace(defaultEvent, "Sample", "SummarySample", -1)
				} else {
					defaultEvent += "SummarySample"
				}
				metric["event_type"] = defaultEvent
				prometheusMakeQuantiles(m, &metric, dtoMF, true)
				*dataStore = append(*dataStore, metric)
			}
		} else if dtoMF.GetType() == dto.MetricType_HISTOGRAM {
			if (*api).Prometheus.Unflatten {
				metric["count"] = fmt.Sprint(m.GetSummary().GetSampleCount())
				metric["sum"] = fmt.Sprint(m.GetSummary().GetSampleSum())
				prometheusMakeBuckets(m, &metric, dtoMF, api.Prometheus.Unflatten)
				*dataStore = append(*dataStore, metric)
			} else if customSample != "" {
				workingSample[dtoMF.GetName()] = fmt.Sprint(getValue(m))
				workingSample[dtoMF.GetName()+".count"] = fmt.Sprint(m.GetSummary().GetSampleCount())
				workingSample[dtoMF.GetName()+".sum"] = fmt.Sprint(m.GetSummary().GetSampleSum())
				prometheusMakeBuckets(m, &metric, dtoMF, api.Prometheus.Unflatten)
				(*sampleKeys)[customSample] = workingSample
			} else if api.Prometheus.Histogram {
				metric["count"] = fmt.Sprint(m.GetSummary().GetSampleCount())
				metric["sum"] = fmt.Sprint(m.GetSummary().GetSampleSum())
				defaultEvent := api.Name
				if api.Prometheus.FlattenedEvent != "" {
					defaultEvent = api.Prometheus.FlattenedEvent
				}
				if strings.Contains(defaultEvent, "Sample") {
					defaultEvent = strings.Replace(defaultEvent, "Sample", "HistogramSample", -1)
				} else {
					defaultEvent += "HistogramSample"
				}
				metric["event_type"] = defaultEvent
				prometheusMakeBuckets(m, &metric, dtoMF, true)
				*dataStore = append(*dataStore, metric)
			}
		} else { // gauge or counter
			metric["value"] = fmt.Sprint(getValue(m))

			if (*api).Prometheus.Unflatten {
				*dataStore = append(*dataStore, metric)
			} else if customSample != "" {
				workingSample[dtoMF.GetName()] = fmt.Sprint(getValue(m))
				(*sampleKeys)[customSample] = workingSample
			} else if len(m.Label) > 0 && useSecondary {
				key := dtoMF.GetName()
				if dtoMF.GetType() == dto.MetricType_GAUGE {
					key += ".gauge"
				} else if dtoMF.GetType() == dto.MetricType_COUNTER {
					key += ".counter"
				}
				(*sampleKeys)[buildKey][key] = fmt.Sprint(getValue(m))
			} else {
				key := dtoMF.GetName()
				for _, keyMerge := range api.Prometheus.KeyMerge {
					if metric[keyMerge] != nil {
						key = key + "." + fmt.Sprintf("%v", metric[keyMerge])
						break
					}
				}
				(*flattenedSample)["name"] = "main"
				(*flattenedSample)[key] = fmt.Sprint(getValue(m))
			}
		}
	}
	// load secondary datastore into main datastore
	// for _, sample := range secondaryDataStore {
	// 	*dataStore = append(*dataStore, sample)
	// }
}

func getValue(m *dto.Metric) float64 {
	if m.Gauge != nil {
		return m.GetGauge().GetValue()
	}
	if m.Counter != nil {
		return m.GetCounter().GetValue()
	}
	if m.Untyped != nil {
		return m.GetUntyped().GetValue()
	}
	return 0.
}

func prometheusMakeLabels(m *dto.Metric, metric *map[string]interface{}) {
	for _, lp := range m.Label {
		(*metric)[lp.GetName()] = lp.GetValue()
	}
}

func prometheusMakeQuantiles(m *dto.Metric, metric *map[string]interface{}, dtoMF *dto.MetricFamily, unflatten bool) {
	for _, q := range m.GetSummary().Quantile {
		if !math.IsNaN(q.GetValue()) {
			if unflatten {
				(*metric)[fmt.Sprintf("%f", q.GetQuantile())] = fmt.Sprint(q.GetValue())
			} else {
				(*metric)[(*dtoMF).GetName()+fmt.Sprintf(".%f", q.GetQuantile())] = fmt.Sprint(q.GetValue())
			}
		}
	}
}

func prometheusMakeBuckets(m *dto.Metric, metric *map[string]interface{}, dtoMF *dto.MetricFamily, unflatten bool) {
	for _, b := range m.GetHistogram().Bucket {
		if unflatten {
			(*metric)[fmt.Sprintf("%f", b.GetUpperBound())] = fmt.Sprint(b.GetCumulativeCount())
		} else {
			// (*metric)[fmt.Sprintf("%f", b.GetUpperBound())] = fmt.Sprint(b.GetCumulativeCount())
			(*metric)[(*dtoMF).GetName()+fmt.Sprintf(".%f", b.GetUpperBound())] = fmt.Sprint(b.GetCumulativeCount())
		}
	}
}

// ParseReader consumes an io.Reader and pushes it to the MetricFamily
// channel. It returns when all MetricFamilies are parsed and put on the
// channel.
func ParseReader(in io.Reader, ch chan<- *dto.MetricFamily) error {
	defer close(ch)
	// We could do further content-type checks here, but the
	// fallback for now will anyway be the text format
	// version 0.0.4, so just go for it and see if it works.
	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(in)
	if err != nil {
		return fmt.Errorf("reading text format failed: %v", err)
	}
	for _, mf := range metricFamilies {
		ch <- mf
	}
	return nil
}

// applyCustomAttributes applies custom attributes to the provided sample
func applyCustomAttributes(sample *map[string]interface{}, attributes *map[string]string) {
	for key, val := range *attributes {
		(*sample)[key] = val
	}
}

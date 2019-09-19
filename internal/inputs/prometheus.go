package inputs

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"

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
func Prometheus(dataStore *[]interface{}, input io.Reader, cfg *load.Config, api *load.API) {
	load.Logrus.WithFields(logrus.Fields{
		"name": cfg.Name,
	}).Debug("prometheus: running parser")

	mfChan := make(chan *dto.MetricFamily, 1024)
	go func() {
		if err := ParseReader(input, mfChan); err != nil {
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"name": cfg.Name,
					"err":  err,
				}).Error("prometheus: parsing failure")
			}
		}
	}()

	if (*cfg).MetricAPI {
		prometheusMetricAPI(api, &mfChan, cfg.Name)
	} else {
		prometheusStandard(api, &mfChan, dataStore, cfg.Name)
	}
}

func prometheusStandard(api *load.API, mfChan *chan *dto.MetricFamily, dataStore *[]interface{}, cfgName string) {
	load.Logrus.WithFields(logrus.Fields{
		"name": cfgName,
	}).Debug("prometheus: parser generating standard event output")

	// kept for temporary backwards compatibility
	if (*api).Prometheus.Unflatten {
		(*api).Prometheus.Raw = true
	}

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
	for mf := range *mfChan {
		prometheusNewFamily(mf, dataStore, api, &flattenedSample, &sampleKeys)
	}
	// anything sampled add into datastore
	for sample := range sampleKeys {
		*dataStore = append(*dataStore, sampleKeys[sample])
	}
	// add flattened sample into datastore
	if len(flattenedSample) > 0 && !api.Prometheus.Raw {
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

		metric := map[string]interface{}{}
		metric["name"] = dtoMF.GetName()
		metric["help"] = dtoMF.GetHelp()
		metric["type"] = dtoMF.GetType().String()
		applyCustomAttributes(&metric, &api.Prometheus.CustomAttributes)
		prometheusMakeLabels(m, &metric)

		if dtoMF.GetType() == dto.MetricType_SUMMARY {
			if (*api).Prometheus.Raw {
				metric["count"] = fmt.Sprint(m.GetSummary().GetSampleCount())
				metric["sum"] = fmt.Sprint(m.GetSummary().GetSampleSum())
				prometheusMakeQuantiles(m, &metric, dtoMF, api.Prometheus.Raw)
				*dataStore = append(*dataStore, metric)
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
			if len(m.Label) > 0 && !api.Prometheus.Summary && !api.Prometheus.Raw {
				sampleKey := prometheusMakeMergedMeta(sampleKeys, m)
				key := dtoMF.GetName() + ".summary"
				(*sampleKeys)[sampleKey][key+".count"] = fmt.Sprint(m.GetSummary().GetSampleCount())
				(*sampleKeys)[sampleKey][key+".sum"] = fmt.Sprint(m.GetSummary().GetSampleSum())
			}
		} else if dtoMF.GetType() == dto.MetricType_HISTOGRAM {
			if (*api).Prometheus.Raw {
				metric["count"] = fmt.Sprint(m.GetHistogram().GetSampleCount())
				metric["sum"] = fmt.Sprint(m.GetHistogram().GetSampleSum())
				prometheusMakeBuckets(m, &metric, dtoMF, api.Prometheus.Raw)
				*dataStore = append(*dataStore, metric)
			} else if api.Prometheus.Histogram {
				metric["count"] = fmt.Sprint(m.GetHistogram().GetSampleCount())
				metric["sum"] = fmt.Sprint(m.GetHistogram().GetSampleSum())
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
			if len(m.Label) > 0 && !api.Prometheus.Histogram && !api.Prometheus.Raw {
				sampleKey := prometheusMakeMergedMeta(sampleKeys, m)
				key := dtoMF.GetName() + ".histogram"
				(*sampleKeys)[sampleKey][key+".count"] = fmt.Sprint(m.GetHistogram().GetSampleCount())
				(*sampleKeys)[sampleKey][key+".sum"] = fmt.Sprint(m.GetHistogram().GetSampleSum())
			}
		} else { // gauge or counter
			metric["value"] = fmt.Sprint(getValue(m))

			if (*api).Prometheus.Raw {
				*dataStore = append(*dataStore, metric)
			} else if len(m.Label) > 0 {
				sampleKey := prometheusMakeMergedMeta(sampleKeys, m)
				key := dtoMF.GetName()
				if dtoMF.GetType() == dto.MetricType_GAUGE {
					key += ".gauge"
				} else if dtoMF.GetType() == dto.MetricType_COUNTER {
					key += ".counter"
				}
				(*sampleKeys)[sampleKey][key] = fmt.Sprint(getValue(m))
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
}

func prometheusMetricAPI(api *load.API, mfChan *chan *dto.MetricFamily, cfgName string) {
	load.Logrus.WithFields(logrus.Fields{
		"name": cfgName,
	}).Debug("prometheus: parser generating standard event output")

	for mf := range *mfChan {
		for _, m := range mf.Metric {

			attributes := map[string]interface{}{"help": mf.GetHelp()}
			applyCustomAttributes(&attributes, &api.Prometheus.CustomAttributes)
			prometheusMakeLabels(m, &attributes)

			if mf.GetType() == dto.MetricType_SUMMARY {
				attributes["prometheusType"] = "summary"
				summaryMetrics := []map[string]interface{}{}
				for _, q := range m.GetSummary().Quantile {
					if !math.IsNaN(q.GetValue()) {
						quantileMetric := map[string]interface{}{
							"name":  mf.GetName() + "_quantile",
							"type":  "gauge",
							"value": q.GetValue(),
							"attributes": map[string]interface{}{
								"quantile": q.GetQuantile(),
							},
						}
						summaryMetrics = append(summaryMetrics, quantileMetric)
					}
				}
				summaryMetrics = append(summaryMetrics, map[string]interface{}{
					"name":  mf.GetName() + "_count",
					"type":  "gauge",
					"value": m.GetSummary().GetSampleCount(),
				})
				summaryMetrics = append(summaryMetrics, map[string]interface{}{
					"name":  mf.GetName() + "_sum",
					"type":  "gauge",
					"value": m.GetSummary().GetSampleSum(),
				})
				Metrics := load.Metrics{
					TimestampMs:      time.Now().UnixNano() / 1e+6,
					CommonAttributes: attributes,
					Metrics:          summaryMetrics,
				}
				load.MetricsStoreAppend(Metrics)
			} else if mf.GetType() == dto.MetricType_HISTOGRAM {
				attributes["prometheusType"] = "histogram"
				histogramMetrics := []map[string]interface{}{}
				for _, b := range m.GetHistogram().Bucket {
					bucketAttributes := map[string]interface{}{}
					bucketVal := fmt.Sprint(b.GetUpperBound())
					bucketParsedVal, err := strconv.ParseFloat(bucketVal, 64)
					if err != nil || bucketVal == "+Inf" {
						bucketAttributes["le"] = bucketVal
					} else {
						bucketAttributes["le"] = bucketParsedVal
					}
					bucketMetric := map[string]interface{}{
						"name":       mf.GetName() + "_bucket",
						"type":       "gauge",
						"value":      b.GetCumulativeCount(),
						"attributes": bucketAttributes,
					}
					histogramMetrics = append(histogramMetrics, bucketMetric)
				}
				histogramMetrics = append(histogramMetrics, map[string]interface{}{
					"name":  mf.GetName() + "_count",
					"type":  "gauge",
					"value": m.GetHistogram().GetSampleCount(),
				})
				histogramMetrics = append(histogramMetrics, map[string]interface{}{
					"name":  mf.GetName() + "_sum",
					"type":  "gauge",
					"value": m.GetHistogram().GetSampleSum(),
				})
				Metrics := load.Metrics{
					TimestampMs:      time.Now().UnixNano() / 1e+6,
					CommonAttributes: attributes,
					Metrics:          histogramMetrics,
				}
				load.MetricsStoreAppend(Metrics)
			} else if mf.GetType() == dto.MetricType_GAUGE {
				attributes["prometheusType"] = "gauge"
				Metrics := load.Metrics{
					TimestampMs:      time.Now().UnixNano() / 1e+6,
					CommonAttributes: attributes,
					Metrics: []map[string]interface{}{
						{
							"name":  mf.GetName(),
							"type":  "gauge",
							"value": getValue(m),
						}},
				}
				load.MetricsStoreAppend(Metrics)
			} else if mf.GetType() == dto.MetricType_COUNTER {
				attributes["prometheusType"] = "counter"
				Metrics := load.Metrics{
					TimestampMs:      time.Now().UnixNano() / 1e+6,
					CommonAttributes: attributes,
					Metrics: []map[string]interface{}{
						{
							"name":  mf.GetName(),
							"type":  "gauge",
							"value": getValue(m),
						}},
				}
				load.MetricsStoreAppend(Metrics)
			}
		}
	}
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

func prometheusMakeMergedMeta(sampleKeys *map[string]map[string]interface{}, m *dto.Metric) string {
	sampleKey := ""
	sample := map[string]interface{}{}
	for _, label := range m.Label {
		sample[label.GetName()] = label.GetValue()
		sampleKey += label.GetValue()
	}
	if (*sampleKeys)[sampleKey] == nil {
		(*sampleKeys)[sampleKey] = sample
	}
	return sampleKey
}

// applyCustomAttributes applies custom attributes to the provided sample
func applyCustomAttributes(sample *map[string]interface{}, attributes *map[string]string) {
	for key, val := range *attributes {
		(*sample)[key] = val
	}
}

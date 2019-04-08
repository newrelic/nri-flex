package outputs

import (
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// CreateStatusSample creates flexStatusSample
func CreateStatusSample() {
	flexStatusSample := load.Entity.NewMetricSet("flexStatusSample")
	logger.Flex("debug", flexStatusSample.SetMetric("eventCount", load.EventCount, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("eventDropCount", load.EventDropCount, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("configsProcessed", load.ConfigsProcessed, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("metricApiGauges", load.GaugeMetrics, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("metricApiCounts", load.CounterMetrics, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("metricApiSummary", load.SummaryMetrics, metric.GAUGE), "", false)
	for sample, count := range load.EventDistribution {
		logger.Flex("debug", flexStatusSample.SetMetric(sample+"_count", count, metric.GAUGE), "", false)
	}
}

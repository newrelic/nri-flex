package outputs

import (
	"nri-flex/internal/load"
	"nri-flex/internal/logger"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// CreateStatusSample creates flexStatusSample
func CreateStatusSample() {
	flexStatusSample := load.Entity.NewMetricSet("flexStatusSample")
	logger.Flex("debug", flexStatusSample.SetMetric("eventCount", load.EventCount, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("eventDropCount", load.EventDropCount, metric.GAUGE), "", false)
	logger.Flex("debug", flexStatusSample.SetMetric("configsProcessed", load.ConfigsProcessed, metric.GAUGE), "", false)
	for sample, count := range load.EventDistribution {
		logger.Flex("debug", flexStatusSample.SetMetric(sample+"_count", count, metric.GAUGE), "", false)
	}
}

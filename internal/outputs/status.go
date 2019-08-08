package outputs

import (
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// StatusSample creates flexStatusSample
func StatusSample() {
	flexStatusSample := load.Entity.NewMetricSet("flexStatusSample")
	logger.Flex("error", flexStatusSample.SetMetric("flex.IntegrationVersion", load.IntegrationVersion, metric.ATTRIBUTE), "", false)
	if load.Args.GitRepo != "" {
		logger.Flex("error", flexStatusSample.SetMetric("flex.GitRepo", load.Args.GitRepo, metric.ATTRIBUTE), "", false)
		if load.Args.GitBranch != "" && load.Args.GitCommit == "" {
			logger.Flex("error", flexStatusSample.SetMetric("flex.GitBranch", load.Args.GitBranch, metric.ATTRIBUTE), "", false)
		}
	}
	if load.Hostname != "" {
		logger.Flex("error", flexStatusSample.SetMetric("flex.Hostname", load.Hostname, metric.ATTRIBUTE), "", false)
	}
	if load.ContainerID != "" {
		logger.Flex("error", flexStatusSample.SetMetric("flex.IsContainer", "true", metric.ATTRIBUTE), "", false)
	}
	if load.IsKubernetes {
		logger.Flex("error", flexStatusSample.SetMetric("flex.IsKubernetes", "true", metric.ATTRIBUTE), "", false)
	}
	if load.IsFargate {
		logger.Flex("error", flexStatusSample.SetMetric("flex.IsFargate", "true", metric.ATTRIBUTE), "", false)
	}
	if load.LambdaName != "" {
		logger.Flex("error", flexStatusSample.SetMetric("flex.LambdaName", load.LambdaName, metric.ATTRIBUTE), "", false)
	}
	if load.AWSExecutionEnv != "" {
		logger.Flex("error", flexStatusSample.SetMetric("flex.AWSExecutionEnv", load.AWSExecutionEnv, metric.ATTRIBUTE), "", false)
	}
	for counter, value := range load.FlexStatusCounter.M {
		logger.Flex("error", flexStatusSample.SetMetric("flex.counter."+counter, value, metric.GAUGE), "", false)
	}
	for pid, val := range load.DiscoveredProcesses {
		logger.Flex("error", flexStatusSample.SetMetric("flex.pd."+pid, val, metric.ATTRIBUTE), "", false)
	}
}

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-flex/internal/load"
)

// StatusSample creates flexStatusSample
func StatusSample() {
	flexStatusSample := load.Entity.NewMetricSet("flexStatusSample")
	endTimeNs := load.MakeTimestamp()
	statusLog(flexStatusSample.SetMetric("flex.time.endMs", endTimeNs, metric.GAUGE))
	statusLog(flexStatusSample.SetMetric("flex.time.startMs", load.StartTime, metric.GAUGE))
	statusLog(flexStatusSample.SetMetric("flex.time.elaspedMs", endTimeNs-load.StartTime, metric.GAUGE))

	statusLog(flexStatusSample.SetMetric("flex.IntegrationVersion", load.IntegrationVersion, metric.ATTRIBUTE))
	if load.Args.GitRepo != "" {
		statusLog(flexStatusSample.SetMetric("flex.GitRepo", load.Args.GitRepo, metric.ATTRIBUTE))
		if load.Args.GitBranch != "" && load.Args.GitCommit == "" {
			statusLog(flexStatusSample.SetMetric("flex.GitBranch", load.Args.GitBranch, metric.ATTRIBUTE))
		}
	}
	if load.Hostname != "" {
		statusLog(flexStatusSample.SetMetric("flex.Hostname", load.Hostname, metric.ATTRIBUTE))
	}
	if load.ContainerID != "" {
		statusLog(flexStatusSample.SetMetric("flex.ContainerId", load.ContainerID, metric.ATTRIBUTE))
	}
	if load.IsKubernetes {
		statusLog(flexStatusSample.SetMetric("flex.IsKubernetes", "true", metric.ATTRIBUTE))
	}
	if load.IsFargate {
		statusLog(flexStatusSample.SetMetric("flex.IsFargate", "true", metric.ATTRIBUTE))
	}
	if load.LambdaName != "" {
		statusLog(flexStatusSample.SetMetric("flex.LambdaName", load.LambdaName, metric.ATTRIBUTE))
	}
	if load.AWSExecutionEnv != "" {
		statusLog(flexStatusSample.SetMetric("flex.AWSExecutionEnv", load.AWSExecutionEnv, metric.ATTRIBUTE))
	}
	for counter, value := range load.FlexStatusCounter.M {
		statusLog(flexStatusSample.SetMetric("flex.counter."+counter, value, metric.GAUGE))
	}
	for pid, val := range load.DiscoveredProcesses {
		statusLog(flexStatusSample.SetMetric("flex.pd."+pid, val, metric.ATTRIBUTE))
	}
}

func statusLog(err error) {
	if err != nil {
		load.Logrus.WithError(err).Error("status: failed to set metric")
	}
}

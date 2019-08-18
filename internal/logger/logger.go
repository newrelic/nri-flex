package logger

import (
	"fmt"

	"github.com/newrelic/nri-flex/internal/load"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// Flex generic log handler to support force logging and creating additional events for insights debugging
func Flex(logType string, err error, message interface{}, createEvent bool) {
	if createEvent || load.Args.ForceLogEvent {
		flexEvent := "flexDebug"
		if logType == "fatal" {
			flexEvent = "flexFatal"
		} else if logType == "info" {
			flexEvent = "flexInfo"
		} else if logType == "error" {
			flexEvent = "flexError"
		}
		metricSet := load.Entity.NewMetricSet(flexEvent)
		msErr := metricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE)
		if msErr != nil {
			load.Logrus.Debug(msErr)
		}
		msErr = metricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE)
		if msErr != nil {
			load.Logrus.Debug(msErr)
		}
		if err != nil {
			msErr := metricSet.SetMetric("error", err.Error(), metric.ATTRIBUTE)
			if msErr != nil {
				load.Logrus.Debug(msErr)
			}
		}
		if message != "" {
			msErr := metricSet.SetMetric("message", message, metric.ATTRIBUTE)
			if msErr != nil {
				load.Logrus.Debug(msErr)
			}
		}
	}

	switch logType {
	case "fatal":
		if err != nil {
			if message != "" {
				load.Logrus.Debug(message)
			}
			load.Logrus.Fatal(err)
		}
	case "error":
		if err != nil {
			completeMsg := err.Error()
			if message != "" {
				completeMsg += " - " + fmt.Sprintf("%v", message)
			}
			load.Logrus.Error(completeMsg)
		}
	case "debug":
		if err != nil {
			completeMsg := err.Error()
			if message != "" {
				completeMsg += " - " + fmt.Sprintf("%v", message)
			}
			load.Logrus.Debug(completeMsg)
		} else {
			load.Logrus.Debug(message)
		}
	case "info":
		if err != nil {
			load.Logrus.Info(err)
		}
		if message != "" {
			load.Logrus.Info(message)
		}
	}
}

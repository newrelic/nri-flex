package logger

import (
	"nri-flex/internal/load"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// Flex generic log handler to support force logging and creating additional events for insights debugging
func Flex(logType string, err error, message string, createEvent bool) {
	log.SetupLogging(load.Args.Verbose)

	if createEvent || load.Args.ForceLogEvent {
		flexEvent := "flexDebug"
		if logType == "fatal" {
			flexEvent = "flexFatal"
		} else if logType == "info" {
			flexEvent = "flexInfo"
		}
		metricSet := load.Entity.NewMetricSet(flexEvent)
		msErr := metricSet.SetMetric("integration_version", load.IntegrationVersion, metric.ATTRIBUTE)
		if msErr != nil {
			log.Debug(msErr.Error())
		}
		msErr = metricSet.SetMetric("integration_name", load.IntegrationName, metric.ATTRIBUTE)
		if msErr != nil {
			log.Debug(msErr.Error())
		}
		if err != nil {
			msErr := metricSet.SetMetric("error", err.Error(), metric.ATTRIBUTE)
			if msErr != nil {
				log.Debug(msErr.Error())
			}
		}
		if message != "" {
			msErr := metricSet.SetMetric("message", message, metric.ATTRIBUTE)
			if msErr != nil {
				log.Debug(msErr.Error())
			}
		}
	}

	switch logType {
	case "fatal":
		if err != nil {
			if message != "" {
				log.Debug(message)
			}
			log.Fatal(err)
		}
	case "debug":
		if err != nil {
			completeMsg := err.Error()
			if message != "" {
				completeMsg += " - " + message
			}
			log.Debug(completeMsg)
		}
	case "info":
		if err != nil {
			log.Info(err.Error())
		}
		if message != "" {
			log.Debug(message)
		}
	}
}

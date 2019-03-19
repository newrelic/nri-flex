package outputs

import (
	"nri-flex/internal/load"
	"nri-flex/internal/logger"

	Integration "github.com/newrelic/infra-integrations-sdk/integration"
)

// CreateIntegration Creates Infrastructure SDK Integration
func CreateIntegration() {
	var err error
	load.Integration, err = Integration.New(load.IntegrationName, load.IntegrationVersion, Integration.Args(&load.Args))
	logger.Flex("fatal", err, "", false)

	if load.Args.Local {
		load.Entity = load.Integration.LocalEntity()
	} else if load.Args.Entity != "" {
		load.Entity, err = load.Integration.Entity(load.Args.Entity, "nri-flex")
		logger.Flex("fatal", err, "", false)
	} else {
		load.Entity, err = load.Integration.Entity(load.Hostname, "nri-flex")
		logger.Flex("fatal", err, "", false)
	}

}

// CreateRemote Creates Infrastructure Remote Entity
func CreateRemote(host string) {
	var err error
	setHost := load.Hostname
	if host != "" {
		setHost = host
	}
	load.Entity, err = load.Integration.Entity(setHost, "nri-flex")
	logger.Flex("fatal", err, "", false)
}

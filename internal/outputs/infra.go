package outputs

import (
	"os"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	Integration "github.com/newrelic/infra-integrations-sdk/integration"
)

// InfraIntegration Creates Infrastructure SDK Integration
func InfraIntegration() {
	load.Hostname, _ = os.Hostname() // set hostname

	var err error
	load.Integration, err = Integration.New(load.IntegrationName, load.IntegrationVersion, Integration.Args(&load.Args))
	logger.Flex("fatal", err, "", false)

	if load.Args.Local {
		load.Entity = load.Integration.LocalEntity()
	} else {
		InfraRemoteEntity()
	}
}

// InfraRemoteEntity Creates Infrastructure Remote Entity
func InfraRemoteEntity() {
	var err error
	setEntity := load.Hostname // default hostname
	if load.Args.Entity != "" {
		setEntity = load.Args.Entity
	}
	load.Entity, err = load.Integration.Entity(setEntity, "nri-flex")
	logger.Flex("fatal", err, "", false)
}

package outputs

import (
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
	"os"

	Integration "github.com/newrelic/infra-integrations-sdk/integration"
)

// CreateIntegration Creates Infrastructure SDK Integration
func CreateIntegration() {
	load.Hostname, _ = os.Hostname() // set hostname

	var err error
	load.Integration, err = Integration.New(load.IntegrationName, load.IntegrationVersion, Integration.Args(&load.Args))
	logger.Flex("fatal", err, "", false)

	if load.Args.Local {
		load.Entity = load.Integration.LocalEntity()
	} else {
		CreateRemoteEntity()
	}
}

// CreateRemoteEntity Creates Infrastructure Remote Entity
func CreateRemoteEntity() {
	var err error
	setEntity := load.Hostname // default hostname
	if load.Args.Entity != "" {
		setEntity = load.Args.Entity
	}
	load.Entity, err = load.Integration.Entity(setEntity, "nri-flex")
	logger.Flex("fatal", err, "", false)
}

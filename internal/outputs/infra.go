/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"fmt"
	"os"
	"sync"

	Integration "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
)

// InfraIntegration Creates Infrastructure SDK Integration
var initInfraOnce sync.Once

func InfraIntegration() error {
	var err error
	load.Hostname, err = os.Hostname() // set hostname
	if err != nil {
		load.Logrus.
			WithError(err).
			Debug("flex: failed to get the hostname while creating integration")
	}

	// Do this only once so flag doesn't throw a redefined panic
	initInfraOnce.Do(func() {
		load.Integration, err = Integration.New(load.IntegrationName, load.IntegrationVersion, Integration.Args(&load.Args))
	})

	if err != nil {
		return fmt.Errorf("flex: failed to create integration %v", err)
	}

	// Accepts ConfigPath as alias for ConfigFile. This will allow the Infrastructure Agent
	// passing an embedded config via the default CONFIG_PATH environment variable
	if load.Args.ConfigPath != "" {
		load.Args.ConfigFile = load.Args.ConfigPath
	}

	load.Entity, err = createEntity(load.Args.Local, load.Args.Entity)
	if err != nil {
		return fmt.Errorf("flex: failed create entity: %v", err)
	}
	return nil
}

func createEntity(isLocalEntity bool, entityName string) (*Integration.Entity, error) {
	if isLocalEntity {
		return load.Integration.LocalEntity(), nil
	}

	if entityName == "" {
		entityName = load.Hostname // default hostname
	}

	return load.Integration.Entity(entityName, "nri-flex")
}

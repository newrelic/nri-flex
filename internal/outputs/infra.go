/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"os"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"

	Integration "github.com/newrelic/infra-integrations-sdk/integration"
)

// InfraIntegration Creates Infrastructure SDK Integration
func InfraIntegration() {
	load.Hostname, _ = os.Hostname() // set hostname

	var err error
	load.Integration, err = Integration.New(load.IntegrationName, load.IntegrationVersion, Integration.Args(&load.Args))

	if err != nil {
		load.Logrus.WithFields(logrus.Fields{"err": err}).Fatal("flex: create integration")
	}

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
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{"err": err}).Fatal("flex: create remote entity")
	}
}

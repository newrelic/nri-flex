/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"fmt"
	"os"
	"time"

	Integration "github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/persist"
	"github.com/newrelic/nri-flex/internal/load"
)

// InfraIntegration Creates Infrastructure SDK Integration
func InfraIntegration() error {
	var err error
	load.Hostname, err = os.Hostname() // set hostname
	if err != nil {
		return fmt.Errorf("flex: failed to get the hostname while creating integration")
	}

	storer, err := createStorer()
	if err != nil {
		return fmt.Errorf("can't create custom store: %s", err)
	}

	load.Integration, err = Integration.New(load.IntegrationName, load.IntegrationVersion, Integration.Args(&load.Args), Integration.Storer(storer))
	if err != nil {
		return fmt.Errorf("flex: failed to create integration %v", err)
	}

	// when arguments have been set re-run logger setup to check verbose flag
	load.SetupLogger()

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

// create custom storer with custom STORER_NAME and STORER_TTL
func createStorer() (persist.Storer, error) {
	storerAttributes := os.Getenv("STORER_NAME")
	storerName := load.IntegrationName + storerAttributes

	ttl := persist.DefaultTTL
	// storerTTL, err := strconv.Atoi(os.Getenv("STORER_TTL"))
	storerTTL, err := time.ParseDuration(os.Getenv("STORER_TTL"))
	if err == nil && storerTTL > 0 {
		ttl = storerTTL
	}

	if os.Getenv("VERBOSE") == "true" || os.Getenv("VERBOSE") == "1" {
		fmt.Fprintf(os.Stderr, "custom storer name: %s and TTL: %s\n", storerName, ttl)
	}
	return persist.NewFileStore(persist.DefaultPath(storerName), load.Logrus, ttl)
}

// +build integration

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"runtime"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/assert"
)

func TestDrivers(t *testing.T) {
	drivers := map[string]string{
		"postgres":  load.DefaultPostgres,
		"pg":        load.DefaultPostgres,
		"pq":        load.DefaultPostgres,
		"mssql":     load.DefaultMSSQLServer,
		"sqlserver": load.DefaultMSSQLServer,
		"mysql":     load.DefaultMySQL,
		"mariadb":   load.DefaultMySQL,
		"hana":      load.DefaultHANA,
		"db2":       load.DefaultDB2,
		"ibm_db2":   load.DefaultDB2,
		"ora":       load.DefaultOracle,
		"oracle":    load.DefaultOracle,
		"godror":    load.DefaultOracle,
		"ase":       load.DefaultSybase,
		"sybase":    load.DefaultSybase,
		"unknown":   "",
	}

	// test switch
	for db, driver := range drivers {
		detectedDriver := setDatabaseDriver(db, "")
		if detectedDriver != driver {
			t.Errorf("expected %v got %v", driver, detectedDriver)
		}
	}

	// test manual driver
	detectedDriver := setDatabaseDriver("", "superdb")
	if detectedDriver != "superdb" {
		t.Errorf("expected superdb got %v", detectedDriver)
	}
}

func TestDatabase(t *testing.T) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		t.Skip("Darwin not supported yet")
	}

	load.Refresh()
	config := load.Config{
		Name: "postgresDbFlex",
		APIs: []load.API{
			{
				Name:     "postgres",
				Database: "postgres",
				DBConn:   "user=postgres host=postgres-db sslmode=disable password=flex port=5432",
				CustomAttributes: map[string]string{
					"parentAttr": "myDbServer",
				},
				DBQueries: []load.Command{
					{
						Name: "pgStatActivitySample",
						Run:  "select * FROM pg_stat_activity LIMIT 2",
						CustomAttributes: map[string]string{
							"nestedAttr": "nestedVal",
						},
					},
				},
			},
			{
				Name:     "postgres",
				Database: "pg",
				DBConn:   "user=postgres host=postgres-db sslmode=disable password=flex port=5432",
				CustomAttributes: map[string]string{
					"parentAttr": "myDbServer",
				},
				DBQueries: []load.Command{
					{
						Name: "pgStatActivitySample",
						Run:  "select * FROM pg_stat_activity LIMIT 2",
						CustomAttributes: map[string]string{
							"nestedAttr": "nestedVal",
						},
					},
				},
			},
			{
				Name:     "postgres",
				Database: "pq",
				DBConn:   "user=postgres host=postgres-db sslmode=disable password=flex port=5433",
				DBQueries: []load.Command{
					{
						Name: "pgStatActivitySample",
						Run:  "select * FROM pg_stat_activity LIMIT 2",
					},
				},
			},
		},
	}

	dataStore := []interface{}{}
	ProcessQueries(&dataStore, &config, 0)
	ProcessQueries(&dataStore, &config, 1)
	ProcessQueries(&dataStore, &config, 2)

	assert.Lenf(t, dataStore, 4, "expected 4 samples, got %d", len(dataStore))

	sampleName := dataStore[0].(map[string]interface{})["queryLabel"]
	assert.Equalf(t, "pgStatActivitySample", sampleName,
		"expected label %v, got %v", "pgStatActivitySample", sampleName)
}

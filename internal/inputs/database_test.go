// +build integration

package inputs

import (
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
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
	load.Refresh()
	config := load.Config{
		Name: "postgresDbFlex",
		APIs: []load.API{
			{
				Name:     "postgres",
				Database: "postgres",
				DbConn:   "user=postgres host=postgres-db sslmode=disable password=flex port=5432",
				CustomAttributes: map[string]string{
					"parentAttr": "myDbServer",
				},
				DbQueries: []load.Command{
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
				DbConn:   "user=postgres host=postgres-db sslmode=disable password=flex port=5432",
				CustomAttributes: map[string]string{
					"parentAttr": "myDbServer",
				},
				DbQueries: []load.Command{
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
				DbConn:   "user=postgres host=postgres-db sslmode=disable password=flex port=5433",
				DbQueries: []load.Command{
					{
						Name: "pgStatActivitySample",
						Run:  "select * FROM pg_stat_activity LIMIT 2",
					},
				},
			},
		},
	}

	ProcessQueries(config.APIs[0])
	ProcessQueries(config.APIs[1])
	ProcessQueries(config.APIs[2])

	if len(load.Store.Data) != 4 {
		t.Errorf("expected 4 samples, got %d", len(load.Store.Data))
	} else {
		if load.Store.Data[0].(map[string]interface{})["queryLabel"] != "pgStatActivitySample" {
			t.Errorf("incorrect label %v", load.Store.Data[0].(map[string]interface{})["queryLabel"])
		}
	}
}

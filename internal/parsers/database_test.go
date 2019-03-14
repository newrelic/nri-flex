package parser

import (
	"nri-flex/internal/load"
	"testing"
)

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

	dataStore := []interface{}{}
	ProcessQueries(config.APIs[0], &dataStore)
	ProcessQueries(config.APIs[1], &dataStore)
	ProcessQueries(config.APIs[2], &dataStore)

	if len(dataStore) != 4 {
		t.Errorf("expected 4 samples, got %d", len(dataStore))
	}

	if dataStore[0].(map[string]interface{})["queryLabel"] != "pgStatActivitySample" {
		t.Errorf("incorrect label %v", dataStore[0].(map[string]interface{})["queryLabel"])
	}
}

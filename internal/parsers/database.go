package parser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
	"strconv"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"

	//Database Drivers
	_ "github.com/denisenkom/go-mssqldb" //mssql | sql-server
	_ "github.com/go-sql-driver/mysql"   //mysql
	_ "github.com/lib/pq"                //postgres
	//
)

// ProcessQueries processes database queries
func ProcessQueries(api load.API, dataStore *[]interface{}) {
	logger.Flex("debug", fmt.Errorf("running %v queries", api.Database), "", false)

	//sql.Open doesn't open the connection, use a generic Ping() to test the connection
	db, err := sql.Open(setDatabaseDriver(api.Database, api.DbDriver), api.DbConn)

	// commenting out as db.Ping is not reliable currently
	// https://stackoverflow.com/questions/41618428/golang-ping-succeed-the-second-time-even-if-database-is-down/41619206#41619206
	var pingError error
	if db != nil {
		dbPingWithTimeout(db, &pingError)
	}

	if err != nil || pingError != nil {
		if err != nil {
			logger.Flex("debug", err, "", false)
			if api.Logging.Open {
				errorLogToInsights(err, api.Database, api.Name, "")
			}
		}
		if pingError != nil {
			logger.Flex("debug", pingError, "ping error", false)
			if api.Logging.Open {
				errorLogToInsights(pingError, api.Database, api.Name, "")
			}
		}
	} else {
		for _, query := range api.DbQueries {
			if query.Name == "" {
				logger.Flex("debug", fmt.Errorf("missing name for: %v", query.Run), "", false)
				break
			}
			if query.Run == "" {
				logger.Flex("debug", fmt.Errorf("query ('run') parameter not defined"), "", false)
				break
			}

			rows, err := db.Query(query.Run)
			if err != nil {
				logger.Flex("debug", err, "query: "+query.Run, false)
				errorLogToInsights(err, api.Database, api.Name, query.Name)
			} else {
				logger.Flex("info", nil, fmt.Sprintf("running query: %v", query.Run), false)

				cols, err := rows.Columns()
				if err != nil {
					logger.Flex("debug", err, "", false)
					errorLogToInsights(err, api.Database, api.Name, query.Name)
				} else {
					values := make([]sql.RawBytes, len(cols))
					scanArgs := make([]interface{}, len(values))
					for i := range values {
						scanArgs[i] = &values[i]
					}

					// Fetch rows
					rowNo := 1
					for rows.Next() {
						rowSet := map[string]interface{}{
							"rowIdentifier": query.Name + "_" + strconv.Itoa(rowNo),
							"queryLabel":    query.Name,
							"event_type":    query.Name,
						}
						// apply event type override if set (this is useful to set if needing to group multiples under one event type)
						if query.EventType != "" {
							rowSet["event_type"] = query.EventType
						}

						// get RawBytes
						err = rows.Scan(scanArgs...)
						if err != nil {
							logger.Flex("debug", err, "", false)
						} else {
							// Loop through each column
							for i, col := range values {
								// If value nil == null
								if col == nil {
									rowSet[cols[i]] = "NULL"
								} else {
									rowSet[cols[i]] = string(col)
								}
							}
							*dataStore = append(*dataStore, rowSet)
							rowNo++
						}
					}
				}
			}
		}
	}
}

// setDatabaseDriver returns driver if set, otherwise sets a default driver based on database
func setDatabaseDriver(database, driver string) string {
	if driver != "" {
		return driver
	}
	switch database {
	case "postgres":
		return load.DefaultPostgres
	case "pg":
		return load.DefaultPostgres
	case "pq":
		return load.DefaultPostgres
	case "mssql":
		return load.DefaultMSSQLServer
	case "sqlserver":
		return load.DefaultMSSQLServer
	case "mysql":
		return load.DefaultMySQL
	case "mariadb":
		return load.DefaultMySQL
	}
	return ""
}

// errorLogToInsights log errors to insights, useful to debug
func errorLogToInsights(err error, database, name, queryLabel string) {
	errorMetricSet := load.Entity.NewMetricSet(database + "Error")
	load.EventDistribution[database+"Error"]++
	load.EventCount++

	logger.Flex("debug", errorMetricSet.SetMetric("errorMsg", err.Error(), metric.ATTRIBUTE), "", false)
	if name != "" {
		logger.Flex("debug", errorMetricSet.SetMetric("name", name, metric.ATTRIBUTE), "", false)
	}
	if queryLabel != "" {
		logger.Flex("debug", errorMetricSet.SetMetric("queryLabel", queryLabel, metric.ATTRIBUTE), "", false)
	}
}

// dbPingWithTimeout Database Ping() with Timeout
func dbPingWithTimeout(db *sql.DB, pingError *error) {
	ctx := context.Background()

	// Create a channel for signal handling
	c := make(chan struct{})
	// Define a cancellation after 1s in the context
	ctx, cancel := context.WithTimeout(ctx, load.DefaultPingTimeout*time.Millisecond)
	defer cancel()

	// Run ping via a goroutine
	go func() {
		pingWrapper(db, c, pingError)
	}()

	// Listen for signals
	select {
	case <-ctx.Done():
		*pingError = errors.New("Ping failed: " + ctx.Err().Error())
	case <-c:
		logger.Flex("info", nil, "db.Ping finished", false)
	}
}

func pingWrapper(db *sql.DB, c chan struct{}, pingError *error) {
	*pingError = db.Ping()
	c <- struct{}{}
}

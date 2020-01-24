/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"

	"github.com/newrelic/infra-integrations-sdk/data/metric"

	//Database Drivers
	_ "github.com/SAP/go-hdb/driver"      //SAP HANA
	_ "github.com/denisenkom/go-mssqldb"  //mssql | sql-server
	_ "github.com/go-sql-driver/mysql"    //mysql
	_ "github.com/lib/pq"                 //postgres
	_ "github.com/vertica/vertica-sql-go" //HP Vertica
	//
)

// ProcessQueries processes database queries
func ProcessQueries(dataStore *[]interface{}, yml *load.Config, apiNo int) {
	api := yml.APIs[apiNo]

	load.Logrus.WithFields(logrus.Fields{
		"name":     yml.Name,
		"database": api.Database,
	}).Debug("database: finding flex container id")

	// sql.Open doesn't open the connection, use a generic Ping() to test the connection
	db, err := sql.Open(setDatabaseDriver(api.Database, api.DbDriver), api.DbConn)

	// wrapping dbPingWithTimeout out as db.Ping is not reliable currently
	// https://stackoverflow.com/questions/41618428/golang-ping-succeed-the-second-time-even-if-database-is-down/41619206#41619206
	var pingError error
	if db != nil {
		dbPingWithTimeout(db, &pingError)
	}

	if err != nil || pingError != nil {
		if err != nil {

			load.Logrus.WithFields(logrus.Fields{
				"err":      err,
				"name":     yml.Name,
				"database": api.Database,
			}).Debug("database: unable to connect")

			if api.Logging.Open {
				errorLogToInsights(err, api.Database, api.Name, "")
			}
		}
		if pingError != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err":      err,
				"name":     yml.Name,
				"database": api.Database,
			}).Debug("database: ping error")

			if api.Logging.Open {
				errorLogToInsights(pingError, api.Database, api.Name, "")
			}
		}
	} else {
		// execute queries async else do synchronously
		if api.DbAsync {
			var wg sync.WaitGroup
			wg.Add(len(api.DbQueries))
			for _, query := range api.DbQueries {
				go func(query load.Command) {
					defer wg.Done()
					if query.Name == "" || query.Run == "" {
						if query.Name == "" {
							load.Logrus.WithFields(logrus.Fields{"query": query.Run}).Error("database: query missing name")
						}
						if query.Run == "" {
							load.Logrus.WithFields(logrus.Fields{"query": query.Run, "name": yml.Name, "database": api.Database}).Error("database: run parameter not defined")
						}
					} else {
						runQuery(db, query, api, yml, dataStore)
					}
				}(query)
			}
			wg.Wait()
		} else {
			for _, query := range api.DbQueries {
				if query.Name == "" {
					load.Logrus.WithFields(logrus.Fields{"query": query.Run}).Error("database: query missing name")
					break
				}
				if query.Run == "" {
					load.Logrus.WithFields(logrus.Fields{"query": query.Run, "name": yml.Name, "database": api.Database}).Error("database: run parameter not defined")
					break
				}
				runQuery(db, query, api, yml, dataStore)
			}
		}
	}
}

func runQuery(db *sql.DB, query load.Command, api load.API, yml *load.Config, dataStore *[]interface{}) {
	queryStartTime := load.TimestampMs()
	rows, err := db.Query(query.Run)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"query":    query.Run,
			"name":     yml.Name,
			"database": api.Database,
		}).Error("database: query failed")

		errorLogToInsights(err, api.Database, api.Name, query.Name)
	} else {

		load.Logrus.WithFields(logrus.Fields{
			"configName": yml.Name,
			"database":   api.Database,
			"query":      query.Run,
		}).Debug("database: running query")

		cols, err := rows.Columns()
		if err != nil {

			load.Logrus.WithFields(logrus.Fields{
				"configName": yml.Name,
				"apiName":    api.Name,
				"database":   api.Database,
				"query":      query.Run,
			}).Debug("database: column return failed")

			errorLogToInsights(err, api.Database, api.Name, query.Name)
		} else {
			// Use interface{} type instead of original sql.RawBytes, parsing the value ourselves instead of using sql scan convert routine,
			// which does not hanlde sql.NullString, sql.NullBool, sql.NullFloat64, sql.NullInt64 conversion to sql.RawBytes.
			values := make([]interface{}, len(cols))
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
					load.Logrus.WithFields(logrus.Fields{
						"err": err,
					}).Error("database: row scan failed")
				} else {
					// Loop through each column
					for i, col := range values {
						// If value nil == null
						if col == nil {
							rowSet[cols[i]] = ""
						} else {
							rowSet[cols[i]] = asString(col)
						}
					}
					queryEndTime := load.TimestampMs()
					rowSet["flex.QueryStartMs"] = queryStartTime
					rowSet["flex.QueryTimeMs"] = queryEndTime - queryStartTime
					*dataStore = append(*dataStore, rowSet)
					// load.StoreAppend(rowSet)
					rowNo++
				}
			}
			err = rows.Err()
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"configName": yml.Name,
					"apiName":    api.Name,
					"database":   api.Database,
					"query":      query.Run,
				}).Debug("database: rows return failed")
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
	case "postgres", "pg", "pq":
		return load.DefaultPostgres
	case "mssql", "sqlserver":
		return load.DefaultMSSQLServer
	case "mysql", "mariadb":
		return load.DefaultMySQL
	case "hana", "go-hdb", "hdb":
		return load.DefaultHANA
	case "vertica", "hpvertica":
		return load.DefaultVertica
	}
	return ""
}

// errorLogToInsights log errors to insights, useful to debug
func errorLogToInsights(err error, database, name, queryLabel string) {
	errorMetricSet := load.Entity.NewMetricSet(database + "Error")

	load.StatusCounterIncrement("EventCount")
	load.StatusCounterIncrement(database + "Error")

	set(errorMetricSet.SetMetric("errorMsg", err.Error(), metric.ATTRIBUTE))
	if name != "" {
		set(errorMetricSet.SetMetric("name", name, metric.ATTRIBUTE))
	}
	if queryLabel != "" {
		set(errorMetricSet.SetMetric("queryLabel", queryLabel, metric.ATTRIBUTE))
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
		load.Logrus.Debug("database: db.Ping finished")
	}
}

func pingWrapper(db *sql.DB, c chan struct{}, pingError *error) {
	*pingError = db.Ping()
	c <- struct{}{}
}

func set(err error) {
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{"err": err}).Error("flex: failed to set")
	}
}

func asString(src interface{}) string {

	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case sql.NullString, sql.NullBool, sql.NullFloat64, sql.NullInt64:
		return ""
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())

	}

	return fmt.Sprintf("%v", src)
}

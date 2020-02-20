# Database Queries

**Disclaimer**: this function is bundled in Alpha status. That means that it is not yet supported by New Relic.

Flex has several database drivers available, to help you run any arbitrary/custom queries against those databases.

* https://github.com/denisenkom/go-mssqldb //mssql | sql-server
* https://github.com/go-sql-driver/mysql   //mysql
* https://github.com/lib/pq                //postgres

The below example shows us being able to run multiple queries against one database. 
But also being able to define another database to send queries too. 
It is perfectly fine to use multiple database types in a single config file if you wish to do so.

```
name: postgresDbFlex
apis: 
  - database: postgres
    db_conn: user=postgres host=postgres-db.com sslmode=disable password=flex port=5432
    logging:
      open: true
    custom_attributes: # applies to all queries
      host: myDbServer
    db_queries: 
      - name: pgStatActivitySample
        run: select * FROM pg_stat_activity
        custom_attributes: # can apply additional at a nested level
          nestedAttr: nestedVal
      - name: pgStatAnotherSample
        run: select * FROM some_otherTable
  - database: postgres
    db_conn: user=abc host=myhost.ap-southeast-2.rds.amazonaws.com sslmode=disable password=mypass port=5432 # could be another DB
    queries: 
      - name: pgStatDbSample
        run: select * FROM pg_stat_database LIMIT 1
```
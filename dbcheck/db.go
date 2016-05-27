package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/palette-software/insight-tester/common/logging"
	"time"
)

func getConnectionString(config Database) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", config.User, config.Password, config.Database, config.Host, config.Port)
}

var DBConnection *sql.DB

func getConnection(database Database) *sql.DB {
	if DBConnection == nil {
		connectionString := getConnectionString(database)
		var err error
		DBConnection, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Errorf("Error connecting to db: %v", err)
			return nil
		}
	}
	return DBConnection
}

func closeDB() {
	DBConnection.Close()
}

func check(database Database, test Test) bool {
	db := getConnection(database)
	if db == nil {
		return false
	}
	start := time.Now()
	rows, err := db.Query(test.Sql)
	if err != nil {
		log.Errorf("Error getting rows: %v", err)
		return false
	}
	defer rows.Close()
	rowCount := 0
	for rows.Next() {
		rowCount++
		var count int
		var hostName string
		rows.Scan(&count, &hostName)
		if !checkTest(count, test) {
			log.Errorf("Failed DB check: %v - host: %v - expected: %s - actual: %v - elapsed: %v", test.Description, hostName, fmt.Sprintf("%s%d", test.Result.Operation, test.Result.Count), count, time.Since(start))
			return false
		}
	}
	log.Infof("Successful test: %v - rows: %v - elapsed: %v", test.Description, rowCount, time.Since(start))
	return true
}

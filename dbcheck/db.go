package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/palette-software/insight-tester/common/logging"
	"time"
)

func getConnectionString(host string, config ConstParams) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", config.User, config.Password, config.Database, host, config.Port)
}

var DBConnection *sql.DB

func getConnection(host string, dbParams ConstParams) *sql.DB {
	if DBConnection == nil {
		connectionString := getConnectionString(host, dbParams)
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

func check(host string, dbParams ConstParams, test Test) bool {
	db := getConnection(host, dbParams)
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
			expected := fmt.Sprintf("%s%d", test.Result.Operation, test.Result.Count)
			log.Errorf("FAILED: [HOST:%v] [MACHINE:%v] [TEST:%v] [EXPECTED:%v] [ACTUAL:%v] [DURATION:%v]", host, hostName, test.Description, expected, count, time.Since(start))
			return false
		}
	}
	log.Infof("OK: [HOST:%v] [TEST:%v] [COUNT:%v] [DURATION:%v]", host, test.Description, rowCount, time.Since(start))
	return true
}

package main

import (
	"database/sql"
	"fmt"
	"github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
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
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error connecting to DB:")
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
	rows, err := db.Query(test.Sql)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error getting rows:")
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		rows.Scan(&count)
		if !checkTest(count, test) {
			logrus.WithFields(logrus.Fields{
				"description": test.Description,
				"expected":    fmt.Sprintf("%s%d", test.Result.Operation, test.Result.Count),
				"actual":      count,
			}).Error("Failed DB sanity check.")
			return false
		}
		logrus.WithFields(logrus.Fields{
			"description": test.Description,
			"expected":    fmt.Sprintf("%s%d", test.Result.Operation, test.Result.Count),
			"actual":      count,
		}).Info("Successful test")
	}
	return true
}

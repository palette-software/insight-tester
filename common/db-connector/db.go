package db_connector

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/palette-software/insight-tester/common/logging"
	"sync"
	"time"
)

var DBConnection *sql.DB
var dbc_mutex sync.Mutex

type DbConnector struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Database string `yaml:"Database"`
	Schema   string `yaml:"Schema"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
}

func (dbc *DbConnector) getConnectionString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", dbc.User, dbc.Password,
		dbc.Database, dbc.Host, dbc.Port)
}

func (dbc *DbConnector) getConnection() (*sql.DB, error) {
	dbc_mutex.Lock()
	defer dbc_mutex.Unlock()

	if DBConnection == nil {
		connectionString := dbc.getConnectionString()
		var err error
		DBConnection, err = sql.Open("postgres", connectionString)
		if err != nil {
			return nil, fmt.Errorf("Error in connecting to db: %v", err)
		}
	}
	return DBConnection, nil
}

func CloseDB() {
	DBConnection.Close()
}

// This function is going to be called on each result row of the SQL statement
type ProcessRowFunc func([]string, []interface{}) error

func (dbc *DbConnector) Query(sql_statement string, handler ProcessRowFunc) error {
	db, err := dbc.getConnection()
	if err != nil {
		return err
	}
	start := time.Now()
	rows, err := db.Query(sql_statement)
	if err != nil {
		return fmt.Errorf("Failed to execute SQL statement: %v! Error: %v", sql_statement, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("Failed to get columns for rows %v!", rows)
	}
	column_count := len(columns)

	var values = make([]interface{}, column_count)
	for i, _ := range values {
		values[i] = new(interface{})
	}

	rowCount := 0
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("Failed to scan values of row! Error: %v", err)
		}
		rowCount++
		handler(columns, values)
	}

	log.Debugf("SQL STATEMENT: %v", sql_statement)
	log.Debugf("OK: [HOST:%v] [COUNT:%v] [DURATION:%v]", dbc.Host, rowCount, time.Since(start))
	return nil
}

package db_connector

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/palette-software/insight-tester/common/logging"
	"sync"
	"time"
)

type DbConnector struct {
	Host         string `yaml:"Host"`
	Port         int    `yaml:"Port"`
	Database     string `yaml:"Database"`
	Schema       string `yaml:"Schema"`
	User         string `yaml:"User"`
	Password     string `yaml:"Password"`

	dbConnection *sql.DB
	dbcMutex     sync.Mutex
}

func (dbc *DbConnector) getConnectionString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", dbc.User, dbc.Password,
		dbc.Database, dbc.Host, dbc.Port)
}

func (dbc *DbConnector) getConnection() (*sql.DB, error) {
	dbc.dbcMutex.Lock()
	defer dbc.dbcMutex.Unlock()

	if dbc.dbConnection == nil {
		connectionString := dbc.getConnectionString()
		var err error
		dbc.dbConnection, err = sql.Open("postgres", connectionString)
		if err != nil {
			return nil, fmt.Errorf("Error in connecting to DB! Connection string: %v! Error: %v", connectionString, err)
		}
		log.Info("Created DB connection with connection string:", connectionString)
	}
	return dbc.dbConnection, nil
}

func (dbc *DbConnector) CloseDB() {
	dbc.dbcMutex.Lock()
	defer dbc.dbcMutex.Unlock()

	if dbc.dbConnection == nil {
		return
	}
	dbc.dbConnection.Close()
}

// This function is going to be called on each result row of the SQL statement
type ProcessRowFunc func(columns []string) error

func (dbc *DbConnector) Query(sql_statement string, handler ProcessRowFunc, valuesToFill... interface{}) error {
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

	rowCount := 0
	for rows.Next() {
		err = rows.Scan(valuesToFill...)
		if err != nil {
			return fmt.Errorf("Failed to scan values of row! Error: %v", err)
		}
		rowCount++
		handler(columns)
	}

	log.Debugf("SQL STATEMENT: %v", sql_statement)
	log.Debugf("OK: [HOST:%v] [COUNT:%v] [DURATION:%v]", dbc.Host, rowCount, time.Since(start))
	return nil
}

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
	Schema	 string `yaml:"Schema"`
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
		//var ii interface{}
		//values[i] = &ii
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

	//values := make([]interface{}, column_count)
	//valuePtrs := make([]interface{}, column_count)
	//
	//for rows.Next() {
	//
	//	for i, _ := range columns {
	//		valuePtrs[i] = &values[i]
	//	}
	//
	//	rows.Scan(valuePtrs...)
	//
	//	for i, col := range columns {
	//
	//		var v interface{}
	//
	//		val := values[i]
	//
	//		b, ok := val.([]byte)
	//
	//		if (ok) {
	//			v = string(b)
	//		} else {
	//			v = val
	//		}
	//
	//		fmt.Println(col, v)
	//	}
	//}

	//rowCount := 0
	//for rows.Next() {
	//	rowCount++
	//	handler(rows)
	//
	//	//var count int
	//	//var hostName string
	//	//rows.Scan(&count, &hostName)
	//	//if !checkTest(count, test) {
	//	//	expected := fmt.Sprintf("%s%d", test.Result.Operation, test.Result.Count)
	//	//	log.Errorf("FAILED: [HOST:%v] [MACHINE:%v] [TEST:%v] [EXPECTED:%v] [ACTUAL:%v] [DURATION:%v]", host, hostName, test.Description, expected, count, time.Since(start))
	//	//	return false
	//	//}
	//}

	log.Debugf("SQL STATEMENT: %v", sql_statement)
	log.Debugf("OK: [HOST:%v] [COUNT:%v] [DURATION:%v]", dbc.Host, rowCount, time.Since(start))
	return nil
}

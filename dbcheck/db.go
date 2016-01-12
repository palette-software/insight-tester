package main

import (
    "database/sql"
    "fmt"
    "strings"
    _ "github.com/lib/pq"
)

func getConnectionString(config Database) (string) {
    return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d", config.User.Login, config.User.Password, config.Name, config.Server.Host, config.Server.Port)
}

var DBConnection *sql.DB

func getConnection(config Config) (*sql.DB) {
    if DBConnection == nil {
        connectionString := getConnectionString(config.Database)
        var err error
        DBConnection, err = sql.Open(strings.ToLower(config.Database.Type), connectionString)
        if err != nil {
            Error.Println("Error connecting to DB:", err)
            return nil
        }
    }
    return DBConnection
}

func closeDB() {
    DBConnection.Close()
}

func check(config Config, test Test) (bool) {
    db := getConnection(config)
    if db == nil {
        return false
    }
    rows, err := db.Query(test.Sql)
    if err != nil {
        Error.Println("Error getting rows:", err)
        return false
    }
    defer rows.Close()
    for rows.Next() {
        var count int
        rows.Scan(&count)
        if !checkTest(count, test) {
            return false
        }
    }
    return true
}

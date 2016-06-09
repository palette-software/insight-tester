package main

import (
	"crypto/tls"
	"fmt"
	log "github.com/palette-software/insight-tester/common/logging"
	"net/http"
	"os"
	dbconnector "github.com/palette-software/insight-tester/common/db-connector"
	"time"
)

func main() {
	os.Exit(mainWithExitCode())
}

func mainWithExitCode() int {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s test_config_json palmon_config_xml\n", os.Args[0])
		return 1
	}
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Read up config
	config, err := dbconnector.ParseConfig(os.Args[2])
	if err != nil {
		log.Errorf("Error while loading config: %v", err)
		return 1
	}

	licenseOwner := "Palette Software"
	log.Info("Creating splunk target")
	splunkLogger, err := log.NewSplunkTarget(config.SplunkServerAddress, config.SplunkToken, licenseOwner)
	log.Infof("Created with error: %v", err)
	if err == nil {
		defer splunkLogger.Close()
		log.AddTarget(splunkLogger, log.LevelDebug)
		log.AddTarget(os.Stdout, log.LevelDebug)
	} else {
		fmt.Printf("Faield to create Splunk target.")
		log.Error("Failed to create Splunk target for watchdog! Error: ", err)
	}
	exitCode := 0

	tests, err := getTests(os.Args[1])
	if err != nil {
		log.Errorf("Error while loading test definitions: %v", err)
		return 1
	}
	for _, host := range config.Databases.Hosts {
		for _, test := range tests {
			if !check(host, config.Databases.Params, test) {
				exitCode = 1
			}
		}
	}
	dbconnector.CloseDB()
	return exitCode
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

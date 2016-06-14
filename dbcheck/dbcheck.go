package main

import (
	"crypto/tls"
	"fmt"
	dbconn "github.com/palette-software/insight-tester/common/db-connector"
	log "github.com/palette-software/insight-tester/common/logging"
	"net/http"
	"os"
	"time"
)

func main() {
	os.Exit(mainWithExitCode())
}

func mainWithExitCode() int {
	log.AddTarget(os.Stdout, log.LevelDebug)

	if len(os.Args) < 3 {
		log.Errorf("Usage: %s test_yml config_yml\n", os.Args[0])
		return 1
	}
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Read up config
	config, err := dbconn.ParseConfig(os.Args[2])
	if err != nil {
		log.Errorf("Error while loading config: %v", err)
		return 1
	}

	dbConnector := config.DbConnector
	defer dbConnector.CloseDB()

	licenseOwner := "Palette Software"
	log.Info("Creating splunk target")
	splunkLogger, err := log.NewSplunkTarget(config.SplunkServerAddress, config.SplunkToken, licenseOwner)
	if err == nil {
		defer splunkLogger.Close()
		log.AddTarget(splunkLogger, log.LevelDebug)
	} else {
		fmt.Printf("Faield to create Splunk target.")
		log.Error("Failed to create Splunk target! Error: ", err)
	}
	exitCode := 0

	tests, err := getTests(os.Args[1])
	if err != nil {
		log.Errorf("Error while loading test definitions: %v", err)
		return 1
	}

	for _, test := range tests {
		if !check(dbConnector, test) {
			exitCode = 1
		}
	}
	return exitCode
}

func check(dbConnector dbconn.DbConnector, test Test) bool {
	start := time.Now()
	rowCount := 0

	var count int
	var hostName string

	err := dbConnector.Query(test.Sql, func(columns []string) error {
		rowCount++
		if len(columns) != 2 {
			return fmt.Errorf("Exactly 2 columns are expected during check! Got %v instead. SQL statement: %v",
				len(columns), test.Sql)
		}

		if !checkTest(count, test) {
			expected := fmt.Sprintf("%s%d", test.Result.Operation, test.Result.Count)
			return fmt.Errorf("FAILED: [HOST:%v] [MACHINE:%v] [TEST:%v] [EXPECTED:%v] [ACTUAL:%v] [DURATION:%v]",
				dbConnector.Host, hostName, test.Description, expected, count, time.Since(start))

		}

		return nil
	}, &count, &hostName)

	if err != nil {
		log.Errorf("Test query failed! Query: %v Error: %v", test.Sql, err)
		return false
	}

	log.Infof("OK: [HOST:%v] [TEST:%v] [COUNT:%v] [DURATION:%v]", dbConnector.Host, test.Description, rowCount, time.Since(start))
	return true
}

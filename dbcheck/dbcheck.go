package main

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"regexp"

	dbconn "github.com/palette-software/insight-tester/common/db-connector"
	insight_server "github.com/palette-software/insight-server/lib"
	log "github.com/palette-software/insight-tester/common/logging"
)

func main() {
	os.Exit(mainWithExitCode())
}

func mainWithExitCode() int {
	log.AddTarget(os.Stdout, log.LevelDebug)

	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s test_yml config_yml\n", os.Args[0])
	}

	// Read up config
	config, err := dbconn.ParseConfig(os.Args[2])
	if err != nil {
		log.Fatalf("Error while loading config: %v", err)
	}

	dbConnector := config.DbConnector
	defer dbConnector.CloseDB()

	license, err := getCustomerName()
	if err != nil {
		log.Error(err)
		// Without license, we don't know the customer name, so there is
		//  no point in logging to Splunk
		log.Error("Performing Sanity checks without Splunk is in vain!")
		return 1
	} else {
		log.Info("Acquired customer name from license: ", license.Name)

		log.Info("Creating splunk target")
		splunkLogger, err := log.NewSplunkTarget(config.SplunkServerAddress, config.SplunkToken, license.Name)
		if err == nil {
			defer splunkLogger.Close()
			log.AddTarget(splunkLogger, log.LevelDebug)
		} else {
			log.Error("Failed to create Splunk target! Error: ", err)
			return 1
		}
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

func getCustomerName() (*insight_server.LicenseData, error) {
	// Open the Insight Server config file to read the license
	contentBytes, err := ioutil.ReadFile("/etc/palette-insight-server/server.config")
	if err != nil {
		return nil, fmt.Errorf("Failed to open Insight Server config file to read the license key!")
	}

	var licenseKeyExpr = regexp.MustCompile(`license_key=([\S]+)`)
	content := string(contentBytes)
	matchGroups := licenseKeyExpr.FindStringSubmatch(content)
	// The first group is the entire match, this is why we expect more than one
	if len(matchGroups) < 2 {
		return nil, fmt.Errorf("No license key found in Insight Server config file!")
	}

	license := insight_server.UpdateLicense(matchGroups[1])
	if license == nil {
		return nil, fmt.Errorf("License is nil!")
	}

	return license, nil
}

func check(dbConnector dbconn.DbConnector, test Test) bool {
	start := time.Now()

	var count int
	var hostName string

	err := dbConnector.Query(test.Sql, func(columns []string) error {
		if len(columns) != 2 {
			return fmt.Errorf("Exactly 2 columns are expected during check! Got %v instead. SQL statement: %v",
				len(columns), test.Sql)
		}

		expected := fmt.Sprintf("%s%d", test.Result.Operation, test.Result.Count)
		if !checkTest(count, test) {
			return fmt.Errorf("FAILED: [HOST:%v] [MACHINE:%v] [TEST:%v] [EXPECTED:%v] [ACTUAL:%v] [DURATION:%v]",
				dbConnector.Host, hostName, test.Description, expected, count, time.Since(start))
		}

		log.Infof("OK: [HOST:%v] [MACHINE:%v] [TEST:%v] [EXPECTED:%v] [ACTUAL:%v] [DURATION:%v]",
			dbConnector.Host, hostName, test.Description, expected, count, time.Since(start))

		return nil
	}, &count, &hostName)

	if err != nil {
		log.Errorf("Test query failed! Query: %v Error: %v", test.Sql, err)
		return false
	}

	return true
}

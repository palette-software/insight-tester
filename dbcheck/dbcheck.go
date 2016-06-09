package main

import (
	"crypto/tls"
	"fmt"
	dbconn "github.com/palette-software/insight-tester/common/db-connector"
	log "github.com/palette-software/insight-tester/common/logging"
	"net/http"
	"os"
	"reflect"
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
	config, err := dbconn.ParseConfig(os.Args[2])
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

	defer dbconn.CloseDB()

	for _, test := range tests {
		if !check(config.DbConnector, test) {
			exitCode = 1
		}
	}
	return exitCode
}

func check(dbconnector dbconn.DbConnector, test Test) bool {
	start := time.Now()
	rowCount := 0
	err := dbconnector.Query(test.Sql, func(columns []string, values []interface{}) error {
		rowCount++
		if len(values) < 2 {
			return fmt.Errorf("Not enough values returned during check! SQL statement: %v", test.Sql)
		}

		count := *(values[0].(*interface{}))
		var countType = reflect.TypeOf(count)

		if countType.Kind() != reflect.Int {
			return fmt.Errorf("Count value is expeted to be an integer, but it is %v! SQL statement: %v",
				countType, test.Sql)
		}

		hostName := *(values[1].(*interface{}))
		var hostNameType = reflect.TypeOf(hostName)

		if hostNameType.Kind() != reflect.String {
			return fmt.Errorf("Host name value is expeted to be a string, but it is %v! SQL statement: %v",
				hostNameType, test.Sql)
		}

		hostName = hostName.(string)

		if !checkTest(count.(int), test) {
			expected := fmt.Sprintf("%s%d", test.Result.Operation, test.Result.Count)
			return fmt.Errorf("FAILED: [HOST:%v] [MACHINE:%v] [TEST:%v] [EXPECTED:%v] [ACTUAL:%v] [DURATION:%v]",
				dbconnector.Host, hostName, test.Description, expected, count, time.Since(start))

		}

		return nil
	})

	if err != nil {
		log.Errorf("Test query failed! Query: %v Error: %v", test.Sql, err)
		return false
	}

	log.Infof("OK: [HOST:%v] [TEST:%v] [COUNT:%v] [DURATION:%v]", dbconnector.Host, test.Description, rowCount, time.Since(start))
	return true
}

package main

import (
	"crypto/tls"
	"fmt"
	log "github.com/palette-software/insight-tester/common/logging"
	"net/http"
	"os"
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
	config, err := parseConfig(os.Args[2])
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
	closeDB()
	return exitCode
}

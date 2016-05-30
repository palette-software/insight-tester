package main

import (
	"crypto/tls"
	"fmt"
	log "github.com/palette-software/insight-tester/common/logging"
	"net/http"
	"os"
)

const SplunkServerAddress = "splunk-insight.palette-software.net"
const SplunkToken = "F59DC3B4-AAB9-4330-ACED-682E1681B507"

func main() {
	os.Exit(mainWithExitCode())
}

func mainWithExitCode() int {
	if len(os.Args) < 3 {
		log.Error("Usage: %s test_config_json palmon_config_xml\n", os.Args[0])
		os.Exit(1)
	}
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	licenseOwner := "Palette Software"
	log.Info("Creating splunk target")
	splunkLogger, err := log.NewSplunkTarget(SplunkServerAddress, SplunkToken, licenseOwner)
	log.Infof("Created with error: %v", err)
	if err == nil {
		defer splunkLogger.Close()
		log.AddTarget(splunkLogger, log.DebugLevel)
		log.AddTarget(os.Stdout, log.DebugLevel)
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
	database, err := getResultDBConfig(os.Args[2])
	if err != nil {
		log.Errorf("Error while loading config: %v", err)
		return 1
	}

	for _, host := range database.Hosts {
		for _, test := range tests {
			if !check(host, database.Params, test) {
				exitCode = 1
			}
		}
	}
	closeDB()
	return exitCode
}

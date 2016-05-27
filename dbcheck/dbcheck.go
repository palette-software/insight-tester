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
	} else {
		fmt.Printf("Faield to create Splunk target.")
		log.Error("Failed to create Splunk target for watchdog! Error: ", err)
	}
	exitCode := 0

	tests := getTests(os.Args[1])
	database := getResultDBConfig(os.Args[2])
	for _, test := range tests {
		if !check(database, test) {
			exitCode = 1
		}
	}
	closeDB()
	return exitCode
}

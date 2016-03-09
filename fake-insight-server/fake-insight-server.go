package main

import (
	"net/http"
	log "github.com/palette-software/insight-tester/common/logging"
	"fmt"
	"flag"
	"io/ioutil"
	"os"
)

// Just respond yes to every request
func okToEverything(w http.ResponseWriter, r *http.Request, name string) {
	// signal that everything went ok
	log.Info.Printf("Responding OK to %s request.", name)
	http.Error(w, "", http.StatusOK)
}

func main() {
	// Initialize the log to write into file instead of stderr
	// open output file
	logFile, err := os.OpenFile("fake_insight_server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Failed to open log file! ", err)
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := logFile.Close(); err != nil {
			fmt.Println("Failed to close log file! ", err)
			panic(err)
		}
	}()

	// Set the levels to be ignored to ioutil.Discard
	// Levels:  TRACE           INFO     WARNING  ERROR    FATAL
	log.InitLog(ioutil.Discard, logFile, logFile, logFile, logFile)

	log.Info.Println("Starting up the fake insight server...")

	// Variables for holding the server address and port
	var bindAddress string
	var bindPort int

	flag.IntVar(&bindPort, "port", 9000, "The port the server is binding itself to")
	flag.StringVar(&bindAddress, "bind_address", "localhost", "The address to bind to. Leave empty for default which is localhost .")

	// create the upload endpoint
	authenticatedUploadHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			okToEverything(w, r, "upload")
	})

	// create the maxid handler
	maxIdHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		okToEverything(w, r, "maxid")
	})

	// Create the fake endpoints
	http.HandleFunc("/upload", authenticatedUploadHandler)
	http.HandleFunc("/maxid", maxIdHandler)

	bindAddressWithPort := fmt.Sprintf("%s:%v", bindAddress, bindPort)
	log.Info.Println("[http] Webservice starting on ", bindAddressWithPort)

	//if useTls {
	//	err := http.ListenAndServeTLS(bindAddressWithPort, tlsCert, tlsKey, nil)
	//	log.Fatal(err)
	//} else {

		err = http.ListenAndServe(bindAddressWithPort, nil)
		log.Fatal(err)
	//}
}

package main

import (
	"net/http"
	log "github.com/palette-software/insight-tester/common/logging"
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"encoding/json"
)

// The base structure for a SemVer like version
type Version struct {
	// The version according to SemVer
	Major, Minor, Patch int
}

// Combines a version with an actual product and a file
type UpdateVersion struct {
	Version
	// The name of the product
	Product string
	// The Md5 checksum of this update
	Md5 string
	// The url where this update can be downloaded from
	Url string
}

// Just respond yes to every request
func okToEverything(w http.ResponseWriter, r *http.Request, name string) {
	// signal that everything went ok
	log.Infof("Responding OK to %s request.", name)
	http.Error(w, "", http.StatusOK)
}

func makeFakeHandler(endpoint string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		okToEverything(w, r, endpoint)
	})
}

func main() {
	log.Init()
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
	log.AddTarget(logFile, log.InfoLevel)

	log.Info("Starting up the fake insight server...")

	// Variables for holding the server address and port
	var bindAddress string
	var bindPort int

	flag.IntVar(&bindPort, "port", 9000, "The port the server is binding itself to")
	flag.StringVar(&bindAddress, "bind_address", "localhost", "The address to bind to. Leave empty for default which is localhost .")

	// create the fake handlers
	authenticatedUploadHandler := makeFakeHandler("upload")
	maxIdHandler := makeFakeHandler("maxid")
	updateCheckHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Expires", "0")

		// Put a fake version in the response
		fakeVersion := Version{6, 3, 2}
		fakeUrl := fmt.Sprintf("https://%s:%d/updates/products/agent/versions/v1.3.2", bindAddress, bindPort)
		fakeUpdateVersion := UpdateVersion{
			Version:    fakeVersion,
			Product:	"agent",
			Md5:		"cool-md5-hash",
			Url:		fakeUrl,
		}

		if err := json.NewEncoder(w).Encode(fakeUpdateVersion); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// And report that everything was OK
		http.Error(w, "", http.StatusOK)
	})

	// Create the fake endpoints
	http.HandleFunc("/upload", authenticatedUploadHandler)
	http.HandleFunc("/maxid", maxIdHandler)
	http.HandleFunc("/updates/latest-version", updateCheckHandler)

	bindAddressWithPort := fmt.Sprintf("%s:%v", bindAddress, bindPort)
	log.Info("[http] Webservice starting on ", bindAddressWithPort)

	//if useTls {
	//	err := http.ListenAndServeTLS(bindAddressWithPort, tlsCert, tlsKey, nil)
	//	log.Fatal(err)
	//} else {
		err = http.ListenAndServe(bindAddressWithPort, nil)
		log.FatalLevel(err)
	//}
}

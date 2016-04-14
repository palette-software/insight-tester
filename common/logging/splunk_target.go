package logging

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/palette-software/insight-server"

	"github.com/kardianos/osext"
	"github.com/zfjagann/golang-ring"
)

type SplunkTarget struct {
	Protocol     string
	Owner        string
	Host         string
	Port         string
	Token        string
	Ring         *ring.Ring
	Ticker       *time.Ticker
	TickInterval int
	Capacity     int
}

type Message struct {
	Event string `json:"event"`
}

func formatSplunkMessage(p string) []byte {
	m := Message{
		Event: p,
	}
	jsonObject, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return jsonObject
}

func (t SplunkTarget) Write(p []byte) (n int, err error) {
	// This conversion is needed as otherwise we overwrite the enqueued items.
	message := fmt.Sprintf("[OW:%s]  %s", t.Owner, p)
	t.Ring.Enqueue(message)
	return n, nil
}

func (t SplunkTarget) Start() {
	t.Ring.SetCapacity(t.Capacity)
	t.Ticker = time.NewTicker(time.Duration(t.TickInterval) * time.Millisecond)
	go func(t *SplunkTarget) {
		for range t.Ticker.C {
			t.SendLogs()
		}
	}(&t)
}

func (t SplunkTarget) SendLogs() {
	var b bytes.Buffer
	for {
		next := t.Ring.Dequeue()
		if next == nil {
			break
		}
		line := fmt.Sprintf("%s", next)
		formattedMessage := formatSplunkMessage(line)
		if formattedMessage != nil {
			b.Write(formatSplunkMessage(line))
		}
	}
	// Return if there's no new records
	if b.Len() == 0 {
		return
	}

	// Send the records to Splunk
	endpoint := fmt.Sprintf("%s://%s:%s/services/collector", t.Protocol, t.Host, t.Port)
	request, err := http.NewRequest("POST", endpoint, &b)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Splunk %s", t.Token))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	client.Do(request)
}

func NewSplunkTarget(Host, Token string) *SplunkTarget {
	ownerName, err := getOwner()
	if err != nil {
		// Without an owner name there is no point in sending logs to Splunk, otherwise we will
		// not be able to identify the source of the log files in Splunk.
		return nil
	}

	st := SplunkTarget{
		Owner:        ownerName,
		Host:         Host,
		Token:        Token,
		Protocol:     "https",
		Port:         "443",
		Ring:         &ring.Ring{},
		TickInterval: 300,
		Capacity:     65000,
	}
	st.Start()
	return &st
}

func getOwner() (string, error) {
	execFolder, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Println("Failed to get executable folder for Splunk target: ", err)
		return "", err
	}

	// Check for license
	files, err := ioutil.ReadDir(execFolder)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("Failed to read exec dir for license files: ", err)
		return "", err
	}

	ownerName := ""
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".license") {
			ownerName, err = queryOwnerOfLicense(filepath.Join(execFolder, file.Name()))
			if err != nil {
				continue
			}
			break
		}
	}

	if ownerName == "" {
		fmt.Println("No valid license file found!")
		return "", err
	}

	return ownerName, nil
}

func queryOwnerOfLicense(licenseFile string) (string, error) {
	request, err := newfileUploadRequest("http://localhost:9000/license-check", "file", licenseFile)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("Fatal error while creating new file upload request: ", err)
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		//log.Fatal(err)
		fmt.Printf("Client do request failed! Request: %v. Error message: %v", request, err)
		return "", err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("Fatal error while reading from response body: ", err)
		return "", err
	}
	resp.Body.Close()
	fmt.Println("Status code:", resp.StatusCode)
	fmt.Println("Header:", resp.Header)
	fmt.Println("Body:", body)

	// Decode the JSON in the response
	var licenseCheck insight_server.LicenseCheckResponse
	if err := json.NewDecoder(body).Decode(&licenseCheck); err != nil {
		//log.Error.Printf("Error while deserializing license check response body! Error message: %v", err)
		fmt.Printf("Error while deserializing license check response body! Error message: %v", err)
		return "", err
	}

	if !licenseCheck.Valid {
		err = fmt.Errorf("License: %v is invalid! Although owner name is %v", licenseFile, licenseCheck.OwnerName)
		fmt.Println(err)
		return "", err
	}

	// FIXME: Remove println-s!!
	fmt.Println("Owner is", licenseCheck.OwnerName)
	return licenseCheck.OwnerName, nil
}

// Creates a new file upload http request with multipart file
func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		fmt.Println("Failed to create new request! Error message: %v", err)
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

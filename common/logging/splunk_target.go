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
	"sync"
	"time"

	"github.com/palette-software/insight-server"

	"github.com/kardianos/osext"
	"github.com/zfjagann/golang-ring"
)

type SplunkTarget struct {
	Protocol     string
	Owner        string
	MachineName  string
	Host         string
	Port         string
	Token        string
	Ring         *ring.Ring
	Ticker       *time.Ticker
	TickInterval int
	Capacity     int
	ringMutex    sync.Mutex
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

func (t *SplunkTarget) Write(p []byte) (n int, err error) {
	// Remove the newline characters from the end of the stream, if there are any, as Splunk does not need them.
	p = bytes.TrimSuffix(p, []byte("\n"))

	// This conversion is needed as otherwise we overwrite the enqueued items. And we need to show the owner and
	// the machine name in the log lines.
	message := fmt.Sprintf("[OW:%s] [MN:%s]  %s", t.Owner, t.MachineName, p)
	t.ringMutex.Lock()
	defer t.ringMutex.Unlock()
	t.Ring.Enqueue(message)
	return n, nil
}

func (t *SplunkTarget) Start() {
	t.Ring.SetCapacity(t.Capacity)
	t.Ticker = time.NewTicker(time.Duration(t.TickInterval) * time.Millisecond)
	go func(t *SplunkTarget) {
		for range t.Ticker.C {
			t.SendLogs()
		}
	}(t)
}

func (t *SplunkTarget) DequeueLines() (b bytes.Buffer) {
	t.ringMutex.Lock()
	defer t.ringMutex.Unlock()
	for {
		next := t.Ring.Dequeue()
		if next == nil {
			break
		}
		line := fmt.Sprintf("%s", next)
		formattedMessage := formatSplunkMessage(line)
		if formattedMessage != nil {
			b.Write(formattedMessage)
		}
	}

	return b
}

func (t *SplunkTarget) SendLogs() {
	b := t.DequeueLines()
	// Return if there are no new records
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
	_, err = client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request! Exception message: ", err)
	}
}

func (t *SplunkTarget) Close() error {
	// FIXME: Wait 3 seconds, so that Splunk logger has some time to upload the logs
	tickUpdate := time.Tick(3 * time.Second)

	select {
	case <-tickUpdate:
		// Time is up. Do not wait any longer.
		return fmt.Errorf("Not all messages are sent to Splunk, time is up. Force close.")
	}
	return nil
}

func NewSplunkTarget(Host, Token string) (*SplunkTarget, error) {
	ownerName, err := getOwner()
	if err != nil {
		// Without an owner name there is no point in sending logs to Splunk, otherwise we will
		// not be able to identify the source of the log files in Splunk.
		return nil, err
	}

	machineName, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("Failed to get machine name! Error: %v", err)
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
		MachineName:  machineName,
	}
	st.Start()
	return &st, nil
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

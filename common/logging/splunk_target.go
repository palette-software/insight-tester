package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

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
	var client http.Client
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

func NewSplunkTarget(Host, Token, Owner string) (*SplunkTarget, error) {
	if Owner == "" {
		// Without an owner name there is no point in sending logs to Splunk, otherwise we will
		// not be able to identify the source of the log files in Splunk.
		return nil, fmt.Errorf("Empty owner name is not allowed for Splunk target!")
	}

	machineName, err := os.Hostname()
	if err != nil {
		// Go on anyway. It's better to have something, than nothing.
		machineName = "UNKNOWN_MACHINE"
	}

	st := SplunkTarget{
		Owner:        Owner,
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

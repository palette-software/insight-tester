package logging

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/zfjagann/golang-ring"
)

type SplunkTarget struct {
    Protocol string
    Host string
    Port string
    Token string
    Ring* ring.Ring
    Ticker* time.Ticker
    TickInterval int
    Capacity int
}

type Message struct {
    Event string  `json:"event"`
}

func formatSplunkMessage(p string) ([] byte) {
    m := Message{
        Event: p,
    }
    jsonObject, _ := json.Marshal(m) 
    return jsonObject
}

func (t SplunkTarget) Write(p [] byte) (n int, err error) {
    // This conversion is needed as otherwise we overwrite the enqueued items.
    message := fmt.Sprintf("%s", p)
    t.Ring.Enqueue(message)
    return n, nil
}

func (t SplunkTarget) Start() {
    t.Ring.SetCapacity(t.Capacity);
    t.Ticker = time.NewTicker(time.Duration(t.TickInterval) * time.Millisecond)
    go func(t* SplunkTarget) {
        for range t.Ticker.C {
            t.SendLogs()
        }
    }(&t)
}

func (t SplunkTarget) SendLogs() {
    endpoint := fmt.Sprintf("%s://%s:%s/services/collector", t.Protocol, t.Host, t.Port) 
    var b bytes.Buffer
    for {
        next := t.Ring.Dequeue()
        if next == nil {
            break;
        }
        line := fmt.Sprintf("%s", next)
        b.Write(formatSplunkMessage(line))
    }
    // Return if there's no new records
    if b.Len() == 0 {
        return
    }

    // Send the records to Splunk
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
    st := SplunkTarget{
        Host: Host,
        Token: Token,
        Protocol: "https",
        Port: "443",
        Ring: &ring.Ring{}, 
        TickInterval: 300, 
        Capacity: 10,
    }
    st.Start() 
    return &st
}


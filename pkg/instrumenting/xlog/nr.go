package xlog

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var hostname string

const (
	requestTimeout = time.Second * 10

	NewRelicLogEndpoint = "https://log-api.newrelic.com/log/v1"
)

// NewRelicClient maintains the authorization token for NR
// and is the entry point for sending entries upstream
type NewRelicClient struct {
	LicenseKey string
}

// isValid ensures the new relic client is not nil or empty
func (c *NewRelicClient) isValid() bool {
	if c != nil {
		return c.LicenseKey != ""
	}

	return false
}

// request sends our payload to New Relic
func (c *NewRelicClient) request(jsonData []byte) error {
	req, err := http.NewRequest("POST", NewRelicLogEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("could not make a request to the New Relic Log API: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	req.Header.Set("X-License-Key", c.LicenseKey)

	client := &http.Client{Timeout: requestTimeout}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return err
}

// setHostName sets the global hostname once.
func setHostName() {
	hostnameOnce.Do(func() {
		h, err := os.Hostname()
		if err != nil {
			log.Println("WARN: unable to fetch hostname")
		}

		hostname = h
	})
}

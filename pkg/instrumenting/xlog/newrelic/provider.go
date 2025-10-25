package newrelic

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/michaeldelorenzo/x/pkg/instrumenting/xlog/types"
)

const (
	requestTimeout       = time.Second * 10
	NewRelicLogEndpoint = "https://log-api.newrelic.com/log/v1"
)

// Provider implements the LogProvider interface for New Relic
type Provider struct {
	licenseKey       string
	reportableLevels []zapcore.Level
	m                sync.Mutex
}

// NewProvider creates a new New Relic log provider
func NewProvider(conf *Config) *Provider {
	levels := conf.ReportableLevels
	if levels == nil {
		levels = allLevels()
	}

	return &Provider{
		licenseKey:       conf.LicenseKey,
		reportableLevels: levels,
	}
}

// SendLog sends a log entry to New Relic
func (p *Provider) SendLog(entry zapcore.Entry, message string) error {
	p.m.Lock()
	defer p.m.Unlock()

	req, err := http.NewRequest("POST", NewRelicLogEndpoint, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return fmt.Errorf("could not make request to New Relic Log API: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	req.Header.Set("X-License-Key", p.licenseKey)

	client := &http.Client{Timeout: requestTimeout}
	_, err = client.Do(req)
	return err
}

// IsValid checks if the provider is properly configured
func (p *Provider) IsValid() bool {
	return p.licenseKey != ""
}

// Type returns the provider type
func (p *Provider) Type() types.ProviderType {
	return types.ProviderNewRelic
}

// ShouldSend determines if this log level should be sent
func (p *Provider) ShouldSend(level zapcore.Level) bool {
	for _, l := range p.reportableLevels {
		if l == level {
			return true
		}
	}
	return false
}

// allLevels returns all supported log levels
func allLevels() []zapcore.Level {
	return []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.FatalLevel,
		zapcore.PanicLevel,
	}
}

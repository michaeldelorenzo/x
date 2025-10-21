package newrelic

import (
	"log"
	"sync"

	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

var (
	// Once register to ensure we don't allocate unnecessary resources
	initOnce sync.Once
)

// Config defines the available configuration options for New Relic agent.
type Config struct {
	LicenseKey        string
	ApplicationName   string
	TrustedAccountKey string
	AdditionalConfig  newrelic.ConfigOption
}

// Init initializes a New Relic APM application with the provided configuration
// and configLogger for echoing New Relic agent output.
func Init(conf *Config, configLogger *zap.Logger) *Provider {
	var app *newrelic.Application

	initOnce.Do(func() {
		// If a custom logger is not provided, we use zap's opinionated
		// production logger. This is only used for logging newrelic agent
		// messages and not publishing logs. See xlog.
		if configLogger == nil {
			configLogger, _ = zap.NewProduction()
		}

		var err error
		app, err = newrelic.NewApplication(
			conf.AdditionalConfig,
			newrelic.ConfigEnabled(conf.LicenseKey != ""),
			newrelic.ConfigAppName(conf.ApplicationName),
			newrelic.ConfigLicense(conf.LicenseKey),
			newrelic.ConfigDistributedTracerEnabled(true),
			nrzap.ConfigLogger(configLogger),
		)

		if err != nil {
			log.Println("WARN: encountered error initializing New Relic client:", err.Error())
		}
	})

	return NewProvider(app)
}

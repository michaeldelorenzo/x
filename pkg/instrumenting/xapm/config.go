package xapm

import (
	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"

	"log"
	"sync"
	"time"
)

var (
	//Once register to ensure we don't allocate unnecessary resources
	initOnce sync.Once
	// Apm is our new relic agent
	Apm *newrelic.Application
)

type DBType = newrelic.DatastoreProduct

// Commonly Supported databases
const (
	DBRedis    DBType = "Redis"
	DBMongoDB  DBType = "MongoDB"
	DBPostgres DBType = "Postgres"
)

// Common Operations used to group db & categorize db segments
type DBOperation string

const (
	DBSelect DBOperation = "SELECT"
	DBInsert DBOperation = "INSERT"
	DBUpdate DBOperation = "UPDATE"
	DBUpsert DBOperation = "UPSERT"
	DBDelete DBOperation = "DELETE"
	DBCount  DBOperation = "COUNT"
)

// Config defines the available configuration options for new relic agent.
type Config struct {
	LicenseKey        string
	ApplicationName   string
	TrustedAccountKey string
	AdditionalConfig  newrelic.ConfigOption
}

// Init initializes our Apm application with the provided configuration
// and configLogger for echoing new relic agent output.
func Init(conf *Config, configLogger *zap.Logger) *newrelic.Application {
	initOnce.Do(func() {
		// If a custom logger is not provided, we use zap's opinionated
		// production logger. This is only used for logging newrelic agent
		// messages and not publishing logs. See xlog.
		if configLogger == nil {
			configLogger, _ = zap.NewProduction()
		}

		app, err := newrelic.NewApplication(
			conf.AdditionalConfig,
			newrelic.ConfigEnabled(conf.LicenseKey != ""),
			newrelic.ConfigAppName(conf.ApplicationName),
			newrelic.ConfigLicense(conf.LicenseKey),
			newrelic.ConfigDistributedTracerEnabled(true),
			nrzap.ConfigLogger(configLogger),
		)

		if err != nil {
			log.Println("WARN: encountered error initializing new relic client:", err.Error())
		}

		Apm = app
	})

	return Apm
}

// Shutdown graceful shuts down the global Apm instance
func Shutdown(timeout time.Duration) {
	Apm.Shutdown(timeout)
}

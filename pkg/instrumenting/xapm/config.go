package xapm

import (
	"fmt"
	"log"

	"github.com/michaeldelorenzo/x/v2/pkg/instrumenting/xapm/newrelic"
	"github.com/michaeldelorenzo/x/v2/pkg/instrumenting/xapm/noop"
	"github.com/michaeldelorenzo/x/v2/pkg/instrumenting/xapm/sentry"
	"go.uber.org/zap"
)

// Config contains the configuration for initializing an APM provider
type Config struct {
	// Provider specifies which APM provider to use
	Provider ProviderType

	// NewRelic specific configuration
	NewRelic *newrelic.Config

	// Sentry specific configuration
	Sentry *sentry.Config

	// Logger is used for New Relic agent logging (optional)
	Logger *zap.Logger
}

// Init initializes the global APM provider based on the configuration
func Init(conf *Config) error {
	var provider Provider

	switch conf.Provider {
	case ProviderNewRelic:
		if conf.NewRelic == nil {
			return fmt.Errorf("NewRelic configuration is required when Provider is %s", ProviderNewRelic)
		}
		provider = newrelic.Init(conf.NewRelic, conf.Logger)
		log.Printf("Initialized New Relic APM provider")

	case ProviderSentry:
		if conf.Sentry == nil {
			return fmt.Errorf("Sentry configuration is required when Provider is %s", ProviderSentry)
		}
		provider = sentry.Init(conf.Sentry)
		log.Printf("Initialized Sentry APM provider")

	case ProviderNoop:
		provider = noop.NewProvider()
		log.Printf("Initialized no-op APM provider")

	default:
		return fmt.Errorf("unsupported APM provider: %s", conf.Provider)
	}

	SetProvider(provider)
	return nil
}

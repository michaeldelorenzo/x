package sentry

import (
	"log"
	"sync"

	"github.com/getsentry/sentry-go"
)

var (
	// Once register to ensure we don't allocate unnecessary resources
	initOnce sync.Once
)

// Config defines the available configuration options for Sentry.
type Config struct {
	DSN              string
	Environment      string
	Release          string
	TracesSampleRate float64
	Debug            bool
}

// Init initializes a Sentry APM provider with the provided configuration.
func Init(conf *Config) *Provider {
	var hub *sentry.Hub

	initOnce.Do(func() {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              conf.DSN,
			Environment:      conf.Environment,
			Release:          conf.Release,
			TracesSampleRate: conf.TracesSampleRate,
			Debug:            conf.Debug,
			EnableTracing:    true,
		})

		if err != nil {
			log.Println("WARN: encountered error initializing Sentry client:", err.Error())
		}

		// Get the current hub or create a new one
		hub = sentry.CurrentHub()
	})

	return NewProvider(hub)
}

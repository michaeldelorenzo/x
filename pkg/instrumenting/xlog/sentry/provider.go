package sentry

import (
	"errors"
	"log"
	"sync"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"

	"github.com/michaeldelorenzo/x/pkg/instrumenting/xlog/types"
)

var (
	// initOnce ensures we only initialize the Sentry SDK once
	initOnce sync.Once
)

// Provider implements the LogProvider interface for Sentry
type Provider struct {
	hub         *sentry.Hub
	eventLevels []zapcore.Level
	logLevels   []zapcore.Level
}

// NewProvider creates a new Sentry log provider
func NewProvider(conf *Config) *Provider {
	initOnce.Do(func() {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         conf.DSN,
			Environment: conf.Environment,
			Release:     conf.Release,
			Debug:       conf.Debug,
		})

		if err != nil {
			log.Println("WARN: error initializing Sentry client:", err)
		}
	})

	hub := sentry.CurrentHub()
	if hub.Client() == nil {
		hub = nil
	}

	eventLevels := conf.EventLevels
	if eventLevels == nil {
		eventLevels = []zapcore.Level{zapcore.ErrorLevel, zapcore.FatalLevel, zapcore.PanicLevel}
	}

	logLevels := conf.LogLevels
	if logLevels == nil {
		logLevels = []zapcore.Level{
			zapcore.DebugLevel,
			zapcore.InfoLevel,
			zapcore.WarnLevel,
			zapcore.ErrorLevel,
			zapcore.FatalLevel,
			zapcore.PanicLevel,
		}
	}

	return &Provider{
		hub:         hub,
		eventLevels: eventLevels,
		logLevels:   logLevels,
	}
}

// SendLog sends a log entry to Sentry
func (p *Provider) SendLog(entry zapcore.Entry, message string) error {
	if p.hub == nil {
		return errors.New("sentry provider: hub is nil, provider not properly initialized")
	}

	sentryLevel := zapLevelToSentry(entry.Level)

	// Check if this should be an Event (creates an issue in Sentry)
	if p.shouldSendAsEvent(entry.Level) {
		event := sentry.NewEvent()
		event.Level = sentryLevel
		event.Message = entry.Message
		event.Logger = entry.LoggerName
		event.Timestamp = entry.Time
		// sentry-go removed Event.Extra (the legacy "Additional Data" field) in
		// favour of structured contexts. sentry.NewEvent initialises Contexts,
		// so this index assignment is safe without a nil check.
		event.Contexts["extra"] = sentry.Context{"log_json": message}

		if p.hub.CaptureEvent(event) == nil {
			return errors.New("sentry provider: event was discarded (client not configured, rate-limited, or filtered by BeforeSend)")
		}
	}

	// Send as breadcrumb (structured log data that attaches to future events)
	if p.shouldSendAsLog(entry.Level) {
		p.hub.AddBreadcrumb(&sentry.Breadcrumb{
			Level:     sentryLevel,
			Message:   message,
			Timestamp: entry.Time,
			Category:  entry.LoggerName,
		}, nil)
	}

	return nil
}

// IsValid checks if the provider is properly configured
func (p *Provider) IsValid() bool {
	return p.hub != nil
}

// Type returns the provider type
func (p *Provider) Type() types.ProviderType {
	return types.ProviderSentry
}

// ShouldSend determines if this log level should be sent
func (p *Provider) ShouldSend(level zapcore.Level) bool {
	return p.shouldSendAsEvent(level) || p.shouldSendAsLog(level)
}

// shouldSendAsEvent checks if the level should create a Sentry event
func (p *Provider) shouldSendAsEvent(level zapcore.Level) bool {
	for _, l := range p.eventLevels {
		if l == level {
			return true
		}
	}
	return false
}

// shouldSendAsLog checks if the level should be sent as a log entry
func (p *Provider) shouldSendAsLog(level zapcore.Level) bool {
	for _, l := range p.logLevels {
		if l == level {
			return true
		}
	}
	return false
}

// zapLevelToSentry converts zap log levels to Sentry levels
func zapLevelToSentry(level zapcore.Level) sentry.Level {
	switch level {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.FatalLevel, zapcore.PanicLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelInfo
	}
}

package types

import "go.uber.org/zap/zapcore"

// ProviderType represents the log provider type
type ProviderType string

const (
	ProviderNewRelic ProviderType = "newrelic"
	ProviderSentry   ProviderType = "sentry"
	ProviderNoop     ProviderType = "noop"
)

// LogProvider defines the interface that all log providers must implement
type LogProvider interface {
	// SendLog sends a log entry to the provider
	SendLog(entry zapcore.Entry, message string) error

	// IsValid checks if the provider is properly configured
	IsValid() bool

	// Type returns the provider type
	Type() ProviderType

	// ShouldSend determines if this log level should be sent
	ShouldSend(level zapcore.Level) bool
}

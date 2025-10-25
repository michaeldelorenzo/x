package newrelic

import "go.uber.org/zap/zapcore"

// Config defines the available configuration options for New Relic log forwarding
type Config struct {
	// LicenseKey is the New Relic license key for authentication
	LicenseKey string

	// ReportableLevels specifies which log levels to forward to New Relic.
	// If nil, defaults to all levels.
	ReportableLevels []zapcore.Level
}

package sentry

import "go.uber.org/zap/zapcore"

// Config defines the available configuration options for Sentry log forwarding
type Config struct {
	// DSN is the Sentry Data Source Name for authentication
	DSN string

	// Environment specifies the environment (e.g., "production", "staging")
	Environment string

	// Release specifies the release version
	Release string

	// Debug enables debug mode for the Sentry SDK
	Debug bool

	// EventLevels specifies which levels create Sentry Events (show as issues).
	// Defaults to Error and Fatal if not specified.
	EventLevels []zapcore.Level

	// LogLevels specifies which levels are sent as Log entries (structured log data).
	// Defaults to Debug, Info, Warn, Error, Fatal if not specified.
	LogLevels []zapcore.Level
}

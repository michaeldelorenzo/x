package xlog

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/michaeldelorenzo/x/pkg/instrumenting/xlog/newrelic"
	"github.com/michaeldelorenzo/x/pkg/instrumenting/xlog/noop"
	"github.com/michaeldelorenzo/x/pkg/instrumenting/xlog/sentry"
	"github.com/michaeldelorenzo/x/pkg/instrumenting/xlog/types"
)

// isInit maintains logger init state
var isInit bool

// config maintains a valid configuration state
var config *Config

// NewRelicClient maintains the authorization token for NR
// Deprecated: Use newrelic.Config instead
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

// Config provides shared configuration options for our logger
// and our internal zapcore.Core implementation. IsAsync is only
// evaluated if LogUpstream is set to true.
type Config struct {
	// AppName will be used as both the service and logger name
	AppName string
	Level   string
	GitHash string
	Version string

	// Provider specifies which log provider to use (if any)
	Provider types.ProviderType

	// IsAsync determines whether upstream log entries should be processed
	// in an asynchronous manner. Only evaluated if LogUpstream is true.
	IsAsync bool

	// LogUpstream determines whether log entries should be sent upstream
	LogUpstream bool

	// ReportableLevels specifies which log levels should be forwarded
	// If nil, all levels are forwarded (filtered by provider)
	ReportableLevels []zapcore.Level

	// NewRelic specific configuration
	NewRelic *newrelic.Config

	// Sentry specific configuration
	Sentry *sentry.Config

	// Deprecated: NRClient is deprecated. Use NewRelic field instead.
	NRClient *NewRelicClient
}

// setLogger sets our global logger based on the provided configuration
func setLogger(conf *Config) *zap.SugaredLogger {
	initOnce.Do(func() {
		// Set our global config state
		config = conf

		// Initialize the log provider if LogUpstream is enabled
		var provider types.LogProvider
		if conf.LogUpstream {
			provider = initProvider(conf)
		}

		// Our default core for stdout writes
		consoleCore := getStdOutCore(conf)

		// If LogUpstream is true and we have a valid provider, use multicore
		if provider != nil && provider.IsValid() {
			// Our custom upstream core which encodes upstream
			upstreamCore := getUpstreamCore(conf, provider)
			// Our multicore core which encodes to both console & upstream
			core := zapcore.NewTee(consoleCore, upstreamCore)
			// Our logger, initialized with our multicore core
			logger = zap.New(core).Named(conf.AppName).Sugar()
		} else {
			// Otherwise simply use our console core
			logger = zap.New(consoleCore).Named(conf.AppName).Sugar()
		}

		// Set our base fields
		logger = logger.With(
			"service",
			conf.AppName,
			"hostname",
			hostname,
			"version",
			conf.Version,
			"commit_sha",
			conf.GitHash,
		)

		isInit = true
	})

	//: If our log level is SilentLevel, we will append a custom hook to our logger.
	if logLevelFromConfig(conf).String() == CustomLevelString(SilentLevel) {
		logger = logger.Desugar().WithOptions(getSilentHook()).Sugar()
	}

	// Returns a copy of our global logger
	return logger.With()
}

// initProvider initializes the log provider based on the configuration
func initProvider(conf *Config) types.LogProvider {
	// Handle deprecated NRClient field for backward compatibility
	if conf.NRClient != nil && conf.NRClient.isValid() {
		log.Println("WARN: NRClient is deprecated. Please use NewRelic config field instead.")
		if conf.NewRelic == nil {
			conf.NewRelic = &newrelic.Config{
				LicenseKey:       conf.NRClient.LicenseKey,
				ReportableLevels: conf.ReportableLevels,
			}
			conf.Provider = types.ProviderNewRelic
		}
	}

	switch conf.Provider {
	case types.ProviderNewRelic:
		if conf.NewRelic == nil {
			log.Println("WARN: NewRelic provider selected but config is nil")
			return nil
		}
		// Override ReportableLevels if specified at top level
		if conf.ReportableLevels != nil {
			conf.NewRelic.ReportableLevels = conf.ReportableLevels
		}
		return newrelic.NewProvider(conf.NewRelic)

	case types.ProviderSentry:
		if conf.Sentry == nil {
			log.Println("WARN: Sentry provider selected but config is nil")
			return nil
		}
		return sentry.NewProvider(conf.Sentry)

	case types.ProviderNoop:
		return noop.NewProvider()

	default:
		log.Printf("WARN: unsupported log provider: %s", conf.Provider)
		return nil
	}
}

// getUpstreamCore returns an initialized instance of our
// custom zapcore.Core
func getUpstreamCore(conf *Config, provider types.LogProvider) zapcore.Core {
	level := logLevelFromConfig(conf)

	// Since our upstream core is complimented by a console core
	// in setLogger, we set Out to a no-op sync.
	return &xlogCore{
		Async:        conf.IsAsync,
		Provider:     provider,
		LevelEnabler: zap.NewAtomicLevelAt(level),
		enc:          zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		out:          zapcore.AddSync(ioutil.Discard),
	}
}

// getStdOutCore returns an initialized Core with a console encoder
func getStdOutCore(conf *Config) zapcore.Core {
	level := zap.NewAtomicLevelAt(logLevelFromConfig(conf))
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleDebugging := zapcore.Lock(os.Stdout)
	return zapcore.NewCore(encoder, consoleDebugging, level)
}

// logLevelFromConfig fetches the log level from the provided logger Config
func logLevelFromConfig(conf *Config) zapcore.Level {
	switch strings.ToLower(conf.Level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "silent":
		return SilentLevel
	default:
		return zapcore.InfoLevel
	}
}

func getSilentHook() zap.Option {
	return zap.Hooks(func(entry zapcore.Entry) error {
		return silentErr{}
	})
}

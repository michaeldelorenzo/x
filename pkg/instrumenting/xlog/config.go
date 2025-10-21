package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"io/ioutil"
	"os"
	"strings"
)

// isInit maintains logger init state
var isInit bool

// config maintains a valid configuration state
var config *Config

// Config provides shared configuration options for our logger
// and our internal zapcore.Core implementation. IsAsync is only
// evaluated if LogUpstream is set to true.
type Config struct {
	// AppName will be used as both the service and logger name
	AppName string
	Level   string
	GitHash string
	Version string

	// IsAsync determines whether upstream log entries should be processed
	// in an asynchronous manner. This is only evaluated if LogUpstream is true.
	IsAsync bool

	// LogUpstream determines whether log entries should be sent to New Relic.
	LogUpstream bool

	NRClient *NewRelicClient
}

// setLogger sets our global logger based on the provided configuration
func setLogger(conf *Config) *zap.SugaredLogger {
	initOnce.Do(func() {
		// Set our global config state
		config = conf

		// Our default core for stdout writes
		consoleCore := getStdOutCore(conf)

		// If LogUpstream is true, fetch a multicore core...console and upstream
		if conf.LogUpstream && conf.NRClient.isValid() {
			// Our custom upstream core which encodes upstream
			upstreamCore := getUpstreamCore(conf)
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

// getUpstreamCore returns an initialized instance of our
// custom zapcore.Core
func getUpstreamCore(conf *Config) zapcore.Core {
	level := logLevelFromConfig(conf)

	// Since our upstream core is complimented by a console core
	// in setLogger, we set Out to a no-op sync.
	return &xlogCore{
		Async:            conf.IsAsync,
		ReportableLevels: LevelThreshold(level),
		LevelEnabler:     zap.NewAtomicLevelAt(level),
		enc:              zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		out:              zapcore.AddSync(ioutil.Discard),
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

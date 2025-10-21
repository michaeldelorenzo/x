package xlog

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// xlogContextKey is an internal key type
// for storing loggers in context.Context
type xlogContextKey string

// CustomLevel is a backwards compatible type alias
// extension for log levels not supported by zapcore.Level
type CustomLevel = zapcore.Level

const (
	// SilentLevel is typically used when running unit
	// tests to ensure tests are not overly verbose. We
	// increment the iota level for zap compatability.
	SilentLevel CustomLevel = iota + 6
)

// CustomLevelString returns the string representation of our
// custom level, since unfortunately, since we cannot implement the
// Stringer interface on CustomLevel since it is aliased from a non-local type.
func CustomLevelString(level CustomLevel) string {
	switch level {
	case SilentLevel:
		return "silent"
	default:
		return "info"
	}
}

// silentErr implements Error interface.
type silentErr struct{}

// Error returns error message for silentErr.
// Invoked implicitly by silent level hook.
func (e silentErr) Error() string {
	return "log output silenced."
}

var (
	// logger is our global logger
	logger *zap.SugaredLogger

	// LoggerKey avoids context name collisions due to its protected nature.
	LoggerKey = xlogContextKey("xlog_logger")

	//Once registers to ensure we don't allocate unnecessary resources
	initOnce     sync.Once
	defaultOnce  sync.Once
	hostnameOnce sync.Once
)

// AllLevels is a collection of all supported log levels
var AllLevels = []zapcore.Level{
	zapcore.DebugLevel,
	zapcore.InfoLevel,
	zapcore.WarnLevel,
	zapcore.ErrorLevel,
	zapcore.FatalLevel,
	zapcore.PanicLevel,
	SilentLevel,
}

func init() {
	// Set our reusable hostname once
	setHostName()
}

// Init returns an initialized zap.SugaredLogger with
// a custom core attached if upstreaming logging is enabled.
func Init(conf *Config) *zap.SugaredLogger {
	return setLogger(conf)
}

// L fetches a copy of the previously initialized zap.SugaredLogger.
// If the global logger has not been initialized, a sensibly
// configured logger is initialized and a copy returned.
func L() *zap.SugaredLogger {
	if !isInit {
		return getDefaultLogger()
	}

	return logger.With()
}

// LL creates a contextual logger using the provided fields, and embeds that
// contextual logger in the provided context. Returns both the context and logger.
func LL(ctx context.Context, fields Fields) (context.Context, *zap.SugaredLogger) {
	l := LoggerFromCtx(ctx)
	for k, v := range fields {
		l = l.With(k, v)
	}

	return context.WithValue(ctx, LoggerKey, l), l
}

// LoggerFromCtx attempts to retrieve a *zap.SugaredLogger from the provided context. Returns
// an initialized logger, even if no logger was found in the provided context, so some expected
// context might be lost from returned logger if a logger was not found in the provided context.
func LoggerFromCtx(ctx context.Context) *zap.SugaredLogger {
	t := ctx.Value(LoggerKey)
	switch logger := t.(type) {
	case *zap.SugaredLogger:
		return logger
	default:
		return L()
	}
}

// EmbedLogger returns a context with the provided logger set as a value
func EmbedLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// Fields stores contextual key-value pairs.
type Fields map[string]interface{}

// ErrWith logs provided error data and fields.
// Intended to be used with contextual, structured logging:
//
//	xlog.ErrWith(errors.New("failed to send event"), xlog.Fields{
//		"topic": topic,
//		"event": event,
//	})
func ErrWith(logger *zap.SugaredLogger, title interface{}, fields Fields) {
	if logger == nil {
		logger = L().With()
	}

	for k, v := range fields {
		logger = logger.With(k, v)
	}

	logger.Error(errMsg(title))
}

// Err logs the provided error data using the configured logger.
func Err(logger *zap.SugaredLogger, title string, message interface{}) {
	if logger == nil {
		logger = L().With()
	}

	msg := errMsg(message)

	logger.Error(fmt.Sprintf("%s: %s", title, msg))
}

// EEWith is a convenience function which will log the provided
// error data & fields, and if the provided transaction is not nil,
// register the error data as a handled transaction error.
func EEWith(logger *zap.SugaredLogger, title string, fields Fields, txn *newrelic.Transaction) {
	if logger == nil {
		logger = L().With()
	}

	for k, v := range fields {
		logger = logger.With(k, v)
	}
	msg := errMsg(title)

	// If a valid new relic transaction is provided, we will have new relic notice the error.
	if txn != nil {
		fields["hostname"] = hostname
		fields["trace"] = txn.GetLinkingMetadata().TraceID
		fields["txn"] = txn.GetLinkingMetadata().EntityGUID
		logger = logger.With("trace", txn.GetLinkingMetadata().TraceID, "txn", txn.GetLinkingMetadata().EntityGUID)

		txn.NoticeError(newrelic.Error{
			Message:    msg,
			Class:      trace(),
			Attributes: fields,
			Stack:      newrelic.NewStackTrace(),
		})
	}

	logger.Error(errMsg(title))
}

// EE is a convenience function which will log the provided
// error, and if the provided transaction is not nil,
// register the error data as a handled transaction error.
func EE(logger *zap.SugaredLogger, error interface{}, txn *newrelic.Transaction) {
	msg := errMsg(error)
	logger.Error(msg)

	// If a valid new relic transaction is provided, we will have new relic notice the error.
	if txn != nil {
		txn.NoticeError(newrelic.Error{
			Message: msg,
			Class:   trace(),
			Attributes: Fields{
				"trace":    txn.GetLinkingMetadata().TraceID,
				"txn":      txn.GetLinkingMetadata().EntityGUID,
				"hostname": hostname,
			},
			Stack: newrelic.NewStackTrace(),
		})
	}
}

// errMsg plucks the error message based in the err type.
func errMsg(err interface{}) string {
	var msg string
	switch m := err.(type) {
	case error:
		msg = m.Error()
	default:
		msg = fmt.Sprintf("%+v", m)
	}

	return msg
}

// xlogCore is our internal implementation of zapcore.Core
type xlogCore struct {
	ReportableLevels []zapcore.Level
	Async            bool

	m sync.Mutex
	zapcore.LevelEnabler
	enc zapcore.Encoder
	out zapcore.WriteSyncer
}

// With is an overload of zapcore.Core With method for our custom
// core implementation.
func (c *xlogCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc, fields)
	return clone
}

// Check is an overload of zapcore.Core Check method for our custom
// core implementation.
func (c *xlogCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

// Write is an overload of zapcore.Core Write method for our custom
// core implementation.
func (c *xlogCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}

	// Inject our upstream (new relic) writer
	err = c.newRelicWriter(ent, buf.String())
	buf.Free()

	if err != nil {
		return err
	}

	if ent.Level > zapcore.ErrorLevel {
		// Since the program may be crashing, sync the output. Sync errors can be safely ignored.
		c.Sync()
	}

	return nil
}

// Sync is an overload of zapcore.Core Sync method for our custom
// core implementation.
func (c *xlogCore) Sync() error {
	return c.out.Sync()
}

// LevelThreshold returns the given parameter level
// and every logging level above it.
func LevelThreshold(l zapcore.Level) []zapcore.Level {
	for i := range AllLevels {
		if AllLevels[i] == l {
			return AllLevels[i:]
		}
	}
	return []zapcore.Level{}
}

// clone creates a referenced copy of our base xlogCore receiver.
// This is necessary for context logging when calling With.
func (c *xlogCore) clone() *xlogCore {
	return &xlogCore{
		Async:            c.Async,
		ReportableLevels: c.ReportableLevels,
		LevelEnabler:     c.LevelEnabler,
		enc:              c.enc.Clone(),
		out:              c.out,
	}
}

// addFields copies our logger fields to the specified Encoder
func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

// sendEvent publishes the specified log entry to new relic.
// An error will be returned
func (c *xlogCore) sendEvent(e *zapcore.Entry, msg string) error {
	c.m.Lock()
	defer c.m.Unlock()

	return config.NRClient.request([]byte(msg))
}

// Levels returns which levels are to be sent to New Relic
func (c *xlogCore) Levels() []zapcore.Level {
	if c.ReportableLevels == nil {
		return AllLevels
	}

	return c.ReportableLevels
}

// isAcceptedLevel returns whether the specified level is valid
func (c *xlogCore) isAcceptedLevel(level zapcore.Level) bool {
	for _, lv := range c.Levels() {
		if lv == level {
			return true
		}
	}
	return false
}

// newRelicWriter determines whether an event can and should
// be sent upstream to new relic based on the provided level
func (c *xlogCore) newRelicWriter(e zapcore.Entry, msg string) error {
	// Check if we should send this entry to new relic based
	// on the configured level
	if !c.isAcceptedLevel(e.Level) {
		return nil
	}

	// If we use Async upstream writing, we sacrifice having an error returned
	// since channel usage seemed excessive here. The sendEvent function will output
	// the WARN to stdout since failing to write upstream is not considered a fatal error.
	if c.Async {
		go c.sendEvent(&e, msg)
		return nil
	}

	return c.sendEvent(&e, msg)
}

// getDefaultLogger returns a sensible default logger using
// zap's production configuration
func getDefaultLogger() *zap.SugaredLogger {
	defaultOnce.Do(func() {
		logConfig := zap.NewProductionConfig()
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

		// Using the config we generated above, build the actual Logger instance.
		l, err := logConfig.Build(
			zap.Fields(
				zap.Bool("usingDefaultLogger", true),
			),
		)

		// If any errors are encountered, panic the program. This should never
		// happen in the normal course of operation.
		if err != nil {
			panic(err.Error())
		}

		logger = l.Sugar()
	})

	return logger.With()
}

// trace builds a readable string with a stack caller info for our error register
func trace() string {
	pc := make([]uintptr, 15)

	// We currently skip 3 stack frames since it
	// seems like the most sensible error class name
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:L%d", frame.Function, frame.Line)
}

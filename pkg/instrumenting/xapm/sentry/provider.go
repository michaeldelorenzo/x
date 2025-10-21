package sentry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/michaeldelorenzo/x/pkg/instrumenting/xapm/types"
)

// Provider implements the xapm.Provider interface for Sentry
type Provider struct {
	hub *sentry.Hub
}

// NewProvider creates a new Sentry APM provider
func NewProvider(hub *sentry.Hub) *Provider {
	return &Provider{hub: hub}
}

// StartTransaction creates a new Sentry transaction
func (p *Provider) StartTransaction(name string) types.Transaction {
	ctx := context.Background()
	span := sentry.StartSpan(ctx, "transaction", sentry.WithTransactionName(name))

	return &Transaction{
		span: span,
		hub:  p.hub,
	}
}

// RecordCustomEvent records a custom event in Sentry
func (p *Provider) RecordCustomEvent(eventType string, params map[string]interface{}) {
	event := sentry.NewEvent()
	event.Message = eventType
	event.Extra = params
	event.Level = sentry.LevelInfo

	p.hub.CaptureEvent(event)
}

// Shutdown gracefully shuts down the Sentry client
func (p *Provider) Shutdown(timeout time.Duration) {
	sentry.Flush(timeout)
}

// Type returns the provider type
func (p *Provider) Type() types.ProviderType {
	return types.ProviderSentry
}

// Transaction wraps a Sentry span (which represents a transaction)
type Transaction struct {
	span *sentry.Span
	hub  *sentry.Hub
}

// End completes the transaction
func (t *Transaction) End() {
	t.span.Finish()
}

// StartSegment begins a new timing segment
func (t *Transaction) StartSegment(name string) types.Segment {
	child := t.span.StartChild(name)
	return &Segment{span: child}
}

// StartExternalSegment begins tracking an external HTTP call
func (t *Transaction) StartExternalSegment(url string) types.Segment {
	child := t.span.StartChild("http.client")
	child.Description = url
	child.SetData("url", url)
	return &Segment{span: child}
}

// StartDBSegment begins tracking a database operation
func (t *Transaction) StartDBSegment(params *types.DBSegParams) types.Segment {
	op := fmt.Sprintf("db.%s", params.Operation)
	child := t.span.StartChild(op)
	child.Description = params.QueryString

	child.SetData("db.system", params.DatabaseType)
	child.SetData("db.name", params.DatabaseName)
	child.SetData("server.address", params.Host)
	if params.Collection != "" {
		child.SetData("db.collection.name", params.Collection)
	}
	if params.QueryString != "" {
		child.SetData("db.query", params.QueryString)
	}

	// Add query parameters as separate data fields
	for k, v := range params.QueryParams {
		child.SetData(fmt.Sprintf("db.param.%s", k), v)
	}

	return &Segment{span: child}
}

// StartPSQLSegment begins tracking a PostgreSQL operation
func (t *Transaction) StartPSQLSegment(params *types.PSQLDBSegParams) types.Segment {
	child := t.span.StartChild("db.query")
	child.Description = params.QueryString

	child.SetData("db.system", "postgresql")
	child.SetData("db.name", params.DatabaseName)
	child.SetData("server.address", params.Host)
	if params.QueryString != "" {
		child.SetData("db.query", params.QueryString)
	}

	// Add query parameters as separate data fields
	for k, v := range params.QueryParams {
		child.SetData(fmt.Sprintf("db.param.%s", k), v)
	}

	return &Segment{span: child}
}

// StartKafkaSegment begins tracking a Kafka message production
func (t *Transaction) StartKafkaSegment(topicName string) types.Segment {
	child := t.span.StartChild("message.publish")
	child.Description = topicName

	child.SetData("messaging.system", "kafka")
	child.SetData("messaging.destination.name", topicName)
	child.SetData("messaging.operation.type", "publish")

	return &Segment{span: child}
}

// SetWebRequestHTTP marks this transaction as a web transaction and captures request details
func (t *Transaction) SetWebRequestHTTP(req *http.Request) {
	t.span.SetData("http.method", req.Method)
	t.span.SetData("http.url", req.URL.String())
	t.span.SetData("http.route", req.URL.Path)
	t.span.SetData("http.scheme", req.URL.Scheme)
	t.span.SetData("server.address", req.Host)

	// Add headers if needed (be careful with sensitive data)
	if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
		t.span.SetData("http.user_agent", userAgent)
	}
}

// SetWebResponse wraps the response writer to capture response details
func (t *Transaction) SetWebResponse(w http.ResponseWriter) http.ResponseWriter {
	// Sentry's approach to capturing response status is different
	// We'll wrap the ResponseWriter to capture the status code
	return &responseWriter{
		ResponseWriter: w,
		span:           t.span,
	}
}

// NoticeError records an error in this transaction
func (t *Transaction) NoticeError(err error) {
	t.hub.CaptureException(err)
	t.span.Status = sentry.SpanStatusInternalError
}

// AddAttribute adds a custom attribute to the transaction
func (t *Transaction) AddAttribute(key string, value interface{}) {
	t.span.SetData(key, value)
}

// Context returns a context with this transaction embedded
func (t *Transaction) Context(ctx context.Context) context.Context {
	return t.span.Context()
}

// Segment wraps a Sentry span (used for segments within a transaction)
type Segment struct {
	span *sentry.Span
}

// End completes the segment
func (s *Segment) End() {
	s.span.Finish()
}

// AddAttribute adds a custom attribute to the segment
func (s *Segment) AddAttribute(key string, value interface{}) {
	s.span.SetData(key, value)
}

// responseWriter wraps http.ResponseWriter to capture status codes
type responseWriter struct {
	http.ResponseWriter
	span         *sentry.Span
	statusCode   int
	wroteHeader  bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.statusCode = code
		rw.span.SetData("http.response.status_code", code)
		rw.span.Status = httpStatusToSentryStatus(code)
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// httpStatusToSentryStatus converts HTTP status codes to Sentry span statuses
func httpStatusToSentryStatus(code int) sentry.SpanStatus {
	switch {
	case code >= 200 && code < 300:
		return sentry.SpanStatusOK
	case code >= 400 && code < 500:
		if code == 404 {
			return sentry.SpanStatusNotFound
		}
		if code == 403 {
			return sentry.SpanStatusPermissionDenied
		}
		if code == 401 {
			return sentry.SpanStatusUnauthenticated
		}
		return sentry.SpanStatusInvalidArgument
	case code >= 500:
		return sentry.SpanStatusInternalError
	default:
		return sentry.SpanStatusUnknown
	}
}

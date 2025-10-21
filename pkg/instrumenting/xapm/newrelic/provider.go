package newrelic

import (
	"context"
	"net/http"
	"time"

	"github.com/michaeldelorenzo/x/v3/pkg/instrumenting/xapm/types"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/go-agent/v3/newrelic/sqlparse"
)

// Provider implements the xapm.Provider interface for New Relic
type Provider struct {
	app *newrelic.Application
}

// NewProvider creates a new New Relic APM provider
func NewProvider(app *newrelic.Application) *Provider {
	return &Provider{app: app}
}

// StartTransaction creates a new New Relic transaction
func (p *Provider) StartTransaction(name string) types.Transaction {
	return &Transaction{
		tx: p.app.StartTransaction(name),
	}
}

// RecordCustomEvent records a custom event in New Relic
func (p *Provider) RecordCustomEvent(eventType string, params map[string]interface{}) {
	p.app.RecordCustomEvent(eventType, params)
}

// Shutdown gracefully shuts down the New Relic agent
func (p *Provider) Shutdown(timeout time.Duration) {
	p.app.Shutdown(timeout)
}

// Type returns the provider type
func (p *Provider) Type() types.ProviderType {
	return types.ProviderNewRelic
}

// Transaction wraps a New Relic transaction
type Transaction struct {
	tx *newrelic.Transaction
}

// End completes the transaction
func (t *Transaction) End() {
	t.tx.End()
}

// StartSegment begins a new timing segment
func (t *Transaction) StartSegment(name string) types.Segment {
	return &Segment{
		seg: &newrelic.Segment{
			Name:      name,
			StartTime: t.tx.StartSegmentNow(),
		},
	}
}

// StartExternalSegment begins tracking an external HTTP call
func (t *Transaction) StartExternalSegment(url string) types.Segment {
	return &Segment{
		seg: &newrelic.ExternalSegment{
			StartTime: t.tx.StartSegmentNow(),
			URL:       url,
		},
	}
}

// StartDBSegment begins tracking a database operation
func (t *Transaction) StartDBSegment(params *types.DBSegParams) types.Segment {
	return &Segment{
		seg: &newrelic.DatastoreSegment{
			StartTime:          t.tx.StartSegmentNow(),
			Host:               params.Host,
			DatabaseName:       params.DatabaseName,
			Collection:         params.Collection,
			Operation:          string(params.Operation),
			Product:            newrelic.DatastoreProduct(params.DatabaseType),
			ParameterizedQuery: params.QueryString,
			QueryParameters:    sanitizeQueryParams(params.QueryParams),
		},
	}
}

// StartPSQLSegment begins tracking a PostgreSQL operation with query parsing
func (t *Transaction) StartPSQLSegment(params *types.PSQLDBSegParams) types.Segment {
	seg := &newrelic.DatastoreSegment{
		Host:               params.Host,
		DatabaseName:       params.DatabaseName,
		Product:            newrelic.DatastorePostgres,
		ParameterizedQuery: params.QueryString,
		QueryParameters:    sanitizeQueryParams(params.QueryParams),
	}

	// mutates the segment and sets the `operation` and `collection` by analyzing the raw query
	sqlparse.ParseQuery(seg, params.QueryString)

	seg.StartTime = t.tx.StartSegmentNow()

	return &Segment{seg: seg}
}

// StartKafkaSegment begins tracking a Kafka message production
func (t *Transaction) StartKafkaSegment(topicName string) types.Segment {
	return &Segment{
		seg: &newrelic.MessageProducerSegment{
			StartTime:       t.tx.StartSegmentNow(),
			Library:         "kafka",
			DestinationType: newrelic.MessageTopic,
			DestinationName: topicName,
		},
	}
}

// SetWebRequestHTTP marks this transaction as a web transaction and captures request details
func (t *Transaction) SetWebRequestHTTP(req *http.Request) {
	t.tx.SetWebRequestHTTP(req)
}

// SetWebResponse wraps the response writer to capture response details
func (t *Transaction) SetWebResponse(w http.ResponseWriter) http.ResponseWriter {
	return t.tx.SetWebResponse(w)
}

// NoticeError records an error in this transaction
func (t *Transaction) NoticeError(err error) {
	t.tx.NoticeError(err)
}

// AddAttribute adds a custom attribute to the transaction
func (t *Transaction) AddAttribute(key string, value interface{}) {
	t.tx.AddAttribute(key, value)
}

// Context returns a context with this transaction embedded
func (t *Transaction) Context(ctx context.Context) context.Context {
	return newrelic.NewContext(ctx, t.tx)
}

// Segment wraps a New Relic segment
type Segment struct {
	seg interface{ End() }
}

// End completes the segment
func (s *Segment) End() {
	s.seg.End()
}

// AddAttribute adds a custom attribute to the segment
func (s *Segment) AddAttribute(key string, value interface{}) {
	// New Relic segments don't support attributes directly
	// This is a no-op for compatibility
}

// sanitizeQueryParams is a helper that's shared with the main package
func sanitizeQueryParams(params map[string]interface{}) map[string]interface{} {
	// For now, we'll just return the params as-is
	// The main xapm package has a more sophisticated implementation
	// that we could expose if needed
	result := make(map[string]interface{}, len(params))
	for k, v := range params {
		result[k] = v
	}
	return result
}

package noop

import (
	"context"
	"net/http"
	"time"

	"github.com/michaeldelorenzo/x/pkg/instrumenting/xapm/types"
)

// Provider implements a no-op APM provider for testing or when APM is disabled
type Provider struct{}

// NewProvider creates a new no-op APM provider
func NewProvider() *Provider {
	return &Provider{}
}

// StartTransaction creates a no-op transaction
func (p *Provider) StartTransaction(name string) types.Transaction {
	return &Transaction{}
}

// RecordCustomEvent does nothing in the no-op provider
func (p *Provider) RecordCustomEvent(eventType string, params map[string]interface{}) {
	// no-op
}

// Shutdown does nothing in the no-op provider
func (p *Provider) Shutdown(timeout time.Duration) {
	// no-op
}

// Type returns the provider type
func (p *Provider) Type() types.ProviderType {
	return types.ProviderNoop
}

// Transaction is a no-op transaction implementation
type Transaction struct{}

// End does nothing
func (t *Transaction) End() {}

// StartSegment returns a no-op segment
func (t *Transaction) StartSegment(name string) types.Segment {
	return &Segment{}
}

// StartExternalSegment returns a no-op segment
func (t *Transaction) StartExternalSegment(url string) types.Segment {
	return &Segment{}
}

// StartDBSegment returns a no-op segment
func (t *Transaction) StartDBSegment(params *types.DBSegParams) types.Segment {
	return &Segment{}
}

// StartPSQLSegment returns a no-op segment
func (t *Transaction) StartPSQLSegment(params *types.PSQLDBSegParams) types.Segment {
	return &Segment{}
}

// StartKafkaSegment returns a no-op segment
func (t *Transaction) StartKafkaSegment(topicName string) types.Segment {
	return &Segment{}
}

// SetWebRequestHTTP does nothing
func (t *Transaction) SetWebRequestHTTP(req *http.Request) {}

// SetWebResponse returns the original writer unchanged
func (t *Transaction) SetWebResponse(w http.ResponseWriter) http.ResponseWriter {
	return w
}

// NoticeError does nothing
func (t *Transaction) NoticeError(err error) {}

// AddAttribute does nothing
func (t *Transaction) AddAttribute(key string, value interface{}) {}

// Context returns the original context unchanged
func (t *Transaction) Context(ctx context.Context) context.Context {
	return ctx
}

// Segment is a no-op segment implementation
type Segment struct{}

// End does nothing
func (s *Segment) End() {}

// AddAttribute does nothing
func (s *Segment) AddAttribute(key string, value interface{}) {}

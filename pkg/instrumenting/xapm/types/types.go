package types

import (
	"context"
	"net/http"
	"time"
)

// ProviderType represents the APM provider type
type ProviderType string

const (
	ProviderNewRelic ProviderType = "newrelic"
	ProviderSentry   ProviderType = "sentry"
	ProviderNoop     ProviderType = "noop"
)

// Provider defines the interface that all APM providers must implement
type Provider interface {
	// StartTransaction creates a new transaction with the given name
	StartTransaction(name string) Transaction

	// RecordCustomEvent records a custom event with the given type and attributes
	RecordCustomEvent(eventType string, params map[string]interface{})

	// Shutdown gracefully shuts down the provider
	Shutdown(timeout time.Duration)

	// Type returns the provider type
	Type() ProviderType
}

// Transaction represents a single unit of work (e.g., HTTP request, background job)
type Transaction interface {
	// End completes the transaction
	End()

	// StartSegment begins a new timing segment
	StartSegment(name string) Segment

	// StartExternalSegment begins tracking an external HTTP call
	StartExternalSegment(url string) Segment

	// StartDBSegment begins tracking a database operation
	StartDBSegment(params *DBSegParams) Segment

	// StartPSQLSegment begins tracking a PostgreSQL operation with query parsing
	StartPSQLSegment(params *PSQLDBSegParams) Segment

	// StartKafkaSegment begins tracking a Kafka message production
	StartKafkaSegment(topicName string) Segment

	// SetWebRequestHTTP marks this transaction as a web transaction and captures request details
	SetWebRequestHTTP(req *http.Request)

	// SetWebResponse wraps the response writer to capture response details
	SetWebResponse(w http.ResponseWriter) http.ResponseWriter

	// NoticeError records an error in this transaction
	NoticeError(err error)

	// AddAttribute adds a custom attribute to the transaction
	AddAttribute(key string, value interface{})

	// Context returns a context with this transaction embedded
	Context(ctx context.Context) context.Context
}

// Segment represents a timing measurement within a transaction
type Segment interface {
	// End completes the segment
	End()

	// AddAttribute adds a custom attribute to the segment
	AddAttribute(key string, value interface{})
}

// DBSegParams contains parameters for database segments
type DBSegParams struct {
	Host         string
	DatabaseName string
	Collection   string
	Operation    DBOperation
	DatabaseType DBType
	QueryString  string
	QueryParams  map[string]interface{}
}

// PSQLDBSegParams contains parameters for PostgreSQL segments
type PSQLDBSegParams struct {
	Host         string
	DatabaseName string
	QueryString  string
	QueryParams  map[string]interface{}
}

// DBType represents a database product type
type DBType string

const (
	DBRedis    DBType = "Redis"
	DBMongoDB  DBType = "MongoDB"
	DBPostgres DBType = "Postgres"
)

// DBOperation represents a database operation type
type DBOperation string

const (
	DBSelect DBOperation = "SELECT"
	DBInsert DBOperation = "INSERT"
	DBUpdate DBOperation = "UPDATE"
	DBUpsert DBOperation = "UPSERT"
	DBDelete DBOperation = "DELETE"
	DBCount  DBOperation = "COUNT"
)

// EventReporter is an interface that allows the receiver model to be processed as a custom event.
type EventReporter interface {
	GetEventType() string
}

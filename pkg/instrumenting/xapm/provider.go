package xapm

// Re-export types from the types package to maintain backward compatibility
// and provide a convenient import path for users

import "github.com/michaeldelorenzo/x/v2/pkg/instrumenting/xapm/types"

type (
	// ProviderType represents the APM provider type
	ProviderType = types.ProviderType

	// Provider defines the interface that all APM providers must implement
	Provider = types.Provider

	// Transaction represents a single unit of work
	Transaction = types.Transaction

	// Segment represents a timing measurement within a transaction
	Segment = types.Segment

	// DBSegParams contains parameters for database segments
	DBSegParams = types.DBSegParams

	// PSQLDBSegParams contains parameters for PostgreSQL segments
	PSQLDBSegParams = types.PSQLDBSegParams

	// DBType represents a database product type
	DBType = types.DBType

	// DBOperation represents a database operation type
	DBOperation = types.DBOperation

	// EventReporter is an interface for custom events
	EventReporter = types.EventReporter
)

const (
	ProviderNewRelic = types.ProviderNewRelic
	ProviderSentry   = types.ProviderSentry
	ProviderNoop     = types.ProviderNoop

	DBRedis    = types.DBRedis
	DBMongoDB  = types.DBMongoDB
	DBPostgres = types.DBPostgres

	DBSelect = types.DBSelect
	DBInsert = types.DBInsert
	DBUpdate = types.DBUpdate
	DBUpsert = types.DBUpsert
	DBDelete = types.DBDelete
	DBCount  = types.DBCount
)

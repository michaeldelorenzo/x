package mockpgdb

import (
	"context"
	"database/sql"

	"github.com/michaeldelorenzo/x/v3/pkg/datastores/pg/db"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Mock mocks the pgdb
type Mock struct {
	QueryError error
	PingError  error
}

// Connect is required for the PGDB interface
func (m Mock) Connect(_ string) (db.PGDBConnection, error) {
	return &MockConnection{
		queryError: m.QueryError,
		pingError:  m.PingError,
	}, nil
}

// MockConnection mocks the pgdb connection
type MockConnection struct {
	queryError error
	pingError  error
}

func (c MockConnection) GetSqlxTx() *sqlx.Tx {
	return &sqlx.Tx{}
}

// Ping is required for the PGDBConnection interface
func (c MockConnection) Ping(ctx context.Context) error {
	return c.pingError
}

// Close is required for the PGDBConnection interface
func (c MockConnection) Close() error {
	return nil
}

// GetContext is required for the PGDBConnection interface
func (c *MockConnection) GetContext(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return nil
}

// QueryContext is required for the PGDBConnection interface
func (c *MockConnection) QueryContext(_ context.Context, _ string, _ ...interface{}) (*sqlx.Rows, error) {
	if c.queryError == nil {
		return &sqlx.Rows{
			Rows: &sql.Rows{},
		}, nil
	}

	return nil, c.queryError
}

// NamedQueryContext is required for the PGDBConnection interface
func (c *MockConnection) NamedQueryContext(_ context.Context, _ string, _ interface{}) (*sqlx.Rows, error) {
	if c.queryError == nil {
		return &sqlx.Rows{
			Rows: &sql.Rows{},
		}, nil
	}

	return nil, c.queryError
}

func (c *MockConnection) BeginTx(_ context.Context) (*sql.Tx, error) {
	return &sql.Tx{}, nil
}

// ExecContext is required for the PGDBConnection interface
func (c *MockConnection) ExecContext(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
	return nil, c.queryError
}

// NamedExecContext is required for the PGDBConnection interface
func (c *MockConnection) NamedExecContext(_ context.Context, _ string, _ interface{}) (sql.Result, error) {
	return nil, c.queryError
}

// GetTransaction is required for the PGDBConnection interface
func (c *MockConnection) GetTransaction() *sqlx.Tx {
	return c.GetSqlxTx()
}

// WithinTx is required for the PGDBConnection interface
func (c *MockConnection) WithinTx(fn func(*sqlx.Tx) error) error {
	tx := c.GetTransaction()
	return fn(tx)
}

// NewPGAuthError returns a mock authentication error
func NewPGAuthError() error {
	return &pq.Error{Code: "28P01"}
}

// NewPGUniqueViolationError returns a mock unique violation error
func NewPGUniqueViolationError() error {
	return &pq.Error{Code: "23505"}
}

package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
)

var basedb PGDB

func init() {
	SetBaseDB(PgDB{})
}

// SetBaseDB allows setting a base DB that conforms to the PGDB interface.
// This is useful for setting a mock base DB when testing.
//
// The package is initialized with the default base DB that wraps the go-pg library,
// so this method should not be used other than when testing.
func SetBaseDB(p PGDB) {
	basedb = p
}

// PGDB is an interface for wrapping methods from the sqlx library
type PGDB interface {
	Connect(string) (PGDBConnection, error)
}

type PgDB struct{}

// Connect wraps the sqlx Connect method
func (p PgDB) Connect(connectionString string) (PGDBConnection, error) {
	return basedb.Connect(connectionString)
}

// PGDBConnection is an interface for wrapping methods from the go-pg db connection
type PGDBConnection interface {
	Ping(ctx context.Context) error
	Close() error
	GetContext(ctx context.Context, destination interface{}, query string, params ...interface{}) error
	QueryContext(ctx context.Context, query string, params ...interface{}) (*sqlx.Rows, error)
	NamedQueryContext(ctx context.Context, query string, params interface{}) (*sqlx.Rows, error)
	ExecContext(ctx context.Context, query string, params ...interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, params interface{}) (sql.Result, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	GetSqlxTx() *sqlx.Tx
	GetTransaction() *sqlx.Tx
	WithinTx(fn func(*sqlx.Tx) error) error
}

type PgDBConnection struct {
	db *sqlx.DB
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (c PgDBConnection) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// GetSqlxTx creates and returns a sqlx transaction
func (c PgDBConnection) GetSqlxTx() *sqlx.Tx {
	return c.db.MustBegin()
}

// BeginTx begins a database transaction by wrapping the sql BeginTx method.
func (c PgDBConnection) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, nil)
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server
// to finish.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (c PgDBConnection) Close() error {
	return c.db.Close()
}

func (c PgDBConnection) GetContext(ctx context.Context, destination interface{}, query string, params ...interface{}) error {
	return c.db.GetContext(ctx, destination, query, params...)
}

// QueryContext wraps the sql-x QueryxContext method
// https://pkg.go.dev/github.com/jmoiron/sqlx#DB.QueryxContext
func (c PgDBConnection) QueryContext(ctx context.Context, query string, params ...interface{}) (*sqlx.Rows, error) {
	return c.db.QueryxContext(ctx, query, params...)
}

func (c PgDBConnection) NamedQueryContext(ctx context.Context, query string, params interface{}) (*sqlx.Rows, error) {
	return c.db.NamedQueryContext(ctx, query, params)
}

// ExecContext wraps the sql ExecContext method
func (c PgDBConnection) ExecContext(ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, params...)
}

func (c PgDBConnection) NamedExecContext(ctx context.Context, query string, params interface{}) (sql.Result, error) {
	return c.db.NamedExecContext(ctx, query, params)
}

// GetTransaction returns a sqlx tx
func (c PgDBConnection) GetTransaction() *sqlx.Tx {
	return c.GetSqlxTx()
}

// WithinTx wraps the provided function in a sqlx transaction
func (c PgDBConnection) WithinTx(fn func(*sqlx.Tx) error) error {
	var err error

	tx := c.GetTransaction()
	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
		} else if err != nil {
			err = tx.Rollback()
		} else {
			// err is nil; if Commit returns error update err
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

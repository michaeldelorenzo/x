package db_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/michaeldelorenzo/x/v2/pkg/datastores/pg/db"
)

func Test_PgDB_Connect(t *testing.T) {
	t.Run("should connect to the database", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func(mockDB *sql.DB) {
			_ = mockDB.Close()
		}(mockDB)

		xdb := sqlx.NewDb(mockDB, "sqlmock")

		// Mock the Connect method of the base DB

		p := &mockPGDB{conn: &mockPGDBConnection{db: xdb}}
		_, err = p.Connect("sqlmock")
		assert.NoError(t, err)
	})
}

func Test_PgDBConnection_Ping(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should ping the database", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		pingErr := c.Ping(context.Background())
		assert.NoError(t, pingErr)
	})
}

func Test_PgDBConnection_ExecContext(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should execute a query", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectExec("INSERT INTO users").WithArgs("test").WillReturnResult(sqlmock.NewResult(1, 1))

		_, pingErr := c.ExecContext(context.Background(), "INSERT INTO users", "test")
		assert.NoError(t, pingErr)
	})

	t.Run("should return an error", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectExec("INSERT INTO users").WithArgs("test").WillReturnError(errors.New("some error"))

		_, pingErr := c.ExecContext(context.Background(), "INSERT INTO users", "test")
		assert.Error(t, pingErr)
	})
}

func Test_PgDBConnection_QueryContext(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should return rows", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test")

		mock.ExpectQuery("SELECT id, name FROM users").WillReturnRows(rows)

		_, pingErr := c.QueryContext(context.Background(), "SELECT id, name FROM users")
		assert.NoError(t, pingErr)
	})

	t.Run("should return an error", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectQuery("SELECT id, name FROM users").WillReturnError(errors.New("some error"))

		_, pingErr := c.QueryContext(context.Background(), "SELECT id, name FROM users")
		assert.Error(t, pingErr)
	})
}

type mockPGDB struct {
	conn db.PGDBConnection
}

func (m *mockPGDB) Connect(string) (db.PGDBConnection, error) {
	return m.conn, nil
}

type mockPGDBConnection struct {
	db *sqlx.DB
}

func (m *mockPGDBConnection) Ping(ctx context.Context) error {
	err := m.db.PingContext(ctx)
	if err != nil {
		// Log the error, or handle it as appropriate for a mock
		return err
	}
	return nil
}

func (m *mockPGDBConnection) Close() error {
	err := m.db.Close()
	if err != nil {
		// Log the error, or handle it as appropriate for a mock
		return err
	}
	return nil
}

func (m *mockPGDBConnection) GetContext(ctx context.Context, destination interface{}, query string, params ...interface{}) error {
	err := m.db.GetContext(ctx, destination, query, params...)
	if err != nil {
		// Log the error, or handle it as appropriate for a mock
		return err
	}
	return nil
}

func (m *mockPGDBConnection) QueryContext(ctx context.Context, query string, params ...interface{}) (*sqlx.Rows, error) {
	rows, err := m.db.QueryxContext(ctx, query, params...)
	if err != nil {
		// Log the error, or handle it as appropriate for a mock
		return nil, err
	}
	return rows, nil
}

func (m *mockPGDBConnection) NamedQueryContext(ctx context.Context, query string, params interface{}) (*sqlx.Rows, error) {
	rows, err := m.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		// Log the error, or handle it as appropriate for a mock
		return nil, err
	}
	return rows, nil
}

func (m *mockPGDBConnection) NamedExecContext(ctx context.Context, query string, params interface{}) (sql.Result, error) {
	res, err := m.db.NamedExecContext(ctx, query, params)
	if err != nil {
		// Log the error, or handle it as appropriate for a mock
		return nil, err
	}
	return res, nil
}

func (m *mockPGDBConnection) ExecContext(ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	res, err := m.db.ExecContext(ctx, query, params...)
	if err != nil {
		// Log the error, or handle it as appropriate for a mock
		return nil, err
	}
	return res, nil
}

func (m *mockPGDBConnection) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		// Log the error, or handle it as appropriate for a mock
		return nil, err
	}
	return tx, nil
}

func (m *mockPGDBConnection) GetSqlxTx() *sqlx.Tx {
	return m.db.MustBegin()
}

func (m *mockPGDBConnection) GetTransaction() *sqlx.Tx {
	return m.GetSqlxTx()
}

func (m *mockPGDBConnection) WithinTx(fn func(*sqlx.Tx) error) error {
	var err error

	tx := m.GetTransaction()
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

func Test_PgDBConnection_GetContext(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should get context", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test")

		mock.ExpectQuery("SELECT id, name FROM users").WillReturnRows(rows)

		var dest struct {
			Id   int    `db:"id"`
			Name string `db:"name"`
		}
		err := c.GetContext(context.Background(), &dest, "SELECT id, name FROM users")
		assert.NoError(t, err)
	})

	t.Run("should return an error", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectQuery("SELECT id, name FROM users").WillReturnError(errors.New("some error"))

		var dest struct {
			Id   int    `db:"id"`
			Name string `db:"name"`
		}
		err := c.GetContext(context.Background(), &dest, "SELECT id, name FROM users")
		assert.Error(t, err)
	})
}

func Test_PgDBConnection_BeginTx(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should begin a transaction", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectBegin()

		_, err := c.BeginTx(context.Background())
		assert.NoError(t, err)
	})

	t.Run("should return an error", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectBegin().WillReturnError(errors.New("some error"))

		_, err := c.BeginTx(context.Background())
		assert.Error(t, err)
	})
}

func Test_PgDBConnection_GetSqlxTx(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should get a sqlx transaction", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectBegin()

		tx := c.GetSqlxTx()
		assert.NotNil(t, tx)
	})
}

func Test_PgDBConnection_GetTransaction(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should get a transaction", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectBegin()

		tx := c.GetTransaction()
		assert.NotNil(t, tx)
	})
}

func Test_PgDBConnection_WithinTx(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should execute within a transaction", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectBegin()
		mock.ExpectCommit()

		err := c.WithinTx(func(tx *sqlx.Tx) error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("should rollback on error", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectBegin()
		mock.ExpectRollback()

		err := c.WithinTx(func(tx *sqlx.Tx) error {
			return errors.New("some error")
		})
		assert.Error(t, err)
	})
}

func Test_PgDBConnection_Close(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should close the database", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectClose()
		err := c.Close()
		assert.NoError(t, err)
	})
}

func Test_PgDBConnection_NamedExecContext(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should execute a named query", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectExec("INSERT INTO users").WithArgs("test").WillReturnResult(sqlmock.NewResult(1, 1))

		_, err := c.NamedExecContext(context.Background(), "INSERT INTO users (name) VALUES (:name)", map[string]interface{}{"name": "test"})
		assert.NoError(t, err)
	})

	t.Run("should return an error", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectExec("INSERT INTO users").WithArgs("test").WillReturnError(errors.New("some error"))

		_, err := c.NamedExecContext(context.Background(), "INSERT INTO users (name) VALUES (:name)", map[string]interface{}{"name": "test"})
		assert.Error(t, err)
	})
}

func Test_PgDBConnection_NamedQueryContext(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	xdb := sqlx.NewDb(mockDB, "sqlmock")

	t.Run("should return rows from a named query", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test")
		mock.ExpectQuery("SELECT id, name FROM users WHERE name = ?").WithArgs("test").WillReturnRows(rows)

		_, err := c.NamedQueryContext(context.Background(), "SELECT id, name FROM users WHERE name = :name", map[string]interface{}{"name": "test"})
		assert.NoError(t, err)
	})

	t.Run("should return an error from a named query", func(t *testing.T) {
		c := &mockPGDBConnection{db: xdb}
		mock.ExpectQuery("SELECT id, name FROM users WHERE name = ?").WithArgs("test").WillReturnError(errors.New("some error"))

		_, err := c.NamedQueryContext(context.Background(), "SELECT id, name FROM users WHERE name = :name", map[string]interface{}{"name": "test"})
		assert.Error(t, err)
	})
}

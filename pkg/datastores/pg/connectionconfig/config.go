package connectionconfig

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/michaeldelorenzo/x/pkg/errors/validationerrors"
	"github.com/pkg/errors"
)

const dbDriverName = "postgres"

var db *sqlx.DB

type SSLMode string

const (
	SSLModeDisable    SSLMode = "disable"
	SSLModeAllow      SSLMode = "allow"
	SSLModePrefer     SSLMode = "prefer"
	SSLModeRequire    SSLMode = "require"
	SSLModeVerifyCA   SSLMode = "verify-ca"
	SSLModeVerifyFull SSLMode = "verify-full"
)

// PostgresConfig holds the Postgres database connection credentials and supports parsing of a JSON object holding the configuration
type PostgresConfig struct {
	UserName string  `json:"username"`
	Password string  `json:"password"`
	Host     string  `json:"host"`
	Port     int     `json:"port"`
	DbName   string  `json:"dbname"`
	SSLMode  SSLMode `json:"sslMode"`

	dbMonitorID string
}

// SetDatastore allows the datastore to be set externally. Useful for testing.
func SetDatastore(database *sqlx.DB) {
	db = database
}

// Datastore returns the connected database (initializes the connection if it's not already connected)
func (pcfg PostgresConfig) Datastore() (*sqlx.DB, error) {
	if db == nil {
		c, err := sqlx.Connect(dbDriverName, pcfg.ConnectionString())
		if err != nil {
			return nil, err
		}

		SetDatastore(c)
	}

	return db, nil
}

// ConnectionString returns a properly formatted Postgres connection string
func (pcfg PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pcfg.UserName,
		url.PathEscape(pcfg.Password),
		pcfg.GetHost(),
		pcfg.GetPort(),
		pcfg.GetDBName(),
		pcfg.SSLMode,
	)
}

// GetHost returns db hostname
func (pcfg PostgresConfig) GetHost() string {
	return pcfg.Host
}

// GetPort returns db port
func (pcfg PostgresConfig) GetPort() string {
	return fmt.Sprintf("%d", pcfg.Port)
}

// GetDBName returns db name
func (pcfg PostgresConfig) GetDBName() string {
	return pcfg.DbName
}

// SetMonitorID sets the instrumenting database monitor ID for the Postgres connection
func (pcfg *PostgresConfig) SetMonitorID(id string) {
	pcfg.dbMonitorID = id
}

// GetMonitorID returns the instrumenting database monitor ID for the Postgres connection
func (pcfg PostgresConfig) GetMonitorID() string {
	return pcfg.dbMonitorID
}

// Validate returns a validation error when PostgresConfig is invalid
func (pcfg *PostgresConfig) Validate() error {
	err := validationerrors.NewValidationError("postgres configuration")

	if pcfg.Host == "" {
		err.AddMessage("host must be present")
	}
	if pcfg.Port == 0 {
		err.AddMessage("port must be present")
	}
	if pcfg.UserName == "" {
		err.AddMessage("user must be present")
	}
	if pcfg.Password == "" {
		err.AddMessage("password must be present")
	}
	if pcfg.DbName == "" {
		err.AddMessage("dbname must be present")
	}

	if err.Present() {
		return err
	}

	return nil
}

// NewPGConfigFromJSON returns a PGConfig from a json string,
// or an error when it fails to parse
func NewPGConfigFromJSON(jsonStr string) (PostgresConfig, error) {
	var pgConfig PostgresConfig

	err := json.Unmarshal([]byte(jsonStr), &pgConfig)
	if err != nil {
		return pgConfig, errors.Wrap(err, "could not parse Postgres config json")
	}

	// Set default SSLMode if not provided
	if pgConfig.SSLMode == "" {
		pgConfig.SSLMode = SSLModeDisable
	}

	return pgConfig, nil
}

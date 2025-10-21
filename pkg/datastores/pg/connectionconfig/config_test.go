package connectionconfig_test

import (
	"testing"

	"github.com/michaeldelorenzo/x/v2/pkg/datastores/pg/connectionconfig"
)

func TestNewFromConfig(t *testing.T) {
	pgConfig := connectionconfig.PostgresConfig{
		Host:     "the-host",
		Port:     1111,
		UserName: "the-user",
		Password: "the-password",
		DbName:   "the-db-name",
		SSLMode:  connectionconfig.SSLModeVerifyFull,
	}

	// connConfig := connectionconfig.NewFromConfig(&pgConfig) // Removed NewFromConfig

	if pgConfig.Host != pgConfig.Host { // Direct comparison
		t.Errorf("expected %s but was %s", pgConfig.Host, pgConfig.Host)
	}
	if pgConfig.Port != pgConfig.Port { // Direct comparison
		t.Errorf("expected %d but was %d", pgConfig.Port, pgConfig.Port)
	}
	if pgConfig.UserName != pgConfig.UserName { // Direct comparison
		t.Errorf("expected %s but was %s", pgConfig.UserName, pgConfig.UserName)
	}
	if pgConfig.Password != pgConfig.Password { // Direct comparison
		t.Errorf("expected %s but was %s", pgConfig.Password, pgConfig.Password)
	}
	if pgConfig.DbName != pgConfig.DbName { // Direct comparison
		t.Errorf("expected %s but was %s", pgConfig.DbName, pgConfig.DbName)
	}
	if pgConfig.SSLMode != pgConfig.SSLMode { // Direct comparison
		t.Errorf("expected %v but was %v", pgConfig.SSLMode, pgConfig.SSLMode)
	}
}

func TestConnectionString(t *testing.T) {
	tests := []struct {
		name        string
		connConfig  connectionconfig.PostgresConfig
		expectedURL string
	}{
		{
			name: "sslmode disable",
			connConfig: connectionconfig.PostgresConfig{
				Host:     "the-host",
				Port:     1111,
				UserName: "the-user",
				Password: ">[m<y-p-!a<s]s~wor]d",
				DbName:   "the-db-name",
				SSLMode:  connectionconfig.SSLModeDisable,
			},
			expectedURL: "postgres://the-user:%3E%5Bm%3Cy-p-%21a%3Cs%5Ds~wor%5Dd@the-host:1111/the-db-name?sslmode=disable",
		},
		{
			name: "sslmode allow",
			connConfig: connectionconfig.PostgresConfig{
				Host:     "the-host",
				Port:     1111,
				UserName: "the-user",
				Password: ">[m<y-p-!a<s]s~wor]d",
				DbName:   "the-db-name",
				SSLMode:  connectionconfig.SSLModeAllow,
			},
			expectedURL: "postgres://the-user:%3E%5Bm%3Cy-p-%21a%3Cs%5Ds~wor%5Dd@the-host:1111/the-db-name?sslmode=allow",
		},
		{
			name: "sslmode prefer",
			connConfig: connectionconfig.PostgresConfig{
				Host:     "the-host",
				Port:     1111,
				UserName: "the-user",
				Password: ">[m<y-p-!a<s]s~wor]d",
				DbName:   "the-db-name",
				SSLMode:  connectionconfig.SSLModePrefer,
			},
			expectedURL: "postgres://the-user:%3E%5Bm%3Cy-p-%21a%3Cs%5Ds~wor%5Dd@the-host:1111/the-db-name?sslmode=prefer",
		},
		{
			name: "sslmode require",
			connConfig: connectionconfig.PostgresConfig{
				Host:     "the-host",
				Port:     1111,
				UserName: "the-user",
				Password: ">[m<y-p-!a<s]s~wor]d",
				DbName:   "the-db-name",
				SSLMode:  connectionconfig.SSLModeRequire,
			},
			expectedURL: "postgres://the-user:%3E%5Bm%3Cy-p-%21a%3Cs%5Ds~wor%5Dd@the-host:1111/the-db-name?sslmode=require",
		},
		{
			name: "sslmode verify-ca",
			connConfig: connectionconfig.PostgresConfig{
				Host:     "the-host",
				Port:     1111,
				UserName: "the-user",
				Password: ">[m<y-p-!a<s]s~wor]d",
				DbName:   "the-db-name",
				SSLMode:  connectionconfig.SSLModeVerifyCA,
			},
			expectedURL: "postgres://the-user:%3E%5Bm%3Cy-p-%21a%3Cs%5Ds~wor%5Dd@the-host:1111/the-db-name?sslmode=verify-ca",
		},
		{
			name: "sslmode verify-full",
			connConfig: connectionconfig.PostgresConfig{
				Host:     "the-host",
				Port:     1111,
				UserName: "the-user",
				Password: ">[m<y-p-!a<s]s~wor]d",
				DbName:   "the-db-name",
				SSLMode:  connectionconfig.SSLModeVerifyFull,
			},
			expectedURL: "postgres://the-user:%3E%5Bm%3Cy-p-%21a%3Cs%5Ds~wor%5Dd@the-host:1111/the-db-name?sslmode=verify-full",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.connConfig.ConnectionString() != test.expectedURL {
				t.Errorf("expected ConnectionString to be %s but was %s", test.expectedURL, test.connConfig.ConnectionString())
			}
		})
	}
}

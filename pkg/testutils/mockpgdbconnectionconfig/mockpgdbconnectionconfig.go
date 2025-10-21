package mockpgdbconnectionconfig

// Mock mocks the db connection config
type Mock struct{}

// ConnectionString returns the postgres connection string
func (m Mock) ConnectionString() string {
	return "the-mock-connection-url"
}

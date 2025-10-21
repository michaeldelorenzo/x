# xapm - Cross-Platform APM Provider

The `xapm` package provides a unified interface for Application Performance Monitoring (APM) that supports multiple providers including New Relic and Sentry.

## Features

- **Provider-agnostic API** - Write code once, switch providers with configuration
- **Multiple provider support**:
  - New Relic
  - Sentry
  - No-op (for testing/development)
- **Comprehensive instrumentation**:
  - HTTP transactions
  - Database queries (including PostgreSQL-specific parsing)
  - External HTTP calls
  - Kafka message production
  - Custom segments
  - Custom events
  - Error tracking

## Installation

```bash
go get github.com/michaeldelorenzo/x/pkg/instrumenting/xapm
```

## Quick Start

### Initialize with New Relic

```go
import (
    "github.com/michaeldelorenzo/x/pkg/instrumenting/xapm"
    "github.com/michaeldelorenzo/x/pkg/instrumenting/xapm/newrelic"
)

func main() {
    err := xapm.Init(&xapm.Config{
        Provider: xapm.ProviderNewRelic,
        NewRelic: &newrelic.Config{
            LicenseKey:      "your-license-key",
            ApplicationName: "my-app",
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer xapm.Apm.Shutdown(10 * time.Second)
}
```

### Initialize with Sentry

```go
import (
    "github.com/michaeldelorenzo/x/pkg/instrumenting/xapm"
    "github.com/michaeldelorenzo/x/pkg/instrumenting/xapm/sentry"
)

func main() {
    err := xapm.Init(&xapm.Config{
        Provider: xapm.ProviderSentry,
        Sentry: &sentry.Config{
            DSN:              "your-sentry-dsn",
            Environment:      "production",
            Release:          "v1.0.0",
            TracesSampleRate: 1.0,
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer xapm.Apm.Shutdown(10 * time.Second)
}
```

### Initialize with No-op (Testing/Development)

```go
err := xapm.Init(&xapm.Config{
    Provider: xapm.ProviderNoop,
})
```

## Usage

### HTTP Middleware

Use the provided middleware to automatically instrument HTTP requests:

```go
import (
    "github.com/michaeldelorenzo/x/pkg/instrumenting/xapm/http/middleware"
)

func main() {
    // Standard library http
    http.Handle("/", middleware.Transaction(http.HandlerFunc(myHandler)))

    // Or with Echo framework
    e := echo.New()
    e.Use(echo.WrapMiddleware(middleware.Transaction))
}
```

### Manual Transactions

```go
func doWork() {
    tx := xapm.StartTransaction("background-job")
    defer tx.End()

    // Add custom attributes
    tx.AddAttribute("user_id", 123)

    // Your code here...

    // Track errors
    if err := someOperation(); err != nil {
        tx.NoticeError(err)
    }
}
```

### Database Instrumentation

#### Generic Database

```go
tx := xapm.TxFromCtx(ctx)
if tx != nil {
    seg := tx.StartDBSegment(&xapm.DBSegParams{
        Host:         "localhost:5432",
        DatabaseName: "mydb",
        Collection:   "users",
        Operation:    xapm.DBSelect,
        DatabaseType: xapm.DBPostgres,
        QueryString:  "SELECT * FROM users WHERE id = $1",
        QueryParams:  map[string]interface{}{"id": userId},
    })
    defer seg.End()
}

// Execute your query...
```

#### PostgreSQL (with query parsing)

```go
seg := tx.StartPSQLSegment(&xapm.PSQLDBSegParams{
    Host:         "localhost:5432",
    DatabaseName: "mydb",
    QueryString:  "SELECT name, email FROM users WHERE id = $1",
    QueryParams:  map[string]interface{}{"id": userId},
})
defer seg.End()
```

### External HTTP Calls

```go
tx := xapm.TxFromCtx(ctx)
if tx != nil {
    seg := tx.StartExternalSegment("https://api.example.com/users")
    defer seg.End()
}

// Make your HTTP request...
```

### Kafka Instrumentation

```go
tx := xapm.TxFromCtx(ctx)
if tx != nil {
    seg := tx.StartKafkaSegment("user-events")
    defer seg.End()
}

// Produce your Kafka message...
```

### Custom Events

Define an event type:

```go
type UserSignupEvent struct {
    UserID   int
    Email    string
    Plan     string
}

func (e *UserSignupEvent) GetEventType() string {
    return "UserSignup"
}
```

Send the event:

```go
xapm.SendCustomEvent(&UserSignupEvent{
    UserID: 123,
    Email:  "user@example.com",
    Plan:   "premium",
})
```

### Context-Based Transactions

Embed a transaction in a context:

```go
tx := xapm.StartTransaction("my-operation")
defer tx.End()

ctx := xapm.CtxFromTx(context.Background(), tx)

// Pass ctx to other functions
doSomething(ctx)
```

Retrieve a transaction from context:

```go
func doSomething(ctx context.Context) {
    tx := xapm.TxFromCtx(ctx)
    if tx != nil {
        seg := tx.StartSegment("sub-operation")
        defer seg.End()

        // Your code...
    }
}
```

## Provider Comparison

| Feature | New Relic | Sentry | No-op |
|---------|-----------|--------|-------|
| Transactions | ✅ | ✅ | ✅ (no-op) |
| Custom Events | ✅ | ✅ | ✅ (no-op) |
| Error Tracking | ✅ | ✅ | ✅ (no-op) |
| DB Instrumentation | ✅ | ✅ | ✅ (no-op) |
| HTTP Instrumentation | ✅ | ✅ | ✅ (no-op) |
| Distributed Tracing | ✅ | ✅ | ❌ |
| Query Parsing | ✅ (PostgreSQL) | ❌ | ❌ |

## Configuration Reference

### New Relic Config

```go
type Config struct {
    LicenseKey        string                 // Required: Your New Relic license key
    ApplicationName   string                 // Required: Application name in New Relic
    TrustedAccountKey string                 // Optional: For distributed tracing
    AdditionalConfig  newrelic.ConfigOption  // Optional: Additional New Relic options
}
```

### Sentry Config

```go
type Config struct {
    DSN              string   // Required: Sentry DSN
    Environment      string   // Optional: Environment name (e.g., "production")
    Release          string   // Optional: Release version
    TracesSampleRate float64  // Required: Sample rate for traces (0.0 to 1.0)
    Debug            bool     // Optional: Enable debug mode
}
```

## Best Practices

1. **Always defer transaction End()** - Ensures proper cleanup and data submission
2. **Use context propagation** - Pass transactions through context for automatic instrumentation
3. **Sanitize sensitive data** - Be careful not to log passwords or tokens in query parameters
4. **Use appropriate sample rates** - In production, consider lowering Sentry's `TracesSampleRate` to reduce volume
5. **Test with no-op provider** - Use the no-op provider in tests to avoid APM overhead

## Examples

See the `xapm_test.go` file for comprehensive examples of all features.

## Switching Providers

To switch from New Relic to Sentry (or vice versa), simply change the configuration:

```go
// Before (New Relic)
xapm.Init(&xapm.Config{
    Provider: xapm.ProviderNewRelic,
    NewRelic: &newrelic.Config{...},
})

// After (Sentry)
xapm.Init(&xapm.Config{
    Provider: xapm.ProviderSentry,
    Sentry: &sentry.Config{...},
})
```

No code changes required! The abstraction layer handles the differences between providers.

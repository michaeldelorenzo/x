# Instrumenting

Contains common interfaces for instrumentation, e.g. APM and logging, for the Go programming language.

## Example Usage

The following is some example usage. This code should be runnable copy/paste standalone Go application:

```go
package main

import (
	"github.com/michaeldelorenzo/x/pkg/instrumenting/xapm"
	"github.com/michaeldelorenzo/x/pkg/instrumenting/xlog"
	"github.com/labstack/echo/v4"

	"context"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Init Logging
	xlog.Init(&xlog.Config{
		AppName: "GolangTest",
		Level:   "info",
		GitHash: "GIT_HASH",
		Version: "2.1",
		// We will push our logs to New Relic as well
		LogUpstream: true,
		// We can also set IsAsync to true to make our NR
		// calls in the background. This is false by default.
		// IsAsync: true,
		NRClient: &xlog.NewRelicClient{
			LicenseKey: "<REDACTED>",
		},
	}).Sync()

	// Init APM
	defer xapm.Init(&xapm.Config{
		ApplicationName: "GolangTest",
		LicenseKey:      "<REDACTED>",
	},
		xlog.L().Desugar(),
	).Shutdown(5 * time.Second)

	// Calling L() returns our global logger.
	// Most With() calls will start from xlog.L().With()
	xlog.L().Info("Starting our Server...")

	//: Start routing
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		txn := xapm.StartTransaction("HandlingNewRequest")
		defer txn.End()

		makeStandardSegment(xapm.CtxFromTx(c.Request().Context(), txn))

		time.Sleep(3 * time.Second)
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/err", func(c echo.Context) error {
		txn := xapm.StartTransaction("HandlingNewErrorRequest")
		defer txn.End()

		xlog.EE("error", "uh-oh....another error came through :(", txn)

		time.Sleep(3 * time.Second)
		return c.String(http.StatusOK, "Uh Oh")
	})

	e.GET("/ext", func(c echo.Context) error {
		txn := xapm.StartTransaction("MakingExternalCall")
		defer txn.End()

		makeExternalSegment(xapm.CtxFromTx(c.Request().Context(), txn))

		time.Sleep(3 * time.Second)
		return c.String(http.StatusOK, "Made external call!")
	})

	e.GET("/db", func(c echo.Context) error {
		txn := xapm.StartTransaction("MakingDBCall")
		defer txn.End()

		makeDBSegment(xapm.CtxFromTx(c.Request().Context(), txn))

		time.Sleep(3 * time.Second)
		return c.String(http.StatusOK, "Making request to the DB!")
	})

	e.GET("/kafka", func(c echo.Context) error {
		txn := xapm.StartTransaction("ProducingKafkaMessage")
		defer txn.End()

		makeKafkaSegment(xapm.CtxFromTx(c.Request().Context(), txn))

		time.Sleep(3 * time.Second)
		return c.String(http.StatusOK, "Producing Message To Kafka!")
	})

	e.GET("/logging", func(c echo.Context) error {
		txn := xapm.StartTransaction("NestedContextualLogging")
		defer txn.End()

		complexFn(xapm.CtxFromTx(c.Request().Context(), txn))

		time.Sleep(3 * time.Second)
		return c.String(http.StatusOK, "Running some complex, multi-call function!")
	})

	//: End Routing

	e.Logger.Fatal(e.Start(":1323"))

	return

}

func makeStandardSegment(ctx context.Context) {
	// Fetch our transaction from context
	txn := xapm.TxFromCtx(ctx)
	defer txn.StartSegment("makeStandardSegment/segment").End()

	time.Sleep(2 * time.Second)
}

func makeExternalSegment(ctx context.Context) {
	// Fetch our transaction from context
	txn := xapm.TxFromCtx(ctx)
	defer xapm.StartExternalSegment(txn, "Measurements Service").End()

	time.Sleep(2 * time.Second)
}

func complexFn(ctx context.Context) {
	// Call convenience method LL() to embed and return both our
	// context (with embedded logger) and the logger itself, ready for use.
	newCtx, logger := xlog.LL(ctx, xlog.Fields{
		"someComplexValue": "Hello",
		"someMoreContext":  "World!",
	})

	// We can call zap's native functionality and if our core is
	//configured for upstream, it will still perform that upstream push.
	logger.Info("Performing some additional complex task")

	time.Sleep(3 * time.Second)

	complexFn2(newCtx)
}

func complexFn2(ctx context.Context) {
	logger := xlog.LoggerFromCtx(ctx)

	// We can call zap's native functionality and if our core is
	//configured for upstream, it will still perform that upstream push.
	logger.Info("More Complex Work...Lets Add More Context To Logger..")
	logger = logger.With(xlog.Fields{
		"Even More Context": "Foo",
		"More Data":         "Bar",
	})

	time.Sleep(3 * time.Second)
}

func makeKafkaSegment(ctx context.Context) {
	// Fetch our transaction from context
	txn := xapm.TxFromCtx(ctx)
	defer xapm.StartKafkaSegment(txn, "measurements_topic").End()

	// Simulate a random error producing to kafka
	if (rand.Intn(10-1)+1)%2 == 0 {
		// EEWith logs and reports the error data to
		//newrelic, if a valid transaction is provided.
		xlog.EEWith(
			"unable to produce to topic",
			xlog.Fields{
				"topic": "measurements_topic",
				"event": "TEST_EVENT",
			},
			txn,
		)
	}

	time.Sleep(3 * time.Second)
}

func makeDBSegment(ctx context.Context) {
	// Fetch our transaction from context
	txn := xapm.TxFromCtx(ctx)
	defer xapm.StartDBSegment(txn, &xapm.DBSegParams{
		Host:         "http://rds.aws.com/accountname",
		DatabaseName: "data_db",
		DatabaseType: xapm.DBPostgres,
		Collection:   "measurements",
		QueryString:  "SELECT * FROM measurements WHERE 1=1;",
	}).End()

	time.Sleep(2 * time.Second)
}
```

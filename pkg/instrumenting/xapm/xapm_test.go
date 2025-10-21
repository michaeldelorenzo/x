package xapm_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michaeldelorenzo/x/v2/pkg/instrumenting/xapm"
	"github.com/michaeldelorenzo/x/v2/pkg/instrumenting/xapm/noop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoopProvider(t *testing.T) {
	provider := noop.NewProvider()
	xapm.SetProvider(provider)

	t.Run("StartTransaction", func(t *testing.T) {
		tx := xapm.StartTransaction("test-transaction")
		require.NotNil(t, tx)
		defer tx.End()

		// No-op provider should not panic
		tx.AddAttribute("test", "value")
		tx.NoticeError(assert.AnError)
	})

	t.Run("SendCustomEvent", func(t *testing.T) {
		event := &testEvent{
			Name:  "TestEvent",
			Count: 42,
		}

		// Should not panic
		xapm.SendCustomEvent(event)
	})

	t.Run("ContextOperations", func(t *testing.T) {
		ctx := context.Background()
		tx := xapm.StartTransaction("test-transaction")
		defer tx.End()

		// Embed transaction in context
		ctx = xapm.CtxFromTx(ctx, tx)

		// Retrieve transaction from context
		retrievedTx := xapm.TxFromCtx(ctx)
		assert.NotNil(t, retrievedTx)
	})

	t.Run("Segments", func(t *testing.T) {
		tx := xapm.StartTransaction("test-transaction")
		defer tx.End()

		seg := tx.StartSegment("test-segment")
		require.NotNil(t, seg)
		seg.AddAttribute("key", "value")
		seg.End()

		extSeg := tx.StartExternalSegment("https://example.com")
		require.NotNil(t, extSeg)
		extSeg.End()

		dbSeg := tx.StartDBSegment(&xapm.DBSegParams{
			Host:         "localhost",
			DatabaseName: "testdb",
			Collection:   "users",
			Operation:    xapm.DBSelect,
			DatabaseType: xapm.DBPostgres,
			QueryString:  "SELECT * FROM users WHERE id = $1",
			QueryParams:  map[string]interface{}{"id": 123},
		})
		require.NotNil(t, dbSeg)
		dbSeg.End()
	})

	t.Run("WebTransaction", func(t *testing.T) {
		tx := xapm.StartTransaction("test-web")
		defer tx.End()

		req := httptest.NewRequest("GET", "/test", nil)
		tx.SetWebRequestHTTP(req)

		w := httptest.NewRecorder()
		wrappedWriter := tx.SetWebResponse(w)
		require.NotNil(t, wrappedWriter)

		// Should not panic
		wrappedWriter.WriteHeader(http.StatusOK)
		_, err := wrappedWriter.Write([]byte("test"))
		assert.NoError(t, err)
	})
}

func TestInit(t *testing.T) {
	t.Run("NoopProvider", func(t *testing.T) {
		conf := &xapm.Config{
			Provider: xapm.ProviderNoop,
		}

		err := xapm.Init(conf)
		require.NoError(t, err)

		provider := xapm.GetProvider()
		require.NotNil(t, provider)
		assert.Equal(t, xapm.ProviderNoop, provider.Type())
	})

	t.Run("MissingNewRelicConfig", func(t *testing.T) {
		conf := &xapm.Config{
			Provider: xapm.ProviderNewRelic,
			// NewRelic config is nil
		}

		err := xapm.Init(conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "NewRelic configuration is required")
	})

	t.Run("MissingSentryConfig", func(t *testing.T) {
		conf := &xapm.Config{
			Provider: xapm.ProviderSentry,
			// Sentry config is nil
		}

		err := xapm.Init(conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Sentry configuration is required")
	})

	t.Run("UnsupportedProvider", func(t *testing.T) {
		conf := &xapm.Config{
			Provider: "unsupported",
		}

		err := xapm.Init(conf)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported APM provider")
	})
}

// testEvent is a test implementation of EventReporter
type testEvent struct {
	Name  string
	Count int
}

func (e *testEvent) GetEventType() string {
	return "TestEvent"
}

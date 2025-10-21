package middleware

import (
	"net/http"

	"github.com/michaeldelorenzo/x/v3/pkg/instrumenting/xapm"
)

// Transaction is a HTTP middleware function which creates a new APM transaction for each request.
//
//	e := echo.New()
//	e.Use(echo.WrapMiddleware(middleware.Transaction))
func Transaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + " " + r.URL.Path

		tx := xapm.StartTransaction(name)
		defer tx.End()

		// Set web request details
		tx.SetWebRequestHTTP(r)

		// Add the transaction into context for both the xapm package and provider-specific packages to use
		ctx := xapm.CtxFromTx(r.Context(), tx)

		// Wrap the response writer to capture response details
		writer := tx.SetWebResponse(w)

		// Create new request with updated context
		req := r.WithContext(ctx)

		next.ServeHTTP(writer, req)
	})
}

package middleware

import (
	"net/http"

	"github.com/michaeldelorenzo/x/pkg/instrumenting/xapm"
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

		// add the txn into context for both the xapm and newrelic packages to use
		ctx := xapm.CtxFromTx(r.Context(), tx)

		req := r.WithContext(ctx)
		tx.SetWebRequestHTTP(req)
		writer := tx.SetWebResponse(w)

		next.ServeHTTP(writer, req)
	})
}

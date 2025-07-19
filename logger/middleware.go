// gowebcore/logger/middleware.go
package logger

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// ChiLogger logs method, path, status, bytes, and request-ID for every request.
func ChiLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		// Use the new helper introduced in logger.go
		Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"size", ww.BytesWritten(),
			"request_id", middleware.GetReqID(r.Context()),
		)
	})
}

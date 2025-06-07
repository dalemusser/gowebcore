package auth

import (
	"context"
	"crypto/subtle"
	"net/http"
)

type apiKeyCtx struct{}

var keyCtx apiKeyCtx = struct{}{}

// APIKeyMiddleware compares header X-API-Key with expectedKey (constant-time).
func APIKeyMiddleware(expectedKey string) func(http.Handler) http.Handler {
	ek := []byte(expectedKey)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := []byte(r.Header.Get("X-API-Key"))
			if len(got) == 0 || subtle.ConstantTimeCompare(got, ek) != 1 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), keyCtx, string(got))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// APIKeyFromContext returns the validated key if present.
func APIKeyFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyCtx).(string)
	return v, ok
}

package auth

import (
	"net/http"          // ← add
	"net/http/httptest"
	"testing"
)

func TestAPIKeyMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "secret123")

	handler := APIKeyMiddleware("secret123")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if k, ok := APIKeyFromContext(r.Context()); !ok || k != "secret123" {
			t.Fatalf("key missing in ctx")
		}
	}))

	handler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("unexpected status %d", rec.Code)
	}
}

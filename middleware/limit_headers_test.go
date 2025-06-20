// middleware/limit_headers_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestRateLimitHeaders(t *testing.T) {
	mw := RateLimit(1, time.Second, 1) // 1 req / s
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "9.8.7.6:12345"

	// first request – should pass
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if got := w.Header().Get("X-RateLimit-Limit"); got != "1" {
		t.Fatalf("want Limit=1, got %s", got)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}

	// second immediate request – should be 429 with Retry-After
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("want 429, got %d", w.Code)
	}
	ra := w.Header().Get("Retry-After")
	if ra == "" {
		t.Fatal("Retry-After missing")
	}
	if sec, _ := strconv.Atoi(ra); sec < 1 {
		t.Fatalf("Retry-After too small: %s", ra)
	}
}

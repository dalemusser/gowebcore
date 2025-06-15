// middleware/csrf_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinas/nosurf"
)

func TestCSRFTokenValidation(t *testing.T) {
	var token string // will be filled inside first request

	// Protect a handler that grabs the token from the request context
	protected := CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token == "" { // first (GET) request
			token = nosurf.Token(r) // guaranteed non-empty after middleware
		}
		w.WriteHeader(http.StatusOK)
	}))

	// 1) GET ─ obtain Secure cookie + token
	getReq := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
	getRes := httptest.NewRecorder()
	protected.ServeHTTP(getRes, getReq)

	var csrfCookie *http.Cookie
	for _, c := range getRes.Result().Cookies() {
		if c.Name == "csrf_token" {
			csrfCookie = c
			break
		}
	}
	if csrfCookie == nil {
		t.Fatal("csrf cookie missing")
	}
	if token == "" {
		t.Fatal("nosurf token missing")
	}

	// 2) POST ─ send cookie, Referer, and the captured token
	postReq := httptest.NewRequest(http.MethodPost, "https://example.com/", nil)
	postReq.AddCookie(csrfCookie)
	postReq.Header.Set("Referer", "https://example.com/")
	postReq.Header.Set("X-CSRF-Token", token)

	postRes := httptest.NewRecorder()
	protected.ServeHTTP(postRes, postReq)

	if postRes.Code != http.StatusOK {
		t.Fatalf("csrf validation failed, got %d", postRes.Code)
	}
}

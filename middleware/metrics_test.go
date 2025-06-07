package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateCSRF(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h := CSRFCookie(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	h.ServeHTTP(w, r)
	c, _ := r.Cookie(csrfCookie)
	r.Header.Set("X-CSRF-Token", c.Value)
	if !ValidateCSRF(r) {
		t.Fatal("csrf validation failed")
	}
}

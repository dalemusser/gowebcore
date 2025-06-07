package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

const csrfCookie = "csrf_token"

func CSRFCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(csrfCookie)
		if err != nil || token.Value == "" {
			b := make([]byte, 32)
			_, _ = rand.Read(b)
			val := base64.RawURLEncoding.EncodeToString(b)
			http.SetCookie(w, &http.Cookie{
				Name:  csrfCookie,
				Value: val,
				Path:  "/",
			})
			r.AddCookie(&http.Cookie{Name: csrfCookie, Value: val})
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateCSRF checks that header "X-CSRF-Token" equals the cookie.
// Call in your POST/PUT/PATCH handlers.
func ValidateCSRF(r *http.Request) bool {
	c, err := r.Cookie(csrfCookie)
	if err != nil {
		return false
	}
	return r.Header.Get("X-CSRF-Token") == c.Value
}

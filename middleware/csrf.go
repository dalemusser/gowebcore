package middleware

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// CSRF wraps the provided handler with nosurf CSRF protection.
//   - A cookie named "csrf_token" is set (SameSite=Lax, Secure when TLS).
//   - For state-changing requests (POST/PUT/PATCH/DELETE) the token must be
//     supplied in a header "X-CSRF-Token" or form field "csrf_token".
//   - Call nosurf.Token(r) in templates to embed the token.
//
// Example:
//
//	r := chi.NewRouter()
//	r.Use(middleware.CSRF)               // global
//	r.Get("/form", showForm)
//	r.Post("/submit", handleSubmit)
func CSRF(next http.Handler) http.Handler {
	csrf := nosurf.New(next)

	// Cookie settings
	csrf.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true, // automatically downgraded by nosurf on non-TLS
	})

	return csrf
}

// Token returns the request's CSRF token (helper for templates).
func Token(r *http.Request) string {
	return nosurf.Token(r)
}

package middleware

import (
	"context"
	"net/http"

	"github.com/dalemusser/gowebcore/auth"
	"github.com/dalemusser/gowebcore/auth/oauth"
	"github.com/go-chi/chi/v5"
)

type ctxKey struct{}

func RequireAuth(sess *auth.Session) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var claims map[string]any
			if err := sess.Get(r, "session", &claims); err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func User(r *http.Request) map[string]any {
	if v, ok := r.Context().Value(ctxKey{}).(map[string]any); ok {
		return v
	}
	return nil
}

func Routes(r chi.Router, provider *oauth.Provider, sess *auth.Session) {
	r.Get("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, provider.AuthURL("state123"), http.StatusFound)
	})

	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := provider.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "exchange failed", http.StatusBadRequest)
			return
		}
		info, err := provider.UserInfo(r.Context(), token)
		if err != nil {
			http.Error(w, "userinfo failed", http.StatusBadRequest)
			return
		}
		_ = sess.Set(w, "session", info)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	r.Get("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		sess.Clear(w, "session")
		http.Redirect(w, r, "/", http.StatusFound)
	})
}

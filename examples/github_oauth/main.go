package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dalemusser/gowebcore/auth"
	"github.com/dalemusser/gowebcore/auth/oauth"
	"github.com/dalemusser/gowebcore/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	provider := oauth.NewGitHub(
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		"http://localhost:8080/auth/callback",
		// The helper already requests "read:user user:email".
		// If you need custom scopes, extend the helper to accept a variadic slice.
	)

	// 64-byte hash key – replace in production
	hashKey := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789ABCDEF")
	sess := auth.NewSession(hashKey, nil)

	r := chi.NewRouter()
	middleware.Routes(r, provider, sess) // /auth/login, /auth/callback, /auth/logout

	// Protected dashboard
	r.With(middleware.RequireAuth(sess)).Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := middleware.User(r)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, " + user["login"].(string) + "\n"))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

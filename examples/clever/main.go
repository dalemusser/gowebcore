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
	provider := oauth.NewClever(
		os.Getenv("CLEVER_CLIENT_ID"),
		os.Getenv("CLEVER_CLIENT_SECRET"),
		"http://localhost:8080/auth/callback",
	)

	// 64-byte hash key — replace with secure random string in production
	hashKey := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789ABCDEF")
	sess := auth.NewSession(hashKey, nil)

	r := chi.NewRouter()

	// Helpers that already exist in middleware/oauth_session.go
	middleware.Routes(r, provider, sess) // /auth/login|callback|logout

	r.With(middleware.RequireAuth(sess)).Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := middleware.User(r) // map[string]any
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, " + user["id"].(string) + "\n"))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

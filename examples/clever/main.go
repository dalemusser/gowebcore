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
	// Pull credentials from env vars
	provider := oauth.NewClever(
		os.Getenv("CLEVER_CLIENT_ID"),
		os.Getenv("CLEVER_CLIENT_SECRET"),
		"http://localhost:8080/auth/callback",
	)

	// 64-byte hash, 32-byte block key for AES; use securecookie.GenerateRandomKey once
	sess := auth.NewSession(
		[]byte("a-64-byte-hash-key-………………….……………………………………………………"), // replace
		nil, // blockKey optional (nil = no encryption, only HMAC-sign)
	)

	r := chi.NewRouter()
	middleware.Routes(r, provider, sess) // /auth/login, /auth/callback, /auth/logout

	// Protected dashboard
	r.With(middleware.RequireAuth(sess)).Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := middleware.User(r)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, " + user["id"].(string) + "\n"))
	})

	log.Printf("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

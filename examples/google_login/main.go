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
	provider := oauth.NewGoogle(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		"http://localhost:8080/auth/callback",
	)

	// 64-byte hash key (replace in prod!)
	hashKey := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789ABCDEF")
	sess := auth.NewSession(hashKey, nil)

	r := chi.NewRouter()
	middleware.Routes(r, provider, sess) // /auth/login|callback|logout

	r.With(middleware.RequireAuth(sess)).Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := middleware.User(r)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hi, " + user["email"].(string) + "\n"))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

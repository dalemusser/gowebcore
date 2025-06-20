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

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s not set", k)
	}
	return v
}

func main() {
	provider := oauth.NewGoogle(
		mustEnv("GOOGLE_CLIENT_ID"),
		mustEnv("GOOGLE_CLIENT_SECRET"),
		"http://localhost:8080/auth/callback",
	)

	// session key from env
	hashKey := []byte(mustEnv("SESSION_HASH_KEY"))
	sess := auth.NewSession(hashKey, nil)

	r := chi.NewRouter()
	middleware.Routes(r, provider, sess)

	r.With(middleware.RequireAuth(sess)).Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := middleware.User(r)
		w.Write([]byte("Hi, " + user["email"].(string) + "\n"))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

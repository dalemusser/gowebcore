package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/dalemusser/gowebcore/auth"
	"github.com/dalemusser/gowebcore/auth/oauth"
	"github.com/dalemusser/gowebcore/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	provider, err := oauth.NewClassLinkOIDC(
		ctx,
		os.Getenv("CL_SUBDOMAIN"), // e.g. "mydistrict"
		os.Getenv("CL_OIDC_CLIENT_ID"),
		os.Getenv("CL_OIDC_CLIENT_SECRET"),
		"http://localhost:8080/auth/callback",
	)
	if err != nil {
		log.Fatal(err)
	}

	// 64-byte hash key (replace with secure random string in production)
	hashKey := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789ABCDEF")
	sess := auth.NewSession(hashKey, nil) // no AES block key

	r := chi.NewRouter()
	middleware.Routes(r, provider, sess) // /auth/login, /auth/callback, /auth/logout

	// Require login for home page
	r.With(middleware.RequireAuth(sess)).Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := middleware.User(r) // map[string]any from ID-token claims
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, " + user["email"].(string) + "\n"))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

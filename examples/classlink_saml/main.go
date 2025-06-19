package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	samlhelper "github.com/dalemusser/gowebcore/auth/samlsp"
	"github.com/dalemusser/gowebcore/middleware"
)

func main() {
	pubURL, _ := url.Parse("http://localhost:8080")

	sp, err := samlhelper.NewClassLinkMiddleware(
		os.Getenv("CL_SUBDOMAIN"), pubURL, nil, nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	middleware.SAMLRoutes(mux, sp) // /saml/* + /auth/login|logout

	mux.Handle("/dashboard",
		sp.RequireAccount(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You are logged in via ClassLink SAML\n"))
		})),
	)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

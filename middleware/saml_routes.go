package middleware

import (
	"net/http"

	"github.com/crewjam/saml/samlsp"
)

// SAMLRoutes registers ClassLink SAML endpoints & session helpers.
//
//	mux – can be chi.Router *or* http.ServeMux (anything that implements
//	      Handle/HandleFunc).  Replace http.ServeMux if you prefer chi.
func SAMLRoutes(mux *http.ServeMux, sp *samlsp.Middleware) {
	// 1) IdP metadata + ACS are auto-served under /saml/
	mux.Handle("/saml/", sp)

	// 2) /auth/login  – start SAML flow
	mux.HandleFunc("/auth/login", sp.HandleStartAuthFlow)

	// 3) /auth/logout – clear session
	mux.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		sp.Session.DeleteSession(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	})
}

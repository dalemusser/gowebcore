package middleware

import (
	"fmt"
	"net/http"

	"github.com/dalemusser/gowebcore/config"
	"github.com/go-chi/cors"
)

// CORSFromConfig builds a cors.Handler using cfg.CORSOrigins.
// If none are provided, it whitelists the server's own domain.
func CORSFromConfig(cfg config.Base) func(http.Handler) http.Handler {
	origins := cfg.CORSOrigins
	if len(origins) == 0 && cfg.Domain != "" {
		origins = []string{
			fmt.Sprintf("https://%s", cfg.Domain),
			fmt.Sprintf("http://%s", cfg.Domain),
		}
	}
	if len(origins) == 0 { // no domain either → fallback to "*"
		origins = []string{"*"}
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// DefaultCORS returns a ready-to-use Chi middleware.
// Example: r.Use(middleware.DefaultCORS())
func DefaultCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	})
}

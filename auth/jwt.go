package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub   string `json:"sub"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// key type to avoid ctx collisions
type ctxKey struct{}

var claimsKey ctxKey = struct{}{}

// ParseAndVerify parses token string with either:
//   * HMAC secret (if rsaPub is nil)
//   * RSA-public key (if provided)
func ParseAndVerify(token string, hmacSecret []byte, rsaPub *rsa.PublicKey) (*Claims, error) {
	claims := &Claims{}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256", "RS256"}))

	_, err := parser.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		switch t.Method.Alg() {
		case "HS256":
			return hmacSecret, nil
		case "RS256":
			return rsaPub, nil
		default:
			return nil, errors.New("unexpected signing method")
		}
	})
	return claims, err
}

// JWTMiddleware validates Authorization header and stores claims in ctx.
func JWTMiddleware(hmacSecret []byte, rsaPub *rsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			prefix := "Bearer "
			if !strings.HasPrefix(authz, prefix) {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			c, err := ParseAndVerify(strings.TrimPrefix(authz, prefix), hmacSecret, rsaPub)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey, c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// FromContext retrieves Claims if present.
func FromContext(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsKey).(*Claims)
	return c, ok
}

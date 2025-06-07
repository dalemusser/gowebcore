package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"          // ← add
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestParseAndVerify_HS256(t *testing.T) {
	secret := []byte("supersecret")
	claims := &Claims{Sub: "123"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString(secret)

	got, err := ParseAndVerify(s, secret, nil)
	if err != nil || got.Sub != "123" {
		t.Fatalf("failed HS256 verify: %v", err)
	}
}

func TestJWTMiddleware(t *testing.T) {
	secret := []byte("k")
	claims := &Claims{Sub: "abc", RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour))}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString(secret)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+s)

	handler := JWTMiddleware(secret, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := FromContext(r.Context()); !ok || c.Sub != "abc" {
			t.Fatalf("claims missing in ctx")
		}
	}))

	handler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("unexpected status %d", rec.Code)
	}
}

func TestParseAndVerify_RS256(t *testing.T) {
	priv, _ := rsa.GenerateKey(rand.Reader, 1024)
	claims := &Claims{Sub: "xyz"}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	s, _ := token.SignedString(priv)

	got, err := ParseAndVerify(s, nil, &priv.PublicKey)
	if err != nil || got.Sub != "xyz" {
		t.Fatalf("failed RS256 verify: %v", err)
	}
}
